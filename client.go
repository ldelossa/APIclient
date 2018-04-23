import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
)

// do will return a http.Response along with an error if the http.Client fails to
// perform it's request. If a pointer 'v' is passed in the response body will be
// json decoded into the passed pointer.
func (c *HTTPClient) do(req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	if v != nil {
		err = json.NewDecoder(resp.Body).Decode(v)
	}
	return resp, err
}

// newRequest is a helper functions which makes any supported HTML request. If 'body' is provided the
// passed in value will be encoded to json and placed into the http.Request's body. path excepts a
// url.URL with the path field defined. path is appended to c.BaseURL to from the location of the http.Request
func (c *HTTPClient) newRequest(method string, path *url.URL, body interface{}) (*http.Request, error) {
	u := c.BaseURL.ResolveReference(path)

	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.UserAgent)
	return req, nil
}

