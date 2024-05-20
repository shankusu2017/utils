package log

import (
	"testing"
)

func TestLog(t *testing.T) {
	SetFile("./test.log")

	Info("utils.log test!")
	Info("[hello world!%s]", "args")
}
