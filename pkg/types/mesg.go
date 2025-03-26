package types

type MesgType uint64

const (
	UserTokenUpdate MesgType = 6000 + iota
	ConvsessionCreate
	ConvsessionJoin
)

type MesgInfo struct {
	MsgType MesgType
	Key     any
	Content []byte
}
