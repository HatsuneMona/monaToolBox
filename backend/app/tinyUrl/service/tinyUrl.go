package service

import (
	"errors"
	"monaToolBox/global"
	"monaToolBox/models"
)

type tinyUrlService struct{}

var TinyUrlService = new(tinyUrlService)

func (s *tinyUrlService) GetById(tinyId int) (err error, tinyUrl models.TinyUrl) {
	err = global.DB.Where(" id = ? ", tinyId).First(&tinyUrl).Error
	return
}

// AddNewTinyUrlBatch 批量添加新短链
func (s *tinyUrlService) AddNewTinyUrlBatch(tinyUrls []models.TinyUrl) (error, []models.TinyUrl) {
	if len(tinyUrls) < 1 {
		return errors.New("input empty list"), nil
	}

	// tinyUrl 有唯一索引
	if err := global.DB.Create(&tinyUrls).Error; err != nil {
		return err, nil
	}

	return nil, tinyUrls
}

// GetByTinyRouteList 通过 短链的 route 查询这个 短链 的信息
func (s *tinyUrlService) GetByTinyRouteList(tinyRouteList []string) (err error, tinyUrlInfoList []models.TinyUrl) {
	if len(tinyRouteList) < 1 {
		return errors.New("input empty list"), []models.TinyUrl{}
	}

	err = global.DB.Where("tiny_url IN ?", tinyRouteList).Find(&tinyUrlInfoList).Error
	if err != nil {
		return
	}

	return
}
