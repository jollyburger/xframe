package log

import (
	"errors"
)

type Hook interface {
	Levels() []string
	Fire(string, []byte) error
}

type Hooks map[string][]Hook

func (hks Hooks) AddHook(hook Hook) error {
	for _, level := range hook.Levels() {
		hks[level] = append(hks[level], hook)
	}
	return nil
}

func (hks Hooks) Fire(level string, msg []byte) error {
	if len(hks) == 0 {
		return errors.New(level)
	}
	for _, hook := range hks[level] {
		if err := hook.Fire(level, msg); err != nil {
			return err
		}
	}
	return nil
}
