package functest

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/valyala/fastjson"
)

func TestParsez(t *testing.T) {
	ss := `{"type":1,"data":{"name":"zzz"}}`
	v, err := fastjson.Parse(ss)
	if err != nil {
		t.Error(err)
	}
	typz := v.GetUint("type")
	t.Logf("type:%d", typz)
	name := v.GetStringBytes("data", "name")
	t.Logf("name:%s", name)
	name1 := fastjson.GetString([]byte(ss), "data", "name")
	t.Logf("name1:%s", name1)
}

func TestPop(t *testing.T) {
	arr := make([]int, 0, 10)
	for i := 0; i < 30; i++ {
		if len(arr) >= cap(arr) {
			copy(arr[:], arr[1:])
			arr[cap(arr)-1] = i
		} else {
			arr = append(arr, i)
		}
		t.Logf("v=%+v,len:%d,cap:%d", arr, len(arr), cap(arr))
	}
	t.Logf("====v=%+v,len:%d,cap:%d", arr, len(arr), cap(arr))
}

func TestMin(t *testing.T) {
	//一分钟的秒数
	var interval int64 = 60
	//距离下次到时的秒数
	now := time.Now()
	realExpire := interval - (now.Unix() % interval)
	t.Logf("now:%s,e:%d", now.Format("2006-01-02 15:04:05"), realExpire)
}

type Response struct {
	Type int `json:"type"`
	//=0:默认成功码
	//<>0其他错误码
	Code    int         `json:"code"`
	CodeMsg string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"` // 数据 json
}

func (r *Response) toJson() []byte {
	if byteData, err := json.Marshal(r); err == nil {
		return byteData
	}
	return nil
}

func TestJsonMarshal(t *testing.T) {
	rsp := &Response{
		Type: 1,
		Code: 0,
		Data: &struct {
			RoomID int64 `json:"roomID"`
			// UserName string `json:"userName"`
		}{
			RoomID: 3,
			// UserName: "zhangsan",
		},
	}
	str := string((rsp).toJson())
	v := fastjson.MustParse(str)
	name := "lisi"
	v.Get("data").Set("userName", fastjson.MustParse(fmt.Sprintf("%q", name)))
	v.Get("data").Set("sendTime", fastjson.MustParse(fmt.Sprintf("%q", time.Now().Format("2006-01-02 15:04:05"))))
	v.Set("type", fastjson.MustParse("2"))
	t.Log(str)
	t.Log(v.String())
}
