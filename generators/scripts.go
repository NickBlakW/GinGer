package generators

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/NickBlakW/ginger/types"
)

func GenerateLocalApiScripts(config types.Config) {
	path := fmt.Sprintf("%s/scripts/", config.LocalApiPath)

	err := os.MkdirAll(filepath.Dir(path), 0750)
	if err != nil {
		log.Fatal("Could not create filepath")
	}

	filename := fmt.Sprintf("%sginger.js", path)

	file, err := os.Create(filename)
	if err != nil {
		log.Fatal("Could not create file")
	}

	file.WriteString(noIndent("function printHello() {"))
	file.WriteString(withIndent("const msg = 'Hello from GinGer';", 1))
	file.WriteString(withIndent("console.log(msg);", 1))
	file.WriteString(withIndent("document.querySelector(body).innerHTML = msg;", 1))
	file.WriteString(noIndent("}"))
}

func noIndent(jsLine string) string {
	return fmt.Sprint(jsLine + "\n")
}

func withIndent(jsLine string, indents int) string {
	indent := ""

	for i := 0; i < indents; i++ {
		indent += "\t"
	}

	return fmt.Sprintf("%s%s%s", indent, jsLine, "\n")
}
