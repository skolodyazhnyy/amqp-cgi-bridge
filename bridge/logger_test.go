package bridge

type nilLogger struct {
}

func (nilLogger) Debug(msg string, args map[string]interface{}) {
}

func (nilLogger) Debugf(fmt string, args ...interface{}) {
}

func (nilLogger) Infof(fmt string, args ...interface{}) {
}

func (nilLogger) Error(msg string, args map[string]interface{}) {
}

func (nilLogger) Errorf(fmt string, args ...interface{}) {
}
