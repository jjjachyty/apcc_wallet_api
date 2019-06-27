package commonSrv

import (
	"apcc_wallet_api/models"
	"time"
)

type Version struct {
	VersionCode        string
	ReleaseTime        time.Time
	ReleaseNote        string
	AndroidDownloadUrl string
	IosDownloadUrl     string
}

func GetLastVersion() (*Version, error) {
	var version = new(Version)

	err := models.SQLBean(version, "select * from version where version_code = (select max(version_code) from version)")
	return version, err
}
