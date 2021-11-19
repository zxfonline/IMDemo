package strutil

import (
	"fmt"
	"math"
	"math/rand"
	"reflect"
	"strconv"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"
	"unsafe"
)

var (
	// Strings accepted as boolean.
	boolString = map[string]bool{
		"t":     true,
		"true":  true,
		"y":     true,
		"yes":   true,
		"on":    true,
		"1":     true,
		"f":     false,
		"false": false,
		"n":     false,
		"no":    false,
		"off":   false,
		"0":     false,
	}
)

//Stou64 没找到并且没有默认值则 panic抛错
func Stou64(v string, def ...uint64) uint64 {
	if n, err := strconv.ParseUint(strings.TrimSpace(v), 10, 64); err == nil {
		return n
	} else if len(def) > 0 {
		return def[0]
	} else {
		panic(fmt.Errorf("invalid uint64 param:%v ,err:%v", v, err))
	}
}

//Stou32 没找到并且没有默认值则 panic抛错
func Stou32(v string, def ...uint32) uint32 {
	if n, err := strconv.ParseUint(strings.TrimSpace(v), 10, 64); err == nil {
		return uint32(n)
	} else if len(def) > 0 {
		return def[0]
	} else {
		panic(fmt.Errorf("invalid uint32 param:%v ,err:%v", v, err))
	}
}

//Stou16 没找到并且没有默认值则 panic抛错
func Stou16(v string, def ...uint16) uint16 {
	if n, err := strconv.ParseUint(strings.TrimSpace(v), 10, 64); err == nil {
		return uint16(n)
	} else if len(def) > 0 {
		return def[0]
	} else {
		panic(fmt.Errorf("invalid uint16 param:%v ,err:%v", v, err))
	}
}

//Stou8 没找到并且没有默认值则 panic抛错
func Stou8(v string, def ...uint8) uint8 {
	if n, err := strconv.ParseUint(strings.TrimSpace(v), 10, 64); err == nil {
		return uint8(n)
	} else if len(def) > 0 {
		return def[0]
	} else {
		panic(fmt.Errorf("invalid uint8 param:%v ,err:%v", v, err))
	}
}

//Stou 没找到并且没有默认值则 panic抛错
func Stou(v string, def ...uint) uint {
	if n, err := strconv.ParseUint(strings.TrimSpace(v), 10, 64); err == nil {
		return uint(n)
	} else if len(def) > 0 {
		return def[0]
	} else {
		panic(fmt.Errorf("invalid uint param:%v ,err:%v", v, err))
	}
}

//Stoi64 没找到并且没有默认值则 panic抛错
func Stoi64(v string, def ...int64) int64 {
	if n, err := strconv.ParseInt(strings.TrimSpace(v), 0, 0); err == nil {
		return n
	} else if len(def) > 0 {
		return def[0]
	} else {
		panic(fmt.Errorf("invalid int64 param:%v ,err:%v", v, err))
	}
}

//Stoi32 没找到并且没有默认值则 panic抛错
func Stoi32(v string, def ...int32) int32 {
	if n, err := strconv.ParseInt(strings.TrimSpace(v), 0, 0); err == nil {
		return int32(n)
	} else if len(def) > 0 {
		return def[0]
	} else {
		panic(fmt.Errorf("invalid int32 param:%v ,err:%v", v, err))
	}
}

//Stoi16 没找到并且没有默认值则 panic抛错
func Stoi16(v string, def ...int16) int16 {
	if n, err := strconv.ParseInt(strings.TrimSpace(v), 0, 0); err == nil {
		return int16(n)
	} else if len(def) > 0 {
		return def[0]
	} else {
		panic(fmt.Errorf("invalid int16 param:%v ,err:%v", v, err))
	}
}

//Stoi8 没找到并且没有默认值则 panic抛错
func Stoi8(v string, def ...int8) int8 {
	if n, err := strconv.ParseInt(strings.TrimSpace(v), 0, 0); err == nil {
		return int8(n)
	} else if len(def) > 0 {
		return def[0]
	} else {
		panic(fmt.Errorf("invalid int8 param:%v ,err:%v", v, err))
	}
}

