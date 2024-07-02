package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/NickBlakW/ginger/generators"
	"github.com/NickBlakW/ginger/types"
	"github.com/gin-gonic/gin"
)

// func main() {
// 	router := gin.Default()

// 	config := types.Config{
// 		LocalApiPath: "./web",
// 	}

// 	UseGingerUi(config, router)

// 	router.Run()
// }

// #region GinGer UI
func UseGingerUi(config types.Config, ginEngine *gin.Engine) {
	path := fmt.Sprintf("%s/", config.LocalApiPath)

	err := os.MkdirAll(filepath.Dir(fmt.Sprintf("%sscripts/", path)), 0750)
	if err != nil {
		log.Fatal("Could not create filepath")
	}

	err = os.MkdirAll(filepath.Dir(fmt.Sprintf("%sstyles/", path)), 0750)
	if err != nil {
		log.Fatal("Could not create filepath")
	}

	fmt.Println("Generating HTML Template...")
	generators.GenerateHTMLTemplate()
	fmt.Println("Done!")

	fmt.Println("Generating CSS...")
	generators.GenerateUiStyles(fmt.Sprintf("%sstyles/", path))
	fmt.Println("Done!")

	fmt.Println("Generating API script...")
	generators.GenerateLocalApiScripts(fmt.Sprintf("%sscripts/", path), ginEngine)
	fmt.Println("Done!")
	
	fmt.Println("Generating element generator script...")
	generators.GenerateUiElementScript(fmt.Sprintf("%sscripts/", path), ginEngine)
	fmt.Println("Done!")
	
	basePath := strings.Split(path, ".")

	//#region HTTP config
	ginEngine.LoadHTMLGlob("templates/*")
	ginEngine.Static(basePath[1], config.LocalApiPath)

	ginEngine.GET("/ginger", func(ctx *gin.Context) {
		base := strings.Split(path, ".")

		scriptsPath := fmt.Sprintf("%sscripts", base[1])
		stylesPath := fmt.Sprintf("%sstyles/ginger.css", base[1])

		templateData := gin.H{
			"title": "GinGer API Viewer",
			"scriptsPath": scriptsPath,
			"stylesPath": stylesPath,
		}

		ctx.HTML(http.StatusOK, "ginger.html", templateData)
	})
	//#endregion
}
// #endregion
