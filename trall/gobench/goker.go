package main

import (
	"fmt"
	"github.com/timmyyuan/gobench"
	"os"
	"strings"
	"time"
	"trall/utils"
)

func positiveCheckFunc(subSubType string) func(r *gobench.SingleRunResult) bool {
	switch subSubType {
	case "Data race":
		return func(r *gobench.SingleRunResult) bool {
			return strings.Contains(string(r.Logs), "DATA RACE")
		}
	default:
		return nil
	}
}

func main() {
	for _, subBenchType := range []gobench.SubBenchType{gobench.GoKerNonBlocking, gobench.GoKerBlocking} {
		allBugs := gobench.GoBenchBugSet.ListByTypes(subBenchType)
		for _, b := range allBugs {
			s := gobench.NewSuite(gobench.SuiteConfig{
				ExecEnvConfig: gobench.ExecEnvConfig{
					Count:             1,
					Timeout:           5 * time.Second,
					Repeat:            1,
					PositiveCheckFunc: positiveCheckFunc(b.SubSubType),
				},
				Type:   subBenchType,
				BugIDs: []string{b.ID},
			})
			s.Run()

			result := s.GetResult(b.ID)

			if result.IsPositive() {
				fmt.Printf("OK, we reproduced %s in GoBench (GoKer)\n", b.ID)
			} else {
				fmt.Printf("Sorry, we failed to reproduce %s in GoBench (GoKer)\n", b.ID)
				if err := os.Remove(utils.PathToTrace(b.SubType, b.SubSubType, b.ID) + ".trace"); err != nil {
					panic(err)
				}
			}
		}
	}
}
