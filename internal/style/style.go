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

// File type icons mapping
var fileIcons = map[string]string{
    ".go":     "󰈔",
    ".py":     "",
    ".js":     "",
    ".json":   "󰘦",
    ".md":     "",
    ".txt":    "",
    ".yml":    "",
    ".yaml":   "",
    ".cpp":    "",
    ".h":      "",
    ".svelte": "󰎔",
    "":        "󰈔",
}

func getFileIcon(ext string) string {
    if icon, exists := fileIcons[ext]; exists {
        return icon + " "
    }
    return fileIcons[""] + " "
}

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

func FormatSplitPath(path string) string {
    filename := filepath.Base(path)
    dirname := filepath.Dir(path)
    
    if dirname == "." {
        dirname = ""
    } else if strings.HasPrefix(dirname, "./") {
        dirname = dirname[2:]
    }

    ext := strings.ToLower(filepath.Ext(filename))
    icon := getFileIcon(ext)

    var fileColor string
    switch ext {
    case ".go":
        fileColor = Cyan
    case ".py", ".pyc":
        fileColor = Blue
    case ".cpp", ".h":
        fileColor = Green
    case ".svelte":
        fileColor = Red
    default:
        fileColor = Reset
    }

    dirColor := Blue

    if dirname == "" {
        return fmt.Sprintf("%s%s%s%s", 
            icon,
            fileColor, filename, Reset)
    }

    return fmt.Sprintf("%s%s%-30s%s %s%s%s",
        icon,
        fileColor, filename, Reset,
        dirColor, dirname, Reset)
}
