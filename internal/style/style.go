package style

import (
    "fmt"
    "path/filepath"
    "strings"
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

// FormatPath returns a colored path based on its extension
func FormatPath(path string) string {
    if strings.HasPrefix(path, "./") {
        path = path[2:]
    }

    ext := strings.ToLower(filepath.Ext(path))
    var color string

    switch ext {
    case ".go":
        color = Cyan
    case ".py", ".pyc":
        color = Blue
    case ".js", ".ts", ".jsx", ".tsx":
        color = Yellow
    case ".html", ".css", ".scss":
        color = Magenta
    case ".md", ".txt":
        color = Green
    case ".json", ".yaml", ".yml":
        color = Red
    default:
        color = Reset
    }

    return fmt.Sprintf("%s%s%s", color, path, Reset)
}
