package main

import (
	"flag"
	"fmt"
	"os"
	"zreader7/zinterpreter"
)

func printHelp() {
	fmt.Println("2025 : See my blog https://jeandi7.github.io/jeandi7blog/")
	fmt.Println()
	fmt.Println("Usage: zreader [options]")
	fmt.Println("Options:")
	flag.PrintDefaults()
	os.Exit(0)
}

func writeOutFile(content string, outname string) {
	filename := outname + ".puml"

	file, err := os.Create(filename)
	if err != nil {
		fmt.Println("Erreur lors de la création du fichier:", err)
		return
	}
	defer file.Close() // Assurer la fermeture du fichier à la fin du programme

	_, err = file.WriteString(content)
	if err != nil {
		fmt.Println("Erreur lors de l'écriture dans le fichier:", err)
		return
	}
}

func writeOutFileForArchi(content string, outname string) {
	filename := outname + ".xml"

	file, err := os.Create(filename)
	if err != nil {
		fmt.Println("Erreur lors de la création du fichier:", err)
		return
	}
	defer file.Close() // Assurer la fermeture du fichier à la fin du programme

	_, err = file.WriteString(content)
	if err != nil {
		fmt.Println("Erreur lors de l'écriture dans le fichier:", err)
		return
	}
}

func main() {

	// examples :
	// input := `definition monsujet { }`
	// input := `definition monsujet { relation marelation1: monsujet2 }`
	// input := `definition monsujet { relation marelation1: monsujet2 | monsujet3 }`
	// input := `definition monsujet { relation marelation1: monsujet2 | monsujet3  relation marelation2: monsujet2 }`
	// input := `definition monsujet { } definition monsujet2 { } definition maressource { relation marelation: monsujet | monsujet2   }`
	// input := `definition monsujet { } definition monsujet2 { } definition maressource { relation marelation: monsujet | monsujet2  relation mr2: monsujet | msj3  }`

	var input string = ""
	var schema string = ""
	var fschema string = ""
	var out string = ""
	var showHelp bool

	flag.StringVar(&schema, "schema", "", "Read schema")
	flag.StringVar(&fschema, "fschema", "", "Read schema file")
	flag.StringVar(&out, "out", "out", "Archimate plantUML generated file name")
	flag.BoolVar(&showHelp, "help", false, "Show help message")
	flag.Parse()

	if (schema == "" && fschema == "") || (schema != "" && fschema != "") {
		fmt.Println("you must provide either -schema or -fschema, but not both.")
		printHelp()
		return
	}

	if schema != "" {
		input = schema
	}

	if fschema != "" {
		fileContent, err := os.ReadFile(fschema)
		if err != nil {
			fmt.Println("Erreur lors de la lecture du fichier : ", err)
		}
		input = string(fileContent)
	}

	if showHelp {
		printHelp()
		return
	}

	lexer := zinterpreter.NewLexer(input)
	lexer.NextToken()
	zschema, err := lexer.ReadZSchema()

	if err != nil {
		fmt.Println("syntax error:", err)
	} else {
		// fmt.Println("parsed schema OK:", zschema)
		fmt.Println("parsed schema is done.")
	}

	mydraw := zinterpreter.PlantUMLArchimateSchema{Zdefs: zschema}
	archimatePlantUml := mydraw.Generate(out)
	writeOutFile(archimatePlantUml, out)
	fmt.Println("Generating " + out + ".puml is done.")

	archimateArchi := mydraw.GenerateForArchimate(out)
	writeOutFileForArchi(archimateArchi, out)
	fmt.Println("Generating " + out + ".xml is done.")

}
