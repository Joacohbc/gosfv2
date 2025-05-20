-- Add 'google_id' column to 'users' table
ALTER TABLE `users` ADD COLUMN `google_id` TEXT;
CREATE UNIQUE INDEX IF NOT EXISTS `users_google_id` ON `users` (`google_id`);

-- Add 'email' column to 'users' table
ALTER TABLE `users` ADD COLUMN `email` TEXT;
CREATE UNIQUE INDEX IF NOT EXISTS `users_email` ON `users` (`email`);
