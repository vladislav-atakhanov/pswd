package cli

import "strings"

func color(code string, elem ...string) string {
	return code + strings.Join(elem, "") + "\033[0m"
}

func gray(elem ...string) string    { return color("\033[90m", elem...) }
func blue(elem ...string) string    { return color("\033[34m", elem...) }
func green(elem ...string) string   { return color("\033[32m", elem...) }
func red(elem ...string) string     { return color("\033[31m", elem...) }
func yellow(elem ...string) string  { return color("\033[33m", elem...) }
func cyan(elem ...string) string    { return color("\033[36m", elem...) }
func magenta(elem ...string) string { return color("\033[35m", elem...) }
func white(elem ...string) string   { return color("\033[97m", elem...) }
func black(elem ...string) string   { return color("\033[30m", elem...) }

func reset(elem ...string) string     { return color("\033[0m", elem...) }
func bold(elem ...string) string      { return color("\033[1m", elem...) }
func underline(elem ...string) string { return color("\033[4m", elem...) }
