package yml_config

import (
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"hello/schoolMission/core/container"
	"hello/schoolMission/global/my_errors"
	"hello/schoolMission/global/variable"
	"hello/schoolMission/utils/yml_config/ymlconfig_interf"
	"log"
	"time"
)

// 由于 viper 包本身对于文件的变化事件有一个bug，相关事件会被回调两次
// 常年未彻底解决，相关的 issue 清单：https://github.com/spf13/viper/issues?q=OnConfigChange
// 设置一个内部全局变量，记录配置文件变化时的时间点，如果两次回调事件事件差小于1秒，我们认为是第二次回调事件，而不是人工修改配置文件
// 这样就避免了 viper 包的这个bug

var lastChangeTime time.Time

type ymlConfig struct {
	viper *viper.Viper
}



func init() {
	lastChangeTime = time.Now()
}

// CreateYamlFactory 创建一个yaml配置文件工厂
// 参数设置为可变参数的文件名, 我们只取第一个参数作为文件名
func CreateYamlFactory(fileName ...string) ymlconfig_interf.YmlConfigInterf  {
	yamlConfig := viper.New()
	// 配置文件所在目录
	yamlConfig.AddConfigPath(variable.BasePath + "/config")
	// 需要读取的文件名，默认为: config
	if len(fileName) == 0 {
		yamlConfig.SetConfigName("config")
	} else {
		yamlConfig.SetConfigName(fileName[0])
	}
	// 设置配置文件类型(后缀)为yml
	yamlConfig.SetConfigType("yml")

	if err := yamlConfig.ReadInConfig(); err != nil {
		log.Fatal(my_errors.ErrorsConfigInitFail + ": " + err.Error())
	}

	return  &ymlConfig{
		yamlConfig,
	}
}


func (y ymlConfig) ConfigFileChangeListen() {
	y.viper.OnConfigChange(func(changeEvent fsnotify.Event) {
		if time.Now().Sub(lastChangeTime).Seconds() >= 1 {
			if changeEvent.Op.String() == "WRITE" {
				y.clearCache()
				lastChangeTime = time.Now()
			}
		}
	})
	y.viper.WatchConfig()
}

// Get 一个原始值
func (y ymlConfig) Get(keyName string) interface{} {
	if y.keyIsCache(keyName){
		return y.getValueFromCache(keyName)
	} else {
		value := y.viper.Get(keyName)
		y.cache(keyName, value)
		return value
	}
}

func (y ymlConfig) GetString(keyName string) string {
	if y.keyIsCache(keyName){
		return y.getValueFromCache(keyName).(string)
	} else {
		value := y.viper.GetString(keyName)
		y.cache(keyName, value)
		return value
	}
}

func (y ymlConfig) GetBool(keyName string) bool {
	if y.keyIsCache(keyName){
		return y.getValueFromCache(keyName).(bool)
	} else {
		value := y.viper.GetBool(keyName)
		y.cache(keyName, value)
		return value
	}
}

func (y ymlConfig) GetInt(keyName string) int {
	if y.keyIsCache(keyName){
		return y.getValueFromCache(keyName).(int)
	} else {
		value := y.viper.GetInt(keyName)
		y.cache(keyName, value)
		return value
	}
}

func (y ymlConfig) GetInt32(keyName string) int32 {
	if y.keyIsCache(keyName){
		return y.getValueFromCache(keyName).(int32)
	} else {
		value := y.viper.GetInt32(keyName)
		y.cache(keyName, value)
		return value
	}
}

func (y ymlConfig) GetInt64(keyName string) int64 {
	if y.keyIsCache(keyName){
		return y.getValueFromCache(keyName).(int64)
	} else {
		value := y.viper.GetInt64(keyName)
		y.cache(keyName, value)
		return value
	}
}

func (y ymlConfig) GetFloat64(keyName string) float64 {
	if y.keyIsCache(keyName){
		return y.getValueFromCache(keyName).(float64)
	} else {
		value := y.viper.GetFloat64(keyName)
		y.cache(keyName, value)
		return value
	}
}

func (y ymlConfig) GetDuration(keyName string) time.Duration {
	if y.keyIsCache(keyName){
		return y.getValueFromCache(keyName).(time.Duration)
	} else {
		value := y.viper.GetDuration(keyName)
		y.cache(keyName, value)
		return value
	}
}

func (y ymlConfig) GetStringSlice(keyName string) []string {
	if y.keyIsCache(keyName){
		return y.getValueFromCache(keyName).([]string)
	} else {
		value := y.viper.GetStringSlice(keyName)
		y.cache(keyName, value)
		return value
	}
}

// 清空已经更改的配置项信息
func (y ymlConfig) clearCache() {
	container.CreateContainersFactory().FuzzyDelete(variable.ConfigKeyPrefix)
}

// 判断相关的键是否已经缓存
func (y *ymlConfig) keyIsCache(key string)bool  {
	 _, exists := container.CreateContainersFactory().KeyIsExists(key)
	return exists
}

// 对键值进行缓存
func (y *ymlConfig) cache(keyName string, value interface{})  {
	container.CreateContainersFactory().Set(keyName,value)
}

// 通过键获取缓存的值
func (y *ymlConfig) getValueFromCache(keyName string)interface{}{
	return container.CreateContainersFactory().Get(keyName)
}

// 允许clone 一个相同功能的结构体
func (y *ymlConfig) Clone(fileName string)ymlconfig_interf.YmlConfigInterf  {
	// 这里存在一个深拷贝，需要注意，避免拷贝的结构体操作对原始结构体造成影响
	var ymlC = *y
	var ymlConfViper = *(y.viper)
	(&ymlC).viper = &ymlConfViper

	(&ymlC).viper.SetConfigName(fileName)
	if err := (&ymlC).viper.ReadInConfig(); err != nil {
		variable.ZapLog.Error(my_errors.ErrorsConfigInitFail, zap.Error(err))
	}
	return &ymlC
}





