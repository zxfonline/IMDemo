package nametrie

import "sync"

//默认空间
const INITCAP = 8

type NameMatchTrie struct {
	sync.Mutex
	key  rune
	next map[rune]*NameMatchTrie
	word bool
}

func NewNameMatchTrie() *NameMatchTrie {
	return &NameMatchTrie{next: make(map[rune]*NameMatchTrie, INITCAP)}
}

//添加记录
func (t *NameMatchTrie) Add(words string) {
	t.Lock()
	defer t.Unlock()
	for _, c := range words {
		if t.next[c] == nil {
			st := NewNameMatchTrie()
			st.key = c
			t.next[c] = st
		}
		t = t.next[c]
	}
	t.word = true
}

//查询是否完整包含words的记录
func (t *NameMatchTrie) FullMatch(words string) bool {
	t.Lock()
	defer t.Unlock()
	for _, c := range words {
		if t.next[c] == nil {
			return false
		}
		t = t.next[c]
	}
	return t.word
}
