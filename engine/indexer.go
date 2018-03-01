package engine

import (
	"git.oschina.net/gdou-geek-bbs/models"
	"git.oschina.net/gdou-geek-bbs/utils"
	"github.com/astaxie/beego"
	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/analysis/analyzer/standard"
	"github.com/blevesearch/bleve/analysis/lang/en"
	"github.com/blevesearch/bleve/mapping"
	"log"
	"strconv"
	"sync"
	"time"
)

type indexer struct {
	NextDocId     int64 // 下一次更新时的DocId
	indexPath     string
	i             bleve.Index
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
	utils.LogError("打开存储的索引文件", err)
	if err == bleve.ErrorIndexPathDoesNotExist { // 存储文件不存在
		// 1. 记录下一次需要index的id
		topicIndex, err = bleve.New(self.indexPath, self.indexMapping)
		utils.LogError("新建索引文件", err)
	} else {
			// 1.记录下一次index的文档的ID
			nextDocId = self.NextDocId
			needRecord = true

	}
	self.i = topicIndex
	// 先通过count寻找是否需要添加增量索引
	count := models.CountTopicFromID(int(nextDocId))
	log.Printf("共需要建立%d条增量索引\n", count)

	if count > 0 {// 若是count>0，则证明需要添加增量索引
		log.Printf("从ID为%d开始建立增量索引\n", nextDocId)
		// 添加增量索引之后，把配置的engine.next.doc.id记录下来，留待下一次初始化的时候启用
		topics := models.FindTopicFrom(int(nextDocId), 0)
		go func() {
			err = self.indexTopic(extraDocCount, needRecord, topics)
			utils.LogError("建立话题索引", err)
		}()
	}

}

func buildIndexMapping() (mapping.IndexMapping, error) {

	numericFieldMapping := bleve.NewNumericFieldMapping()
	numericFieldMapping.Analyzer = standard.Name
	textFieldMapping := bleve.NewTextFieldMapping()
	textFieldMapping.Analyzer = standard.Name
	dataTimeFieldMapping := bleve.NewDateTimeFieldMapping()
	dataTimeFieldMapping.Analyzer = standard.Name

	topicMapping := bleve.NewDocumentMapping()
	topicMapping.AddFieldMappingsAt("Id", numericFieldMapping)
	topicMapping.AddFieldMappingsAt("Title", textFieldMapping)
	topicMapping.AddFieldMappingsAt("Content", textFieldMapping)
	topicMapping.AddFieldMappingsAt("InTime", dataTimeFieldMapping)
	topicMapping.AddFieldMappingsAt("View", numericFieldMapping)
	topicMapping.AddFieldMappingsAt("ReplyCount", numericFieldMapping)
	topicMapping.AddFieldMappingsAt("CollectCount", numericFieldMapping)
	//topicMapping.AddFieldMappingsAt("LastReplyTime", dataTimeFieldMapping)

	userMapping := bleve.NewDocumentMapping()
	userMapping.AddFieldMappingsAt("User.Id", numericFieldMapping)
	userMapping.AddFieldMappingsAt("User.Username", textFieldMapping)
	//userMapping.AddFieldMappingsAt("User.Password", textFieldMapping)
	//userMapping.AddFieldMappingsAt("User.Token", textFieldMapping)
	//userMapping.AddFieldMappingsAt("User.Avatar", textFieldMapping)
	//userMapping.AddFieldMappingsAt("User.Email", textFieldMapping)
	//userMapping.AddFieldMappingsAt("User.Url", textFieldMapping)
	//userMapping.AddFieldMappingsAt("User.Signature", textFieldMapping)
	//userMapping.AddFieldMappingsAt("User.InTime", dataTimeFieldMapping)
	topicMapping.AddSubDocumentMapping("User", userMapping)

	sectionMapping := bleve.NewDocumentMapping()
	//sectionMapping.AddFieldMappingsAt("Section.Id", textFieldMapping)
	sectionMapping.AddFieldMappingsAt("Section.Name", textFieldMapping)
	topicMapping.AddSubDocumentMapping("Section", sectionMapping)

	lastReplyUserMapping := bleve.NewDocumentMapping()
	//topicMapping.AddFieldMappingsAt("LastReplyUser.Id", numericFieldMapping)
	topicMapping.AddFieldMappingsAt("LastReplyUser.Username", textFieldMapping)
	//topicMapping.AddFieldMappingsAt("LastReplyUser.Password", textFieldMapping)
	//topicMapping.AddFieldMappingsAt("LastReplyUser.Token", textFieldMapping)
	//topicMapping.AddFieldMappingsAt("LastReplyUser.Avatar", textFieldMapping)
	//topicMapping.AddFieldMappingsAt("LastReplyUser.Email", textFieldMapping)
	//topicMapping.AddFieldMappingsAt("LastReplyUser.Url", textFieldMapping)
	//topicMapping.AddFieldMappingsAt("LastReplyUser.Signature", textFieldMapping)
	//topicMapping.AddFieldMappingsAt("LastReplyUser.InTime", dataTimeFieldMapping)
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

	// flush the last batch
	if batchCount > 0 {
		err := self.i.Batch(batch)
		utils.LogError("批量索引化", err)
	}
	if i := len(topics) - 1; i >= 0 {
			// 将检索到的最后一个文章的id的值记录下来
			nextDocId := strconv.Itoa(topics[i].Id)
			log.Printf("记录到配置文件的nextDocId为%s\n", nextDocId)
			err := beego.AppConfig.Set("engine.next.doc.id", nextDocId)
			utils.LogError("记录nextDocId", err)
		if err := beego.AppConfig.SaveConfigFile("conf/app.conf"); err != nil {
			utils.LogError("保存配置文件", err)
		}
	}

	indexDuration := time.Since(startTime)
	indexDurationSeconds := float64(indexDuration) / float64(time.Second)
	timePerDoc := float64(indexDuration) / float64(count)
	log.Printf("Indexed %d documents, in %.2fs (average %.2fms/doc)", count, indexDurationSeconds, timePerDoc/float64(time.Millisecond))
	return nil
}
