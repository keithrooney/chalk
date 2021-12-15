package db

import (
	"errors"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type DialectorFactory interface {
	Create(dsn string) gorm.Dialector
}

type SqliteDialectorFactory struct {
}

func (factory SqliteDialectorFactory) Create(dsn string) gorm.Dialector {
	return sqlite.Open(dsn)
}

type DatabaseConfiguration struct {
	DSN string
}

type Database struct {
	DSN     string
	Factory DialectorFactory
}

func NewDatabase(dsn string, factory DialectorFactory) *Database {
	return &Database{
		DSN:     dsn,
		Factory: factory,
	}
}

func (database *Database) Connect() (*gorm.DB, error) {
	return gorm.Open(database.Factory.Create(database.DSN), &gorm.Config{})
}

func (database *Database) Create(model interface{}) error {
	db, err := database.Connect()
	if err != nil {
		return err
	}
	return db.Create(model).Error
}

func (database *Database) Update(model interface{}) error {
	db, err := database.Connect()
	if err != nil {
		return err
	}
	return db.Save(model).Error
}

func (database *Database) Get(id uint, model interface{}) error {
	db, err := database.Connect()
	if err != nil {
		return err
	}
	return db.First(&model, id).Error
}

func (database *Database) Delete(model interface{}) error {
	db, err := database.Connect()
	if err != nil {
		return err
	}
	if db.Delete(model).RowsAffected != 1 {
		return errors.New("record not found")
	}
	return nil
}

func (database *Database) Query(fields map[string]interface{}, model interface{}) error {
	db, err := database.Connect()
	if err != nil {
		return nil
	}
	return db.Where(fields).Find(model).Error
}

func (database *Database) AutoMigrate(models ...interface{}) error {
	db, err := database.Connect()
	if err != nil {
		return err
	}
	return db.AutoMigrate(models...)
}
