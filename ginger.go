package main

import (
	"fmt"

	"github.com/NickBlakW/ginger/generators"
	"github.com/NickBlakW/ginger/types"
)

func main() {
	config := types.Config{
		LocalApiPath: "./web",
	}

	fmt.Println("Generating API script...")

	generators.GenerateLocalApiScripts(config)

	fmt.Println("Done!")
}
