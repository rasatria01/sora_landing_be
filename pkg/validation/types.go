package validation

type Enum interface {
	IsValidEnum() bool
}

const (
	CheckOnlyAlphabet = `^[A-Za-z\s._\-\/]+$`
	Digits            = `^\d+$`
)
