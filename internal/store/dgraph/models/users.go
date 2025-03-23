package models

type User struct {
	Nickname     string   `json:"nickname,omitempty"`
	Email        string   `json:"email,omitempty"`
	ID           string   `json:"uuid,omitempty"`
	PasswordHash []byte   `json:"pass_hash,omitempty"`
	DType        []string `json:"dgraph.type,omitempty"`
}
