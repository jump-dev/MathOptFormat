package main

import (
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"github.com/qri-io/jsonschema"
)

//go:generate go run generate-static-schema.go ..\\..\\mof.schema.json
func main() {
	validateFlag := flag.Bool(
		"validate", true,
		"Validate the input JSON file against the MathOptFormat schema")
	flag.Parse()

	tail := flag.Args()
	if len(tail) != 1 {
		fmt.Println("Invalid arguments")
		flag.PrintDefaults()
		return
	}
	filename := tail[0]

	if *validateFlag {
		if err := validateFile(filename); err != nil {
			fmt.Printf("%s is not a valid MOF file\nThe error is:\n", filename)
			fmt.Println(err)
		} else {
			fmt.Printf("Success! %s conforms to the MathOptFormat schema", filename)
		}
	}
	return
}

func validateFile(filename string) error {
	schemaBytes, err := base64.StdEncoding.DecodeString(jsonSchema64)
	if err != nil {
		fmt.Println("Unable to decode JSON schema")
		return err
	}
	rs := &jsonschema.RootSchema{}
	if err := json.Unmarshal(schemaBytes, rs); err != nil {
		fmt.Println("Unable to unmarshall schema")
		return err
	}
	modelData, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Printf("Unable to read %s\n", filename)
		return err
	}
	if errs, _ := rs.ValidateBytes(modelData); len(errs) > 0 {
		fmt.Printf("Error validating file")
		for _, err = range errs {
			fmt.Println(err.Error())
		}
		return errs[0]
	}
	return nil
}
