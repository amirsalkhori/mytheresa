package domain

type Pagination struct {
	Next     uint32
	Prev     uint32
	PageSize uint32
}

type HashedPagination struct {
	Next     string
	Prev     string
	PageSize uint32
}
