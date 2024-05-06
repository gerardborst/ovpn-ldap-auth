package report

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/gerardborst/ovpn-ldap-auth/internal/logging"
)

func TestMain(m *testing.M) {
	logger = logging.NewLogger(&logging.LogConfiguration{LogToFile: false})
	m.Run()
}

func TestReport(t *testing.T) {
	var reporter *Reporter
	var out *bytes.Buffer
	var tests = []struct {
		arg  bool
		want string
	}{
		{true, "1"},
		{false, "0"},
	}
	for _, test := range tests {
		descr := fmt.Sprintf("Report(%v)",
			test.arg)

		out = new(bytes.Buffer) // captured output
		reporter = NewReporter(out)
		reporter.Report(test.arg)
		got := out.String()
		if got != test.want {
			t.Errorf("%s = %q, want %q", descr, got, test.want)
		}
	}
}
