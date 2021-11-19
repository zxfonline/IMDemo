package nametrie

import (
	"testing"
)

func Test_Timetrie_HotTopX(t *testing.T) {
	root := NewNameMatchTrie()
	root.Add("张三")
	root.Add("李四")
	root.Add("zhangsan")
	t.Log(root.FullMatch("李三四"))
	t.Log(root.FullMatch("李四"))
	t.Log(root.FullMatch("ZHANGSAN"))
	t.Log(root.FullMatch("zhangsan"))
	t.Log(root.FullMatch("zhangsan2"))
}
