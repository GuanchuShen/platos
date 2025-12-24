package cmd

import (
	"github.com/shenguanchu/platos/client"
	"github.com/spf13/cobra"
)

// init 注册 client 子命令到 rootCmd
// 由 Go 运行时在 main() 执行前自动调用
func init() {
	rootCmd.AddCommand(clientCmd)
}

// clientCmd 客户端子命令
// 用法: platos client --config=./platos.yaml
var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "启动客户端",
	Run:   ClientHandle,
}

func ClientHandle(cmd *cobra.Command, args []string) {
	client.RunMain()
}
