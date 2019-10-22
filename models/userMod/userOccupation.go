package userMod

import "time"

type UserOccupation struct {
	UUID                    string `xorm:"varchar(26) notnull unique pk 'uuid'"`
	UserID                  string `xorm:"notnull pk 'user_id'"`
	JobType                 int    `form:"JobType" json:"JobType"  binding:"required"`
	CertificateCode         string `form:"CertificateCode" json:"CertificateCode"  binding:"required"`
	CertificateOrganization string `form:"CertificateOrganization" json:"CertificateOrganization"  binding:"required"`
	CorkingOrganization     string `form:"CorkingOrganization" json:"CorkingOrganization"  binding:"required"`
	Department              string `form:"Department" json:"Department"  binding:"required"`
	Seniority               int
	Title                   string `form:"Title" json:"Title"  binding:"required"`
	Province                string
	City                    string
	County                  string
	Address                 string
	Introduction            string
	Status                  int
	CreateAt                time.Time
	ReviewTime              time.Time
	Reviewer                string
	ReviewComments          string
	JobSubject              string
}
