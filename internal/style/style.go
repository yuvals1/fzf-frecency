package style

import (
    "fmt"
    "path/filepath"
    "regexp"
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

// FileIcon represents an icon with its associated color
type FileIcon struct {
    Icon  string
    Color string
}

var (
    // Extension to icon mapping
    extensionIcons = map[string]FileIcon{
        // Programming Languages
        ".go":    {Icon: "", Color: Cyan},
        ".py":    {Icon: "", Color: Blue},
        ".js":    {Icon: "", Color: Yellow},
        ".ts":    {Icon: "", Color: Blue},
        ".jsx":   {Icon: "", Color: Blue},
        ".tsx":   {Icon: "", Color: Blue},
        ".html":  {Icon: "", Color: Red},
        ".css":   {Icon: "", Color: Blue},
        ".scss":  {Icon: "", Color: Magenta},
        ".json":  {Icon: "", Color: Yellow},
        ".yaml":  {Icon: "", Color: Yellow},
        ".yml":   {Icon: "", Color: Yellow},
        ".xml":   {Icon: "", Color: Yellow},
        ".sh":    {Icon: "", Color: Green},
        ".bash":  {Icon: "", Color: Green},
        ".zsh":   {Icon: "", Color: Green},

        // Documents
        ".md":    {Icon: "", Color: Blue},
        ".txt":   {Icon: "", Color: Blue},
        ".pdf":   {Icon: "", Color: Red},
        ".doc":   {Icon: "", Color: Blue},
        ".docx":  {Icon: "", Color: Blue},

        // Images
        ".png":   {Icon: "", Color: Magenta},
        ".jpg":   {Icon: "", Color: Magenta},
        ".jpeg":  {Icon: "", Color: Magenta},
        ".gif":   {Icon: "", Color: Magenta},
        ".svg":   {Icon: "", Color: Magenta},

        // Archives
        ".zip":   {Icon: "", Color: Red},
        ".tar":   {Icon: "", Color: Red},
        ".gz":    {Icon: "", Color: Red},
        ".7z":    {Icon: "", Color: Red},

        // Config files
        ".toml":  {Icon: "", Color: Yellow},
        ".conf":  {Icon: "", Color: Yellow},
        ".ini":   {Icon: "", Color: Yellow},
    }

    // Special filenames to icon mapping
    filenameIcons = map[string]FileIcon{
        "go.mod":      {Icon: "󰟓", Color: Cyan},
        "go.sum":      {Icon: "󰟓", Color: Cyan},
        "Makefile":    {Icon: "", Color: Yellow},
        "README":      {Icon: "", Color: Blue},
        "readme.md":   {Icon: "", Color: Blue},
        "README.md":   {Icon: "", Color: Blue},
        ".gitignore":  {Icon: "", Color: Red},
        "Dockerfile":  {Icon: "", Color: Blue},
        "docker-compose.yml": {Icon: "", Color: Blue},
    }

    // Default icon for unknown files
    defaultIcon = FileIcon{Icon: "", Color: Reset}

    // ANSI regex for stripping color codes
    ansiRegex = regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]`)
)

// GetFileIcon returns the appropriate icon and color for a file
func GetFileIcon(path string) FileIcon {
    // Check special filenames first
    if icon, ok := filenameIcons[filepath.Base(path)]; ok {
        return icon
    }

    // Check file extension
    ext := strings.ToLower(filepath.Ext(path))
    if icon, ok := extensionIcons[ext]; ok {
        return icon
    }

    // Return default icon
    return defaultIcon
}

// FormatPath returns a colored path with an appropriate icon
func FormatPath(path string) string {
    icon := GetFileIcon(path)
    return fmt.Sprintf("%s%s %s%s%s", 
        icon.Color, 
        icon.Icon,
        Bold,
        path,
        Reset)
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

// StripAnsi removes ANSI color codes from a string
func StripAnsi(s string) string {
    return ansiRegex.ReplaceAllString(s, "")
}
