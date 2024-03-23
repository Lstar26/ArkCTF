package system

import (
	"database/sql"
	"encoding/json"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

// CreateUser 创建用户
func CreateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "仅支持POST方法", http.StatusMethodNotAllowed)
		return
	}

	// 解析请求体中的JSON数据
	var user struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Email    string `json:"email"`
		Role     int    `json:"role"`
	}
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "解析JSON数据出错", http.StatusBadRequest)
		return
	}

	// 连接数据库
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/arkctf")
	if err != nil {
		http.Error(w, "数据库连接错误", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// 插入用户数据
	_, err = db.Exec("INSERT INTO Users(username, password, email, role) VALUES (?, ?, ?, ?)", user.Username, user.Password, user.Email, user.Role)
	if err != nil {
		http.Error(w, "插入用户数据出错", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("用户创建成功"))
}

// DeleteUser 删除用户
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != "DELETE" {
		http.Error(w, "仅支持DELETE方法", http.StatusMethodNotAllowed)
		return
	}

	// 获取URL参数中的id
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "缺少用户ID", http.StatusBadRequest)
		return
	}

	// 连接数据库
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/arkctf")
	if err != nil {
		http.Error(w, "数据库连接错误", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// 删除用户数据
	_, err = db.Exec("DELETE FROM Users WHERE id = ?", id)
	if err != nil {
		http.Error(w, "删除用户数据出错", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("用户删除成功"))
}

// UpdateUser 修改用户信息
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != "PUT" {
		http.Error(w, "仅支持PUT方法", http.StatusMethodNotAllowed)
		return
	}

	// 解析请求体中的JSON数据
	var user struct {
		ID       string `json:"id"`
		Username string `json:"username"`
		Password string `json:"password"`
		Email    string `json:"email"`
		Role     int    `json:"role"`
	}
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "解析JSON数据出错", http.StatusBadRequest)
		return
	}

	// 连接数据库
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/arkctf")
	if err != nil {
		http.Error(w, "数据库连接错误", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// 更新用户数据
	_, err = db.Exec("UPDATE Users SET username = ?, password = ?, email = ?, role = ? WHERE id = ?", user.Username, user.Password, user.Email, user.Role, user.ID)
	if err != nil {
		http.Error(w, "更新用户数据出错", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("用户信息更新成功"))
}

// QueryUser 查询用户信息
func QueryUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "仅支持GET方法", http.StatusMethodNotAllowed)
		return
	}

	// 获取URL参数中的username
	username := r.URL.Query().Get("username")
	if username == "" {
		http.Error(w, "缺少用户名", http.StatusBadRequest)
		return
	}

	// 连接数据库
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/arkctf")
	if err != nil {
		http.Error(w, "数据库连接错误", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// 查询用户数据
	var user struct {
		ID       string `json:"id"`
		Username string `json:"username"`
		Password string `json:"password"`
		Email    string `json:"email"`
		Role     int    `json:"role"`
	}
	err = db.QueryRow("SELECT id, username, password, email, role FROM Users WHERE username = ?", username).Scan(&user.ID, &user.Username, &user.Password, &user.Email, &user.Role)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "用户不存在", http.StatusNotFound)
		} else {
			http.Error(w, "查询用户数据出错", http.StatusInternalServerError)
		}
		return
	}

	// 返回查询结果
	respJSON, err := json.Marshal(user)
	if err != nil {
		http.Error(w, "生成响应失败", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(respJSON)
}

func ListUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "仅支持GET方法", http.StatusMethodNotAllowed)
		return
	}

	// 连接数据库
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/arkctf")
	if err != nil {
		http.Error(w, "数据库连接错误", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// 查询所有用户数据
	rows, err := db.Query("SELECT id, username, password, email, role FROM Users")
	if err != nil {
		http.Error(w, "查询用户数据出错", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// 构建用户列表
	var users []struct {
		ID       string `json:"id"`
		Username string `json:"username"`
		Password string `json:"password"`
		Email    string `json:"email"`
		Role     int    `json:"role"`
	}
	for rows.Next() {
		var user struct {
			ID       string `json:"id"`
			Username string `json:"username"`
			Password string `json:"password"`
			Email    string `json:"email"`
			Role     int    `json:"role"`
		}
		err = rows.Scan(&user.ID, &user.Username, &user.Password, &user.Email, &user.Role)
		if err != nil {
			http.Error(w, "扫描用户数据出错", http.StatusInternalServerError)
			return
		}
		users = append(users, user)
	}

	// 返回查询结果
	respJSON, err := json.Marshal(users)
	if err != nil {
		http.Error(w, "生成响应失败", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(respJSON)
}

func SumUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "仅支持GET方法", http.StatusMethodNotAllowed)
		return
	}

	// 连接数据库
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/arkctf")
	if err != nil {
		http.Error(w, "数据库连接错误", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// 查询注册用户总数
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM Users").Scan(&count)
	if err != nil {
		http.Error(w, "查询用户总数出错", http.StatusInternalServerError)
		return
	}

	// 返回查询结果
	respJSON, err := json.Marshal(map[string]int{"count": count})
	if err != nil {
		http.Error(w, "生成响应失败", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(respJSON)
}
