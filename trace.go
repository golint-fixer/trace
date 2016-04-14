// Package trace implements a featured and structured generic logger
// interface designed to be used in core and third-party vinxi components.
package trace

import (
	"io"
	"net/http"
	"sync"
	"time"

	logrus "gopkg.in/Sirupsen/logrus.v0"
	"gopkg.in/vinxi/log.v0"
)

// TracerFunc represents the log trace function.
type TracerFunc func(l log.Interface, w http.ResponseWriter, r *http.Request) log.Interface

// mu provides thread synchronization for the public singleton API.
var mu = sync.Mutex{}

// Logger represents the standard default vinxi log based logger.
var Logger = log.New()

// Default tracer preconfigured instance for convenience.
var Default = New()

// init configures the logger instance.
func init() {
	SetFormatter(&logrus.JSONFormatter{})
}

// SetOutput sets the standard logger trace output.
func SetOutput(out io.Writer) {
	mu.Lock()
	defer mu.Unlock()
	Logger.Out = out
}

// SetFormatter sets the standard logger trace formatter.
// Defaults to JSON formatter.
func SetFormatter(formatter logrus.Formatter) {
	mu.Lock()
	defer mu.Unlock()
	Logger.Formatter = formatter
}

// SetLevel sets the standard logger level.
func SetLevel(level logrus.Level) {
	mu.Lock()
	defer mu.Unlock()
	Logger.Level = level
}

// GetLevel returns the standard logger level.
func GetLevel() logrus.Level {
	mu.Lock()
	defer mu.Unlock()
	return Logger.Level
}

// AddHook adds a hook to the standard logger hooks.
func AddHook(hook logrus.Hook) {
	mu.Lock()
	defer mu.Unlock()
	Logger.Hooks.Add(hook)
}

// Tracer provides HTTP tracing capabilities to incoming traffic.
type Tracer struct {
	tracers []TracerFunc
	logger  *logrus.Logger
}

// New creates a new Tracer with default settings.
func New() *Tracer {
	return &Tracer{
		logger:  Logger,
		tracers: []TracerFunc{DefaultTracer},
	}
}

// AddTracer adds a new custom tracer function.
func (t *Tracer) AddTracer(tracer ...TracerFunc) {
	t.tracers = append(t.tracers, tracer...)
}

// SetTracer sets one or multiple tracer functions, removing the old ones.
func (t *Tracer) SetTracer(tracers ...TracerFunc) {
	t.tracers = tracers
}

// HandleHTTP read and report an incoming HTTP request.
func (t *Tracer) HandleHTTP(w http.ResponseWriter, r *http.Request, h http.Handler) {
	// Run tracer functions, chaining the logger across calls
	var l log.Interface = t.logger
	for _, tracer := range t.tracers {
		if logger := tracer(l, w, r); logger != nil {
			l = logger
		}
	}

	// Continue with next middleware
	h.ServeHTTP(w, r)
}

// DefaultTracer reads the request info and writes a structured data entry to the default logger.
// It's used as default tracer function. Returns an entry logger for subsequent writes.
func DefaultTracer(l log.Interface, w http.ResponseWriter, r *http.Request) log.Interface {
	entry := logrus.NewEntry(Logger)
	logger := entry.WithFields(logrus.Fields{
		"protocol":      r.Proto,
		"method":        r.Method,
		"host":          r.Host,
		"path":          r.RequestURI,
		"ip":            r.RemoteAddr,
		"headers":       r.Header,
		"contentlength": r.ContentLength,
		"time":          time.Now(),
	})
	logger.Message = "trace"
	logger.Time = time.Now()
	logger.Infof("%s %s [%s] - %s%s", r.Proto, r.Method, r.RemoteAddr, r.Host, r.RequestURI)
	return logger
}
