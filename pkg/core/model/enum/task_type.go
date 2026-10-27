package enum

type TaskType string

const (
	IngestContendersTaskType  TaskType = "ingest_contenders"
	SeasonMatchmakingTaskType TaskType = "season_matchmaking"
)

func (t TaskType) Valid() bool {
	switch t {
	case IngestContendersTaskType, SeasonMatchmakingTaskType:
		return true
	default:
		return false
	}
}
