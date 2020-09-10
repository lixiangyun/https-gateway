package data

import "github.com/lixiangyun/https-gateway/etcdsdk"

var DBAPI etcdsdk.DBAPI
var KvAPI etcdsdk.BaseAPI

func DataInit(endpoints []string) error {
	baseAPI, err := etcdsdk.ClientInit(2, endpoints )
	if err != nil {
		return err
	}
	DBAPI, err = etcdsdk.NewDBInit(baseAPI, "https-gateway/db")
	if err != nil {
		return err
	}
	KvAPI = baseAPI
	return nil
}
