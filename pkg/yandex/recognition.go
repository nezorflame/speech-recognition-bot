package yandex

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"path/filepath"

	"github.com/pkg/errors"
	stt "github.com/yandex-cloud/go-genproto/yandex/cloud/ai/stt/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// STTAPIEndpoint is used for all voice recognition requests
const STTAPIEndpoint = "stt.api.cloud.yandex.net:443"

// RecognitionClient wraps SttService_StreamingRecognizeClient
type RecognitionClient struct {
	conn     *grpc.ClientConn
	stt      stt.SttService_StreamingRecognizeClient
	iamToken string
	folder   string
}

// NewRecognitionClient creates new RecognitionClient
func (sdk *SDK) NewRecognitionClient(ctx context.Context) (*RecognitionClient, error) {
	iamToken, err := sdk.IAMToken(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create IAM token")
	}

	conn, err := grpc.DialContext(ctx, STTAPIEndpoint,
		grpc.WithTransportCredentials(credentials.NewTLS(nil)),
		grpc.WithPerRPCCredentials(tokenAuth{token: iamToken}),
	)
	if err != nil {
		return nil, errors.Wrap(err, "unable to establish gRPC connection")
	}

	sttClient, err := stt.NewSttServiceClient(conn).StreamingRecognize(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create STT service client")
	}

	return &RecognitionClient{conn: conn, stt: sttClient, folder: sdk.folder, iamToken: iamToken}, nil
}

// Close closes the RecognitionClient gRPC connection
func (rc *RecognitionClient) Close() error {
	return rc.conn.Close()
}

// SimpleRecognize sends an audiofile to recognize it through Yandex SpeechKit with default parameters
func (rc *RecognitionClient) SimpleRecognize(filePath string) (string, error) {
	confReq := rc.NewConfigRequest()
	if err := rc.Send(confReq); err != nil && err != io.EOF {
		log.Fatalf("Unable to send config request: %v", err)
	}

	contentReq, err := rc.NewAudioRequest(filePath)
	if err != nil {
		log.Fatalf("Unable to create the audio request: %v", err)
	}
	if err = rc.Send(contentReq); err != nil && err != io.EOF {
		log.Fatalf("Unable to send audio request: %v", err)
	}

	if err = rc.CloseSend(); err != nil {
		log.Fatalf("Unable to close recognition sending: %v", err)
	}

	sttResp, err := rc.Recv()
	if err != nil {
		return "", errors.Wrap(err, "unable to receive recognition response")
	}

	for _, c := range sttResp.GetChunks() {
		if c.GetFinal() {
			alt := c.GetAlternatives()
			return alt[0].GetText(), nil
		}
	}

	return "", errors.New("no final result was found")
}

// NewConfigRequest returns a properly set StreamingRecognitionRequest for config
func (rc *RecognitionClient) NewConfigRequest() *stt.StreamingRecognitionRequest {
	return &stt.StreamingRecognitionRequest{StreamingRequest: &stt.StreamingRecognitionRequest_Config{
		Config: &stt.RecognitionConfig{
			Specification: &stt.RecognitionSpec{
				AudioEncoding:  stt.RecognitionSpec_OGG_OPUS,
				LanguageCode:   "ru-RU",
				PartialResults: false,
			},
			FolderId: rc.folder,
		}},
	}
}

// NewAudioRequest returns a properly set StreamingRecognitionRequest for audiofile
func (rc *RecognitionClient) NewAudioRequest(audioFilePath string) (*stt.StreamingRecognitionRequest, error) {
	audioFile, err := ioutil.ReadFile(filepath.Clean(audioFilePath))
	if err != nil {
		return nil, fmt.Errorf("unable to open the audio file: %v", err)
	}

	return &stt.StreamingRecognitionRequest{
		StreamingRequest: &stt.StreamingRecognitionRequest_AudioContent{AudioContent: audioFile},
	}, nil
}

// Send wraps SttService_StreamingRecognizeClient.Send
func (rc *RecognitionClient) Send(req *stt.StreamingRecognitionRequest) error {
	return rc.stt.Send(req)
}

// CloseSend wraps SttService_StreamingRecognizeClient.CloseSend
func (rc *RecognitionClient) CloseSend() error {
	return rc.stt.CloseSend()
}

// Recv wraps SttService_StreamingRecognizeClient.Recv
func (rc *RecognitionClient) Recv() (*stt.StreamingRecognitionResponse, error) {
	return rc.stt.Recv()
}
