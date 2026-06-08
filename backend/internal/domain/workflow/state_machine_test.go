package workflow_test

import (
	"testing"

	"github.com/gin-mam/backend/internal/domain/workflow"
	"github.com/stretchr/testify/require"
)

func TestApplyAuditAction_Reject(t *testing.T) {
	level, status := workflow.ApplyAuditAction(1, 3, false)
	require.Equal(t, 1, level)
	require.Equal(t, workflow.AuditStatusRejected, status)
}

func TestApplyAuditAction_PassNotLast(t *testing.T) {
	level, status := workflow.ApplyAuditAction(0, 3, true)
	require.Equal(t, 1, level)
	require.Equal(t, workflow.AuditStatusAuditing, status)
}

func TestApplyAuditAction_PassLastLevel(t *testing.T) {
	level, status := workflow.ApplyAuditAction(2, 3, true)
	require.Equal(t, 3, level)
	require.Equal(t, workflow.AuditStatusPassed, status)
}

func TestApplyAuditAction_SingleLevelPass(t *testing.T) {
	level, status := workflow.ApplyAuditAction(0, 1, true)
	require.Equal(t, 1, level)
	require.Equal(t, workflow.AuditStatusPassed, status)
}
