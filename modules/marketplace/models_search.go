package marketplace

import (
	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/mapping"
)

const textFieldAnalyzer = "standard"

var (
	BleveIndex bleve.Index
)

func buildIndexMapping() *mapping.IndexMappingImpl {

	enTextFieldMapping := bleve.NewTextFieldMapping()
	enTextFieldMapping.Analyzer = textFieldAnalyzer

	storeFieldOnlyMapping := bleve.NewTextFieldMapping()
	storeFieldOnlyMapping.Index = false
	storeFieldOnlyMapping.IncludeTermVectors = false
	storeFieldOnlyMapping.IncludeInAll = false

	itemMapping := bleve.NewDocumentMapping()
	itemMapping.AddFieldMappingsAt("Name", enTextFieldMapping)
	itemMapping.AddFieldMappingsAt("Description", enTextFieldMapping)
	itemMapping.AddFieldMappingsAt("Category", enTextFieldMapping)
	itemMapping.AddFieldMappingsAt("SubCategory", enTextFieldMapping)

	indexMapping := bleve.NewIndexMapping()
	indexMapping.AddDocumentMapping("item", itemMapping)
	indexMapping.DefaultAnalyzer = textFieldAnalyzer

	return indexMapping
}

func init() {
	var err error
	BleveIndex, err = bleve.New("./data/index.bleve", buildIndexMapping())
	if err != nil {
		BleveIndex, err = bleve.Open("./data/index.bleve")
		if err != nil {
			panic(err)
		}
	}
}

func SearchItems(text string) []string {
	query := bleve.NewMatchQuery(text)
	search := bleve.NewSearchRequest(query)
	search.Size = 500

	searchResults, err := BleveIndex.Search(search)

	if err != nil {
		panic(err)
	}
	ids := []string{}
	for _, hit := range searchResults.Hits {
		ids = append(ids, hit.ID)
	}
	return ids
}
