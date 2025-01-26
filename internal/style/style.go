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

var iconMap *IconMap

func init() {
    iconMap = NewIconMap()
}

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

// FormatSplitPath formats a path by separating filename and directory with colors and icons
func FormatSplitPath(path string) string {
    filename := filepath.Base(path)
    dirname := filepath.Dir(path)
    
    // Clean up directory path
    if dirname == "." {
        dirname = ""
    } else if strings.HasPrefix(dirname, "./") {
        dirname = dirname[2:]
    }

    // Get icon and color for the file
    iconDef := iconMap.Get(filename)

    // Directory always uses blue
    dirColor := Blue

    // Format the output based on whether we have a directory component
    if dirname == "" {
        return fmt.Sprintf("%s %s%s%s", 
            iconDef.Icon,
            iconDef.Color, filename, Reset)
    }

    return fmt.Sprintf("%s %s%-30s%s %s%s%s",
        iconDef.Icon,
        iconDef.Color, filename, Reset,
        dirColor, dirname, Reset)
}

// FormatPath returns a colored path based on its extension (legacy function)
func FormatPath(path string) string {
    if strings.HasPrefix(path, "./") {
        path = path[2:]
    }

    iconDef := iconMap.Get(path)
    return fmt.Sprintf("%s %s%s%s", iconDef.Icon, iconDef.Color, path, Reset)
}

// GetIcon is a helper function to get icon for a path
func GetIcon(path string) string {
    return iconMap.Get(path).Icon
}

// GetColor is a helper function to get color for a path
func GetColor(path string) string {
    return iconMap.Get(path).Color
}

// ColorizeFilename applies appropriate color to a filename
func ColorizeFilename(filename string) string {
    iconDef := iconMap.Get(filename)
    return fmt.Sprintf("%s%s%s", iconDef.Color, filename, Reset)
}

// ColorizeDirectory applies directory color to a path
func ColorizeDirectory(path string) string {
    return fmt.Sprintf("%s%s%s", Blue, path, Reset)
}

// StripANSI removes ANSI escape sequences efficiently
func StripANSI(s string) string {
    var b strings.Builder
    b.Grow(len(s))
    inEscape := false

    for i := 0; i < len(s); i++ {
        if s[i] == '\x1b' {
            inEscape = true
            continue
        }
        if inEscape {
            if (s[i] >= 'A' && s[i] <= 'Z') || (s[i] >= 'a' && s[i] <= 'z') {
                inEscape = false
            }
            continue
        }
        b.WriteByte(s[i])
    }
    return b.String()
}
