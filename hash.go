package confident

import (
	"crypto/md5"
	"encoding/json"
)

func CalculateHash(value interface{}) [16]byte {
	arrBytes := []byte{}
	jsonBytes, _ := json.Marshal(value)
	arrBytes = append(arrBytes, jsonBytes...)
	return md5.Sum(arrBytes)
}
