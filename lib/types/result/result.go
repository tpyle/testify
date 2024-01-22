package result

import "time"

type TestStatus string

const (
	TestStatusPassed  TestStatus = "passed"
	TestStatusFailed  TestStatus = "failed"
	TestStatusErrored TestStatus = "errored"
)

type Result struct {
	Name      string        `json:"name"`
	Status    string        `json:"status"`
	Error     string        `json:"error"`
	StartTime time.Time     `json:"startTime"`
	EndTime   time.Time     `json:"endTime"`
	Duration  time.Duration `json:"duration"`
	Metadata  map[string]string
}

type ResultGroup struct {
	Name         string    `json:"name"`
	StartTime    time.Time `json:"startTime"`
	EndTime      time.Time `json:"endTime"`
	Duration     time.Time `json:"duration"`
	Results      []*Result
	ResultGroups []*ResultGroup
	Metadata     map[string]string
}
