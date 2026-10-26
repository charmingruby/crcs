package enum

type WeightClass string

const (
	Flyweight        WeightClass = "Flyweight"
	Bantamweight     WeightClass = "Bantamweight"
	Featherweight    WeightClass = "Featherweight"
	Lightweight      WeightClass = "Lightweight"
	Welterweight     WeightClass = "Welterweight"
	Middleweight     WeightClass = "Middleweight"
	LightHeavyweight WeightClass = "Light Heavyweight"
	Heavyweight      WeightClass = "Heavyweight"
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
