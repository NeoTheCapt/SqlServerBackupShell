package utils

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// HttpPost send post request and return response
func HttpPost(baseUrl string, body string, headers []string, proxyUrl string, timeout int) (string, error) {
	var client *http.Client
	if proxyUrl != "" {
		proxy, _ := url.Parse(proxyUrl)
		tr := &http.Transport{
			Proxy:           http.ProxyURL(proxy),
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client = &http.Client{
			Timeout:   time.Second * time.Duration(timeout),
			Transport: tr,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		}
	} else {
		client = &http.Client{
			//Transport: tr,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		}
	}
	req, err := http.NewRequest("POST", baseUrl, strings.NewReader(body))
	if err != nil {
		return "", err
	}
	for _, header := range headers {
		kv := strings.Split(header, ":")
		if len(kv) != 2 {
			continue
		}
		//bl.headers.Add(kv[0], kv[1])
		req.Header[kv[0]] = append(req.Header[kv[0]], kv[1])
	}
	//for k, v := range headers {
	//	req.Header[k] = v
	//}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(respBody), nil
}

// HttpGet send get request and return response
func HttpGet(baseUrl string, headers []string, proxyUrl string, timeout int) (string, error) {
	var client *http.Client
	if proxyUrl != "" {
		proxy, _ := url.Parse(proxyUrl)
		tr := &http.Transport{
			Proxy:           http.ProxyURL(proxy),
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client = &http.Client{
			Timeout:   time.Second * time.Duration(timeout),
			Transport: tr,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		}
	} else {
		client = &http.Client{
			//Transport: tr,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		}
	}
	req, err := http.NewRequest("GET", baseUrl, nil)
	if err != nil {
		return "", err
	}
	for _, header := range headers {
		kv := strings.Split(header, ":")
		if len(kv) != 2 {
			continue
		}
		//bl.headers.Add(kv[0], kv[1])
		req.Header[kv[0]] = append(req.Header[kv[0]], kv[1])
	}
	//for k, v := range headers {
	//	req.Header[k] = v
	//}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(respBody), nil
}

func GetFileContent(path string) string {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Print(err)
	}
	return string(b)
}
