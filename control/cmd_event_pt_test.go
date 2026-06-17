package control

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParsePTLogEvent(t *testing.T) {
	evt := ParsePTLogEvent(`PT=/usr/bin/obfs4proxy SEVERITY=warning MESSAGE="connection to bridge failed: timeout"`)
	require.Equal(t, EventCodePTLog, evt.Code())
	require.Equal(t, "/usr/bin/obfs4proxy", evt.PT)
	require.Equal(t, "warning", evt.Severity)
	require.Equal(t, "connection to bridge failed: timeout", evt.Message)
}

func TestParsePTStatusEvent(t *testing.T) {
	evt := ParsePTStatusEvent(`PT=/usr/bin/obfs4proxy TRANSPORT=obfs4 ADDRESS=1.2.3.4:443 CONNECT=SUCCESS`)
	require.Equal(t, EventCodePTStatus, evt.Code())
	require.Equal(t, "/usr/bin/obfs4proxy", evt.PT)
	require.Equal(t, "obfs4", evt.Transport)
	require.Equal(t, "1.2.3.4:443", evt.Values["ADDRESS"])
	require.Equal(t, "SUCCESS", evt.Values["CONNECT"])
}

func TestPTEventsRecognized(t *testing.T) {
	require.Contains(t, recognizedEventCodesByCode, EventCodePTLog)
	require.Contains(t, recognizedEventCodesByCode, EventCodePTStatus)
	// ParseEvent dispatches to the typed parsers rather than UnrecognizedEvent.
	require.IsType(t, &PTLogEvent{}, ParseEvent(EventCodePTLog, "PT=x MESSAGE=y", nil))
	require.IsType(t, &PTStatusEvent{}, ParseEvent(EventCodePTStatus, "PT=x TRANSPORT=y", nil))
}
