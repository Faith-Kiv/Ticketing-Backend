package database

import (
	"os"
	"strings"

	"github.com/Faith-Kiv/Ticketing-Backend/models"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"gorm.io/plugin/dbresolver"
)

var DbR *gorm.DB

// var DbReader *gorm.DB

func init() {

	if strings.EqualFold(os.Getenv("DATABASE_URL"), "") {
		os.Setenv("DATABASE_URL", "file::memory:?cache=shared")
		database, err := gorm.Open(sqlite.Open(os.Getenv("DATABASE_URL")), &gorm.Config{
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true,
			},
		})
		if err != nil {
			panic(err)
		}
		DbR = database
		DbR.Table("channels").AutoMigrate(models.Tickets{})
		DbR.Table("screens").AutoMigrate(models.TicketMessage{})

	} else {
		database, err := gorm.Open(postgres.Open(os.Getenv("DATABASE_URL")), &gorm.Config{
			TranslateError: true,
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true,
			},
		})
		if err != nil {
			panic(err)
		}
		err = database.Use(
			dbresolver.Register(dbresolver.Config{
				Replicas:          []gorm.Dialector{postgres.Open(os.Getenv("DATABASE_URL_READER"))},
				Policy:            dbresolver.RandomPolicy{},
				TraceResolverMode: true,
			}),
		)

		if err != nil {
			panic(err)
		}
		DbR = database
		m, err := migrate.New("file://database/migrations", os.Getenv("DATABASE_URL"))
		if err != nil {
			panic(err)
		}

		if err = m.Up(); err != nil && err != migrate.ErrNoChange {
			panic(err)
		}

	}

	logrus.Info("Completed migration")

}
