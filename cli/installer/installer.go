package installer

import (
	cp "github.com/otiai10/copy"
	"github.com/syncloud/golib/linux"
	"github.com/syncloud/golib/platform"
	"github.com/syncloud/golib/config"

	"os"
	"path"
)

const (
	App       = "peertube"
	AppDir    = "/snap/peertube/current"
	DataDir   = "/var/snap/peertube/current"
	CommonDir = "/var/snap/peertube/common"
)

type Variables struct {
	Domain string
	AppUrl string
}

type Installer struct {
	newVersionFile     string
	currentVersionFile string
	configDir          string
	platformClient     *platform.Client
}

func New() *Installer {
	configDir := path.Join(DataDir, "config")

	return &Installer{
		newVersionFile:     path.Join(AppDir, "version"),
		currentVersionFile: path.Join(DataDir, "version"),
		configDir:          configDir,
		platformClient:     platform.New(),
	}
}

func (i *Installer) Install() error {
	err := linux.CreateUser(App)
	if err != nil {
		return err
	}

	err = os.Mkdir(path.Join(DataDir, "nginx"), 0755)
	if err != nil {
		return err
	}

	err = i.UpdateConfigs()
	if err != nil {
		return err
	}

	err = i.FixPermissions()
	if err != nil {
		return err
	}

	err = i.StorageChange()
	if err != nil {
		return err
	}
	return nil
}

func (i *Installer) Configure() error {
	return i.UpdateVersion()
}

func (i *Installer) PreRefresh() error {
	return nil
}

func (i *Installer) PostRefresh() error {
	err := i.UpdateConfigs()
	if err != nil {
		return err
	}

	err = i.ClearVersion()
	if err != nil {
		return err
	}

	err = i.FixPermissions()
	if err != nil {
		return err
	}
	return nil

}
func (i *Installer) StorageChange() error {
	storageDir, err := i.platformClient.InitStorage(App, App)
	if err != nil {
		return err
	}

	err = os.Mkdir(path.Join(storageDir, "media"), 0755)
	if err != nil {
		if !os.IsExist(err) {
			return err
		}
	}
	err = os.Mkdir(path.Join(storageDir, "cache"), 0755)
	if err != nil {
		if !os.IsExist(err) {
			return err
		}
	}
	err = linux.Chown(storageDir, App)
	if err != nil {
		return err
	}

	return nil
}

func (i *Installer) ClearVersion() error {
	return os.RemoveAll(i.currentVersionFile)
}

func (i *Installer) UpdateVersion() error {
	return cp.Copy(i.newVersionFile, i.currentVersionFile)
}

func (i *Installer) UpdateConfigs() error {

	variables := Variables{}

	err := config.Generate(
		path.Join(AppDir, "config"),
		path.Join(DataDir, "config"),
		variables,
	)
	if err != nil {
		return err
	}

	return nil

}

func (i *Installer) FixPermissions() error {
	err := linux.Chown(DataDir, App)
	if err != nil {
		return err
	}
	err = linux.Chown(CommonDir, App)
	if err != nil {
		return err
	}
	return nil
}
