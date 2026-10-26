package enum

type Gender string

const (
	MaleGender   Gender = "M"
	FemaleGender Gender = "F"
)

func (g Gender) Valid() bool {
	return g == MaleGender || g == FemaleGender
}
