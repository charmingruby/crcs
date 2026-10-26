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
	EventRating     int
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
	CartelMatches      []CartelMatch
	Reach              int
	Height             int
	Weight             int
}

type Message struct {
	SeasonID   string
	Contenders []Contender
}
