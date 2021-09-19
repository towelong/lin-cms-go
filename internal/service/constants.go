package service

type LoginType int
type GroupLevel int

const (
	UserPassword LoginType = iota
	Mini
)

const (
	Root  = 1
	Guest = 2
	User  = 3
)

func (l LoginType) String() string {
	switch l {
	case UserPassword:
		return "USERNAME_PASSWORD"
	case Mini:
		return "MINI"
	}
	return ""
}
