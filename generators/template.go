package generators

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/nickblakw/ginger/generators/utils"
)

func GenerateHTMLTemplate() {
	path := "./templates/"

	err := os.MkdirAll(filepath.Dir(path), 0750)
	if err != nil {
		log.Fatal("Could not create filepath")
	}

	filename := fmt.Sprintf("%sginger.html", path)

	file, err := os.Create(filename)
	if err != nil {
		log.Fatal("Could not generate HTML template")
	}

	defer file.Close()

	file.WriteString(utils.NoIndent("<!DOCTYPE html>"))
	file.WriteString(utils.NoIndent("<html>"))

	//#region html head
	file.WriteString(utils.WithIndent("<head>", 1))
	file.WriteString(utils.WithIndent("<link rel=\"stylesheet\" href=\"{{ .stylesPath }}\">", 2))
	file.WriteString(utils.WithIndent("<script type=\"text/javascript\" src=\"{{ .scriptsPath }}/ginger.js\"></script>", 2))
	file.WriteString(utils.WithIndent("<script type=\"text/javascript\" src=\"{{ .scriptsPath }}/ginger.elements.js\"></script>", 2))
	file.WriteString(utils.WithIndent("</head>", 1))
	//#endregion

	//#region html body
	file.WriteString(utils.WithIndent("<body>", 1))
	file.WriteString(utils.WithIndent("<h1>{{ .title }}</h1>", 2))
	file.WriteString(utils.WithIndent("</body>", 1))
	//#endregion

	file.WriteString(utils.NoIndent("</html>"))
}