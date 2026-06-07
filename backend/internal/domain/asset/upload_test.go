package asset_test

import (
	"bytes"
	"context"
	"io"
	"testing"

	"github.com/gin-mam/backend/internal/domain/asset"
	"github.com/gin-mam/backend/internal/domain/asset/mock"
	infraconfig "github.com/gin-mam/backend/internal/infra/config"
	configmock "github.com/gin-mam/backend/internal/infra/config/mock"
	storagemock "github.com/gin-mam/backend/internal/infra/storage/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestUploadService_Init_UsesDefaultChunkSize(t *testing.T) {
	ctrl := gomock.NewController(t)
	store := mock.NewMockUploadSessionStore(ctrl)
	cfg := configmock.NewMockService(ctrl)
	stor := storagemock.NewMockStorage(ctrl)
	assets := mock.NewMockService(ctrl)

	cfg.EXPECT().Get(gomock.Any(), asset.ConfigKeyUploadChunkSize).Return("", infraconfig.ErrNotFound)
	store.EXPECT().SaveMeta(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

	svc := asset.NewUploadService(store, cfg, stor, assets)
	res, err := svc.Init(context.Background(), asset.UploadInitInput{
		FileName: "demo.mp4",
		FileSize: 10,
		MimeType: "video/mp4",
	})
	require.NoError(t, err)
	require.NotEmpty(t, res.UploadID)
	require.Equal(t, int64(asset.DefaultUploadChunkSize), res.ChunkSize)
}

func TestUploadService_Complete_MergesChunksAndCreatesAsset(t *testing.T) {
	ctrl := gomock.NewController(t)
	store := mock.NewMockUploadSessionStore(ctrl)
	cfg := configmock.NewMockService(ctrl)
	stor := storagemock.NewMockStorage(ctrl)
	assets := mock.NewMockService(ctrl)

	uploadID := "up001"
	payload := []byte("hello")
	meta := asset.UploadMeta{
		FileName:  "demo.txt",
		FileSize:  int64(len(payload)),
		MimeType:  "text/plain",
		ChunkSize: int64(len(payload)),
	}

	store.EXPECT().GetMeta(gomock.Any(), uploadID).Return(&meta, nil)
	store.EXPECT().GetChunk(gomock.Any(), uploadID, 0).Return(payload, nil)
	stor.EXPECT().Put(gomock.Any(), "uploads/up001/demo.txt", gomock.Any(), int64(len(payload)), "text/plain").
		DoAndReturn(func(_ context.Context, _ string, body io.Reader, size int64, _ string) error {
			b, err := io.ReadAll(body)
			require.NoError(t, err)
			require.Equal(t, payload, b)
			require.Equal(t, int64(len(payload)), size)
			return nil
		})
	assets.EXPECT().Create(gomock.Any(), gomock.Any()).DoAndReturn(
		func(_ context.Context, in asset.CreateInput) (*asset.AssetVO, error) {
			require.Equal(t, "title", in.Title)
			require.Equal(t, "uploads/up001/demo.txt", in.StoragePath)
			return &asset.AssetVO{ID: "asset1", Title: in.Title}, nil
		},
	)
	store.EXPECT().DeleteSession(gomock.Any(), uploadID).Return(nil)

	svc := asset.NewUploadService(store, cfg, stor, assets)
	vo, err := svc.Complete(context.Background(), asset.UploadCompleteInput{
		UploadID:   uploadID,
		AssetTitle: "title",
		Type:       "other",
	})
	require.NoError(t, err)
	require.Equal(t, "asset1", vo.ID)
}

func TestUploadService_SaveChunk_WritesToStore(t *testing.T) {
	ctrl := gomock.NewController(t)
	store := mock.NewMockUploadSessionStore(ctrl)
	cfg := configmock.NewMockService(ctrl)
	stor := storagemock.NewMockStorage(ctrl)
	assets := mock.NewMockService(ctrl)

	meta := asset.UploadMeta{FileName: "a.bin", FileSize: 4, ChunkSize: 4}
	store.EXPECT().GetMeta(gomock.Any(), "up").Return(&meta, nil)
	store.EXPECT().SaveChunk(gomock.Any(), "up", 0, []byte("data")).Return(nil)

	svc := asset.NewUploadService(store, cfg, stor, assets)
	err := svc.SaveChunk(context.Background(), asset.UploadChunkInput{
		UploadID:   "up",
		ChunkIndex: 0,
		Body:       bytes.NewReader([]byte("data")),
	})
	require.NoError(t, err)
}