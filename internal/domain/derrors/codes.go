package derrors

const (
	CodeUnhandledError ErrorCode = -1
)

// Please Note that all Codes must be here
const (
	CodeDBIsDown    = ErrorCode(1000)
	ProductNotFound = ErrorCode(1001)
)
