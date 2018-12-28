package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

const (
	apiAddress = `https://api.64clouds.com/v1/getServiceInfo`
	byte2gb    = 1024 * 1024 * 1024
)

var flags struct {
	show   bool
	save   bool //save to config
	list   bool //list
	help   bool
	veid   string
	apikey string
	delete string
}

func apiCall(veid, apiKey string) {
	resp, err := http.Get(apiAddress + "?veid=" + veid + "&api_key=" + apiKey)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	j := NewJson(body)
	if j.Get("error").Number() != 0 {
		fmt.Println("error: veid or apikey dot match.")
		return
	}
	dataCounter := j.Get("data_counter").Number() * j.Get("monthly_data_multiplier").Number() / byte2gb
	dataAll := j.Get("plan_monthly_data").Number() * j.Get("monthly_data_multiplier").Number() / byte2gb
	fmt.Println("vm_type: ", j.Get("vm_type").String())
	fmt.Println("node_location: ", j.Get("node_location").String())
	fmt.Println("ip_address: ", j.Get("ip_addresses").Item(0).String())
	fmt.Printf("data_usage:  %.2f/%.2f (GB)\n", dataCounter, dataAll)
}

func lines2map(data []byte) (m map[string][]byte) {
	lines := bytes.Split(data, []byte("\n"))
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}
		items := bytes.Split(line, []byte(","))
		if items[0] != nil && items[1] != nil {
			if m == nil {
				m = make(map[string][]byte)
			}
			m[string(items[0])] = items[1]
		}
	}
	return
}
func map2Lines(m map[string][]byte) (lines []byte) {
	for key, val := range m {
		temp := []byte(key + "," + string(val) + "\n")
		lines = append(lines, temp...)
	}
	return
}

func main() {
	configFile, err := os.OpenFile("./config.txt", os.O_RDWR, 0666)
	if err != nil {

		cfile, err := os.Create("./config.txt")
		if err != nil {
			panic(err)
		}
		configFile = cfile
	}
	defer configFile.Close()
	data, err := ioutil.ReadAll(configFile)
	m := lines2map(data)

	if err != nil {
		panic(err)
	}

	flag.BoolVar(&flags.show, "show", false, "show server's detail.")
	flag.BoolVar(&flags.save, "save", false, "save veid and apikey in config.txt")
	flag.BoolVar(&flags.list, "l", false, "list all the saved servers")
	flag.StringVar(&flags.delete, "d", "", "remove veid from list config")
	flag.BoolVar(&flags.help, "h", false, "help")
	flag.StringVar(&flags.veid, "veid", "", "server's veid")
	flag.StringVar(&flags.apikey, "apikey", "", "server's api key")
	flag.Parse()

	if flags.help || flag.NFlag() == 0 {
		flag.PrintDefaults()
		return
	}

	if flags.veid != "" && flags.apikey != "" {
		if flags.show {
			//show
			apiCall(flags.veid, flags.apikey)

		} else if flags.save {
			//save

			if len(data) == 0 {
				//new file
				_, err := fmt.Fprintf(configFile, "%s,%s\n", flags.veid, flags.apikey)
				if err != nil {
					panic(err)
				}
			} else {

				m[flags.veid] = []byte(flags.apikey)
				waitToWrite := map2Lines(m)
				configFile.Truncate(0)
				configFile.Write(waitToWrite)
				apiCall(flags.veid, flags.apikey)
			}

		}
		return
	}

	if flags.list {
		if m != nil {
			for key, val := range m {
				fmt.Printf("veid:%s apiKey:%s\n", string(key), val)
			}
		}
		return
	}

	if flags.delete != "" {
		if m != nil {
			delete(m, flags.delete)
			waitToWrite := map2Lines(m)
			configFile.Truncate(0)
			configFile.Write(waitToWrite)
			for key, val := range m {
				fmt.Printf("veid:%s apiKey:%s\n", string(key), val)
			}
		}
		return
	}

	if flags.show {
		if len(m) == 1 {
			for key, val := range m {
				apiCall(key, string(val))
				return
			}
		}
		if flags.veid != "" {
			if val, ok := m[flags.veid]; ok == true {
				apiCall(flags.veid, string(val))
			} else {
				fmt.Println("veid not found in config file")
			}
		} else {
			fmt.Println("miss veid(-show -veid xxx)")
		}
		return
	}

}
