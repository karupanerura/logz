package logger

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/glassonion1/logz/go111/internal/config"
	"github.com/glassonion1/logz/go111/internal/severity"
	"github.com/glassonion1/logz/go111/internal/spancontext"
	"github.com/glassonion1/logz/go111/internal/types"
)

var NowFunc = time.Now

const traceFmt = "projects/%s/traces/%s"

// Looger is for GCP
type Logger struct {
}

// WriteAppLog writes a application log to stdout
func (l *Logger) WriteApplicationLog(ctx context.Context, severity severity.Severity, format string, a ...interface{}) {
	// Gets the traceID and spanID
	sc := spancontext.Extract(ctx)

	// gets the source location
	var location types.SourceLocation
	if pc, file, line, ok := runtime.Caller(2); ok {
		if function := runtime.FuncForPC(pc); function != nil {
			location.Function = function.Name()
		}
		location.Line = fmt.Sprintf("%d", line)
		parts := strings.Split(file, "/")
		location.File = parts[len(parts)-1] // use short file name
	}

	trace := fmt.Sprintf(traceFmt, config.ProjectID, sc.TraceID)
	msg := fmt.Sprintf(format, a...)
	ety := &types.ApplicationLog{
		Severity:       severity.String(),
		Message:        msg,
		Time:           NowFunc(),
		SourceLocation: location,
		Trace:          trace,
		SpanID:         sc.SpanID,
		TraceSampled:   sc.TraceSampled,
	}

	if err := json.NewEncoder(os.Stdout).Encode(ety); err != nil {
		fmt.Printf("failed to write log: %v", err)
	}
}
