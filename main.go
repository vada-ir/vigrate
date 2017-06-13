// App vigrate is a sql migration tool
// Note: this is a tool
package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"sync"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	_ "github.com/lib/pq"
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

var configMap = make(map[string]*Config)

var once sync.Once

func main() {
	app := cli.NewApp()
	app.Name = "Vigrate"
	app.Usage = "migration tool"
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
					logrus.Panic(err)
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
					logrus.Panic(err)
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
					logrus.Panic(err)
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
					logrus.Panic(err)
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
					logrus.Panic(err)
				}
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		logrus.Panic(err)
	}
}

func createMigration(name string) error {
	config := getConfig()

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
	err := ioutil.WriteFile(filename, data, 0644)
	if err == nil {

		logrus.Printf("migration '%s' successfully created", filename)
	}

	return err
}

func doMigrateUp(schema string) error {
	config := getConfig()
	n, err := doMigrate(schema, config.Dir, migrate.Up, 0)

	if err != nil {
		return err
	}

	logrus.Printf("'%d' migration applied", n)

	return nil
}

func doMigrateRollback(schema string, step int) error {
	config := getConfig()
	n, err := doMigrate(schema, config.Dir, migrate.Down, step)

	if err != nil {
		return err
	}

	logrus.Printf("'%d' migration rollbacked", n)

	return err
}

func doMigrateReset(schema string) error {
	config := getConfig()
	n, err := doMigrate(schema, config.Dir, migrate.Down, 0)

	if err != nil {
		return err
	}

	logrus.Printf("'%d' migration rollbacked", n)

	return err
}

func doMigrateRefresh(schema string, step int) error {
	config := getConfig()
	n, err := doMigrate(schema, config.Dir, migrate.Down, step)

	if err != nil {
		return err
	}

	logrus.Printf("'%d' migration rollbacked", n)

	n, err = doMigrate(schema, config.Dir, migrate.Up, step)

	if err != nil {
		return err
	}

	logrus.Printf("'%d' migration applied again", n)

	return nil
}

func doMigrate(schema, path string, dir migrate.MigrationDirection, step int) (int, error) {
	config := getConfig()

	db, err := sql.Open(config.DriverName, config.DataSourceName)
	if err != nil {
		return 0, err
	}

	defer func() {
		err := db.Close()
		if err != nil {
			logrus.Panic(err)
		}
	}()

	migrations := &migrate.FileMigrationSource{
		Dir: path,
	}

	migrate.SetSchema(schema)

	return migrate.ExecMax(db, "postgres", migrations, dir, step)
}

func loadConfig() {
	once.Do(func() {
		configMap[env] = &Config{
			DriverName:     "postgres",
			DataSourceName: "postgres://postgres:@localhost/postgres?sslmode=verify-full",
			Dir:            "db/migrations",
		}

		file, err := ioutil.ReadFile(configPath)
		if err != nil {
			logrus.Fatal(err)
			return
		}

		err = yaml.Unmarshal(file, &configMap)
		if err != nil {
			logrus.Print(err)
			return
		}
	})
}

func getConfig() *Config {
	loadConfig()
	return configMap[env]
}
