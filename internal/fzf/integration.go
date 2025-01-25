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

// DefaultOptions returns default FZF options
func DefaultOptions() *FzfOptions {
    return &FzfOptions{
        Preview:     "bat -n --color=always {2..}",
        Height:      "100%",
        MultiSelect: false,
        KeyBindings: map[string]string{
            "shift-up":   "preview-page-up",
            "shift-down": "preview-page-down",
        },
        ExtraArgs: []string{
            "--ansi",
            "--border",
            "--delimiter=\t",
            "--with-nth=1,2",
        },
    }
}

// FormatScoredFile formats a scored file for FZF display with colors and icons
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

    args := []string{
        "--height", opts.Height,
        "--preview", opts.Preview,
    }
    
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
            return nil, nil
        }
        return nil, fmt.Errorf("fzf process: %w", err)
    }

    var results []string
    for _, line := range strings.Split(strings.TrimSpace(string(output)), "\n") {
        if line == "" {
            continue
        }
        parts := strings.SplitN(line, "\t", 2)
        if len(parts) == 2 {
            // Strip ANSI codes from the path before returning
            cleanPath := style.StripAnsi(parts[1])
            results = append(results, cleanPath)
        }
    }

    return results, nil
}
