package setup

import (
	toxiproxy "github.com/Shopify/toxiproxy/client"
)

var toxiClient *toxiproxy.Client

//ref: https://github.com/Shopify/toxiproxy/tree/master/client
func init() {
	toxiClient = toxiproxy.NewClient("localhost:8474")
}

// AddPGProxy for postgresql
func AddPGProxy(listen string, upstream string) (*toxiproxy.Proxy, error) {
	// Alternatively, create the proxies manually with
	// return toxiClient.CreateProxy("pgsql", "[::]:32777", "localhost:5432")
	return toxiClient.CreateProxy("pgsql", listen, upstream)
}

// InjectLatency helper
func InjectLatency(name string, delay int) *toxiproxy.Proxy {
	proxy, err := toxiClient.Proxy(name)

	if err != nil {
		panic(err)
	}

	proxy.AddToxic("", "latency", "", 1, toxiproxy.Attributes{
		"latency": delay,
	})
	// defer proxies[name].RemoveToxic("latency_downstream")
	return proxy
}
