package main

import (
	"bufio"
	"encoding/json"
	"golang.org/x/xerrors"
	"io"
	"time"
)

// testEvent is a single message emitted from test2json.
// Stolen from https://golang.org/cmd/test2json/.
type testEvent struct {
	Time    time.Time // encodes as an RFC3339-format string
	Action  string
	Package string
	Test    string
	Elapsed float64 // seconds
	Output  string
}

// parseTestEvents parses a newline delimited list of JSON test events.
func parseTestEvents(r io.Reader) ([]testEvent, error) {
	var events []testEvent

	sc := bufio.NewScanner(r)
	for sc.Scan() {
		var te testEvent
		err := json.Unmarshal(sc.Bytes(), &te)
		if err != nil {
			return nil, xerrors.Errorf("failed to unmarshal %s: %w", sc.Text(), err)
		}
		events = append(events, te)
	}
	return events, nil
}
