// +build ignore

package main

import (
	"encoding/base64"
	"io/ioutil"
	"os"
)

func main() {
	filename := os.Args[1]
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	base64str := base64.StdEncoding.EncodeToString(bytes)
	packageString := "// Code generated by go generate; DO NOT EDIT.\n" +
		"package main\n" +
		"\n" +
		"var jsonSchema64 = \"" + base64str + "\""
	if err := ioutil.WriteFile("staticSchema.go", []byte(packageString), os.ModePerm); err != nil {
		panic(err)
	}
}
