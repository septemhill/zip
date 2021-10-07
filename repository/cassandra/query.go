package cassandra

type query struct {
	stmt  string
	names []string
}

type statement struct {
	del query
	ins query
	sel query
	up  query
}
