package services

import (
	"fmt"
	"io/ioutil"

	"tools/iris/common"

	"github.com/emicklei/go-restful"
	"github.com/kataras/iris/v12"
	"k8s.io/client-go/tools/clientcmd"
)

//var kubeConfigPath = "/Users/sino/.kube/config"
var kubeConfigPath = "/Users/sino/Documents/luckymrwang/kubernetes/.kube/k8s1.config"

type ContainerService struct {
	Ctx iris.Context
}

func GetContainerService(ctx iris.Context) *ContainerService {
	return &ContainerService{Ctx: ctx}
}

func (s *ContainerService) ExecShell(namespace, pod, container string) (string, error) {
	kubeConfig, err := ioutil.ReadFile(kubeConfigPath)
	if err != nil {
		return "", fmt.Errorf("读取 kube config 失败：%s", err.Error())
	}
	restConfig, err := clientcmd.RESTConfigFromKubeConfig(kubeConfig)
	if err != nil {
		return "", fmt.Errorf("获取 k8s config 失败：%s", err.Error())
	}
	k8sClient, err := common.GetK8sClient(string(kubeConfig))
	if err != nil {
		return "", fmt.Errorf("获取 k8s 客户端失败：%s", err.Error())
	}
	request := restful.NewRequest(nil)
	request.SetAttribute("namespace", namespace)
	request.SetAttribute("pod", pod)
	request.SetAttribute("container", container)
	return common.ExecShell(k8sClient, restConfig, request)
}

func (s *ContainerService) CopyFromPod(namespace, pod, container, srcPath string) (string, error) {
	destPath := "./uploads"
	prefix := `require('skyapm-nodejs').start({ serviceName: 'nodejs-demo-code', directServers: '10.48.51.135:21594' });`

	kubeConfig, err := ioutil.ReadFile(kubeConfigPath)
	if err != nil {
		return "", fmt.Errorf("读取 kube config 失败：%s", err.Error())
	}
	restConfig, err := clientcmd.RESTConfigFromKubeConfig(kubeConfig)
	if err != nil {
		return "", fmt.Errorf("获取 k8s config 失败：%s", err.Error())
	}
	k8sClient, err := common.GetK8sClient(string(kubeConfig))
	if err != nil {
		return "", fmt.Errorf("获取 k8s 客户端失败：%s", err.Error())
	}
	err = common.CopyFromPod(k8sClient, restConfig, namespace, pod, container, srcPath, destPath)
	if err != nil {
		return "", err
	}
	fileInfo := common.Pathinfo(srcPath)
	localPath := fmt.Sprintf("%s/%v.%v", destPath, fileInfo["filename"], fileInfo["extension"])
	content := common.ReadFile(localPath)
	mergeContent := prefix + "\n" + content
	err = common.WriteFile(localPath, []byte(mergeContent))
	if err != nil {
		return "", err
	}

	return mergeContent, nil
}

func (s *ContainerService) PublishNodeJS(namespace, pod, container, srcPath string) (string, error) {
	data, err := s.CopyFromPod(namespace, pod, container, srcPath)
	if err != nil {
		return "", err
	}
	deployName := "banana"
	fileInfo := common.Pathinfo(srcPath)
	configmapName := fmt.Sprintf("%s-nodejs", deployName)

	_ = new(ConfigMapService).Delete(namespace, configmapName)
	// create configmap
	_, err = new(ConfigMapService).Add(namespace, configmapName, fmt.Sprintf("%v.%v", fileInfo["filename"], fileInfo["extension"]), data)
	if err != nil {
		return "", err
	}
	// create deployment
	_, err = new(DeploymentService).Update(namespace, deployName, srcPath)
	if err != nil {
		return "", err
	}

	return "", nil
}
