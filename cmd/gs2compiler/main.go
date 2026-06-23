package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/MorenoLand/GScript.gs2parser-go"
)

type arguments struct {
	inputs    []string
	output    string
	help      bool
	verbose   bool
	directory bool
	multi     bool
	err       string
}

const helpText = `
GS2 Script Compiler

Usage:
  %s [OPTIONS] INPUT [OUTPUT]
  %s INPUT -o OUTPUT
  %s --help

Arguments:
  INPUT              Input file (.gs2 or .txt) or directory
  OUTPUT             Output file (.gs2bc)

Options:
  -o, --output FILE  Specify output file
  -v, --verbose      Verbose output
  -h, --help         Show this help message

Examples:
  %s script.gs2                    # Creates script.gs2bc
  %s script.gs2 output.gs2bc       # Creates output.gs2bc
  %s script.gs2 -o output.gs2bc    # Creates output.gs2bc
  %s scripts/                      # Process directory
  %s file1.gs2 file2.gs2 file3.gs2 # Process multiple files (drag & drop)
`

func main() {
	os.Exit(run(os.Args, os.Stdout, os.Stderr))
}

func run(argv []string, stdout, stderr io.Writer) int {
	args := parseArgs(argv)
	prog := "gs2compiler"
	if len(argv) > 0 && argv[0] != "" {
		prog = argv[0]
	}
	if args.help {
		fmt.Fprintf(stdout, helpText, prog, prog, prog, prog, prog, prog, prog, prog)
		return 0
	}
	if args.err != "" {
		fmt.Fprintf(stderr, "Error: %s\nUse --help for usage information.\n", args.err)
		return 1
	}
	if args.directory {
		return processDirectory(args.inputs[0], args.verbose, stdout, stderr)
	}
	if args.multi {
		processFileList(args.inputs, args.verbose, "Multi-file", "", stdout, stderr)
		return 0
	}
	processFileList(args.inputs, args.verbose, "", args.output, stdout, stderr)
	return 0
}

func parseArgs(argv []string) arguments {
	var args arguments
	if len(argv) < 2 {
		args.err = "No input file specified. Use --help for usage information."
		return args
	}
	for i := 1; i < len(argv); i++ {
		arg := argv[i]
		switch arg {
		case "--help", "-h":
			args.help = true
			return args
		case "--verbose", "-v":
			args.verbose = true
		case "--output", "-o":
			i++
			if i >= len(argv) {
				args.err = "Missing output file after " + arg
				return args
			}
			args.output = argv[i]
		default:
			if strings.HasPrefix(arg, "-") {
				args.err = "Unknown option: " + arg
				return args
			}
			args.inputs = append(args.inputs, arg)
		}
	}
	if len(args.inputs) == 0 {
		args.err = "No input file specified"
		return args
	}
	if len(args.inputs) == 2 && args.output == "" {
		args.output = args.inputs[1]
		args.inputs = args.inputs[:1]
	}
	if len(args.inputs) == 1 {
		info, err := os.Stat(args.inputs[0])
		if err == nil && info.IsDir() {
			args.directory = true
			if args.output != "" {
				args.err = "Output file cannot be specified for directory mode"
			}
		} else if args.output == "" {
			args.output = defaultOutput(args.inputs[0])
		}
		return args
	}
	args.multi = true
	if args.output != "" {
		args.err = "Output file cannot be specified when processing multiple files"
		return args
	}
	for _, input := range args.inputs {
		info, err := os.Stat(input)
		if err == nil && info.IsDir() {
			args.err = "Cannot mix files and directories in multi-file mode"
			return args
		}
	}
	return args
}

func processDirectory(input string, verbose bool, stdout, stderr io.Writer) int {
	info, err := os.Stat(input)
	if err != nil || !info.IsDir() {
		fmt.Fprintf(stderr, "Error: Invalid directory: %s\n", input)
		return 1
	}
	if verbose {
		fmt.Fprintf(stdout, "Scanning directory: %s\n", input)
	}
	entries, err := os.ReadDir(input)
	if err != nil {
		fmt.Fprintf(stderr, "Error: %s\n", err)
		return 1
	}
	var files []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		path := filepath.Join(input, entry.Name())
		ext := filepath.Ext(path)
		if ext == ".gs2" || ext == ".txt" {
			files = append(files, path)
		} else if verbose {
			fmt.Fprintf(stdout, "Skipping file %s\n", path)
		}
	}
	processFileList(files, verbose, "Directory", "", stdout, stderr)
	return 0
}

func processFileList(files []string, verbose bool, mode, singleOutput string, stdout, stderr io.Writer) {
	processed, errors := 0, 0
	if mode != "" {
		fmt.Fprintf(stdout, "Processing %d files (%s mode):\n\n", len(files), mode)
	}
	for _, file := range files {
		if mode != "" {
			fmt.Fprintf(stdout, "Processing: %s\n", filepath.Base(file))
		}
		output := ""
		if len(files) == 1 && singleOutput != "" {
			output = singleOutput
		}
		if compileAndReport(file, output, verbose, len(files) == 1 && mode == "", stdout, stderr) {
			processed++
		} else {
			errors++
		}
	}
	if mode != "" {
		fmt.Fprintf(stdout, "\n%s processing complete: %d files processed, %d errors\n", mode, processed, errors)
	}
}

func compileAndReport(input, output string, verbose, single bool, stdout, stderr io.Writer) bool {
	if _, err := os.Stat(input); err != nil {
		fmt.Fprintln(stdout, " -> [ERROR] File does not exist")
		return false
	}
	if verbose {
		fmt.Fprintf(stdout, "Compiling file %s\n", input)
	}
	if output == "" {
		output = defaultOutput(input)
	}
	src, err := os.ReadFile(input)
	if err != nil {
		fmt.Fprintf(stdout, " -> [ERROR] %s\n", err)
		return false
	}
	res, err := gs2parser.Compile(string(src))
	if err != nil {
		fmt.Fprintf(stdout, " -> [ERROR] %s\n", err)
		return false
	}
	if err := os.WriteFile(output, res.Bytecode, 0644); err != nil {
		fmt.Fprintf(stderr, "Error: %s\n", err)
		return false
	}
	if verbose {
		fmt.Fprintf(stdout, " -> saved to %s\n", output)
	}
	if single && !verbose {
		fmt.Fprintf(stdout, "Compilation successful\n -> saved to %s\n", output)
	}
	return true
}

func defaultOutput(input string) string {
	ext := filepath.Ext(input)
	return filepath.Join(filepath.Dir(input), strings.TrimSuffix(filepath.Base(input), ext)+".gs2bc")
}
