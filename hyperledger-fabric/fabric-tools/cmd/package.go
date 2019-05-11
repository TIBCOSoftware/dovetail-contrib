// Copyright © 2018. TIBCO Software Inc.
//
// This file is subject to the license terms contained
// in the license file that is distributed with this file.

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// packageCmd represents the package command
var packageCmd = &cobra.Command{
	Use:   "package",
	Short: "Package chaincode to cds format",
	Long:  "Package chaincode to cds format",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("not implemented. you can create package using 'docker exec cli peer chaincode package'")
	},
}

func init() {
	rootCmd.AddCommand(packageCmd)
}
