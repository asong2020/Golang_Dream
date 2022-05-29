package main

import (
	"asong.cloud/Golang_Dream/code_demo/custom_linter/firstparamcontext"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(firstparamcontext.Analyzer)
}
