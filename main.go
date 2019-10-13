package main

import (
	"flag"

	"github.com/ShrimpF/jackCompiler/analyzer"
)

func main() {
	flag.Parse()
	analyzer.Analyze(flag.Arg(0))
}
