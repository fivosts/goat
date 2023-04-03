package main

import (
	"fmt"
	"github.com/timmyyuan/gobench"
	"strings"
	"time"
)

func main() {
	for _, type_ := range []gobench.SubBenchType{gobench.GoKerNonBlocking, gobench.GoKerBlocking, gobench.GoRealNonBlocking, gobench.GoRealBlocking} {
		s := gobench.NewSuite(gobench.SuiteConfig{
			ExecEnvConfig: gobench.ExecEnvConfig{
				Count:   1,
				Timeout: 5 * time.Second,
				Repeat:  1,
				PositiveCheckFunc: func(r *gobench.SingleRunResult) bool {
					if r.Bug.SubType == "Resource Deadlock" {
						return strings.Contains(string(r.Logs), "semacquire")
					}
					switch r.Bug.SubSubType {
					case "Data race":
						return strings.Contains(string(r.Logs), "DATA RACE")
					default:
						return r.ExitCode != 0
					}
				},
			},
			Type: type_,
		})
		s.Run()
		fmt.Println("Finished", type_)
	}
}
