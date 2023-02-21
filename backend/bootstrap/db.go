package bootstrap

import (
	"fmt"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"io"
	"log"
	"monaToolBox/global"
	"monaToolBox/models"
	"os"
	"strconv"
	"time"
)

// 初始化数据库
func InitDatabase() *gorm.DB {
	switch global.Config.Database.Driver {
	case "mysql", "mariadb":
		return initMysqlGorm()
	default:
		global.Log.Fatal(fmt.Sprintf("sql derive not support. input database.driver:%s", global.Config.Database.Driver))
		return nil
	}
}

func initMysqlGorm() *gorm.DB {
	dbConfig := global.Config.Database

	if dbConfig.Schema == "" {
		global.Log.Fatal("mysql schema not define.")
		return nil
	}

	dsn := dbConfig.UserName + ":" + dbConfig.Password + "@tcp(" + dbConfig.Host + ":" + strconv.Itoa(dbConfig.Port) + ")/" +
		dbConfig.Schema + "?charset=" + dbConfig.Charset + "&parseTime=True&loc=Local"

	mysqlConfig := mysql.Config{
		DSN:                       dsn,   // DSN data source name
		DefaultStringSize:         256,   // string 类型字段的默认长度
		DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false, // 根据版本自动配置
	}

	if db, err := gorm.Open(
		mysql.New(mysqlConfig), &gorm.Config{
			DisableForeignKeyConstraintWhenMigrating: true,            // 禁用自动创建外键约束
			Logger:                                   getGormLogger(), // 使用自定义 Logger
		},
	); err != nil {
		global.Log.Error("mysql connect failed, err:", zap.Any("err", err))
		return nil
	} else {
		sqlDB, _ := db.DB()
		sqlDB.SetMaxIdleConns(dbConfig.MaxIdleConns) // 空闲连接池中连接的最大数量
		sqlDB.SetMaxOpenConns(dbConfig.MaxOpenConns) // 打开数据库连接的最大数量
		initDbTables(db)
		return db
	}

}

func initDbTables(db *gorm.DB) {
	if err := db.Set("gorm:table_options", "ENGINE=InnoDB CHARSET=utf8mb4").AutoMigrate(
		models.User{},
	); err != nil {
		global.Log.Error("migrate tables failed.", zap.Any("err", err))
		os.Exit(-1)
	}
}

// getGormLogWriter 获取一个自定义的logger
func getGormLogWriter() logger.Writer {
	var writer io.Writer

	if global.Config.Database.EnableFileLogWriter {
		writer = &lumberjack.Logger{
			Filename:   global.Config.Log.RootDir + "/" + global.Config.Database.LogFilename,
			MaxSize:    global.Config.Log.MaxSize,
			MaxBackups: global.Config.Log.MaxBackups,
			MaxAge:     global.Config.Log.MaxAge,
			Compress:   global.Config.Log.Compress,
		}
	} else {
		writer = os.Stdout
	}

	return log.New(writer, "\r\n", log.LstdFlags)
}

// getFormLogger 定义一个gorm的logger 并返回
func getGormLogger() logger.Interface {
	var logMode logger.LogLevel

	switch global.Config.Database.LogMode {
	case "silent":
		logMode = logger.Silent
	case "error":
		logMode = logger.Error
	case "warn":
		logMode = logger.Warn
	case "info":
		logMode = logger.Info
	default:
		logMode = logger.Info
	}

	return logger.New(
		getGormLogWriter(), logger.Config{
			SlowThreshold:             200 * time.Millisecond,                      // 慢 SQL 阈值
			LogLevel:                  logMode,                                     // 日志级别
			IgnoreRecordNotFoundError: false,                                       // 忽略ErrRecordNotFound（记录未找到）错误
			Colorful:                  !global.Config.Database.EnableFileLogWriter, // 禁用彩色打印
		},
	)

}
