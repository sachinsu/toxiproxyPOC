package setup

import (
	toxiproxy "github.com/Shopify/toxiproxy/client"
)

var toxiClient *toxiproxy.Client
var proxies map[string]*toxiproxy.Proxy

//ref: https://github.com/Shopify/toxiproxy/tree/master/client
func init() {
	var err error
	toxiClient = toxiproxy.NewClient("localhost:8474")
	proxies, err = toxiClient.Populate([]toxiproxy.Proxy{{
		Name:     "pgsql",
		Listen:   "localhost:32323",
		Upstream: "localhost:5432",
	}})
	if err != nil {
		panic(err)
	}
	// Alternatively, create the proxies manually with
	// toxiClient.CreateProxy("redis", "localhost:26379", "localhost:6379")
}

func InjectSlowness(name string, delay int32) *toxiproxy.Proxy {
	proxies[name].AddToxic("", "latency", "", 1, toxiproxy.Attributes{
		"latency": delay,
	})
	// defer proxies[name].RemoveToxic("latency_downstream")
	return proxies[name]
}
