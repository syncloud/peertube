package installer

import (
	"fmt"
	"github.com/google/uuid"
	cp "github.com/otiai10/copy"
	"github.com/syncloud/golib/config"
	"github.com/syncloud/golib/linux"
	"github.com/syncloud/golib/platform"
	"go.uber.org/zap"

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
	Domain      string
	Secret      string
	DatabaseDir string
}

type Installer struct {
	newVersionFile     string
	currentVersionFile string
	configDir          string
	platformClient     *platform.Client
	database           *Database
	installFile        string
	executor           *Executor
	logger             *zap.Logger
}

func New(logger *zap.Logger) *Installer {
	configDir := path.Join(DataDir, "config")

	executor := NewExecutor(logger)
	return &Installer{
		newVersionFile:     path.Join(AppDir, "version"),
		currentVersionFile: path.Join(DataDir, "version"),
		configDir:          configDir,
		platformClient:     platform.New(),
		database:           NewDatabase(AppDir, DataDir, configDir, App, executor, logger),
		installFile:        path.Join(CommonDir, "installed"),
		executor:           executor,
		logger:             logger,
	}
}

func (i *Installer) Install() error {

	err := i.UpdateConfigs()
	if err != nil {
		return err
	}

	err = i.database.Init()
	if err != nil {
		return err
	}
	err = i.database.InitConfig()
	if err != nil {
		return err
	}

	return nil
}

func (i *Installer) Configure() error {
	if i.IsInstalled() {
		err := i.Upgrade()
		if err != nil {
			return err
		}
	} else {
		err := i.Initialize()
		if err != nil {
			return err
		}
	}

	err := i.executor.Run("snap",
		"run", "peertube.node",
		fmt.Sprintf("%s/peertube/app/dist/scripts/plugin/install.js", AppDir),
		"-p", fmt.Sprintf("%s/peertube/app/plugins/peertube-plugin-auth-openid-connect", AppDir),
	)
	if err != nil {
		i.logger.Error("failed to install plugin", zap.Error(err))
		return err
	}

	return nil
}

func (i *Installer) IsInstalled() bool {
	_, err := os.Stat(i.installFile)
	return os.IsExist(err)
}

func (i *Installer) Initialize() error {
	err := i.StorageChange()
	if err != nil {
		return err
	}

	err = i.database.Execute(
		"postgres",
		fmt.Sprintf("ALTER USER %s WITH PASSWORD '%s'", App, App),
	)
	if err != nil {
		return err
	}
	err = i.database.createDbIfMissing(App)
	if err != nil {
		return err
	}
	err = i.database.Execute("postgres", fmt.Sprintf("GRANT CREATE ON SCHEMA public TO %s", App))
	if err != nil {
		return err
	}
	err = os.WriteFile(i.installFile, []byte("installed"), 0644)
	if err != nil {
		return err
	}

	return i.UpdateVersion()
}

func (i *Installer) Upgrade() error {
	err := i.database.Restore()
	if err != nil {
		return err
	}
	err = i.StorageChange()
	if err != nil {
		return err
	}
	err = i.database.createDbIfMissing(App)
	if err != nil {
		return err
	}

	return i.UpdateVersion()
}

func (i *Installer) PreRefresh() error {
	return i.database.Backup()
}

func (i *Installer) PostRefresh() error {
	err := i.UpdateConfigs()
	if err != nil {
		return err
	}
	err = i.database.Remove()
	if err != nil {
		return err
	}
	err = i.database.Init()
	if err != nil {
		return err
	}
	err = i.database.InitConfig()
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
func (i *Installer) AccessChange() error {
	err := i.UpdateConfigs()
	if err != nil {
		return err
	}
	err = i.executor.Run("snap", "restart", "peertube")
	return err
}

func (i *Installer) StorageChange() error {
	storageDir, err := i.platformClient.InitStorage(App, App)
	if err != nil {
		return err
	}
	err = i.createMissingDirs(
		path.Join(storageDir, "tmp"),
		path.Join(storageDir, "tmp-persistent"),
		path.Join(storageDir, "bin"),
		path.Join(storageDir, "avatars"),
		path.Join(storageDir, "web-videos"),
		path.Join(storageDir, "streaming-playlists"),
		path.Join(storageDir, "original-video-files"),
		path.Join(storageDir, "redundancy"),
		path.Join(storageDir, "logs"),
		path.Join(storageDir, "previews"),
		path.Join(storageDir, "thumbnails"),
		path.Join(storageDir, "storyboards"),
		path.Join(storageDir, "torrents"),
		path.Join(storageDir, "captions"),
		path.Join(storageDir, "cache"),
		path.Join(storageDir, "plugins"),
		path.Join(storageDir, "well-known"),
		path.Join(storageDir, "client-overrides"),
	)
	if err != nil {
		return err
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
	err := linux.CreateUser(App)
	if err != nil {
		return err
	}

	err = i.StorageChange()
	if err != nil {
		return err
	}

	err = createMissingDir(path.Join(DataDir, "nginx"))
	if err != nil {
		return err
	}

	domain, err := i.platformClient.GetAppDomainName(App)
	if err != nil {
		return err
	}

	secret, err := getOrCreateUuid(path.Join(DataDir, "peertube.secret"))
	if err != nil {
		return err
	}

	variables := Variables{
		Domain:      domain,
		Secret:      secret,
		DatabaseDir: i.database.DatabaseDir(),
	}

	err = config.Generate(
		path.Join(AppDir, "config"),
		path.Join(DataDir, "config"),
		variables,
	)
	if err != nil {
		return err
	}

	err = i.FixPermissions()
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

func (i *Installer) BackupPreStop() error {
	return i.PreRefresh()
}

func (i *Installer) RestorePreStart() error {
	return i.PostRefresh()
}

func (i *Installer) RestorePostStart() error {
	return i.Configure()
}

func (i *Installer) createMissingDirs(dirs ...string) error {
	for _, dir := range dirs {
		err := createMissingDir(dir)
		if err != nil {
			i.logger.Error("cannot create dir", zap.String("dir", dir), zap.Error(err))
			return err
		}
	}
	return nil
}

func createMissingDir(dir string) error {
	_, err := os.Stat(dir)
	if os.IsNotExist(err) {
		err = os.Mkdir(dir, 0755)
		if err != nil {
			return err
		}
	}
	return nil
}

func getOrCreateUuid(file string) (string, error) {
	_, err := os.Stat(file)
	if os.IsNotExist(err) {
		secret := uuid.New().String()
		err = os.WriteFile(file, []byte(secret), 0644)
		return secret, err
	}
	content, err := os.ReadFile(file)
	if err != nil {
		return "", err
	}
	return string(content), nil
}
