package app

import "os"

// App interface
type App interface {
	Up() error
	Down() error
	Remove() error
	Log(follow bool) error
	Port() string
}

// Container interface
type Container interface {
	Name() string
	Created() bool
	Running() bool
}

// New to create application
func New() (App, error) {
	dir, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	return newApp(dir)
}
