package utils

import (
	"os"
	"path/filepath"
	"strings"
)

func PathToTrace(type_ string, subType string, subSubType string, id string, success bool) string {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	var result string
	if success {
		result = "positive"
	} else {
		result = "negative"
	}
	traceDir := filepath.Join(
		filepath.Dir(wd),
		"results",
		result,
		"gobench",
		type_,
		strings.ReplaceAll(subType, " ", "_"),
		strings.ReplaceAll(subSubType, " ", "_"),
	)
	if err = os.MkdirAll(traceDir, os.ModePerm); err != nil {
		panic(err)
	}
	file, err := os.CreateTemp(traceDir, id+"_*.trace")
	if err != nil {
		panic(err)
	}
	return file.Name()
}
