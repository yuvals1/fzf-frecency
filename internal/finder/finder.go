package finder

import (
    "bufio"
    "fmt"
    "os"
    "os/exec"
    "strings"

    "github.com/yuvals1/fzf-frecency/internal/log"
)

// FileFinder handles file discovery using fd
type FileFinder struct {
    fdPath          string
    excludePatterns []string
    includeHidden   bool
    maxDepth        int
}

// NewFileFinder creates a new FileFinder with default settings
func NewFileFinder() *FileFinder {
    return &FileFinder{
        fdPath:          "fd",
        excludePatterns: []string{"*.mypy", "*.git", "*.mypy_cache"},
        includeHidden:   true,
        maxDepth:        8,
    }
}

// buildFdCommand constructs the fd command with appropriate flags
func (f *FileFinder) buildFdCommand(root string) *exec.Cmd {
    args := []string{
        "--type", "f",         // files only
        "--strip-cwd-prefix",  // remove ./ prefix
        "--follow",           // follow symlinks
        "--color", "always",  // ensure colored output
    }

    // Add max depth
    args = append(args, "--max-depth", fmt.Sprintf("%d", f.maxDepth))

    // Include hidden files if specified
    if f.includeHidden {
        args = append(args, "--hidden")  // Removed --no-ignore to respect .gitignore
    }

    // Add exclude patterns
    for _, pattern := range f.excludePatterns {
        args = append(args, "--exclude", pattern)
    }

    // The search pattern is "." to match everything
    args = append(args, ".")

    cmd := exec.Command(f.fdPath, args...)
    cmd.Dir = root // Set working directory instead of passing as argument

    // Minimal logging that was in the original version
    log.Debug("Running fd command: %s %s (in directory %s)", 
        f.fdPath, strings.Join(args, " "), root)

    return cmd
}

// FindFiles returns a channel of found files
func (f *FileFinder) FindFiles(root string) (<-chan string, error) {
    // Verify fd is installed
    if _, err := exec.LookPath(f.fdPath); err != nil {
        return nil, fmt.Errorf("fd command not found. Please install fd-find: %w", err)
    }

    // Create and configure fd command
    cmd := f.buildFdCommand(root)
    
    // Get stdout pipe
    stdout, err := cmd.StdoutPipe()
    if err != nil {
        return nil, fmt.Errorf("creating stdout pipe: %w", err)
    }

    // Redirect stderr to os.Stderr for debugging
    cmd.Stderr = os.Stderr

    // Create output channel
    filesChan := make(chan string)

    // Start command
    if err := cmd.Start(); err != nil {
        return nil, fmt.Errorf("starting fd command: %w", err)
    }

    // Read output in goroutine
    go func() {
        defer close(filesChan)

        scanner := bufio.NewScanner(stdout)
        for scanner.Scan() {
            path := strings.TrimSpace(scanner.Text())
            if path != "" {
                filesChan <- path
            }
        }

        // Wait for command to complete
        if err := cmd.Wait(); err != nil {
            fmt.Fprintf(os.Stderr, "Error running fd: %v\n", err)
        }
    }()

    return filesChan, nil
}

// SetMaxDepth sets the maximum directory depth to search
func (f *FileFinder) SetMaxDepth(depth int) {
    f.maxDepth = depth
}

// SetExcludePatterns sets the patterns to exclude from search
func (f *FileFinder) SetExcludePatterns(patterns []string) {
    f.excludePatterns = patterns
}

// SetIncludeHidden sets whether to include hidden files
func (f *FileFinder) SetIncludeHidden(include bool) {
    f.includeHidden = include
}

// SetFdPath sets the path to the fd executable
func (f *FileFinder) SetFdPath(path string) {
    f.fdPath = path
}
