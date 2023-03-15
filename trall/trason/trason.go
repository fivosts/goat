package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"trall/trace"
)

func main() {
	filename, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Print("error opening file\n")
		return
	}
	reader := bufio.NewReader(filename)
	parsed, err := trace.Parse(reader, "")
	if err != nil {
		fmt.Print("Error parsing trace: ", err, "\n")
		return
	}
	w := os.Stdout
	_, err = w.WriteString(`{"events": [`)
	if err != nil {
		return
	}
	for i, event := range parsed.Events {
		if i > 0 {
			_, err = w.WriteString(",")
			if err != nil {
				return
			}
		}
		marshaled, err := json.Marshal(event)
		_, err = w.WriteString(string(marshaled))
		if err != nil {
			return
		}
	}
	_, err = w.WriteString(`]}`)
	if err != nil {
		return
	}
}
