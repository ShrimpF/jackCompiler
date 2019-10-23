package main

import (
	"os"

	"github.com/ShrimpF/jackCompiler/compiler"
)

func main() {
	path := "./test.jack"
	file, err := os.Create("output.vm")
	if err != nil {
		panic(err)
	}

	c := compiler.New(path, file)
	c.Start()
}
