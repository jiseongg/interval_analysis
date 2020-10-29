package main

import (
	"flag"
	"fmt"
	"interval/internal/analyzer"

	"github.com/llir/llvm/asm"
)

func main() {
	flag.Parse()
	args := flag.Args()
	filename := args[0]
	m, err := asm.ParseFile(filename)
	if err != nil {
		panic(err)
	}
	md := analyzer.NewModule(m)
	for _, cfg := range md.Cfgs {
		fmt.Printf("Analysis of %s begins...\n", cfg.GetFid())
		tbl := analyzer.Analyze(cfg)
		fmt.Println()
		fmt.Println("Analysis Results for:", cfg.GetFid())
		fmt.Println(tbl.String())
		fmt.Println()
	}
}
