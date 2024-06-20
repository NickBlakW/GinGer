package utils

import "fmt"

func NoIndent(jsLine string) string {
	return fmt.Sprint(jsLine + "\n")
}

func WithIndent(jsLine string, indents int) string {
	indent := ""

	for i := 0; i < indents; i++ {
		indent += "\t"
	}

	return fmt.Sprintf("%s%s%s", indent, jsLine, "\n")
}
