package badword

import (
	"bufio"
	"io"
	"os"
	"unicode/utf8"

	"github.com/zxfonline/IMDemo/core/fileutil"
)

var _G *BadWordTrie

func init() {
	_G = NewBadWordTrie()
	fi, err := fileutil.FindFile("runtime/badword.txt", os.O_RDONLY, 0)
	if err != nil {
		panic(err)
	}
	defer fi.Close()
	br := bufio.NewReader(fi)
	for {
		a, _, c := br.ReadLine()
		if c == io.EOF {
			break
		}
		_G.Add(string(a))
	}
}

//关键字查询
func BadWordSearch(str string) bool {
	return _G.Search(str)
}

//关键字替换
func BadWordReplace(str string) string {
	return _G.Replace(str)
}

//设置全局关键字
func BadWordGolbal(key *BadWordTrie) {
	_G = key
}

type BadWordTrie struct {
	next map[rune]*BadWordTrie
	word bool
}

func NewBadWordTrie() *BadWordTrie {
	return &BadWordTrie{next: make(map[rune]*BadWordTrie)}
}

//添加脏字
func (bt *BadWordTrie) Add(badword string) {
	for _, c := range badword {
		if bt.next[c] == nil {
			bt.next[c] = NewBadWordTrie()
		}
		bt = bt.next[c]
	}
	bt.word = true
}

func (bt *BadWordTrie) Search(words string) bool {
	sbt := bt
	key := []rune(words)
	var chars []rune
	slen := len(key)
	for i := 0; i < slen; i++ {
		if _, exists := sbt.next[key[i]]; exists {
			sbt = sbt.next[key[i]]
			for j := i + 1; j < slen; j++ {
				if _, exists := sbt.next[key[j]]; exists {
					sbt = sbt.next[key[j]]
					if sbt.word {
						if chars == nil {
							chars = key
						}
						for t := i; t <= j; t++ {
							return true
						}
						i = j
						sbt = bt
						break
					}
				}
			}
			sbt = bt
		}
	}
	return false
}

func (bt *BadWordTrie) Replace(words string) string {
	sbt := bt
	key := []rune(words)
	var chars []rune
	slen := len(key)
	for i := 0; i < slen; i++ {
		if _, exists := sbt.next[key[i]]; exists {
			sbt = sbt.next[key[i]]
			for j := i + 1; j < slen; j++ {
				if _, exists := sbt.next[key[j]]; exists {
					sbt = sbt.next[key[j]]
					if sbt.word {
						if chars == nil {
							chars = key
						}
						for t := i; t <= j; t++ {
							c, _ := utf8.DecodeRuneInString("*")
							chars[t] = c
						}
						i = j
						sbt = bt
						break
					}
				}
			}
			sbt = bt
		}
	}
	if chars == nil {
		return words
	} else {
		return string(chars)
	}
}
