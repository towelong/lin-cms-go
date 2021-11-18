package vo

type FileVo struct {
	ID        int    `json:"id"`
	Path      string `json:"path"`
	Type      string `json:"type"`
	Name      string `json:"name"`
	Extension string `json:"extension"`
	Size      int    `json:"size"`
	Md5       string `json:"md5"`
	URL       string `json:"url"`
}
