package derrors

var (
	ErrDatabaseIssue   = baseError{code: CodeDBIsDown}
	ErrProductNotFound = baseError{code: ProductNotFound}
)
