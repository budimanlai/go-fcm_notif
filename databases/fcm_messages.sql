CREATE TABLE `fcm_messages` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `to_phone` varchar(25) NOT NULL DEFAULT '',
  `token` varchar(255) NOT NULL DEFAULT '',
  `title` varchar(256) NOT NULL DEFAULT '',
  `body` varchar(256) NOT NULL DEFAULT '',
  `data` text DEFAULT NULL,
  `status` varchar(15) NOT NULL DEFAULT 'pending',
  `response_log` text DEFAULT NULL,
  `created_at` datetime NOT NULL,
  `sended_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;