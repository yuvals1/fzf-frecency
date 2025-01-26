package style

import (
    "fmt"
)

// ANSI color codes
const (
    Reset     = "\033[0m"
    Bold      = "\033[1m"
    Blue      = "\033[34m"
    Green     = "\033[32m"
    Red       = "\033[31m"
    Cyan      = "\033[36m"
    Yellow    = "\033[33m"
    Magenta   = "\033[35m"
)

// FormatScore returns a colored score based on its value
func FormatScore(score int) string {
    var color string
    switch {
    case score >= 400:
        color = Green
    case score >= 200:
        color = Yellow
    case score > 0:
        color = Blue
    default:
        color = Reset
    }
    return fmt.Sprintf("%s%5d%s", color, score, Reset)
}
