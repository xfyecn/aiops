package controller

import (
	"encoding/json"
	"fmt"
	"github.com/golang/glog"
	"github.com/yarntime/aiops/pkg/mysql"
	v1 "github.com/yarntime/aiops/pkg/types"
	"net/http"
)

type Controller struct {
	dbWorker      *mysql.Worker
	jobController *JobController
	customConfig  v1.CustomConfig
	appConfig     v1.ApplicationConfig
}

func NewController(c *v1.Config) *Controller {
	return &Controller{
		dbWorker:      mysql.NewDBWorker(c),
		jobController: NewJobController(c),
		customConfig:  c.CustomCfg,
		appConfig:     c.AppCfg,
	}
}

func (c *Controller) Create(w http.ResponseWriter, req *http.Request) {
	glog.V(3).Info("scan the monitor objects, create the cron jobs")

	monitorObjects := c.dbWorker.List()
	for _, monitorObject := range monitorObjects {
		fmt.Printf("%v\n", monitorObject)
		for _, appConf := range c.appConfig.App {
			if monitorObject.MonitorTypes&appConf.Id != 0 {
				c.jobController.CreateTrainingJob(monitorObject, c.customConfig, appConf)
			}
		}
	}

	res := &v1.ApiResponse{
		Code:    200,
		Message: "Successful to scan monitor objects.",
	}

	r, _ := json.Marshal(res)
	w.Write(r)
}