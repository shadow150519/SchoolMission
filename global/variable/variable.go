package variable

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
	"hello/schoolMission/global/my_errors"
	"hello/schoolMission/utils/snow_flake/snowflake_interf"
	"hello/schoolMission/utils/yml_config/ymlconfig_interf"
	"log"
	"os"
	"strings"
)

var (
	BasePath string 						    // 定义项目根路径
	EventDestroyPrefix  = "Destroy_"            // 程序退出时需要销毁的事件前缀
	ConfigKeyPrefix 	= "Config_" 			// 配置文件键值缓存时，键的前缀
	DateFormat          = "2006-01-02 15:04:05"	// 时间格式

	// 全局日志指针
	ZapLog *zap.Logger
	// 全局配置文件
	ConfigYml ymlconfig_interf.YmlConfigInterf 		 // 全局配置文件指针
	ConfigGormv2Yml ymlconfig_interf.YmlConfigInterf // 数据库配置文件指针

	GormDbMysql      *gorm.DB 			// 全局gorm的客户端连接
	GormDbSqlserver  *gorm.DB 		    // 全局gorm的客户端连接
	GormDbPostgreSql *gorm.DB           // 全局gorm的客户端连接

	// 雪花算法全局变量
	SnowFlake snowflake_interf.InterfaceSnowFlake
)

func init() {
	// 初始化程序根目录
	if curPath, err := os.Getwd(); err == nil {
		// 路径进行处理，兼容单元测试程序启动时的奇怪路径
		if len(os.Args) > 1 && strings.HasPrefix(os.Args[1],"-test"){
			BasePath = strings.Replace(strings.Replace(curPath,`\test`,"", 1),`/test`,"",1)
		} else {
			BasePath = curPath
		}
	} else {
		log.Fatal(my_errors.ErrorsBasePath)
	}
}
