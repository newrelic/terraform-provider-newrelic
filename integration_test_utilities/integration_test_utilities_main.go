package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

func main() {
	refreshFlag := flag.Bool("r", false, "Refresh flag to generate/update YAML content")
	flag.Parse()

	yamlFile := "integration_test_mappings.yaml"

	if *refreshFlag {
		integrationTestUtilitiesAlteredPrintln("🟠🟠🟠  Running this script in \"local\" mode, as the `-r` flag has been used to run the script...  🟠🟠🟠")
	} else {
		integrationTestUtilitiesAlteredPrintln("⚙️  Running this script in \"GitHub workflows\" mode...")
	}

	content, err := ioutil.ReadFile(yamlFile)
	if err != nil && !os.IsNotExist(err) {
		integrationTestUtilitiesAlteredPrintln(fmt.Sprintf("Error reading YAML file: %v", err))
		os.Exit(1)
	}

	// Case 1: Empty YAML file with -r flag
	if len(content) == 0 && *refreshFlag {
		generateInitialYAML(yamlFile)
		return
	}

	// Case 2: Empty YAML file without -r flag
	if len(content) == 0 && !*refreshFlag {
		integrationTestUtilitiesAlteredPrintln("Error: YAML file is empty; also, this program is not running in a local environment, but is running on GitHub workflows; please run the make compile job locally")
		os.Exit(1)
	}

	// Case 3: Non-empty YAML file
	if len(content) > 0 {
		processExistingYAML(yamlFile, *refreshFlag)
	}
}

func generateInitialYAML(yamlFile string) {
	if len(productMappingMapKeysSorted) == 0 {
		getProductMappingSortedKeys()
	}

	files := scanDirectory("../newrelic")
	mappings := make(FileMappings)

	for _, file := range files {
		if !shouldExcludeFile(file) {
			mappings[file] = FileMapping{
				Test:           strings.HasSuffix(file, "_test.go"),
				ProductMapping: assignProductMapping(file, productMappingMapKeysSorted),
			}
		}
	}

	writeYAMLFile(yamlFile, mappings)
	integrationTestUtilitiesAlteredPrintln("✅  Successfully generated initial YAML content")
}

func processExistingYAML(yamlFile string, refreshFlag bool) {
	var mappings FileMappings
	content, _ := ioutil.ReadFile(yamlFile)
	err := yaml.Unmarshal(content, &mappings)
	if err != nil {
		return
	}

	var unknownMappings []string
	var missingInYAML []string
	var notInDirectory []string

	for filename, mapping := range mappings {
		if mapping.ProductMapping == "UNKNOWN" {
			unknownMappings = append(unknownMappings, filename)
		}
	}

	currentFiles := scanDirectory("../newrelic")
	currentFilesMap := make(map[string]bool)

	for _, file := range currentFiles {
		if !shouldExcludeFile(file) {
			currentFilesMap[file] = true
			if _, exists := mappings[file]; !exists {
				missingInYAML = append(missingInYAML, file)
			}
		}
	}

	for filename := range mappings {
		if !currentFilesMap[filename] {
			notInDirectory = append(notInDirectory, filename)
		}
	}

	if !refreshFlag {
		reportErrors(unknownMappings, missingInYAML, notInDirectory, false)
	} else {
		if len(unknownMappings) > 0 {
			integrationTestUtilities_PrintYAMLProductMappingNotFound(unknownMappings)
			os.Exit(1)
		}
		reportErrors(unknownMappings, missingInYAML, notInDirectory, true)
		updateYAML(yamlFile, mappings, missingInYAML, notInDirectory)
	}
}

func scanDirectory(dir string) []string {
	var files []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(path, ".go") {
			relPath, _ := filepath.Rel(dir, path)
			files = append(files, relPath)
		}
		return nil
	})
	if err != nil {
		return nil
	}
	return files
}

func shouldExcludeFile(filename string) bool {
	return strings.HasPrefix(filename, "helpers_") ||
		strings.HasPrefix(filename, "provider_") ||
		filename == "config.go"
}

func writeYAMLFile(filename string, data FileMappings) {
	yamlData, err := yaml.Marshal(data)
	if err != nil {
		integrationTestUtilitiesAlteredPrintln(fmt.Sprintf("Error marshaling YAML: %v", err))
		os.Exit(1)
	}

	err = ioutil.WriteFile(filename, yamlData, 0644)
	if err != nil {
		integrationTestUtilitiesAlteredPrintln(fmt.Sprintf("Error writing YAML file: %v", err))
		os.Exit(1)
	}
}

func reportErrors(unknownMappings, missingInYAML, notInDirectory []string, isRefreshFlagEnabled bool) {
	var errors []string

	if len(unknownMappings) > 0 {
		errors = append(errors, fmt.Sprintf("❌  Error: The following files have unknown product mappings: %v", unknownMappings))
	}
	if len(missingInYAML) > 0 {
		errors = append(errors, fmt.Sprintf("Files missing in YAML: %v", missingInYAML))
	}
	if len(notInDirectory) > 0 {
		errors = append(errors, fmt.Sprintf("Files in YAML but not in directory: %v", notInDirectory))
	}

	if len(errors) > 0 {
		integrationTestUtilitiesAlteredPrintln("Errors found:")
		for _, err := range errors {
			integrationTestUtilitiesAlteredPrintln(err)
		}
		if !isRefreshFlagEnabled {
			// if errors are printed via this function without the -r mode, exit immediately after the print to avoid YAML changes
			// however, if the script is run in the -r mode, the script will continue to update the YAML file, hence no os.Exit()
			os.Exit(1)
		}
	} else {
		integrationTestUtilities_PrintYAMLNoChangesNeeded()
	}
}

func updateYAML(yamlFile string, mappings FileMappings, missingInYAML, notInDirectory []string) {
	if len(productMappingMapKeysSorted) == 0 {
		getProductMappingSortedKeys()
	}

	for _, file := range notInDirectory {
		delete(mappings, file)
	}

	for _, file := range missingInYAML {
		mappings[file] = FileMapping{
			Test:           strings.HasSuffix(file, "_test.go"),
			ProductMapping: assignProductMapping(file, productMappingMapKeysSorted),
		}
	}

	if len(missingInYAML) == 0 && len(notInDirectory) == 0 {
		integrationTestUtilities_PrintYAMLNoChangesNeeded()
		return
	}
	writeYAMLFile(yamlFile, mappings)
	integrationTestUtilitiesAlteredPrintln("✅  Successfully updated YAML content")
}

func assignProductMapping(file string, productMappingKeys []ProductMapping) string {
	for _, product := range productMappingKeys {
		patterns := productMappings[product]
		for _, pattern := range patterns {
			if strings.Contains(file, pattern) {
				return string(product)
			}
		}
	}
	return "UNKNOWN"
}
