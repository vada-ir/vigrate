// App migrate is a sql migration tool for sdp
// Note: this is a tool
package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"database/sql"

	"github.com/codegangsta/cli"
	_ "github.com/lib/pq" // import for init
	"github.com/rubenv/sql-migrate"
	"gopkg.in/yaml.v2"
)

const emptyMigration = `-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied


-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back

`

var (
	migrationTitle = "something"
	step           = 1
	schema         = "public"
	env            = "development"
	configPath     = "dbconfig.yml"
)

// Config is structure of config parameters in yaml under environments
type Config struct {
	DataSourceName string `yaml:"dsn"`
	DriverName     string `yaml:"driver"`
	Dir            string `yaml:"dir"`
}

func main() {
	app := cli.NewApp()
	app.Name = "Migrate"
	app.Usage = "migration tool for sdp"
	app.Version = "0.0.1"

	stepFlag := cli.IntFlag{
		Name:        "step",
		Value:       step,
		EnvVar:      "STEP",
		Usage:       "step to rollback or refresh",
		Destination: &step,
	}

	schemaFlag := cli.StringFlag{
		Name:        "schema",
		Value:       schema,
		EnvVar:      "SCHEMA",
		Usage:       "selected schema for migration",
		Destination: &schema,
	}

	envFlag := cli.StringFlag{
		Name:        "env",
		Value:       env,
		EnvVar:      "ENV",
		Usage:       "migration environment",
		Destination: &env,
	}

	configPathFlag := cli.StringFlag{
		Name:        "config",
		Value:       configPath,
		EnvVar:      "CONFIG",
		Usage:       "migration configuration file",
		Destination: &configPath,
	}

	app.Commands = []cli.Command{
		{
			Name:  "create",
			Usage: "to create a migration",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:        "name",
					Value:       migrationTitle,
					EnvVar:      "NAME",
					Usage:       "name of migration to create",
					Destination: &migrationTitle,
				},
				envFlag,
				configPathFlag,
			},
			Action: func(c *cli.Context) {
				err := createMigration(migrationTitle)
				if err != nil {
					log.Panic(err)
				}
			},
		},
		{
			Name:  "up",
			Usage: "to apply all new migrations",
			Flags: []cli.Flag{
				schemaFlag,
				envFlag,
				configPathFlag,
			},
			Action: func(c *cli.Context) {
				err := doMigrateUp(schema)
				if err != nil {
					log.Panic(err)
				}
			},
		},
		{
			Name:  "rollback",
			Usage: "to rollback migrations",
			Flags: []cli.Flag{
				schemaFlag,
				stepFlag,
				envFlag,
				configPathFlag,
			},
			Action: func(c *cli.Context) {
				err := doMigrateRollback(schema, step)
				if err != nil {
					log.Panic(err)
				}
			},
		},
		{
			Name:  "reset",
			Usage: "to reset database",
			Flags: []cli.Flag{
				schemaFlag,
				envFlag,
				configPathFlag,
			},
			Action: func(c *cli.Context) {
				err := doMigrateReset(schema)
				if err != nil {
					log.Panic(err)
				}
			},
		},
		{
			Name:  "refresh",
			Usage: "to redo migrations",
			Flags: []cli.Flag{
				schemaFlag,
				stepFlag,
				envFlag,
				configPathFlag,
			},
			Action: func(c *cli.Context) {
				err := doMigrateRefresh(schema, step)
				if err != nil {
					log.Panic(err)
				}
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Panic(err)
	}
}

func createMigration(name string) error {
	config, err := getConfig()
	if err != nil {
		return err
	}
	data := []byte(emptyMigration)
	now := time.Now()
	filename := fmt.Sprintf(
		"%s/%04d%02d%02d%02d%02d%02d_%s.sql",
		config.Dir,
		now.Year(),
		now.Month(),
		now.Day(),
		now.Hour(),
		now.Minute(),
		now.Second(),
		name,
	)
	err = ioutil.WriteFile(filename, data, 0644)
	if err == nil {
		log.Printf("migration '%s' successfully created", filename)
	}

	return err
}

func doMigrateUp(schema string) error {
	n, err := doMigrate(schema, migrate.Down, 0)

	if err == nil {
		log.Printf("'%d' migration applied", n)
	}

	return err
}

func doMigrateRollback(schema string, step int) error {
	n, err := doMigrate(schema, migrate.Down, step)

	if err == nil {
		log.Printf("'%d' migration rollbacked", n)
	}

	return err
}

func doMigrateReset(schema string) error {
	n, err := doMigrate(schema, migrate.Down, 0)

	if err == nil {
		log.Printf("'%d' migration rollbacked", n)
	}

	return err
}

func doMigrateRefresh(schema string, step int) error {
	n, err := doMigrate(schema, migrate.Down, step)

	if err == nil {
		log.Printf("'%d' migration rollbacked", n)

		n, err = doMigrate(schema, migrate.Up, step)

		if err == nil {
			log.Printf("'%d' migration applied again", n)
		}
	}

	return err
}

func doMigrate(schema string, dir migrate.MigrationDirection, step int) (int, error) {
	config, err := getConfig()
	if err != nil {
		return 0, err
	}

	db, err := sql.Open(config.DriverName, config.DataSourceName)
	if err != nil {
		return 0, err
	}

	defer func() {
		err := db.Close()
		if err != nil {
			log.Panic(err)
		}
	}()

	migrations := &migrate.FileMigrationSource{
		Dir: config.Dir,
	}

	migrate.SetSchema(schema)

	return migrate.ExecMax(db, "postgres", migrations, dir, step)
}

func getConfig() (*Config, error) {
	config := make(map[string]*Config)
	config[env] = &Config{
		DriverName:     "postgres",
		DataSourceName: "postgres://postgres:@localhost/postgres?sslmode=verify-full",
		Dir:            "db/migrations",
	}

	file, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(file, &config)
	if err != nil {
		return nil, err
	}

	return config[env], nil
}
