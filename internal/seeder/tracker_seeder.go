package seeder

import (
	"cesizen/api/internal/database/prisma/db"
	"context"
	"fmt"
	"time"
)

func SeedTrackers(client *db.PrismaClient, users []db.UserModel, emotions []db.EmotionModel, n int) ([]db.TrackerModel, error) {
	background := context.Background()

	var trackers []db.TrackerModel
	for i := 0; i < n; i++ {
		user := users[i%len(users)]
		emotion := emotions[i%len(emotions)]

		tracker, err := client.Tracker.CreateOne(
			db.Tracker.User.Link(
				db.User.ID.Equals(user.ID),
			),
			db.Tracker.Emotion.Link(
				db.Emotion.ID.Equals(emotion.ID),
			),
			db.Tracker.Description.Set(fmt.Sprintf("Tracker %d", i)),
			db.Tracker.CreateAt.Set(time.Now().AddDate(0, 0, -i)), // pour variation temporelle
		).Exec(background)
		if err != nil {
			return nil, err
		}
		trackers = append(trackers, *tracker)
	}
	return trackers, nil
}
