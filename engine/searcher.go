package engine

import (
	"encoding/gob"
	"github.com/astaxie/beego"
	"github.com/huichen/wukong/engine"
	"github.com/huichen/wukong/types"
	"github.com/vincent3i/beego-blog/models"
	"github.com/vincent3i/beego-blog/models/blog"
	"path/filepath"
	"reflect"
	"strings"
)

const (
	SecondsInADay     = 86400
	MaxTokenProximity = 2
)

var (
	searcher = engine.Engine{}
)

/*******************************************************************************
索引
*******************************************************************************/
func indexBlog() {
	var bbs []models.Blog
	blog.Blogs().Limit(-1).All(&bbs)
	beego.Debug("针对当前所有博客进行索引....")
	for _, bb := range bbs {
		labels := strings.Split(bb.Keywords, ",")
		beego.BeeLogger.Debug("Get [%d] key words from blog", len(labels))
		searcher.IndexDocument(uint64(bb.Id), types.DocumentIndexData{
			Content: bb.Title + blog.ReadBlogContent(&bb).Content,
			Labels:  labels,
			Fields: BlogScoringFields{
				BlogLastUpdate: uint64(bb.BlogContentLastUpdate),
				BlogViews:      uint64(bb.Views),
			},
		})
	}

	searcher.FlushIndex()
	beego.BeeLogger.Debug("索引了%d条", len(bbs))
}

/*******************************************************************************
评分
*******************************************************************************/
type BlogScoringFields struct {
	BlogLastUpdate uint64
	BlogViews      uint64
}

type BlogScoringCriteria struct {
}

func (criteria BlogScoringCriteria) Score(doc types.IndexedDocument, fields interface{}) []float32 {
	if reflect.TypeOf(fields) != reflect.TypeOf(BlogScoringFields{}) {
		beego.Warn("当前排序类型不匹配，匹配的类型为", reflect.TypeOf(BlogScoringFields{}))
		return []float32{}
	}

	wsf := fields.(BlogScoringFields)
	output := make([]float32, 3)
	if doc.TokenProximity > MaxTokenProximity {
		output[0] = 1.0 / float32(doc.TokenProximity)
	} else {
		output[0] = 1.0
	}
	output[1] = float32(wsf.BlogLastUpdate / (SecondsInADay * 3))
	output[2] = float32(doc.BM25 * (1 + float32(wsf.BlogViews)/1000))
	return output
}

/*******************************************************************************
search
*******************************************************************************/
func SearchResult(query string, isAdmin bool) []models.Blog {
	output := searcher.Search(types.SearchRequest{
		Text: query,
		RankOptions: &types.RankOptions{
			ScoringCriteria: &BlogScoringCriteria{},
			OutputOffset:    0,
			MaxOutputs:      100,
		},
	})

	// 整理为输出格式
	docs := []models.Blog{}
	for _, doc := range output.Docs {
		bb_id := int64(doc.DocId)
		if !isAdmin && blog.IsAuth(bb_id) {
			beego.Debug("You can not read this blog since you don't have permission.")
			continue
		}
		wb := blog.OneById(bb_id)
		//		wb.Content = blog.ReadBlogContent(wb)
		//		for _, t := range output.Tokens {
		//			beego.Debug("token--->>", t)
		//			wb.SnapShot = strings.Replace(wb.SnapShot, t, "<font color=red>"+t+"</font>", -1)
		//		}
		docs = append(docs, *wb)
	}

	return docs
}

/*******************************************************************************
主函数
*******************************************************************************/
func InitSearcher() error {
	// 初始化
	gob.Register(BlogScoringFields{})

	segmDictFilePath := filepath.Join(beego.AppPath, "data", "dictionary.txt")
	beego.Debug("SegmenterDictionaries path ---->>> ", segmDictFilePath)
	stopTokenFilePath := filepath.Join(beego.AppPath, "data", "stop_tokens.txt")
	beego.Debug("StopTokenFile path ---->>> ", stopTokenFilePath)
	searcher.Init(types.EngineInitOptions{
		SegmenterDictionaries: segmDictFilePath,
		StopTokenFile:         stopTokenFilePath,
		IndexerInitOptions: &types.IndexerInitOptions{
			IndexType: types.LocationsIndex,
		},
		UsePersistentStorage:    true,
		PersistentStorageFolder: "db",
	})

	// 索引
	go indexBlog()

	return nil
}

func CloseSearcher() {
	searcher.Close()
	ElasticClient.Stop()
	beego.Debug("Closed searcher......")
}
