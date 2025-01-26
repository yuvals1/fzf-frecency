package style

import (
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

func NewIconMap() *IconMap {
    im := &IconMap{
        icons: make(map[string]IconDef),
    }

    // Load default icons
    im.loadDefaults()

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

    // Try extension match
    if val, ok := im.icons[ext]; ok {
        return val
    }

    // Return default icon
    return im.icons[""]
}

func (im *IconMap) loadDefaults() {
    // Default icons mapping
    defaults := map[string]IconDef{
        ".md":     {Icon: "", Color: Blue},
        ".go":     {Icon: "", Color: Cyan},
        ".py":     {Icon: "", Color: Blue},
        ".js":     {Icon: "", Color: Yellow},
        ".json":   {Icon: "", Color: Yellow},
        ".toml":   {Icon: "", Color: Red},
        ".lua":    {Icon: "", Color: Blue},
        ".zsh":    {Icon: "", Color: Green},
        "":        {Icon: "", Color: Reset}, // default
    }

    for k, v := range defaults {
        im.icons[k] = v
    }
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

        im.icons[parts[0]] = IconDef{Icon: parts[1], Color: Reset}
    }
}
