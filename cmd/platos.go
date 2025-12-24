// Package cmd 提供命令行入口，基于 cobra 框架实现子命令管理。
// 初始化顺序：
//  1. 各子命令文件的 init() 自动注册子命令到 rootCmd（如 ipconf.go、client.go）
//  2. cobra.OnInitialize 注册 initConfig，在命令执行前加载配置文件
//  3. main.go 调用 Execute() 启动命令解析
package cmd

import (
	"fmt"
	"os"

	"github.com/shenguanchu/platos/common/config"
	"github.com/spf13/cobra"
)

// ConfigPath 配置文件路径，通过 --config 参数指定
var ConfigPath string

// init 注册全局 flags 和配置初始化回调
// 注意：此函数由 Go 运行时自动调用，无需手动调用
func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&ConfigPath, "config", "./plato.yaml", "config file (default is ./plato.yaml)")
}

// rootCmd 根命令，所有子命令都注册到此命令下
var rootCmd = &cobra.Command{
	Use:   "plato",
	Short: "这是一个超牛逼的IM系统",
	Run:   Plato,
}

// Execute 启动命令行解析，由 main.go 调用
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func Plato(cmd *cobra.Command, args []string) {
}

// initConfig 加载配置文件，由 cobra.OnInitialize 在命令执行前自动调用
func initConfig() {
	if ConfigPath != "" {
		config.Init(ConfigPath)
	}
}
