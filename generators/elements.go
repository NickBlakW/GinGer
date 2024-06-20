package generators

import (
	"fmt"
	"log"
	"os"

	"github.com/NickBlakW/ginger/generators/utils"
)

func GenerateUiElementScript(path string) {
	filename := fmt.Sprintf("%sginger.elements.js", path)

	file, err := os.Create(filename)
	if err != nil {
		log.Fatal("Could not create file")
	}

	defer file.Close()

	file.WriteString(utils.NoIndent("window.addEventListener('load', () => {"))
	
	file.WriteString(utils.NoIndent("});"))
}