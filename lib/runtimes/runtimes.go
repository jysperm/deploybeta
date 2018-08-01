package runtimes

import (
	"bytes"
	"errors"
)

var ErrInvalidRuntimeType = errors.New("invalid runtime type")

type RuntimeDecider = func(context *BuildContext) (bool, error)

type Runtime interface {
	Dockerfile() (*bytes.Buffer, error)
}

func DecideRuntime(context *BuildContext) (Runtime, error) {
	deciders := map[string]RuntimeDecider{
		"nodejs": DecideNodejsRuntime,
		"golang": DecideGolangRuntime,
	}

	for runtime, decider := range deciders {
		isThisRuntime, err := decider(context)

		if err != nil {
			return nil, err
		} else if isThisRuntime {
			return NewRuntime(runtime, context), nil
		}
	}

	return nil, ErrInvalidRuntimeType
}

func NewRuntime(runtimeType string, context *BuildContext) Runtime {
	if runtimeType == "nodejs" {
		return &NodejsRuntime{
			BuildContext: *context,
		}
	} else if runtimeType == "golang" {
		return &GolangRuntime{
			BuildContext: *context,
		}
	} else {
		panic(ErrInvalidRuntimeType)
	}
}
