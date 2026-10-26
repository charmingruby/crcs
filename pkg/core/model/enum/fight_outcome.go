package enum

type FightOutcome string

const (
	WinFightOutcome       FightOutcome = "W"
	LossFightOutcome      FightOutcome = "L"
	DrawFightOutcome      FightOutcome = "D"
	NoContestFightOutcome FightOutcome = "NC"
)

func (o FightOutcome) Valid() bool {
	switch o {
	case WinFightOutcome, LossFightOutcome, DrawFightOutcome, NoContestFightOutcome:
		return true
	default:
		return false
	}
}
