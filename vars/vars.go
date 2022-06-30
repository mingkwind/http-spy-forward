package vars

import (
	"http-spy-forward/models"
)

var (
	// 50是线程数
	DataChan = make(chan *models.HttpReq, 50)
	// 目标转发url
	TargetURL string
)
