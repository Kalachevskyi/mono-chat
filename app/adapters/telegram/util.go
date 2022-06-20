package telegram

// Logger - represents the application's logger interface.
type Logger interface {
	Errorf(template string, args ...interface{})
}
