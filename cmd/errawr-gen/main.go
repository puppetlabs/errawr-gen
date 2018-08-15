package main

import (
	"flag"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/puppetlabs/errawr-gen/doc"
	"github.com/puppetlabs/errawr-gen/golang"
	"github.com/xeipuuv/gojsonschema"
	yaml "gopkg.in/yaml.v2"
)

type Generator interface {
	Generate(pkg string, document *doc.Document, output io.Writer) error
}

type Language string

var (
	LanguageGo Language = "go"
)

var (
	Package        string
	OutputPath     string
	OutputLanguage Language
	InputPath      string
)

func init() {
	flag.StringVar(&Package, "package", os.Getenv("GOPACKAGE"), "the package to write")
	flag.StringVar(&OutputPath, "output-path", "-", "the path to write output to")
	flag.StringVar((*string)(&OutputLanguage), "output-language", string(LanguageGo), "the language to write errors for")
	flag.StringVar(&InputPath, "input-path", "-", "the path to read input from")
}

func main() {
	flag.Parse()

	if len(Package) == 0 {
		log.Fatalf("Package name could not be determined; specify one.")
	}

	var generator Generator

	switch OutputLanguage {
	case LanguageGo:
		generator = golang.NewGenerator()
	default:
		log.Fatalf("Language %q is not supported.", OutputLanguage)
	}

	var input, output *os.File
	var err error

	if len(InputPath) > 0 && InputPath != "-" {
		input, err = os.Open(InputPath)
		if err != nil {
			log.Fatalf("Could not open input file: %+v", err)
		}
		defer input.Close()
	} else {
		input = os.Stdin
	}

	y, err := ioutil.ReadAll(input)
	if err != nil {
		log.Fatalf("Could not read file: %+v", err)
	}

	// Pull out the document version.
	var version doc.DocumentVersionFragment
	if err := yaml.Unmarshal(y, &version); err != nil {
		log.Fatalf("Could not read version from YAML: %+v", err)
	}

	if version.Version != "1" {
		log.Fatalf(`Unexpected version %q; expected "1"`, version)
	}

	var document doc.Document
	if err := yaml.UnmarshalStrict(y, &document); err != nil {
		log.Fatalf("Could not parse YAML: %+v", err)
	}

	result, err := doc.Schema.Validate(gojsonschema.NewGoLoader(document))
	if err != nil {
		log.Fatalf("Could not generate YAML validation: %+v", err)
	} else if !result.Valid() {
		for _, err := range result.Errors() {
			log.Println(err)
		}

		log.Fatalf("Validation errors occurred.")
	}

	if len(OutputPath) > 0 && OutputPath != "-" {
		output, err = os.Create(OutputPath)
		if err != nil {
			log.Fatalf("Could not open output file: %+v", err)
		}
		defer output.Close()
	} else {
		output = os.Stdout
	}

	if err := generator.Generate(Package, &document, output); err != nil {
		log.Fatalf("Failed to generate Go file: %+v", err)
	}
}
