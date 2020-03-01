package crawler

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/angrymuskrat/event-monitoring-system/services/insta-crawler/crawler/data"
	"github.com/angrymuskrat/event-monitoring-system/services/insta-crawler/crawler/parser"
	"github.com/angrymuskrat/event-monitoring-system/services/insta-crawler/crawler/storage"
	"github.com/corpix/uarand"
	"github.com/visheratin/unilog"
	"go.uber.org/zap"
	"golang.org/x/net/proxy"
)

type worker struct {
	id         int
	inCh       chan entity
	outCh      chan entity
	postsCh    chan []data.Post
	entitiesCh chan data.Entity
	mediaCh    chan []data.Media
	paramsCh   chan Parameters
	fixer      storage.Fixer
	mu         sync.Mutex
	params     Parameters
	agent      string
	http       http.Client
	tor        http.Client
}

const useTor = true

func (w *worker) init(port int) {
	//w.inCh = inCh
	//w.outCh = outCh
	w.http = http.Client{
		Timeout: 30 * time.Second,
	}
	tbProxyURL, err := url.Parse("socks5://127.0.0.1:" + strconv.Itoa(port+w.id))
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
	fixer, err := storage.NewFixer("./locations.json")
	if err == nil {
		w.fixer = fixer
	}
	go w.paramsEdit()
	unilog.Logger().Info("started worker", zap.Int("id", w.id))
}

