package logster

import (
	"fmt"
	"runtime"
)

// VersionStr is version of logster
const VersionStr = "1.0.0"

// Version returns the version string of logster
func Version(app string) string {
	return fmt.Sprintf("%s v%s built with %s", app, VersionStr, runtime.Version())
}
