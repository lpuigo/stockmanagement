package article

import (
	"fmt"
	"github.com/lpuig/batec/stockmanagement/src/backend/model/date"
	"github.com/lpuig/batec/stockmanagement/src/backend/persist"
	"io"
	"sort"
)

type ArticlesPersister struct {
	*persist.Persister
}

func NewArticlesPersister(dir string) (*ArticlesPersister, error) {
	up := &ArticlesPersister{
		Persister: persist.NewPersister("Articles", dir),
	}
	err := up.CheckDirectory()
	if err != nil {
		return nil, err
	}
	return up, nil
}

func (ap ArticlesPersister) NbArticles() int {
	return ap.NbRecords()
}

// LoadDirectory loads all persisted Articles Records
func (ap *ArticlesPersister) LoadDirectory() error {
	return ap.Persister.LoadDirectory("article", func(file string) (persist.Recorder, error) {
		return NewArticleRecordFromFile(file)
	})
}

// Add adds the given ArticleRecord to the USersPersister and return its (updated with new id) ArticleRecord
func (ap *ArticlesPersister) Add(nar *ArticleRecord) *ArticleRecord {
	nar.SetCreateDate()
	ap.Persister.Add(nar)
	return nar
}

// Update updates the given ArticleRecord
func (ap *ArticlesPersister) Update(uar *ArticleRecord) error {
	uar.SetUpdateDate()
	return ap.Persister.Update(uar)
}

// Remove removes the given ArticleRecord from the ArticlesPersister (pertaining file is moved to deleted dir)
func (ap *ArticlesPersister) Remove(rar *ArticleRecord) error {
	rar.SetDeleteDate()
	return ap.Persister.Remove(rar)
}

// GetById returns the ArticleRecord with given Id (or nil if Id not found)
func (ap *ArticlesPersister) GetById(id int) *ArticleRecord {
	return ap.Persister.GetById(id).(*ArticleRecord)
}

// GetByRef returns the ArticleRecord with given Ref (or nil if Ref not found)
func (ap *ArticlesPersister) GetByRef(ref string) *ArticleRecord {
	for _, r := range ap.GetRecords() {
		if ar, ok := r.(*ArticleRecord); ok && ar.Ref == ref {
			return ar
		}
	}
	return nil
}

func (ap *ArticlesPersister) GetArticles() []*Article {
	res := []*Article{}
	for _, r := range ap.GetRecords() {
		if ar, ok := r.(*ArticleRecord); ok {
			res = append(res, ar.Article)
		}
	}
	return res
}

func (ap *ArticlesPersister) UpdateArticles(updatedArticles []*Article) error {
	for _, updAct := range updatedArticles {
		updActRec := NewArticleRecordFromArticle(updAct)
		if updActRec.Id < 0 { // New Article, add it instead of update
			ap.Add(updActRec)
			continue
		}
		err := ap.Update(updActRec)
		if err != nil {
			fmt.Errorf("could not update article '%s' (id: %d)", updActRec.Ref, updActRec.Id)
		}
	}
	return nil
}

// Export / Import to XLSx methods
func (ap *ArticlesPersister) ExportName() string {
	return fmt.Sprintf("%s Catalogue Articles.xlsx", date.Today().String())
}

func (ap *ArticlesPersister) XLSExport(writer io.Writer) error {
	articles := ap.GetArticles()
	sort.Slice(articles, func(i, j int) bool {
		if articles[i].Manufacturer == articles[j].Manufacturer {
			return articles[i].Designation < articles[j].Designation
		}
		return articles[i].Manufacturer < articles[j].Manufacturer
	})
	return WriteArticlesToXlsx(writer, articles)
}
