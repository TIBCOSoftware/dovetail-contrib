// Copyright © 2018. TIBCO Software Inc.
//
// This file is subject to the license terms contained
// in the license file that is distributed with this file.

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// enrollCmd represents the enroll command
var enrollCmd = &cobra.Command{
	Use:   "enroll",
	Short: "Enroll a client user and download key and certs from CA server",
	Long:  "Enroll a client user and download key and certs from CA server",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("not implemented. enroll user using 'fabric-ca-client enroll'")
	},
}

func init() {
	ibpCmd.AddCommand(enrollCmd)
}
