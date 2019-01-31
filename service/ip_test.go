package service

import (
	"testing"
)

func TestConvertIP2Uint32(t *testing.T) {
	t.Run("valid ip", func(t *testing.T) {
		ip, err := ConvertIP2Uint32("1.2.127.255")
		if err != nil || ip != 16941055 {
			t.Fatalf("convert ip to uint32 fail, %v", err)
		}
	})

	t.Run("invalid ip", func(t *testing.T) {
		_, err := ConvertIP2Uint32("1.2.127")
		if err != ErrIPAddressInvalid {
			t.Fatalf("invalid ip should return error")
		}
		_, err = ConvertIP2Uint32("1.2.127.256")
		if err != ErrIPAddressInvalid {
			t.Fatalf("invalid ip should return error")
		}
		_, err = ConvertIP2Uint32("1.2.127.a")
		if err == nil {
			t.Fatalf("invalid ip should return error")
		}

	})

}

func TestGetLocationByIP(t *testing.T) {

	fn := func(t *testing.T, ipAddr string) *IPLocation {
		ip, err := ConvertIP2Uint32(ipAddr)
		if err != nil {
			t.Fatalf("convert ip fail")
		}
		return GetLocationByIP(ip)
	}

	t.Run("0.255.255.255", func(t *testing.T) {
		info := fn(t, "0.255.255.255")
		if info.ISP != "内网IP" {
			t.Fatalf("get location fail")
		}
	})

	t.Run("1.2.7.255", func(t *testing.T) {
		info := fn(t, "1.2.7.255")
		if info.Country != "中国" ||
			info.Province != "福建省" ||
			info.City != "福州市" ||
			info.ISP != "电信" {
			t.Fatalf("get location fail")
		}
	})

	t.Run("first ip", func(t *testing.T) {
		info := fn(t, "0.0.0.0")
		if info.ISP != "内网IP" {
			t.Fatalf("get the first location fail")
		}
	})
	t.Run("last ip", func(t *testing.T) {
		info := fn(t, "1.3.255.255")
		if info.Country != "中国" ||
			info.Province != "广东省" ||
			info.City != "广州市" ||
			info.ISP != "电信" {
			t.Fatalf("get location fail")
		}
	})
}

func BenchmarkConvertIP2Uint32(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		ConvertIP2Uint32("1.2.7.255")
	}
}

func BenchmarkGetLocationByIP(b *testing.B) {
	b.ReportAllocs()
	ip, _ := ConvertIP2Uint32("1.2.7.255")
	for i := 0; i < b.N; i++ {
		GetLocationByIP(ip)
	}
}
