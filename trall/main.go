package main

import (
	"fmt"
	"github.com/timmyyuan/gobench"
	"os"
	"strings"
	"time"
	"trall/trason"
	"trall/utils"
)

func main() {
	for _, type_ := range []gobench.SubBenchType{gobench.GoKerNonBlocking, gobench.GoKerBlocking, gobench.GoRealNonBlocking, gobench.GoRealBlocking} {
		s := gobench.NewSuite(gobench.SuiteConfig{
			ExecEnvConfig: gobench.ExecEnvConfig{
				Count:   1,
				Timeout: 5 * time.Second,
				Repeat:  1,
				PositiveCheckFunc: func(r *gobench.SingleRunResult) bool {
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
		bugs := s.BugSet.ListByTypes(type_)
		for _, b := range bugs {
			pathToTrace := utils.PathToTrace(b.Type.String(), b.SubType, b.SubSubType, b.ID)
			if s.GetResult(b.ID).IsPositive() {
				fmt.Printf("OK, we reproduced %s in GoBench\n", b.ID)
				trason.Trason(pathToTrace)
			} else {
				fmt.Printf("Sorry, we failed to reproduce %s in GoBench\n", b.ID)
				if err := os.Remove(pathToTrace); err != nil {
					panic(err)
				}
			}
		}
		fmt.Println("Finished", type_)
	}
}
