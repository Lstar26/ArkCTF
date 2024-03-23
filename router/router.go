package router

import (
	"ArkCTF/admin"
	"ArkCTF/kube"
	"ArkCTF/system"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func InitRoutes() {
	r := mux.NewRouter()

	// 设置跨域访问
	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	r.Use(handlers.CORS(headersOk, originsOk, methodsOk))

	r.HandleFunc("/admin/admin-login", admin.AdminLoginHandler) // 管理员后台登录

	r.HandleFunc("/system/create-user", system.CreateUser) // 创建用户信息
	r.HandleFunc("/system/update-user", system.UpdateUser) // 更新用户信息
	r.HandleFunc("/system/delete-user", system.DeleteUser) // 删除用户信息
	r.HandleFunc("/system/query-user", system.QueryUser)   // 查询单个用户信息
	r.HandleFunc("/system/list-users", system.ListUsers)   // 查询所有用户信息
	r.HandleFunc("/system/sum-users", system.SumUsers)     // 统计注册用户总数

	r.HandleFunc("/kube/pods", kube.PodsHandler)            // 查询Pod信息
	r.HandleFunc("/kube/create-pod", kube.CreatePodHandler) // 创建Deployment
	r.HandleFunc("/kube/update-pod", kube.UpdatePodHandler) // 更新Deployment
	r.HandleFunc("/kube/delete-pod", kube.DeletePodHandler) // 删除Deployment
	fmt.Println("ApiServer is running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", r))
}
