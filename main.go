package main

import (
	"embed"
	"encoding/json"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
	"webssh/controller"
	"webssh/core"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

//go:embed web/dist/*
var f embed.FS

var (
	port       = flag.Int("p", 5032, "服务运行端口")
	v          = flag.Bool("v", false, "显示版本号")
	authInfo   = flag.String("a", "", "开启账号密码登录验证, '-a user:pass'的格式传参")
	timeout    int
	savePass   bool
	version    string
	buildDate  string
	goVersion  string
	gitVersion string
	username   string
	password   string
)

func init() {
	flag.IntVar(&timeout, "t", 120, "ssh连接超时时间(min)")
	flag.BoolVar(&savePass, "s", true, "保存ssh密码")
	if envVal, ok := os.LookupEnv("savePass"); ok {
		if b, err := strconv.ParseBool(envVal); err == nil {
			savePass = b
		}
	}
	if envVal, ok := os.LookupEnv("authInfo"); ok {
		*authInfo = envVal
	}
	if envVal, ok := os.LookupEnv("port"); ok {
		if b, err := strconv.Atoi(envVal); err == nil {
			*port = b
		}
	}
	flag.Parse()
	if *v {
		fmt.Printf("Version: %s\n\n", version)
		fmt.Printf("BuildDate: %s\n\n", buildDate)
		fmt.Printf("GoVersion: %s\n\n", goVersion)
		fmt.Printf("GitVersion: %s\n\n", gitVersion)
		os.Exit(0)
	}
	
	// 首先检查命令行参数
	if *authInfo != "" {
		accountInfo := strings.Split(*authInfo, ":")
		if len(accountInfo) != 2 || accountInfo[0] == "" || accountInfo[1] == "" {
			fmt.Println("请按'user:pass'的格式来传参或设置环境变量, 且账号密码都不能为空!")
			os.Exit(0)
		}
		username, password = accountInfo[0], accountInfo[1]
	} else {
		// 如果命令行参数没有设置，则检查配置文件
		authConfig := core.GetAuthConfig()
		if authConfig.Enable && authConfig.Username != "" && authConfig.Password != "" {
			username = authConfig.Username
			password = authConfig.Password
			log.Println("Using authentication settings from config file")
		}
	}
}

func staticRouter(router *gin.Engine) {
	if password != "" {
		accountList := map[string]string{
			username: password,
		}
		authorized := router.Group("/", gin.BasicAuth(accountList))
		authorized.GET("", func(c *gin.Context) {
			indexHTML, _ := f.ReadFile("web/dist/" + "index.html")
			c.Writer.Write(indexHTML)
		})
	} else {
		router.GET("/", func(c *gin.Context) {
			indexHTML, _ := f.ReadFile("web/dist/" + "index.html")
			c.Writer.Write(indexHTML)
		})
	}
	staticFs, _ := fs.Sub(f, "web/dist/static")
	router.StaticFS("/static", http.FS(staticFs))
}

func getDefaultConfig(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	// 获取默认配置
	defaultConfig := core.GetDefaultSSHConfig()
	
	// 转换为JSON
	jsonData, err := json.Marshal(defaultConfig)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	// 返回JSON数据
	w.Write(jsonData)
}

func main() {
	server := gin.Default()
	server.SetTrustedProxies(nil)
	server.Use(gzip.Gzip(gzip.DefaultCompression))
	staticRouter(server)
	
	// 添加获取默认配置的API路由
	server.GET("/api/default-config", func(c *gin.Context) {
		// 获取默认配置
		defaultConfig := core.GetDefaultSSHConfig()
		
		// 添加日志
		log.Printf("Returning default config: %+v\n", defaultConfig)
		
		// 返回JSON数据
		c.JSON(200, defaultConfig)
	})
	
	// 添加获取服务器列表的API路由
	server.GET("/api/server-list", func(c *gin.Context) {
		// 获取服务器列表
		serverList := core.GetServerList()
		
		// 返回JSON数据
		c.JSON(200, serverList)
	})
	
	server.GET("/term", func(c *gin.Context) {
		controller.TermWs(c, time.Duration(timeout)*time.Minute)
	})
	server.GET("/check", func(c *gin.Context) {
		responseBody := controller.CheckSSH(c)
		responseBody.Data = map[string]interface{}{
			"savePass": savePass,
		}
		c.JSON(200, responseBody)
	})
	file := server.Group("/file")
	{
		file.GET("/list", func(c *gin.Context) {
			c.JSON(200, controller.FileList(c))
		})
		file.GET("/download", func(c *gin.Context) {
			controller.DownloadFile(c)
		})
		file.POST("/upload", func(c *gin.Context) {
			c.JSON(200, controller.UploadFile(c))
		})
		file.GET("/progress", func(c *gin.Context) {
			controller.UploadProgressWs(c)
		})
	}
	server.Run(fmt.Sprintf(":%d", *port))
}
