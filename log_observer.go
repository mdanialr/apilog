package apilog

import (
	"bytes"
	"encoding/json"
	"io"
	"strings"
	"sync"
	"time"
)

// ObservedLog is a concurrency-safe, ordered collection of observed Log(s).
type ObservedLog struct {
	mu   sync.RWMutex
	logs []loggedLog
	lvl  Level
	out  Output
}

func (o *ObservedLog) Writer() io.Writer { return o }

func (o *ObservedLog) Output() Output { return o.out }

func (o *ObservedLog) Level() Level { return o.lvl }

func (o *ObservedLog) Wait(_ time.Duration) {}

func (o *ObservedLog) Flush(_ time.Duration) {}

func (o *ObservedLog) Write(p []byte) (n int, err error) {
	m := make(map[string]any)
	_ = json.Unmarshal(bytes.TrimSpace(p), &m)
	var l loggedLog
	// grab level if possible
	if v, ok := m["level"].(string); ok {
		l.level = ParseLevel(v)
	}
	// grab message if possible
	if v, ok := m["msg"].(string); ok {
		l.msg = v
	}
	// put the rest to context
	l.context = m

	o.mu.Lock()
	o.logs = append(o.logs, l)
	o.mu.Unlock()

	return len(p), nil
}

// Len returns the number of items in the collection.
func (o *ObservedLog) Len() int {
	o.mu.RLock()
	n := len(o.logs)
	o.mu.RUnlock()
	return n
}

// All returns a copy of all the observed logs.
func (o *ObservedLog) All() []loggedLog {
	o.mu.RLock()
	ret := make([]loggedLog, len(o.logs))
	copy(ret, o.logs)
	o.mu.RUnlock()
	return ret
}

// TakeAll returns a copy of all the observed logs, and truncates the observed
// slice.
func (o *ObservedLog) TakeAll() []loggedLog {
	o.mu.Lock()
	ret := o.logs
	o.logs = nil
	o.mu.Unlock()
	return ret
}

// NewObserverWriter return new Writer implementer that write logs to memory
// and also return ObservedLog to help assert and check logged Log(s).
func NewObserverWriter(lvl Level, out Output) (Writer, *ObservedLog) {
	ol := &ObservedLog{
		lvl: lvl,
		out: out,
		mu:  sync.RWMutex{},
	}
	return ol, ol
}

type loggedLog struct {
	level   Level
	msg     string
	context map[string]any
}

// EqualLevel return true if given lvl is equal with level.
func (l *loggedLog) EqualLevel(lvl Level) bool {
	return l.level == lvl
}

// EqualMsg return true if given s is equal with message.
func (l *loggedLog) EqualMsg(s string) bool {
	return l.msg == s
}

// ContainMsg return true if message contain given s.
func (l *loggedLog) ContainMsg(s string) bool {
	return strings.Contains(l.msg, s)
}

// Get grab a data from context using given key.
func (l *loggedLog) Get(k string) any {
	return l.context[k]
}
