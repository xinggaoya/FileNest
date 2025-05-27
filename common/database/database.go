package database

import (
	"FileNest/common/glog"
	"FileNest/internal/config"
	"FileNest/internal/model"
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// InitDB 初始化数据库连接
func InitDB() error {
	cfg := config.Database

	// 构建DSN
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Database)

	glog.Infof("正在连接数据库: %s@%s:%d/%s", cfg.Username, cfg.Host, cfg.Port, cfg.Database)

	// 配置GORM
	config := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent), // 生产环境使用Silent，开发环境可以用Info
	}

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), config)
	if err != nil {
		return fmt.Errorf("连接数据库失败: %v", err)
	}

	// 配置连接池
	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("获取数据库实例失败: %v", err)
	}

	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	// 测试连接
	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("数据库连接测试失败: %v", err)
	}

	glog.Infof("数据库连接成功")

	// 自动迁移
	if err := autoMigrate(); err != nil {
		return fmt.Errorf("数据库迁移失败: %v", err)
	}

	glog.Infof("数据库初始化完成")
	return nil
}

// autoMigrate 自动迁移数据库表结构
func autoMigrate() error {
	glog.Info("开始数据库表结构迁移...")

	// 迁移所有模型
	err := DB.AutoMigrate(
		&model.Favorite{},
		&model.UploadLog{},
		&model.FileShare{},
	)

	if err != nil {
		return fmt.Errorf("自动迁移失败: %v", err)
	}

	glog.Info("数据库表结构迁移完成")
	return nil
}

// GetDB 获取数据库实例
func GetDB() *gorm.DB {
	return DB
}

// Close 关闭数据库连接
func Close() error {
	if DB != nil {
		sqlDB, err := DB.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}

// Transaction 执行事务
func Transaction(fn func(tx *gorm.DB) error) error {
	return DB.Transaction(fn)
}

// IsTableEmpty 检查表是否为空
func IsTableEmpty(tableName string) (bool, error) {
	var count int64
	err := DB.Table(tableName).Count(&count).Error
	return count == 0, err
}

// TruncateTable 清空表数据
func TruncateTable(tableName string) error {
	return DB.Exec(fmt.Sprintf("TRUNCATE TABLE %s", tableName)).Error
}

// BackupTable 备份表数据（创建备份表）
func BackupTable(tableName string) error {
	backupTableName := fmt.Sprintf("%s_backup_%s", tableName, time.Now().Format("20060102150405"))
	sql := fmt.Sprintf("CREATE TABLE %s AS SELECT * FROM %s", backupTableName, tableName)
	return DB.Exec(sql).Error
}
