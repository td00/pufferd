package migration

import (
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/pufferpanel/pufferd/logging"
	"encoding/json"
	"github.com/pufferpanel/pufferd/utils"
	"os"
	"github.com/pufferpanel/pufferd/programs"
	"github.com/pufferpanel/pufferd/data/templates"
)

const Scales = "/srv/scales/data"

func MigrateFromScales() {
	templates.CopyTemplates()
	os.MkdirAll(programs.ServerFolder, 0755)

	programFiles, err := ioutil.ReadDir(Scales)
	if err != nil {
		logging.Critical("Error reading from old Scales folder", err)
		return
	}
	for _, element := range programFiles {
		if element.IsDir() {
			continue
		}
		id := strings.TrimSuffix(element.Name(), filepath.Ext(element.Name()))
		logging.Infof("Attempting to migrate %s", id)
		data, err := ioutil.ReadFile(utils.JoinPath(Scales, element.Name()))
		if err != nil {
			logging.Error("Error read server config", err)
			continue
		}
		scales := scalesServer{}
		err = json.Unmarshal(data, &scales)
		if err != nil {
			logging.Error("Error read server config", err)
			continue
		}
		newPath := utils.JoinPath("data", "servers", scales.Name)
		err = os.Rename(utils.JoinPath("/home", scales.User), newPath)
		if err != nil {
			logging.Error("Error moving folder", err);
			continue
		}

		err = filepath.Walk(newPath, func(name string, info os.FileInfo, err error) error {
			if err == nil {
				err = os.Chown(name, os.Getuid(), os.Getgid())
			}
			return err
		})
		if err != nil {
			logging.Error("Error changing owner of folder", err);
			continue
		}
		serverData := make(map[string]interface{})
		serverData["ip"] = scales.Gamehost
		serverData["port"] = scales.Gameport
		if scales.Plugin == "minecraft" {
			serverData["memory"] = scales.Build.Memory
		} else if scales.Plugin == "srcds" {
			serverData["appid"] = scales.Startup.Variables.Build_Params
			serverData["gametype"] = scales.Startup.Variables.Game
			serverData["map"] = scales.Startup.Variables.Map
		}
		programs.Create(scales.Name, scales.Plugin, serverData)
	}
	logging.Info("Migration complete, please restart pufferd to have it recognize the changes");
}

type scalesServer struct {
	Name     string            `json:"name,omitempty"`
	User     string            `json:"user,omitempty"`
	Build    scalesServerBuild `json:"build,omitempty"`
	Gameport int               `json:"gameport,omitempty"`
	Gamehost string            `json:"gamehost,omitempty"`
	Plugin   string            `json:"plugin,omitempty"`
	Startup  scalesServerStartup `json:"startup,omitempty"`
}

type scalesServerBuild struct {
	Memory int `json:"memory,omitempty"`
}

type scalesServerStartup struct {
	Variables scalesServerStartupVariables `json:"variables,omitempty"`
}

type scalesServerStartupVariables struct {
	Build_Params string `json:"build_params,omitempty"`
	Game         string `json:"game,omitempty"`
	Map          string `json:"map,omitempty"`
	Players      string `json:"players,omitempty"`
}
