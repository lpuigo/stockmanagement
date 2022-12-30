package manager

type ManagerConfig struct {
	UsersDir       string
	ActorsDir      string
	ArticlesDir    string
	SaveArchiveDir string

	SessionKey string
}
