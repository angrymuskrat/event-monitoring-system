package crawler

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"

	"github.com/angrymuskrat/event-monitoring-system/services/insta-crawler/crawler/data"
	"github.com/angrymuskrat/event-monitoring-system/services/insta-crawler/crawler/parser"
	"github.com/corpix/uarand"
	"github.com/visheratin/unilog"
	"go.uber.org/zap"
	"golang.org/x/net/proxy"
)

type worker struct {
	id         int
	tid        int
	inCh       chan entity
	outCh      chan entity
	postsCh    chan []data.Post
	entitiesCh chan data.Entity
	mediaCh    chan []data.Media
	paramsCh   chan Parameters
	fixer      Fixer
	mu         sync.Mutex
	params     Parameters
	agent      string
	http       http.Client
	tor        http.Client
	cl         *client
}

func (w *worker) init(port int) {
	w.http = http.Client{
		Timeout: 30 * time.Second,
	}
	tbProxyURL, err := url.Parse("socks5://127.0.0.1:" + strconv.Itoa(port))
	if err != nil {
		return
	}
	tbDialer, err := proxy.FromURL(tbProxyURL, proxy.Direct)
	if err != nil {
		return
	}
	tbTransport := &http.Transport{
		Dial:                tbDialer.Dial,
		MaxIdleConnsPerHost: 1,
	}
	w.tor = http.Client{
		Transport: tbTransport,
		Timeout:   30 * time.Second,
	}
	w.agent = uarand.GetRandom()
	go w.paramsEdit()
	unilog.Logger().Info("started worker", zap.Int("id", w.id))
}

func (w *worker) paramsEdit() {
	for p := range w.paramsCh {
		w.mu.Lock()
		w.params = p
		fixer, err := NewFixer(p.Locations)
		if err == nil {
			w.fixer = fixer
		}
		w.mu.Unlock()
	}
}

func (w *worker) start() {
	for e := range w.inCh {
		w.proceedLocation(e)
		time.Sleep(2500 * time.Millisecond)
	}
}

func (w *worker) proceedLocation(e entity) {
	defer func() {
		w.outCh <- e
	}()
	e.id = url.QueryEscape(e.id)
	var req string
	loadEntity := false
	if e.checkpoint == "" {
		req = "https://www.instagram.com/graphql/query/?query_hash=1b84447a4d8b6d6d0426fefb34514485&variables=%7B%22id%22%3A%22" +
			e.id + "%22%2C%22first%22%3A50%7D"
		loadEntity = true
	} else {
		req = "https://www.instagram.com/graphql/query/?query_hash=1b84447a4d8b6d6d0426fefb34514485&variables=%7B%22id%22%3A%22" +
			e.id + "%22%2C%22first%22%3A50%2C%22after%22%3A%22" + e.checkpoint + "%22%7D"
	}
	f, cp := w.extractData(req, loadEntity, e.id)
	e.finished = f
	e.checkpoint = cp
}

func (w *worker) extractData(req string, loadEntity bool, entityID string) (bool, string) {
	rawData, err := w.makeRequest(req, true)
	if err != nil {
		return false, ""
	}
	cursor, hasNext, timestamp, zeroPosts, err := w.proceedResponse(rawData, loadEntity, entityID)
	if err != nil {
		return true, ""
	}
	if zeroPosts {
		//unilog.Logger().Info("zero posts", zap.String("req", req))
		return true, ""
	}
	cp := cursor
	f := false
	if timestamp < w.params.FinishTimestamp {
		//unilog.Logger().Info("before finish", zap.String("req", req))
		f = true
	}
	if !hasNext {
		f = true
	}
	return f, cp
}

func filterPosts(posts []data.Post, finish int64) []data.Post {
	res := make([]data.Post, 0, len(posts))
	for _, p := range posts {
		if p.Timestamp >= finish {
			res = append(res, p)
		}
	}
	return res
}

func (w *worker) makeRequest(request string, useTor bool) ([]byte, error) {
	req, err := http.NewRequest("GET", request, nil)
	if err != nil {
		unilog.Logger().Error("unable to create request", zap.String("URL", request), zap.Error(err))
		return nil, err
	}
	req.Header.Set("user-agent", w.agent)

	var resp *http.Response
	if useTor {
		resp, err = w.tor.Do(req)
	} else {
		resp, err = w.http.Do(req)
	}
	if err != nil {
		unilog.Logger().Error("unable to make request", zap.String("URL", request), zap.Error(err))
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == 429 {
		msg := fmt.Sprintf("too many requests from worker %d", w.id)
		// unilog.Logger().Error(msg)
		err = errors.New(msg)
		//time.Sleep(10 * time.Second)
		return nil, err
	}
	if resp.StatusCode == 500 {
		return w.cl.makeRequest(request, w.tid)
	}
	if resp.StatusCode == 404 {
		msg := "entity page was not found"
		unilog.Logger().Error(msg, zap.String("URL", request))
		err = errors.New(msg)
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	return body, nil
}

func (w *worker) proceedResponse(d []byte, loadEntity bool, entityID string) (endCursor string, hasNext bool, timestamp int64,
	zeroPosts bool, err error) {
	var posts []data.Post
	var location data.Location
	posts, location, endCursor, hasNext, timestamp, err = parser.ParseFromLocationRequest(d)
	if err != nil {
		unilog.Logger().Error("error during parsing response",
			zap.String("data", string(d)), zap.String("entity", entityID), zap.Error(err))
		return
	}
	if loadEntity {
		w.entitiesCh <- &location
	}
	if w.params.DetailedPosts {
		for i := 0; i < len(posts); i++ {
			detailedPost, err := w.detailedPost(posts[i])
			time.Sleep(50 * time.Millisecond)
			if err == nil {
				posts[i] = detailedPost
			}
		}
	}
	if w.params.LoadMedia {
		media := make([]data.Media, len(posts))
		for i := 0; i < len(posts); i++ {
			imgData, err := w.makeRequest(posts[i].ImageURL, false)
			time.Sleep(50 * time.Millisecond)
			if err != nil {
				continue
			}
			media[i] = data.Media{
				PostID: posts[i].ID,
				Data:   imgData,
			}
		}
		w.mediaCh <- media
	}
	if w.fixer.Init {
		posts = w.fixer.Fix(posts)
	}
	posts = filterPosts(posts, w.params.FinishTimestamp)
	w.postsCh <- posts
	zeroPosts = len(posts) == 0
	return
}

func (w worker) detailedPost(post data.Post) (data.Post, error) {
	request := "https://www.instagram.com/Params/" + post.Shortcode + "?__a=1"
	rawData, err := w.makeRequest(request, false)
	if err != nil {
		return data.Post{}, err
	}
	detailedPost, err := parser.ParseFromPostRequest(rawData)
	if err != nil {
		return data.Post{}, err
	}
	return detailedPost, nil
}
