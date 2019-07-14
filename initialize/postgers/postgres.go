package postgers

import (
	"account-sync/models"
	"fmt"
	"github.com/go-errors/errors"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/xo/dburl"
	"os"
)

const DefaultUrl = "postgres://mbympgbxovcaec:5c254085dca2140af8553b3c941abe44b47f7569e63d782c8db52b3e40970205@ec2-54-228-246-214.eu-west-1.compute.amazonaws.com:5432/d5idksj6ro3iuo"

var AccountDatabase *gorm.DB

func init() {
	url := os.Getenv("DATABASE_URL")
	if url == "" {
		url = DefaultUrl
	}
	AccountDatabase = connectToPostgres(url)

	migrateStructure()
}

func connectToPostgres(url string) *gorm.DB {
	var err error
	c, err := dburl.Parse(url)
	if err != nil {
		panic(err)
	}

	password, isSet := c.User.Password()
	if !isSet {
		panic(errors.New("password not set for postgres connection"))
	}

	uri := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s", c.Host, c.Port(), c.User.Username(), c.Scheme, password)

	fmt.Println(fmt.Sprintf("Connecting to mysql... "))
	client, err := gorm.Open("postgres", uri)
	if err != nil {
		panic(err)
	}
	fmt.Println(fmt.Sprintf("[success]\n"))
	return client
}

func migrateStructure() {
	AccountDatabase.Begin()
	AccountDatabase.AutoMigrate(&models.Account{})
	AccountDatabase.Commit()
}
