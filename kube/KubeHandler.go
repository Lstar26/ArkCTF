package kube

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"gopkg.in/yaml.v2"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// PodInfo 用于存储Pod的信息，以便于转换成JSON
type PodInfo struct {
	Name       string `json:"name"`
	Status     string `json:"status"`
	CreateTime string `json:"createTime"`
}

// initClient 用于初始化Kubernetes客户端
func initClient() (*kubernetes.Clientset, error) {
	// 配置Kubernetes客户端
	config, err := clientcmd.BuildConfigFromFlags("", "/Users/starl/GolandProjects/ArkCTF/kube/kube-conifg.yaml")
	if err != nil {
		log.Printf("无法构建Kubernetes配置：%v", err)
		return nil, err
	}
	// 创建Kubernetes核心客户端
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Printf("无法创建Kubernetes客户端：%v", err)
		return nil, err
	}
	return clientset, nil
}

func GetPods() ([]PodInfo, error) {
	clientset, err := initClient()
	if err != nil {
		log.Printf("初始化Kubernetes客户端失败：%v", err)
		return nil, err
	}

	// 获取Pod列表
	pods, err := clientset.CoreV1().Pods("default").List(context.TODO(), v1.ListOptions{})
	if err != nil {
		log.Printf("无法获取Pod列表：%v", err)
		return nil, err
	}
	// 解析并处理获取到的Pod信息
	var podInfos []PodInfo
	for _, pod := range pods.Items {
		podInfos = append(podInfos, PodInfo{
			Name:       pod.Name,
			Status:     string(pod.Status.Phase),
			CreateTime: pod.CreationTimestamp.String(),
		})
	}

	// 如果podInfos为空，则打印相应的日志
	if len(podInfos) == 0 {
		log.Println("未查询到任何Pod信息")
	} else {
		log.Printf("查询到%d个Pod信息", len(podInfos))
	}

	return podInfos, nil
}

func PodsHandler(w http.ResponseWriter, r *http.Request) {
	podInfos, err := GetPods()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// 设置响应的内容类型为JSON
	w.Header().Set("Content-Type", "application/json")
	// 将Pod信息编码为JSON并写入响应
	json.NewEncoder(w).Encode(podInfos)
}

// createPodHandler 处理创建Deployment的HTTP请求
func CreatePodHandler(w http.ResponseWriter, r *http.Request) {
	// 解析HTTP请求体中的YAML数据
	var deploymentYAML string
	if err := json.NewDecoder(r.Body).Decode(&deploymentYAML); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// 将YAML转换为Deployment对象
	deployment := &appsv1.Deployment{}
	if err := yaml.Unmarshal([]byte(deploymentYAML), deployment); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	clientset, err := initClient()
	if err != nil {
		http.Error(w, "Failed to initialize Kubernetes client: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 使用clientset创建Deployment
	_, err = clientset.AppsV1().Deployments(deployment.Namespace).Create(context.TODO(), deployment, v1.CreateOptions{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Deployment created successfully"))
}

// updatePodHandler 处理更新Deployment的HTTP请求
func UpdatePodHandler(w http.ResponseWriter, r *http.Request) {
	// 解析HTTP请求体中的YAML数据
	var deploymentYAML string
	if err := json.NewDecoder(r.Body).Decode(&deploymentYAML); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// 将YAML转换为Deployment对象
	deployment := &appsv1.Deployment{}
	if err := yaml.Unmarshal([]byte(deploymentYAML), deployment); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	clientset, err := initClient()
	if err != nil {
		http.Error(w, "Failed to initialize Kubernetes client: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 获取Deployment的命名空间
	namespace := deployment.Namespace
	if namespace == "" {
		namespace = "default" // 如果没有指定命名空间，则使用默认命名空间
	}

	// 使用clientset修改Deployment
	_, err = clientset.AppsV1().Deployments(namespace).Update(context.TODO(), deployment, v1.UpdateOptions{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Deployment updated successfully"))
}

// deletePodHandler 处理删除Deployment的HTTP请求
func DeletePodHandler(w http.ResponseWriter, r *http.Request) {
	// 解析HTTP请求体中的YAML数据
	var deploymentYAML string
	if err := json.NewDecoder(r.Body).Decode(&deploymentYAML); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// 将YAML转换为Deployment对象
	deployment := &appsv1.Deployment{}
	if err := yaml.Unmarshal([]byte(deploymentYAML), deployment); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	clientset, err := initClient()
	if err != nil {
		http.Error(w, "Failed to initialize Kubernetes client: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 从Deployment对象中获取name和namespace
	name := deployment.Name
	namespace := deployment.Namespace
	if namespace == "" {
		namespace = "default" // 如果没有指定命名空间，则使用默认命名空间
	}

	// 使用clientset删除Deployment
	err = clientset.AppsV1().Deployments(namespace).Delete(context.TODO(), name, v1.DeleteOptions{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Deployment deleted successfully"))
}
