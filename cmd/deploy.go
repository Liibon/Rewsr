package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var deployCmd = &cobra.Command{
	Use:   "deploy [EIF_FILE]",
	Short: "Deploy EIF to AWS Nitro Enclave",
	Long:  `Deploy launches an EIF in an AWS Nitro Enclave.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		eifPath := args[0]
		port, _ := cmd.Flags().GetInt("port")
		cpuCount, _ := cmd.Flags().GetInt("cpu-count")
		memory, _ := cmd.Flags().GetInt("memory")

		// Add resource constraints
		if cpuCount < 1 || cpuCount > 16 {
			return fmt.Errorf("invalid cpu-count: must be between 1 and 16")
		}
		if memory < 512 || memory > 16384 {
			return fmt.Errorf("invalid memory: must be between 512 and 16384 MB")
		}

		color.Cyan("> rewsr deploy %s", eifPath)

		err := deployToNitro(eifPath, cpuCount, memory, port)
		if err != nil {
			color.Red("Failed: %v", err)
			return err
		}

		color.Green("Enclave running")
		if port > 0 {
			color.White("Port: %d", port)
		}
		return nil
	},
}

func deployToNitro(eifPath string, cpuCount, memory, port int) error {
	// Verify EIF file
	if stat, err := os.Stat(eifPath); err != nil {
		return fmt.Errorf("EIF file not found: %s", eifPath)
	} else if stat.Size() == 0 {
		return fmt.Errorf("EIF file is empty: %s", eifPath)
	}

	// Check dependencies
	if _, err := exec.LookPath("nitro-cli"); err != nil {
		return fmt.Errorf("nitro-cli required: sudo yum install aws-nitro-enclaves-cli")
	}
	if port > 0 {
		if _, err := exec.LookPath("vsock-proxy"); err != nil {
			return fmt.Errorf("vsock-proxy required for port mapping")
		}
	}

	// Terminate existing enclaves
	terminateExistingEnclaves()

	// Run enclave
	runArgs := []string{
		"nitro-cli", "run-enclave",
		"--cpu-count", strconv.Itoa(cpuCount),
		"--memory", strconv.Itoa(memory),
		"--eif-path", eifPath,
		"--debug-mode",
	}

	runCmd := exec.Command(runArgs[0], runArgs[1:]...)
	output, err := runCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to start enclave:\n%s", string(output))
	}

	// Parse enclave info for vsock proxy
	var enclaveInfo struct {
		EnclaveID  string `json:"EnclaveID"`
		EnclaveCID int    `json:"EnclaveCID"`
	}

	if err := json.Unmarshal(output, &enclaveInfo); err == nil && port > 0 && enclaveInfo.EnclaveCID > 0 {
		go startVsockProxy(port, enclaveInfo.EnclaveCID)
	}

	time.Sleep(2 * time.Second)
	return nil
}

func terminateExistingEnclaves() {
	describeCmd := exec.Command("nitro-cli", "describe-enclaves", "--json")
	if output, err := describeCmd.CombinedOutput(); err == nil {
		var enclaves []struct{ EnclaveID string }
		if json.Unmarshal(output, &enclaves) == nil {
			for _, enclave := range enclaves {
				exec.Command("nitro-cli", "terminate-enclave", "--enclave-id", enclave.EnclaveID).Run()
			}
		}
	}
}

func startVsockProxy(localPort, enclaveCID int) {
	proxyArgs := []string{
		"vsock-proxy",
		strconv.Itoa(localPort),
		fmt.Sprintf("%d:80", enclaveCID),
	}

	ctx := context.Background()
	proxyCmd := exec.CommandContext(ctx, proxyArgs[0], proxyArgs[1:]...)
	proxyCmd.Start()
}

func init() {
	deployCmd.Flags().IntP("port", "p", 0, "Local port for vsock proxy")
	deployCmd.Flags().Int("cpu-count", 2, "Number of CPUs")
	deployCmd.Flags().Int("memory", 2048, "Memory in MB")
}