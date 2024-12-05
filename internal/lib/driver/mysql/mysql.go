package mysql

import (
	"fmt"
	"go-admin/internal/lib/config"
	"reflect"
	"sync"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type driver struct {
	clients sync.Map
}

var driverObj *driver

func DBInstance(key string, isShowSQL bool) (*gorm.DB, error) {
	if driverObj == nil {
		driverObj = &driver{}
	}

	// get cache
	dbClient := driverObj.getClient(key)
	if dbClient != nil {
		return dbClient, nil
	}

	// new db
	mysqlConfigInterface := GetFieldByTagName(config.Settings.MySQL, "mapstructure", key)
	if mysqlConfigInterface == nil {
		return nil, fmt.Errorf("DBInstance error, can't get mysql key=%s config", key)
	}
	mysqlConfig := mysqlConfigInterface.(config.MySQLConfig)
	db, err := NewDB(&mysqlConfig, isShowSQL)
	if err != nil {
		return nil, fmt.Errorf("DBInstance error: %v", err)
	}

	// set cache
	driverObj.cacheClient(key, db)

	return db, nil
}

func NewDB(conf *config.MySQLConfig, isShowSQL bool) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(conf.Conn[0]), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return db, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return db, err
	}
	if conf.MaxIdle != 0 {
		sqlDB.SetMaxIdleConns(conf.MaxIdle)
	}
	if conf.MaxOpen != 0 {
		sqlDB.SetMaxOpenConns(conf.MaxOpen)
	}
	if conf.MaxLifetime != 0 {
		sqlDB.SetConnMaxLifetime(conf.MaxLifetime)
	}

	if isShowSQL {
		db.Debug()
	} else {
		db.Logger = logger.Default.LogMode(logger.Silent)
	}

	return db, nil
}

func (d *driver) getClient(key string) *gorm.DB {
	db, exist := d.clients.Load(key)
	if exist {
		return db.(*gorm.DB)
	}
	return nil
}

func (d *driver) cacheClient(key string, e *gorm.DB) {
	d.clients.Store(key, e)
}

func GetFieldByTagName(obj interface{}, tagName, targetTag string) interface{} {
	objValue := reflect.ValueOf(obj)
	if objValue.Kind() == reflect.Ptr {
		objValue = objValue.Elem()
	}

	objType := objValue.Type()
	for i := 0; i < objType.NumField(); i++ {
		field := objType.Field(i)
		if tagValue, ok := field.Tag.Lookup(tagName); ok && tagValue == targetTag {
			return objValue.Field(i).Interface()
		}
	}

	return nil
}
