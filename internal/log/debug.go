package log

import (
    "fmt"
    "os"
)

var isDebugEnabled bool

func init() {
    _, isDebugEnabled = os.LookupEnv("DEBUG_FZF_FRECENCY")
}

// Debug prints a debug message if DEBUG_FZF_FRECENCY is set
func Debug(format string, args ...interface{}) {
    if isDebugEnabled {
        fmt.Fprintf(os.Stderr, "Debug: "+format+"\n", args...)
    }
}
