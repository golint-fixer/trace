package trace

import (
	"bytes"
	"net/http"
	"strings"
	"testing"

	"github.com/nbio/st"
	log "gopkg.in/Sirupsen/logrus.v0"
)

func TestTrace(t *testing.T) {
	var buffer bytes.Buffer
	SetOutput(&buffer)
	SetLevel(log.InfoLevel)
	SetFormatter(&log.TextFormatter{DisableColors: true, DisableTimestamp: true, DisableSorting: true})

	tracer := New()

	headers := http.Header{"foo": []string{"bar"}}
	req := &http.Request{
		Proto:         "HTTP/1.1",
		Method:        "GET",
		Host:          "localhost",
		RequestURI:    "/foo",
		Header:        headers,
		ContentLength: 10,
		RemoteAddr:    "127.0.0.1:32131",
	}

	var called bool
	tracer.HandleHTTP(nil, req, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	}))

	st.Expect(t, called, true)

	entry := safeString(buffer)
	st.Expect(t, strings.Contains(entry, "level=info"), true)
	st.Expect(t, strings.Contains(entry, "GET [127.0.0.1:32131] -"), true)
	st.Expect(t, strings.Contains(entry, "- localhost/foo\""), true)
	st.Expect(t, strings.Contains(entry, "time="), true)
	st.Expect(t, strings.Contains(entry, "protocol=\"HTTP/1.1\""), true)
	st.Expect(t, strings.Contains(entry, "host=localhost"), true)
	st.Expect(t, strings.Contains(entry, "path=\"/foo\""), true)
	st.Expect(t, strings.Contains(entry, "ip=\"127.0.0.1:32131\""), true)
	st.Expect(t, strings.Contains(entry, "headers=map[foo:[bar]]"), true)
	st.Expect(t, strings.Contains(entry, "contentlength=10"), true)
}

func safeString(buf bytes.Buffer) string {
	b := buf.Bytes()
	return string(b[:len(b)-2])
}