//Stoi 没找到并且没有默认值则 panic抛错
func Stoi(v string, def ...int) int {
	if n, err := strconv.ParseInt(strings.TrimSpace(v), 0, 0); err == nil {
		return int(n)
	} else if len(def) > 0 {
		return def[0]
	} else {
		panic(fmt.Errorf("invalid int param:%v ,err:%v", v, err))
	}
}

//StoBol 没找到并且没有默认值则 panic抛错
func StoBol(v string, def ...bool) bool {
	if value, ok := boolString[strings.ToLower(strings.TrimSpace(v))]; ok {
		return value
	} else if len(def) > 0 {
		return def[0]
	} else {
		panic(fmt.Errorf("invalid bool param:%v", v))
	}
}

//Stof32 没找到并且没有默认值则 panic抛错
func Stof32(v string, def ...float32) float32 {
	if val, err := strconv.ParseFloat(strings.TrimSpace(v), 64); err == nil {
		return float32(val)
	} else if len(def) > 0 {
		return def[0]
	} else {
		panic(fmt.Errorf("invalid float32 param:%v ,err:%v", v, err))
	}
}

//Stof64 没找到并且没有默认值则 panic抛错
func Stof64(v string, def ...float64) float64 {
	if val, err := strconv.ParseFloat(strings.TrimSpace(v), 64); err == nil {
		return val
	} else if len(def) > 0 {
		return def[0]
	} else {
		panic(fmt.Errorf("invalid float64 param:%v ,err:%v", v, err))
	}
}

func Radians(degrees float64) float64 {
	return degrees * math.Pi / 180
}

func Degrees(radians float64) float64 {
	return radians * 180 / math.Pi
}

func Median(items []float64) float64 {
	n := len(items)
	switch {
	case n == 0:
		return 0
	case n%2 == 1:
		return items[n/2]
	default:
		a := items[n/2-1]
		b := items[n/2]
		return (a + b) / 2
	}
}

func DurationString(d time.Duration) string {
	h := int(d.Hours())
	m := int(d.Minutes()) % 60
	s := int(d.Seconds()) % 60
	return fmt.Sprintf("%d:%02d:%02d", h, m, s)
}

//ParseFloat64s 忽略空字符串
func ParseFloat64s(items []string) ([]float64, error) {
	result := make([]float64, 0, len(items))
	for i := 0; i < len(items); i++ {
		t := strings.TrimSpace(items[i])
		if len(t) != 0 {
			if f, err := strconv.ParseFloat(t, 64); err == nil {
				result = append(result, f)
			} else {
				return nil, fmt.Errorf("invalid []float64 param:%v ,err:%v", items, err)
			}
		}
	}
	return result, nil
}

//ParseFloat32s 忽略空字符串
func ParseFloat32s(items []string) ([]float32, error) {
	result := make([]float32, 0, len(items))
	for i := 0; i < len(items); i++ {
		t := strings.TrimSpace(items[i])
		if len(t) != 0 {
			if f, err := strconv.ParseFloat(t, 64); err == nil {
				result = append(result, float32(f))
			} else {
				return nil, fmt.Errorf("invalid []float32 param:%v ,err:%v", items, err)
			}
		}
	}
	return result, nil
}

//ParseInts 忽略空字符串
func ParseInts(items []string) ([]int, error) {
	result := make([]int, 0, len(items))
	for i := 0; i < len(items); i++ {
		t := strings.TrimSpace(items[i])
		if len(t) != 0 {
			if f, err := strconv.ParseInt(t, 0, 0); err == nil {
				result = append(result, int(f))
			} else {
				return nil, fmt.Errorf("invalid []int param:%v ,err:%v", items, err)
			}
		}
	}
	return result, nil
}

//ParseInt64s 忽略空字符串
func ParseInt64s(items []string) ([]int64, error) {
	result := make([]int64, 0, len(items))
	for i := 0; i < len(items); i++ {
		t := strings.TrimSpace(items[i])
		if len(t) != 0 {
			if f, err := strconv.ParseInt(t, 0, 0); err == nil {
				result = append(result, f)
			} else {
				return nil, fmt.Errorf("invalid []int64 param:%v ,err:%v", items, err)
			}
		}
	}
	return result, nil
}

