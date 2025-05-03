/*
  Warnings:

  - A unique constraint covering the columns `[name]` on the table `Emotion` will be added. If there are existing duplicate values, this will fail.
  - A unique constraint covering the columns `[name]` on the table `EmotionBase` will be added. If there are existing duplicate values, this will fail.

*/
-- DropForeignKey
ALTER TABLE "Emotion" DROP CONSTRAINT "Emotion_emotionBaseID_fkey";

-- CreateIndex
CREATE UNIQUE INDEX "Emotion_name_key" ON "Emotion"("name");

-- CreateIndex
CREATE UNIQUE INDEX "EmotionBase_name_key" ON "EmotionBase"("name");

-- AddForeignKey
ALTER TABLE "Emotion" ADD CONSTRAINT "Emotion_emotionBaseID_fkey" FOREIGN KEY ("emotionBaseID") REFERENCES "EmotionBase"("id") ON DELETE CASCADE ON UPDATE CASCADE;
