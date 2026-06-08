package transcode_test

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/gin-mam/backend/internal/domain/transcode"
	transcodemodel "github.com/gin-mam/backend/internal/domain/transcode/model"
	"github.com/gin-mam/backend/internal/domain/transcode/mock"
	"github.com/gin-mam/backend/internal/domain/transcode/client"
	configmock "github.com/gin-mam/backend/internal/infra/config/mock"
	"github.com/gin-mam/backend/internal/infra/tx"
	"github.com/gin-mam/backend/internal/pkg/httpx"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func openTranscodeTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(fmt.Sprintf("file:%s?mode=memory&cache=private", t.Name())), &gorm.Config{
		NowFunc: func() time.Time { return time.Date(2026, 6, 8, 12, 0, 0, 0, time.UTC) },
	})
	if err != nil {
		t.Skipf("sqlite driver unavailable: %v", err)
	}
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&transcodemodel.TranscodeGroup{}, &transcodemodel.TranscodeTask{}))
	return db
}

func TestCreateTask_SubmitsExternalJobAndRunning(t *testing.T) {
	db := openTranscodeTestDB(t)
	now := db.NowFunc()
	require.NoError(t, db.Create(&transcodemodel.TranscodeGroup{
		ID: "g1", TenantID: "default", Name: "g", AvailableFlag: 1, CreatedBy: "u1", CreatedAt: now, UpdatedAt: now,
	}).Error)

	ctrl := gomock.NewController(t)
	cfg := configmock.NewMockService(ctrl)
	cfg.EXPECT().Get(gomock.Any(), transcode.ConfigKeyTranscodeURL).Return("http://transcode/jobs", nil)
	cfg.EXPECT().Get(gomock.Any(), transcode.ConfigKeyTranscodeCallbackURL).Return("http://mam/callback", nil)

	tc := mock.NewMockClient(ctrl)
	tc.EXPECT().SubmitJob(gomock.Any(), "http://transcode/jobs", gomock.Any()).
		DoAndReturn(func(_ context.Context, _ string, req client.SubmitJobRequest) (*client.SubmitJobResponse, error) {
			require.Equal(t, "a1", req.AssetID)
			require.Equal(t, "g1", req.TranscodeGroupID)
			return &client.SubmitJobResponse{JobID: "job-1"}, nil
		})

	svc := transcode.NewService(transcode.NewMySQLRepository(db), cfg, tx.NewManager(db), tc)
	out, err := svc.CreateTask(context.Background(), "u1", transcode.CreateTaskInput{
		AssetID: "a1", TranscodeGroupID: "g1",
	})
	require.NoError(t, err)
	require.Equal(t, transcodemodel.TaskStatusRunning, out.Status)
	require.Equal(t, "job-1", out.ExternalJobID)

	var row transcodemodel.TranscodeTask
	require.NoError(t, db.First(&row, "id = ?", out.ID).Error)
	require.Equal(t, transcodemodel.TaskStatusRunning, row.Status)
}

func TestHandleCallback_InvalidSignForbidden(t *testing.T) {
	ctrl := gomock.NewController(t)
	cfg := configmock.NewMockService(ctrl)
	cfg.EXPECT().Get(gomock.Any(), transcode.ConfigKeyTranscodeCallbackSecret).Return("secret", nil)

	svc := transcode.NewService(transcode.NewMySQLRepository(openTranscodeTestDB(t)), cfg, tx.NewManager(nil), mock.NewMockClient(ctrl))
	body := []byte(`{"taskId":"t1","status":2}`)
	err := svc.HandleCallback(context.Background(), body, "bad-sign")
	require.Error(t, err)
	biz, ok := err.(*httpx.BizError)
	require.True(t, ok)
	require.Equal(t, 403, biz.HTTPStatus)
	require.Equal(t, transcode.ErrCodeForbidden, biz.ErrCode)
}

func TestHandleCallback_ValidSignUpdatesTask(t *testing.T) {
	db := openTranscodeTestDB(t)
	now := db.NowFunc()
	require.NoError(t, db.Create(&transcodemodel.TranscodeTask{
		ID: "t1", TenantID: "default", AssetID: "a1", TranscodeGroupID: "g1",
		Status: transcodemodel.TaskStatusRunning, CreatedBy: "u1", CreatedAt: now, UpdatedAt: now,
	}).Error)

	secret := "secret-key"
	body, err := json.Marshal(transcode.CallbackInput{TaskID: "t1", Status: transcodemodel.TaskStatusSuccess, OutputPath: "/out.mp4"})
	require.NoError(t, err)
	sign := transcode.SignCallbackBody(body, secret)

	ctrl := gomock.NewController(t)
	cfg := configmock.NewMockService(ctrl)
	cfg.EXPECT().Get(gomock.Any(), transcode.ConfigKeyTranscodeCallbackSecret).Return(secret, nil)

	svc := transcode.NewService(transcode.NewMySQLRepository(db), cfg, tx.NewManager(db), mock.NewMockClient(ctrl))
	require.NoError(t, svc.HandleCallback(context.Background(), body, sign))

	var row transcodemodel.TranscodeTask
	require.NoError(t, db.First(&row, "id = ?", "t1").Error)
	require.Equal(t, transcodemodel.TaskStatusSuccess, row.Status)
	require.Equal(t, "/out.mp4", row.OutputPath)
}