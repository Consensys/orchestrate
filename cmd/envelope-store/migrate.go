package envelopestore

import (
	"fmt"

	"github.com/go-pg/pg"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/database/postgres"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/services/envelope-store/pg/migrations"
)

// newMigrateCmd create migrate command
func newMigrateCmd() *cobra.Command {
	var db *pg.DB

	migrateCmd := &cobra.Command{
		Use:   "migrate",
		Short: "Migrate database",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			// Set database connection
			db = pg.Connect(postgres.NewOptions())
		},
		Run: func(cmd *cobra.Command, args []string) {
			migrate(db)
		},
		PersistentPostRun: func(cmd *cobra.Command, args []string) {
			err := db.Close()
			if err != nil {
				log.WithError(err).Warn("could not close Postgres connection")
			}
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
		Long:  "Runs all available migrations or up to [target] if argument is provided",
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
		Long:  "Set database version without running migrations",
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
