-- --------------------------------------------------------
-- Host:                         127.0.0.1
-- Server version:               8.0.30 - MySQL Community Server - GPL
-- Server OS:                    Win64
-- HeidiSQL Version:             12.1.0.6537
-- --------------------------------------------------------

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET NAMES utf8 */;
/*!50503 SET NAMES utf8mb4 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;


-- Dumping database structure for cinema_tickets
CREATE DATABASE IF NOT EXISTS `cinema_tickets` /*!40100 DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci */ /*!80016 DEFAULT ENCRYPTION='N' */;
USE `cinema_tickets`;

-- Dumping structure for table cinema_tickets.cinemas
CREATE TABLE IF NOT EXISTS `cinemas` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  `location` varchar(255) NOT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- Dumping data for table cinema_tickets.cinemas: ~2 rows (approximately)
INSERT INTO `cinemas` (`id`, `name`, `location`, `created_at`, `updated_at`) VALUES
	(1, 'Studio A', 'Lt 1', '2024-12-05 12:16:58.905', '2024-12-05 12:16:58.905'),
	(2, 'Studio B', 'Lt 2', '2024-12-05 12:17:08.640', '2024-12-06 06:40:50.934');

-- Dumping structure for table cinema_tickets.films
CREATE TABLE IF NOT EXISTS `films` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `title` varchar(255) NOT NULL,
  `description` text,
  `duration` bigint NOT NULL,
  `genre` varchar(255) NOT NULL,
  `rating` varchar(5) NOT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- Dumping data for table cinema_tickets.films: ~4 rows (approximately)
INSERT INTO `films` (`id`, `title`, `description`, `duration`, `genre`, `rating`, `created_at`, `updated_at`) VALUES
	(1, 'Agak Laen 2', 'Film komedi dari podcast Agak Laen kedua', 102, 'komedi', 'PG-13', '2024-12-05 06:53:26.591', '2024-12-05 07:01:44.011'),
	(2, 'Agak Laen', 'Film komedi dari podcast Agak Laen', 130, 'komedi', 'R', '2024-12-05 06:53:56.147', '2024-12-05 06:53:56.147'),
	(3, 'Kang Mak from Pee Mak', 'Film komedi horror adaptasi dari thailand', 110, 'komedi', 'R', '2024-12-05 06:55:36.804', '2024-12-05 06:55:36.804');

-- Dumping structure for table cinema_tickets.schedules
CREATE TABLE IF NOT EXISTS `schedules` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `film_id` bigint unsigned NOT NULL,
  `cinema_id` bigint unsigned NOT NULL,
  `time` varchar(250) NOT NULL,
  `price` bigint NOT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `fk_schedules_film` (`film_id`),
  KEY `fk_schedules_cinema` (`cinema_id`),
  CONSTRAINT `fk_schedules_cinema` FOREIGN KEY (`cinema_id`) REFERENCES `cinemas` (`id`),
  CONSTRAINT `fk_schedules_film` FOREIGN KEY (`film_id`) REFERENCES `films` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- Dumping data for table cinema_tickets.schedules: ~1 rows (approximately)
INSERT INTO `schedules` (`id`, `film_id`, `cinema_id`, `time`, `price`, `created_at`, `updated_at`) VALUES
	(1, 1, 1, '2024-12-10 15:01:00', 40000, '2024-12-06 13:27:49.906', '2024-12-07 07:50:20.022'),
	(2, 2, 1, '2024-12-05 15:01:00', 40000, '2024-12-07 12:40:22.806', '2024-12-07 12:40:22.806');

-- Dumping structure for table cinema_tickets.seats
CREATE TABLE IF NOT EXISTS `seats` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `number` varchar(10) NOT NULL,
  `cinema_id` bigint unsigned NOT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `fk_seats_cinema` (`cinema_id`),
  KEY `idx_seats_deleted_at` (`deleted_at`),
  CONSTRAINT `fk_seats_cinema` FOREIGN KEY (`cinema_id`) REFERENCES `cinemas` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=9 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- Dumping data for table cinema_tickets.seats: ~5 rows (approximately)
INSERT INTO `seats` (`id`, `number`, `cinema_id`, `created_at`, `updated_at`, `deleted_at`) VALUES
	(1, 'A1', 1, NULL, NULL, NULL),
	(2, 'A1', 2, '2024-12-06 06:22:06.975', '2024-12-06 06:22:06.975', NULL),
	(5, 'A2', 1, NULL, NULL, NULL),
	(6, 'A3', 1, NULL, NULL, NULL),
	(7, 'A4', 1, NULL, NULL, NULL),
	(8, 'A5', 1, NULL, NULL, NULL);

