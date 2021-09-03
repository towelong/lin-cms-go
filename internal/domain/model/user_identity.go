package model

type UserIdentity struct {
	BaseModel
	UserID       int    `json:"user_id"`
	IdentityType string `json:"identity_type"`
	Identifier   string `json:"identifier"`
	Credential   string `json:"credential"`
}
