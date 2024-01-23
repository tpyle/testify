package types

import (
	"io"

	"github.com/tpyle/testify/lib/types/result"
)

type Setup interface {
	Validate() error
	Setup(context map[string]string, logFile *io.Writer) (interface{}, error)
	Teardown() error
}

type ReadyCheck interface {
	Validate() error
	WaitForReady(context interface{}, logFile *io.Writer) error
}

type Runner interface {
	Validate() error
	Run(context interface{}, logFile *io.Writer) (result.ResultGroup, error)
}
