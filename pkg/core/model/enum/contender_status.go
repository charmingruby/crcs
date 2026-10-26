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
