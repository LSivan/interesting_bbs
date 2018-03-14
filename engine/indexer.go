package engine

import (
	"git.oschina.net/gdou-geek-bbs/models"
	"git.oschina.net/gdou-geek-bbs/utils"
	"github.com/astaxie/beego"
	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/analysis/analyzer/standard"
	"github.com/blevesearch/bleve/mapping"
	"log"
	"strconv"
	"sync"
	"time"
)

type indexer struct {
	NextDocId    int64 // 下一次更新时的DocId
	indexPath    string
	i            bleve.Index
	indexMapping mapping.IndexMapping
	batchSize    int
	InsertChan   chan *models.Topic
	UpdateChan   chan *models.Topic
	DeleteChan   chan *models.Topic
	Exit         chan struct{}
	Fin          chan struct{}
	lock         chan struct{}
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
			NextDocId:    beego.AppConfig.DefaultInt64("engine.next.doc.id", 0),
			indexPath:    beego.AppConfig.String("engine.index.path"),
			indexMapping: indexMapping,
			batchSize:    beego.AppConfig.DefaultInt("engine.batch.size", 10),
			InsertChan:   make(chan *models.Topic, 500),
			UpdateChan:   make(chan *models.Topic, 50),
			DeleteChan:   make(chan *models.Topic, 50),
			lock:         make(chan struct{}, 1),
			Exit:         make(chan struct{}, 1),
			Fin:          make(chan struct{}, 0),
		}
	})
}

// 负责索引的建立，更新
func (self *indexer) Index() {

	topicIndex, err := bleve.Open(self.indexPath)
	var needRecord bool
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
	beego.BeeLogger.Info("共需要建立%d条增量索引\n", count)

	if count > 0 { // 若是count>0，则证明需要添加增量索引
		beego.BeeLogger.Info("从ID为%d开始建立增量索引\n", nextDocId)
		// 添加增量索引之后，把配置的engine.next.doc.id记录下来，留待下一次初始化的时候启用
		topics := models.FindTopicFrom(int(nextDocId), 0)
		go func() {
			err = self.indexTopic(needRecord, topics)
			utils.LogError("建立话题索引", err)
		}()
	}
	batch := self.i.NewBatch()
	batchCount := 0
	go func() {
		beego.BeeLogger.Info("监听文章变化中....\n")
		select {
		case topic := <-self.InsertChan:
			beego.BeeLogger.Debug("存在新的文章，需要建立新索引\n")
			batch.Index(strconv.Itoa(topic.Id), topic)
			self.lock <- struct{}{}
			batchCount++
			if batchCount >= self.batchSize {
				beego.AppConfig.Set("engine.next.doc.id", strconv.Itoa(topic.Id))
				err := self.i.Batch(batch)
				utils.LogError("批量添加新索引", err)
				batch = self.i.NewBatch()
				batchCount = 0
			}
			<-self.lock
		case topic := <-self.UpdateChan:
			beego.BeeLogger.Debug("存在更新的文章，需要改变索引\n")
			batch.Index(strconv.Itoa(topic.Id), topic)
		case topic := <-self.DeleteChan:
			beego.BeeLogger.Debug("存在删除的文章，需要删除索引\n")
			err := self.i.Delete(strconv.Itoa(topic.Id))
			utils.LogError("从索引中移除文档", err)
		case <-self.Exit:
			// 保存配置，将那些batch推送到索引文件
			self.i.Batch(batch)
			beego.AppConfig.SaveConfigFile("conf/app.conf")
			beego.BeeLogger.Warning("捕获到退出信号，保存一些必要的东西\n")
			self.Fin <- struct{}{}
		}
	}()
}

