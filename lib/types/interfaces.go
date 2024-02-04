package types

import (
	"io"

	"github.com/tpyle/testamint/lib/types/result"
)

type ReadyCheck interface {
	Validate() error
	WaitForReady(context interface{}, logFile *io.Writer) error
}

type Runner interface {
	Validate() error
	Run(context interface{}, logFile *io.Writer) (result.ResultGroup, error)
}
