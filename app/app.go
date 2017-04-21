package app

import (
	"fmt"
	"runtime"

	"strings"

	"github.com/robjporter/go-functions/cisco/ucs"
)

const (
	VERSION = "0.0.0.1"
)

func New() *Application {
	runtime.GOMAXPROCS(runtime.NumCPU())
	return &Application{Version: VERSION}
}

func (a *Application) Run() {
	a.displayBanner()
	a.processResponse(ProcessCommandLineArguments())
}

func (a *Application) runAll() {
	totalServers := 0
	a.getAllSystems()
	if len(a.UCS) > 0 {
		for i := 0; i < len(a.UCS); i++ {
			a.LogInfo("Attempting to connect to UCS System", map[string]interface{}{"System": a.UCS[i].ip}, false)
			myucs := ucs.New()
			myucs.Login(a.UCS[i].ip, a.UCS[i].username, a.DecryptPassword(a.UCS[i].password))
			if myucs.LastResponse.Errors == nil {
				a.LogInfo("Successfully connected to UCS System", map[string]interface{}{"System": a.UCS[i].ip}, false)
				a.LogInfo("Getting UCS System Version", map[string]interface{}{"System": a.UCS[i].ip}, false)
				a.UCS[i].version = myucs.GetVersion()
				a.LogInfo("Getting UCS System Version", map[string]interface{}{"System": a.UCS[i].ip, "Version": a.UCS[i].version}, false)
				a.LogInfo("Getting UCS System Name", map[string]interface{}{"System": a.UCS[i].ip}, false)
				a.UCS[i].name = myucs.GetSystemName()
				a.LogInfo("Getting UCS System Version", map[string]interface{}{"System": a.UCS[i].ip, "Name": a.UCS[i].name}, false)
				a.LogInfo("Getting UCS System Servers", map[string]interface{}{"System": a.UCS[i].ip}, false)
				a.UCS[i].blades = myucs.GetSystemBlades()
				a.LogInfo("Gained UCS System Servers", map[string]interface{}{"System": a.UCS[i].ip, "Servers": len(a.UCS[i].blades)}, false)
				totalServers += len(a.UCS[i].blades)
				a.UCS[i].status = true
			} else {
				a.UCS[i].status = false
				a.LogInfo("Failed to connect to UCS System", map[string]interface{}{"System": a.UCS[i].ip}, false)
			}
			a.LogInfo("Logging out of UCS System", map[string]interface{}{"System": a.UCS[i].ip}, false)
			myucs.Logout()
		}
		a.LogInfo("Indexed all discovered UCS System Servers", map[string]interface{}{"Servers": totalServers}, false)
		a.processAllServers()
	} else {
		fmt.Println("No UCS Systems detected in the config file.  Please trying adding one and try again.")
	}
}

func (a *Application) processAllServers() {
	csv := "Active,Associated,Powered,Serial,Model,Chassis,Slot,Name,Label,Description,CPU,Memory,Associated To\n"
	if len(a.UCS) > 0 {
		for i := 0; i < len(a.UCS); i++ {
			for j := 0; j < len(a.UCS[i].blades); j++ {
				csv += isActive(a.UCS[i].blades[j].BladeAssociation, a.UCS[i].blades[j].BladePower) + ","
				csv += a.UCS[i].blades[j].BladeAssociation + ","
				csv += a.UCS[i].blades[j].BladePower + ","
				csv += a.UCS[i].blades[j].BladeSerial + ","
				csv += a.UCS[i].blades[j].BladeModel + ","
				csv += a.UCS[i].blades[j].BladeChassis + ","
				csv += a.UCS[i].blades[j].BladeSlot + ","
				csv += a.UCS[i].blades[j].BladeName + ","
				csv += a.UCS[i].blades[j].BladeLabel + ","
				csv += a.UCS[i].blades[j].BladeDescr + ","
				csv += a.UCS[i].blades[j].BladeSockets + ","
				csv += a.UCS[i].blades[j].BladeMemory + ","
				csv += a.UCS[i].blades[j].BladeAssociatedTo + ","
				csv += "\n"
			}
		}
	}
	fmt.Println(csv)
}

func isActive(assocated string, power string) string {
	if strings.ToLower(assocated) == "associated" && strings.ToLower(power) == "on" {
		return "true"
	}
	return "false"
}
