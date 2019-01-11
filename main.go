package main

import (
	"fmt"
	"html/template"
	"os"
	"strings"

	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/bitrise-tools/go-steputils/input"
	"github.com/bitrise-tools/go-steputils/tools"
	"github.com/tamasbazs/steps-apk-info/apkutils"
)

var fileBaseNamesToSkip = []string{".DS_Store"}

// ConfigsModel ...
type ConfigsModel struct {
	ApkPath string
}

func createConfigsModelFromEnvs() ConfigsModel {
	return ConfigsModel{
		ApkPath: os.Getenv("apk_path"),
	}
}

func (configs ConfigsModel) validate() error {
	if err := input.ValidateIfPathExists(configs.ApkPath); err != nil {
		return fmt.Errorf("ApkPath - %s", err)
	}
	return nil
}

func (configs ConfigsModel) print() {
	log.Infof("Configs:")
	log.Printf("- ApkPath: %s", configs.ApkPath)
}

func fail(format string, v ...interface{}) {
	log.Errorf(format, v...)
	os.Exit(1)
}

func main() {
	configs := createConfigsModelFromEnvs()
	configs.ApkPath = strings.TrimSpace(configs.ApkPath)

	fmt.Println()
	configs.print()

	if err := configs.validate(); err != nil {
		fail("Issue with input: %s", err)
	}

	absApkPth, err := pathutil.AbsPath(configs.ApkPath)
	if err != nil {
		fail("Failed to expand path: %s, error: %s", configs.ApkPath, err)
	}

	// Collect files to deploy
	isApkPathDir, err := pathutil.IsDirExists(absApkPth)
	if err != nil {
		fail("Failed to check if file (%s), error: %s", absApkPth, err)
	}

	if isApkPathDir {
		fmt.Println()
		log.Infof("Input is a folder")

		os.Exit(3)
	}

	apkInfo, err := apkutils.GetAPKInfo(absApkPth)
	if err != nil {
		fail("Get APK info failed, error: %s", err)
		os.Exit(2)
	}
	if err := tools.ExportEnvironmentWithEnvman("ANDROID_APP_PACKAGE_NAME", apkInfo.PackageName); err != nil {
		fail("Failed to export ANDROID_APP_PACKAGE_NAME, error: %s", err)
	}
	log.Printf("ANDROID_APP_PACKAGE_NAME value: %s", apkInfo.PackageName)

	if err := tools.ExportEnvironmentWithEnvman("ANDROID_APP_NAME", apkInfo.AppName); err != nil {
		fail("Failed to export ANDROID_APP_NAME, error: %s", err)
	}
	log.Printf("ANDROID_APP_NAME value: %s", apkInfo.AppName)

	if err := tools.ExportEnvironmentWithEnvman("ANDROID_APP_VERSION_NAME", apkInfo.VersionName); err != nil {
		fail("Failed to export ANDROID_APP_VERSION_NAME, error: %s", err)
	}
	log.Printf("ANDROID_APP_VERSION_NAME value: %s", apkInfo.VersionName)

	if err := tools.ExportEnvironmentWithEnvman("ANDROID_APP_VERSION_CODE", apkInfo.VersionCode); err != nil {
		fail("Failed to export ANDROID_APP_VERSION_CODE, error: %s", err)
	}
	log.Printf("ANDROID_APP_VERSION_CODE value: %s", apkInfo.VersionCode)

	fmt.Println()
	log.Donef("Success")
	// log.Printf("You can find the Artifact on Bitrise, on the Build's page: %s", configs.BuildURL)
}

func validateGoTemplate(publicInstallPageMapFormat string) error {
	temp := template.New("Public Install Page Map template")

	_, err := temp.Parse(publicInstallPageMapFormat)
	return err
}
