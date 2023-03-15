package main

import (
	"github.com/timmyyuan/gobench"
	"strings"
	"time"
)

func bugString(subSubType string) string {
	switch subSubType {
	case "Data race":
		return "DATA RACE"
	}
	return ""
}

func main() {
	allBugs := gobench.GoBenchBugSet.ListByTypes(gobench.GoKerNonBlocking)
	for _, b := range allBugs {
		tmp := bugString(b.SubSubType)
		var positiveCheckFunc func(r *gobench.SingleRunResult) bool
		if tmp != "" {
			positiveCheckFunc = func(r *gobench.SingleRunResult) bool {
				return strings.Contains(string(r.Logs), tmp)
			}
		}
		isPositive := false
		for !isPositive {
			s := gobench.NewSuite(gobench.SuiteConfig{
				ExecEnvConfig: gobench.ExecEnvConfig{
					Count:             1,
					Timeout:           5 * time.Second,
					Repeat:            1,
					PositiveCheckFunc: positiveCheckFunc,
				},
				Type:   gobench.GoKerNonBlocking,
				BugIDs: []string{b.ID},
			})
			s.Run()

			result := s.GetResult(b.ID)

			//if result.IsPositive() {
			//	fmt.Printf("OK, we reproduced %s in GoBench (GoKer)\n", b.ID)
			//} else {
			//	fmt.Printf("Sorry, we failed to reproduce %s in GoBench (GoKer)\n", b.ID)
			//}
			isPositive = result.IsPositive()
		}
	}
}
