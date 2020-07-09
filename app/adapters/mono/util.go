package mono

import "io"

// Logger - represents the application's logger interface
type Logger interface {
	Errorf(template string, args ...interface{})
}

func closeBody(c io.Closer, log Logger) {
	if err := c.Close(); err != nil {
		log.Errorf("%+v", err)
	}
}
