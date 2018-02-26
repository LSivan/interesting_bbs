package engine

import (
	"github.com/blevesearch/bleve"
	"log"
)

// 负责检索

type Searcher struct {
}

func (searcher *Searcher)Search(keyword string) (*bleve.SearchResult, error) {
	log.Println("keyword : ",keyword)
	termQuery := bleve.NewMatchQuery(keyword)
	//termQuery.SetField("Title")
	//termQuery.SetField("name")

	termSearchRequest := bleve.NewSearchRequest(termQuery)
	termSearchRequest.Explain = true
	termSearchRequest.Highlight = bleve.NewHighlight()
	termSearchResult, err := Indexer.i.Search(termSearchRequest)
	if err != nil {
		return nil, err
	}
	return termSearchResult,nil
	//return &termSearchResult.Hits, nil
}
