package generators

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/NickBlakW/ginger/generators/utils"
	"github.com/gin-gonic/gin"
)

func GenerateUiElementScript(path string, engine *gin.Engine) {
	filename := fmt.Sprintf("%sginger.elements.js", path)

	file, err := os.Create(filename)
	if err != nil {
		log.Fatal("Could not create file")
	}

	defer file.Close()

	file.WriteString(utils.NoIndent("window.addEventListener('load', () => {"))
	
	routes := engine.Routes()
	
	for i, route := range routes {
		if strings.Contains(route.Path, "/ginger") {
			continue
		} else if route.Method == "HEAD" {
			continue
		}

		var jsLine string

		if route.Method == "GET" {
			jsLine = generateGetElement(route.Path, i)
		}

		file.WriteString(utils.WithIndent(jsLine, 1))
	}

	file.WriteString(utils.NoIndent("});\n"))
}

func createElementJS(element string) string {
	elem := fmt.Sprintf("document.createElement('%s')", element)

	return elem
}

func generateGetElement(path string, index int) string {
	jsLine := fmt.Sprintf("// #region GET-req%d\n", index)

	jsLine += utils.WithIndent(fmt.Sprintf("const apiTitle%d = %s;", index, createElementJS("h2")), 1)
	jsLine += utils.WithIndent(fmt.Sprintf("apiTitle%d.innerText = '%s';\n", index, path), 1)

	div := fmt.Sprintf("div%d", index)
	button := fmt.Sprintf("btnSend%d", index)

	jsLine += utils.WithIndent(fmt.Sprintf("const %s = %s;", div, createElementJS("div")), 1)
	jsLine += utils.WithIndent(fmt.Sprintf("const %s = %s;", button, createElementJS("button")), 1)
	jsLine += utils.WithIndent(fmt.Sprintf("%s.innerText = 'Send';\n", button), 1)
	jsLine += utils.WithIndent(fmt.Sprintf("%s.appendChild(%s)", div, button), 1)

	jsLine += utils.WithIndent("// #endregion", 1)

	return jsLine
}