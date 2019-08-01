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
	Short: "Generate contract metadata and graphql definition from a flogo app json file",
	Long:  "Generate contract metadata and graphql definition from a flogo app json file",
	RunE: func(cmd *cobra.Command, args []string) error {
		appfile, _ := cmd.Flags().GetString("appfile")
		gqlfile, _ := cmd.Flags().GetString("gqlfile")
		if appfile == "" && gqlfile != "" {
			return errors.New("flogo app json file must be specified")
		}

		overridefile, _ := cmd.Flags().GetString("override")
		if overridefile == "" {
			overridefile = "override.json"
		}
		setSchemaOverride(overridefile)

		// generate metadata json file
		metafile, _ := cmd.Flags().GetString("metafile")
		if metafile == "" {
			metafile = "metadata.json"
		}
		if appfile != "" {
			if err := generateMetadata(appfile, metafile); err != nil {
				return err
			}
		}

		// generate graphql file from metadata json file
		if gqlfile == "" {
			gqlfile = "metadata.gql"
		}
		return generateGqlfile(metafile, gqlfile)
	},
}

func init() {
	rootCmd.AddCommand(metadataCmd)
	metadataCmd.Flags().StringP("appfile", "f", "", "Path of the flogo app json file")
	metadataCmd.Flags().StringP("metafile", "m", "", "name of the output metadata file, defaults to metadata.json")
	metadataCmd.Flags().StringP("gqlfile", "g", "", "name of the output graphql file, defaults to metadata.gql")
	metadataCmd.Flags().StringP("override", "o", "", "name of the config file for overriding type ID of schema objects, defaults to override.json")
}
