package ipconf

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/shenguanchu/platos/ipconf/domain"
)

type Response struct {
	Message string      `json:"message"`
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
}

func GetIpInfoList(c context.Context, ctx *app.RequestContext) {
	defer func() {
		if err := recover(); err != nil {
			ctx.JSON(consts.StatusBadRequest, utils.H{"err": err})
		}
	}()
	// Step0: 构建客户请求信息
	ipConfCtx := domain.BuildIpConfContext(&c, ctx)
	// Step1：进行ip调度
	eps := domain.Dispatch(ipConfCtx)
	// Step2: 根据得分取 top5 返回
	ipConfCtx.AppCtx.JSON(consts.StatusOK, packRes(eps))
}
