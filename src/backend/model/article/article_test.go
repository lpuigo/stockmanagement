package article

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadArticlesFromXlsx(t *testing.T) {
	const (
		articleTestDir string = `C:\Users\Laurent\Google Drive\Travail\Batec\2022-11-04 - Desc Probleme & besoin`
		articleTestXls string = `Catalogue Articles`
	)

	inF, err := os.Open(filepath.Join(articleTestDir, articleTestXls+".xlsx"))
	if err != nil {
		t.Fatalf("open in file returned unexpected: %s", err.Error())
	}
	defer inF.Close()
	articles, err := LoadArticlesFromXlsx(inF)
	if err != nil {
		t.Fatalf("LoadArticlesFromXlsx returned unexpected: %s", err.Error())
	}

	t.Logf("read %d articles\n", len(articles))

	outF, err := os.Create(filepath.Join(articleTestDir, articleTestXls+"_res.xlsx"))
	if err != nil {
		t.Fatalf("create out file returned unexpected: %s", err.Error())
	}
	defer outF.Close()
	err = WriteArticlesToXlsx(outF, articles)
	if err != nil {
		t.Fatalf("WriteArticlesToXlsx returned unexpected: %s", err.Error())
	}
}

func Test_UpdateArticlePersisterFromXls(t *testing.T) {
	const (
		articlePersisterDir string = `C:\Users\Laurent\Golang\src\github.com\lpuig\batec\stockmanagement\Ressources\Articles`
		articleTestDir      string = `C:\Users\Laurent\Google Drive\Travail\Batec\2022-11-04 - Desc Probleme & besoin`
		articleTestXls      string = `Catalogue Articles`
	)

	ap, err := NewArticlesPersister(articlePersisterDir)
	if err != nil {
		t.Fatalf("NewArticlesPersister returned unexpected: %s", err.Error())
	}
	ap.NoDelay()
	err = ap.LoadDirectory()
	if err != nil {
		t.Fatalf("LoadDirectory returned unexpected: %s", err.Error())
	}

	inF, err := os.Open(filepath.Join(articleTestDir, articleTestXls+".xlsx"))
	if err != nil {
		t.Fatalf("open in file returned unexpected: %s", err.Error())
	}
	defer inF.Close()
	articles, err := LoadArticlesFromXlsx(inF)
	if err != nil {
		t.Fatalf("LoadArticlesFromXlsx returned unexpected: %s", err.Error())
	}

	t.Logf("read %d articles\n", len(articles))

	for _, article := range articles {
		ar := NewArticleRecordFromArticle(article)

		if ar.Id < 0 {
			// it a new Article, add it
			ap.Add(ar)
			continue
		}
		// its an already existing article, update it
		err = ap.Update(ar)
		if err != nil {
			t.Fatalf("Update on id %d returned unexpected: %s", ar.Id, err.Error())
		}
	}
}
