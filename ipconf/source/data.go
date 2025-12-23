package source

import (
	"context"

	"github.com/bytedance/gopkg/util/logger"
	"github.com/shenguanchu/platos/common/discovery"
)

func Init() {
	eventChan = make(chan *Event)
	ctx := context.Background()
	go DataHandler(&ctx)
	//if config.IsDebug() {
	//	ctx := context.Background()
	//	test
	//}
}

// DataHandler 服务发现处理
func DataHandler(ctx *context.Context) {
	dis := discovery.NewServiceDiscovery(ctx, []string{"localhost:2379"})
	defer dis.Close()
	setFunc := func(key, value string) {
		if ep, err := discovery.UnMarshal([]byte(value)); err == nil {
			if event := NewEvent(ep); ep != nil {
				event.Type = AddNodeEvent
				eventChan <- event
			}
		} else {
			logger.CtxErrorf(*ctx, "DataHandler.setFunc.err: %s", err.Error())
		}
	}
	delFunc := func(key, value string) {
		if ep, err := discovery.UnMarshal([]byte(value)); err == nil {
			if event := NewEvent(ep); ep != nil {
				event.Type = DelNodeEvent
				eventChan <- event
			}
		} else {
			logger.CtxErrorf(*ctx, "DataHandler.setFunc.err: %s", err.Error())
		}
	}
	err := dis.WatchService("/platos/ip_dispatcher", setFunc, delFunc)
	if err != nil {
		panic(err)
	}
}
