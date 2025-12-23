package discovery

import "encoding/json"

type EndpointInfo struct {
	IP       string                 `json:"ip"`
	Port     string                 `json:"port"`
	MetaData map[string]interface{} `json:"meta"`
}

func UnMarshal(data []byte) (*EndpointInfo, error) {
	ep := &EndpointInfo{}
	err := json.Unmarshal(data, ep)
	if err != nil {
		return nil, err
	}
	return ep, nil
}

func (epi *EndpointInfo) Marshal() string {
	data, err := json.Marshal(epi)
	if err != nil {
		panic(err)
	}
	return string(data)
}
