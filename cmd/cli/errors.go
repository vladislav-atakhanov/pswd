package cli

import (
	"errors"
	"fmt"
	"strings"
)

var (
	tooManyArgumentsErr = errors.New("to many arguments")
)

func TooManyArgumentsErr() error {
	return tooManyArgumentsErr
}

func PassArgumentsErr(args ...string) error {
	return fmt.Errorf("pass %s", humanJoin(args))
}

func humanJoin(items []string) string {
	n := len(items)
	switch n {
	case 0:
		return ""
	case 1:
		return items[0]
	case 2:
		return items[0] + " and " + items[1]
	default:
		return strings.Join(items[:n-1], ", ") + " and " + items[n-1]
	}
}
