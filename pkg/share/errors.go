package share

import "fmt"

var (
	Success            = fmt.Errorf("success")
	ErrInvalidRequest  = fmt.Errorf("invalid request")
	ErrNotSupported    = fmt.Errorf("not supported event type")
	ErrInternalError   = fmt.Errorf("internal error")
	ErrAtLeastTwoUsers = fmt.Errorf("at least two users are required")
)

const (
	SuccessCode uint64 = 4000 + iota
	ErrInvalidRequestCode
	ErrNotSupportedCode
	ErrInternalErrorCode
	ErrAtLeastTwoUsersCode
)
