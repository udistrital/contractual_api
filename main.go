package main

import (
	"fmt"
	"time"

	"github.com/astaxie/beego"

	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/plugins/cors"
	_ "github.com/lib/pq"
	_ "github.com/udistrital/administrativa_crud_api/routers"
	apistatus "github.com/udistrital/utils_oas/apiStatusLib"
	"github.com/udistrital/utils_oas/customerror"
)

func init() {
	orm.Debug = true
	// orm.DefaultTimeLoc = time.UTC
	loc, err := time.LoadLocation("America/Bogota")
	if err != nil {
		fmt.Println(err)
	}
	orm.DefaultTimeLoc = loc
	q := "postgres://" + beego.AppConfig.String("PGuser") + ":" + beego.AppConfig.String("PGpass") + "@" + beego.AppConfig.String("PGurls") + ":5431/" + beego.AppConfig.String("PGdb") + "?sslmode=disable&search_path=" + beego.AppConfig.String("PGschemas") + "&timezone=UTC"
	//fmt.Println(q)
	if err := orm.RegisterDataBase("default", "postgres", q); err != nil {
		panic(err) //Nunca deberia pasar si están bien descargados los paquetes del repo
	}
}

func main() {
	orm.Debug = true
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}

	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"PUT", "PATCH", "GET", "POST", "OPTIONS", "DELETE"},
		AllowHeaders: []string{"Origin", "x-requested-with",
			"content-type",
			"accept",
			"origin",
			"authorization",
			"x-csrftoken"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Custom JSON error pages
	beego.ErrorHandler("400", badRequestJSONPage)
	beego.ErrorHandler("403", forgivenJSONPage)
	beego.ErrorHandler("404", notFoundJSONPage)

	if err := logs.SetLogger(logs.AdapterFile, `{"filename":"/var/log/beego/administrativa_crud_api.log"}`); err != nil {
		beego.Info(err)
	}
	beego.ErrorController(&customerror.CustomErrorController{})
	apistatus.Init()
	beego.Run()
}
