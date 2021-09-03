package service

type LoginType int

const (
	UserPassword LoginType = iota
	Mini
)

const (
	Root  = 1
	Guest = 2
)

func (l LoginType) String() string{
	switch l {
	case UserPassword:
		return "USERNAME_PASSWORD"
	case Mini:
		return "MINI"
	}
	return ""
}