package client

import "github.com/willoma/keepakonf/internal/log"

func (c *client) logs(a ...any) {
	if len(a) < 2 {
		callback(a, log.GetPage(0))
		return
	}

	options, ok := a[0].(map[string]any)
	if !ok {
		callback(a, log.GetPage(0))
		return
	}

	offsetIface, ok := options["offset"]
	if !ok {
		callback(a, log.GetPage(0))
		return
	}

	offset, ok := offsetIface.(float64)
	if !ok {
		offset = 0
	}

	callback(a, log.GetPage(int(offset)))
}
