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
	result, err := exec.Command("/bin/sh", "-c", cmd).Output()
	return string(result), err
}

func doSomething(whichImage string) {
	// var status = map[string]string{
	// 	"126": "命令不可执行",
	// 	"127": "没找到命令",
	// 	"128": "无效退出参数",
	// 	"130": "(Ctrl+C)终止",
	// }

	cmdPath, err := exec.LookPath("docker")
	if err != nil {
		log.Fatal("docker is no install\n")
	}

	//转换路径，下载
	turned := pathToPath(whichImage)
	cmd := fmt.Sprintf("%s %s %s", cmdPath, "pull", turned)
	_, err = execShellCmd(cmd)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
	}

	//tag改名
	cmd = fmt.Sprintf("%s %s %s %s", cmdPath, "tag", turned, whichImage)
	_, err = execShellCmd(cmd)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
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
func pathToPath(path string) string {
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
