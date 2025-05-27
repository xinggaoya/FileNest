package config

import "time"

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Host            string        `mapstructure:"host"`
	Port            int           `mapstructure:"port"`
	Username        string        `mapstructure:"username"`
	Password        string        `mapstructure:"password"`
	Database        string        `mapstructure:"database"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`
	MaxOpenConns    int           `mapstructure:"max_open_conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
}

// 默认数据库配置
var Database = &DatabaseConfig{
	Host:            "moncn.cn",
	Port:            8206,
	Username:        "xinggao",
	Password:        "xinggaocode",
	Database:        "db_filenest",
	MaxIdleConns:    10,
	MaxOpenConns:    100,
	ConnMaxLifetime: time.Hour,
}
