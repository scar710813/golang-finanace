package i_query

type IQueryHandler[Query any, Result any] interface {
	Handle(query Query) (Result, error)
}
