package conf

import (
	"fmt"
	"github.com/spf13/viper"
)

func ViperConf() {
	viper.SetConfigName("valueConfig") // 设置配置文件名（不包含扩展名）
	viper.SetConfigType("yaml")        // 设置配置文件类型，可选
	viper.AddConfigPath("conf")        // 添加配置文件搜索路径
	err := viper.ReadInConfig()        // 加载配置文件
	if err != nil {
		panic(fmt.Errorf("无法加载配置文件: %s", err))
	} // 添加配置文件搜索路径，可选
}
