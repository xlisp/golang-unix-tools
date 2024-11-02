package main

import (
    "fmt"
    "os"
    "path/filepath"
    "sort"
)

type dirSize struct {
    path string
    size int64
}

func main() {
    entries, err := os.ReadDir(".")
    if err != nil {
        fmt.Printf("Error reading directory: %v\n", err)
        return
    }

    var sizes []dirSize
    for _, entry := range entries {
        size, err := getDirSize(entry.Name())
        if err != nil {
            fmt.Printf("Error getting size for %s: %v\n", entry.Name(), err)
            continue
        }
        sizes = append(sizes, dirSize{entry.Name(), size})
    }

    // Sort by size in descending order
    sort.Slice(sizes, func(i, j int) bool {
        return sizes[i].size > sizes[j].size
    })

    // Print sizes in human-readable format
    for _, ds := range sizes {
        fmt.Printf("%s\t%s\n", formatSize(ds.size), ds.path)
    }
}

func getDirSize(path string) (int64, error) {
    var size int64
    err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        if !info.IsDir() {
            size += info.Size()
        }
        return nil
    })
    return size, err
}

func formatSize(bytes int64) string {
    const (
        B  = 1
        KB = B * 1024
        MB = KB * 1024
        GB = MB * 1024
        TB = GB * 1024
    )

    switch {
    case bytes >= TB:
        return fmt.Sprintf("%.2fT", float64(bytes)/float64(TB))
    case bytes >= GB:
        return fmt.Sprintf("%.2fG", float64(bytes)/float64(GB))
    case bytes >= MB:
        return fmt.Sprintf("%.2fM", float64(bytes)/float64(MB))
    case bytes >= KB:
        return fmt.Sprintf("%.2fK", float64(bytes)/float64(KB))
    default:
        return fmt.Sprintf("%dB", bytes)
    }
}
