package main

import (
	"errors"
	"net/http"
)

type Config struct {
	Port int
}

type ConfigBuilder struct {
	port *int // prevent zero value
}

func (b *ConfigBuilder) Port(port int) *ConfigBuilder {
	b.port = &port
	return b
}

// validation must be delayed to the Build() method to allow for method chaining
func (b *ConfigBuilder) Build() (Config, error) {
	cfg := Config{}

	if b.port == nil {
		cfg.Port = 8000
	} else {
		if *b.port == 0 {
			cfg.Port = 7000
		} else if *b.port < 0 {
			return Config{}, errors.New("port should be positive")
		} else {
			cfg.Port = *b.port
		}
	}

	return cfg, nil
}

// functional options allows for handling optional configurations
type options struct {
	port *int
}

type Option func(options *options) error

func WithPort(port int) Option {
	return func(options *options) error {
		if port < 0 {
			return errors.New("port should be positive")
		}
		options.port = &port
		return nil
	}
}

func NewServer(addr string, opts ...Option) (*http.Server, error) {
	var options options
	for _, opt := range opts {
		err := opt(&options)
		if err != nil {
			return nil, err
		}
	}

	var port int
	if options.port == nil {
		port = 8000
	} else if *options.port == 0 {
		port = 7000
	} else {
		port = *options.port
	}

	// ...

	_ = port
	return nil, nil
}

func main() {
	builder := ConfigBuilder{}
	builder.Port(8080)
	_, _ = builder.Build()

	_, _ = NewServer("localhost")
	_, _ = NewServer("localhost", WithPort(9000))
}
