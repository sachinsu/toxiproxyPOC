# Resiliency Testing using ToxiProxy

This is [Go](https://golang.org) lang project that demonstrates how to use [Toxiproxy](https://toxiproxy.io) from Shopify for chaos and resiliency testing. 

## Setup
  * Clone the repository
  * Create database in PostgreSQL like ` psql -h localhost -U postgres -w -c "create database proxytest;"`
  * Edit `cmd/proxytest/main.go` for database connection and setting up proxy,

        ```
        func main() {
            dbConnString := "postgres://postgres:@localhost:32777/proxytest?sslmode=disable"

            tproxy, err := setup.AddPGProxy("[::]:32777", "localhost:5432")

        ```
        (Todo: for me to move these settings to configuration). Update `dbConnString` as well as ports in call to `AddPGProxy` as per your environment

  * Project consists of REST API that reads from PostgreSQL database. It has two end-points, 
    * /api/cities - This end-point returns list of cities in JSON format but without any Toxics
    * /api/cities/latency - This end-point returns list of cities in JSON format with induced latency. One can specify latency in miliseconds using `/api/cities/latency?delay=200` (default is 100).
    
 * Download Toxiproxy server for your OS from [here](https://toxiproxy.io).
 * Run the server
 * Run Benchmark tests `go test -v -bench=. -benchtime=10s  ./...` from root folder of repo.

 * or Use benchmarking tools like Rakyll's excellent [hey](https://github.com/rakyll/hey) to generate traffic and observe impact. 
    * Run Web server using `go run cmd/proxytest/main.go`
    * Run hey, `hey http://localhost:8080/api/cities`
    * On my Laptop, 4-core I5, 8GB RAM with Both API + Postgresql + Toxiproxy running, it shows below, 
      * Without Latency,

        ```
        Summary:
        Total:        23.5590 secs
        Slowest:      7.2657 secs
        Fastest:      0.9188 secs
        Average:      5.4242 secs
        Requests/sec: 8.4893


        Response time histogram:
        0.919 [1]     |■
        1.554 [7]     |■■■■
        2.188 [2]     |■
        2.823 [9]     |■■■■■
        3.458 [3]     |■■
        4.092 [2]     |■
        4.727 [13]    |■■■■■■■■
        5.362 [22]    |■■■■■■■■■■■■■
        5.996 [69]    |■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
        6.631 [49]    |■■■■■■■■■■■■■■■■■■■■■■■■■■■■
        7.266 [23]    |■■■■■■■■■■■■■


            Latency distribution:
            10% in 3.2869 secs
            25% in 5.2367 secs
            50% in 5.7531 secs
            75% in 6.2558 secs
            90% in 6.6610 secs
            95% in 6.9587 secs
            99% in 7.1657 secs

            Details (average, fastest, slowest):
            DNS+dialup:   0.0096 secs, 0.9188 secs, 7.2657 secs
            DNS-lookup:   0.0077 secs, 0.0000 secs, 0.0450 secs
            req write:    0.0003 secs, 0.0000 secs, 0.0023 secs
            resp wait:    5.3926 secs, 0.8664 secs, 7.2541 secs
            resp read:    0.0216 secs, 0.0050 secs, 0.0900 secs

            Status code distribution:
            [200] 200 responses

        ```

      * With Latency of 300 Ms,

        ```
            Summary:
            Total:        23.4831 secs
            Slowest:      8.2054 secs
            Fastest:      2.1162 secs
            Average:      5.3499 secs
            Requests/sec: 8.5168


            Response time histogram:
            2.116 [1]     |■
            2.725 [20]    |■■■■■■■■■■■■■■■■■■■
            3.334 [1]     |■
            3.943 [6]     |■■■■■■
            4.552 [23]    |■■■■■■■■■■■■■■■■■■■■■
            5.161 [28]    |■■■■■■■■■■■■■■■■■■■■■■■■■■
            5.770 [35]    |■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
            6.379 [43]    |■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
            6.988 [23]    |■■■■■■■■■■■■■■■■■■■■■
            7.597 [9]     |■■■■■■■■
            8.205 [11]    |■■■■■■■■■■


            Latency distribution:
            10% in 2.6683 secs
            25% in 4.4909 secs
            50% in 5.5572 secs
            75% in 6.3346 secs
            90% in 7.1538 secs
            95% in 7.6392 secs
            99% in 8.1987 secs

            Details (average, fastest, slowest):
            DNS+dialup:   0.0096 secs, 2.1162 secs, 8.2054 secs
            DNS-lookup:   0.0063 secs, 0.0000 secs, 0.0326 secs
            req write:    0.0011 secs, 0.0000 secs, 0.0131 secs
            resp wait:    5.2940 secs, 1.7847 secs, 8.1922 secs
            resp read:    0.0451 secs, 0.0050 secs, 0.4180 secs

            Status code distribution:
            [200] 200 responses

        ```


## Todo 
    * Implement other toxics(like down, bandwidth, timeouts) provided by Toxiproxy.