package cmd

import (
	"github.com/shenguanchu/platos/ipconf"
	"github.com/spf13/cobra"
)

// init 注册 ipconf 子命令到 rootCmd
// 由 Go 运行时在 main() 执行前自动调用
func init() {
	rootCmd.AddCommand(ipConfCmd)
}

// ipConfCmd IP 调度服务子命令
// 用法: platos ipconf --config=./platos.yaml
var ipConfCmd = &cobra.Command{
	Use:   "ipconf",
	Short: "启动 IP 调度服务",
	Run:   IpConfHandle,
}

func IpConfHandle(cmd *cobra.Command, args []string) {
	ipconf.RunMain()
}
