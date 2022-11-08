package cmd

import (
	"fmt"
	"os"

	"github.com/giantas/gotu/indexer"
	"github.com/giantas/gotu/storage"
	"github.com/spf13/cobra"
)

var (
	dbName string = "gotu.db"
	initDb bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gotu",
	Short: "A file indexer",
	Long:  `A long description of a file indexer with a cli`,
}

var dbCmd = &cobra.Command{
	Use:   "db",
	Short: "Database configuration",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := storage.StoreConfig{
			Init: initDb,
			URI:  dbName,
		}
		db, err := storage.ConnectDatabase(cfg)
		if err != nil {
			exitWithError(err)
		}
		defer db.Close()
	},
}

var indexCmd = &cobra.Command{
	Use:   "index",
	Short: "Indexer",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := storage.StoreConfig{
			Init: initDb,
			URI:  dbName,
		}
		db, err := storage.ConnectDatabase(cfg)
		if err != nil {
			exitWithError(err)
		}
		defer db.Close()

		fileStore := storage.NewFileStore(db)
		err = indexer.Run(fileStore)
		if err != nil {
			exitWithError(err)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// DB command
	rootCmd.AddCommand(dbCmd)
	dbCmd.Flags().BoolVar(&initDb, "init", false, "Initialise the database")

	// Index command
	rootCmd.AddCommand(indexCmd)
}

func exitWithError(err error) {
	fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	os.Exit(1)
}
