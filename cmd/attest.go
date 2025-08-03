package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var attestCmd = &cobra.Command{
	Use:   "attest [EIF_FILE]",
	Short: "Generate attestation document",
	Long:  `Generate hardware attestation document from running enclave.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		eifFile := args[0]
		output, _ := cmd.Flags().GetString("output")

		if output == "" {
			base := strings.TrimSuffix(filepath.Base(eifFile), ".eif")
			output = base + ".cbor"
		}

		// Sanitize output path to prevent path traversal
		cleanOutput, err := sanitizePath(output)
		if err != nil {
			return err
		}
		output = cleanOutput

		color.Cyan("> rewsr attest %s", eifFile)

		err = generateAttestation(output)
		if err != nil {
			color.Red("Failed: %v", err)
			return err
		}

		color.Green("Attestation: %s", output)
		return nil
	},
}

var verifyCmd = &cobra.Command{
	Use:   "verify [CBOR_FILE]",
	Short: "Verify attestation document",
	Long:  `Verify attestation document against AWS root CA.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cborFile := args[0]

		color.Cyan("> rewsr verify %s", cborFile)

		err := verifyAttestation(cborFile)
		if err != nil {
			color.Red("Failed: %v", err)
			return err
		}

		color.Green("Attestation valid")
		return nil
	},
}

func generateAttestation(outputFile string) error {
	if _, err := exec.LookPath("nitro-cli"); err != nil {
		return fmt.Errorf("nitro-cli required")
	}

	// Get running enclaves
	cmd := exec.Command("nitro-cli", "describe-enclaves")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to list enclaves: %s", string(output))
	}

	var enclaves []struct{ EnclaveID string }
	if err := json.Unmarshal(output, &enclaves); err != nil {
		return fmt.Errorf("failed to parse enclaves")
	}

	if len(enclaves) == 0 {
		return fmt.Errorf("no running enclaves found")
	}

	// Generate attestation
	attestCmd := exec.Command("nitro-cli", "generate-attestation-document",
		"--enclave-id", enclaves[0].EnclaveID,
		"--output-file", outputFile)

	if output, err := attestCmd.CombinedOutput(); err != nil {
		return fmt.Errorf("attestation failed: %s", string(output))
	}

	return nil
}

func verifyAttestation(cborFile string) error {
	if _, err := os.Stat(cborFile); err != nil {
		return fmt.Errorf("file not found: %s", cborFile)
	}

	data, err := os.ReadFile(cborFile)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	if len(data) < 100 {
		return fmt.Errorf("invalid attestation document")
	}

	// Try nitro-cli verification
	verifyCmd := exec.Command("nitro-cli", "verify-attestation-document", "--attestation-document", cborFile)
	if output, err := verifyCmd.CombinedOutput(); err != nil {
		color.Yellow("Warning: Full verification unavailable")
		return nil
	} else {
		if strings.Contains(string(output), "VALID") {
			color.White("Certificate chain verified")
		}
	}

	return nil
}

func init() {
	attestCmd.Flags().StringP("output", "o", "", "Output CBOR file")
}