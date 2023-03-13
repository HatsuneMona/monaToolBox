package service

import (
	"encoding/json"
	"errors"
	"go.uber.org/zap"
	"monaToolBox/global"
	"monaToolBox/global/cache"
	"monaToolBox/models"
)

type tinyUrlService struct{}

var TinyUrlService = new(tinyUrlService)

const (
	TinyurlCachePrefix = "tinyUrlCache_"
)

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

	// 添加以后，删除缓存
	var routes []string
	for _, url := range tinyUrls {
		routes = append(routes, url.TinyUrl)
	}
	if err := s.RefreshCache(routes, true); err != nil {
		global.Log.Warn("refresh tiny url Cache err!", zap.Error(err))
	}

	return nil, tinyUrls
}

// GetByTinyRouteList 通过 短链的 route 查询这个 短链 的信息
func (s *tinyUrlService) GetByTinyRouteList(tinyRouteList []string, passCache bool) (err error, tinyUrlInfoList []models.TinyUrl) {
	if len(tinyRouteList) < 1 {
		return errors.New("input empty list"), []models.TinyUrl{}
	}

	// 先取 cache
	if !passCache {
		if err, cacheInfos, emptyList := s.getTinyRouteListByCache(tinyRouteList); err != nil {
			global.Log.Warn("getTinyRouteListByCache error: ", zap.Error(err))
		} else {
			tinyUrlInfoList = cacheInfos
			tinyRouteList = emptyList // 没打中缓存的，需要重新查
		}
	}

	// 没有打中缓存的情况
	if len(tinyRouteList) > 0 {
		var dbInfos []models.TinyUrl
		err = global.DB.Where("tiny_url IN ?", tinyRouteList).Find(&dbInfos).Error
		if err != nil {
			return
		}
		tinyUrlInfoList = append(tinyUrlInfoList, dbInfos...)
		// 写入缓存
		_ = s.addTinyRoutesToCache(dbInfos)
	}

	return
}

// RefreshCache 刷新缓存，可选择是否重新更新缓存，或添加最新数据
func (s *tinyUrlService) RefreshCache(tinyRoutes []string, deleteOnly bool) error {
	if len(tinyRoutes) < 1 {
		return errors.New("input empty list")
	}

	if err := s.deleteTinyRoutesCache(tinyRoutes); err != nil {
		global.Log.Warn("refresh Cache error! Delete cache error.", zap.Error(err))
		return err
	}

	if !deleteOnly {
		if err, _ := s.GetByTinyRouteList(tinyRoutes, true); err != nil {
			global.Log.Warn("refresh Cache error! Add cache again error.", zap.Error(err))
			return err
		}
	}

	return nil
}

// getTinyRouteListByCache 从Cache里找到缓存了的 短链 信息
func (s *tinyUrlService) getTinyRouteListByCache(tinyRouteList []string) (err error, tinyUrlInfoList []models.TinyUrl, missList []string) {

	if len(tinyRouteList) < 1 {
		return errors.New("input empty list"), tinyUrlInfoList, missList
	}

	mc := cache.NewMemcachedClient(TinyurlCachePrefix)
	result, err := mc.GetMulti(tinyRouteList)
	if err != nil {
		return err, nil, nil
	}

	// 缓存里查到信息了
	if len(result) > 0 {
		var Info models.TinyUrl
		for _, routeStr := range tinyRouteList {
			if jsonInfo, ok := result[routeStr]; ok {
				// found cache
				global.Log.Debug("catch cache!", zap.String("memCacheKey", routeStr))
				if err := json.Unmarshal(jsonInfo, &Info); err != nil {
					global.Log.Warn("unmarshal memcached result error! ",
						zap.String("memCacheKey", routeStr),
						zap.String("jsonInfo", string(jsonInfo)),
					)
					// 这条记录不要了
					missList = append(missList, routeStr)
					continue
				}
				tinyUrlInfoList = append(tinyUrlInfoList, Info)
			} else {
				// cache miss
				missList = append(missList, routeStr)
			}
		}
	} else {
		// 一条信息都没找到
		missList = tinyRouteList
	}
	return
}

// addTinyRoutesToCache 添加缓存
func (s *tinyUrlService) addTinyRoutesToCache(tinyRoutes []models.TinyUrl) error {
	mc := cache.NewMemcachedClient(TinyurlCachePrefix)

	// TODO change to multi func
	for _, tinyRoute := range tinyRoutes {
		if err := mc.Set(tinyRoute.TinyUrl, tinyRoute, 3600); err != nil {
			global.Log.Warn("tinyRoute add memcached error!",
				zap.String("key", tinyRoute.TinyUrl),
				zap.Any("value(struct)", tinyRoute),
			)
		}
	}

	return nil
}

// deleteTinyRoutesCache 删除缓存
func (s *tinyUrlService) deleteTinyRoutesCache(tinyRoutes []string) error {
	mc := cache.NewMemcachedClient(TinyurlCachePrefix)

	// TODO change to multi func
	for _, tinyRoute := range tinyRoutes {
		if err := mc.Delete(tinyRoute); err != nil {
			global.Log.Warn("tinyRoute delete memcached error!",
				zap.String("key", tinyRoute),
			)
		}
	}

	return nil
}
