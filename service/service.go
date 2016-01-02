package service

import (
	"github.com/categolj/categolj3-protos/categolj3"
	"github.com/golang/protobuf/proto"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

type Service struct {
	client *http.Client
	apiURL string
}

type Pageable struct {
	page int64
	size int64
}

func NewService(apiURL string) *Service {
	service := &Service{client: &http.Client{}, apiURL: apiURL}
	return service
}

func (s *Service) checkPageable(pageable *Pageable) {
	if pageable.page < 0 {
		pageable.page = 0
	}
	if pageable.size < 3 {
		pageable.size = 3
	}
}

func (s *Service) doRequest(path string, pageable *Pageable, excludeContent bool) ([]byte, error) {
	target, err := url.Parse(s.apiURL + "/" + path)
	if err != nil {
		return nil, err
	}
	q, err := url.ParseQuery(target.RawQuery)
	if err != nil {
		return nil, err
	}
	if excludeContent {
		q.Set("excludeContent", "true")
	}
	s.checkPageable(pageable)
	q.Set("page", strconv.FormatInt(pageable.page, 10))
	q.Set("size", strconv.FormatInt(pageable.size, 10))

	target.RawQuery = q.Encode()
	req, err := http.NewRequest("GET", target.String(), nil)

	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/x-protobuf")

	res, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func (s *Service) getEntriePage(path string, pageable *Pageable, excludeContent bool) (*categolj3.EntryPage, error) {
	body, err := s.doRequest(path, pageable, excludeContent)
	if err != nil {
		return nil, err
	}
	page := &categolj3.EntryPage{}
	err = proto.Unmarshal(body, page)
	if err != nil {
		return nil, err
	}
	return page, nil

}

func (s *Service) GetEntries(pageable *Pageable, excludeContent bool) (*categolj3.EntryPage, error) {
	return s.getEntriePage("entries", pageable, excludeContent)
}

func (s *Service) SearchEntries(pageable *Pageable, q string, excludeContent bool) (*categolj3.EntryPage, error) {
	return s.getEntriePage("entries?q="+q, pageable, excludeContent)
}

func (s *Service) GetEntriesByCreatedBy(pageable *Pageable, createdBy string, excludeContent bool) (*categolj3.EntryPage, error) {
	return s.getEntriePage("users/"+createdBy+"/entries", pageable, excludeContent)
}

func (s *Service) GetEntriesByUpdatedBy(pageable *Pageable, updatedBy string, excludeContent bool) (*categolj3.EntryPage, error) {
	return s.getEntriePage("users/"+updatedBy+"/entries?updated=true", pageable, excludeContent)
}

func (s *Service) GetEntriesByTag(pageable *Pageable, tag string, excludeContent bool) (*categolj3.EntryPage, error) {
	return s.getEntriePage("tags/"+tag+"/entries", pageable, excludeContent)
}

func (s *Service) GetEntriesByCategories(pageable *Pageable, categories string, excludeContent bool) (*categolj3.EntryPage, error) {
	return s.getEntriePage("categories/"+categories+"/entries", pageable, excludeContent)
}

func (s *Service) GetEntry(entryId int64) (*categolj3.Entry, error) {
	body, err := s.doRequest("entries/"+strconv.FormatInt(entryId, 10), &Pageable{}, false)
	if err != nil {
		return nil, err
	}
	entry := &categolj3.Entry{}
	err = proto.Unmarshal(body, entry)
	if err != nil {
		return nil, err
	}
	return entry, nil

}
