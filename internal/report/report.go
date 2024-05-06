package report

import (
	"io"
	"log/slog"
	"os"

	"github.com/gerardborst/ovpn-ldap-auth/internal/logging"
)

type Reporter struct {
	ControleFile io.Writer
}

var authControlReporter Reporter

var logger *slog.Logger

func NewReporter(af io.Writer) *Reporter {
	// logger is already created with config in main
	logger = logging.GetLogger()

	authControlReporter = Reporter{
		ControleFile: af,
	}
	return &authControlReporter
}

func (ar *Reporter) Report(authSuccess bool) {
	if authSuccess {
		_, err := ar.ControleFile.Write([]byte("1"))
		if err != nil {
			logger.Error("ReportSuccess errored", "error", err)
			os.Exit(1)
		}
		return
	}
	_, err := ar.ControleFile.Write([]byte("0"))
	if err != nil {
		logger.Error("ReportSuccess errored", "error", err)
		os.Exit(1)
	}
}
