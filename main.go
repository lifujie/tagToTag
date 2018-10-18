package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/urfave/cli"
)

func execShellCmd(cmd string) (string, error) {
	fmt.Printf("cmd: %s\n", cmd)
	result, err := exec.Command("/bin/sh", "-c", cmd).Output()
	return string(result), err
}

func pullAndTurnTag(turned, whichImage string) error {
	cmdPath, err := exec.LookPath("docker")
	if err != nil {
		log.Fatal("docker is no install\n")
	}

	cmd := fmt.Sprintf("%s %s %s", cmdPath, "pull", turned)
	_, err = execShellCmd(cmd)
	if err != nil {
		return err
	}

	//tag改名
	cmd = fmt.Sprintf("%s %s %s %s", cmdPath, "tag", turned, whichImage)
	_, err = execShellCmd(cmd)
	if err != nil {
		return err
	}
	return nil
}

var (
	registryList = map[string]func(string) string{"ali": pathToPathAli, "anjia": pathToPathAnJia}
)

func doSomething(whichImage string) {
	// var status = map[string]string{
	// 	"126": "命令不可执行",
	// 	"127": "没找到命令",
	// 	"128": "无效退出参数",
	// 	"130": "(Ctrl+C)终止",
	// }

	//转换路径，下载
	var err error
	var turned string
	for k, v := range registryList {
		turned = v(whichImage)
		err = pullAndTurnTag(turned, whichImage)
		if err == nil {
			fmt.Printf("From registry %s Get %s\n", k, whichImage)
			break
		}
	}
	//fmt.Println(result, status)
}

/*
docker pull anjia0532/google-containers.federation-controller-manager-arm64:v1.3.1-beta.1
# eq
docker pull gcr.io/google-containers/federation-controller-manager-arm64:v1.3.1-beta.1

# special
# eq
docker pull k8s.gcr.io/federation-controller-manager-arm64:v1.3.1-beta.1
*/
func pathToPathAnJia(path string) string {
	var turned string

	pathSplit := strings.Split(path, "/")

	switch pathSplit[0] {
	case "gcr.io":
		turned = "anjia0532" + "/" + pathSplit[1] + "." + pathSplit[2]
	case "k8s.gcr.io":
		turned = "anjia0532" + "/" + "google-containers" + "." + pathSplit[1]
	}

	return turned
}

// # 下载阿里云镜像
// docker pull registry.cn-hangzhou.aliyuncs.com/google-containers/pause-amd64:3.0

// # 本地命名为 gcr.io/google_containers/pause-amd64:3.0
// docker tag registry.cn-hangzhou.aliyuncs.com/google-containers/pause-amd64:3.0 gcr.io/google_containers/pause-amd64:3.0
//
func pathToPathAli(path string) string {
	var turned string

	pathSplit := strings.Split(path, "/")

	switch pathSplit[0] {
	case "gcr.io":
		turned = "registry.cn-hangzhou.aliyuncs.com" + "/" + pathSplit[1] + "/" + pathSplit[2]
	default:
		fmt.Printf("ali not reachable\n")
	}

	return turned
}

func tagToTag(path string, which int) string {
	var turned string

	pathSplit := strings.Split(path, "/")
	imageNameIndex := strings.Index(pathSplit[1], "google-containers")

	switch which {
	case 0:
		turned = "gcr.io" + "/" + pathSplit[1][0:imageNameIndex] + "/" + pathSplit[1][imageNameIndex:]
	case 1:
		turned = "k8s.gcr.io" + "/" + pathSplit[1][imageNameIndex:]
	}

	return turned
}

func main() {
	app := cli.NewApp()

	app.Action = func(c *cli.Context) error {
		doSomething(c.Args().Get(0))
		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
