package utils

import (
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"encoding"
	"encoding/json"
	"fmt"
	"io"
	"os"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func FileExists(filePath string) bool {
	info, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func GUnzipData(data []byte) []byte {
	b := bytes.NewBuffer(data)
	var r io.Reader
	r, err := gzip.NewReader(b)
	check(err)

	var resB bytes.Buffer
	resB.ReadFrom(r)
	check(err)
	return resB.Bytes()
}

func GZipData(data []byte) []byte {
	var b bytes.Buffer
	gz := gzip.NewWriter(&b)
	_, err := gz.Write(data)
	check(err)
	if err = gz.Flush(); err != nil {
		panic(err)
	}
	if err = gz.Close(); err != nil {
		panic(err)
	}
	compressedData := b.Bytes()
	return compressedData
}

func Hash(input string) string {
	hash := sha256.New()
	hash.Write([]byte(input))
	val := fmt.Sprintf("%x", hash.Sum(nil))
	marshaler, ok := hash.(encoding.BinaryMarshaler)
	if !ok {
		panic(ok)
	}
	_, err := marshaler.MarshalBinary()
	if err != nil {
	}
	return val

}

//func MarshalYaml(obj) []byte {
//	return yaml.Marshal(&obj)
//}
func JSONMarshal(obj interface{}) []byte {
	bytes, err := json.Marshal(obj)
	check(err)
	return bytes
}

func JSONUnmarshal(data *[]byte) interface{} {
	var obj interface{}
	err := json.Unmarshal(*data, &obj)
	check(err)
	return obj
}

func ByteCountSI(b int64) string {
	const unit = 1000
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB",
		float64(b)/float64(div), "kMGTPE"[exp])
}
