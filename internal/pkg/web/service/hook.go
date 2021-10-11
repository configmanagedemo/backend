package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"main/config"
	"main/internal/pkg/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	ErrWebHookRsp = errors.New("not get a 200 response")
)

func Notify(data *gin.H) error {
	url := config.Conf.Svr.Webhook
	body, err := json.Marshal(data)
	if err != nil {
		logger.ERROR(err.Error())
		return err
	}
	logger.DEBUG(fmt.Sprintf("webhook:%s, data:%s", url, string(body)))
	res, err := http.Post(url, "Content-Type: application/json", bytes.NewBuffer(body)) //nolint: gosec,noctx
	if err != nil {
		logger.ERROR(err.Error())
		return err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		logger.ERROR(fmt.Sprintf("status not ok, code is %d, status is %s", res.StatusCode, res.Status))
		return ErrWebHookRsp
	}

	return nil
}
