package util

import (
	"github.com/astaxie/beego/logs"
	"net/http"
)

type HealthCheck struct {
}

func (HealthCheck) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("system is running normally!"))
	w.WriteHeader(http.StatusOK)
	return
}

func HealthCheckInit(address string)  {
	go func() {
		err := http.ListenAndServe(address, &HealthCheck{})
		if err != nil {
			panic(err.Error())
		}
	}()
	logs.Info("health check start success", address)
}