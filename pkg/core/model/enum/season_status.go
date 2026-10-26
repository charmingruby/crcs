package enum

type SeasonStatus string

const (
	DraftSeasonStatus SeasonStatus = "draft"
)

func (s SeasonStatus) Valid() bool {
	return s == DraftSeasonStatus
}
