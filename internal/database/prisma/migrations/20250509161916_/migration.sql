/*
  Warnings:

  - Added the required column `image` to the `Emotion` table without a default value. This is not possible if the table is not empty.

*/
-- AlterTable
ALTER TABLE "Emotion" ADD COLUMN     "image" TEXT NOT NULL;
