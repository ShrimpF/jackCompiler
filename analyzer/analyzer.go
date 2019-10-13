package analyzer

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/ShrimpF/jackCompiler/compiler"
	"github.com/ShrimpF/jackcompiler/tokenizer"
)

// Analyze -- start compiling jack file
func Analyze(filePath string) {
	var t *tokenizer.Base
	var c *compiler.Base

	fileInfo, err := os.Stat(filePath)
	checkErr(err)

	if fileInfo.IsDir() {
		outputName := fileInfo.Name() + ".xml"
		output, err := os.Create(outputName)
		checkErr(err)

		jackFilePaths := getJackFilesFromPath(filePath)

		for _, jackFilePath := range jackFilePaths {
			t = tokenizer.NewTokenizer(jackFilePath)
			c = compiler.New(t, output)
			c.Start()
		}
	} else {
		outputName, err := converFileExt(fileInfo.Name(), ".jack", ".xml")
		checkErr(err)

		output, err := os.Create(outputName)
		checkErr(err)

		t = tokenizer.NewTokenizer(filePath)
		c = compiler.New(t, output)
		c.Start()
	}

	fmt.Println("compile end successfully")
}

func getJackFilesFromPath(path string) []string {
	dir := filepath.Dir(path)
	var jackFilePaths []string

	files, err := ioutil.ReadDir(path)
	checkErr(err)

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".jack" {
			jackFilePaths = append(jackFilePaths, dir+"/"+file.Name())
		}
	}

	return jackFilePaths
}

func converFileExt(fileName, prevExt, nextExt string) (string, error) {
	if filepath.Ext(fileName) == ".jack" {
		return strings.Replace(fileName, prevExt, nextExt, 1), nil
	}
	return fileName, errors.New("invalid file extension")
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
