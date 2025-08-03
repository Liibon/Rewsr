package cmd

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

// TestVerifyAttestation_FileNotFound ensures that verifyAttestation returns an error for non-existent files.
func TestVerifyAttestation_FileNotFound(t *testing.T) {
	err := verifyAttestation("nonexistent_file.cbor")
	if err == nil || !strings.Contains(err.Error(), "file not found") {
		t.Errorf("expected file not found error; got: %v", err)
	}
}

// TestVerifyAttestation_InvalidDocument creates a temporary file with content < 100 bytes.
func TestVerifyAttestation_InvalidDocument(t *testing.T) {
	tmpFile, err := ioutil.TempFile("", "invalid*.cbor")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	// Write an invalid (too short) document.
	if _, err := tmpFile.Write([]byte("short content")); err != nil {
		t.Fatalf("failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	err = verifyAttestation(tmpFile.Name())
	if err == nil || !strings.Contains(err.Error(), "invalid attestation document") {
		t.Errorf("expected invalid document error; got: %v", err)
	}
}

// TestGenerateAttestation_NoEnclaves simulates no running enclaves; if nitro-cli is available, skip the test.
func TestGenerateAttestation_NoEnclaves(t *testing.T) {
	// In real environment this test would depend on external commands.
	t.Skip("Skipping TestGenerateAttestation_NoEnclaves because it requires system-specific nitro-cli output")

	outFile := os.TempDir() + "/dummy_attestation.cbor"
	err := generateAttestation(outFile)
	if err == nil || !strings.Contains(err.Error(), "no running enclaves found") {
		t.Errorf("expected no running enclaves error; got: %v", err)
	}
}

// TestDeployToNitro_InvalidEIFFile tests deployToNitro with a non-existing EIF file.
func TestDeployToNitro_InvalidEIFFile(t *testing.T) {
	err := deployToNitro("nonexistent.eif", 2, 2048, 0)
	if err == nil || !strings.Contains(err.Error(), "EIF file not found") {
		t.Errorf("expected error about missing EIF file; got: %v", err)
	}
}

// TestDeployToNitro_EmptyEIFFile tests deployToNitro with an empty EIF file.
func TestDeployToNitro_EmptyEIFFile(t *testing.T) {
	tmpFile, err := ioutil.TempFile("", "empty*.eif")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	err = deployToNitro(tmpFile.Name(), 2, 2048, 0)
	if err == nil || !strings.Contains(err.Error(), "EIF file is empty") {
		t.Errorf("expected error about empty EIF file; got: %v", err)
	}
}

// TestReplaceSpecialChars tests the replaceSpecialChars function from pack.go
func TestReplaceSpecialChars(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"nginx:alpine", "nginx-alpine"},
		{"ubuntu", "ubuntu"},
		{"registry.io/user/app:v1.0", "registry.io-user-app-v1.0"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := replaceSpecialChars(tt.input)
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

// ...existing tests in cmd_test.go already cover flag checks and command existence.
