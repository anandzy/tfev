package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Reads the variable names from the variable list file
func readVarList(filename string) ([]string, error) {
	var vars []string
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" && !strings.HasPrefix(line, "#") {
			vars = append(vars, line)
		}
	}
	return vars, scanner.Err()
}

func main() {
	if len(os.Args) == 2 && (os.Args[1] == "-v" || os.Args[1] == "--version") {
		fmt.Println("filtervars version 1.0.0")
		os.Exit(0)
	}

	if len(os.Args) != 4 {
		fmt.Printf("Usage: %s <input_tfvars_file> <variable_list_file> <output_tfvars_file>\n", os.Args[0])
		os.Exit(1)
	}

	tfvarsFile := os.Args[1]
	varFile := os.Args[2]
	outFile := os.Args[3]

	// Step 1: Get list of variables to filter
	filterVars, err := readVarList(varFile)
	if err != nil {
		fmt.Printf("Failed to read %s: %s\n", varFile, err)
		return
	}

	// Step 2: Read all lines from tfvars
	tfLines := []string{}
	f, err := os.Open(tfvarsFile)
	if err != nil {
		fmt.Printf("Failed to open %s: %s\n", tfvarsFile, err)
		return
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		tfLines = append(tfLines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		fmt.Printf("Failed to read lines: %s\n", err)
		return
	}

	// Step 3: Parse and extract variable blocks/assignments
	varMap := make(map[string][]string)
	i := 0
	for i < len(tfLines) {
		line := tfLines[i]
		trimmed := strings.TrimSpace(line)
		// skip comments and empty
		if trimmed == "" || strings.HasPrefix(trimmed, "#") || strings.HasPrefix(trimmed, "//") {
			i++
			continue
		}

		// Check for variable assignment: name = ...
		eqIdx := strings.Index(trimmed, "=")
		if eqIdx > 0 {
			varName := strings.TrimSpace(trimmed[:eqIdx])
			valueStart := strings.TrimSpace(trimmed[eqIdx+1:])
			block := []string{line}

			// Handle multi-line values (lists, maps)
			if strings.HasPrefix(valueStart, "[") && !strings.Contains(valueStart, "]") {
				i++
				for i < len(tfLines) {
					block = append(block, tfLines[i])
					if strings.Contains(tfLines[i], "]") {
						break
					}
					i++
				}
			} else if strings.HasPrefix(valueStart, "{") && !strings.Contains(valueStart, "}") {
				i++
				for i < len(tfLines) {
					block = append(block, tfLines[i])
					if strings.Contains(tfLines[i], "}") {
						break
					}
					i++
				}
			}
			varMap[varName] = block
		}
		i++
	}

	// Step 4: Write filtered variables in order to output file
	out, err := os.Create(outFile)
	if err != nil {
		fmt.Printf("Failed to create %s: %s\n", outFile, err)
		return
	}
	defer out.Close()

	wroteAny := false
	for _, v := range filterVars {
		if block, ok := varMap[v]; ok {
			for _, l := range block {
				out.WriteString(l + "\n")
			}
			out.WriteString("\n")
			wroteAny = true
		}
	}
	if !wroteAny {
		fmt.Println("None of the variables from the variable list were found in tfvars!")
	} else {
		fmt.Printf("Filtered vars written to %s\n", outFile)
	}
}
