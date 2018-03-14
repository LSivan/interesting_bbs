package engine

import (
	"git.oschina.net/gdou-geek-bbs/utils"
	"github.com/astaxie/beego"
	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/analysis/analyzer/standard"
	"github.com/blevesearch/bleve/search"
	"github.com/blevesearch/bleve/search/query"
	"github.com/huichen/sego"
	"strings"
	"sync"
)

// 负责检索

type Searcher struct {
}

var searchRequest *bleve.SearchRequest
var pageSize int
var segmenter sego.Segmenter

func init() {
	once := &sync.Once{}

	once.Do(func() {
		pageSize = beego.AppConfig.DefaultInt("page.size", 8)
		searchRequest = bleve.NewSearchRequestOptions(
			nil,
			pageSize,
			0,
			beego.AppConfig.DefaultBool("engine.is.explain", true),
		)
		searchRequest.Highlight = bleve.NewHighlightWithStyle("html")

		segmenter.LoadDictionary("engine/dictionary.txt")

	})
}

/**
 * q 查询的关键字
 * p 页码
 */
func (searcher *Searcher) Search(q string, p int) (*bleve.SearchResult, error) {
	mpQuery := bleve.NewMatchPhraseQuery(q)
	searchRequest.From = pageSize * (p - 1)

	mpQuery.Analyzer = standard.Name
	searchRequest.Sort = search.SortOrder{
		&search.SortScore{
			Desc: true,
		},
		&search.SortField{
			Field: "CollectCount",
			Desc:  true,
		},
		&search.SortField{
			Field: "ReplyCount",
			Desc:  true,
		},
		&search.SortField{
			Field: "View",
			Desc:  true,
		},
	}
	searchRequest.Query = mpQuery

	searchResult, err := Indexer.i.Search(searchRequest)
	utils.LogError("从索引查找", err)
	if len(searchResult.Hits) == 0 { // 假设在不切割的情况下就已经能搜索到了，直接将结果返回
		// 不切割的时候搜索不到，则进行切割
		terms := sego.SegmentsToSlice(segmenter.Segment([]byte(strings.TrimSpace(q))), true)
		//terms := strings.Split(q, " ")
		queries := make([]query.Query, 0, len(terms))
		for _, term := range terms {
			mpQuery := bleve.NewMatchPhraseQuery(term)
			mpQuery.Analyzer = standard.Name
			queries = append(queries, mpQuery)
		}
		conjunctionQuery := bleve.NewDisjunctionQuery(queries...)
		searchRequest.Query = conjunctionQuery
		searchResult, err = Indexer.i.Search(searchRequest)
		utils.LogError("从索引查找", err)
		if err != nil {
			return nil, err
		}
	}
	for _, v := range searchResult.Hits {
		doc, err := Indexer.i.Document(v.ID)
		utils.LogError("获取索引的文档", err)

		if doc != nil {
			for _, v := range doc.Fields {
				beego.BeeLogger.Debug("name :%s, value : %s\n", v.Name(), string(v.Value()))
			}
		}
	}
	return searchResult, nil
}
