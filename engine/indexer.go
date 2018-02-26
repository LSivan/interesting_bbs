package engine

import (
	"fmt"
	"git.oschina.net/gdou-geek-bbs/models"
	"github.com/astaxie/beego"
	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/analysis/analyzer/keyword"
	"github.com/blevesearch/bleve/mapping"
	"log"
	"os"
	"strconv"
	"sync"
	"time"
	"github.com/blevesearch/bleve/analysis/lang/en"
)

type indexer struct {
	NextDocId     int64 // 下一次更新时的DocId
	indexPath     string
	i             bleve.Index
	extraDocCount int64 // 已建立好索引的存储之外的文档数量，超过1000就重新建立存储
	indexMapping  mapping.IndexMapping
	batchSize     int
}

var Indexer *indexer

func init() {
	once := sync.Once{}
	once.Do(func() {
		indexMapping, err := buildIndexMapping()
		if err != nil {
			log.Fatal(err)
		}

		Indexer = &indexer{
			NextDocId:     beego.AppConfig.DefaultInt64("engine.next.doc.id", 0),
			indexPath:     beego.AppConfig.String("engine.index.path"),
			indexMapping:  indexMapping,
			extraDocCount: beego.AppConfig.DefaultInt64("engine.extra.doc.count", 0),
			batchSize:     beego.AppConfig.DefaultInt("engine.batch.size", 100),
		}
	})
}

// 第一种，非持久化存储，每次重启应用直接Index
// 第二种，持久化存储，重启应用只需要加载存储的文件，有新的文档需要建立索引的时候直接往内存中index，
// 			同时记录需要index操作的文档数，若是当index的文档数达到某个阈值的时候（index的时间较长，还不如重新建立索引文件，
// 			一次性加载来得快），便重新建立索引文件
//
// 负责索引的建立，更新
func (self *indexer) Index() {
	topicIndex, err := bleve.Open(self.indexPath)
	var needRecord bool
	var extraDocCount int64
	var nextDocId int64
	var flag bool
	log.Println("err : ", err)
	if err == bleve.ErrorIndexPathDoesNotExist { // 存储文件不存在
		// 1. 记录下一次需要index的id
		topicIndex, err = bleve.New(self.indexPath, self.indexMapping)
		if err != nil {
			log.Fatal(err)
		}
		flag = true
	} else {
		log.Println("self.extraDocCount", self.extraDocCount)
		if self.extraDocCount >= 1000 { // 如果需要额外加入索引的文档数量达到了阈值，就重新建立索引存储
			// 1.记录下一次index的ID
			// 2.将需要额外index的文档数量置为0
			os.RemoveAll(self.indexPath)
			topicIndex, err = bleve.New(self.indexPath, self.indexMapping)
			if err != nil {
				log.Fatal(err)
			}
			flag = true
		} else if self.extraDocCount <= 0 { // 直接读取完存储的索引文件即可
			flag = false
		} else { // 直接追加
			// 1.记录下一次index的文档的ID
			// 2.记录需要append的文档的数量
			nextDocId = self.NextDocId
			extraDocCount = self.extraDocCount
			needRecord = true
			flag = true
		}

	}
	self.i = topicIndex
	if flag {
		topics := models.FindTopicFrom(int(nextDocId), 0)
		go func() {
			err = self.indexTopic(extraDocCount, needRecord, topics)
			if err != nil {
				log.Fatal(err)
			}
		}()
	}
	doc, err := self.i.Document("1")
	if err != nil {
		log.Println("get doc err ", err)
	}
	if doc != nil {
		for _, v := range doc.Fields {
			//s, _ := json.Marshal(v)
			fmt.Printf("topicMapping.AddFieldMappingsAt(\"%s\", textFieldMapping)\n", v.Name())
		}
	}
}

