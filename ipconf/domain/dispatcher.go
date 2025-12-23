package domain

import (
	"sort"
	"sync"

	"github.com/shenguanchu/platos/ipconf/source"
)

type Dispatcher struct {
	candidateTable map[string]*Endpoint
	sync.RWMutex
}

var dp *Dispatcher

func Init() {
	dp = &Dispatcher{}
	dp.candidateTable = make(map[string]*Endpoint)
	go func() {
		// 从 channel 持续接收数据的惯用写法（for range channel）
		// 持续监听事件
		for event := range source.EventChan() {
			switch event.Type {
			case source.AddNodeEvent:
				dp.addNode(event)
			case source.DelNodeEvent:
				dp.delNode(event)
			}
		}
	}()
}

func Dispatch(ctx *IpConfContext) []*Endpoint {
	// 1. get candidate endpoint
	eps := dp.getCandidateEndpoint(ctx)
	// 2. calculate score
	for _, ep := range eps {
		ep.CalculateScore(ctx)
	}
	// 3. global sort, 动静结合的排序策略
	sort.Slice(eps, func(i, j int) bool {
		// 优先基于活跃分数进行排序
		if eps[i].ActiveScore > eps[j].ActiveScore {
			return true
		}
		// 如果活跃分数相同，则使用静态分数排序
		if eps[i].ActiveScore == eps[j].ActiveScore {
			if eps[i].StaticScore > eps[j].StaticScore {
				return true
			}
			return false
		}
		return false
	})
	return eps
}

func (dp *Dispatcher) getCandidateEndpoint(ctx *IpConfContext) []*Endpoint {
	dp.RLock()
	defer dp.RUnlock()
	candidateList := make([]*Endpoint, 0, len(dp.candidateTable))
	for _, ep := range dp.candidateTable {
		candidateList = append(candidateList, ep)
	}
	return candidateList
}

func (dp *Dispatcher) addNode(event *source.Event) {
	dp.Lock()
	defer dp.Unlock()
	var (
		ep *Endpoint
		ok bool
	)
	if ep, ok = dp.candidateTable[event.Key()]; !ok {
		ep = NewEndpoint(event.IP, event.Port)
		dp.candidateTable[event.Key()] = ep
	}
	ep.UpdateStat(&Stat{
		ConnectNum:   event.ConnectNum,
		MessageBytes: event.MessageBytes,
	})
}

func (dp *Dispatcher) delNode(event *source.Event) {
	dp.Lock()
	defer dp.Unlock()
	delete(dp.candidateTable, event.Key())
}
