package server

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	toxiproxy "github.com/Shopify/toxiproxy/client"
	"github.com/sachinsu/chaosengg/internal/setup"
)

var result int

func Setup(proxyport int) (*Server, *toxiproxy.Proxy, error) {
	if proxyport > 0 {
		tproxy, err := setup.AddPGProxy(fmt.Sprintf("[::]:%d", proxyport), "localhost:5432")
		if err != nil {
			return nil, nil, err
		}
		srv, err := NewServer(fmt.Sprintf("postgres://postgres:@localhost:%d/proxytest?sslmode=disable", proxyport), ":8080", "")
		return srv, tproxy, err
	} else {
		srv, err := NewServer("postgres://postgres:@localhost:5432/proxytest?sslmode=disable", ":8080", "")
		return srv, nil, err
	}
}
func TestCitiesList(t *testing.T) {
	srv, _, err := Setup(0)
	if err != nil {
		t.Fatal(err)
	}
	srv.SetupRoutes()

	req, _ := http.NewRequest("GET", "/cities", nil)
	rr := httptest.NewRecorder()

	srv.Router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected %d, got %d", http.StatusOK, rr.Code)
	}
	if rr.Body.Len() == 0 {
		t.Errorf("expected > 0, got 0")
	}
}

func BenchmarkWithNoProxy(b *testing.B) {

	srv, _, err := Setup(0)
	if err != nil {
		b.Fatal(err)
	}
	srv.SetupRoutes()
	rr := httptest.NewRecorder()

	for n := 0; n < b.N; n++ {
		count := n
		if n*100 < 10000 {
			count = n * 100
		}
		req, _ := http.NewRequest("GET", fmt.Sprintf("/cities?limit=%d", count), nil)
		srv.Router.ServeHTTP(rr, req)
		// always store the result to a package level variable
		// so the compiler cannot eliminate the Benchmark itself.
		result = rr.Code
	}
}

func BenchmarkWithProxy(b *testing.B) {

	srv, tproxy, err := Setup(32555)
	if err != nil {
		b.Fatal(err)
	}
	defer tproxy.Delete()

	tproxy = setup.InjectLatency("pgsql", 100)

	srv.SetupRoutes()
	rr := httptest.NewRecorder()

	for n := 0; n < b.N; n++ {
		count := n
		if n*100 < 10000 {
			count = n * 100
		}
		req, _ := http.NewRequest("GET", fmt.Sprintf("/cities?limit=%d", count), nil)
		srv.Router.ServeHTTP(rr, req)
		// always store the result to a package level variable
		// so the compiler cannot eliminate the Benchmark itself.
		result = rr.Code
	}
}
