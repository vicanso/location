package service

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/vicanso/hes"
)

type (
	// IPLocation ip location
	IPLocation struct {
		Country  string `json:"country"`
		Province string `json:"province"`
		City     string `json:"city"`
		ISP      string `json:"isp"`
	}
)

var (
	errCategoryIPLocation = "ip-location"
	// ErrIPAddressInvalid ip address invalid
	ErrIPAddressInvalid = &hes.Error{
		Message:    "ip address is invalid",
		Category:   errCategoryIPLocation,
		StatusCode: http.StatusBadRequest,
	}
)

func getDesc(index uint32) string {
	v := WordMapList[int(index)]
	if v == "0" {
		return ""
	}
	return v
}

// GetLocationByIP get location by ip
func GetLocationByIP(ip uint32) (ipLocation *IPLocation) {
	start := 0
	end := IPCount - 1
	for {
		if start > end {
			break
		}
		mid := (end-start)/2 + start
		ipInfo := IPInfos[mid]
		ipStart := ipInfo[0]
		ipEnd := ipInfo[1]
		if ip < ipStart {
			end = mid - 1
		} else if ip > ipEnd {
			start = mid + 1
		} else {
			ipLocation = &IPLocation{
				Country:  getDesc(ipInfo[2]),
				Province: getDesc(ipInfo[3]),
				City:     getDesc(ipInfo[4]),
				ISP:      getDesc(ipInfo[5]),
			}
			break
		}
	}
	return
}

// ConvertIP2Uint32 convert ip to uint32
func ConvertIP2Uint32(ip string) (value uint32, err error) {
	arr := strings.SplitN(ip, ".", -1)
	if len(arr) != 4 {
		err = ErrIPAddressInvalid
		return
	}
	offset := 8
	max := 3
	for index, item := range arr {
		v, e := strconv.Atoi(item)
		if e != nil {
			he := hes.New(e.Error())
			he.Category = errCategoryIPLocation
			return 0, he
		}
		if v < 0 || v > 255 {
			err = ErrIPAddressInvalid
			return 0, err
		}
		value += uint32(v) << uint32(offset*(max-index))
	}
	return
}
