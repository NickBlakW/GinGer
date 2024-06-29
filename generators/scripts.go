package generators

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/NickBlakW/ginger/generators/utils"
	"github.com/NickBlakW/ginger/requests"
	"github.com/gin-gonic/gin"
)

func GenerateLocalApiScripts(path string, engine *gin.Engine) {
	filename := fmt.Sprintf("%sginger.js", path)

	file, err := os.Create(filename)
	if err != nil {
		log.Fatal("Could not create file")
	}

	defer file.Close()

	file.WriteString(generateDefaultFuncsJS())

	routes := engine.Routes()

	generateAsyncFunctionsJS(file, routes)
}

func generateDefaultFuncsJS() string {
	function := "function hideElement(id) {\n"
	function += utils.WithIndent("document.querySelector(`#${id}`).style.display = 'none';", 1)
	function += utils.NoIndent("}\n")

	return function
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
	
	jsLine += generateGetJSONLines(index)

	return jsLine
}

func generatePostRequest(route string, index int) string {
	var dtoFields utils.DTOFields

	for _, req := range requests.APIRequestRegistry {
		if req.Path != route {
			return ""
		}

		dtoFields = utils.GetDTOFields(req.DTO)
	}
	
	var inputs string
	var inputFields string

	for i, name := range dtoFields.Names {
		var tempInput string
		name = strings.ToLower(name[:1]) + name[1:] // lowercase first letter, rest remains the same

		fieldName := fmt.Sprintf("%s%d", name, index)

		for _, typ := range dtoFields.Types {
			if typ == "number" {
				tempInput = fmt.Sprintf("%s: %v,", name, name)
			} else {
				tempInput = fmt.Sprintf("%s: '%s',", name, name)
			}
		}

		input := fmt.Sprintf("const %s = document.querySelector('#%s').value", name, fieldName)

		// formatting
		if i > 0 {
			input = utils.WithIndent(input, 1)
		}

		inputs += utils.NoIndent(input)
		inputFields += utils.WithIndent(tempInput, 3)
	}

	jsLine := utils.NoIndent(inputs)

	fetchLine := fmt.Sprintf("const res = await fetch('%s', {", route)

	jsLine += utils.WithIndent(fetchLine, 1)
	jsLine += utils.WithIndent("method: 'POST',", 2)

	reqBody := utils.WithIndent("body: JSON.stringify({", 2)

	reqBody += inputFields
	
	// handle request
	reqBody += utils.WithIndent("}),", 2)
	jsLine += reqBody
	jsLine += utils.WithIndent("});\n", 1)

	jsLine += generateGetJSONLines(index)

	return jsLine
}

func generateGetJSONLines(index int) string {
	jsLine := utils.WithIndent("const jsonBody = await res.json();\n", 1)

	querySelect := fmt.Sprintf("document.querySelector('#in_div%d')\n\t\t.style.display = 'flex';\n", index)
	jsLine += utils.WithIndent(querySelect, 1)

	renderPre := fmt.Sprintf("\tdocument.querySelector('#pre%d').innerHTML =\n\t\tJSON.stringify(jsonBody, null, 2);", index)
	jsLine += renderPre

	return jsLine
}