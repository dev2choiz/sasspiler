package cmd

import (
	"github.com/dev2choiz/sasspiler/transpiler"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"log"
	"os"
	"strings"
)

var (
	// Contains source flag
	// only available in root command
	sourceDir string

	// Contains dest flag
	// only available in root command
	destDir string

	// Contains impDir flag
	// only available in root command
	impDir string

	// Contains verbose flag
	// available also in sub commands
	verbose bool

	// Cobra main command
	rootCmd = &cobra.Command{
		Use:   "sasspiler",
		Short: "sasspiler",
		Long:  `sasspiler`,
		Run: func(cmd *cobra.Command, args []string) {
			checkFlags()
			t := transpiler.New()
			t.SetVerbose(verbose)

			// use real filesystem
			fs := afero.NewOsFs()
			t.SetFileSystem(fs)
			files := t.GetFiles(sourceDir)
			imp := make([]string, 0)
			if impDir != "" {
				imp = strings.Split(impDir, ",")
			}
			err := t.Run(sourceDir, destDir, files, imp)
			if err != nil {
				log.Fatalln(err)
			}
		},
	}
)

// define flags and add cobra sub commands
func Execute() error {
	rootCmd.Flags().StringVarP(&sourceDir, "source", "s", "", "scss source directory")
	rootCmd.Flags().StringVarP(&destDir, "dest", "d", "", "directory where css files will be generated")
	rootCmd.Flags().StringVarP(&impDir, "importDir", "i", "", "directory where are located imported scss files")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbosity")
	rootCmd.AddCommand(versionCmd)
	return rootCmd.Execute()
}

// check rootCmd flags
func checkFlags() {
	if sourceDir == "" {
		log.Fatalln("--source flag is required")
	}
	if destDir == "" {
		log.Fatalln("--dest flag is required")
	}

	if _, err := os.Stat(sourceDir); os.IsNotExist(err) {
		log.Fatalln(sourceDir, "directory does not exist")
	}
}
