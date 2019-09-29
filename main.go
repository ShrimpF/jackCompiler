package main

import (
	"os"

	"github.com/ShrimpF/jackCompiler/tokenizer"
)

func main() {
	filePath := os.Args[1]
	t := tokenizer.NewTokenizer(filePath)
	t.WriteXML()
}
