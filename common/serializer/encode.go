package serializer

import (
	"bytes"
	"encoding/gob"
)

// GobEncode 使用 gob 序列化指定的对象
func GobEncode(src interface{}) ([]byte, error) {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	if err := encoder.Encode(src); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// GobDecode 使用 gob 反序列化指定的对象
func GobDecode(dest interface{}, buf []byte) error {
	decoder := gob.NewDecoder(bytes.NewBuffer(buf))
	if err := decoder.Decode(dest); err != nil {
		return err
	}
	return nil
}
