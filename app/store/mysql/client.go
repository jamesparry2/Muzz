package mysql

import (
	"fmt"

	"github.com/jamesparry2/Muzz/app/store"
	officalmysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type DB interface {
	Find(dest interface{}, conds ...interface{}) (tx *gorm.DB)
	Save(value interface{}) (tx *gorm.DB)
	Where(query interface{}, args ...interface{}) (tx *gorm.DB)
	Preload(query string, args ...interface{}) (tx *gorm.DB)
}

type ClientOptions struct {
	Host     string
	Port     string
	Username string
	Password string
	Database string
}

type Client struct {
	db DB
}

func NewClientWithConection(opts *ClientOptions) (*Client, error) {
	if opts.Database == "" || opts.Port == "" || opts.Username == "" || opts.Password == "" || opts.Host == "" {
		return nil, store.ErrNewClientMissingRequiredOptions
	}

	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", opts.Username, opts.Password, opts.Host, opts.Port, opts.Database)

	db, err := gorm.Open(officalmysql.Open(connectionString), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(&store.User{}, &store.Location{}, &store.Preferences{}, &store.Swipe{}); err != nil {
		return nil, err
	}

	return NewClient(db), nil
}

func NewClient(db DB) *Client {
	return &Client{db: db}
}
