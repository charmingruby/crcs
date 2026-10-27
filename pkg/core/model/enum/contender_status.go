package enum

type ContenderStatus string

const (
	IdleContenderStatus       ContenderStatus = "idle"
	ActiveContenderStatus     ContenderStatus = "active"
	ContractedContenderStatus ContenderStatus = "contracted"
	RejectedContenderStatus   ContenderStatus = "rejected"
)

func (s ContenderStatus) Valid() bool {
	switch s {
	case ActiveContenderStatus, ContractedContenderStatus, RejectedContenderStatus, IdleContenderStatus:
		return true
	default:
		return false
	}
}

func (s ContenderStatus) CanTransitionTo(next ContenderStatus) bool {
	if s == next {
		return s.Valid()
	}

	switch s {
	case IdleContenderStatus:
		return next == ActiveContenderStatus || next == ContractedContenderStatus || next == RejectedContenderStatus
	case ActiveContenderStatus:
		return next == ContractedContenderStatus || next == RejectedContenderStatus
	case ContractedContenderStatus, RejectedContenderStatus:
		return next == ContractedContenderStatus || next == RejectedContenderStatus
	default:
		return false
	}
}
