package news

import (
	"encoding/json"
	"strconv"

	"github.com/priobike/priobike-biker-swarm/common"
)

type NewsArticle struct {
	CategoryId int `json:"category_id"`
}

func FetchNews(deployment common.Deployment) {
	urlNews := "https://" + deployment.BaseUrl() + "/news-service/news/articles"

	responseBody := common.Get(urlNews, "News Articles")
	newsArticles := []NewsArticle{}
	json.Unmarshal(responseBody, &newsArticles)

	for _, newsArticle := range newsArticles {
		urlCategory := "https://" + deployment.BaseUrl() + "/news-service/news/category/" + strconv.Itoa(newsArticle.CategoryId)
		common.Get(urlCategory, "News Category")
	}
}
