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

func Trason(filepath string) {
	traceFile, err := os.Open(filepath)
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
			if err = os.Remove(filepath); err != nil {
				panic(err)
			}
		}
	}
	if err := traceFile.Close(); err != nil {
		panic(err)
	}
	jsonFile, err := os.Create(strings.Replace(filepath, ".trace", ".json", 1))
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
