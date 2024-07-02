package generators

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/nickblakw/ginger/generators/utils"
	"github.com/nickblakw/ginger/requests"
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

	var codeLines []string
	
	for i, route := range routes {
		if strings.Contains(route.Path, "/ginger") {
			continue
		} else if route.Method == "HEAD" {
			continue
		}

		var jsLine string

		if route.Method == "GET" {
			function, elements := generateGetElement(route.Path, i)

			elemFunc := strings.Split(function, " ")[1]
			file.WriteString(utils.WithIndent(elemFunc, 1))

			jsLine = elements
		} else if route.Method == "POST" {
			function, elements := generatePostElementJS(route.Path, i)

			elemFunc := strings.Split(function, " ")[1]
			file.WriteString(utils.WithIndent(elemFunc, 1))

			jsLine = elements
		}

		codeLines = append(codeLines, jsLine)
	}

	file.WriteString(utils.NoIndent("});\n"))

	for _, line := range codeLines {
		file.WriteString(utils.NoIndent(line))
	}
}

func createElementJS(element string) string {
	elem := fmt.Sprintf("document.createElement('%s');", element)

	return elem
}

func appendChildrenJS(elem string, inserts []string) string {
	var js string
	
	for _, el := range inserts {
		js += utils.WithIndent(fmt.Sprintf("%s.appendChild(%s);", elem, el), 1)
	}

	return js
}

func createButtonJS(btnName string, text string) string {
	jsBtn := utils.WithIndent(fmt.Sprintf("const %s = %s", btnName, createElementJS("button")), 1)
	jsBtn += utils.WithIndent(fmt.Sprintf("%s.id = '%s';", btnName, btnName), 1)
	jsBtn += utils.WithIndent(fmt.Sprintf("%s.innerText = '%s';", btnName, text), 1)

	return jsBtn
}

func addClickEventListenerJS(button string, funcName string) string {
	btnFunc := fmt.Sprintf("%s.addEventListener('click', () => {\n", button)
	btnFunc += fmt.Sprintf("\t\t%s\n", funcName)
	btnFunc += "\t});\n"

	return btnFunc
}

func createElementTitle(title string, path string) string {
	elem := utils.WithIndent(fmt.Sprintf("const %s = %s", title, createElementJS("h2")), 1)
	elem += utils.WithIndent(fmt.Sprintf("%s.innerText = '%s';\n", title, path), 1)

	return elem
}

func generateFuncNameJS(path string) (string, string) {
	name := strings.Replace(path, "/", "_", -1)[1:]
	name += "()"

	funcName := fmt.Sprintf("function elem_%s {\n", name)

	return name, funcName
}

func createInputFields(path string, index int) (string, []string) {
	var dtoFields utils.DTOFields

	for _, req := range requests.APIRequestRegistry {
		if req.Path != path {
			return "", []string{}
		}

		dtoFields = utils.GetDTOFields(req.DTO)
	}

	var inputs string
	var inputNames []string

	for i, name := range dtoFields.Names {
		var tempInput string
		formattedName := strings.ToLower(name[:1]) + name[1:]

		fieldName := fmt.Sprintf("%s%d", formattedName, index)
		inputNames = append(inputNames, fieldName)

		label := fmt.Sprintf("label_%s", fieldName)
		inputNames = append(inputNames, label)
		
		elem := fmt.Sprintf("const %s = %s\n", label, createElementJS("label"))
		tempInput += utils.WithIndent(elem, 1)

		tempInput += utils.WithIndent(fmt.Sprintf("%s.htmlFor = '%s';\n", label, fieldName), 1)
		tempInput += utils.WithIndent(fmt.Sprintf("%s.innerHTML = '%s';\n", label, name), 1)
		tempInput += appendChildrenJS(fmt.Sprintf("in_div%d", index), []string{label})

		tempInput += utils.WithIndent(fmt.Sprintf("const %s = %s\n", fieldName, createElementJS("input")), 1)
		tempInput += utils.WithIndent(fmt.Sprintf("%s.id = '%s';", fieldName, fieldName), 1)

		for j, typ := range dtoFields.Types {

			if i == j && typ == "number" {
				tempInput += utils.WithIndent(fmt.Sprintf("%s.setAttribute('type', 'number');\n", fieldName), 1)
			} else if i == j && typ != "number" {
				tempInput += utils.WithIndent(fmt.Sprintf("%s.setAttribute('type', 'text');\n", fieldName), 1)
			} else {
				continue
			}
		}

		inputs += tempInput
		inputNames = append(inputNames, fieldName)
	}

	return inputs, inputNames
}

