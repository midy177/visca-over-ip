package visca

type Logger interface {
	Debugf(format string, a ...interface{})
	Infof(format string, a ...interface{})
	Warnf(format string, a ...interface{})
	Errorf(format string, a ...interface{})
}

func (c *Camera) debugf(format string, a ...interface{}) {
	if c.logger != nil {
		c.logger.Debugf(format, a...)
	}
}

func (c *Camera) infof(format string, a ...interface{}) {
	if c.logger != nil {
		c.logger.Infof(format, a...)
	}
}

func (c *Camera) warnf(format string, a ...interface{}) {
	if c.logger != nil {
		c.logger.Warnf(format, a...)
	}
}

func (c *Camera) errorf(format string, a ...interface{}) {
	if c.logger != nil {
		c.logger.Errorf(format, a...)
	}
}

