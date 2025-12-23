package ipconf

import "github.com/shenguanchu/platos/ipconf/domain"

func top5Endpoints(eps []*domain.Endpoint) []*domain.Endpoint {
	if len(eps) < 5 {
		return eps
	}
	return eps[:5]
}

func packRes(ep []*domain.Endpoint) Response {
	return Response{
		Message: "ok",
		Code:    0,
		Data:    ep,
	}
}
