package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
	"xframe/cmd/templates"
)

type handleFunc func(args []string) error
type renderFunc func(*appConfInfo, string) error

const (
	craeteCmd          = "create"
	appTypeWeb         = "web"
	deployMethodDocker = "docker"
	templatesPathPre   = "cmd/templates"
)

var (
	appTypeTplMap = map[string]string{
		appTypeWeb: "cmd/templates/web",
	}
	deployTplMap = map[string]renderFunc{
		deployMethodDocker: renderDockerDeploy,
	}
)

type appConfInfo struct {
	AppName         string
	AppType         string
	AppDeployMethod string
	AppCreatePath   string
	BaseImage       string
	BaseImageTag    string
	Ports           []string
}

func NewAppConf(appName string, appType string, deployMethod string) *appConfInfo {
	appConf := appConfInfo{
		AppType:         appType,
		AppDeployMethod: deployMethod,
		AppName:         filepath.Base(appName),
	}
	if filepath.IsAbs(appName) {
		appConf.AppCreatePath = filepath.Dir(appName)
	} else {
		currPath, _ := os.Getwd()
		dir := filepath.Dir(appName)
		appConf.AppCreatePath = filepath.Join(currPath, dir)
	}
	return &appConf
}

func (c *appConfInfo) extendConfInfo() error {
	var err error
	switch c.AppDeployMethod {
	case deployMethodDocker:
		baseImge, _ := getInput("please input docker base image ")
		componts := strings.Split(strings.Trim(baseImge, "\n"), ":")
		c.BaseImage = componts[0]
		if len(componts) < 2 || componts[1] == "" {
			c.BaseImageTag = "latest"
		} else {
			c.BaseImageTag = componts[1]
		}
	}
	return err
}

func (c *appConfInfo) Init() error {
	err := c.extendConfInfo()
	if err != nil {
		return err
	}

	if out, err := initProj(filepath.Join(c.AppCreatePath, c.AppName)); err != nil {
		return err
	} else {
		fmt.Println(string(out))
	}
	return c.RenderApp()
}

func (c *appConfInfo) RenderApp() error {
	tplPath := templatesPathPre + "/" + "app.yml.tpl"
	tp := template.New("app.yml")
	content, err := templates.Asset(tplPath)
	if err != nil {
		return err
	}
	tp, err = tp.Parse(string(content))
	if err != nil {
		return err
	}
	f, err := os.Create(filepath.Join(c.AppCreatePath, c.AppName, "app.yml"))
	if err != nil {
		return err
	}
	defer f.Close()
	return tp.Execute(f, c)
}

func getInput(prompt string) (string, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(prompt + ":")
	return reader.ReadString('\n')
}

func usage(topic string) {
	switch topic {
	case craeteCmd:
		fmt.Println("Usage: xframe create [-d deploy method] [-t app type]  <app name>\n" +
			"-d flag instructs the way to deploy the app, default is docker, now support:" +
			" docker\n" +
			"-t flag instructs the app type, default is web service, now support: web\n")
	default:
		fmt.Println("Usage: \n" +
			"\t xframe command [arguments]\n" +
			"The commands are:\n" +
			"\t create \t create app use specify template\n")
	}
}

func createApp(args []string) error {
	var (
		err          error
		deployMethod = deployMethodDocker
		appType      = appTypeWeb
		appName      string
	)

	for idx, length := 0, len(args); idx < length; {
		arg := args[idx]
		if arg[0] == '-' {
			if length-idx <= 1 {
				return errors.New("parse args error")
			}
			switch arg {
			case "-d":
				deployMethod = args[idx+1]
			case "-t":
				appType = args[idx+1]
			}
			idx += 2
		} else {
			appName = arg
			idx++
		}
	}
	if _, ok := appTypeTplMap[appType]; !ok {
		return errors.New("app type not support")
	}
	if _, ok := deployTplMap[deployMethod]; !ok {
		return errors.New("deploy method not support")
	}

	appConf := NewAppConf(appName, appType, deployMethod)
	err = appConf.Init()
	if err != nil {
		return err
	}

	if err := createProjectTemplate(appConf); err != nil {
		return err
	}

	return renderDeployTemplate(appConf)
}

func initProj(appName string) ([]byte, error) {
	fmt.Println("init project using git")
	cmd := exec.Command("git", "init", appName)
	return cmd.CombinedOutput()
}

func createProjectTemplate(appConf *appConfInfo) error {
	tplPath := appTypeTplMap[appConf.AppType]
	return copyFilesFromAsset(tplPath, filepath.Join(appConf.AppCreatePath, appConf.AppName))
}

func copyFilesFromAsset(assetPath string, dst string) error {
	var err error
	fileNames, err := templates.AssetDir(assetPath)
	if err != nil {
		return err
	}
	for _, fileName := range fileNames {
		tsrc := assetPath + "/" + fileName
		tdst := filepath.Join(dst, fileName)

		_, err := templates.AssetInfo(tsrc)
		if err != nil {
			if _, err := os.Stat(tdst); os.IsNotExist(err) {
				os.Mkdir(tdst, 0755)
			}
			err = copyFilesFromAsset(tsrc, tdst)
		} else {
			dstFile, err := os.Create(tdst)
			if err != nil {
				return err
			}
			content, err := templates.Asset(tsrc)
			if err != nil {
				return err
			}
			_, err = dstFile.Write(content)
			dstFile.Close()
		}
	}
	return err
}

// render path: "d-"+AppDeployMethod
func renderDeployTemplate(appConf *appConfInfo) error {
	handler, _ := deployTplMap[appConf.AppDeployMethod]
	return handler(appConf, templatesPathPre+"/d-"+appConf.AppDeployMethod)
}

func cmdFactory(cmd string) (handleFunc, bool) {
	handlerMap := map[string]handleFunc{
		"create": createApp,
	}
	if handler, ok := handlerMap[cmd]; ok {
		return handler, ok
	}
	return nil, false
}

func renderDockerDeploy(appConf *appConfInfo, path string) error {
	fileNames, err := templates.AssetDir(path)
	for _, filename := range fileNames {
		fn := strings.Split(filename, ".tpl")[0]
		dstf := filepath.Join(appConf.AppCreatePath, appConf.AppName, fn)
		f, err := os.Create(dstf)
		content, err := templates.Asset(path + "/" + filename)
		if err != nil {
			return err
		}

		tp := template.New(filename)
		tp, err = tp.Parse(string(content))
		if err != nil {
			f.Close()
			return err
		}
		err = tp.Execute(f, appConf)
		if err != nil {
			f.Close()
			return err
		}
		f.Close()
	}
	return err
}

func main() {
	if len(os.Args) < 3 {
		usage("")
		return
	}
	command := os.Args[1]
	args := os.Args[2:]
	handler, ok := cmdFactory(command)
	if !ok {
		fmt.Println("no such command")
		usage("")
	}
	err := handler(args)
	if err != nil {
		fmt.Println(err, "\n\r")
		usage(command)
	}
}
