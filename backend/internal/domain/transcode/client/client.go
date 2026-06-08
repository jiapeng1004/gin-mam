package client

import "context"

//go:generate go run go.uber.org/mock/mockgen@v0.6.0 -source=client.go -destination=../mock/client_mock.go -package=mock
type Client interface {
	SubmitJob(ctx context.Context, endpoint string, req SubmitJobRequest) (*SubmitJobResponse, error)
}

type SubmitJobRequest struct {
	TaskID          string `json:"taskId"`
	CallbackURL     string `json:"callbackUrl"`
	AssetID         string `json:"assetId"`
	AssetFileID     string `json:"assetFileId,omitempty"`
	TranscodeGroupID string `json:"transcodeGroupId"`
	GroupParam      string `json:"groupParam,omitempty"`
}

type SubmitJobResponse struct {
	JobID string `json:"jobId"`
}