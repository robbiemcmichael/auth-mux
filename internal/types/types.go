package types

import (
	"fmt"
	"strings"
)

type Validation struct {
	Valid  bool   `json:"valid"`
	Error  string `json:"error"`
	Claims Claims `json:"claims"`
}

type Claims struct {
	ID      string      `json:"id"`
	Subject string      `json:"subject"`
	Groups  []string    `json:"groups"`
	Extra   interface{} `json:"extra"`
}

type Assertion struct {
	IDPrefix      string
	SubjectPrefix string
	GroupPrefix   string
}

func (v *Validation) Assert(a Assertion) {
	if !strings.HasPrefix(v.Claims.ID, a.IDPrefix) {
		v.Valid = false
		v.Error += fmt.Sprintf("Expected ID %q to have prefix %q\n", v.Claims.ID, a.IDPrefix)
	}

	if !strings.HasPrefix(v.Claims.Subject, a.SubjectPrefix) {
		v.Valid = false
		v.Error += fmt.Sprintf("Expected subject %q to have prefix %q\n", v.Claims.Subject, a.SubjectPrefix)
	}

	for _, group := range v.Claims.Groups {
		if !strings.HasPrefix(group, a.SubjectPrefix) {
			v.Valid = false
			v.Error += fmt.Sprintf("Expected group %q to have prefix %q\n", group, a.GroupPrefix)
		}
	}
}
