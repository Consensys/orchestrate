package scripts

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/go-pg/pg/v9"
	"os"
)

type columnSchema struct {
	ColumnName string
	DataType   string
}

func MigrateAPIDB(apiDB *pg.DB) error {
	var tables []string
	switch os.Getenv("DB_MIGRATION_SERVICE") {
	case "contract-registry":
		tables = []string{"repositories", "artifacts", "tags", "events", "methods", "codehashes"}
	case "transaction-scheduler":
		tables = []string{"schedules", "transaction_requests", "transactions", "jobs", "logs"}
	case "chain-registry":
		tables = []string{"chains", "faucets", "private_tx_managers"}
	default:
		return errors.New("unknown service")
	}

	return migrate(apiDB, tables)
}

func migrate(apiDB *pg.DB, tables []string) error {
	oldDB := pg.Connect(&pg.Options{
		User:     os.Getenv("DB_MIGRATION_USERNAME"),
		Password: os.Getenv("DB_MIGRATION_PASSWORD"),
		Database: os.Getenv("DB_MIGRATION_DATABASE"),
		Addr:     os.Getenv("DB_MIGRATION_ADDRESS"),
	})
	defer oldDB.Close()

	for _, tableName := range tables {
		err := compareSchemas(oldDB, apiDB, tableName)
		if err != nil {
			return err
		}

		err = dumpToFile(oldDB, tableName)
		if err != nil {
			return err
		}

		err = importFromFile(apiDB, tableName)
		if err != nil {
			return err
		}

		err = os.Remove(getFileName(tableName))
		if err != nil {
			return err
		}
	}

	return nil
}

func compareSchemas(oldDB *pg.DB, apiDB *pg.DB, tableName string) error {
	oldSchema, err := getSchema(oldDB, tableName)
	if err != nil {
		return err
	}

	newSchema, err := getSchema(apiDB, tableName)
	if err != nil {
		return err
	}

	for i, column := range oldSchema {
		if column.ColumnName != newSchema[i].ColumnName && column.DataType != newSchema[i].DataType {
			return errors.New("schemas differ. aborting DB copy")
		}
	}

	return nil
}

func getSchema(db *pg.DB, tableName string) ([]*columnSchema, error) {
	var schema []*columnSchema
	_, err := db.Query(
		&schema,
		"SELECT column_name, data_type FROM information_schema.columns WHERE table_name = ?",
		tableName,
	)
	if err != nil {
		return nil, err
	}

	return schema, nil
}

func dumpToFile(oldDB *pg.DB, tableName string) error {
	out, err := os.Create(getFileName(tableName))
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = oldDB.CopyTo(out, fmt.Sprintf(`COPY %s TO STDOUT WITH CSV`, tableName))
	if err != nil {
		return err
	}

	return nil
}

func importFromFile(apiDB *pg.DB, tableName string) error {
	f, err := os.Open(getFileName(tableName))
	if err != nil {
		return err
	}
	r := bufio.NewReader(f)

	_, err = apiDB.CopyFrom(r, fmt.Sprintf(`COPY %s FROM STDIN WITH CSV`, tableName))
	if err != nil {
		return err
	}

	return nil
}

func getFileName(tableName string) string {
	return os.TempDir() + "dump-" + tableName + ".csv"
}
