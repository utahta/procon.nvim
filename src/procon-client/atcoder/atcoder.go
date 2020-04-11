package atcoder

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/pkg/errors"
	"github.com/utahta/procon.nvim/src/procon-client/base"
)

type (
	Client struct {
		httpClient *base.HTTPClient
		URL        *url.URL

		examplesCache map[string][]Example
	}

	Example struct {
		Input  string
		Output string
	}
)

const (
	LanguageIDCPP14gcc = "3003"
	LanguageIDCPP17gcc = "4003"
)

func New(cacheDir string) (*Client, error) {
	u, err := url.Parse("https://atcoder.jp")
	if err != nil {
		return nil, err
	}
	c, err := base.NewHTTPClient(u, cacheDir)
	if err != nil {
		return nil, err
	}

	return &Client{
		httpClient: c,
		URL:        u,

		examplesCache: make(map[string][]Example),
	}, nil
}

func (c *Client) IsLoggedIn() bool {
	session := c.httpClient.Cookie(c.URL, "REVEL_SESSION")
	if session == nil {
		return false
	}
	return strings.Contains(session.Value, "UserName%3A")
}

func (c *Client) Login(ctx context.Context, username, password string) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	loginURL := url.URL{
		Scheme: c.URL.Scheme,
		Host:   c.URL.Host,
		Path:   "login",
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, loginURL.String(), nil)
	if err != nil {
		return errors.Wrap(err, "failed to create a request")
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "failed to request")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.Errorf("invalid http status code was returned: %v", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return errors.Wrap(err, "failed to create a goquery document")
	}
	csrfToken, _ := doc.Find(`input[name="csrf_token"]`).Attr("value")

	postForm := url.Values{}
	postForm.Set("username", username)
	postForm.Set("password", password)
	postForm.Set("csrf_token", csrfToken)
	req, err = http.NewRequestWithContext(ctx, http.MethodPost, loginURL.String(), strings.NewReader(postForm.Encode()))
	if err != nil {
		return err
	}
	req.Header.Add("Referer", "https://atcoder.jp/login")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, err = c.httpClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "failed to request")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.Errorf("invalid status code was returned: %d", resp.StatusCode)
	}

	session := c.httpClient.Cookie(c.URL, "REVEL_SESSION")
	if session == nil || !strings.Contains(session.Value, fmt.Sprintf("UserName%%3A%s", username)) {
		return errors.Errorf("failed to login: session is nil or `UserName:%s` does not contain", username)
	}
	return nil
}

func (c *Client) ListExamples(ctx context.Context, contest, task string) ([]Example, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	taskURL := url.URL{
		Scheme: c.URL.Scheme,
		Host:   c.URL.Host,
		Path:   path.Join("contests", contest, "tasks", task),
	}
	if v, ok := c.examplesCache[taskURL.String()]; ok {
		return v, nil
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, taskURL.String(), nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create a request")
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to do a request")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.Errorf("http status code is %v", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create a goquery document")
	}

	var examples []Example
	doc.Find("div.part").Each(func(_ int, sel *goquery.Selection) {
		if strings.HasPrefix(sel.Find("h3").Text(), "入力例") {
			examples = append(examples, Example{
				Input: strings.TrimSpace(sel.Find("pre").Text()),
			})
		}
		if strings.HasPrefix(sel.Find("h3").Text(), "出力例") {
			examples[len(examples)-1].Output = strings.TrimSpace(sel.Find("pre").Text())
		}
	})

	if len(examples) == 0 {
		doc.Find("section").Each(func(_ int, sel *goquery.Selection) {
			if strings.HasPrefix(sel.Find("h3").Text(), "入力例") {
				examples = append(examples, Example{
					Input: strings.TrimSpace(sel.Find("pre").Text()),
				})
			}
			if strings.HasPrefix(sel.Find("h3").Text(), "出力例") {
				examples[len(examples)-1].Output = strings.TrimSpace(sel.Find("pre").Text())
			}
		})
	}

	if len(examples) == 0 {
		return nil, errors.New("example is not found")
	}
	c.examplesCache[taskURL.String()] = examples
	return examples, nil
}

func (c *Client) Submit(ctx context.Context, contest string, task string, languageID string, sourceCode string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	taskURL := url.URL{
		Scheme: c.URL.Scheme,
		Host:   c.URL.Host,
		Path:   path.Join("contests", contest, "tasks", task),
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, taskURL.String(), nil)
	if err != nil {
		return "", errors.Wrap(err, "failed to create a request")
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", errors.Wrap(err, "failed to request")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errors.Errorf("invalid http status code was returned: %v", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", errors.Wrap(err, "failed to create a goquery document")
	}
	action, _ := doc.Find(`form[class="form-horizontal form-code-submit"]`).Attr("action")
	dataTaskScreenName, _ := doc.Find(`input[name="data.TaskScreenName"]`).Attr("value")
	csrfToken, _ := doc.Find(`input[name="csrf_token"]`).Attr("value")

	var languageIDfound bool
	doc.Find(`option[data-mime="text/x-c++src"]`).Each(func(_ int, sec *goquery.Selection) {
		if v, _ := sec.Attr("value"); v == languageID {
			languageIDfound = true
		}
	})
	if !languageIDfound {
		return "", errors.Errorf("%s is not supported in this task. please use ProconResetBuilder command to change it to appropriate builder", languageID)
	}

	submitURL := url.URL{
		Scheme: c.URL.Scheme,
		Host:   c.URL.Host,
		Path:   action,
	}
	postForm := url.Values{}
	postForm.Set("data.TaskScreenName", dataTaskScreenName)
	postForm.Set("data.LanguageId", languageID)
	postForm.Set("csrf_token", csrfToken)
	postForm.Set("sourceCode", sourceCode)
	req, err = http.NewRequestWithContext(ctx, http.MethodPost, submitURL.String(), strings.NewReader(postForm.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Add("Referer", taskURL.String())
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	c.httpClient.CheckRedirect(func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	})
	resp, err = c.httpClient.Do(req)
	if err != nil {
		return "", errors.Wrap(err, "failed to request")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusFound {
		return "", errors.Errorf("invalid http status code was returned: %d", resp.StatusCode)
	}

	redirectURL := url.URL{
		Scheme: c.URL.Scheme,
		Host:   c.URL.Host,
		Path:   resp.Header.Get("location"),
	}
	return redirectURL.String(), nil
}
