package model

type Status string

const (
	StatusPass Status = "pass"
	StatusWarn Status = "warn"
	StatusFail Status = "fail"
	StatusInfo Status = "info"
)

func (s Status) IsProblem() bool {
	return s == StatusWarn || s == StatusFail
}
