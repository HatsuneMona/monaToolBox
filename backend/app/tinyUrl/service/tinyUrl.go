package service

import (
	"monaToolBox/global"
	"monaToolBox/models"
)

type tinyUrlService struct{}

var TinyUrlService = new(tinyUrlService)

func (s *tinyUrlService) GetById(tinyId int) (err error, tinyUrl models.TinyUrl) {
	err = global.DB.Where(" id = ? ", tinyId).First(&tinyUrl).Error

	return
}
