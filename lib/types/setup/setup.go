package setup

import "io"

type Setup interface {
	Validate() error
	Setup(context map[string]string, logFile io.Writer) (interface{}, error)
	Teardown(logFile io.Writer) error
}
