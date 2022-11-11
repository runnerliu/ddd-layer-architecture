package es

import (
	"bytes"
	"ddd-demo/common/consts"
	"encoding/json"
	"io/ioutil"
	"sync"

	"github.com/elastic/go-elasticsearch/v7"
)

// ESClient ES 客户端
type ESClient struct {
	esClient *elasticsearch.Client
}

// QueryIndexData 查询索引数据
func (e *ESClient) QueryIndexData(index string, params map[string]interface{}) (map[string]interface{}, error) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(params); err != nil {
		return nil, err
	}
	resp, err := e.esClient.Search(
		e.esClient.Search.WithIndex(index),
		e.esClient.Search.WithDocumentType("_doc"),
		e.esClient.Search.WithBody(&buf))
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			panic(err)
		}
	}()
	if resp.StatusCode < 200 || resp.StatusCode > 300 {
		return nil, consts.ErrESQueryIndexData
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	result := make(map[string]interface{})
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return result, err
}

var (
	defaultESClientOnce sync.Once
	defaultESClient     *ESClient
)

// NewESClient 创建 ES 客户端
func NewESClient(host, username, password string) *ESClient {
	defaultESClientOnce.Do(func() {
		client, err := elasticsearch.NewClient(elasticsearch.Config{
			Addresses:  []string{host},
			Username:   username,
			Password:   password,
			MaxRetries: 3,
		})
		if err == nil {
			defaultESClient = &ESClient{
				esClient: client,
			}
		}
	})
	return defaultESClient
}
