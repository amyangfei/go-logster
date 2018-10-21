package logster

import (
	"fmt"
	"runtime"
)

const VersionStr = "1.0.0"

func Version(app string) string {
	return fmt.Sprintf("%s v%s built with %s", app, VersionStr, runtime.Version())
}
