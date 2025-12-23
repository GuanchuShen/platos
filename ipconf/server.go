package ipconf

import (
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/shenguanchu/platos/ipconf/domain"
	"github.com/shenguanchu/platos/ipconf/source"
)

func RunMain() {
	//config.Init()
	source.Init()
	domain.Init()
	s := server.Default(server.WithHostPorts(":6789"))
	s.GET("/ip/list", GetIpInfoList)
	s.Spin()
}
