package engine

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/olivere/elastic"
	"github.com/vincent3i/beego-blog/models"
	"github.com/vincent3i/beego-blog/models/blog"
	"github.com/vincent3i/beego-blog/nsq/producer"
	"strconv"
)

const Blog_Index_Name = "ik_g_blogs"

var (
	ElasticClient *elastic.Client
)

type ElasticBlog struct {
	Id      string `json:"id"`
	Content string `json:"content"`
}

func InitElasticSearch() error {
	var err error
	// Obtain a client and connect to the default Elasticsearch installation
	// on 127.0.0.1:9200. Of course you can configure your client to connect
	// to other hosts and configure it in various other ways.
	ElasticClient, err = elastic.NewClient(
		elastic.SetURL("http://localhost:9200"),
	)
	if err != nil {
		return err
	}

	// Ping the Elasticsearch server to get e.g. the version number
	info, code, err := ElasticClient.Ping().Do()
	if err != nil {
		// Handle error
		return err
	}
	beego.BeeLogger.Debug("Elasticsearch returned with code [%d] and version [%s]", code, info.Version.Number)

	// Use the IndexExists service to check if a specified index exists.
	exists, err := ElasticClient.IndexExists(Blog_Index_Name).Do()
	if err != nil {
		// Handle error
		return err
	}
	if !exists {
		// Create a new index.
		createIndex, err := ElasticClient.CreateIndex(Blog_Index_Name).Do()
		if err != nil {
			// Handle error
			return err
		}
		if !createIndex.Acknowledged {
			// Not acknowledged
			beego.BeeLogger.Debug("Create index %s failed!", Blog_Index_Name)
		}
	}

	// Index (using JSON serialization)
	var bbs []models.Blog
	blog.Blogs().Limit(-1).All(&bbs)
	beego.Debug("Elastic search 针对当前所有博客进行索引....")
	vv := make([]interface{}, len(bbs))

	for i, bb := range bbs {
		vv[i] = bb
	}
	producer.PublishMsg("elastic-blog", vv...)

	// Flush to make sure the documents got written.
	_, err = ElasticClient.Flush().Index(Blog_Index_Name).Do()
	if err != nil {
		return err
	}

	return nil
}

func ElasticSearch(query string, isAdmin bool) []models.Blog {
	// Search with a term query
	docs := []models.Blog{}

	termQuery := elastic.NewTermQuery("content", query)
	searchResult, err := ElasticClient.Search().
		Index(Blog_Index_Name). // search in index "twitter"
		Query(&termQuery).      // specify the query
		Sort("content", false). // sort by "user" field, ascending
		//From(0).Size(10).       // take documents 0-9
		//Debug(true).        // print request and response to stdout
		//Pretty(true). // pretty print request and response JSON
		Do() // execute
	if err != nil {
		// Handle error
		beego.Error(err)
		return docs
	}
	// searchResult is of type SearchResult and returns hits, suggestions,
	// and all kinds of other information from Elasticsearch.
	beego.BeeLogger.Debug("Query took %d milliseconds", searchResult.TookInMillis)

	// Here's how you iterate through results with full control over each step.

	if searchResult.Hits != nil {
		beego.BeeLogger.Debug("Found a total of %d tweets", searchResult.Hits.TotalHits)

		// Iterate through results
		for _, hit := range searchResult.Hits.Hits {
			// hit.Index contains the name of the index
			v, err := json.Marshal(hit)
			if err != nil {
				beego.Error(err)
			}
			beego.BeeLogger.Debug("Found hit %s by current search keyworld [%s]", v, query)

			// Deserialize hit.Source into a Tweet (could also be just a map[string]interface{}).
			var t ElasticBlog
			err = json.Unmarshal(*hit.Source, &t)
			if err != nil {
				// Deserialization failed
				beego.Error(err)
			}

			bb_id, err := strconv.ParseInt(t.Id, 10, 64)
			if err != nil {
				beego.Error(err)
				continue
			}
			if !isAdmin && blog.IsAuth(bb_id) {
				beego.Debug("You can not read this blog since you don't have permission.")
				continue
			}
			wb := blog.OneById(bb_id)

			docs = append(docs, *wb)
		}
	} else {
		// No hits
		beego.Debug("Found no blogs")
	}

	return docs
}
