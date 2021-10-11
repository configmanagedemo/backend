package main

import (
	"net/http"
	"os"
	"strconv"

	config "main/config"
	_ "main/internal/pkg/cache"
	"main/internal/pkg/logger"
	model "main/internal/pkg/web/model"
	"main/internal/pkg/web/router"
)

func main() {
	// first init
	if len(os.Args) == 3 && os.Args[2] == "init" {
		model.FirstInitData()
	}

	addr := config.Conf.Svr.IP + ":" + strconv.Itoa(config.Conf.Svr.Port)
	logger.INFO("server listen on " + addr)

	r := router.InitRouter()
	s := &http.Server{
		Addr:    addr,
		Handler: r,
	}
	err := s.ListenAndServe()
	if err != nil {
		logger.ERROR(err.Error())
	}
}
