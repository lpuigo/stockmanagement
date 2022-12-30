package manager

import (
	"encoding/json"
	"github.com/lpuig/batec/stockmanagement/src/backend/model/article"
	"io"
)

func (m Manager) GetArticles(writer io.Writer) error {
	return json.NewEncoder(writer).Encode(m.Articles.GetArticles())
}

func (m Manager) UpdateArticles(updatedArticles []*article.Article) error {
	return m.Articles.UpdateArticles(updatedArticles)
}
