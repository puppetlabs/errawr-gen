package main

import (
	"flag"
	"log"
	"os"

	"github.com/puppetlabs/errawr-gen/generator"
)

var (
	conf generator.Config
)

func init() {
	flag.StringVar(&conf.Package, "package", os.Getenv("GOPACKAGE"), "the package to write")
	flag.StringVar(&conf.OutputPath, "output-path", "-", "the path to write output to")
	flag.StringVar((*string)(&conf.OutputLanguage), "output-language", string(generator.LanguageGo), "the language to write errors for")
	flag.StringVar(&conf.InputPath, "input-path", "-", "the path to read input from")
}

func main() {
	flag.Parse()

	if err := generator.Generate(conf); err != nil {
		log.Fatalln(err)
	}
}
