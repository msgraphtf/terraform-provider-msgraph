package transform

import (
	"strings"
)

func upperFirst(s string) string {
	return strings.ToUpper(s[0:1]) + s[1:]
}
