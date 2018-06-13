package bridge

type logger interface {
	Debugf(fmt string, args ...interface{})
	Infof(fmt string, args ...interface{})
	Errorf(fmt string, args ...interface{})
}
