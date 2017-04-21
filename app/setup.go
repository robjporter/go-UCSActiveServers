package app

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	functions2 "github.com/robjporter/go-functions"
	"github.com/robjporter/go-functions/banner"
	"github.com/robjporter/go-functions/colors"
	"github.com/robjporter/go-functions/logrus"
	"github.com/robjporter/go-functions/terminal"
	"github.com/robjporter/go-functions/viper"
	yaml "github.com/robjporter/go-functions/yaml"
)

func (a *Application) init() {
	a.Config = viper.New()
	a.Logger = logrus.New()
	a.Logger.Level = logrus.DebugLevel
	customFormatter := new(logrus.TextFormatter)
	customFormatter.TimestampFormat = "02-01-2006 15:04:05.000"
	customFormatter.FullTimestamp = true
	a.Logger.Formatter = customFormatter
	a.Logger.Out = os.Stdout

	a.Key = []byte("random123456")
}

func (a *Application) displayBanner() {
	_, w, err := terminal.GetTerminalSize()
	if err == nil {
		terminal.ClearScreen()
		banner.PrintNewFigure("UCS Service Profiles", "rounded", true)
		fmt.Println(colors.Color("Cisco Unified Computing System Service Profile Utilisation v"+a.Version, colors.BRIGHTYELLOW))
		banner.BannerPrintLineS("=", w)
	}
}

func (a *Application) LoadConfig(filename string) {
	a.init()
	a.Log("Loading Configuration File.", nil, true)
	a.ConfigFile = filename
	configName := ""
	configExtension := ""
	configPath := ""

	splits := strings.Split(filepath.Base(a.ConfigFile), ".")
	if len(splits) == 2 {
		configName = splits[0]
		configExtension = splits[1]
	}
	configPath = filepath.Dir(a.ConfigFile)

	a.Config.SetConfigName(configName)
	a.Config.SetConfigType(configExtension)
	a.Config.AddConfigPath(configPath)

	a.Log("Configuration File defined", map[string]interface{}{"Path": configPath, "Name": configName, "Extension": configExtension}, true)

	if functions2.Exists(a.ConfigFile) {
		err := a.Config.ReadInConfig()
		if err != nil {
			panic(fmt.Errorf("Fatal error config file: %s \n", err))
			os.Exit(0)
		}
		a.Log("Configuration File read successfully.", nil, true)
	} else {
		a.LogInfo("Configuration File not found.", nil, true)
	}
}

func (a *Application) LogInfo(message string, fields map[string]interface{}, infoMessage bool) {
	if infoMessage && a.Debug || !infoMessage {
		if fields != nil {
			a.Logger.WithFields(fields).Info(message)
		} else {
			a.Logger.Info(message)
		}
	}
}

func (a *Application) Log(message string, fields map[string]interface{}, debugMessage bool) {
	if debugMessage && a.Debug || !debugMessage {
		if fields != nil {
			a.Logger.WithFields(fields).Info(message)
		} else {
			a.Logger.Info(message)
		}
	}
}

func (a *Application) EncryptPassword(password string) string {
	return functions2.Encrypt(a.Key, []byte(password))
}

func (a *Application) DecryptPassword(password string) string {
	return functions2.Decrypt(a.Key, password)
}

func (a *Application) saveConfig() {
	a.LogInfo("Saving configuration file.", nil, false)
	if len(a.UCS) > 0 {
		items := a.processSystems()
		a.Config.Set("ucs.systems", items)
	}
	out, err := yaml.Marshal(a.Config.AllSettings())
	if err == nil {
		fp, err := os.Create(a.ConfigFile)
		if err == nil {
			defer fp.Close()
			_, err = fp.Write(out)
		}
	}
	a.Log("Saving configuration file complete.", nil, true)
}

func (a *Application) processSystems() []interface{} {
	var items []interface{}
	var item map[string]interface{}
	for i := 0; i < len(a.UCS); i++ {

		item = make(map[string]interface{})
		item["url"] = a.UCS[i].ip
		item["username"] = a.UCS[i].username
		item["password"] = a.UCS[i].password
		items = append(items, item)
	}
	return items
}

func (a *Application) saveFile(filename, data string) bool {
	ret := false
	f, err := os.Create(filename)
	if err == nil {
		_, err := f.Write([]byte(data))
		if err == nil {
			a.LogInfo("File has been saved successfully.", map[string]interface{}{"Filename": filename}, false)
			ret = true
		} else {
			a.LogInfo("There was a problem saving the file.", map[string]interface{}{"Error": err}, false)
		}
	}
	defer f.Close()
	return ret
}
