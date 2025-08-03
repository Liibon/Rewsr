package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var packCmd = &cobra.Command{
	Use:   "pack [IMAGE]",
	Short: "Build an Enclave Image File (EIF) from Docker image",
	Long:  `Pack converts a Docker image into an EIF file for AWS Nitro Enclaves.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		image := args[0]
		outputEIF, _ := cmd.Flags().GetString("output")
		entrypoint, _ := cmd.Flags().GetString("entry")

		if outputEIF == "" {
			outputEIF = replaceSpecialChars(image) + ".eif"
		}

		// Sanitize output path to prevent path traversal
		cleanEIF, err := sanitizePath(outputEIF)
		if err != nil {
			return err
		}
		outputEIF = cleanEIF

		// Prompt-style output
		color.Cyan("> rewsr pack %s", image)

		err = buildEIF(image, outputEIF, entrypoint)
		if err != nil {
			color.Red("Failed: %v", err)
			return err
		}

		color.Green("Success: %s → %s", image, outputEIF)
		return nil
	},
}

// replaceSpecialChars converts image names to safe filenames
func replaceSpecialChars(s string) string {
	s = strings.ReplaceAll(s, ":", "-")
	s = strings.ReplaceAll(s, "/", "-")
	return s
}

// sanitizePath cleans and validates a file path to prevent traversal attacks.
func sanitizePath(path string) (string, error) {
	cleanPath := filepath.Clean(path)
	if strings.Contains(cleanPath, "..") {
		return "", fmt.Errorf("invalid output path: path traversal detected")
	}
	if filepath.IsAbs(cleanPath) {
		return "", fmt.Errorf("invalid output path: absolute paths are not allowed")
	}
	return cleanPath, nil
}

func buildEIF(baseImage, outputEIF, entrypoint string) error {
	// Check dependencies
	if _, err := exec.LookPath("nitro-cli"); err != nil {
		return fmt.Errorf("nitro-cli required: sudo yum install aws-nitro-enclaves-cli")
	}
	if _, err := exec.LookPath("docker"); err != nil {
		return fmt.Errorf("docker required: sudo yum install docker")
	}

	if entrypoint != "" {
		return buildEIFWithCustomEntrypoint(baseImage, outputEIF, entrypoint)
	}

	// Direct EIF build
	cmd := exec.Command("nitro-cli", "build-enclave", "--docker-uri", baseImage, "--output-file", outputEIF)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("nitro-cli build failed:\n%s", string(output))
	}

	return nil
}

func buildEIFWithCustomEntrypoint(baseImage, outputEIF, entrypoint string) error {
	tempDir, err := os.MkdirTemp("", "rewsr-build-")
	if err != nil {
		return fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(tempDir)

	dockerfile := fmt.Sprintf("FROM %s\nENTRYPOINT %s\n", baseImage, entrypoint)
	dockerfilePath := filepath.Join(tempDir, "Dockerfile")
	if err := os.WriteFile(dockerfilePath, []byte(dockerfile), 0644); err != nil {
		return fmt.Errorf("failed to create Dockerfile: %w", err)
	}

	tempImage := "rewsr-temp-" + replaceSpecialChars(outputEIF)

	// Build custom image
	buildCmd := exec.Command("docker", "build", "-t", tempImage, tempDir)
	if output, err := buildCmd.CombinedOutput(); err != nil {
		return fmt.Errorf("docker build failed:\n%s", string(output))
	}

	// Build EIF
	nitroCmd := exec.Command("nitro-cli", "build-enclave", "--docker-uri", tempImage, "--output-file", outputEIF)
	if output, err := nitroCmd.CombinedOutput(); err != nil {
		exec.Command("docker", "rmi", tempImage).Run()
		return fmt.Errorf("nitro-cli build failed:\n%s", string(output))
	}

	exec.Command("docker", "rmi", tempImage).Run()
	return nil
}

func init() {
	packCmd.Flags().StringP("output", "o", "", "Output EIF filename")
	packCmd.Flags().StringP("entry", "e", "", "Override entrypoint (JSON array)")
}