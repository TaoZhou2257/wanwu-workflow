package db

import (
	"fmt"

	wanwu_db "github.com/UnicomAI/wanwu/pkg/db"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func New(cfg wanwu_db.Config) (*gorm.DB, error) {
	var db *gorm.DB
	var err error
	switch cfg.DBName {
	case "mysql", "tidb", "oceanbase":
		var connCfg wanwu_db.ConnConfig
		switch cfg.DBName {
		case "mysql":
			connCfg = cfg.MySQL
		case "tidb":
			connCfg = cfg.TiDB
		case "oceanbase":
			connCfg = cfg.OceanBase
		}
		db, err = gorm.Open(mysql.Open(fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=%t&loc=%s",
			connCfg.User,
			connCfg.Password,
			connCfg.Address,
			connCfg.Database,
			true,
			// "Asia/Shanghai",
			"Local")), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})
		if err != nil {
			break
		}
		err = setPoolParam(db, connCfg.MaxOpenConns, connCfg.MaxIdleConns)
	case "postgres":
		db, err = gorm.Open(postgres.Open(fmt.Sprintf("postgres://%s:%s@%s/%s",
			cfg.PostgreSQL.User,
			cfg.PostgreSQL.Password,
			cfg.PostgreSQL.Address,
			cfg.PostgreSQL.Database,
		)), &gorm.Config{})
		if err != nil {
			break
		}
		err = setPoolParam(db, cfg.PostgreSQL.MaxOpenConns, cfg.PostgreSQL.MaxIdleConns)
	default:
		err = fmt.Errorf("invalid db %v", cfg.DBName)
	}
	if err != nil {
		return nil, err
	}
	return db, err
}

func setPoolParam(db *gorm.DB, maxOpenConn, maxIdleConn int) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	sqlDB.SetMaxOpenConns(maxOpenConn)
	sqlDB.SetMaxIdleConns(maxIdleConn)
	return nil
}
