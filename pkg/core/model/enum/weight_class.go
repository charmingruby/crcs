package enum

type WeightClass string

const (
	Flyweight        WeightClass = "flyweight"
	Bantamweight     WeightClass = "bantamweight"
	Featherweight    WeightClass = "featherweight"
	Lightweight      WeightClass = "lightweight"
	Welterweight     WeightClass = "welterweight"
	Middleweight     WeightClass = "middleweight"
	LightHeavyweight WeightClass = "light_heavyweight"
	Heavyweight      WeightClass = "heavyweight"
)

func (w WeightClass) Valid() bool {
	switch w {
	case Flyweight, Bantamweight, Featherweight, Lightweight,
		Welterweight, Middleweight, LightHeavyweight, Heavyweight:
		return true
	default:
		return false
	}
}
