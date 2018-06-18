package bridge

type logger interface {
	Debug(msg string, rec map[string]interface{})
	Debugf(fmt string, args ...interface{})
	Infof(fmt string, args ...interface{})
	Error(msg string, rec map[string]interface{})
	Errorf(fmt string, args ...interface{})
}
