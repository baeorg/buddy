package storage

type MesgType uint64

const (
	UserTokenUpdate MesgType = 6000 + iota
)

type MesgInfo struct {
	MsgType MesgType
	Key     any
	Content []byte
}

const (
	PermiPrefix string = "tokens:"
)
