package enum

type TaskStatus string

const (
	PendingTaskStatus   TaskStatus = "pending"
	RunningTaskStatus   TaskStatus = "running"
	CompletedTaskStatus TaskStatus = "completed"
	FailedTaskStatus    TaskStatus = "failed"
)

func (s TaskStatus) Valid() bool {
	switch s {
	case PendingTaskStatus, RunningTaskStatus, CompletedTaskStatus, FailedTaskStatus:
		return true
	default:
		return false
	}
}

func (s TaskStatus) CanTransitionTo(target TaskStatus) bool {
	switch s {
	case PendingTaskStatus:
		return target == RunningTaskStatus
	case RunningTaskStatus:
		return target == CompletedTaskStatus || target == FailedTaskStatus
	case FailedTaskStatus:
		return target == RunningTaskStatus
	default:
		return false
	}
}
