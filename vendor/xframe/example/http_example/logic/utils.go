package logic

import (
	"encoding/json"
)

type Result struct {
	Action  string
	RetCode int
	Message string
	Data    string
}

func FormatResponse(action string, ret int, reason string, i string) (buf []byte, err error) {
	out := &Result{action, ret, reason, i}
	buf, err = json.Marshal(out)
	if err != nil {
		return
	}
	return
}
