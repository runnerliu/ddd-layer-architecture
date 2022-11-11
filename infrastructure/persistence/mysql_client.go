package persistence

import (
	"ddd-demo/common/factory"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

// NewMysqlClient 创建 Mysql 客户端
func NewMysqlClient(factory *factory.ObjectFactory, uri, env string, maxOpenConns, maxIdleConns int) (*gorm.DB, error) {
	client, err := factory.Get(
		uri,
		func(uri string) (interface{}, error) {
			client, err := gorm.Open("mysql", uri)
			if err != nil {
				panic(err)
			} else if client.Error != nil {
				panic(client.Error)
			}

			return client, nil
		},
		func(obj interface{}) error {
			client, ok := obj.(*gorm.DB)
			if !ok {
				return nil
			}

			return client.Close()
		},
	)
	if err != nil {
		return nil, err
	}

	mysqlClient := client.(*gorm.DB)
	mysqlClient.DB().SetMaxOpenConns(maxOpenConns)
	mysqlClient.DB().SetMaxIdleConns(maxIdleConns)
	if env != "prod" {
		mysqlClient.LogMode(true)
	}

	return mysqlClient, nil
}
