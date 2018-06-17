package bridge

type nilLogger struct {
}

func (nilLogger) Debugf(fmt string, args ...interface{}) {
}

func (nilLogger) Infof(fmt string, args ...interface{}) {
}

func (nilLogger) Errorf(fmt string, args ...interface{}) {
}
