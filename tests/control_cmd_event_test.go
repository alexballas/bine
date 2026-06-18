package tests

import (
	"testing"

	"github.com/alexballas/bine/control"
)

func TestParseConfChangedEvent(t *testing.T) {
	event := control.ParseConfChangedEvent([]string{
		"LogMessageDomains=1",
		`DataDirectory="/tmp/tor data"`,
		`AccountingStart="month 1 00:00"`,
		"Nickname=",
		"Bridge",
	})

	require := newRequire(t)
	require.Equal(control.EventCodeConfChanged, event.Code())
	require.Equal([]string{
		"LogMessageDomains=1",
		`DataDirectory="/tmp/tor data"`,
		`AccountingStart="month 1 00:00"`,
		"Nickname=",
		"Bridge",
	}, event.Raw)
	require.Equal([]*control.KeyVal{
		control.NewKeyVal("LogMessageDomains", "1"),
		control.NewKeyVal("DataDirectory", "/tmp/tor data"),
		control.NewKeyVal("AccountingStart", "month 1 00:00"),
		{Key: "Nickname", ValSetAndEmpty: true},
		{Key: "Bridge"},
	}, event.Changes)
}

func TestParseEventDispatchesConfChangedEvent(t *testing.T) {
	event, ok := control.ParseEvent(control.EventCodeConfChanged, "", []string{
		"SocksPort=0",
	}).(*control.ConfChangedEvent)
	if !ok {
		t.Fatalf("ParseEvent(%q) = %T, want *control.ConfChangedEvent", control.EventCodeConfChanged, event)
	}
	if len(event.Changes) != 1 || event.Changes[0].Key != "SocksPort" || event.Changes[0].Val != "0" {
		t.Fatalf("Changes = %#v, want SocksPort=0", event.Changes)
	}
}

func TestParseCellStatsEvent(t *testing.T) {
	event := control.ParseCellStatsEvent(
		"ID=7 InboundQueue=q1 InboundConn=c1 InboundAdded=relay:1,created:2 " +
			"InboundRemoved=relay:3 InboundTime=relay:4 " +
			"OutboundQueue=q2 OutboundConn=c2 OutboundAdded=relay:5 " +
			"OutboundRemoved=relay:6 OutboundTime=relay:7",
	)

	require := newRequire(t)
	require.Equal("7", event.CircuitID)
	require.Equal("q1", event.InboundQueueID)
	require.Equal("c1", event.InboundConnID)
	require.Equal(map[string]int{"relay": 1, "created": 2}, event.InboundAdded)
	require.Equal(map[string]int{"relay": 3}, event.InboundRemoved)
	require.Equal(map[string]int{"relay": 4}, event.InboundTime)
	require.Equal("q2", event.OutboundQueueID)
	require.Equal("c2", event.OutboundConnID)
	require.Equal(map[string]int{"relay": 5}, event.OutboundAdded)
	require.Equal(map[string]int{"relay": 6}, event.OutboundRemoved)
	require.Equal(map[string]int{"relay": 7}, event.OutboundTime)
}
