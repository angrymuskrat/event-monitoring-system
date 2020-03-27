package service

import (
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net/http"
)

const (
	cookieName    = "session"
	adminUser     = "admin"
	adminPassword = "admin!pwd"
)

type AuthManager struct {
	store  *sessions.CookieStore
	logger *zap.Logger
}

func newAuthManager(key string, logPath string) (*AuthManager, error) {
	m := AuthManager{}
	if len(key) < 32 {
		return nil, errors.New("key must be at least 32 symbols long")
	}
	m.store = sessions.NewCookieStore([]byte(key))
	if m.store == nil {
		return nil, errors.New("unable to initialize cookies store")
	}
	cfg := zap.Config{
		Encoding:         "console",
		Level:            zap.NewAtomicLevelAt(zapcore.InfoLevel),
		OutputPaths:      []string{logPath},
		ErrorOutputPaths: []string{logPath},
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:  "message",
			LevelKey:    "level",
			EncodeLevel: zapcore.CapitalLevelEncoder,
			TimeKey:     "time",
			EncodeTime:  zapcore.ISO8601TimeEncoder,
		},
	}
	log, err := cfg.Build()
	if err != nil {
		return nil, err
	}
	m.logger = log
	return &m, nil
}

func (m *AuthManager) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.RequestURI == "/login" {
			next.ServeHTTP(w, r)
		}
		sess, err := m.store.Get(r, cookieName)
		if err != nil {
			m.logger.Error("unable to get session from the store", zap.Error(err))
			http.Error(w, "server error", http.StatusInternalServerError)
			return
		}
		if ar, ok := sess.Values["auth"]; !ok {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		} else {
			auth, ok := ar.(bool)
			if !ok || !auth {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		}
	})
}

func (m *AuthManager) login(w http.ResponseWriter, r *http.Request) {
	req, err := m.decodeLoginRequest(r)
	if err != nil {
		m.logger.Error("unable to decode request", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if req.Login != adminUser || req.Password != adminPassword {
		http.Error(w, "login/password incorrect", http.StatusUnauthorized)
		return
	}
	sess, err := m.store.Get(r, cookieName)
	if err != nil {
		m.logger.Error("unable to get session from the store", zap.Error(err))
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}
	uid := uuid.New().String()
	id := uid
	sess.Values["id"] = id
	sess.Values["auth"] = true
	sess.Save(r, w)
}

func (m *AuthManager) decodeLoginRequest(r *http.Request) (LoginRequest, error) {
	defer r.Body.Close()
	var req LoginRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		m.logger.Error("unable to decode request", zap.Error(err))
	}
	return req, err
}