//ParseUint64s 忽略空字符串
func ParseUint64s(items []string) ([]uint64, error) {
	result := make([]uint64, 0, len(items))
	for i := 0; i < len(items); i++ {
		t := strings.TrimSpace(items[i])
		if len(t) != 0 {
			if f, err := strconv.ParseUint(t, 0, 0); err == nil {
				result = append(result, f)
			} else {
				return nil, fmt.Errorf("invalid []uint64 param:%v ,err:%v", items, err)
			}
		}
	}
	return result, nil
}

//ParseInt32s 忽略空字符串
func ParseInt32s(items []string) ([]int32, error) {
	result := make([]int32, 0, len(items))
	for i := 0; i < len(items); i++ {
		t := strings.TrimSpace(items[i])
		if len(t) != 0 {
			if f, err := strconv.ParseInt(t, 0, 0); err == nil {
				result = append(result, int32(f))
			} else {
				return nil, fmt.Errorf("invalid []int32 param:%v ,err:%v", items, err)
			}
		}
	}
	return result, nil
}

//ParseInt16s 忽略空字符串
func ParseInt16s(items []string) ([]int16, error) {
	result := make([]int16, len(items))
	for i := 0; i < len(items); i++ {
		t := strings.TrimSpace(items[i])
		if len(t) != 0 {
			if f, err := strconv.ParseInt(t, 0, 0); err == nil {
				result = append(result, int16(f))
			} else {
				return nil, fmt.Errorf("invalid []int16 param:%v ,err:%v", items, err)
			}
		}
	}
	return result, nil
}

//ParseInt8s 忽略空字符串
func ParseInt8s(items []string) ([]int8, error) {
	result := make([]int8, len(items))
	for i := 0; i < len(items); i++ {
		t := strings.TrimSpace(items[i])
		if len(t) != 0 {
			if f, err := strconv.ParseInt(t, 0, 0); err == nil {
				result = append(result, int8(f))
			} else {
				return nil, fmt.Errorf("invalid []int8 param:%v ,err:%v", items, err)
			}
		}
	}
	return result, nil
}

//ParseBools 忽略空字符串
func ParseBools(items []string) ([]bool, error) {
	result := make([]bool, 0, len(items))
	for i := 0; i < len(items); i++ {
		t := strings.TrimSpace(items[i])
		if len(t) != 0 {
			if f, ok := boolString[strings.ToLower(t)]; ok {
				result = append(result, f)
			} else {
				return nil, fmt.Errorf("invalid []bool param:%v", items)
			}
		}
	}
	return result, nil
}

func Float64sToStrs(items []float64) []string {
	result := make([]string, len(items))
	for i, item := range items {
		result[i] = fmt.Sprintf("%v", item)
	}
	return result
}

func Float32sToStrs(items []float32) []string {
	result := make([]string, len(items))
	for i, item := range items {
		result[i] = fmt.Sprintf("%v", item)
	}
	return result
}

func IntsToStrs(items []int) []string {
	result := make([]string, len(items))
	for i, item := range items {
		result[i] = fmt.Sprintf("%v", item)
	}
	return result
}

func Int64sToStrs(items []int64) []string {
	result := make([]string, len(items))
	for i, item := range items {
		result[i] = fmt.Sprintf("%v", item)
	}
	return result
}

func Int32sToStrs(items []int32) []string {
	result := make([]string, len(items))
	for i, item := range items {
		result[i] = fmt.Sprintf("%v", item)
	}
	return result
}

func Int16sToStrs(items []int16) []string {
	result := make([]string, len(items))
	for i, item := range items {
		result[i] = fmt.Sprintf("%v", item)
	}
	return result
}

func Int8sToStrs(items []int8) []string {
	result := make([]string, len(items))
	for i, item := range items {
		result[i] = fmt.Sprintf("%v", item)
	}
	return result
}

func BoolsToStrs(items []bool) []string {
	result := make([]string, len(items))
	for i, item := range items {
		result[i] = fmt.Sprintf("%v", item)
	}
	return result
}

func M2string(i interface{}) string {
	if i == nil {
		return ""
	}
	switch t := i.(type) {
	case string:
		return t
	case []byte:
		return string(t)
	default:
		return fmt.Sprintf("%v", i)
	}
}

