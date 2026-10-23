package enum

type FightMethod string

const (
	SubmissionFightMethod FightMethod = "SUB"
	KOFightMethod         FightMethod = "KO"
	TKOFightMethod        FightMethod = "TKO"
	DecisionFightMethod   FightMethod = "DEC"
	NoContestFightMethod  FightMethod = "NOC"
)
