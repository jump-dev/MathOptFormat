package main

import (
	"io/ioutil"
	"testing"
)

func TestExamples(t *testing.T) {
	files, err := ioutil.ReadDir("../examples")
	if err != nil {
		t.Errorf("%s", err)
	}
	for _, file := range files {
		if err := ValidateFile("../examples/" + file.Name()); err != nil {
			t.Errorf("%s failed to validate", file.Name())
		}
	}
}

func TestSummary(t *testing.T) {
	summary, err := SummarizeSchema()
	if err != nil || len(summary) < 100 {
		t.Errorf("Failed to summarize schema: %s", err)
	}
}

func TestHelp(t *testing.T) {
	PrintHelp()
}
