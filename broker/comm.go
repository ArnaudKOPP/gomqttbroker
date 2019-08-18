package broker

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"io"
	"reflect"
	"time"
)

const (
	// AcceptMinSleep is the minimum acceptable sleep times on temporary errors.
	AcceptMinSleep = 100 * time.Millisecond
	// AcceptMaxSleep is the maximum acceptable sleep times on temporary errors
	AcceptMaxSleep = 10 * time.Second
	// DefaultRouteConnect Route solicitation intervals.
	DefaultRouteConnect = 5 * time.Second
	// DefaultTLSTimeout timeout of tls
	DefaultTLSTimeout = 5 * time.Second
)

const (
	// CONNECT 1
	CONNECT = uint8(iota + 1)
	// CONNACK 2
	CONNACK
	// PUBLISH 3
	PUBLISH
	// PUBACK 4
	PUBACK
	// PUBREC 5
	PUBREC
	// PUBREL 6
	PUBREL
	// PUBCOMP 7
	PUBCOMP
	// SUBSCRIBE 8
	SUBSCRIBE
	// SUBACK 9
	SUBACK
	// UNSUBSCRIBE 10
	UNSUBSCRIBE
	// UNSUBACK 11
	UNSUBACK
	// PINGREQ 12
	PINGREQ
	// PINGRESP 13
	PINGRESP
	// DISCONNECT 14
	DISCONNECT
)
const (
	// QosAtMostOnce 0
	QosAtMostOnce byte = iota
	// QosAtLeastOnce 1
	QosAtLeastOnce
	// QosExactlyOnce 2
	QosExactlyOnce
	// QosFailure fail value
	QosFailure = 0x80
)

func equal(k1, k2 interface{}) bool {
	if reflect.TypeOf(k1) != reflect.TypeOf(k2) {
		return false
	}

	if reflect.ValueOf(k1).Kind() == reflect.Func {
		return &k1 == &k2
	}

	if k1 == k2 {
		return true
	}
	switch k1 := k1.(type) {
	case string:
		return k1 == k2.(string)
	case int64:
		return k1 == k2.(int64)
	case int32:
		return k1 == k2.(int32)
	case int16:
		return k1 == k2.(int16)
	case int8:
		return k1 == k2.(int8)
	case int:
		return k1 == k2.(int)
	case float32:
		return k1 == k2.(float32)
	case float64:
		return k1 == k2.(float64)
	case uint:
		return k1 == k2.(uint)
	case uint8:
		return k1 == k2.(uint8)
	case uint16:
		return k1 == k2.(uint16)
	case uint32:
		return k1 == k2.(uint32)
	case uint64:
		return k1 == k2.(uint64)
	case uintptr:
		return k1 == k2.(uintptr)
	}
	return false
}

// GenUniqueID generate uid
func GenUniqueID() string {
	b := make([]byte, 48)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	h := md5.New()
	h.Write([]byte(base64.URLEncoding.EncodeToString(b)))
	return hex.EncodeToString(h.Sum(nil))
	// return GetMd5String()
}
