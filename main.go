package main

import (
    "fmt"
    "os"
    "flag"
    "github.com/diplombmstu/rest-server-template/core"
)

var pars core.Parameters

func usage() {
    fmt.Fprintf(os.Stderr, "usage: example -stderrthreshold=[INFO|WARN|FATAL] -log_dir=[string] -cfg=[string]\n")
    flag.PrintDefaults()
    os.Exit(2)
}

func init() {
    flag.Usage = usage
    flag.StringVar(&pars.ConfigFile, "cfg", "", "The path to config file. Mandatory.")

    required := []string{"cfg"}
    flag.Parse()

    seen := make(map[string]bool)
    flag.Visit(func(f *flag.Flag) { seen[f.Name] = true })
    for _, req := range required {
        if !seen[req] {
            fmt.Fprintf(os.Stderr, "Missing required -%s argument/flag\n", req)
            usage()
        }
    }
}

func main() {
    core.Start(pars)
}
