package common

import (
	"time"
)

// RawTestStep models single line from "go test -json" output.
type RawTestStep struct {
	Time    time.Time `json:"Time"`
	Action  string    `json:"Action"`
	Package string    `json:"Package"`
	Test    string    `json:"Test"`
	Output  string    `json:"Output"`
	Elapsed float32   `json:"Elapsed"`
}

// TestRun models full level 1 summary of a test run from "go test -json".
type TestRun struct {
	TestMap map[string][]TestResult `json:"-"`
	Rows    []TestResultRow         `json:"rows"`
}

type TestResultRow struct {
	TestResult TestResult `json:"json"`
}

// TestResult models result of a single test
type TestResult struct {
	// data that spans multiple tests - it's added at the test level because it will be used
	// by BigQuery tables and will need to be flattened
	CommitSha  string    `json:"commit_sha"`
	CommitDate time.Time `json:"commit_date"`
	JobRunDate time.Time `json:"job_run_date"`

	// test specific data
	Test    string `json:"test"`
	Package string `json:"package"`
	Output  []struct {
		Item string `json:"item"`
	} `json:"output"`
	Result  string  `json:"result"`
	Elapsed float32 `json:"elapsed"`
}
