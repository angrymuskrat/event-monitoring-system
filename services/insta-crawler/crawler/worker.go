package crawler

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
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
	id            int
	entities      *entities
	params        Parameters
	sessionID     string
	sessionStatus *Status
	rootDir       string
	pCh           chan bool
	oCh           chan bool
	paused        bool
	agent         string
	cookies       []*http.Cookie
	token         string
	rhx           string
	checkpoints   map[string]string
	http          http.Client
	tor           http.Client
	savePosts     bool        // use data storage or not
	posts         []data.Post // tmp array of posts for sending to data-storage
}

const useTor = true

func (w *worker) init(port int) {
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
}

func (w *worker) start() {
	go func() {
		for {
			select {
			case p := <-w.pCh:
				w.paused = p
			default:
			}
			time.Sleep(1 * time.Second)
		}
	}()
	for {
		fmt.Print("") // this is important, this line fix freeze bug
		if w.paused {
			time.Sleep(1 * time.Second)
			continue
		}
		for i := range w.entities.data {
			if w.paused {
				// w.oCh <- true
				break
			}
			w.proceedLocation(i)
			time.Sleep(2500 * time.Millisecond)
		}
	}
}

func (w *worker) proceedLocation(i int) error {
	st, err := storage.Instance()
	if err != nil {
		unilog.Logger().Error("unable to get storage", zap.Error(err))
		return err
	}
	var cursor string
	var hasNext bool
	var timestamp int64
	var zeroPosts bool
	requestTemplatePt1 := "https://www.instagram.com/graphql/query/?query_hash=1b84447a4d8b6d6d0426fefb34514485&variables=%7B%22id%22%3A%22"
	requestTemplatePt2 := "%22%2C%22first%22%3A50%2C%22after%22%3A%22"
	requestTemplatePt3 := "%22%7D"
	var cp string
	id := w.entities.get(i)
	if id == "" {
		return nil
	}
	cp, ok := w.checkpoints[id]
	if !ok {
		cp = st.Checkpoint(w.sessionID, id)
	}
	if cp == "" {
		initRequest := "https://www.instagram.com/graphql/query/?query_hash=1b84447a4d8b6d6d0426fefb34514485&variables=%7B%22id%22%3A%22" + id +
			"%22%2C%22first%22%3A50%7D"
		referer := "https://www.instagram.com/explore/locations/" + id
		rawData, err := w.makeRequest(initRequest, true, "", referer, false)
		if err != nil {
			if rawData != nil {
				w.removeEntity(i)
			}
			return err
		}
		cursor, hasNext, _, zeroPosts, err = w.proceedResponse(rawData, id)
		if err != nil {
			return err
		}
		w.checkpoints[id] = cursor
		err = st.WriteCheckpoint(w.sessionID, id, cursor)
		if err != nil {
			return err
		}
	} else {
		cursor = cp
		hasNext = true
	}
	if hasNext {
		var newRequest string
		var referer string
		switch w.params.Type {
		case data.LocationsType:
			newRequest = requestTemplatePt1 + id + requestTemplatePt2 + cursor + requestTemplatePt3
			referer = "https://www.instagram.com/explore/locations/" + id
		case data.ProfilesType:
			//profile, err := worker.writerInstance.ReadProfile(entityID)
			//if err != nil {
			//	return
			//}
			//newRequest = requestTemplatePt1 + profile.ID + requestTemplatePt2 + cursor + requestTemplatePt3
			//referer = "https://www.instagram.com/" + profile.Username
		case data.StoriesType: // TODO : StoriesType

		}
		variables := "{\"ID\":\"" + id + "\",\"first\":50,\"after\":\"" + cursor + "\"}"
		gisString := w.rhx + ":" + variables
		h := md5.New()
		io.WriteString(h, gisString)
		gis := hex.EncodeToString(h.Sum(nil))
		rawData, err := w.makeRequest(newRequest, useTor, gis, referer, false)
		if err != nil {
			if rawData != nil {
				w.removeEntity(i)
			}
			return err
		}
		cursor, hasNext, timestamp, zeroPosts, err = w.proceedResponse(rawData, id)
		if err != nil {
			if zeroPosts {
				w.removeEntity(i)
			}
			return err
		}
		w.checkpoints[id] = cursor
		err = st.WriteCheckpoint(w.sessionID, id, cursor)
		if err != nil {
			return err
		}
		if timestamp < w.params.FinishTimestamp {
			w.removeEntity(i)
		}
	} else {
		w.removeEntity(i)
	}
	return nil
}

func (w *worker) removeEntity(i int) {
	w.entities.remove(i)
	//w.sessionStatus.updateEntitiesLeft(-1)
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
	if resp.StatusCode == 403 {
		w.getCookies()
		unilog.Logger().Error("cookies have expired", zap.String("URL", request), zap.Error(err))
		return nil, err
	}
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
		return []byte{}, err
	}
	cookies := resp.Cookies()
	w.cookies = cookies
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
		w.rhx = id
	}
	cookies := resp.Cookies()
	w.cookies = cookies
	for _, cookie := range cookies {
		if cookie.Name == "csrftoken" {
			w.token = cookie.Value
			break
		}
	}
	return nil
}

func (w *worker) proceedResponse(d []byte, entityID string) (endCursor string, hasNext bool, timestamp int64,
	zeroPosts bool, err error) {
	var posts []data.Post
	var st storage.Storage
	switch w.params.Type {
	case data.ProfilesType:
		var profile data.Profile
		posts, profile, endCursor, hasNext, timestamp, err = parser.ParseFromProfileRequest(d)
		if err != nil {
			return
		}
		st, err = storage.Instance()
		if err != nil {
			unilog.Logger().Error("unable to get storage", zap.Error(err))
			return
		}
		err = st.WriteEntity(w.sessionID, entityID, &profile)
	case data.LocationsType:
		var location data.Location
		posts, location, endCursor, hasNext, timestamp, err = parser.ParseFromLocationRequest(d)
		if err != nil {
			return
		}
		st, err = storage.Instance()
		if err != nil {
			unilog.Logger().Error("unable to get storage", zap.Error(err))
			return
		}
		err = st.WriteEntity(w.sessionID, entityID, &location)
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
		w.saveMedia(w.sessionID, entityID, media)
	}
	st, err = storage.Instance()
	if err != nil {
		unilog.Logger().Error("unable to get storage", zap.Error(err))
		return
	}
	if w.savePosts { // save posts to tmp array for sending to data storage
		w.posts = append(w.posts, posts...)
	}
	err = st.WritePosts(w.sessionID, posts)
	if len(posts) > 0 {
		w.sessionStatus.updatePostsCollected(len(posts))
	}
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

func (w worker) saveMedia(sessionID, entityID string, media []data.Media) {
	mediaPath := path.Join(w.rootDir, sessionID, "img", entityID)
	err := os.MkdirAll(mediaPath, 0777)
	if err != nil {
		unilog.Logger().Error("unable to create media directory", zap.String("path", mediaPath), zap.Error(err))
	}
	for _, item := range media {
		if item.PostID != "" {
			imgp := path.Join(mediaPath, item.PostID+".png")
			err = ioutil.WriteFile(imgp, item.Data, 0644)
			if err != nil {
				unilog.Logger().Error("unable to write post media", zap.String("path", imgp), zap.Error(err))
				continue
			}
		}
	}
}
