package derrors

var (
	ErrDatabaseIssue        = baseError{code: CodeDBIsDown}
	ErrOrderProductNotFound = baseError{code: ProductNotFound}
)