func M2bool(i interface{}) bool {
	if i == nil {
		return false
	}
	switch t := i.(type) {
	case bool:
		return t
	default:
		return StoBol(M2string(i))
	}
}

func M2float64(i interface{}) float64 {
	if i == nil {
		return 0
	}
	switch t := i.(type) {
	case float64:
		return t
	case float32:
		return float64(t)
	case int64:
		return float64(t)
	case uint64:
		return float64(t)
	case int32:
		return float64(t)
	case uint32:
		return float64(t)
	case int16:
		return float64(t)
	case uint16:
		return float64(t)
	case int8:
		return float64(t)
	case uint8:
		return float64(t)
	case int:
		return float64(t)
	case uint:
		return float64(t)
	default:
		return Stof64(M2string(i), 0)
	}
}
func M2float32(i interface{}) float32 {
	return float32(M2float64(i))
}

func M2uint64(i interface{}) uint64 {
	if i == nil {
		return 0
	}
	switch t := i.(type) {
	case int8:
		return uint64(t)
	case int16:
		return uint64(t)
	case int32:
		return uint64(t)
	case int:
		return uint64(t)
	case int64:
		return uint64(i.(int64))

	case uint8:
		return uint64(t)
	case uint16:
		return uint64(t)
	case uint32:
		return uint64(t)
	case uint:
		return uint64(t)
	case uint64:
		return i.(uint64)

	case float64:
		return uint64(t)
	case float32:
		return uint64(t)
	default:
		return Stou64(M2string(i), 0)
	}
}

func M2int64(i interface{}) int64 {
	return int64(M2uint64(i))
}

func M2int32(i interface{}) int32 {
	return int32(M2uint64(i))
}

func M2uint32(i interface{}) uint32 {
	return uint32(M2uint64(i))
}

func M2int16(i interface{}) int16 {
	return int16(M2uint64(i))
}

func M2uint16(i interface{}) uint16 {
	return uint16(M2uint64(i))
}

func M2int8(i interface{}) int8 {
	return int8(M2uint64(i))
}

func M2uint8(i interface{}) uint8 {
	return uint8(M2uint64(i))
}

func Fract(x float64) float64 {
	_, x = math.Modf(x)
	return x
}

func Clamp(x, lo, hi float64) float64 {
	if x < lo {
		return lo
	}
	if x > hi {
		return hi
	}
	return x
}

func ClampInt(x, lo, hi int) int {
	if x < lo {
		return lo
	}
	if x > hi {
		return hi
	}
	return x
}

func ShowProgress(start time.Time, rays uint64, crtValue, totalValue int) {
	pct := int(100 * float64(crtValue) / float64(totalValue))
	elapsed := time.Since(start)
	fmt.Printf("\r%4d / %d (%3d%%) [", crtValue, totalValue, pct)
	for p := 0; p < 100; p += 3 {
		if pct > p {
			fmt.Print("=")
		} else {
			fmt.Print(" ")
		}
	}
	fmt.Printf("] %s ", DurationString(elapsed))
}

func FirstIsLower(s string) bool {
	if r := s[0]; r < utf8.RuneSelf {
		return 'a' <= r && r <= 'z'
	}
	r, _ := utf8.DecodeRuneInString(s)
	return unicode.IsLower(r)
}

func FirstIsUpper(s string) bool {
	if r := s[0]; r < utf8.RuneSelf {
		return 'A' <= r && r <= 'Z'
	}
	r, _ := utf8.DecodeRuneInString(s)
	return unicode.IsUpper(r)
}

func FirstToUpper(s string) string {
	if s == "" {
		return ""
	}
	r, n := utf8.DecodeRuneInString(s)
	if r == '_' {
		return s
	}
	return string(unicode.ToUpper(r)) + s[n:]
}

func FirstToLower(s string) string {
	if s == "" {
		return ""
	}
	r, n := utf8.DecodeRuneInString(s)
	if r == '_' {
		return s
	}
	return string(unicode.ToLower(r)) + s[n:]
}

