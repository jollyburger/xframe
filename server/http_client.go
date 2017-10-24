package server

import (
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"time"
	"xframe/trace"

	"github.com/juju/ratelimit"
	"github.com/opentracing-contrib/go-stdlib/nethttp"
)

func sendHttpRequest(ctx trace.XContext, url_path string, params map[string]interface{}, timeOut uint32) (res []byte, err error) {
	req_url, err := url.Parse(url_path)
	if err != nil {
		return
	}
	req_params := req_url.Query()
	for k, v := range params {
		req_params.Set(k, v.(string))
	}
	req_url.RawQuery = req_params.Encode()
	// timeout config, 0 means no-timeout
	//client := newTimeoutHTTPClient(time.Duration(timeOut) * time.Second)
	client := &http.Client{Transport: &nethttp.Transport{}, Timeout: time.Duration(timeOut) * time.Second}
	req, err := http.NewRequest("GET", req_url.String(), nil)
	if err != nil {
		return
	}
	if ctx != nil {
		var (
			ht *nethttp.Tracer
		)
		req = req.WithContext(ctx)
		req, ht = nethttp.TraceRequest(trace.GTracer, req, nethttp.OperationName("HTTP GET: "+url_path))
		defer ht.Finish()
	}
	result, err := client.Do(req)
	if err != nil {
		return
	}
	defer result.Body.Close()
	//new_span_id, err := strconv.Atoi(result.Header.Get("X-Trace-Context"))
	res, err = ioutil.ReadAll(result.Body)
	return
}

func sendHttpPostRequest(ctx trace.XContext, url_path string, body_type string, body io.Reader, timeOut uint32) (res []byte, err error) {
	//client := newTimeoutHTTPClient(time.Duration(timeOut) * time.Second)
	client := &http.Client{Transport: &nethttp.Transport{}, Timeout: time.Duration(timeOut) * time.Second}
	req, err := http.NewRequest("POST", url_path, body)
	if err != nil {
		return
	}
	if ctx != nil {
		var (
			ht *nethttp.Tracer
		)
		req = req.WithContext(ctx)
		req, ht = nethttp.TraceRequest(trace.GTracer, req, nethttp.OperationName("HTTP POST: "+url_path))
		defer ht.Finish()
	}
	req.Header.Set("Content-Type", body_type)
	result, err := client.Do(req)
	if err != nil {
		return
	}
	defer result.Body.Close()
	//new_span_id, err := strconv.Atoi(result.Header.Get("X-Trace-Context"))
	res, err = ioutil.ReadAll(result.Body)
	return
}

func sendHttpMethodRequest(ctx trace.XContext, method string, url_path string, body io.Reader, timeOut uint32) (res []byte, err error) {
	req, err := http.NewRequest(method, url_path, body)
	if err != nil {
		return
	}
	if ctx != nil {
		var (
			ht *nethttp.Tracer
		)
		req = req.WithContext(ctx)
		req, ht = nethttp.TraceRequest(trace.GTracer, req, nethttp.OperationName("HTTP "+method+": "+url_path))
		defer ht.Finish()
	}
	//client := newTimeoutHTTPClient(time.Duration(timeOut) * time.Second)
	client := &http.Client{Transport: &nethttp.Transport{}, Timeout: time.Duration(timeOut) * time.Second}
	result, err := client.Do(req)
	if err != nil {
		return
	}
	defer result.Body.Close()
	//new_span_id, err := strconv.Atoi(result.Header.Get("X-Trace-Context"))
	res, err = ioutil.ReadAll(result.Body)
	return
}

func dialHTTPTimeout(timeOut time.Duration) func(net, addr string) (net.Conn, error) {
	return func(network, addr string) (c net.Conn, err error) {
		c, err = net.DialTimeout(network, addr, timeOut)
		if err != nil {
			return
		}
		if timeOut > 0 {
			c.SetDeadline(time.Now().Add(timeOut))
		}
		return
	}
}

func newTimeoutHTTPClient(timeOut time.Duration) *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			Dial: dialHTTPTimeout(timeOut),
		},
	}
}

//http client ratelimiter
func sendHttPRequestBylimit(ctx trace.XContext, method, urlStr string, reader io.Reader, header http.Header, rate float64, capacity int64) ([]byte, error) {
	//define ratelimit reader
	bucket := ratelimit.NewBucketWithRate(rate, capacity)
	bucket_reader := ratelimit.Reader(reader, bucket)
	req, err := http.NewRequest(method, urlStr, bucket_reader)
	if err != nil {
		return nil, err
	}
	if ctx != nil {
		var (
			ht *nethttp.Tracer
		)
		req = req.WithContext(ctx)
		req, ht = nethttp.TraceRequest(trace.GTracer, req, nethttp.OperationName("HTTP "+method+": "+urlStr))
		defer ht.Finish()
	}
	req.Header = header
	//client := new(http.Client)
	client := &http.Client{Transport: &nethttp.Transport{}}
	result, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer result.Body.Close()
	//new_span_id, err := strconv.Atoi(result.Header.Get("X-Trace-Context"))
	res, err := ioutil.ReadAll(result.Body)
	return res, err
}
