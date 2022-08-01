package main

import (
	"testing"
)

var mssageTests = []struct {
	in  string
	out Command
	err error
}{
	{"CONNECT TO RCCMD (SEND) MSG_EXECUTE MSG_TEXT Powerfail on SLC-20-CUBE3 (5): restored.", Command{Restored, 0}, nil},
	{"CONNECT TO RCCMD (SEND) MSG_EXECUTE MSG_TEXT Powerfail on SLC-20-CUBE3 (5): random message", Command{}, ErrUnkownEventType},
	{"CONNECT TO RCCMD (SEND) MSG_EXECUTE MSG_TEXT Powerfail on SLC-20-CUBE3 (5): Autonomietime 125.00 min.", Command{Failure, 125}, nil},
	{"random message", Command{}, ErrMessageParsing},
}

func TestParseMessage(t *testing.T) {
	for _, test := range mssageTests {
		cmd, err := ParseMessage(test.in)
		if cmd != test.out {
			t.Fatalf(`Error expected command = %v, want %v`, cmd, test.out)
		}
		if err != test.err {
			t.Fatalf(`Expected error = %v, want %v`, err, test.err)
		}
	}
}