var (
	validChars = []*unicode.RangeTable{
		unicode.Han,               //汉语 中文集
		unicode.Hangul,            //旧朝鲜文
		unicode.Unified_Ideograph, //中日韩统一表意文字
		unicode.Thai,              //泰文
		unicode.Latin,             //中欧语言：克罗地亚语、捷克语、匈牙利语、波兰语、斯洛伐克语、斯洛文尼亚语、上索布语、下索布语、阿尔巴尼亚语、英语、德语、拉丁语、法语、西班牙语、葡萄牙、意大利、荷兰、印尼
		unicode.Cyrillic,          //西里尔字母：斯拉夫语族的语言，包括俄语、乌克兰语、卢森尼亚语、白俄罗斯语、保加利亚语、塞尔维亚语、马其顿语等
		unicode.Old_Italic,        //古意大利文
		unicode.Old_Turkic,        //古突厥文（土耳其）
		unicode.Hiragana,          //日语：平假名
		unicode.Katakana,          //日语：片假名
		unicode.Upper,             //大写
		unicode.Lower,             //小写
		unicode.Digit,             //十进制数字
	}
)

func CheckVaildName(name string, minLen, maxLen int) bool {
	r := []rune(name)
	length := len(r)
	//1.长度是否合法
	if length < minLen || length > maxLen {
		return false
	}
	//2.是否包含非法字符
	for _, c := range r {
		if n := utf8.RuneLen(c); n > 3 || n == -1 {
			return false
		}
		switch {
		case c >= 'a' && c <= 'z':
		case c >= 'A' && c <= 'Z':
		case c >= '0' && c <= '9':
		case unicode.IsOneOf(validChars, c):
		default:
			return false
		}
	}
	return true
}

//可以包含空格，字符串首尾不能有空格
func CheckVaildName1(name string, minLen, maxLen int) bool {
	r := []rune(name)
	length := len(r)
	//1.长度是否合法
	if length < minLen || length > maxLen {
		return false
	}
	//if ll := len(name); ll < minLen || ll > maxLen { //使用字节数
	//	return false
	//}
	//2.是否包含非法字符
	for i, c := range r {
		if n := utf8.RuneLen(c); n > 3 || n == -1 {
			return false
		}
		switch {
		case c >= 'a' && c <= 'z':
		case c >= 'A' && c <= 'Z':
		case c >= '0' && c <= '9':
		case c == ' ':
			if i == 0 || i == (length-1) { //第一个、最后一个字符不能为空格
				return false
			}
		case unicode.IsOneOf(validChars, c):
		default:
			return false
		}
	}
	return true
}

func CheckVaildAbbr(abbr string, size int) bool {
	r := []rune(abbr)
	length := len(r)
	//1.长度是否合法
	if length != size {
		return false
	}
	//2.是否包含非法字符
	allSpace := true
	for _, c := range r {
		switch {
		case c == ' ':
		case c > ' ' && c <= '~': //abbr支持所有可见字符 《! " # $ % & ' ( ) * + , - . / 0 1 2 3 4 5 6 7 8 9 : ; < = > ? @ A B C D E F G H I J K L M N O P Q R S T U V W X Y Z [ \ ] ^ _ ` a b c d e f g h i j k l m n o p q r s t u v w x y z { | } ~》
			allSpace = false
		default:
			return false
		}
	}
	return !allSpace
}

func CheckStrLen(str string, minLen, maxLen int) bool {
	length := len([]rune(str))
	//1.长度是否合法
	return length >= minLen && length <= maxLen
}

const (
	letterBytes   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

func RandStr(n int) string {
	b := make([]byte, n)
	// A rd.Int63() generates 63 random bits, enough for letterIdxMax letters!
	for i, cache, remain := n-1, rand.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = rand.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}
	return string(b)
}

// no copy to change slice to string
func String(b []byte) (s string) {
	pbytes := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	pstring := (*reflect.StringHeader)(unsafe.Pointer(&s))
	pstring.Data = pbytes.Data
	pstring.Len = pbytes.Len
	return
}

// no copy to change string to slice
func Slice(s string) (b []byte) {
	pbytes := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	pstring := (*reflect.StringHeader)(unsafe.Pointer(&s))
	pbytes.Data = pstring.Data
	pbytes.Len = pstring.Len
	pbytes.Cap = pstring.Len
	return
}
