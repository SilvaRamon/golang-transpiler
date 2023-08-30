package types

type SourceCode struct {
	Line       string
	LineNumber int
}

type Program struct {
	Type string
	Body interface{}
}