package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/urfave/cli"
)

//Config 配置
type Config struct {
	Images []Image
}

//Image 镜像
type Image struct {
	Respository string
	Tag         string
}

func parseConf(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var conf Config

	return &conf, json.NewDecoder(f).Decode(&conf)
}

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

	//删除
	cmd = fmt.Sprintf("%s %s %s", cmdPath, "rmi", turned)
	_, err = execShellCmd(cmd)
	if err != nil {
		return err
	}

	return nil
}

var (
	registryList = map[string]func(string) string{"ali": pathToPathAli, "anjia": pathToPathAnJia}
)

func pullSomething(whichImage string) bool {
	// var status = map[string]string{
	// 	"126": "命令不可执行",
	// 	"127": "没找到命令",
	// 	"128": "无效退出参数",
	// 	"130": "(Ctrl+C)终止",
	// }

	//转换路径，下载
	var err error
	var turned string
	result := true
	for k, v := range registryList {
		turned = v(whichImage)
		err = pullAndTurnTag(turned, whichImage)
		if err == nil {
			fmt.Printf("From registry %s Get %s\n", k, whichImage)
			result = false
			break
		}
	}
	return result
	//fmt.Println(result, status)
}

/*
docker pull gcr.azk8s.cn/google-containers/federation-controller-manager-arm64:v1.3.1-beta.1
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
		turned = "gcr.azk8s.cn" + "/" + pathSplit[1] + "/" + pathSplit[2]
	case "k8s.gcr.io":
		turned = "gcr.azk8s.cn" + "/" + "google-containers" + "/" + pathSplit[1]
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
	case "k8s.gcr.io":
		turned = "registry.cn-hangzhou.aliyuncs.com" + "/" + "google-containers" + "/" + pathSplit[1]
	default:
		turned = "registry.cn-hangzhou.aliyuncs.com" + "/" + "google-containers" + "/" + pathSplit[len(pathSplit)-1]
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
	var configPath string
	var config *Config
	app := cli.NewApp()

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "c, config",
			Value:       "",
			Usage:       "waiting pull image list",
			Destination: &configPath,
		},
	}

	app.Commands = []cli.Command{
		{
			Name:    "pull",
			Aliases: []string{"p"},
			Usage:   "pull image and turn tag",
			Action: func(c *cli.Context) error {
				fmt.Println("pull image and turn tag")
				return nil
			},
		},
		{
			Name:    "search",
			Aliases: []string{"s"},
			Usage:   "search image",
			Action: func(c *cli.Context) error {
				fmt.Println("search image")
				return nil
			},
		},
	}

	app.Action = func(c *cli.Context) error {
		var err error
		if configPath == "" {
			fmt.Println("please config image list")
		}
		config, err = parseConf(configPath)
		// type Image struct {
		// 	Respository string
		// 	Tag         string
		// }
		var success []string
		var failure []string

		for k := range config.Images {
			if config.Images[k].Tag == "" {
				config.Images[k].Tag = "latest"
			}
			image := config.Images[k].Respository + ":" + config.Images[k].Tag
			result := pullSomething(image)
			if result {
				success = append(success, image)
			} else {
				failure = append(failure, image)
			}
		}
		fmt.Printf("success(%d): %v, failure(%d): %v\n", len(success), success, len(failure), failure)
		return err
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
