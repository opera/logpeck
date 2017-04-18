package logpeck

import (
	"errors"
)

// Field format:
//   "Value": "$num"
func (p *PeckField) Check() error {
	if len(p.Value) < 2 || p.Value[0] != '$' {
		return errors.New("Value format error, start with '$': " + p.Value + "\n")
	}

	for _, c := range p.Value[1:] {
		if c > '9' || c < '0' {
			return errors.New("Value format error: " + p.Value + "\n")
		}
	}
	return nil
}

func (p *PeckTaskConfig) Check() error {
	errMsg := ""
	for _, f := range p.Fields {
		err := f.Check()
		if err != nil {
			errMsg += err.Error()
		}
	}
	if len(errMsg) > 0 {
		return errors.New(errMsg)
	}
	return nil
}
