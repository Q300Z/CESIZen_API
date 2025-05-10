package seeder

import (
	"cesizen/api/internal/database/prisma/db"
	"context"
	"fmt"
)

func SeedEmotionBases(client *db.PrismaClient, n int) ([]db.EmotionBaseModel, error) {
	background := context.Background()

	_, err := client.EmotionBase.FindMany().Delete().Exec(background)
	if err != nil {
		return nil, err
	}

	var bases []db.EmotionBaseModel
	for i := 0; i < n; i++ {
		base, err := client.EmotionBase.CreateOne(
			db.EmotionBase.Name.Set(fmt.Sprintf("BaseEmotion%d", i)),
		).Exec(background)
		if err != nil {
			return nil, err
		}
		bases = append(bases, *base)
	}
	return bases, nil
}

func SeedEmotions(client *db.PrismaClient, bases []db.EmotionBaseModel, n int) ([]db.EmotionModel, error) {
	background := context.Background()

	_, err := client.Emotion.FindMany().Delete().Exec(background)
	if err != nil {
		return nil, err
	}

	var emotions []db.EmotionModel
	for i := 0; i < n; i++ {
		base := bases[i%len(bases)]
		emotion, err := client.Emotion.CreateOne(
			db.Emotion.Name.Set(fmt.Sprintf("Emotion%d", i)),
			db.Emotion.Image.Set(fmt.Sprintf("https://example.com/image%d.png", i)),
			db.Emotion.EmotionBase.Link(db.EmotionBase.ID.Equals(base.ID)),
		).Exec(background)
		if err != nil {
			return nil, err
		}
		emotions = append(emotions, *emotion)
	}
	return emotions, nil
}
