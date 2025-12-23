package domain

const (
	windowSize = 5
)

type statWindow struct {
	statQueue []*Stat
	statChan  chan *Stat
	sumStat   *Stat
	idx       int64
}

func newStatWindow() *statWindow {
	return &statWindow{
		statQueue: make([]*Stat, windowSize),
		statChan:  make(chan *Stat),
		sumStat:   &Stat{},
	}
}

func (sw *statWindow) getStat() *Stat {
	res := sw.sumStat.Clone()
	res.Avg(windowSize)
	return res
}

func (sw *statWindow) appendStat(s *Stat) {
	// 减去即将被删除的stat
	sw.sumStat.Sub(sw.statQueue[sw.idx%windowSize])
	// 更新最新的stat
	sw.statQueue[sw.idx%windowSize] = s
	// 计算最新的窗口和
	sw.sumStat.Add(s)
	sw.idx++
}
