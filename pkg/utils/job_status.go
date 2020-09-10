package utils

const (
	StatusCreated    = "CREATED"
	StatusStarted    = "STARTED"
	StatusPending    = "PENDING"
	StatusResending  = "RESENDING"
	StatusStored     = "STORED"
	StatusRecovering = "RECOVERING"
	StatusWarning    = "WARNING"
	StatusFailed     = "FAILED"
	StatusMined      = "MINED"
	StatusNeverMined = "NEVER_MINED"
)

func IsFinalStatus(status string) bool {
	return status == StatusFailed || status == StatusMined || status == StatusNeverMined || status == StatusStored
}
