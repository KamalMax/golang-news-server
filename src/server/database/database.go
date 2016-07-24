package database

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"net/url"
)

func init() {
	CreateTable()
}

type NewsFeedSource struct {
	Context  string `sql:"type:VARCHAR(20);not null" gorm:"primary_key"`
	Category string `sql:"type:VARCHAR(20);not null" gorm:"primary_key"`
	Link     string `sql:"unique;not null"`
}

func (this *NewsFeedSource) isValid() bool {
	result := true
	if _, err := url.Parse(this.Link); err != nil {
		fmt.Println("url error")
		result = false
	}

	if this.Category == "" || this.Context == "" {
		result = false
	}

	return result
}

func CreateTable() {
	db := getConnectionDB()
	defer db.Close()

	if !db.HasTable(&NewsFeedSource{}) {
		db.CreateTable(&NewsFeedSource{})
	}
}

func GetAllContextsAndCategories() *[]NewsFeedSource {
	db := getConnectionDB()
	defer db.Close()

	resultArr := []NewsFeedSource{}
	db.Select("context, category").Find(&resultArr)
	return &resultArr
}

func GetLinkByContextAndCategory(context, category string) string {
	db := getConnectionDB()
	defer db.Close()

	fs := new(NewsFeedSource)
	db.Where("context = ? AND category = ?", context, category).First(fs)
	return fs.Link
}

func Add(feedSource *NewsFeedSource) {
	if feedSource.isValid() {
		db := getConnectionDB()
		defer db.Close()
		db.Create(feedSource)
	} else {
		fmt.Println("Adding forbidden")
	}
}

func GetCategoriesMap() map[string][]string {
	db := getConnectionDB()
	defer db.Close()

	cMap := make(map[string][]string)
	keyRows, _ :=
		db.Raw(
			`SELECT DISTINCT category 
			FROM news_feed_sources`).Rows()

	for keyRows.Next() {
		var key string
		keyRows.Scan(&key)
		valueRows, _ :=
			db.Raw(
				`SELECT DISTINCT context 
				FROM news_feed_sources 
				WHERE category='" + key + "'`).Rows()

		for valueRows.Next() {
			var value string
			valueRows.Scan(&value)
			cMap[key] = append(cMap[key], value)
		}
	}
	return cMap
}

func getConnectionDB() *gorm.DB {
	db, err := gorm.Open(
		"mysql", "iris:iris@max@/news_db?charset=utf8&parseTime=True&loc=Local")
	checkErr(err)
	return db
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
