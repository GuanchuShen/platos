package domain

import "math"

// Stat 对于gateway网关机来说，存在不同时期加入进来的物理机，所以机器的配置是不同的，使用负载来衡量会导致偏差。
// 为更好的应对动态的机器配置变化，我们统计其剩余资源值，来衡量一个机器其是否更适合增加其负载。
// 这里的数值代表的是，此 endpoint 对应的机器其，自身剩余的资源指标。
// 使用场景：
//   - ipconf 服务根据各 Gateway 的剩余资源，为客户端分配最优的连接地址
type Stat struct {
	ConnectNum   float64 // 业务上，im gateway 总体持有的长连接数量 的剩余值（剩余可用连接数（总容量 - 当前连接数））
	MessageBytes float64 // 业务上，im gateway 每秒收发消息的总字节数 的剩余值（剩余可用带宽（字节/秒）（总带宽 - 当前吞吐量））
}

// CalculateActiveScore 计算动态评分（基于实时带宽）
//
// 计算逻辑：
//   - 将剩余带宽（字节）转换为 GB 单位
//   - 剩余带宽越大，评分越高，表示该节点越空闲
//
// 为什么用带宽作为动态指标：
//   - 带宽是 IM 系统的主要瓶颈（消息收发是 IO 密集型）
//   - 带宽变化快，能实时反映节点负载状态
func (s *Stat) CalculateActiveScore() float64 {
	return getGB(s.MessageBytes)
}

// CalculateStaticScore 计算静态评分（基于连接数）
//
// 计算逻辑：
//   - 直接返回剩余可用连接数
//   - 剩余连接数越多，评分越高
//
// 为什么用连接数作为静态指标：
//   - 连接数相对稳定，不像带宽那样波动剧烈
//   - 适合作为兜底指标，当动态指标不可用时使用
func (s *Stat) CalculateStaticScore() float64 {
	return s.ConnectNum
}

// Avg 计算平均值（原地修改）
//
// 使用场景：
//   - 滑动窗口统计时，计算窗口内多个采样点的平均值
//   - 例如：窗口内有 5 个采样点，调用 Avg(5) 得到平均剩余资源
//
// 参数：
//   - num: 采样点数量（除数）
func (s *Stat) Avg(num float64) {
	s.ConnectNum /= num
	s.MessageBytes /= num
}

// Clone 深拷贝，返回一个新的 Stat 副本
//
// 为什么需要深拷贝：
//   - Stat 会被多个 goroutine 并发访问
//   - 拷贝后可以安全地在其他 goroutine 中使用，避免数据竞争
func (s *Stat) Clone() *Stat {
	newStat := &Stat{
		MessageBytes: s.MessageBytes,
		ConnectNum:   s.ConnectNum,
	}
	return newStat
}

// Add 累加另一个 Stat 的值（原地修改）
//
// 使用场景：
//   - 滑动窗口中，将新采样点的数据累加到窗口总和
func (s *Stat) Add(st *Stat) {
	if st == nil {
		return
	}
	s.ConnectNum += st.ConnectNum
	s.MessageBytes += st.MessageBytes
}

// Sub 减去另一个 Stat 的值（原地修改）
//
// 使用场景：
//   - 滑动窗口中，移除过期采样点时，从窗口总和中减去旧数据
func (s *Stat) Sub(st *Stat) {
	if st == nil {
		return
	}
	s.ConnectNum -= st.ConnectNum
	s.MessageBytes -= st.MessageBytes
}

// getGB 将字节数转换为 GB
//
// 计算逻辑：
//   - 1 GB = 2^30 字节 = 1073741824 字节
//   - 1 << 30 是位运算，等价于 2^30
//   - 结果保留两位小数
func getGB(m float64) float64 {
	return decimal(m / (1 << 30))
}

// decimal 保留两位小数（四舍五入）
//
// 计算逻辑：
//   - value * 100：将小数点右移两位
//   - + 0.5：用于四舍五入（Trunc 是截断，加 0.5 后截断等于四舍五入）
//   - math.Trunc：截断小数部分
//   - * 0.01：将小数点左移两位，恢复原来的数量级
//
// 示例：
//   - decimal(3.1415) = Trunc(314.15 + 0.5) * 0.01 = Trunc(314.65) * 0.01 = 314 * 0.01 = 3.14
//   - decimal(3.1465) = Trunc(314.65 + 0.5) * 0.01 = Trunc(315.15) * 0.01 = 315 * 0.01 = 3.15
func decimal(value float64) float64 {
	return math.Trunc(value*1e2+0.5) * 1e-2
}

func min(a, b, c float64) float64 {
	m := func(k, j float64) float64 {
		if k > j {
			return j
		}
		return k
	}
	return m(a, m(a, b))
}
