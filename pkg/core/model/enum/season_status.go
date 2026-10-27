package enum

type SeasonStatus string

const (
	DraftSeasonStatus    SeasonStatus = "draft"
	AcceptedSeasonStatus SeasonStatus = "accepted"
	RejectedSeasonStatus SeasonStatus = "rejected"
)

func (s SeasonStatus) Valid() bool {
	switch s {
	case DraftSeasonStatus, AcceptedSeasonStatus, RejectedSeasonStatus:
		return true
	default:
		return false
	}
}

func (s SeasonStatus) CanTransitionTo(next SeasonStatus) bool {
	if s == next {
		return s.Valid()
	}

	switch s {
	case DraftSeasonStatus:
		return next == AcceptedSeasonStatus || next == RejectedSeasonStatus
	case AcceptedSeasonStatus, RejectedSeasonStatus:
		return next == AcceptedSeasonStatus || next == RejectedSeasonStatus
	default:
		return false
	}
}
