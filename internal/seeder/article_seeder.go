package seeder

import (
	"cesizen/api/internal/database/prisma/db"
	"context"
	"fmt"
)

func SeedArticles(client *db.PrismaClient, users []db.UserModel, n int) ([]db.ArticleModel, error) {
	background := context.Background()

	_, err := client.Article.FindMany().Delete().Exec(background)
	if err != nil {
		return nil, err
	}

	var articles []db.ArticleModel
	for i := 0; i < n; i++ {
		user := users[i%len(users)]
		article, err := client.Article.CreateOne(
			db.Article.Title.Set(fmt.Sprintf("Article %d", i)),
			db.Article.Content.Set(fmt.Sprintf("Contenu de l'article %d", i)),
			db.Article.User.Link(db.User.ID.Equals(user.ID)),
			db.Article.Description.Set(fmt.Sprintf("Description %d", i)),
		).Exec(background)
		if err != nil {
			return nil, err
		}
		articles = append(articles, *article)
	}
	return articles, nil
}
