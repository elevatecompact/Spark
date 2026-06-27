package logger

import (
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
)

type Logger struct {
	zerolog.Logger
}

type Option func(*Options)

type Options struct {
	Level       zerolog.Level
	Output      io.Writer
	ServiceName string
	Environment string
}

func WithLevel(level zerolog.Level) Option {
	return func(o *Options) {
		o.Level = level
	}
}

func WithOutput(w io.Writer) Option {
	return func(o *Options) {
		o.Output = w
	}
}

func New(service, environment string) *Logger {
	var l zerolog.Logger

	if environment == "development" {
		l = zerolog.New(zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
			NoColor:    false,
		})
	} else {
		l = zerolog.New(os.Stdout)
	}

	l = l.With().
		Str("service", service).
		Str("environment", environment).
		Timestamp().
		Logger()

	return &Logger{l}
}

func NewWithOptions(service, environment string, opts ...Option) *Logger {
	options := Options{
		Level:       zerolog.InfoLevel,
		Output:      os.Stdout,
		ServiceName: service,
		Environment: environment,
	}
	for _, opt := range opts {
		opt(&options)
	}

	var l zerolog.Logger

	if environment == "development" {
		l = zerolog.New(zerolog.ConsoleWriter{
			Out:        options.Output,
			TimeFormat: time.RFC3339,
		})
	} else {
		l = zerolog.New(options.Output)
	}

	l = l.Level(options.Level).
		With().
		Str("service", service).
		Str("environment", environment).
		Timestamp().
		Logger()

	return &Logger{l}
}

func (l *Logger) WithComponent(component string) *Logger {
	nl := l.Logger.With().Str("component", component).Logger()
	return &Logger{nl}
}

func (l *Logger) WithRequestID(requestID string) *Logger {
	nl := l.Logger.With().Str("request_id", requestID).Logger()
	return &Logger{nl}
}

func (l *Logger) WithUserID(userID string) *Logger {
	nl := l.Logger.With().Str("user_id", userID).Logger()
	return &Logger{nl}
}
