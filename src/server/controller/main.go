package controller

import (
	"../../feedProvider"
	"../database"
	"encoding/json"
	"github.com/kataras/iris"
	"strings"
)

func StartService() {
	register()
	iris.Listen(":8080")
}

func register() {
	iris.Get("/allpaths", getPaths)

	feedSourcesArr := database.AllContextsAndCategories()
	for _, feedSource := range *feedSourcesArr {
		iris.Get(makePath(feedSource), getNews)
	}
}

func makePath(feedSource database.NewsFeedSource) string {
	return ("/" + feedSource.Context + "/" + feedSource.Category)
}

func getNews(context *iris.Context) {
	fsArgs :=
		strings.Split(string(context.Path())[1:], "/")
	feedSourceLink :=
		database.LinkByContextAndCategory(fsArgs[0], fsArgs[1])

	feed := feedProvider.ReadNewsFeed(feedSourceLink)

	if feed != nil {
		jsonFeed, _ := json.Marshal(feed)

		context.SetHeader("Content-Type", "application/json")
		context.SetHeader("Access-Control-Allow-Origin", "*")
		context.SetHeader("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
		context.Write(string(jsonFeed))
	} else {
		context.NotFound()
	}
}

func getPaths(context *iris.Context) {
	pathsMap := database.GetCategoriesMap()
	jsonPaths, _ := json.Marshal(pathsMap)

	context.SetHeader("Content-Type", "application/json")
	context.SetHeader("Access-Control-Allow-Origin", "*")
	context.SetHeader("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
	context.Write(string(jsonPaths))
}
