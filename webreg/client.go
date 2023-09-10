package webreg

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type WebregClient struct {
	cookies string
	term    Term
}

func NewWebregClient(term Term) *WebregClient {
	return &WebregClient{
		term: term,
	}
}

func (c *WebregClient) GetCookiesFromServer() {
	resp, err := http.Get("http://localhost:3001")
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	c.cookies = string(data)
}

func (c WebregClient) GetCourseInfo(subject string, course string) (*CourseInfo, error) {
	client := &http.Client{}

	url := fmt.Sprintf("https://act.ucsd.edu/webreg2/svc/wradapter/secure/search-load-group-data?subjcode=%s&crsecode=+%s&termcode=%s&_=%d",
		subject,
		course,
		c.term.Code,
		time.Now().UnixMilli(),
	)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}

	req.Header.Set("Cookie", c.cookies)

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	sections := []*SectionInfo{}
	if err := json.Unmarshal(data, &sections); err != nil {
		panic(err)
	}

	return &CourseInfo{
		SubjectCode: subject,
		CourseCode:  course,
		Sections:    sections,
	}, nil
}
