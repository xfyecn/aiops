package main

import (
	"encoding/json"
	"flag"
	"github.com/golang/glog"
	"github.com/gorilla/mux"
	"github.com/yarntime/aiops/pkg/controller"
	v1 "github.com/yarntime/aiops/pkg/types"
	io "io/ioutil"
	"net/http"
	"reflect"
)

var (
	apiserverAddress  string
	globalConfig      string
	applicationConfig string
)

func init() {
	flag.StringVar(&apiserverAddress, "apiserver_address", "", "Kubernetes apiserver address")
	flag.StringVar(&globalConfig, "global_config_file", "/etc/aiops/config.json", "global config file")
	flag.StringVar(&applicationConfig, "applicationConfig", "/etc/aiops/application.json", "application config file")
	flag.Set("alsologtostderr", "true")
	flag.Parse()
}

func main() {
	customConfig := v1.CustomConfig{}
	err := LoadConfig(globalConfig, &customConfig)
	if err != nil {
		glog.Fatal("Failed to load custom config.")
	}

	appConfig := v1.ApplicationConfig{}
	err = LoadConfig(applicationConfig, &appConfig)
	if err != nil {
		glog.Fatal("Failed to load application config." + err.Error())
	}

	initAppConfig(customConfig, appConfig)

	config := &v1.Config{
		CustomCfg: customConfig,
		AppCfg:    appConfig,
		Host:      apiserverAddress,
	}

	c := controller.NewController(config)

	router := mux.NewRouter()
	router.HandleFunc("/health", health).Methods("GET")
	router.HandleFunc("/create", c.Create).Methods("GET")
	router.HandleFunc("/delete", c.Delete).Methods("GET")

	glog.Info("http server started.")
	glog.Fatal(http.ListenAndServe(":8080", router))
}

func LoadConfig(filename string, v interface{}) error {
	data, err := io.ReadFile(filename)
	if err != nil {
		return err
	}

	dataJson := []byte(data)
	err = json.Unmarshal(dataJson, v)
	if err != nil {
		return err
	}

	return nil
}

func initAppConfig(customConfig v1.CustomConfig, appConfig v1.ApplicationConfig) {
	globalFiled := reflect.TypeOf(customConfig.Global)
	globalValue := reflect.ValueOf(customConfig.Global)
	params := []string{}
	for i := 0; i < globalFiled.NumField(); i++ {
		f := globalFiled.Field(i)
		_, skip := f.Tag.Lookup("skip")
		if skip {
			continue
		}
		params = append(params, "--"+f.Tag.Get("json")+"="+globalValue.Field(i).Interface().(string))
	}
	for i := 0; i < len(appConfig.App); i++ {
		appConfig.App[i].Params = append(appConfig.App[i].Params, params...)
	}
}

func health(w http.ResponseWriter, _ *http.Request) {
	w.Write([]byte("ok."))
}
