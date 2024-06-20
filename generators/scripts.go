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
			fetchUrl := fmt.Sprintf("const res = await fetch(\"%s\");", route.Path)

		jsLine += utils.NoIndent(fetchUrl)
		jsLine += utils.WithIndent("const jsonBody = await res.json();\n", 1)
				
		divId := fmt.Sprintf("div%d", i)
		query := fmt.Sprintf("document.querySelector('#%s')", divId)

		jsLine += utils.WithIndent(query, 1)
		} else if route.Method == "POST" {
			jsLine += fmt.Sprintf("const res = await fetch(\"%s\", {", route.Path)
			jsLine += utils.WithIndent("method: 'POST',", 2)

			reqBody := "body: JSON.stringify({\n"
			reqBody += fmt.Sprintf("")

			jsLine += utils.WithIndent(reqBody, 2)

			jsLine += utils.WithIndent("", 1)
		}

		file.WriteString(utils.WithIndent(jsLine, 1))

		file.WriteString(utils.NoIndent("}"))

		fmt.Println(fmt.Sprintf("Generated %d functions", i))
	}
}

func generateGetRequest(route string, jsLine string, index int) {
	fetchUrl := fmt.Sprintf("const res = await fetch(\"%s\");", route)

	jsLine += utils.NoIndent(fetchUrl)
	jsLine += utils.WithIndent("const jsonBody = await res.json();\n", 1)
			
	divId := fmt.Sprintf("div%d", index)
	query := fmt.Sprintf("document.querySelector('#%s')", divId)

	jsLine += utils.WithIndent(query, 1)
}
