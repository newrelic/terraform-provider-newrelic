package main

import "fmt"

const integrationTestUtilitiesLogInitializer = "=== terraform-provider-newrelic === [ int-test-deps    ]: "

func integrationTestUtilitiesAlteredPrintln(string string) {
	fmt.Printf("%s%s\n", integrationTestUtilitiesLogInitializer, string)
}

func integrationTestUtilities_PrintYAMLNoChangesNeeded() {
	integrationTestUtilitiesAlteredPrintln("✅  All files checked, YAML up to date. Exiting...")
}

func integrationTestUtilities_PrintYAMLProductMappingNotFound(unknownMappings []string) {
	integrationTestUtilitiesAlteredPrintln(fmt.Sprintf("❌  Error: The following files have unknown product mappings: %v", unknownMappings))
}
