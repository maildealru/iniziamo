package main

import (
	"log"
	"os"

	"github.com/maildealru/iniziamo/pkg/generator"
	"github.com/maildealru/iniziamo/pkg/parser"
	"github.com/maildealru/iniziamo/pkg/preprocessor"

	"github.com/hokaccha/go-prettyjson"
)

func main() {
	path := os.Getenv("GOFILE")
	if path == "" {
		log.Fatal("no file to parse")
	}

	p := parser.NewParser()
	info := p.ParseFile(path)

	s, _ := prettyjson.Marshal(info)
	//log.Printf("INFO PARSED: \n%s",s)

	pp := preprocessor.NewPreprocessor()
	pp.Preprocess(info)

	s, _ = prettyjson.Marshal(info)
	log.Printf("INFO PREPROCESSED: \n%s", s)

	g := generator.NewGenerator()
	g.WriteFile(info)
}
