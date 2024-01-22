package types

import (
	"github.com/tpyle/testify/lib/types/result"
)

type Setup interface {
	Validate() error
	Setup() (interface{}, error)
	Teardown() error
}

type ReadyCheck interface {
	Validate() error
	WaitForReady(context interface{}) error
}

type Runner interface {
	Validate() error
	Run(context interface{}) (result.ResultGroup, error)
}
