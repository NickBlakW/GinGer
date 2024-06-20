package generators

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/NickBlakW/ginger/generators/utils"
	"github.com/gin-gonic/gin"
)

func GenerateLocalApiScripts(path string, engine *gin.Engine) {
	filename := fmt.Sprintf("%sginger.js", path)

	file, err := os.Create(filename)
	if err != nil {
		log.Fatal("Could not create file")
	}

	defer file.Close()

	file.WriteString(utils.NoIndent("function printHello() {"))
	file.WriteString(utils.WithIndent("const msg = 'Hello from GinGer';", 1))
	file.WriteString(utils.WithIndent("console.log(msg);", 1))
	file.WriteString(utils.WithIndent("document.querySelector(body).innerHTML = msg;", 1))
	file.WriteString(utils.NoIndent("}\n"))

	routes := engine.Routes()

	generateAsyncFunctionsJS(file, routes)
}

func generateAsyncFunctionsJS(file *os.File, routes gin.RoutesInfo) {
	var generatedAmount int

	for i, route := range routes {
		if strings.Contains(route.Path, "/ginger") {
			continue
		} else if route.Method == "HEAD" {
			continue
		}

		funcName := strings.Replace(route.Path, "/", "_", -1)[1:] // replace all '/' with '_' and take slice from index 1
		funcDef := fmt.Sprintf("async function %s() {", funcName)

		file.WriteString(utils.NoIndent(funcDef))

		var jsLine string
		
		if route.Method == "GET" {
			jsLine = generateGetRequest(route.Path, i)
		} else if route.Method == "POST" {
			jsLine = generatePostRequest(route.Path, i)
		}

		file.WriteString(utils.WithIndent(jsLine, 1))
		file.WriteString(utils.NoIndent("}\n"))

		generatedAmount = i+1
	}

	fmt.Printf("\tGenerated %d functions\n", generatedAmount)
}

func generateGetRequest(route string, index int) string {
	fetchUrl := fmt.Sprintf("const res = await fetch(\"%s\");", route)

	jsLine := utils.NoIndent(fetchUrl)
	jsLine += utils.WithIndent("const jsonBody = await res.json();\n", 1)

	querySelect := fmt.Sprintf("document.querySelector('#div%d')\n\t\t.style.visibility = 'visible';\n", index)
	jsLine += utils.WithIndent(querySelect, 1)

	renderPre := fmt.Sprintf("\tdocument.querySelector('#pre%d').innerHTML =\n\t\tJSON.stringify(jsonBody, null, 2);", index)
	jsLine += renderPre

	return jsLine
}

func generatePostRequest(route string, index int) string {
	fetchLine := fmt.Sprintf("const res = await fetch(\"%s\", {", route)

	jsLine := utils.NoIndent(fetchLine)
	jsLine += utils.WithIndent("method: 'POST',", 2)

	reqBody := utils.WithIndent("body: JSON.stringify({", 2)
	reqBody += utils.WithIndent(fmt.Sprintf("\tusername: \"something\""), 2)
	reqBody += utils.WithIndent("}),", 2)

	jsLine += reqBody

	jsLine += utils.WithIndent("});", 1)

	return jsLine
}
