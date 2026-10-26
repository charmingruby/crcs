package enum

type TaskType string

const (
	IngestContendersTaskType TaskType = "ingest_contenders"
)

func (t TaskType) Valid() bool {
	return t == IngestContendersTaskType
}
