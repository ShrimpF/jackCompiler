package main

import (
	"os"

	"github.com/ShrimpF/jackCompiler/compiler"
	"github.com/ShrimpF/jackCompiler/tokenizer"
)

func main() {
	filePath := os.Args[1]
	t := tokenizer.NewTokenizer(filePath)
	file, err := os.Create("sample.xml")
	if err != nil {
		panic(err)
	}
	compiler := compiler.New(t, file)
	compiler.CompileClass()
}
