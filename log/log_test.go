package log

import (
	"testing"
)

func TestLog(t *testing.T) {
	//SetFile("./test.log")

	Info("utils.log test!")
	Debug("log.test.debug.level")
	Info("[hello world!%s]", "args")
	Error("log.test.error.level")
}
