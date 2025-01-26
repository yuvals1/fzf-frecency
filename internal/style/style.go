package style

import (
    "fmt"
    "path/filepath"
    "strings"
)

// ANSI color codes
const (
    Reset    = "\033[0m"
    Bold     = "\033[1m"
    Blue     = "\033[34m"
    Green    = "\033[32m"
    Red      = "\033[31m"
    Cyan     = "\033[36m"
    Yellow   = "\033[33m"
    Magenta  = "\033[35m"
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
    
    // Get icon definition for this file
    iconDef := iconMap.Get(filename)
    
    if dirname == "." {
        dirname = ""
    } else if strings.HasPrefix(dirname, "./") {
        dirname = dirname[2:]
    }

    // Build the icon part only if we have an icon
    var iconPart string
    if iconDef.Icon != "" {
        iconPart = fmt.Sprintf("%s%s%s ", iconDef.Color, iconDef.Icon, Reset)
    }

    // Color the icon, but keep filename text in default color
    if dirname == "" {
        return fmt.Sprintf("%s%s", iconPart, filename)
    }

    // Removed fixed width formatting to avoid overflow
    return fmt.Sprintf("%s%s %s%s%s",
        iconPart,
        filename,
        Blue, dirname, Reset)
}

// FormatPath returns a colored path based on its extension
func FormatPath(path string) string {
    if strings.HasPrefix(path, "./") {
        path = path[2:]
    }

    iconDef := iconMap.Get(path)
    if iconDef.Icon == "" {
        return path
    }
    return fmt.Sprintf("%s%s%s %s", iconDef.Color, iconDef.Icon, Reset, path)
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
    if iconDef.Icon == "" {
        return filename
    }
    return fmt.Sprintf("%s%s%s", iconDef.Color, iconDef.Icon, Reset)
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