func generateStdElements(path string, index int, sendFunc string, isPostReq bool) string {
	title := fmt.Sprintf("apiTitle%d", index)
	jsLine := createElementTitle(title, path)

	div := fmt.Sprintf("div%d", index)
	button := fmt.Sprintf("btnSend%d", index)
	pre := fmt.Sprintf("pre%d", index)
	innerDiv := fmt.Sprintf("in_div%d", index)
	hideBtn := fmt.Sprintf("hideBtn%d", index)

	jsLine += utils.WithIndent(fmt.Sprintf("const %s = %s", div, createElementJS("div")), 1)
	jsLine += utils.WithIndent(fmt.Sprintf("%s.id = '%s';\n", div, div), 1)
	
	jsLine += createButtonJS(button, "Send")
	jsLine += createButtonJS(hideBtn, "Hide")

	jsLine += utils.WithIndent(addClickEventListenerJS(button, fmt.Sprintf("%s;", sendFunc)), 1)
	jsLine += utils.WithIndent(addClickEventListenerJS(hideBtn, fmt.Sprintf("hideElement('%s')", innerDiv)), 1)

	jsLine += utils.WithIndent(fmt.Sprintf("const %s = %s", innerDiv, createElementJS("div")), 1)
	jsLine += utils.WithIndent(fmt.Sprintf("%s.id = '%s';", innerDiv, innerDiv), 1)

	jsLine += utils.WithIndent(fmt.Sprintf("const %s = %s", pre, createElementJS("pre")), 1)
	jsLine += utils.WithIndent(fmt.Sprintf("%s.id = '%s';", pre, pre), 1)
	jsLine += utils.WithIndent(fmt.Sprintf("%s.style.display = 'none';\n", innerDiv), 1)
	
	showBtn := fmt.Sprintf("showBtn%d", index)
	jsLine += createButtonJS(showBtn, "Try it out!")
	jsLine += utils.WithIndent(addClickEventListenerJS(showBtn, fmt.Sprintf("showElement('%s')", innerDiv)), 1)

	if isPostReq {
		inputs := []string{button, hideBtn}

		htmlInput, jsFields := createInputFields(path, index)
		jsLine += htmlInput

		inputs = append(inputs, jsFields...)
		inputs = append(inputs, pre)

		jsLine += utils.NoIndent(appendChildrenJS(innerDiv, inputs))
	} else {
		jsLine += utils.NoIndent(appendChildrenJS(innerDiv, []string{button, hideBtn, pre}))
	}

	jsLine += utils.NoIndent(appendChildrenJS(div, []string{title, showBtn , innerDiv}))
	
	jsLine += utils.WithIndent(fmt.Sprintf("document.body.insertAdjacentElement('beforeEnd', %s);", div), 1)
	jsLine += utils.NoIndent("}")

	return jsLine
}

func generateGetElement(path string, index int) (string, string) {
	jsLine := fmt.Sprintf("// #region GET-request%d\n", index)

	apiReq, funcName := generateFuncNameJS(path)
	jsLine += utils.NoIndent(funcName)

	jsLine += generateStdElements(path, index, apiReq, false)

	jsLine += utils.WithIndent("// #endregion", 1)

	return funcName, jsLine
}

func generatePostElementJS(path string, index int) (string, string) {
	jsLine := fmt.Sprintf("// #region POST-request%d\n", index)

	apiReq, funcName := generateFuncNameJS(path)
	jsLine += utils.NoIndent(funcName)

	jsLine += generateStdElements(path, index, apiReq, true)

	return funcName, jsLine
}