/*
	构建一个支持统计有效时间范围内的热门记录 Trie树
*/
package hotword

import (
	"container/heap"
	"strings"
	"time"
)

//默认空间
const INITCAP = 8

//有效时长 s
const ValidSeconds = 60 * 10

type TimeTrie struct {
	key    rune
	next   map[rune]*TimeTrie
	parent *TimeTrie
	//关键词录入的时间戳列表
	//方便定时器清理过期的统计
	recordTss []int64
}

func NewTimeTrie() *TimeTrie {
	return &TimeTrie{next: make(map[rune]*TimeTrie, INITCAP)}
}

//添加记录
func (t *TimeTrie) Add(words string) {
	t.AddWithTime(words, time.Now().Unix())
}

//添加记录
func (t *TimeTrie) AddWithTime(words string, timeUnix int64) {
	for _, c := range words {
		if t.next[c] == nil {
			st := NewTimeTrie()
			st.key = c
			st.parent = t
			t.next[c] = st
		}
		t = t.next[c]
	}
	t.recordTss = append(t.recordTss, timeUnix)
}

//查询是否完整包含words的记录
func (t *TimeTrie) FullMatch(words string) bool {
	for _, c := range words {
		if t.next[c] == nil {
			return false
		}
		t = t.next[c]
	}
	return t.isWord()
}

func (t *TimeTrie) isWord() bool {
	return len(t.recordTss) > 0
}

func (t *TimeTrie) Points() int {
	return len(t.recordTss)
}

//查询是否有包含prefix前缀的记录
func (t *TimeTrie) StartWith(prefix string) bool {
	for _, v := range prefix {
		if t.next[v] == nil {
			return false
		}
		t = t.next[v]
	}
	return true
}

func (t *TimeTrie) getCompletionInfo(founds *[]string, prefix string) {
	if t.isWord() {
		*founds = append(*founds, prefix)
	}
	for k, st := range t.next {
		st.getCompletionInfo(founds, prefix+string(k))
	}
}

//查询包含prefix前缀的所有记录
func (t *TimeTrie) GetStartWith(prefix string) (founds []string) {
	for _, v := range prefix {
		if t.next[v] == nil {
			return
		}
		t = t.next[v]
	}
	t.getCompletionInfo(&founds, prefix)
	return
}

func (t *TimeTrie) clearOutTimeRecord(outTimeUnix int64) {
	for i := len(t.recordTss) - 1; i >= 0; i-- {
		if t.recordTss[i] < outTimeUnix {
			t.recordTss = t.recordTss[i+1:]
			return
		}
	}
}
func (t *TimeTrie) cleanOutTimeTrie() {
	if t.parent == nil {
		return
	}
	if !t.isWord() && len(t.next) == 0 {
		delete(t.parent.next, t.key)
		t.parent.cleanOutTimeTrie()
	}
}

//清理过期的记录数据，降低内存占用
//outTimeUnix 清理缓存的过期数据时间戳
func (t *TimeTrie) OnTimeout(outTimeUnix int64) {
	for _, s := range t.next {
		s.clearOutTimeRecord(outTimeUnix)
		s.cleanOutTimeTrie()
		s.OnTimeout(outTimeUnix)
	}
}

// 统计当前缓存前x的热门
func (t *TimeTrie) HotTopX(topX int) []string {
	txInfo := NewTopKInfo(topX, nil)
	for _, s := range t.next {
		s.statisTopKInfo(txInfo)
	}
	hotestX := make([]string, 0, topX)
	for txInfo.MinHeap.Len() > 0 {
		tn := heap.Pop(&txInfo.MinHeap).(*TimeTrie)
		hotestX = append(hotestX, tn.printFullTxt())
	}
	for i, j := 0, len(hotestX)-1; i < j; i, j = i+1, j-1 {
		hotestX[i], hotestX[j] = hotestX[j], hotestX[i]
	}
	return hotestX
}

func (t *TimeTrie) printFullTxt() string {
	var sb strings.Builder
	sb.WriteRune(t.key)
	for t.parent != nil && t.parent.key != 0 {
		sb.WriteRune(t.parent.key)
		t = t.parent
	}
	return Reverse(sb.String())
}

func (t *TimeTrie) statisTopKInfo(txInfo *TopKInfo) {
	if t.isWord() {
		txInfo.Add(t)
	}
	for _, st := range t.next {
		st.statisTopKInfo(txInfo)
	}
}

//字符串反转
func Reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	v := string(runes)
	return v
}
