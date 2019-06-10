package types

type Validation struct {
	Valid  bool
	Error  string
	Claims Claims
}

type Claims struct {
	ID      string
	Subject string
	Groups  []string
	Extra   interface{}
}
