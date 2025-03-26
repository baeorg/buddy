package types

type User struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	AvatarURL string `json:"avatar_url"`
	Bio       string `json:"bio"`
}

type Mesg struct {
	ID      uint64 `json:"id"`
	ConvsID uint64 `json:"convs_id"`
	FromID  uint64 `json:"from_id"`
	Payload []byte `json:"payload"`
}
