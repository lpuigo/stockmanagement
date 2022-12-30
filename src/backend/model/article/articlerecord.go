package article

import (
	"encoding/json"
	"github.com/lpuig/batec/stockmanagement/src/backend/persist"
	"io"
	"os"
)

type ArticleRecord struct {
	*persist.Record
	*Article
}

// NewArticleRecord returns a new ArticleRecord
func NewArticleRecord() *ArticleRecord {
	ar := &ArticleRecord{}
	ar.Record = persist.NewRecord(func(w io.Writer) error {
		return json.NewEncoder(w).Encode(ar.Article)
	})
	return ar
}

// NewArticleRecordFrom returns a ArticleRecord populated from the given reader
func NewArticleRecordFrom(r io.Reader) (ar *ArticleRecord, err error) {
	ar = NewArticleRecord()
	err = json.NewDecoder(r).Decode(ar)
	if err != nil {
		ar = nil
		return
	}
	ar.SetId(ar.Id)
	return
}

// NewArticleRecordFromFile returns a ArticleRecord populated from the given file
func NewArticleRecordFromFile(file string) (ar *ArticleRecord, err error) {
	f, err := os.Open(file)
	if err != nil {
		return
	}
	defer f.Close()
	ar, err = NewArticleRecordFrom(f)
	if err != nil {
		ar = nil
		return
	}
	return
}

// NewArticleRecordFromArticle returns a ArticleRecord populated from given Article
func NewArticleRecordFromArticle(act *Article) *ArticleRecord {
	ar := NewArticleRecord()
	ar.Article = act
	ar.SetId(ar.Id)
	return ar
}

func (ar *ArticleRecord) SetId(id int) {
	ar.Record.SetId(id)
	ar.Id = ar.GetId()
}
