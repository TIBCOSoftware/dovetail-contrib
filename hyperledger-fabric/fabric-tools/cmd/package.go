// Copyright © 2018. TIBCO Software Inc.
//
// This file is subject to the license terms contained
// in the license file that is distributed with this file.

package cmd

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// packageCmd represents the package command
// equivalent to 'docker exec cli peer chaincode package'
var packageCmd = &cobra.Command{
	Use:   "package",
	Short: "Package chaincode to cds format",
	Long:  "Package chaincode to cds format",
	RunE: func(cmd *cobra.Command, args []string) error {
		name, _ := cmd.Flags().GetString("name")
		if name == "" {
			return errors.New("chaincode name must be specified")
		}
		path, _ := cmd.Flags().GetString("path")
		if path == "" {
			return errors.New("chaincode source path must be specified")
		}

		version, _ := cmd.Flags().GetString("version")
		outFile, _ := cmd.Flags().GetString("out")
		if outFile == "" {
			outFile = name + ".cds"
		}

		return packageCDS(path, name, version, outFile)
	},
}

func init() {
	rootCmd.AddCommand(packageCmd)

	packageCmd.Flags().StringP("name", "n", "", "Name of the chaincode")
	packageCmd.Flags().StringP("version", "v", "0", "version of the chaincode")
	packageCmd.Flags().StringP("path", "p", "", "path of the source code, it must contain a folder \"src\"")
	packageCmd.Flags().StringP("out", "o", "", "name of the output cds file, defaults to <chaincode name>.cds")
	cobra.MarkFlagRequired(packageCmd.Flags(), "name")
	cobra.MarkFlagRequired(packageCmd.Flags(), "path")
}
