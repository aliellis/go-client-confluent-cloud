package confluentcloud

import (
	"fmt"
	"net/url"
)

type ACLRequest struct {
	PatternFilter *PatternFilter `json:"patternFilter"`
	EntryFilter   *EntryFilter   `json:"entryFilter"`
}
type PatternFilter struct {
	ResourceType string `json:"resourceType"`
	PatternType  string `json:"patternType"`
}
type EntryFilter struct {
	Operation      string `json:"operation"`
	Host           string `json:"host"`
	PermissionType string `json:"permissionType"`
}

type ListACLResponse = []ACL
type ACL struct {
	Pattern Pattern `json:"pattern"`
	Entry   Entry   `json:"entry"`
}
type Pattern struct {
	ResourceType string `json:"resourceType"`
	Name         string `json:"name"`
	PatternType  string `json:"patternType"`
}
type Entry struct {
	Principal      string `json:"principal"`
	Host           string `json:"host"`
	Operation      string `json:"operation"`
	PermissionType string `json:"permissionType"`
}

type ACLCreateRequestW = []ACLCreateRequest
type ACLCreateRequest struct {
	Pattern *Pattern `json:"pattern"`
	Entry   *Entry   `json:"entry"`
}
type CreateACLResponse = []ACLCreateRequest

func (c *Client) ListACLs(apiEndpoint *url.URL, clusterID string, aclRequest *ACLRequest) ([]ACL, error) {
	token, err := c.GetKafkaClusterAccessToken()
	if err != nil {
		return nil, err
	}
	cc := NewKafkaClusterClient(apiEndpoint, clusterID, *token)

	// cannot use url.Parse due to colon being interpreted as scheme
	suffix := url.URL{
		Path: "acls:search",
	}
	u := cc.BaseURL.ResolveReference(&suffix)
	response, err := cc.NewKafkaClusterRequest().
		SetAuthToken(*token).
		SetBody(aclRequest).
		SetResult(&ListACLResponse{}).
		Post(u.String())

	if err != nil {
		return nil, err
	}

	// response is raw, cannot parse from JSON
	if response.IsError() {
		return nil, fmt.Errorf("list_acls: %s", response.Body())
	}

	result := response.Result().(*ListACLResponse)
	return *result, nil
}

func (c *Client) CreateACLs(apiEndpoint *url.URL, clusterID string, aclCreateRequestW *ACLCreateRequestW) (interface{}, error) {
	token, err := c.GetKafkaClusterAccessToken()
	if err != nil {
		return nil, err
	}
	cc := NewKafkaClusterClient(apiEndpoint, clusterID, *token)

	rel, err := url.Parse("acls")
	if err != nil {
		return nil, err
	}

	u := cc.BaseURL.ResolveReference(rel)

	response, err := c.NewRequest().
		SetAuthToken(*token).
		SetBody(aclCreateRequestW).
		SetResult(&CreateACLResponse{}).
		Post(u.String())
	if err != nil {
		return nil, err
	}
	if response.IsError() {
		return nil, fmt.Errorf("create_acls: %s", response.Body())
	}

	result := response.Result().(*CreateACLResponse)
	return *result, nil
}
