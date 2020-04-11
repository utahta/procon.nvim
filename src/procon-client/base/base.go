package base

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/utahta/procon.nvim/src/procon-client/base/internal/cookiejar"
)

type (
	HTTPClient struct {
		httpClient http.Client

		cookiesPath string
		cookies     cookies
	}

	cookies map[string][]*http.Cookie
)

func NewHTTPClient(u *url.URL, cacheDir string) (*HTTPClient, error) {
	cookiesPath := filepath.Join(cacheDir, "cookies")

	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}
	c := &HTTPClient{
		httpClient: http.Client{Jar: jar},

		cookiesPath: cookiesPath,
		cookies:     cookies{},
	}

	if err := c.loadCookies(u); err != nil {
		return nil, err
	}
	return c, nil
}

func (cli *HTTPClient) Do(req *http.Request) (*http.Response, error) {
	resp, err := cli.httpClient.Do(req)
	if err != nil {
		return resp, err
	}

	u := &url.URL{
		Scheme: req.URL.Scheme,
		Host:   req.Host,
	}
	if err := cli.saveCookies(u); err != nil {
		resp.Body.Close()
		return nil, err
	}
	return resp, err
}

func (cli *HTTPClient) CheckRedirect(fn func(req *http.Request, via []*http.Request) error) {
	cli.httpClient.CheckRedirect = fn
}

func (cli *HTTPClient) Cookie(u *url.URL, name string) *http.Cookie {
	cs := cli.httpClient.Jar.Cookies(u)
	for _, c := range cs {
		if c.Name == name {
			return c
		}
	}
	return nil
}

func (cli *HTTPClient) loadCookies(u *url.URL) error {
	fp, err := os.OpenFile(cli.cookiesPath, os.O_RDONLY|os.O_CREATE, 0600)
	if err != nil {
		return errors.Wrap(err, "failed to open or create a cookies file")
	}
	defer fp.Close()

	b, err := ioutil.ReadAll(fp)
	if err != nil {
		return err
	}
	if len(b) == 0 {
		return nil
	}

	if err := json.Unmarshal(b, &cli.cookies); err != nil {
		return err
	}
	cs, ok := cli.cookies[u.String()]
	if !ok {
		return nil
	}
	cli.httpClient.Jar.SetCookies(u, cs)
	return nil
}

func (cli *HTTPClient) saveCookies(u *url.URL) error {
	cs := cli.httpClient.Jar.Cookies(u)
	cli.cookies[u.String()] = cs

	b, err := json.Marshal(cli.cookies)
	if err != nil {
		return err
	}

	fp, err := os.OpenFile(cli.cookiesPath, os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer fp.Close()

	if _, err := fp.Write(b); err != nil {
		return err
	}
	return nil
}
