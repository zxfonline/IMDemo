package hotword

import (
	"testing"
	"time"
)

func Test_Timetrie_GetStartWith(t *testing.T) {
	root := NewTimeTrie()
	root.Add("p")
	root.Add("z")
	root.Add("zx")
	root.Add("zxf")
	root.Add("zxfd")
	root.Add("zxfa")
	root.Add("zxs")
	root.Add("zxs3")
	root.Add("zxs2")
	root.Add("zy")
	now := time.Now().Unix()
	time.Sleep(2 * time.Second)
	root.Add("zy")
	root.Add("zyx")
	root.Add("zya")
	t.Log(root.GetStartWith("z"))
	root.OnTimeout(now + 1)
	t.Log(root.GetStartWith("z"))
	time.Sleep(2 * time.Second)
	root.OnTimeout(time.Now().Unix() + 1)
	t.Log(root.GetStartWith("z"))
}

func Test_Timetrie_HotTopX(t *testing.T) {
	root := NewTimeTrie()
	//2 scores
	root.Add("你好")
	root.Add("你好")
	//3 scores
	root.Add("你好吗")
	root.Add("你好吗")
	root.Add("你好吗")
	//1 scores
	root.Add("我")
	//1 scores
	root.Add("我好")
	//4 scores
	root.Add("我好着呢")
	root.Add("我好着呢")
	root.Add("我好着呢")
	root.Add("我好着呢")
	t.Log(root.HotTopX(3))
}
