package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/ip2location/ip2location-go/v9"
)

type Proxy struct {
	IP   string
	Port string
}

type Geo struct {
	Country string
	Proxy   Proxy
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter the proxy file name: ")
	proxyFileName, _ := reader.ReadString('\n')
	proxyFileName = strings.TrimSpace(proxyFileName)

	db, err := ip2location.OpenDB("./IP2LOCATION-LITE-DB1.BIN")

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	proxies, err := getProxyList(proxyFileName)

	if err != nil {
		log.Fatal(err)
	}

	err = os.RemoveAll("countries")
	if err != nil {
		log.Fatal(err)
	}

	err = os.Mkdir("countries", 0755)
	if err != nil {
		log.Fatal(err)
	}

	for _, proxy := range *proxies {
		results, _ := db.Get_all(proxy.IP)

		filename := fmt.Sprintf("countries/%s.txt", results.Country_long)
		file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		_, err = file.WriteString(fmt.Sprintf("%s:%s\n", proxy.IP, proxy.Port))
		if err != nil {
			log.Fatal(err)
		}
	}
}

func getProxyList(dbFilePath string) (*[]Proxy, error) {
	data, err := os.ReadFile(dbFilePath)

	if err != nil {
		return nil, err
	}

	proxies := make([]Proxy, 0)

	for _, line := range strings.Split(strings.TrimSpace(string(data)), "\n") {
		words := strings.Fields(line)
		if len(words) == 0 {
			continue
		}

		if words[0] == "New" {
			proxyport := strings.Split(words[2], ":")
			proxies = append(proxies, Proxy{
				IP:   proxyport[0],
				Port: proxyport[1],
			})
		}

	}

	return &proxies, nil
}
