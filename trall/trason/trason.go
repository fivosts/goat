package trason

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"trall/trace"
)

func Trason(paths ...string) {
	pathToTrace := paths[0]
	var pathToJSON string
	if len(paths) > 1 {
		pathToJSON = paths[1]
	} else {
		pathToJSON = strings.Replace(pathToTrace, ".trace", ".json", 1)
	}
	traceFile, err := os.Open(pathToTrace)
	if err != nil {
		panic(err)
	}
	fi, err := traceFile.Stat()
	if err != nil {
		panic(err)
	}
	if fi.Size() > 2<<15 {
		println("Trace file is too large, skipping...")
		return
	}
	parsed, err := trace.Parse(io.Reader(traceFile), "")
	if err != nil {
		if err.Error() == "trace is empty" {
			fmt.Println("Trace file is empty, deleting...")
			if err = os.Remove(pathToTrace); err != nil {
				panic(err)
			}
		} else {
			panic(err)
		}
	}
	if err := traceFile.Close(); err != nil {
		panic(err)
	}
	jsonFile, err := os.Create(pathToJSON)
	if err != nil {
		panic(err)
	}
	w := bufio.NewWriter(jsonFile)
	if _, err = w.WriteString(`{"events": [`); err != nil {
		panic(err)
	}
	for i, event := range parsed.Events {
		if i > 0 {
			if _, err = w.WriteString(","); err != nil {
				panic(err)
			}
		}
		marshaled, err := json.Marshal(event)
		if err != nil {
			panic(err)
		}
		if _, err = w.Write(marshaled); err != nil {
			panic(err)
		}
	}
	if _, err = w.WriteString(`]}`); err != nil {
		panic(err)
	}
	if err = w.Flush(); err != nil {
		panic(err)
	}
}
