package commonSrv

import (
	"bytes"
	"context"

	"github.com/qiniu/api.v7/auth/qbox"
	"github.com/qiniu/api.v7/storage"
)

var accessKey = "TvtvJ6rK7JfRmuDQvRfsZZ9zm1JvD_MotrxgT4aN"
var secretKey = "noyP3xm9Ibj_Al1CqCHdS7H6aWP4lsg15ps2IZXy"
var mac = qbox.NewMac(accessKey, secretKey)

var cfg = storage.Config{Zone: &storage.ZoneHuanan, UseHTTPS: false, UseCdnDomains: false}

var formUploader = storage.NewFormUploader(&cfg)

func UploadImage(key string, data []byte) error {
	var putPolicy = storage.PutPolicy{
		Scope: "wallet:" + key,
	}
	var upToken = putPolicy.UploadToken(mac)
	return formUploader.Put(context.Background(), &storage.PutRet{}, upToken, key, bytes.NewReader(data), int64(len(data)), &storage.PutExtra{})
}
