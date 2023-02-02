package visca

import "testing"

var _ Logger = testLogger{}

type testLogger struct {
	*testing.T
}

func (t testLogger) Debugf(format string, a ...interface{}) {
	t.T.Logf("[debug] "+format, a...)
}

func (t testLogger) Infof(format string, a ...interface{}) {
	t.T.Logf("[info] "+format, a...)
}

func (t testLogger) Warnf(format string, a ...interface{}) {
	t.T.Logf("[warn] "+format, a...)
}
