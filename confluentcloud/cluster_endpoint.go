package confluentcloud

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"net/url"
)

const (
	baseURLSuffix = "2.0/kafka/"
)

type KafkaClusterClient struct {
	KafkaApiEndpoint *url.URL
	BaseURLSuffix    string
	BaseURL          *url.URL
	client           *resty.Client
	token            string
}

func NewKafkaClusterClient(kafkaApiEndpoint *url.URL, clusterID string, token string) *KafkaClusterClient {
	_baseURL := fmt.Sprintf("%s/%s%s/", kafkaApiEndpoint, baseURLSuffix, clusterID)
	baseURL, _ := url.Parse(_baseURL)

	client := resty.New()
	client.SetDebug(false)
	c := &KafkaClusterClient{KafkaApiEndpoint: kafkaApiEndpoint, BaseURL: baseURL, BaseURLSuffix: baseURLSuffix}
	c.client = client
	c.token = token

	return c
}

func (c *KafkaClusterClient) NewKafkaClusterRequest() *resty.Request {
	return c.client.R()
}