-- Dumping structure for table cinema_tickets.transactions
CREATE TABLE IF NOT EXISTS `transactions` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `user_id` bigint unsigned NOT NULL,
  `schedule_id` bigint unsigned NOT NULL,
  `total_price` bigint NOT NULL,
  `status` enum('pending','completed','canceled') DEFAULT 'pending',
  `transaction_id` varchar(255) NOT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uni_transactions_transaction_id` (`transaction_id`),
  KEY `fk_transactions_user` (`user_id`),
  KEY `fk_transactions_schedule` (`schedule_id`),
  CONSTRAINT `fk_transactions_schedule` FOREIGN KEY (`schedule_id`) REFERENCES `schedules` (`id`),
  CONSTRAINT `fk_transactions_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=15 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- Dumping data for table cinema_tickets.transactions: ~3 rows (approximately)
INSERT INTO `transactions` (`id`, `user_id`, `schedule_id`, `total_price`, `status`, `transaction_id`, `created_at`, `updated_at`) VALUES
	(12, 1, 1, 80000, 'canceled', '20241208-ugORv8Hm', '2024-12-08 06:04:32.364', '2024-12-08 06:09:17.602'),
	(13, 1, 1, 80000, 'pending', '20241208-A3IcOdra', '2024-12-08 06:09:33.914', '2024-12-08 06:09:33.914'),
	(14, 1, 1, 40000, 'pending', '20241208-iNmFs4q1', '2024-12-08 06:11:20.144', '2024-12-08 06:11:20.144');

-- Dumping structure for table cinema_tickets.transaction_seats
CREATE TABLE IF NOT EXISTS `transaction_seats` (
  `transaction_id` varchar(255) DEFAULT NULL,
  `seat_id` bigint unsigned DEFAULT NULL,
  KEY `fk_transactions_transaction_seat` (`transaction_id`),
  KEY `fk_transaction_seats_seat` (`seat_id`),
  CONSTRAINT `fk_transaction_seats_seat` FOREIGN KEY (`seat_id`) REFERENCES `seats` (`id`),
  CONSTRAINT `fk_transactions_transaction_seat` FOREIGN KEY (`transaction_id`) REFERENCES `transactions` (`transaction_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- Dumping data for table cinema_tickets.transaction_seats: ~5 rows (approximately)
INSERT INTO `transaction_seats` (`transaction_id`, `seat_id`) VALUES
	('20241208-ugORv8Hm', 1),
	('20241208-ugORv8Hm', 5),
	('20241208-A3IcOdra', 1),
	('20241208-A3IcOdra', 5),
	('20241208-iNmFs4q1', 6);

-- Dumping structure for table cinema_tickets.users
CREATE TABLE IF NOT EXISTS `users` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  `email` varchar(255) NOT NULL,
  `password` varchar(255) NOT NULL,
  `role` enum('admin','user') DEFAULT 'user',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uni_users_email` (`email`)
) ENGINE=InnoDB AUTO_INCREMENT=9 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- Dumping data for table cinema_tickets.users: ~6 rows (approximately)
INSERT INTO `users` (`id`, `name`, `email`, `password`, `role`, `created_at`, `updated_at`) VALUES
	(1, 'Admin1', 'admin1@mail.com', '7f95b733f4210c71482904eb422143f8', 'admin', '2024-12-05 13:21:43.637', '2024-12-05 13:21:43.637'),
	(3, 'Admin2', 'admin22@mail.com', '7f95b733f4210c71482904eb422143f8', 'admin', '2024-12-05 15:22:41.130', '2024-12-05 15:22:41.130'),
	(4, 'Admin2', 'admin25@mail.com', '7f95b733f4210c71482904eb422143f8', 'admin', '2024-12-05 15:25:41.251', '2024-12-05 15:25:41.251'),
	(5, 'Admin2', 'admin223@mail.com', '25d55ad283aa400af464c76d713c07ad', 'admin', '2024-12-05 15:25:52.808', '2024-12-05 15:31:17.771'),
	(6, 'User1', 'user1@mail.com', '7f95b733f4210c71482904eb422143f8', 'user', '2024-12-08 08:42:22.704', '2024-12-08 08:42:22.704'),
	(7, 'User2', 'user2@mail.com', '7f95b733f4210c71482904eb422143f8', 'user', '2024-12-08 08:55:01.850', '2024-12-08 08:55:01.850'),
	(8, 'User3', 'user3@mail.com', '7f95b733f4210c71482904eb422143f8', 'user', '2024-12-08 09:57:18.757', '2024-12-08 09:57:18.757');

/*!40103 SET TIME_ZONE=IFNULL(@OLD_TIME_ZONE, 'system') */;
/*!40101 SET SQL_MODE=IFNULL(@OLD_SQL_MODE, '') */;
/*!40014 SET FOREIGN_KEY_CHECKS=IFNULL(@OLD_FOREIGN_KEY_CHECKS, 1) */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40111 SET SQL_NOTES=IFNULL(@OLD_SQL_NOTES, 1) */;
