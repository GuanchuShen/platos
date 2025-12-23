package source

import (
	"fmt"

	"github.com/shenguanchu/platos/common/discovery"
)

var eventChan chan *Event

func EventChan() <-chan *Event {
	return eventChan
}

type EventType string

const (
	AddNodeEvent EventType = "addNode"
	DelNodeEvent EventType = "delNode"
)

type Event struct {
	Type         EventType
	IP           string
	Port         string
	ConnectNum   float64
	MessageBytes float64
}

func NewEvent(ep *discovery.EndpointInfo) *Event {
	if ep == nil || ep.MetaData == nil {
		return nil
	}
	var connNum, msgBytes float64
	if data, ok := ep.MetaData["connect_num"]; ok {
		connNum = data.(float64) // 如果出错，此处应该panic 暴露错误
	}
	if data, ok := ep.MetaData["message_bytes"]; ok {
		msgBytes = data.(float64) // 如果出错，此处应该panic 暴露错误
	}
	return &Event{
		Type:         AddNodeEvent,
		IP:           ep.IP,
		Port:         ep.Port,
		ConnectNum:   connNum,
		MessageBytes: msgBytes,
	}
}

func (e *Event) Key() string {
	return fmt.Sprintf("%s:%s", e.IP, e.Port)
}
