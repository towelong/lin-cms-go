package model

type File struct {
	BaseModel
	Path      string `json:"path"`
	Type      string `json:"type"`
	Name      string `json:"name"`
	Extension string `json:"extension"`
	Size      int    `json:"size"`
	Md5       string `json:"md5"`
}
