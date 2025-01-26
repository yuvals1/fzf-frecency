package style

import (
    "fmt"
    "os"
    "path/filepath"
    "strings"
)

type IconDef struct {
    Icon  string
    Color string
}

type IconMap struct {
    icons map[string]IconDef
}

// Convert 256-color code to ANSI escape sequence
func colorToANSI(color string) string {
    if color == "" || color == "Reset" {
        return Reset
    }
    // Convert numeric color code to ANSI escape sequence
    return fmt.Sprintf("\033[38;5;%sm", color)
}

// NewIconMap now uses the generated defaultIcons
func NewIconMap() *IconMap {
    im := &IconMap{
        icons: make(map[string]IconDef),
    }

    // Copy generated icons
    for k, v := range defaultIcons {
        im.icons[k] = IconDef{
            Icon:  v.Icon,
            Color: colorToANSI(v.Color), // Convert color code to ANSI sequence
        }
    }

    // Load from environment if present
    if env := os.Getenv("FZF_FRECENCY_ICONS"); env != "" {
        im.parseEnv(env)
    }

    return im
}

func (im *IconMap) Get(path string) IconDef {
    ext := strings.ToLower(filepath.Ext(path))
    name := filepath.Base(path)

    // Try exact filename match
    if val, ok := im.icons[name]; ok {
        return val
    }

    // Try extension match (including dot)
    if val, ok := im.icons[ext]; ok {
        return val
    }

    // Try extension match (without dot, only if ext is not empty)
    if ext != "" {
        if val, ok := im.icons[ext[1:]]; ok {
            return val
        }
    }

    // Return default icon
    return im.icons[""]
}

func (im *IconMap) parseEnv(env string) {
    for _, entry := range strings.Split(env, ":") {
        if entry == "" {
            continue
        }

        parts := strings.Split(entry, "=")
        if len(parts) != 2 {
            continue
        }

        im.icons[parts[0]] = IconDef{
            Icon:  parts[1], 
            Color: Reset,
        }
    }
}
