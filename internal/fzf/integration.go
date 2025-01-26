package fzf

import (
    "fmt"
    "io"
    "os"
    "os/exec"
    "strings"

    "github.com/yuvals1/fzf-frecency/internal/finder"
    "github.com/yuvals1/fzf-frecency/internal/style"
)

// FzfOptions contains configuration for fzf
type FzfOptions struct {
    Preview      string
    Height       string
    MultiSelect  bool
    ExtraArgs    []string
    KeyBindings  map[string]string
}

// DefaultOptions returns default FZF options matching user's configuration
func DefaultOptions() *FzfOptions {
    return &FzfOptions{
        Preview:     "bat -n --color=always {2..}",  // Use {2..} to skip the score column
        Height:      "100%",
        MultiSelect: false,
        KeyBindings: map[string]string{
            "shift-up":   "preview-page-up",
            "shift-down": "preview-page-down",
        },
        ExtraArgs: []string{
            "--ansi",            // Enable ANSI color codes
            "--delimiter=\\t",   // Use tab as delimiter
            "--with-nth=1,2",   // Show only score and path columns
        },
    }
}

// FormatScoredFile formats a scored file for FZF display with colors
func FormatScoredFile(file finder.ScoredFile) string {
    score := style.FormatScore(file.Score)
    formattedPath := style.FormatPath(file.Path)
    return fmt.Sprintf("%s\t%s", score, formattedPath)
}

// RunFzf runs fzf with the provided scored files
func RunFzf(files <-chan finder.ScoredFile, opts *FzfOptions) ([]string, error) {
    if opts == nil {
        opts = DefaultOptions()
    }

    args := []string{}
    
    // Add preview command
    if opts.Preview != "" {
        args = append(args, "--preview", opts.Preview)
    }
    
    // Add key bindings
    for key, action := range opts.KeyBindings {
        args = append(args, fmt.Sprintf("--bind=%s:%s", key, action))
    }
    
    if opts.MultiSelect {
        args = append(args, "-m")
    }
    
    args = append(args, opts.ExtraArgs...)
    
    cmd := exec.Command("fzf", args...)
    cmd.Stderr = os.Stderr
    
    stdin, err := cmd.StdinPipe()
    if err != nil {
        return nil, fmt.Errorf("creating stdin pipe: %w", err)
    }
    
    stdout, err := cmd.StdoutPipe()
    if err != nil {
        return nil, fmt.Errorf("creating stdout pipe: %w", err)
    }
    
    if err := cmd.Start(); err != nil {
        return nil, fmt.Errorf("starting fzf: %w", err)
    }

    go func() {
        defer stdin.Close()
        for file := range files {
            fmt.Fprintln(stdin, FormatScoredFile(file))
        }
    }()

    output, err := io.ReadAll(stdout)
    if err != nil {
        return nil, fmt.Errorf("reading fzf output: %w", err)
    }

    if err := cmd.Wait(); err != nil {
        if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 130 {
            return nil, nil // User cancelled with ESC
        }
        return nil, fmt.Errorf("fzf process: %w", err)
    }

    var results []string
    for _, line := range strings.Split(strings.TrimSpace(string(output)), "\n") {
        if line == "" {
            continue
        }
        // Split by tab and take the path part (second column)
        parts := strings.SplitN(line, "\t", 2)
        if len(parts) == 2 {
            path := parts[1]
            // Remove ./ prefix if present
            if strings.HasPrefix(path, "./") {
                path = path[2:]
            }
            results = append(results, path)
        }
    }

    return results, nil
}
