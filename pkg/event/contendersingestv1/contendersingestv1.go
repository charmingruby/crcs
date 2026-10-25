package contendersingestv1

const Topic = "contenders.ingest"

type CartelMatch struct {
	EventName       string
	EventDate       string
	Method          string
	Outcome         string
	OpponentName    string
	WeightClass     string
	FightDuration   string
	EventRating     float64
	ScheduledRounds int
	ActualRounds    int
}

type Contender struct {
	FirstName          string
	LastName           string
	Nickname           string
	Birthdate          string
	PlaceOfBirth       string
	CurrentWeightClass string
	Gender             string
	SeasonID           string
	CartelMatches      []CartelMatch
	Reach              float64
	Height             float64
	Weight             float64
}

type Message struct {
	Contenders []Contender
}
