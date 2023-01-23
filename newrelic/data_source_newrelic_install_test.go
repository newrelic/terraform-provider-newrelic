//go:build integration
// +build integration

package newrelic

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccNewRelicInstallDataSourceBasicLinux(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testNewRelicInstallBasicConfig("linux"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicInstallDataSourceIsCorrectOS("data.newrelic_install.basic", "linux"),
					testAccCheckNewRelicInstallDataSourceHasCLIInstall("data.newrelic_install.basic", "linux", true),
					testAccCheckNewRelicInstallDataSourceHasAssumeYes("data.newrelic_install.basic"),
				),
			},
		},
	})
}

func TestAccNewRelicInstallDataSourceBasicWindows(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testNewRelicInstallBasicConfig("windows"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicInstallDataSourceIsCorrectOS("data.newrelic_install.basic", "windows"),
					testAccCheckNewRelicInstallDataSourceHasCLIInstall("data.newrelic_install.basic", "windows", true),
					testAccCheckNewRelicInstallDataSourceHasAssumeYes("data.newrelic_install.basic"),
				),
			},
		},
	})
}

func TestAccNewRelicInstallDataSourceBasicLinuxWithoutCLIDownload(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testNewRelicInstallBasicConfigWithoutDownload("linux"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicInstallDataSourceIsCorrectOS("data.newrelic_install.basic", "linux"),
					testAccCheckNewRelicInstallDataSourceHasCLIInstall("data.newrelic_install.basic", "linux", false),
					testAccCheckNewRelicInstallDataSourceHasAssumeYes("data.newrelic_install.basic"),
				),
			},
		},
	})
}

func TestAccNewRelicInstallDataSourceBasicWindowsWithoutCLIDownload(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testNewRelicInstallBasicConfigWithoutDownload("windows"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicInstallDataSourceIsCorrectOS("data.newrelic_install.basic", "windows"),
					testAccCheckNewRelicInstallDataSourceHasCLIInstall("data.newrelic_install.basic", "windows", false),
					testAccCheckNewRelicInstallDataSourceHasAssumeYes("data.newrelic_install.basic"),
				),
			},
		},
	})
}

func TestAccNewRelicInstallDataSourceTagsLinux(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testNewRelicInstallBasicConfigWithTags("linux"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicInstallDataSourceIsCorrectOS("data.newrelic_install.basic", "linux"),
					testAccCheckNewRelicInstallDataSourceHasCLIInstall("data.newrelic_install.basic", "linux", true),
					testAccCheckNewRelicInstallDataSourceHasTags("data.newrelic_install.basic"),
					testAccCheckNewRelicInstallDataSourceHasAssumeYes("data.newrelic_install.basic"),
				),
			},
		},
	})
}

func TestAccNewRelicInstallDataSourceRecipeLinux(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testNewRelicInstallBasicConfigWithRecipe("linux"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicInstallDataSourceIsCorrectOS("data.newrelic_install.basic", "linux"),
					testAccCheckNewRelicInstallDataSourceHasCLIInstall("data.newrelic_install.basic", "linux", true),
					testAccCheckNewRelicInstallDataSourceHasRecipe("data.newrelic_install.basic", "samwise"),
					testAccCheckNewRelicInstallDataSourceHasAssumeYes("data.newrelic_install.basic"),
				),
			},
		},
	})
}

func testNewRelicInstallBasicConfig(os string) string {
	return fmt.Sprintf(`
data "newrelic_install" "basic" {
	os = "%s"
}`, os)
}

func testNewRelicInstallBasicConfigWithoutDownload(os string) string {
	return fmt.Sprintf(`
data "newrelic_install" "basic" {
	os = "%s"
	download_cli = false
}`, os)
}

func testNewRelicInstallBasicConfigWithTags(os string) string {
	return fmt.Sprintf(`
data "newrelic_install" "basic" {
	os = "%s"
	tag {
        key = "whatabout"
        value = "second_breakfast"
    }

    tag {
        key = "afternoon_tea"
        value = "false"
    }
}`, os)
}

func testNewRelicInstallBasicConfigWithRecipe(os string) string {
	return fmt.Sprintf(`
data "newrelic_install" "basic" {
	os = "%s"
	recipe = [
		"frodo",
		"samwise",
		"pippin",
		"merry",
	]
}`, os)
}

func testAccCheckNewRelicInstallDataSourceGetCommand(s *terraform.State, n string) string {
	r := s.RootModule().Resources[n]
	return r.Primary.Attributes["command"]
}

func testAccCheckNewRelicInstallDataSourceHasCLIInstall(n string, os string, needsCommand bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		command := testAccCheckNewRelicInstallDataSourceGetCommand(s, n)
		hasCommand := false
		switch os {
		case "linux":
			hasCommand = strings.Contains(command, "curl")
		case "windows":
			hasCommand = strings.Contains(command, "DownloadFile")
		default:
			return fmt.Errorf("unknown OS cannot continue: %s", os)
		}

		if hasCommand == needsCommand {
			return nil
		}

		return fmt.Errorf("download CLI command expecation (%t) does not line up with result (%t): %s", needsCommand, hasCommand, command)
	}
}

func testAccCheckNewRelicInstallDataSourceIsCorrectOS(n string, os string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		command := testAccCheckNewRelicInstallDataSourceGetCommand(s, n)
		switch os {
		case "linux":
			if strings.Contains(command, "/newrelic") {
				return nil
			}
		case "windows":
			if strings.Contains(command, "newrelic.exe") {
				return nil
			}
		default:
			return fmt.Errorf("unknown OS cannot continue: %s", os)
		}

		return fmt.Errorf("incorrect OS found in command, was expecting %s but found: %s", os, command)
	}
}

func testAccCheckNewRelicInstallDataSourceHasTags(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		command := testAccCheckNewRelicInstallDataSourceGetCommand(s, n)
		if strings.Contains(command, "--tag") {
			return nil
		}

		return fmt.Errorf("was expecting tags, no tags found: %s", command)
	}
}

func testAccCheckNewRelicInstallDataSourceHasRecipe(n string, recipe string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		command := testAccCheckNewRelicInstallDataSourceGetCommand(s, n)
		if strings.Contains(command, "-n") && strings.Contains(command, recipe) {
			return nil
		}

		return fmt.Errorf("was expecting a recipe, supplied recipe (%s) not found: %s", recipe, command)
	}
}

func testAccCheckNewRelicInstallDataSourceHasAssumeYes(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		command := testAccCheckNewRelicInstallDataSourceGetCommand(s, n)
		if strings.Contains(command, "-y") {
			return nil
		}

		return fmt.Errorf("assume yes not found in command: %s", command)
	}
}
