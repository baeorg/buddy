package types

type ConvsReq struct {
	Title   string   `json:"title"`
	UserIDs []uint64 `json:"user_ids" validate:"required"`
}
