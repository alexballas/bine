package control

import (
	"testing"
)

func TestParsePTLogEvent(t *testing.T) {
	evt := ParsePTLogEvent(`PT=/usr/bin/obfs4proxy SEVERITY=warning MESSAGE="connection to bridge failed: timeout"`)
	if evt.Code() != EventCodePTLog {
		t.Errorf("Code() = %q, want %q", evt.Code(), EventCodePTLog)
	}
	if evt.PT != "/usr/bin/obfs4proxy" {
		t.Errorf("PT = %q, want %q", evt.PT, "/usr/bin/obfs4proxy")
	}
	if evt.Severity != "warning" {
		t.Errorf("Severity = %q, want %q", evt.Severity, "warning")
	}
	if evt.Message != "connection to bridge failed: timeout" {
		t.Errorf("Message = %q, want %q", evt.Message, "connection to bridge failed: timeout")
	}
}

func TestParsePTStatusEvent(t *testing.T) {
	evt := ParsePTStatusEvent(`PT=/usr/bin/obfs4proxy TRANSPORT=obfs4 ADDRESS=1.2.3.4:443 CONNECT=SUCCESS`)
	if evt.Code() != EventCodePTStatus {
		t.Errorf("Code() = %q, want %q", evt.Code(), EventCodePTStatus)
	}
	if evt.PT != "/usr/bin/obfs4proxy" {
		t.Errorf("PT = %q, want %q", evt.PT, "/usr/bin/obfs4proxy")
	}
	if evt.Transport != "obfs4" {
		t.Errorf("Transport = %q, want %q", evt.Transport, "obfs4")
	}
	if evt.Values["ADDRESS"] != "1.2.3.4:443" {
		t.Errorf("Values[ADDRESS] = %q, want %q", evt.Values["ADDRESS"], "1.2.3.4:443")
	}
	if evt.Values["CONNECT"] != "SUCCESS" {
		t.Errorf("Values[CONNECT] = %q, want %q", evt.Values["CONNECT"], "SUCCESS")
	}
}

func TestPTEventsRecognized(t *testing.T) {
	if _, ok := recognizedEventCodesByCode[EventCodePTLog]; !ok {
		t.Errorf("recognizedEventCodesByCode missing %q", EventCodePTLog)
	}
	if _, ok := recognizedEventCodesByCode[EventCodePTStatus]; !ok {
		t.Errorf("recognizedEventCodesByCode missing %q", EventCodePTStatus)
	}
	// ParseEvent dispatches to the typed parsers rather than UnrecognizedEvent.
	if _, ok := ParseEvent(EventCodePTLog, "PT=x MESSAGE=y", nil).(*PTLogEvent); !ok {
		t.Errorf("ParseEvent(%q) did not return *PTLogEvent", EventCodePTLog)
	}
	if _, ok := ParseEvent(EventCodePTStatus, "PT=x TRANSPORT=y", nil).(*PTStatusEvent); !ok {
		t.Errorf("ParseEvent(%q) did not return *PTStatusEvent", EventCodePTStatus)
	}
}
