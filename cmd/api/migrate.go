package api

import (
	"os"

	"github.com/go-pg/pg/v9"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/cmd/api/scripts"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/database/postgres"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/store/postgres/migrations"
)

// newMigrateCmd create migrate command
func newMigrateCmd() *cobra.Command {
	var db *pg.DB

	migrateCmd := &cobra.Command{
		Use:   "migrate",
		Short: "Migrate database",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// Set database connection
			opts, err := postgres.NewConfig(viper.GetViper()).PGOptions()
			if err != nil {
				return err
			}
			db = pg.Connect(opts)
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return migrate(db)
		},
		PersistentPostRun: func(cmd *cobra.Command, args []string) {
			err := db.Close()
			if err != nil {
				log.WithError(err).Warn("could not close Postgres connection")
			}
		},
	}

	// Postgres flags
	postgres.PGFlags(migrateCmd.Flags())

	// Register Init command
	initCmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize database",
		RunE: func(cmd *cobra.Command, args []string) error {
			_, _, err := migrations.Run(db, "init")
			if err != nil {
				return err
			}
			log.Infof("Database initialized")
			return nil
		},
	}
	migrateCmd.AddCommand(initCmd)

	// Register Up command
	upCmd := &cobra.Command{
		Use:   "up [target]",
		Short: "Upgrade database",
		Long:  "Runs all available migrations or up to [target] if argument is provided",
		RunE: func(cmd *cobra.Command, args []string) error {
			return migrate(db, append([]string{"up"}, args...)...)
		},
	}
	migrateCmd.AddCommand(upCmd)

	// Register Down command
	downCmd := &cobra.Command{
		Use:   "down",
		Short: "Reverts last migration",
		RunE: func(cmd *cobra.Command, args []string) error {
			return migrate(db, "down")
		},
	}
	migrateCmd.AddCommand(downCmd)

	// Register Reset command
	resetCmd := &cobra.Command{
		Use:   "reset",
		Short: "Reverts all migrations",
		RunE: func(cmd *cobra.Command, args []string) error {
			return migrate(db, "reset")
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
			log.Infof("%v", version)
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
			log.Infof("%v", version)
			return nil
		},
	}
	migrateCmd.AddCommand(setVersionCmd)

	// Register Migrate API DB
	migrateAPIDBCmd := &cobra.Command{
		Use:   "copy-db",
		Short: "Copy Database from version 2.5.x to version 21.1.x",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Set database connection
			opts, err := postgres.NewConfig(viper.GetViper()).PGOptions()
			if err != nil {
				return err
			}

			opts.User = os.Getenv("DB_MIGRATION_USERNAME")
			opts.Password = os.Getenv("DB_MIGRATION_PASSWORD")
			opts.Database = os.Getenv("DB_MIGRATION_DATABASE")
			opts.Addr = os.Getenv("DB_MIGRATION_ADDRESS")

			oldDB := pg.Connect(opts)

			return scripts.MigrateAPIDB(db, oldDB)
		},
	}
	migrateCmd.AddCommand(migrateAPIDBCmd)

	return migrateCmd
}

func migrate(db *pg.DB, a ...string) error {
	oldVersion, newVersion, err := migrations.Run(db, a...)
	if err != nil {
		log.WithError(err).Errorf("Migration failed")
		return err
	}

	log.WithField("version", newVersion).WithField("previous_version", oldVersion).Info("All migrations completed")
	return nil
}
