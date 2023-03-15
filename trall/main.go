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
	for _, type_ := range []gobench.SubBenchType{gobench.GoKerNonBlocking, gobench.GoKerBlocking} {
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
		println(bugs[len(bugs)-1].ID)
		for _, b := range bugs {
			result := s.GetResult(b.ID)
			//fmt.Println(b.ID, "logs ->", filepath.Join(result.OutputDir, "full.log"))
			pathToTrace := utils.PathToTrace(b.Type.String(), b.SubType, b.SubSubType, b.ID)
			if result.IsPositive() {
				fmt.Printf("OK, we reproduced %s in GoBench (GoKer)\n", b.ID)
				trason.Trason(pathToTrace)
			} else {
				fmt.Printf("Sorry, we failed to reproduce %s in GoBench (GoKer)\n", b.ID)
				if err := os.Remove(pathToTrace); err != nil {
					panic(err)
				}
			}
		}
		println("Done")
	}
}
