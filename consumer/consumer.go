// Package consumer defines event consumers.
package consumer

// Consumer recives and process incoming events.
type Consumer interface {
	Start() error
}
