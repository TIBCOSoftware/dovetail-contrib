// Copyright © 2018. TIBCO Software Inc.
//
// This file is subject to the license terms contained
// in the license file that is distributed with this file.

package cmd

import (
	"github.com/spf13/cobra"
)

// ibpCmd represents the ibp command
var ibpCmd = &cobra.Command{
	Use:   "ibp",
	Short: "Commands for IBM Blockchain Platform",
	Long:  "Commands for IBM Blockchain Platform",
}

func init() {
	rootCmd.AddCommand(ibpCmd)
}
