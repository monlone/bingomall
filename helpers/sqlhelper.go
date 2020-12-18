package helper

import (
	"fmt"
	"bingomall/constant"
	"bingomall/system"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"os"
	"time"
)

//var DBConnect map[string]*gorm.DB
var DBConnect = make(map[string]*gorm.DB)

// 初始化连接
func init() {
	errEnv := godotenv.Load()
	if errEnv != nil {
		ErrorLogger.Errorln("Error loading .env file")
	}

	// 先读取Token配置文件
	err := system.LoadDatasourceConfig("./conf/datasource.yml")
	if err != nil {
		ErrorLogger.Errorln("读取数据库配置错误：", err)
	}
	datasourceList := system.GetDatasource()
	errorLevel := logger.Info
	//errorLevel := logger.Error
	newLogger := logger.New(
		SQLLogger,
		//log.New(os.Stdout, "\r\n", log.LstdFlags), // io writers
		logger.Config{
			SlowThreshold: time.Second, // Slow SQL threshold
			LogLevel:      errorLevel,  // Log level
			Colorful:      false,       // Disable color
		},
	)
	for _, datasource := range datasourceList {
		dsn := "%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local"
		dsnData := fmt.Sprintf(dsn, datasource.Nodes.Username, datasource.Nodes.Password, datasource.Nodes.Host,
			datasource.Nodes.Port, datasource.Nodes.Database)

		gormDB, err := gorm.Open(mysql.Open(dsnData), &gorm.Config{
			Logger: newLogger,
			NamingStrategy: schema.NamingStrategy{
				//TablePrefix: "t_",   // 表名前缀，`User` 的表名应该是 `t_users`
				SingularTable: true, // 使用单数表名，启用该选项，此时，`User` 的表名应该是 `t_user`
			},
			DisableForeignKeyConstraintWhenMigrating: true,
		})

		DBConnect[datasource.Nodes.Database] = gormDB
		if err != nil {
			fmt.Println("连接数据库失败1：", err)
			ErrorLogger.Errorln("连接数据库失败：", err)
			os.Exit(0)
		}
		sqlDB, err := gormDB.DB()
		if err != nil {
			fmt.Println("连接数据库失败2：", err)
			ErrorLogger.Errorln("连接数据库失败：", err)
			os.Exit(0)
		}

		sqlDB.SetMaxOpenConns(datasource.Nodes.MaxOpenConns)
		sqlDB.SetMaxIdleConns(datasource.Nodes.MaxIdleConns)
		sqlDB.SetConnMaxLifetime(time.Hour)
	}

}

func GetDBByName(database string) *gorm.DB {
	return DBConnect[database]
}

func GetUserDB() *gorm.DB {
	return DBConnect[constant.DBUser]
}
