package sego

import (
	"bytes"
	"fmt"
	"strings"
	"math"
)

// 输出分词结果为字符串
//
// 有两种输出模式，以"中华人民共和国"为例
//
//  普通模式（searchMode=false）输出一个分词"中华人民共和国/ns "
//  搜索模式（searchMode=true） 输出普通模式的再细致切分：
//      "中华/nz 人民/n 共和/nz 共和国/ns 人民共和国/nt 中华人民共和国/ns "
//
// 搜索模式主要用于给搜索引擎提供尽可能多的关键字，详情请见Token结构体的注释。
func SegmentsToString(segs []Segment, searchMode bool) (output string) {
	if searchMode {
		for _, seg := range segs {
			output += tokenToString(seg.token)
		}
	} else {
		for _, seg := range segs {
			output += fmt.Sprintf(
				"%s/%s ", textSliceToString(seg.token.text), seg.token.pos)
		}
	}
	return
}

func tokenToString(token *Token) (output string) {
	hasOnlyTerminalToken := true
	for _, s := range token.segments {
		if len(s.token.segments) > 1 {
			hasOnlyTerminalToken = false
		}
	}

	if !hasOnlyTerminalToken {
		for _, s := range token.segments {
			if s != nil {
				output += tokenToString(s.token)
			}
		}
	}
	output += fmt.Sprintf("%s/%s ", textSliceToString(token.text), token.pos)
	return
}

// 输出分词结果到一个字符串slice
//
// 有两种输出模式，以"中华人民共和国"为例
//
//  普通模式（searchMode=false）输出一个分词"[中华人民共和国]"
//  搜索模式（searchMode=true） 输出普通模式的再细致切分：
//      "[中华 人民 共和 共和国 人民共和国 中华人民共和国]"
//
// 搜索模式主要用于给搜索引擎提供尽可能多的关键字，详情请见Token结构体的注释。

func SegmentsToSlice(segs []Segment, searchMode bool) (output []string) {
	if searchMode {
		for _, seg := range segs {
			output = append(output, tokenToSlice(seg.token)...)
		}
	} else {
		for _, seg := range segs {
			output = append(output, seg.token.Text())
		}
	}
	return
}

func SegmentsToFreqMap(segs []Segment) (output map[string]int) {
	output = make(map[string]int, 20)
	for _, seg := range segs {
		word := seg.token.Text()
		if i, ok := output[word];ok {
			output[word] = i + 1
		} else {
			output[word] = 1
		}
	}
	return
}
/**
salt 影响特征值的一个数
 */
func SegmentsToFeatureSlice(segs []Segment, salt float64) (feature []float64) {
	output := make([]int,len(KeywordMap))
	feature = make([]float64,len(KeywordMap))
	total := 0
	for _, seg := range segs {
		word := seg.token.Text()
		if index,exist := KeywordMap[strings.ToLower(word)];exist {
			i := output[index]
			output[index] = i + 1
			total++
		}
	}
	var log10 = func(num float64) float64 {
		num = math.Log10(num+10+salt)
		// 防止恶意刷关键字
		for num > 2.33 {
			num = math.Log10(num)
		}
		return num - 1
	}
	for index, i := range output {
		feature[index] = log10(float64(i))+float64(i)/float64(total)

	}
	return
}


func tokenToSlice(token *Token) (output []string) {
	hasOnlyTerminalToken := true
	for _, s := range token.segments {
		if len(s.token.segments) > 1 {
			hasOnlyTerminalToken = false
		}
	}
	if !hasOnlyTerminalToken {
		for _, s := range token.segments {
			output = append(output, tokenToSlice(s.token)...)
		}
	}
	output = append(output, textSliceToString(token.text))
	return output
}

// 将多个字元拼接一个字符串输出
func textSliceToString(text []Text) string {
	var output string
	for _, word := range text {
		output += string(word)
	}
	return output
}

// 返回多个字元的字节总长度
func textSliceByteLength(text []Text) (length int) {
	for _, word := range text {
		length += len(word)
	}
	return
}

func textSliceToBytes(text []Text) []byte {
	var buf bytes.Buffer
	for _, word := range text {
		buf.Write(word)
	}
	return buf.Bytes()
}
