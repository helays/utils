package esClose

import (
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"helay.net/go/utils/v3/close/vclose"
)

// CloseResp 关闭esapi.Response
func CloseResp(resp *esapi.Response) {
	if resp != nil {
		vclose.Close(resp.Body)
	}
}
