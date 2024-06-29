package generators

import (
	"fmt"
	"log"
	"os"
)

const styleTemplate = `
body {
    background-color: aliceblue;
    display: flex;
    width: 100vw;
    justify-content: center;
    font-family: "Trebuchet MS", "Lucida Sans Unicode",
        "Lucida Grande", "Lucida Sans", Arial, sans-serif;
    flex-direction: column;
}

h1 {
    align-self: center;
    margin-left: auto;
    margin-right: auto;
}

button {
    width: 20vw;
}

pre {
    background-color: #333;
    color: #aaa;
    width: 80vw;
    min-height: 10%;
    border-radius: 2px;
}
`

func GenerateUiStyles(path string) {
	filename := fmt.Sprintf("%sginger.css", path)

	file, err := os.Create(filename)
	if err != nil {
		log.Fatal("Could not create file")
	}

	defer file.Close()

	file.Write([]byte(styleTemplate))
}