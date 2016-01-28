package sources

import (
	"crypto/tls"
	"io/ioutil"
	"net/http"

	"github.com/zeroturnaround/configo/parsers"
)

type HTTPSource struct {
	URL    string `json:"url"`
	Format string `json:"format"`
	TLS    struct {
		Cert string `json:"cert"`
		Key  string `json:"key"`
	} `json:"tls"`
}

func (httpSource *HTTPSource) Get() (map[string]interface{}, error) {
	transport := &http.Transport{}

	if httpSource.TLS.Cert != "" {
		cert, err := tls.X509KeyPair([]byte(httpSource.TLS.Cert), []byte(httpSource.TLS.Key))
		if err != nil {
			return nil, err
		}

		transport.TLSClientConfig = &tls.Config{
			Certificates: []tls.Certificate{cert},
		}
	}

	client := &http.Client{Transport: transport}

	response, err := client.Get(httpSource.URL)

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
