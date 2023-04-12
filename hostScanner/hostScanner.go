package hostScanner

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"time"
)

func scanner(target string, ports, results chan int) {
	for port := range ports {
		address := fmt.Sprintf("%s:%d", target, port)
		fmt.Printf("%s:%d\n", target, port)
		conn, err := net.DialTimeout("tcp", address, 2*time.Second)
		if err != nil {
			results <- 0
		} else {
			conn.Close()
			results <- port
		}
	}
}

func Start(host string) []string {
	log.Printf("Сканирование портов: %s \n", host)

	ports := make(chan int, 10000)
	results := make(chan int)

	var openPorts []string

	for port := 1; port <= cap(ports); port++ {
		go scanner(host, ports, results)
	}

	go func() {
		for port := 1; port <= 65536; port++ {
			ports <- port
		}
	}()

	for port := 1; port <= 65536; port++ {
		portStatus := <-results
		if portStatus != 0 {
			openPorts = append(openPorts, strconv.Itoa(portStatus))
		}
	}

	close(ports)
	close(results)
	log.Println("Сканирование завершено")
	log.Printf("Открытые порты: %v", openPorts)
	return openPorts
}
