package workflow

const (
	AuditStatusPending  int8 = 0
	AuditStatusAuditing int8 = 1
	AuditStatusPassed   int8 = 2
	AuditStatusRejected int8 = 3
)

const (
	AssetStatusDraft    int8 = 0
	AssetStatusAuditing int8 = 1
	AssetStatusPassed   int8 = 2
	AssetStatusRejected int8 = 3
)

const (
	OperateActionSubmit int8 = 1
	OperateActionPass   int8 = 2
	OperateActionReject int8 = 3
	OperateActionRevoke int8 = 4
)

const (
	ErrCodeNotFound      = 40402
	ErrCodeLockConflict  = 40901
	ErrCodeInvalidState  = 42202
)

func ApplyAuditAction(currentLevel, auditLevel int, pass bool) (newLevel int, newStatus int8) {
	if !pass {
		return currentLevel, AuditStatusRejected
	}
	nextLevel := currentLevel + 1
	if nextLevel >= auditLevel {
		return nextLevel, AuditStatusPassed
	}
	return nextLevel, AuditStatusAuditing
}

func AssigneeLevel(instanceLevel int) int {
	return instanceLevel + 1
}
