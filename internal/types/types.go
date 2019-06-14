package types

import (
	"fmt"
	"strings"
)

type Validation struct {
	Valid  bool   `yaml:"valid"`
	Error  string `yaml:"error"`
	Claims Claims `yaml:"claims"`
}

type Claims struct {
	ID      string      `yaml:"id"`
	Subject string      `yaml:"subject"`
	Groups  []string    `yaml:"groups"`
	Extra   interface{} `yaml:"extra"`
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
