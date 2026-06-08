package workflow_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/gin-mam/backend/internal/domain/asset"
	"github.com/gin-mam/backend/internal/domain/asset/model"
	"github.com/gin-mam/backend/internal/domain/workflow"
	workflowmodel "github.com/gin-mam/backend/internal/domain/workflow/model"
	"github.com/gin-mam/backend/internal/infra/lock"
	"github.com/gin-mam/backend/internal/infra/tx"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func openWorkflowTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(fmt.Sprintf("file:%s?mode=memory&cache=private", t.Name())), &gorm.Config{
		NowFunc: func() time.Time { return time.Date(2026, 6, 8, 10, 0, 0, 0, time.UTC) },
	})
	if err != nil {
		t.Skipf("sqlite driver unavailable: %v", err)
	}
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(
		&model.Asset{},
		&workflowmodel.WorkflowDef{},
		&workflowmodel.WorkflowLevelUser{},
		&workflowmodel.WorkflowInstance{},
		&workflowmodel.WorkflowInstanceLevelUser{},
		&workflowmodel.WorkflowOperate{},
	))
	return db
}

func newWorkflowSvc(t *testing.T, db *gorm.DB) workflow.Service {
	t.Helper()
	assetRepo := asset.NewMySQLRepository(db)
	assetSvc := asset.NewService(assetRepo, nil, tx.NewManager(db), nil)
	wfRepo := workflow.NewMySQLRepository(db)
	return workflow.NewService(wfRepo, assetRepo, assetSvc, tx.NewManager(db), lock.NewMemoryLocker())
}

func TestSubmit_CreatesInstance(t *testing.T) {
	db := openWorkflowTestDB(t)
	now := db.NowFunc()
	require.NoError(t, db.Create(&model.Asset{
		ID: "a1", TenantID: "default", Title: "t", Type: "video", Status: 0, CreatedBy: "u1", CreatedAt: now, UpdatedAt: now,
	}).Error)
	require.NoError(t, db.Create(&workflowmodel.WorkflowDef{
		ID: "d1", TenantID: "default", Name: "def", AuditLevel: 1, Status: 1, CreatedAt: now, UpdatedAt: now,
	}).Error)
	require.NoError(t, db.Create(&workflowmodel.WorkflowLevelUser{
		ID: "lu1", TenantID: "default", WorkflowDefID: "d1", Level: 1, UserID: "auditor1",
	}).Error)

	svc := newWorkflowSvc(t, db)
	out, err := svc.Submit(context.Background(), "u1", workflow.SubmitRequest{AssetID: "a1", WorkflowDefID: "d1"})
	require.NoError(t, err)
	require.Equal(t, "a1", out.AssetID)
	require.Equal(t, workflow.AuditStatusPending, out.AuditStatus)

	var assetRow model.Asset
	require.NoError(t, db.First(&assetRow, "id = ?", "a1").Error)
	require.Equal(t, int8(1), assetRow.Status)
}

func TestAudit_PassLastLevelUpdatesAsset(t *testing.T) {
	db := openWorkflowTestDB(t)
	now := db.NowFunc()
	require.NoError(t, db.Create(&model.Asset{
		ID: "a2", TenantID: "default", Title: "t", Type: "video", Status: 0, CreatedBy: "u1", CreatedAt: now, UpdatedAt: now,
	}).Error)
	require.NoError(t, db.Create(&workflowmodel.WorkflowDef{
		ID: "d2", TenantID: "default", Name: "def", AuditLevel: 1, Status: 1, CreatedAt: now, UpdatedAt: now,
	}).Error)
	require.NoError(t, db.Create(&workflowmodel.WorkflowLevelUser{
		ID: "lu2", TenantID: "default", WorkflowDefID: "d2", Level: 1, UserID: "auditor1",
	}).Error)

	svc := newWorkflowSvc(t, db)
	inst, err := svc.Submit(context.Background(), "u1", workflow.SubmitRequest{AssetID: "a2", WorkflowDefID: "d2"})
	require.NoError(t, err)
	require.NoError(t, svc.Audit(context.Background(), "auditor1", workflow.AuditRequest{InstanceID: inst.ID, Pass: true}))

	var assetRow model.Asset
	require.NoError(t, db.First(&assetRow, "id = ?", "a2").Error)
	require.Equal(t, int8(2), assetRow.Status)
}
