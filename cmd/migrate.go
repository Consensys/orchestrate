package cmd

import (
	"fmt"

	"github.com/go-pg/pg"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gitlab.com/ConsenSys/client/fr/core-stack/api/context-store.git/pg/migrations"
	"gitlab.com/ConsenSys/client/fr/core-stack/common.git/config"
)

// mewMigrateCmd create migrate command
func mewMigrateCmd() *cobra.Command {
	var db *pg.DB

	migrateCmd := &cobra.Command{
		Use:   "migrate",
		Short: "Migrate database",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			// Configura logger
			config.ConfigureLogger()

			// Set database connection
			db = pg.Connect(&pg.Options{
				Addr:     fmt.Sprintf("%v:%v", viper.GetString("db.host"), viper.GetString("db.port")),
				User:     viper.GetString("db.user"),
				Password: viper.GetString("db.password"),
				Database: viper.GetString("db.database"),
				PoolSize: viper.GetInt("db.poolsize"),
			})
		},
		Run: func(cmd *cobra.Command, args []string) {
			migrate(db)
		},
		PersistentPostRun: func(cmd *cobra.Command, args []string) {
			db.Close()
		},
	}	

	// Register Init command
	initCmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize database",
		RunE: func(cmd *cobra.Command, args []string) error {
			_, _, err := migrations.Run(db, "init")
			if err != nil {
				return err
			}
			fmt.Printf("Database initialized\n")
			return nil
		},
	}
	migrateCmd.AddCommand(initCmd)

	// Register Up command
	upCmd := &cobra.Command{
		Use:   "up [target]",
		Short: "Upgrade database",
		Long: "Runs all available migrations or up to [target] if argument is provided",
		Run: func(cmd *cobra.Command, args []string) {
			migrate(db, append([]string{"up"}, args...)...)
		},
	}
	migrateCmd.AddCommand(upCmd)

	// Register Down command
	downCmd := &cobra.Command{
		Use:   "down",
		Short: "Reverts last migration",
		Run: func(cmd *cobra.Command, args []string) {
			migrate(db, "down")
		},
	}
	migrateCmd.AddCommand(downCmd)

	// Register Reset command
	resetCmd := &cobra.Command{
		Use:   "reset",
		Short: "Reverts all migrations",
		Run: func(cmd *cobra.Command, args []string) {
			migrate(db, "reset")
		},
	}
	migrateCmd.AddCommand(resetCmd)

	// Register Reset command
	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Print current database version",
		RunE: func(cmd *cobra.Command, args []string) error {
			version, _, err := migrations.Run(db, "version")
			if err != nil {
				return err
			}
			fmt.Printf("%v\n", version)
			return nil
		},
	}
	migrateCmd.AddCommand(versionCmd)

	// Register set version command
	setVersionCmd := &cobra.Command{
		Use:   "set-version",
		Short: "Set database version",
		Long: "Set database version without running migrations",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			version, _, err := migrations.Run(db, "set_version", args[0])
			if err != nil {
				return err
			}
			fmt.Printf("%v\n", version)
			return nil
		},
	}

	migrateCmd.AddCommand(setVersionCmd)

	return migrateCmd
}

func migrate(db *pg.DB, a ...string) {
	oldVersion, newVersion, err := migrations.Run(db, a...)
	if err != nil {
		log.WithError(err).Errorf("Migration failed")
		return
	}

	if newVersion != oldVersion {
		log.WithFields(log.Fields{
			"version.old": oldVersion,
			"version.new": newVersion,
		}).Infof("Migration completed")
	} else {
		log.WithFields(log.Fields{
			"version": oldVersion,
		}).Warnf("Nothing to migrate")
	}
}
