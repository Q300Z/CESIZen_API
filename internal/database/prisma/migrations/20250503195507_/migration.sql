/*
  Warnings:

  - You are about to drop the column `date` on the `Tracker` table. All the data in the column will be lost.
  - Made the column `emotionID` on table `Tracker` required. This step will fail if there are existing NULL values in that column.

*/
-- DropForeignKey
ALTER TABLE "Tracker" DROP CONSTRAINT "Tracker_emotionID_fkey";

-- AlterTable
ALTER TABLE "Tracker" DROP COLUMN "date",
ALTER COLUMN "emotionID" SET NOT NULL;

-- AddForeignKey
ALTER TABLE "Tracker" ADD CONSTRAINT "Tracker_emotionID_fkey" FOREIGN KEY ("emotionID") REFERENCES "Emotion"("id") ON DELETE RESTRICT ON UPDATE CASCADE;
