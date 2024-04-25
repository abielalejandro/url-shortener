package utils

import (
	"math/big"
	"net/http"
)

func ReadUserIP(r *http.Request) string {
	IPAddress := r.Header.Get("X-Real-Ip")
	if IPAddress == "" {
		IPAddress = r.Header.Get("X-Forwarded-For")
	}
	if IPAddress == "" {
		IPAddress = r.RemoteAddr
	}
	return IPAddress
}

func ToBase62(str string) string {
	var i big.Int
	i.SetBytes([]byte(str))
	return i.Text(62)
}
