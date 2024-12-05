package repository

import (
	"fmt"
	"go-admin/internal/lib/driver/mysql"
	"go-admin/internal/lib/env"

	"gorm.io/gorm"
)

type Base struct {
	DriverKey string
	Table     string
}

// GetDBSession
func (repo *Base) GetDBSession() (*gorm.DB, error) {
	conn, err := mysql.DBInstance(repo.DriverKey, !env.IsProd())
	if err != nil {
		return nil, fmt.Errorf("get mysql db instance: %w", err)
	}
	return conn.Table(repo.Table), nil
}

// Insert
func (repo *Base) Insert(values map[string]interface{}) (int64, error) {
	dbSession, err := repo.GetDBSession()
	if err != nil {
		return 0, err
	}
	dbSession.Create(values)
	return dbSession.RowsAffected, dbSession.Error
}

// Update
func (repo *Base) Update(values map[string]interface{}, condition map[string]interface{}) (int64, error) {
	dbSession, err := repo.GetDBSession()
	if err != nil {
		return 0, err
	}
	dbSession.Where(condition)
	dbSession.Updates(values)
	return dbSession.RowsAffected, dbSession.Error
}

// Delete
func (repo *Base) Delete(condition map[string]interface{}) (int64, error) {
	dbSession, err := repo.GetDBSession()
	if err != nil {
		return 0, err
	}
	dbSession.Where(condition)
	dbSession.Delete(map[string]interface{}{})
	return dbSession.RowsAffected, dbSession.Error
}

// GET one record
func (repo *Base) Get(cols string, condition map[string]interface{}) (map[string]interface{}, error) {
	dbSession, err := repo.GetDBSession()
	if err != nil {
		return nil, err
	}

	if len(cols) > 0 {
		dbSession.Select(cols)
	} else {
		dbSession.Select("*")
	}

	dbSession.Where(condition)

	result := map[string]interface{}{}
	dbSession.Take(&result)

	return result, dbSession.Error
}

// GetList all records
func (repo *Base) GetList(cols string, condition map[string]interface{}, orderBy string) ([]map[string]interface{}, error) {
	dbSession, err := repo.GetDBSession()
	if err != nil {
		return nil, err
	}

	if len(cols) > 0 {
		dbSession.Select(cols)
	} else {
		dbSession.Select("*")
	}

	dbSession.Where(condition)

	if len(orderBy) > 0 {
		dbSession.Order(orderBy)
	} else {
		dbSession.Order("id asc")
	}

	result := []map[string]interface{}{}
	dbSession.Find(&result)

	return result, dbSession.Error
}

// Exist
func (repo *Base) Exist(condition map[string]interface{}) (bool, error) {
	dbSession, err := repo.GetDBSession()
	if err != nil {
		return false, err
	}

	var exist bool
	err = dbSession.Where(condition).Select("count(*) > 0").Find(&exist).Error
	if err != nil {
		return false, err
	}

	return exist, err
}

// Count
func (repo *Base) Count(condition map[string]interface{}) (int64, error) {
	dbSession, err := repo.GetDBSession()
	if err != nil {
		return 0, err
	}

	var count int64
	err = dbSession.Where(condition).Count(&count).Error
	if err != nil {
		return 0, err
	}

	return count, err
}

// QuerySqlCount
func (repo *Base) QuerySqlCount(sql string) (int64, error) {
	dbSession, err := repo.GetDBSession()
	if err != nil {
		return 0, err
	}

	var count int64
	err = dbSession.Raw(sql).Count(&count).Error
	if err != nil {
		return 0, err
	}

	return count, err
}

// ExecSql
func (repo *Base) ExecSql(sql string) (int64, error) {
	dbSession, err := repo.GetDBSession()
	if err != nil {
		return 0, err
	}

	dbSession.Exec(sql)

	return dbSession.RowsAffected, dbSession.Error
}

// GetByID
func (repo *Base) GetByID(ID int, dest interface{}) error {
	db, err := repo.GetDBSession()
	if err != nil {
		return err
	}

	db.Where("id=?", ID).First(&dest)
	return db.Error
}

// GetByUID
func (repo *Base) GetByUID(ID int, dest interface{}) error {
	db, err := repo.GetDBSession()
	if err != nil {
		return err
	}

	db.Where("uid=?", ID).First(&dest)
	return db.Error
}

// FindByID
func (repo *Base) FindByID(ID int, dest interface{}) error {
	db, err := repo.GetDBSession()
	if err != nil {
		return err
	}

	db.Where("id=?", ID).Find(&dest)
	return db.Error
}

// CheckNameExists
func (repo *Base) CheckNameExists(ID int, name string) (bool, error) {
	db, err := repo.GetDBSession()
	if err != nil {
		return false, err
	}

	var exist bool
	if ID == 0 {
		// create
		err = db.Where("name = ? AND is_delete = ?", name, 0).Select("count(*) > 0").Find(&exist).Error
	} else {
		// update
		err = db.Where("id <> ? AND name = ? AND is_delete = ?", ID, name, 0).Select("count(*) > 0").Find(&exist).Error
	}
	if err != nil {
		return false, err
	}

	return exist, err
}
