package mysql_test

import (
	"github.com/DATA-DOG/go-sqlmock"
	customsql "github.com/jamesparry2/Muzz/app/store/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func setupMockDB() (*customsql.Client, sqlmock.Sqlmock, error) {
	db, mock, err := sqlmock.New()
	if err != nil {
		return nil, mock, err
	}

	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	if err != nil {
		return nil, mock, err
	}

	return customsql.NewClient(gormDB), mock, nil
}
