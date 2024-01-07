package main

import (
    "bufio"
    "embed"
    "fmt"
    "os"
    "sort"
    "strings"
    "unicode"
)

//go:embed data/ext.txt
var extFile embed.FS

func main() {
    allowedExtensions, err := readAllowedExtensions()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error reading allowed extensions: %s\n", err)
        os.Exit(1)
    }

    extensionCount := make(map[string]int)
    scanner := bufio.NewScanner(os.Stdin)

    for scanner.Scan() {
        url := scanner.Text()
        ext := extractExtension(url, allowedExtensions)
        if ext != "" {
            extensionCount[ext]++
        }
    }

    if err := scanner.Err(); err != nil {
        fmt.Fprintf(os.Stderr, "Error reading from stdin: %s\n", err)
        os.Exit(1)
    }

    var extList []kv
    for k, v := range extensionCount {
        extList = append(extList, kv{k, v})
    }

    sort.Slice(extList, func(i, j int) bool {
        return extList[i].Value > extList[j].Value
    })

    fmt.Println("File Extension Count (sorted by frequency):")
    for _, kv := range extList {
        fmt.Printf("%s: %d\n", kv.Key, kv.Value)
    }
}

type kv struct {
    Key   string
    Value int
}

func readAllowedExtensions() (map[string]bool, error) {
    file, err := extFile.Open("data/ext.txt")
    if err != nil {
        return nil, err
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    allowedExtensions := make(map[string]bool)
    for scanner.Scan() {
        ext := strings.ToLower(strings.TrimSpace(scanner.Text()))
        allowedExtensions[ext] = true
    }

    if err := scanner.Err(); err != nil {
        return nil, err
    }

    return allowedExtensions, nil
}

func extractExtension(url string, allowedExtensions map[string]bool) string {
    if pos := strings.Index(url, "?"); pos != -1 {
        url = url[:pos]
    }

    if pos := strings.LastIndex(url, "."); pos != -1 && pos != len(url)-1 {
        substring := url[pos:] // Include the '.' in the substring
        if len(substring) < 11 && isValidExtension(substring) {
            ext := strings.ToLower(substring)
            if allowedExtensions[ext] {
                return ext
            }
        }
    }
    return ""
}

func isValidExtension(ext string) bool {
    for _, r := range ext {
        if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '.' && r != '-' && r != '_' {
            return false
        }
    }
    return true
}
