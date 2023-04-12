package portScanner

import (
	"log"
	"time"

	"github.com/Ullaakut/nmap/v2"
)

type Port struct {
	ID       uint16
	Protocol string
	Service
}

type Service struct {
	Name    string
	Product string
	Version string
}

func scanner(host string, ports chan string, res chan Port) {

	for port := range ports {
		scan, err := nmap.NewScanner(
			nmap.WithTargets(host),
			nmap.WithPorts(port),
			nmap.WithServiceInfo(),
		)

		if err != nil {
			log.Fatalf("unable to create nmap scanner: %v", err)
		}

		result, _, err := scan.Run()
		if err != nil {
			log.Fatalf("unable to run nmap scan: %v", err)
		}

		//if warnings != nil {
		//	log.Printf("Warnings: %v", warnings)
		//}

		for _, h := range result.Hosts {
			if len(h.Ports) == 0 || len(h.Addresses) == 0 {
				continue
			}

			for _, p := range h.Ports {
				s := Service{
					Name:    p.Service.Name,
					Product: p.Service.Product,
					Version: p.Service.Version,
				}
				r := Port{
					ID:       p.ID,
					Protocol: p.Protocol,
					Service:  s,
				}

				log.Printf("Port: %s %s %s", port, p.Service.Product, p.Service.Version)

				res <- r
			}
		}
	}
}

func Start(host string, p []string) {
	start := time.Now()

	ports := make(chan string, 10)
	result := make(chan Port)

	var portsInfo []Port

	for i := 1; i <= cap(ports); i++ {
		go scanner(host, ports, result)
	}

	go func() {
		for _, elem := range p {
			ports <- elem
		}
	}()

	for i := 1; i <= cap(ports); i++ {
		r := <-result
		emptyPort := Port{}
		if r != emptyPort {
			portsInfo = append(portsInfo, r)
		}
	}
	close(ports)
	close(result)

	duration := time.Since(start)
	log.Printf("Время сканирования портов: %v", duration)
}
