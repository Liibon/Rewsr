package cmd

import (
	"testing"
)

func TestRootCommand(t *testing.T) {
	// Test that root command executes without error
	if err := rootCmd.Execute(); err != nil {
		t.Errorf("Root command failed: %v", err)
	}
}

func TestCommandsExist(t *testing.T) {
	tests := []struct {
		name        string
		commandName string
	}{
		{"pack command exists", "pack"},
		{"deploy command exists", "deploy"},
		{"attest command exists", "attest"},
		{"verify command exists", "verify"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd, _, err := rootCmd.Find([]string{tt.commandName})
			if err != nil {
				t.Errorf("Command %s not found: %v", tt.commandName, err)
			}
			if cmd == nil {
				t.Errorf("Command %s is nil", tt.commandName)
			}
		})
	}
}

func TestPackCommandFlags(t *testing.T) {
	// Test that pack command has expected flags for EIF building
	outputFlag := packCmd.Flags().Lookup("output")
	if outputFlag == nil {
		t.Error("Pack command missing --output flag")
	}

	entryFlag := packCmd.Flags().Lookup("entry")
	if entryFlag == nil {
		t.Error("Pack command missing --entry flag")
	}
}

func TestDeployCommandFlags(t *testing.T) {
	// Test that deploy command has expected flags for Nitro
	portFlag := deployCmd.Flags().Lookup("port")
	if portFlag == nil {
		t.Error("Deploy command missing --port flag")
	}

	cpuFlag := deployCmd.Flags().Lookup("cpu-count")
	if cpuFlag == nil {
		t.Error("Deploy command missing --cpu-count flag")
	}

	memoryFlag := deployCmd.Flags().Lookup("memory")
	if memoryFlag == nil {
		t.Error("Deploy command missing --memory flag")
	}
}

func TestAttestCommandFlags(t *testing.T) {
	// Test that attest command has expected flags for CBOR output
	outputFlag := attestCmd.Flags().Lookup("output")
	if outputFlag == nil {
		t.Error("Attest command missing --output flag")
	}
}

func TestEIFNameGeneration(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"nginx:alpine", "nginx-alpine.eif"},
		{"ubuntu", "ubuntu.eif"},
		{"registry.io/user/app:v1.0", "registry.io-user-app-v1.0.eif"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			// Test the EIF name generation logic that's now in pack.go
			safeName := replaceSpecialChars(tt.input)
			result := safeName + ".eif"

			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}