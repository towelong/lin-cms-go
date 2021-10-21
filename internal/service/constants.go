package service

type LoginType int
type UploadType string

const (
	UserPassword LoginType = iota
	Mini
)

const (
	Root  = 1
	Guest = 2
	User  = 3
)

const (
	Local  UploadType = "LOCAL"
	Remote UploadType = "REMOTE"
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

func (u UploadType) String() string {
	switch u {
	case Local:
		return "LOCAL"
	case Remote:
		return "REMOTE"
	}
	return ""
}
