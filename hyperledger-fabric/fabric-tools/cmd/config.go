// Copyright © 2018. TIBCO Software Inc.
//
// This file is subject to the license terms contained
// in the license file that is distributed with this file.

package cmd

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Convert IBM Blockchain Flatform network connection file to fabric-sdk-go config file",
	Long:  "Convert IBM Blockchain Flatform network connection file to fabric-sdk-go config file",
	RunE: func(cmd *cobra.Command, args []string) error {
		ibpConfig, _ := cmd.Flags().GetString("ibpConfig")
		if ibpConfig == "" {
			return errors.New("ibpConfig must be specified")
		}
		outConfig, _ := cmd.Flags().GetString("outConfig")
		cryptoPath, _ := cmd.Flags().GetString("cryptoPath")

		return configNetwork(ibpConfig, outConfig, cryptoPath)
	},
}

func init() {
	ibpCmd.AddCommand(configCmd)

	configCmd.Flags().StringP("ibpConfig", "i", "", "Network connection JSON file downloaded from IBM Blockchain Platform")
	configCmd.Flags().StringP("outConfig", "o", "config-ibp.yaml", "Output yaml config file for fabric client")
	configCmd.Flags().StringP("cryptoPath", "p", "./crypto-ibp", "Output folder for crypto data files")
	cobra.MarkFlagRequired(configCmd.Flags(), "ibpConfig")
}
