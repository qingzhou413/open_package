// Copyright 2018. All rights reserved.
//haha

/*
main包是打包程序入口
*/
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"util/file"
	"exception"
	"io/ioutil"
	"strings"
	"github.com/mholt/archiver"
)

var version = flag.String("version", "", "What's the version of js_pen?")

var WorkspaceBe = "d:/work/svn/ideaworkspace/js_open"
var WorkspaceFe = "d:/work/svn/fe/js_open_fe"
var JsOpenDir = "d:/work/projects/js_open"

func main() {
	fmt.Printf("WorkspaceBe: %s\n", WorkspaceBe)
	fmt.Printf("WorkspaceFe: %s\n", WorkspaceFe)
	fmt.Printf("agrs: %s\n", os.Args)
	flag.Parse()
	fmt.Printf("version %s \n", Version())
	fmt.Printf("setup dir %s.\n", SetupName())

	//check command line arg: version
	if Version() == "" {
		log.Fatal("version can not be empty.")
		os.Exit(-1)
	}
	//check fe code
	if !file.FileExist(WorkspaceFe + "/" + Version()) {
		log.Fatalf("js_open_fe of version %s not exists.", Version())
		os.Exit(-1)
	}

	//初始化安装目录
	initDirErr := initDir()
	if initDirErr != nil {
		log.Fatalf("init dir error %s.", initDirErr.Error())
		os.Exit(-1)
	}

	//复制template目录所有文件到setup
	copyTemplateErr := copyTemplate()
	if copyTemplateErr != nil {
		log.Fatalf("copy template files error %s.", copyTemplateErr.Error())
		os.Exit(-1)
	}

	//复制定位引擎相关服务
	copyLocateErr := copyLocateFiles()
	if copyLocateErr != nil {
		log.Fatalf("copy locate files error %s.", copyLocateErr.Error())
		os.Exit(-1)
	}

	//#复制dubbo service服务
	copyDubboServiceErr := copyDubboService()
	if copyDubboServiceErr != nil {
		log.Fatalf("copy dubbo service files error %s.", copyDubboServiceErr.Error())
		os.Exit(-1)
	}

	//#复制dubbo web服务
	copyDubboWebErr := copyDubboWeb()
	if copyDubboWebErr != nil {
		log.Fatalf("copy dubbo web files error %s.", copyDubboWebErr.Error())
		os.Exit(-1)
	}

	zipErr := zip()
	if zipErr != nil {
		log.Fatalf("zipping error %s.", zipErr.Error())
		os.Exit(-1)
	}

}
//打成压缩包
func zip() error {
	setupDir := SetupPath()
	zipPath := ZipPath()
	err := archiver.Zip.Make(zipPath, []string{setupDir})
	return err
}

