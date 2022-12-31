package manager

import (
	"fmt"
	"github.com/lpuig/batec/stockmanagement/src/backend/logger"
	"github.com/lpuig/batec/stockmanagement/src/backend/model/actor"
	"github.com/lpuig/batec/stockmanagement/src/backend/model/article"
	"github.com/lpuig/batec/stockmanagement/src/backend/model/movement"
	"github.com/lpuig/batec/stockmanagement/src/backend/model/session"
	"github.com/lpuig/batec/stockmanagement/src/backend/model/stock"
	"github.com/lpuig/batec/stockmanagement/src/backend/model/user"
	"github.com/lpuig/batec/stockmanagement/src/backend/model/worksite"
	"github.com/lpuig/batec/stockmanagement/src/backend/persist"
)

type Manager struct {
	Users     *user.UsersPersister
	Actors    *actor.ActorsPersister
	Articles  *article.ArticlesPersister
	Movements *movement.MovementsPersister
	Worksites *worksite.WorksitesPersister
	Stocks    *stock.StocksPersister

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
	acp, err := actor.NewActorsPersister(conf.ActorsDir)
	if err != nil {
		return nil, fmt.Errorf("could not create actors: %s", err.Error())
	}

	// Init Articles persister
	arp, err := article.NewArticlesPersister(conf.ArticlesDir)
	if err != nil {
		return nil, fmt.Errorf("could not create articles: %s", err.Error())
	}

	// Init Movements persister
	mp, err := movement.NewMovementsPersister(conf.MovementsDir)
	if err != nil {
		return nil, fmt.Errorf("could not create movements: %s", err.Error())
	}

	// Init Worksites persister
	wp, err := worksite.NewWorksitesPersister(conf.WorksitesDir)
	if err != nil {
		return nil, fmt.Errorf("could not create worksites: %s", err.Error())
	}

	// Init Stocks persister
	sp, err := stock.NewStocksPersister(conf.StocksDir)
	if err != nil {
		return nil, fmt.Errorf("could not create stocks: %s", err.Error())
	}

	// Init manager
	m := &Manager{
		Users:        up,
		Actors:       acp,
		Articles:     arp,
		Movements:    mp,
		Worksites:    wp,
		Stocks:       sp,
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
	logger.Entry("Server").LogInfo(fmt.Sprintf("loaded %d users", m.Users.NbUsers()))

	err = m.Actors.LoadDirectory()
	if err != nil {
		return fmt.Errorf("could not populate actor: %s", err.Error())
	}
	logger.Entry("Server").LogInfo(fmt.Sprintf("loaded %d actors", m.Actors.NbRecords()))

	err = m.Articles.LoadDirectory()
	if err != nil {
		return fmt.Errorf("could not populate article: %s", err.Error())
	}
	logger.Entry("Server").LogInfo(fmt.Sprintf("loaded %d articles", m.Articles.NbRecords()))

	err = m.Movements.LoadDirectory()
	if err != nil {
		return fmt.Errorf("could not populate movement: %s", err.Error())
	}
	logger.Entry("Server").LogInfo(fmt.Sprintf("loaded %d movements", m.Movements.NbRecords()))

	err = m.Worksites.LoadDirectory()
	if err != nil {
		return fmt.Errorf("could not populate worksite: %s", err.Error())
	}
	logger.Entry("Server").LogInfo(fmt.Sprintf("loaded %d worksites", m.Worksites.NbRecords()))

	err = m.Stocks.LoadDirectory()
	if err != nil {
		return fmt.Errorf("could not populate stock: %s", err.Error())
	}
	logger.Entry("Server").LogInfo(fmt.Sprintf("loaded %d stocks", m.Stocks.NbRecords()))

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
	saveArchive(m.Articles)
	saveArchive(m.Actors)
	saveArchive(m.Users)
	saveArchive(m.Movements)
	saveArchive(m.Worksites)
	saveArchive(m.Stocks)
	if failed != "" {
		return fmt.Errorf("failed to save%s", failed)
	}
	return nil
}
