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

			file.WriteString(utils.WithIndent(fmt.Sprintf("elem_%s;", function), 1))

			jsLine = elements
		} else if route.Method == "POST" {
			// for _, req := range requests.APIRequestRegistry {
			// 	if req.Path != route.Path {
			// 		continue
			// 	}

			// 	dto := reflect.ValueOf(req.DTO)
			// 	typeOfDto := dto.Type()

			// 	for i := 0;  i < dto.NumField(); i++ {
			// 		tsType := utils.GenerateTSType(typeOfDto.Field(i).Type)

			// 		jsLine = fmt.Sprintf("DTO Field: %s\nField Type: %s\nTS Type: %s\n\n",
			// 			typeOfDto.Field(i).Name,
			// 			typeOfDto.Field(i).Type,
			// 			tsType,
			// 		)
			// 	}
			// }
		}

		codeLines = append(codeLines, jsLine)
	}

	file.WriteString(utils.NoIndent("});\n"))

	for _, line := range codeLines {
		file.WriteString(utils.NoIndent(line))
	}
}

func createElementJS(element string) string {
	elem := fmt.Sprintf("document.createElement('%s')", element)

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
	jsBtn := utils.WithIndent(fmt.Sprintf("const %s = %s;", btnName, createElementJS("button")), 1)
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
	elem := utils.WithIndent(fmt.Sprintf("const %s = %s;", title, createElementJS("h2")), 1)
	elem += utils.WithIndent(fmt.Sprintf("%s.innerText = '%s';\n", title, path), 1)

	return elem
}

func generateGetElement(path string, index int) (string, string) {
	jsLine := fmt.Sprintf("// #region GET-request%d\n", index)

	funcName := strings.Replace(path, "/", "_", -1)[1:]
	funcName += "()"

	jsLine += utils.NoIndent(fmt.Sprintf("function elem_%s {\n", funcName))

	title := fmt.Sprintf("apiTitle%d", index)
	jsLine += createElementTitle(title, path)

	div := fmt.Sprintf("div%d", index)
	button := fmt.Sprintf("btnSend%d", index)
	pre := fmt.Sprintf("pre%d", index)
	innerDiv := fmt.Sprintf("in_div%d", index)
	hideBtn := fmt.Sprintf("hideBtn%d", index)

	jsLine += utils.WithIndent(fmt.Sprintf("const %s = %s;", div, createElementJS("div")), 1)
	jsLine += utils.WithIndent(fmt.Sprintf("%s.id = '%s';\n", div, div), 1)

	jsLine += createButtonJS(button, "Send")
	jsLine += createButtonJS(hideBtn, "Hide")

	jsLine += utils.WithIndent(addClickEventListenerJS(button, fmt.Sprintf("%s;", funcName)), 1)
	jsLine += utils.WithIndent(addClickEventListenerJS(hideBtn, fmt.Sprintf("hideElement('%s')", innerDiv)), 1)

	jsLine += utils.WithIndent(fmt.Sprintf("const %s = %s;", innerDiv, createElementJS("div")), 1)
	jsLine += utils.WithIndent(fmt.Sprintf("%s.id = '%s';", innerDiv, innerDiv), 1)

	jsLine += utils.WithIndent(fmt.Sprintf("const %s = %s;", pre, createElementJS("pre")), 1)
	jsLine += utils.WithIndent(fmt.Sprintf("%s.id = '%s';", pre, pre), 1)
	jsLine += utils.WithIndent(fmt.Sprintf("%s.style.display = 'none';\n", innerDiv), 1)
	
	jsLine += utils.NoIndent(appendChildrenJS(div, []string{title, button, innerDiv}))
	jsLine += utils.NoIndent(appendChildrenJS(innerDiv, []string{hideBtn, pre}))
	
	jsLine += utils.WithIndent(fmt.Sprintf("document.body.insertAdjacentElement('beforeEnd', %s);", div), 1)
	jsLine += utils.NoIndent("}")

	jsLine += utils.WithIndent("// #endregion", 1)

	return funcName, jsLine
}

// func generatePostElementJS(path string, index int) string {

// }