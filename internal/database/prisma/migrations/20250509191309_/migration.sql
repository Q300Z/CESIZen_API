-- DropForeignKey
ALTER TABLE "Tracker" DROP CONSTRAINT "Tracker_emotionID_fkey";

-- AddForeignKey
ALTER TABLE "Tracker" ADD CONSTRAINT "Tracker_emotionID_fkey" FOREIGN KEY ("emotionID") REFERENCES "Emotion"("id") ON DELETE CASCADE ON UPDATE CASCADE;