func buildIndexMapping() (mapping.IndexMapping, error) {
	numericFieldMapping := bleve.NewNumericFieldMapping()
	numericFieldMapping.Analyzer = standard.Name
	textFieldMapping := bleve.NewTextFieldMapping()
	textFieldMapping.Analyzer = standard.Name

	topicMapping := bleve.NewDocumentMapping()
	topicMapping.AddFieldMappingsAt("Title", textFieldMapping)
	contentFieldMapping := bleve.NewTextFieldMapping()
	contentFieldMapping.Store = false
	topicMapping.AddFieldMappingsAt("Content", contentFieldMapping)
	topicMapping.AddFieldMappingsAt("View", numericFieldMapping)
	topicMapping.AddFieldMappingsAt("ReplyCount", numericFieldMapping)
	topicMapping.AddFieldMappingsAt("CollectCount", numericFieldMapping)
	//topicMapping.AddFieldMappingsAt("LastReplyTime", dataTimeFieldMapping)

	usernameFieldMapping := bleve.NewTextFieldMapping()

	usernameFieldMapping.Store = false
	beego.BeeLogger.Debug("usernameFieldMapping.Options().IsStored():%v\n;",usernameFieldMapping.Options())
	beego.BeeLogger.Debug("	o&StoreField != 0:%v\n;",usernameFieldMapping.Options()&2)
	topicMapping.AddFieldMappingsAt("Username", usernameFieldMapping)

	topicMapping.AddFieldMappingsAt("SectionName", textFieldMapping)

	topicMapping.AddFieldMappingsAt("LastReplyUsername", textFieldMapping)

	indexMapping := bleve.NewIndexMapping()

	indexMapping.AddDocumentMapping("topic", topicMapping)
	//indexMapping.TypeField = "type"
	indexMapping.DefaultAnalyzer = standard.Name

	return indexMapping, nil
}

func (self *indexer) indexTopic(needRecordCount bool, topics []*models.Topic) error {
	beego.BeeLogger.Info("Indexing...\n")
	count := 0
	startTime := time.Now()
	//batch := self.i.NewBatch()
	//batchCount := 0
	for _, topic := range topics {
		self.i.Index(strconv.Itoa(topic.Id), transformIndexTopic(topic))
		time.Sleep(1*time.Second)
		//batch.Index(strconv.Itoa(topic.Id), transformIndexTopic(topic))
		//batchCount++
		//
		//if batchCount >= self.batchSize {
		//	err := self.i.Batch(batch)
		//	if err != nil {
		//		return err
		//	}
		//	batch = self.i.NewBatch()
		//	batchCount = 0
		//}
		count++
		if count%1000 == 0 {
			indexDuration := time.Since(startTime)
			indexDurationSeconds := float64(indexDuration) / float64(time.Second)
			timePerDoc := float64(indexDuration) / float64(count)
			beego.BeeLogger.Info("Indexed %d documents, in %.2fs (average %.2fms/doc)\n", count, indexDurationSeconds, timePerDoc/float64(time.Millisecond))
		}
	}

	// flush the last batch
	//if batchCount > 0 {
	//	err := self.i.Batch(batch)
	//	utils.LogError("批量索引化", err)
	//	batchCount = 0
	//}
	if i := len(topics) - 1; i >= 0 {
		// 将检索到的最后一个文章的id的值记录下来
		nextDocId := strconv.Itoa(topics[i].Id)
		beego.BeeLogger.Info("记录到配置文件的nextDocId为%s\n", nextDocId)
		err := beego.AppConfig.Set("engine.next.doc.id", nextDocId)
		utils.LogError("记录nextDocId", err)
		if err := beego.AppConfig.SaveConfigFile("conf/app.conf"); err != nil {
			utils.LogError("保存配置文件", err)
		}
	}

	indexDuration := time.Since(startTime)
	indexDurationSeconds := float64(indexDuration) / float64(time.Second)
	timePerDoc := float64(indexDuration) / float64(count)
	beego.BeeLogger.Info("Indexed %d documents, in %.2fs (average %.2fms/doc)", count, indexDurationSeconds, timePerDoc/float64(time.Millisecond))

	return nil
}

type indexTopic struct {
	Title             string
	Content           string
	Username          string
	SectionName       string
	LastReplyUsername string
	View              int
	ReplyCount        int
	CollectCount      int
}

func transformIndexTopic(topic *models.Topic) *indexTopic {
	i := &indexTopic{
		Title:        topic.Title,
		Content:      topic.Content,
		Username:     topic.User.Username,
		SectionName:  topic.Section.Name,
		View:         topic.View,
		ReplyCount:   topic.ReplyCount,
		CollectCount: topic.CollectCount,
	}
	if topic.LastReplyUser != nil {
		i.LastReplyUsername = topic.LastReplyUser.Username
	}
	return i
}
