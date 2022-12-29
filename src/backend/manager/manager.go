package manager

import (
	"fmt"
	"github.com/lpuig/batec/stockmanagement/src/backend/logger"
	"github.com/lpuig/batec/stockmanagement/src/backend/model/actor"
	"github.com/lpuig/batec/stockmanagement/src/backend/model/session"
	"github.com/lpuig/batec/stockmanagement/src/backend/model/user"
	"github.com/lpuig/batec/stockmanagement/src/backend/persist"
)

type Manager struct {
	Users  *user.UsersPersister
	Actors *actor.ActorsPersister

	SessionStore *session.SessionStore
	CurrentUser  *user.UserRecord

	Config ManagerConfig
}

func NewManager(conf ManagerConfig) (*Manager, error) {
	// Init Users persister
	up, err := user.NewUsersPersister(conf.UsersDir)
	if err != nil {
		return nil, fmt.Errorf("could not create users: %s", err.Error())
	}

	// Init Actors persister
	ap, err := actor.NewActorsPersister(conf.ActorsDir)
	if err != nil {
		return nil, fmt.Errorf("could not create actors: %s", err.Error())
	}

	// Init manager
	m := &Manager{
		Users:        up,
		Actors:       ap,
		SessionStore: session.NewSessionStore(conf.SessionKey),
		//CurrentUser: is set during session control transaction
		Config: conf,
	}

	err = m.Reload()
	if err != nil {
		return nil, err
	}

	return m, nil
}

func (m Manager) Clone() *Manager {
	return &m
}

func (m *Manager) Reload() error {
	err := m.Users.LoadDirectory()
	if err != nil {
		return fmt.Errorf("could not populate user: %s", err.Error())
	}
	logger.Entry("Server").LogInfo(fmt.Sprintf("loaded %d Users", m.Users.NbUsers()))

	err = m.Actors.LoadDirectory()
	if err != nil {
		return fmt.Errorf("could not populate actor: %s", err.Error())
	}
	logger.Entry("Server").LogInfo(fmt.Sprintf("loaded %d Actors", m.Actors.NbActors()))

	return nil
}

func (m *Manager) SaveArchive() error {
	failed := ""
	saveArchive := func(container persist.ArchivableRecordContainer) {
		err := persist.SaveRecordsArchive(m.Config.SaveArchiveDir, container)
		if err != nil {
			failed += " " + fmt.Sprintf("%s (%s)", container.GetName(), err.Error())
		}
	}
	saveArchive(m.Actors)
	saveArchive(m.Users)
	if failed != "" {
		return fmt.Errorf("failed to save%s", failed)
	}
	return nil
}
