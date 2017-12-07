package api

import (
	"fmt"
	"bytes"
	"encoding/json"
	"time"
)

func PanicWith(a ... interface{}) {
	panic(fmt.Sprint(a...))
}

// 检测 []string 中是否包含特定的元素
func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// BoolAsString 转换 bool 到 string
func BoolAsString(status bool) string {

	m := map[bool]string{
		false: "0",
		true:  "1",
	}

	return m[status]
}

func jsonMarshalWithHTML(t interface{}) ([]byte, error) {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(t)
	return buffer.Bytes(), err
}

// GetSignedTime 获取 用户 oa 系统中的出席时间 时间戳
func GetSignedTime() int64 {
	// 先从 oa_attend 表里查找，有就获取 date+signIn 取回 秒数
	// 没有，就取当前时间秒数
	// 因为从 oa 剥离出来，所以直接去当前时间的时间戳
	return time.Now().Unix()
}

func EncrypFrom(data map[string]interface{}, debug bool) (d []byte, err error) {

	defer func() {
		if rec := recover(); rec != nil {
			d = nil
			err = fmt.Errorf("%v", rec)
		}
	}()

	encrypData, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	if debug {
		fmt.Printf("Respond: <<< %s\n", string(encrypData))
	}
	if d, err = AesEncrypt(encrypData, []byte(ServerToken)); err != nil {
		panic(err)
	}

	return d, err

}