//复制模板文件夹
func copyTemplate() error {
	setupDir := SetupPath()
	templateDir := TemplateDir()
	err := file.Copy(templateDir, setupDir)
	if err != nil {
		return err
	}
	//修改文件里的安装目录
	fromStr := "setup_dir="
	dstStr := "setup_dir=" + SetupName()
	err = file.Replace(setupDir+"/setup.sh", fromStr, dstStr)
	if err != nil {
		return err
	}
	err = file.Replace(setupDir+"/hot-change-ip.sh", fromStr, dstStr)
	if err != nil {
		return err
	}

	upgrPath := setupDir + "/upgrade"
	_, err = os.Stat(upgrPath)
	if err != nil {
		return err
	}
	shs, err := ioutil.ReadDir(upgrPath)
	if err != nil {
		return err
	}
	for _, sh := range shs {
		if strings.Contains(sh.Name(), ".sh") {
			err = file.Replace(setupDir+"/upgrade/"+sh.Name(), fromStr, dstStr)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

//#复制dubbo web服务
func copyDubboWeb() error {
	buildDir := BuildDir()

	webs := [] string{"api", "building", "datacenter", "device", "rail", "system"}
	for _, web := range webs {
		fmt.Printf("copy web be %s\n", web)
		err := file.Copy(WorkspaceBe+"/js-web-"+web+"/target/js-web-"+web+"/WEB-INF", buildDir+"/web/"+web+"/web-"+web+"/webapps/ROOT/WEB-INF")
		if err != nil {
			return err
		}
	}

	//fe code api
	fmt.Println("copy web fe api")
	err := file.Copy(FeCodeDir()+"/首页/onLineDebugTool", buildDir+"/web/api/web-api/webapps/ROOT")
	if err != nil {
		return err
	}
	//fe code building
	fmt.Println("copy web fe building")
	err = file.Copy(FeCodeDir()+"/控制台/Console_ng1", buildDir+"/web/building/web-building/webapps/ROOT")
	if err != nil {
		return err
	}
	//fe code datacenter
	fmt.Println("copy web fe datacenter")
	err = file.Copy(FeCodeDir()+"/控制台/DataAnalysisManagement", buildDir+"/web/datacenter/web-datacenter/webapps/ROOT")
	if err != nil {
		return err
	}
	//fe code device
	fmt.Println("copy web fe device")
	err = file.Copy(FeCodeDir()+"/控制台/DeviceManagement", buildDir+"/web/device/web-device/webapps/ROOT")
	if err != nil {
		return err
	}
	//fe code rail
	fmt.Println("copy web fe rail")
	err = file.Copy(FeCodeDir()+"/控制台/ElectricFenceManagement", buildDir+"/web/rail/web-rail/webapps/ROOT")
	if err != nil {
		return err
	}
	//fe code system
	fmt.Println("copy web fe system")
	err = file.Copy(FeCodeDir()+"/首页/Home", buildDir+"/web/system/web-system/webapps/ROOT")
	if err != nil {
		return err
	}
	// rm system unused files(media, zip, etc)
	err = os.RemoveAll(buildDir + "/web/system/web-system/webapps/ROOT/media")
	if err != nil {
		return err
	}
	err = os.RemoveAll(buildDir + "/web/system/web-system/webapps/ROOT/download/location")
	if err != nil {
		return err
	}
	return nil
}

//#复制dubbo service服务
//rm -rf $setup_dir/build/service
//mkdir $setup_dir/build/service
//
//all_sers=("building" "datacenter" "device" "rail" "system" "nettyserver" "timer" "notify")
//for ser in ${all_sers[@]}
//do
//echo 'copying service '${ser}'...'
//  mkdir $setup_dir/build/service/${ser}
//cp -r $workspace_user/js-service-${ser}/target/lib $workspace_user/js-service-${ser}/target/js-service-${ser}.jar $setup_dir/build/service/${ser}/
//done

func copyDubboService() error {
	buildDir := BuildDir()

	sers := [] string{"building", "datacenter", "device", "rail", "system", "nettyserver", "timer", "notify"}
	for _, ser := range sers {
		fmt.Printf("copy service %s\n", ser)
		err := file.Copy(WorkspaceBe+"/js-service-"+ser+"/target/lib", buildDir+"/service/"+ser+"/lib")
		if err != nil {
			return err
		}
		err = file.Copy(WorkspaceBe+"/js-service-"+ser+"/target/js-service-"+ser+".jar", buildDir+"/service/"+ser+"/js-service-"+ser+".jar")
		if err != nil {
			return err
		}
	}
	return nil
}

//复制定位引擎相关服务
//cp -r template/build/locate $setup_dir/build/
//echo 'copying locate eg...'
//cp -rf $workspace_user/js-engine/target/eg_lib $workspace_user/js-engine/target/eg.jar $setup_dir/build/locate/run/
//echo 'copying locate rd...'
//cp -rf $workspace_user/js-receiverdispatcher/target/rd_lib $workspace_user/js-receiverdispatcher/target/rd.jar $setup_dir/build/locate/run/

//echo 'copying locate web...'
//cp -rf $workspace_user/js-weblocate/target/js-weblocate/app $workspace_user/js-weblocate/target/js-weblocate/WEB-INF $web_tomcat_name/webapps/ROOT/
func copyLocateFiles() error {
	buildDir := BuildDir()

	fmt.Println("copying eg...")
	err := file.Copy(WorkspaceBe+"/js-engine/target/eg_lib", buildDir+"/locate/run/eg_lib")
	if err != nil {
		return err
	}
	err = file.Copy(WorkspaceBe+"/js-engine/target/eg.jar", buildDir+"/locate/run/eg.jar")
	if err != nil {
		return err
	}

	fmt.Println("copying rd...")
	err = file.Copy(WorkspaceBe+"/js-receiverdispatcher/target/rd_lib", buildDir+"/locate/run/rd_lib")
	if err != nil {
		return err
	}
	err = file.Copy(WorkspaceBe+"/js-receiverdispatcher/target/rd.jar", buildDir+"/locate/run/rd.jar")
	if err != nil {
		return err
	}

	fmt.Println("copying web-locate...")
	err = file.Copy(WorkspaceBe+"/js-weblocate/target/js-weblocate/app", buildDir+"/locate/tomcat-locate/webapps/ROOT/app")
	if err != nil {
		return err
	}
	err = file.Copy(WorkspaceBe+"/js-weblocate/target/js-weblocate/WEB-INF", buildDir+"/locate/tomcat-locate/webapps/ROOT/WEB-INF")
	if err != nil {
		return err
	}

	return nil
}

//初始化安装目录
func initDir() error {
	setupDir := SetupPath()
	buildDir := BuildDir()
	//
	if file.FileExist(setupDir) {
		return &exception.PackageError{Msg: fmt.Sprintf("Setup directory %s already exists.", SetupName())}
	}

	mkSetupDirErr := os.Mkdir(setupDir, 0744)
	if mkSetupDirErr != nil {
		return &exception.PackageError{Msg: fmt.Sprintf("Create directory %s err %s", setupDir, mkSetupDirErr.Error())}
	}
	mkBuildDirErr := os.Mkdir(buildDir, 0744)
	if mkBuildDirErr != nil {
		return &exception.PackageError{Msg: fmt.Sprintf("Create directory %s err %s", buildDir, mkBuildDirErr.Error())}
	}
	return nil
}

//获取要打包的版本号
func Version() string {
	return *version
}

//获取打包文件夹名称
func SetupName() string {
	return Version() + "-setup"
}

//目标安装文件夹目录
func SetupPath() string {
	return JsOpenDir + "/" + SetupName()
}

func ZipPath() string {
	return JsOpenDir + "/" + Version() + ".zip"
}

//目标程序目录
func BuildDir() string {
	return SetupPath() + "/" + "build"
}

//安装模板目录
func TemplateDir() string {
	return JsOpenDir + "/deploy/template"
}

func FeCodeDir() string {
	return WorkspaceFe + "/" + Version()
}
