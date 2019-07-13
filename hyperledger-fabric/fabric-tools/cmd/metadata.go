// Copyright © 2018. TIBCO Software Inc.
//
// This file is subject to the license terms contained
// in the license file that is distributed with this file.

package cmd

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// metadataCmd generates contract metadata from a flogo app json file
var metadataCmd = &cobra.Command{
	Use:   "metadata",
	Short: "Generate contract metadata from a flogo app json file",
	Long:  "Generate contract metadata from a flogo app json file",
	RunE: func(cmd *cobra.Command, args []string) error {
		appfile, _ := cmd.Flags().GetString("appfile")
		if appfile == "" {
			return errors.New("flogo app json file must be specified")
		}
		outpath, _ := cmd.Flags().GetString("out")
		if outpath == "" {
			outpath = "metadata.json"
		}

		return generateMetadata(appfile, outpath)
	},
}

func init() {
	rootCmd.AddCommand(metadataCmd)

	metadataCmd.Flags().StringP("appfile", "f", "", "Path of the flogo app json file")
	metadataCmd.Flags().StringP("out", "o", "", "name of the output metadata file, defaults to metadata.json")
	cobra.MarkFlagRequired(metadataCmd.Flags(), "appfile")
}
