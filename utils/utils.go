package utils

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"hash"
	"math"
	"time"
)

func NowTSNorm() int32 {
	ts := time.Now().Unix()
	return int32(math.Floor(float64(ts)/86400) * 86400)
}

func NowTS() int {
	ts := time.Now().Unix()
	return int(float64(ts))
}

func NowMonth() string {
	date_time := time.Now().Format("2006-01")
	return date_time
}

func NowDate() string {
	date_time := time.Now().Format("2006-01-01")
	return date_time
}

func MD5SHA1(link string, h hash.Hash, hasher hash.Hash) string {
	hasher.Write([]byte(link))
	_id := hex.EncodeToString(hasher.Sum(nil))
	h.Write([]byte(_id))
	bs := hex.EncodeToString(h.Sum(nil))
	return bs
}

func SpecialID(link string) string {
	h := sha1.New()
	hasher := md5.New()
	hasher.Write([]byte(link))
	_id := hex.EncodeToString(hasher.Sum(nil))
	h.Write([]byte(_id))
	bs := hex.EncodeToString(h.Sum(nil))
	return bs[:24]
}
