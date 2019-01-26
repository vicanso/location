package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

const (
	// 分隔之后的长度
	divideSize = 7
	template   = `
package service 

var (
	// WordMapList word map list
	WordMapList = []string{%s}
	// IPInfos ip infos
	IPInfos = [][]uint32{%s}
	// IPCount IP count
	IPCount = %d
)
`
)

var (
	wordMapList = []string{}
)

func convertIP(ip string) (value uint32, err error) {
	arr := strings.SplitN(ip, ".", -1)
	offset := 8
	max := 3
	for index, item := range arr {
		v, err := strconv.Atoi(item)
		if err != nil {
			return 0, err
		}
		value += uint32(v) << uint32(offset*(max-index))
	}
	return
}

func addWord(str string) uint32 {
	index := -1
	for i, item := range wordMapList {
		if item == str {
			index = i
		}
	}
	if index == -1 {
		wordMapList = append(wordMapList, str)
		index = len(wordMapList) - 1
	}
	return uint32(index)
}

func main() {
	maxCount := flag.Int("max", -1, "max count")
	flag.Parse()

	url := "https://raw.githubusercontent.com/lionsoul2014/ip2region/master/data/ip.merge.txt"
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	buf := bufio.NewReader(resp.Body)
	allIPInfos := make([][]uint32, 0)
	ipMax := *maxCount
	current := 0
	var currentIP uint32
	for {
		current++
		if ipMax != -1 && current > ipMax {
			break
		}
		line, _, err := buf.ReadLine()
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}
		arr := bytes.SplitN(line, []byte("|"), -1)
		if len(arr) != divideSize {
			log.Println("the line is invalid, " + string(line))
			continue
		}
		// 0 IP起始地址
		// 1 IP结束地址
		// 2 国家
		// 4 省份
		// 5 城市
		// 6 ISP
		country := string(arr[2])
		province := string(arr[4])
		city := string(arr[5])
		isp := string(arr[6])

		ipInfos := make([]uint32, 6)

		ipStart, err := convertIP(string(arr[0]))
		if err != nil {
			panic(err)
		}
		ipInfos[0] = ipStart
		if ipStart < currentIP {
			panic("ip address should be ordered")
		}
		currentIP = ipStart

		ipEnd, err := convertIP(string(arr[1]))
		if err != nil {
			panic(err)
		}
		if ipEnd < currentIP {
			panic("ip address should be ordered")
		}
		currentIP = ipEnd
		ipInfos[1] = ipEnd

		ipInfos[2] = addWord(country)
		ipInfos[3] = addWord(province)
		ipInfos[4] = addWord(city)
		ipInfos[5] = addWord(isp)
		allIPInfos = append(allIPInfos, ipInfos)
	}
	wordBuilder := new(strings.Builder)
	for _, word := range wordMapList {
		wordBuilder.WriteString(fmt.Sprintf(`"%s",`, word))
	}
	ipInfosBuilder := new(strings.Builder)
	for _, ipInfos := range allIPInfos {
		buf, _ := json.Marshal(ipInfos)
		buf[0] = '{'
		buf[len(buf)-1] = '}'
		ipInfosBuilder.WriteString(string(buf) + ",")
	}
	code := fmt.Sprintf(template, wordBuilder.String(), ipInfosBuilder.String(), len(allIPInfos))
	err = ioutil.WriteFile("./service/ip_data.go", []byte(code), 0644)
	if err != nil {
		panic(err)
	}
}
