package main

import (
	"errors"
	"os"
	"testing"

	"github.com/NickBlakW/ginger/types"
	"github.com/gin-gonic/gin"
)

func TestGingerGeneratesFiles(t *testing.T) {
	testConfig := types.Config{
		LocalApiPath: "./web",
	}
	
	UseGingerUi(testConfig, gin.Default())

	if _, err := os.Stat(testConfig.LocalApiPath); err == nil {
		t.Failed()
	} else if errors.Is(err, os.ErrNotExist) {
		t.Failed()
	}
}