package crawler

import (
	"errors"
	"fmt"
	"github.com/visheratin/unilog"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

type client struct {
	mu        sync.Mutex
	cl        http.Client
	cookies   []*http.Cookie
	token     string
	sessionID string
}

func newClient(token string, sessionID string) *client {
	return &client{
		cl: http.Client{
			Timeout: 30 * time.Second,
		},
		cookies:   nil,
		token:     token,
		sessionID: sessionID,
	}
}

func (cl *client) makeRequest(request string) ([]byte, error) {
	cl.mu.Lock()
	defer cl.mu.Unlock()
	req, err := http.NewRequest("GET", request, nil)
	if err != nil {
		unilog.Logger().Error("unable to create request", zap.String("URL", request), zap.Error(err))
		return nil, err
	}
	cookie := fmt.Sprintf("csrftoken=%v; sessionid=%v;", cl.token, cl.sessionID)
	req.Header.Set("cookie", cookie)
	//if cl.cookies == nil {
	//	req.Header.Set("cookie", cl.authCookie)
	//} else {
	//	for _, c := range cl.cookies {
	//		req.AddCookie(c)
	//	}
	//}
	resp, err := cl.cl.Do(req)
	if err != nil {
		unilog.Logger().Error("unable to make request", zap.String("URL", request), zap.Error(err))
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == 404 {
		msg := "entity page was not found"
		unilog.Logger().Error(msg, zap.String("URL", request))
		err = errors.New(msg)
		return nil, err
	}
	if resp.StatusCode == 500 {
		msg := "error during request execution."
		unilog.Logger().Error(msg, zap.String("URL", request))
		err = errors.New(msg)
		return nil, err
	}
	//cl.cookies = resp.Cookies()
	body, err := ioutil.ReadAll(resp.Body)
	time.Sleep(100 * time.Millisecond)
	return body, nil
}
