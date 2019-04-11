package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/qri-io/jsonschema"
	"io/ioutil"
	"os"
	"strings"
)

//go:generate go run generate-static-schema.go ../mof.schema.json

func main() {
	switch os.Args[1] {
	case "validate":
		if len(os.Args) != 3 {
			fmt.Println("Invalid arguments to `mof validate`")
			PrintHelp()
		}
		filename := os.Args[2]
		if err := ValidateFile(filename); err != nil {
			fmt.Printf("%s is not a valid MOF file\nThe error is:\n", filename)
			fmt.Println(err)
		} else {
			fmt.Printf("Success! %s conforms to the MathOptFormat schema", filename)
		}
	case "summarize":
		summary, err := SummarizeSchema()
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(summary)
		}
	case "help":
		PrintHelp()
	default:
		fmt.Println("Invalid arguments to `mof`.")
		PrintHelp()
	}
}

func PrintHelp() {
	fmt.Println("mof [arg] [args...]")
	fmt.Println()
	fmt.Println("mof validate filename.json")
	fmt.Println("    Validate the file `filename.json` using the MathOptFormat schema")
	fmt.Println("mof summarize")
	fmt.Println("    Print a summary of the functions and sets supported by MathOptFormat")
}

func ValidateFile(filename string) error {
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

func SummarizeSchema() (string, error) {
	schemaBytes, err := base64.StdEncoding.DecodeString(jsonSchema64)
	if err != nil {
		fmt.Println("Unable to decode JSON schema")
		return "", err
	}
	var data map[string]interface{}
	if err := json.Unmarshal(schemaBytes, &data); err != nil {
		fmt.Println("Unable to unmarshall schema")
		return "", err
	}
	operators, leaves := summarizeNonlinear(data)
	summary := "## Sets\n\n" +
		"### Scalar Sets\n\n" +
		summarize(data, "scalar_sets") + "\n" +
		"### Vector Sets\n\n" +
		summarize(data, "vector_sets") + "\n" +
		"## Functions\n\n" +
		"### Scalar Functions\n\n" +
		summarize(data, "scalar_functions") + "\n" +
		"### Vector Functions\n\n" +
		summarize(data, "vector_functions") + "\n" +
		"### Nonlinear functions\n\n" +
		"#### Leaf nodes\n\n" +
		leaves + "\n" +
		"#### Operators\n\n" +
		operators
	return summary, nil
}

type Object struct {
	Head        string
	Description string
	Example     string
}

func oneOfToObject(data map[string]interface{}) []Object {
	val, ok := data["description"]
	var description string
	if ok {
		description = strings.ReplaceAll(val.(string), "|", "\\|")
	} else {
		description = ""
	}
	vals, ok := data["examples"]
	var example string
	if !ok {
		example = ""
	} else {
		example = vals.([]interface{})[0].(string)
	}
	properties := data["properties"].(map[string]interface{})
	head := properties["head"].(map[string]interface{})
	if val, ok := head["const"]; ok {
		return []Object{Object{
			Head: val.(string), Description: description, Example: example}}
	} else if val, ok := head["enum"]; ok {
		objects := []Object{}
		for _, name := range val.([]interface{}) {
			objects = append(objects, Object{
				Head: name.(string), Description: description, Example: example})
		}
		return objects
	}
	return []Object{}
}

func summarize(data map[string]interface{}, key string) string {
	summary := "| Name | Description | Example |\n" +
		"| ---- | ----------- | ------- |\n"
	definitions := data["definitions"].(map[string]interface{})
	keyData := definitions[key].(map[string]interface{})
	for _, oneOf := range keyData["oneOf"].([]interface{}) {
		for _, o := range oneOfToObject(oneOf.(map[string]interface{})) {
			summary = summary + fmt.Sprintf(
				"| `\"%s\"` | %s | %s |\n", o.Head, o.Description, o.Example)
		}
	}
	return summary
}

func summarizeNonlinear(data map[string]interface{}) (string, string) {
	definitions := data["definitions"].(map[string]interface{})
	nonlinearTerm := definitions["NonlinearTerm"].(map[string]interface{})
	operators := "| Name | Arity |\n" +
		"| ---- | ----- |\n"
	leaves := "| Name | Description | Example |\n" +
		"| ---- | ----------- | ------- |\n"
	for _, term := range nonlinearTerm["oneOf"].([]interface{}) {
		oneOf := term.(map[string]interface{})
		objects := oneOfToObject(oneOf)
		switch oneOf["description"] {
		case "Unary operators":
			for _, f := range objects {
				operators = operators + fmt.Sprintf(
					"| `\"%s\"` | Unary |\n", f.Head)
			}
		case "Binary operators":
			for _, f := range objects {
				operators = operators + fmt.Sprintf(
					"| `\"%s\"` | Binary |\n", f.Head)
			}
		case "N-ary operators":
			for _, f := range objects {
				operators = operators + fmt.Sprintf(
					"| `\"%s\"` | N-ary |\n", f.Head)
			}
		default:
			if len(objects) == 1 {
				leaves = leaves + fmt.Sprintf(
					"| `\"%s\"` | %s | %s |\n",
					objects[0].Head, objects[0].Description, objects[0].Example)
			} else {
				fmt.Printf("Unsupported object: %s\n", oneOf)
			}
		}
	}
	return operators, leaves
}
