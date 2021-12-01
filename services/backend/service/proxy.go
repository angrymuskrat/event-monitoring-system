package service

import (
	"fmt"
	"github.com/visheratin/unilog"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

func uploadImage(req InstaImageRequest) (image []byte, status int) {
	image, isExist := loadImageFromStorage(req.Shortcode)
	if isExist {
		return image, http.StatusOK
	}
	urlTemplate := "https://www.instagram.com/p/%v/media/?size=m"
	url := fmt.Sprintf(urlTemplate, req.Shortcode)
	resp, err := http.Get(url)
	if err != nil {
		unilog.Logger().Error("unable to get image from Instagram", zap.Error(err))
		return nil, http.StatusInternalServerError
	}
	image, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		unilog.Logger().Error("unable to read image from response", zap.Error(err))
		return nil, http.StatusInternalServerError
	}
	statusStr := resp.Status[:3]
	switch statusStr {
	case "200":
		saveImageToStorage(image, req.Shortcode)
		return image, http.StatusOK
	case "404":
		return nil, http.StatusNotFound
	default:
		unilog.Logger().Error("insta block request", zap.String("status", resp.Status),
			zap.String("code", req.Shortcode), zap.Error(err))
		return nil, http.StatusInternalServerError
	}
}

func saveImageToStorage(image []byte, code string) {
	subDir := code[:3]
	path := filepath.Join(storagePath, subDir)
	fileName := fmt.Sprintf("%v.jpg", code)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.Mkdir(path, 0777)
		if err != nil {
			unilog.Logger().Error("unable to open file system", zap.Error(err))
			return
		}
	}
	err := ioutil.WriteFile(filepath.Join(path, fileName), image, 0777)
	if err != nil {
		unilog.Logger().Error("unable to write image file", zap.String("shortcode", code), zap.Error(err))
		return
	}
}

func loadImageFromStorage(code string) (image []byte, isExist bool) {
	subDir := code[:3]
	path := filepath.Join(storagePath, subDir)
	fileName := fmt.Sprintf("%v.jpg", code)

	image, err := ioutil.ReadFile(filepath.Join(path, fileName))
	if err != nil {
		return nil, false
	}
	if err != nil {
		unilog.Logger().Error("loadImage: unable to load image to file", zap.Error(err))
		return nil, false
	}
	return image, true
}
