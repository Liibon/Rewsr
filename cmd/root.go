package cmd

import (
	"github.com/spf13/cobra"
)

const banner = `
 ██████╗ ███████╗██╗    ██╗███████╗██████╗ 
 ██╔══██╗██╔════╝██║    ██║██╔════╝██╔══██╗
 ██████╔╝█████╗  ██║ █╗ ██║█████╗  ██████╔╝
 ██╔═══╝ ██╔══╝  ██║███╗██║██╔══╝  ██╔═══╝ 
 ██║     ███████╗╚███╔███╔╝███████╗██║     
 ╚═╝     ╚══════╝ ╚══╝╚══╝ ╚══════╝╚═╝     
                                           
CLI tool for AWS Nitro Enclaves
`

var rootCmd = &cobra.Command{
	Use:   "rewsr",
	Short: "CLI tool for AWS Nitro Enclaves",
	Long: `REWSR builds EIF files from any Docker image and runs them 
in AWS Nitro Enclaves with hardware attestation.

Requires AWS Nitro-enabled EC2 instance (e.g., t3.large).

Tips for getting started:
 1. Launch t3.large instance with Nitro Enclaves enabled
 2. Install: sudo yum install aws-nitro-enclaves-cli docker
 3. Run: rewsr pack nginx:alpine && rewsr deploy nginx-alpine.eif`,
	Version: "0.2.0-nitro-preview",
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(packCmd)
	rootCmd.AddCommand(deployCmd)
	rootCmd.AddCommand(attestCmd)
	rootCmd.AddCommand(verifyCmd)
}