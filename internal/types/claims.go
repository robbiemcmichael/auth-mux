package types

type Result struct {
	Valid  bool
	Error  string
	Claims Claims
}

type Claims struct {
	UID    string
	User   string
	Groups []string
	Extra  interface{}
}