func (w *worker) paramsEdit() {
	for p := range w.paramsCh {
		w.mu.Lock()
		w.params = p
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
	var cursor string
	var hasNext bool
	var timestamp int64
	var zeroPosts bool
	requestTemplatePt1 := "https://www.instagram.com/graphql/query/?query_hash=1b84447a4d8b6d6d0426fefb34514485&variables=%7B%22id%22%3A%22"
	requestTemplatePt2 := "%22%2C%22first%22%3A50%2C%22after%22%3A%22"
	requestTemplatePt3 := "%22%7D"
	var cp string
	if e.checkpoint == "" {
		initRequest := "https://www.instagram.com/graphql/query/?query_hash=1b84447a4d8b6d6d0426fefb34514485&variables=%7B%22id%22%3A%22" + e.id +
			"%22%2C%22first%22%3A50%7D"
		//referer := "https://www.instagram.com/explore/locations/" + e.id
		rawData, err := w.makeRequest(initRequest, true, "", "", false)
		if err != nil {
			if rawData != nil {
				e.finished = true
			}
			return
		}
		cursor, hasNext, _, zeroPosts, err = w.proceedResponse(rawData, w.params.FinishTimestamp, true)
		if err != nil {
			e.finished = true
			return
		}
		if zeroPosts {
			e.finished = true
			return
		}
		e.checkpoint = cursor
		if timestamp < w.params.FinishTimestamp {
			e.finished = true
		}
	} else {
		cursor = cp
		hasNext = true
	}
	if hasNext {
		var newRequest string
		//var referer string
		switch w.params.Type {
		case data.LocationsType:
			newRequest = requestTemplatePt1 + e.id + requestTemplatePt2 + cursor + requestTemplatePt3
			//referer = "https://www.instagram.com/explore/locations/" + e.id
		}
		//variables := "{\"ID\":\"" + e.id + "\",\"first\":50,\"after\":\"" + cursor + "\"}"
		//gisString := w.rhx + ":" + variables
		//h := md5.New()
		//io.WriteString(h, gisString)
		//gis := hex.EncodeToString(h.Sum(nil))
		rawData, err := w.makeRequest(newRequest, useTor, "", "", false)
		if err != nil {
			if rawData != nil {
				e.finished = true
			}
			return
		}
		cursor, hasNext, timestamp, zeroPosts, err = w.proceedResponse(rawData, w.params.FinishTimestamp, false)
		if err != nil {
			e.finished = true
			return
		}
		if zeroPosts {
			e.finished = true
			return
		}
		e.checkpoint = cursor
		if timestamp < w.params.FinishTimestamp {
			e.finished = true
		}
	} else {
		e.finished = true
	}
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

func (w *worker) makeRequest(request string, useTor bool, gis string, referer string, auth bool) ([]byte, error) {
	req, err := http.NewRequest("GET", request, nil)
	if err != nil {
		unilog.Logger().Error("unable to create request", zap.String("URL", request), zap.Error(err))
		return nil, err
	}
	req.Header.Set("user-agent", w.agent)
	// req.Header.Set("dnt", "1")
	// req.Header.Set("referer", referer)
	// req.Header.Set("x-requested-with", "XMLHttpRequest")
	//if gis != "" {
	//	req.Header.Set("x-instagram-gis", gis)
	//}
	if w.params.AuthCookie != "" && auth {
		req.Header.Set("cookie", w.params.AuthCookie)
	}
	//for _, cookie := range w.cookies {
	//	req.AddCookie(cookie)
	//}

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
	//if resp.StatusCode == 403 {
	//	w.getCookies()
	//	unilog.Logger().Error("cookies have expired", zap.String("URL", request), zap.Error(err))
	//	return nil, err
	//}
	if resp.StatusCode == 429 {
		msg := fmt.Sprintf("too many requests from worker %d", w.id)
		// unilog.Logger().Error(msg)
		err = errors.New(msg)
		time.Sleep(10 * time.Second)
		return nil, err
	}
	if resp.StatusCode == 404 || resp.StatusCode == 500 {
		msg := "entity page was not found"
		unilog.Logger().Error(msg, zap.String("URL", request))
		err = errors.New(msg)
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	return body, nil
}

func (w *worker) getCookies() error {
	request := "https://www.instagram.com/"
	req, err := http.NewRequest("GET", request, nil)
	if err != nil {
		unilog.Logger().Error("unable to create request", zap.String("URL", request), zap.Error(err))
		return err
	}
	req.Header.Set("user-agent", w.agent)
	var resp *http.Response
	resp, err = w.tor.Do(req)
	if err != nil {
		unilog.Logger().Error("unable to make request", zap.String("URL", request), zap.Error(err))
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	bodyString := string(body)
	re := regexp.MustCompile("\"rhx_gis\":\".*\",*")
	data := re.FindAllString(bodyString, -1)
	if len(data) > 0 {
		id := data[0]
		id = strings.Replace(id, "\"rhx_gis\":\"", "", -1)
		id = strings.Split(id, ",")[0]
		id = strings.Replace(id, "\"", "", -1)
		id = strings.Replace(id, ",", "", -1)
		//w.rhx = id
	}
	//cookies := resp.Cookies()
	//w.cookies = cookies
	//for _, cookie := range cookies {
	//	if cookie.Name == "csrftoken" {
	//		w.token = cookie.Value
	//		break
	//	}
	//}
	return nil
}

func (w *worker) proceedResponse(d []byte, finish int64, loadEntity bool) (endCursor string, hasNext bool, timestamp int64,
	zeroPosts bool, err error) {
	var posts []data.Post
	switch w.params.Type {
	case data.ProfilesType:
		var profile data.Profile
		posts, profile, endCursor, hasNext, timestamp, err = parser.ParseFromProfileRequest(d)
		if err != nil {
			return
		}
		if loadEntity {
			w.entitiesCh <- &profile
		}
	case data.LocationsType:
		var location data.Location
		posts, location, endCursor, hasNext, timestamp, err = parser.ParseFromLocationRequest(d)
		if err != nil {
			return
		}
		if loadEntity {
			w.entitiesCh <- &location
		}
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
			referer := "https://www.instagram.com/Params/" + posts[i].Shortcode
			imgData, err := w.makeRequest(posts[i].ImageURL, false, "", referer, false)
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
	posts = filterPosts(posts, finish)
	w.postsCh <- posts
	zeroPosts = len(posts) == 0
	return
}

func (w worker) detailedPost(post data.Post) (data.Post, error) {
	request := "https://www.instagram.com/Params/" + post.Shortcode + "?__a=1"
	referer := "https://www.instagram.com/Params/" + post.Shortcode
	rawData, err := w.makeRequest(request, false, "", referer, false)
	if err != nil {
		return data.Post{}, err
	}
	detailedPost, err := parser.ParseFromPostRequest(rawData)
	if err != nil {
		return data.Post{}, err
	}
	return detailedPost, nil
}
