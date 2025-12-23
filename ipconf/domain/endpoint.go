package domain

import (
	"sync/atomic"
)

type Endpoint struct {
	IP          string               `json:"ip"`
	Port        string               `json:"port"`
	ActiveScore float64              `json:"-"`
	StaticScore float64              `json:"-"`
	Stats       atomic.Pointer[Stat] `json:"-"`
	window      *statWindow
}

// NewEndpoint 创建一个新的服务端点
//
// 工作流程：
//  1. 创建 Endport 实例，设置 IP 和 Port
//  2. 创建滑动窗口（用于平滑统计数据）
//  3. 启动一个后台 goroutine 持续处理统计数据更新
//
// 架构图：
//
//	                  statChan
//	UpdateStat() ──────────────► goroutine ──► 更新滑动窗口
//	     │                           │
//	     │                           ▼
//	     │                      计算窗口平均值
//	     │                           │
//	     │                           ▼
//	     │                      原子更新 Stats 指针
//	     │                           │
//	     ▼                           ▼
//	外部调用              CalculateScore() 读取 Stats 计算评分
//
// 为什么用 goroutine + channel：
//   - 统计数据更新可能很频繁（每次 etcd 事件都会触发）
//   - 使用 channel 串行化更新操作，避免并发写入问题
//   - goroutine 异步处理，不阻塞调用方
func NewEndpoint(ip, port string) *Endpoint {
	ep := &Endpoint{
		IP:   ip,
		Port: port,
	}
	ep.window = newStatWindow()
	ep.Stats.Store(ep.window.getStat())
	// 启动后台 goroutine 处理统计数据更新
	go func() {
		for stat := range ep.window.statChan {
			ep.window.appendStat(stat)
			newStat := ep.window.getStat()
			//atomic.SwapPointer((*unsafe.Pointer)(unsafe.Pointer(&ep.Stats)), unsafe.Pointer(newStat))
			ep.Stats.Store(newStat)
		}
	}()
	return ep
}

func (ep *Endpoint) UpdateStat(s *Stat) {
	ep.window.statChan <- s
}

func (ep *Endpoint) CalculateScore(ctx *IpConfContext) {
	// 如果 stats 字段是空的，则直接使用上一次计算的结果，此次不更新
	stat := ep.Stats.Load()
	if stat != nil {
		ep.ActiveScore = stat.CalculateActiveScore()
		ep.StaticScore = stat.CalculateStaticScore()
	}
}
