package check

import (
	"encoding/json"
	"fmt"
	"io"
)

type Check interface {
	Validate(context interface{}, logFile io.Writer) error
	WaitForReady(context interface{}, logFile io.Writer) error
}

type CheckType string

const (
	HttpGet   CheckType = "httpGet"
	LogLine   CheckType = "logLine"
	TcpSocket CheckType = "tcpSocket"
)

func UnmarshalCheck(data []byte) (Check, error) {
	type auxC struct {
		CheckType CheckType `json:"type"`
	}

	var aux auxC

	if err := json.Unmarshal(data, &aux); err != nil {
		return nil, err
	}

	var check Check
	switch aux.CheckType {
	case HttpGet:
		var h HttpGetReadyCheck
		check = &h
	default:
		return nil, fmt.Errorf("unknown check type: %s", aux.CheckType)
	}

	err := json.Unmarshal(data, check)
	if err != nil {
		return nil, err
	}

	return check, nil
}
