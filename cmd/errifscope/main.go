package main

import (
	"golang.org/x/tools/go/analysis/unitchecker"

	"github.com/otakakot/errifscope"
)

func main() { unitchecker.Main(errifscope.Analyzer) }
