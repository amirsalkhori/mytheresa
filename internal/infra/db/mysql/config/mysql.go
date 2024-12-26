package mysql

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"mytheresa/configs"

	"github.com/golang-migrate/migrate/v4"
	mysqlMigrate "github.com/golang-migrate/migrate/v4/database/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type MySQLRepository struct {
	DB *gorm.DB
}

func NewMySQLRepository(cfg *configs.Config, encryptionKey string) (*MySQLRepository, error) {
	// Run migrations
	if err := runMigrations(&cfg.Mysql); err != nil {
		return nil, err
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
		cfg.Mysql.User,
		cfg.Mysql.Pass,
		cfg.Mysql.Host,
		cfg.Mysql.Port,
		cfg.Mysql.Name,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	log.Println("Connected to MySQL using GORM!")

	return &MySQLRepository{
		DB: db,
	}, nil
}

func migrationsPath() string {
	return "file:./migrations"
}

func runMigrations(cfg *configs.Mysql) error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local&multiStatements=true",
		cfg.User,
		cfg.Pass,
		cfg.Host,
		cfg.Port,
		cfg.Name,
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}

	// Create a MySQL instance
	driver, err := mysqlMigrate.WithInstance(db, &mysqlMigrate.Config{})
	if err != nil {

		panic(err)
	}

	// Create a new migration instance
	m, err := migrate.NewWithDatabaseInstance(
		migrationsPath(),
		"mysql", driver)
	if err != nil {
		panic(err)
	}

	// Run migrations
	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		fmt.Println("Error during migration.", err)
		panic(err)
	}

	fmt.Println("Migrations run successfully")

	return nil
}
