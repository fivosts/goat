package utils

import (
	"os"
	"path/filepath"
	"strings"
)

func PathToTrace(subType string, subSubType string, id string) string {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	traceDir := filepath.Join(filepath.Dir(wd), "results", "gobench", strings.ReplaceAll(subType, " ", "_"), strings.ReplaceAll(subSubType, " ", "_"))
	if err = os.MkdirAll(traceDir, os.ModePerm); err != nil {
		panic(err)
	}
	return filepath.Join(traceDir, id)
}
