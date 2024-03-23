package admin

import (
	"database/sql"
	"encoding/json"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

// AdminLoginHandler 处理管理员登录请求
func AdminLoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "仅支持POST方法", http.StatusMethodNotAllowed)
		return
	}

	// 解析表单数据
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "解析表单数据出错", http.StatusBadRequest)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	// 连接数据库
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/arkctf")
	if err != nil {
		http.Error(w, "数据库连接错误", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// 检查是否存在具有管理员角色的用户
	var dbUsername, dbPassword, dbEmail string
	err = db.QueryRow("SELECT username, password, email FROM Users WHERE username=? AND role=1", username).Scan(&dbUsername, &dbPassword, &dbEmail)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "凭证无效或不是管理员", http.StatusUnauthorized)
		} else {
			http.Error(w, "查询错误", http.StatusInternalServerError)
		}
		return
	}

	// 检查密码是否匹配
	if password != dbPassword {
		http.Error(w, "凭证无效", http.StatusUnauthorized)
		return
	}

	sessionID, err := GenerateSessionID()
	if err != nil {
		http.Error(w, "生成session_id失败", http.StatusInternalServerError)
		return
	}

	// 设置登录成功的cookie
	http.SetCookie(w, &http.Cookie{
		Name:  "session_id",
		Value: sessionID,
		Path:  "/",
		// 其他cookie选项...
	})

	// 登录成功
	// 定义一个结构体用于响应
	type loginResponse struct {
		Status    string `json:"status"`
		SessionID string `json:"session_id"`
		Data      struct {
			Username string `json:"username"`
			Email    string `json:"email"`
		} `json:"data"`
	}

	respData := loginResponse{
		Status:    "OK",
		SessionID: sessionID,
		Data: struct {
			Username string `json:"username"`
			Email    string `json:"email"`
		}{
			Username: dbUsername,
			Email:    dbEmail,
		},
	}

	respJSON, err := json.Marshal(respData)
	if err != nil {
		http.Error(w, "生成响应失败", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*") //允许跨域访问
	w.WriteHeader(http.StatusOK)
	w.Write(respJSON)

	// 重定向到管理员控制面板
	http.Redirect(w, r, "/admin/dashboard", http.StatusSeeOther)
}
