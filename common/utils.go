package common

import (
	"bytes"
	"crypto/md5"
	"encoding/gob"
	"fmt"
	"strings"
	"time"

	"github.com/go-basic/uuid"
)

// GetMd5 获取 MD5
func GetMd5(data []byte) string {
	return fmt.Sprintf("%x", md5.Sum(data))
}

// Clone 深拷贝
func Clone(source, dest interface{}) error {
	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)
	dec := gob.NewDecoder(buf)
	if err := enc.Encode(source); err != nil {
		return err
	}
	if err := dec.Decode(dest); err != nil {
		return err
	}
	return nil
}

// GetTimeStamp 获取Unix时间戳 精度秒
func GetTimeStamp() int64 {
	return time.Now().Unix()
}

// GetCurDayTimeStamp 获取当前天时间戳
func GetCurDayTimeStamp() int64 {
	cstZone := time.FixedZone("CST", 8*3600)
	curTime := time.Now().In(cstZone)
	curDayUnix := time.Date(curTime.Year(), curTime.Month(), curTime.Day(), 0, 0, 0, 0, curTime.Location()).Unix()
	return curDayUnix
}

// GetTimeStr 获取字符串形式时间
// f: 20060102150405
// f: 2006-01-02 15:04:05
// f: 2006/01/02 15:04:05
func GetTimeStr(f string) string {
	return time.Unix(GetTimeStamp(), 0).Format(f)
}

// MergeMap 合并两个map
func MergeMap(maps ...map[string]interface{}) map[string]interface{} {
	newMap := map[string]interface{}{}
	for _, m := range maps {
		for k, v := range m {
			newMap[k] = v
		}
	}
	return newMap
}

// GetRandomString 获取随机字符串
func GetRandomString() string {
	return strings.ReplaceAll(uuid.New(), "-", "")
}

// HasRepeatElement 是否存在重复元素
func HasRepeatElement(array []interface{}) bool {
	count := make(map[interface{}]int, 0)
	for i := 0; i < len(array); i++ {
		count[array[i]]++
		if count[array[i]] > 1 {
			return true
		}
	}
	return false
}