func buildIndexMapping() (mapping.IndexMapping, error) {

	numericFieldMapping := bleve.NewNumericFieldMapping()
	numericFieldMapping.Analyzer = en.AnalyzerName
	textFieldMapping := bleve.NewTextFieldMapping()
	textFieldMapping.Analyzer = en.AnalyzerName
	dataTimeFieldMapping := bleve.NewDateTimeFieldMapping()
	dataTimeFieldMapping.Analyzer = keyword.Name

	topicMapping := bleve.NewDocumentMapping()
	topicMapping.AddFieldMappingsAt("Id", numericFieldMapping)
	topicMapping.AddFieldMappingsAt("Title", textFieldMapping)
	topicMapping.AddFieldMappingsAt("Content", textFieldMapping)
	topicMapping.AddFieldMappingsAt("InTime", dataTimeFieldMapping)
	topicMapping.AddFieldMappingsAt("View", numericFieldMapping)
	topicMapping.AddFieldMappingsAt("ReplyCount", numericFieldMapping)
	topicMapping.AddFieldMappingsAt("CollectCount", numericFieldMapping)
	topicMapping.AddFieldMappingsAt("LastReplyTime", dataTimeFieldMapping)

	userMapping := bleve.NewDocumentMapping()
	userMapping.AddFieldMappingsAt("User.Id", numericFieldMapping)
	userMapping.AddFieldMappingsAt("User.Username", textFieldMapping)
	userMapping.AddFieldMappingsAt("User.Password", textFieldMapping)
	userMapping.AddFieldMappingsAt("User.Token", textFieldMapping)
	userMapping.AddFieldMappingsAt("User.Avatar", textFieldMapping)
	userMapping.AddFieldMappingsAt("User.Email", textFieldMapping)
	userMapping.AddFieldMappingsAt("User.Url", textFieldMapping)
	userMapping.AddFieldMappingsAt("User.Signature", textFieldMapping)
	userMapping.AddFieldMappingsAt("User.InTime", dataTimeFieldMapping)
	topicMapping.AddSubDocumentMapping("User", userMapping)

	sectionMapping := bleve.NewDocumentMapping()
	sectionMapping.AddFieldMappingsAt("Section.Id", textFieldMapping)
	sectionMapping.AddFieldMappingsAt("Section.Name", textFieldMapping)
	topicMapping.AddSubDocumentMapping("Section", sectionMapping)

	lastReplyUserMapping := bleve.NewDocumentMapping()
	topicMapping.AddFieldMappingsAt("LastReplyUser.Id", numericFieldMapping)
	topicMapping.AddFieldMappingsAt("LastReplyUser.Username", textFieldMapping)
	topicMapping.AddFieldMappingsAt("LastReplyUser.Password", textFieldMapping)
	topicMapping.AddFieldMappingsAt("LastReplyUser.Token", textFieldMapping)
	topicMapping.AddFieldMappingsAt("LastReplyUser.Avatar", textFieldMapping)
	topicMapping.AddFieldMappingsAt("LastReplyUser.Email", textFieldMapping)
	topicMapping.AddFieldMappingsAt("LastReplyUser.Url", textFieldMapping)
	topicMapping.AddFieldMappingsAt("LastReplyUser.Signature", textFieldMapping)
	topicMapping.AddFieldMappingsAt("LastReplyUser.InTime", dataTimeFieldMapping)
	topicMapping.AddSubDocumentMapping("LastReplyUser", lastReplyUserMapping)

	indexMapping := bleve.NewIndexMapping()
	indexMapping.AddDocumentMapping("topic", topicMapping)
	indexMapping.TypeField = "type"
	indexMapping.DefaultAnalyzer = en.AnalyzerName

	return indexMapping, nil
}

func (self *indexer) indexTopic(extraDocCount int64, needRecordCount bool, topics []*models.Topic) error {

	log.Printf("Indexing...")
	count := 0
	startTime := time.Now()
	batch := self.i.NewBatch()
	batchCount := 0
	for _, topic := range topics {
		batch.Index(strconv.Itoa(topic.Id), topic)
		extraDocCount++
		batchCount++

		if batchCount >= self.batchSize {
			err := self.i.Batch(batch)
			if err != nil {
				return err
			}
			batch = self.i.NewBatch()
			batchCount = 0
		}
		count++
		if count%1000 == 0 {
			indexDuration := time.Since(startTime)
			indexDurationSeconds := float64(indexDuration) / float64(time.Second)
			timePerDoc := float64(indexDuration) / float64(count)
			log.Printf("Indexed %d documents, in %.2fs (average %.2fms/doc)", count, indexDurationSeconds, timePerDoc/float64(time.Millisecond))
		}
	}
	if needRecordCount {
		beego.AppConfig.Set("engine.extra.doc.count", strconv.Itoa(int(extraDocCount)))
	}
	// flush the last batch
	if batchCount > 0 {
		err := self.i.Batch(batch)
		if err != nil {
			log.Fatal(err)
		}
	}
	indexDuration := time.Since(startTime)
	indexDurationSeconds := float64(indexDuration) / float64(time.Second)
	timePerDoc := float64(indexDuration) / float64(count)
	log.Printf("Indexed %d documents, in %.2fs (average %.2fms/doc)", count, indexDurationSeconds, timePerDoc/float64(time.Millisecond))
	return nil
}
