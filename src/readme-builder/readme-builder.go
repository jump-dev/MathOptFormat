/*
This file is used to rebuild the list of supported functions and sets for the
README from a provided schema. See `main()` for details.
*/
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

// Build the MathOptFormat documentation based on the provided schema.
//
// If the template is provided as the third argument, a find-and-replace will be
// conducted for [[[AUTOMATICALLY_GENERATED_SET_SUMMARY]]] and
// [[[AUTOMATICALLY_GENERATED_FUNCTION_SUMMARY]]].
//
// If no template is provided, then the set and function summaries will be
// written to stdout.
func main() {
	if len(os.Args) != 2 && len(os.Args) != 3 {
		fmt.Println("Usage: doc_build schema.json readme_template.md")
		return
	}
	setSummary, functionSummary, err := generateDocs(os.Args[1])
	if err != nil {
		fmt.Println(err)
		return
	}
	if len(os.Args) == 2 {
		fmt.Println(setSummary)
		fmt.Println()
		fmt.Println(functionSummary)
	} else {
		template_bytes, err := readFileToBytes(os.Args[2])
		if err != nil {
			fmt.Println(err)
			return
		}
		readme := strings.Replace(
			strings.Replace(
				string(template_bytes),
				"[[[AUTOMATICALLY_GENERATED_SET_SUMMARY]]]",
				setSummary, 1),
			"[[[AUTOMATICALLY_GENERATED_FUNCTION_SUMMARY]]]",
			functionSummary, 1)
		fmt.Println(readme)
	}
}

// Return `[]byte` of the file `filename`.
func readFileToBytes(filename string) ([]byte, error) {
	jsonFile, err := os.Open(filename)
	defer jsonFile.Close()
	if err != nil {
		return []byte{}, errors.New(
			fmt.Sprintf("unable to open file: %s", filename))
	}
	bytes, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return []byte{}, errors.New(
			fmt.Sprintf("unable to read file: %s", filename))
	}
	return bytes, nil
}

type Set struct {
	Head        string
	Description string
	Example     string
}

func processOneOf(data map[string]interface{}) []Set {
	var set Set
	description, ok := data["description"]
	if !ok {
		description = ""
	}
	set.Description = strings.ReplaceAll(description.(string), "|", "\\|")

	example, ok := data["example"]
	if !ok {
		example = ""
	}
	set.Example = strings.ReplaceAll(example.(string), "|", "\\|")

	properties := data["properties"].(map[string]interface{})
	head := properties["head"].(map[string]interface{})
	if val, ok := head["const"]; ok {
		set.Head = val.(string)
		return []Set{set}
	} else if val, ok := head["enum"]; ok {
		sets := []Set{}
		for _, name := range val.([]interface{}) {
			sets = append(sets, Set{name.(string), set.Description, set.Example})
		}
		return sets
	}
	return []Set{}
}

// Generate the Markdown template for the functions and sets based on the JSON
// schema.
func generateDocs(jsonSchemaFilename string) (string, string, error) {
	schema_bytes, err := readFileToBytes(jsonSchemaFilename)
	if err != nil {
		return "", "", err
	}

	var data map[string]interface{}
	if err := json.Unmarshal(schema_bytes, &data); err != nil {
		fmt.Println(err)
		return "", "", errors.New(
			fmt.Sprintf("unable to parse JSON schema: %s", jsonSchemaFilename))
	}

	setSummary := strings.Join([]string{
		processStrings("Scalar Sets", data, "scalar_sets"),
		"",
		processStrings("Vector Sets", data, "vector_sets")}, "\n")

	functionSummary := strings.Join([]string{
		processStrings("Scalar Functions", data, "scalar_functions"),
		"",
		processStrings("Vector Functions", data, "vector_functions")}, "\n")

	nlSummary := processNonlinear(data)
	return setSummary, functionSummary + nlSummary, nil
}

func processStrings(title string, data map[string]interface{}, key string) string {
	setStrings := []string{
		fmt.Sprintf("#### %s", title),
		"",
		"| Name | Description | Example |",
		"| ---- | ----------- | ------- |"}
	definitions := data["definitions"].(map[string]interface{})
	sets := definitions[key].(map[string]interface{})
	for _, setData := range sets["oneOf"].([]interface{}) {
		for _, set := range processOneOf(setData.(map[string]interface{})) {
			setStrings = append(
				setStrings,
				fmt.Sprintf("| `\"%s\"` | %s | %s |", set.Head, set.Description, set.Example))
		}
	}
	return strings.Join(setStrings, "\n")
}

func processNonlinear(data map[string]interface{}) string {
	definitions := data["definitions"].(map[string]interface{})
	nonlinearTerm := definitions["NonlinearTerm"].(map[string]interface{})
	functionStrings := []string{
		"##### Functions",
		"",
		"| Name | Arity |",
		"| ---- | ----- |"}

	leafStrings := []string{
		"#### Nonlinear functions",
		"",
		"##### Leaf nodes",
		"| Name | Description | Example |",
		"| ---- | ----------- | ------- |"}

	for _, setData := range nonlinearTerm["oneOf"].([]interface{}) {
		object := setData.(map[string]interface{})
		objects := processOneOf(object)

		switch object["description"] {
		case "Unary operators":
			for _, f := range objects {
				functionStrings = append(functionStrings,
					fmt.Sprintf("| `\"%s\"` | Unary |", f.Head))
			}
		case "Binary operators":
			for _, f := range objects {
				functionStrings = append(functionStrings,
					fmt.Sprintf("| `\"%s\"` | Binary |", f.Head))
			}
		case "N-ary operators":
			for _, f := range objects {
				functionStrings = append(functionStrings,
					fmt.Sprintf("| `\"%s\"` | N-ary |", f.Head))
			}
		default:
			if len(objects) == 1 {
				leafStrings = append(leafStrings,
					fmt.Sprintf("| `\"%s\"` | %s | %s |", objects[0].Head, objects[0].Description, objects[0].Example))
			} else {
				fmt.Println("Unsupported object: %s", object)
			}
		}
	}
	return "\n" + strings.Join(leafStrings, "\n") + "\n\n" +
		strings.Join(functionStrings, "\n")
}
