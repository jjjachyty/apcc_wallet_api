package utils

import (
	"bytes"
	"crypto/md5"
	"encoding/gob"
	"encoding/hex"
	"net"
	"os"
	"os/user"
	"reflect"
	"runtime"

	"github.com/gin-gonic/gin"

	uuid "github.com/satori/go.uuid"
)

type responseData struct {
	Status bool
	// Code    string
	Message string
	Data    interface{}
}

//WriteMessage API返回消息函数
func Response(c *gin.Context, err error, data interface{}) {
	if err == nil {
		c.JSON(200, responseData{true, "请求成功", data})
	} else {
		c.JSON(200, responseData{false, err.Error(), data})

	}
}

//GetUUID 生成UUID 36位 中间带_
func GetUUID() string {
	uuid, _ := uuid.NewV4()
	return uuid.String()
}

type Exception struct {
	Code      string
	Mesg      string
	Exception error
}

//Data Code 为自定义代码 Mesg 自定义消息 Exception 原生错误 Data 返回数据
type SData struct {
	Except Exception
	Data   interface{}
}

//GetData Service层返回Action层的通用数据格式 code 为自定义代码 mesg 自定义消息 exception 原生错误 data 返回数据
func GetData(code string, mesg string, exception error, data ...interface{}) SData {
	if len(data) > 0 {
		return SData{Except: Exception{Code: code, Mesg: mesg, Exception: exception}, Data: data[0]}
	}
	return SData{Except: Exception{Code: code, Mesg: mesg, Exception: exception}, Data: nil}
}

//Struct2Map 实体转MAP
func Struct2Map(obj interface{}) map[string]interface{} {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)

	var data = make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		data[t.Field(i).Name] = v.Field(i).Interface()
	}
	return data
}

//DeepCopy 实体深度拷贝
func DeepCopy(dst, src interface{}) error {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(src); err != nil {
		return err
	}
	return gob.NewDecoder(bytes.NewBuffer(buf.Bytes())).Decode(dst)
}

//GetGID获取当前的协程ID
func GetGID() string {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	return string(b)
}

func GetMD5(orgStr string) string {
	h := md5.New()
	h.Write([]byte(orgStr))
	return hex.EncodeToString(h.Sum(nil))
}
func GetHostName() string {
	host, err := os.Hostname()
	if err != nil {
		SysLog.Errorln("获取本机名称错误")
	}
	return host
}
func GetUserName() string {
	if u, err := user.Current(); err == nil {
		return u.Name
	}
	return ""

}

func GetIP() string {
	addrs, err := net.InterfaceAddrs()

	if err != nil {
		SysLog.Errorf("获取本机IP错误 %v", err)
	}

	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return "0.0.0.0"
}
