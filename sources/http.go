package sources

import (
	"io/ioutil"
	"net/http"

	"github.com/zeroturnaround/configo/parsers"
)

type HTTPSource struct {
	URL    string `json:"url"`
	Format string `json:"format"`
}

func (httpSource *HTTPSource) Get() (map[string]interface{}, error) {
	response, err := http.Get(httpSource.URL)

	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	content, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return nil, err
	}

	result := make(map[string]interface{})
	parsers.MustGetParser(httpSource.Format).Parse(content, result)

	return result, nil
}
