package commands

import "strings"

func Judge(names []string) string {
	return strings.Join(names, ",")
}
