-- MySQL dump 10.13  Distrib 5.7.22, for macos10.13 (x86_64)
--
-- Host: localhost    Database: eago_task
-- ------------------------------------------------------
-- Server version	5.7.22

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `result_partitions`
--

DROP TABLE IF EXISTS `result_partitions`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `result_partitions`
(
    `id`        int(11) NOT NULL AUTO_INCREMENT,
    `partition` varchar(100) NOT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `result_partitions_id_uindex` (`id`),
    UNIQUE KEY `result_partitions_partition_uindex` (`partition`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `schedules`
--

DROP TABLE IF EXISTS `schedules`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `schedules`
(
    `id`         int(11) NOT NULL AUTO_INCREMENT,
    `task_id`    int(11) NOT NULL,
    `expression` varchar(20)  NOT NULL,
    `created_at` datetime     NOT NULL,
    `created_by` varchar(100) NOT NULL DEFAULT '',
    `updated_at` datetime              DEFAULT NULL,
    `updated_by` varchar(100) NOT NULL DEFAULT '',
    PRIMARY KEY (`id`),
    UNIQUE KEY `schedules_id_uindex` (`id`),
    KEY          `schedules_tasks_id_fk` (`task_id`),
    CONSTRAINT `schedules_tasks_id_fk` FOREIGN KEY (`task_id`) REFERENCES `tasks` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `tasks`
--

DROP TABLE IF EXISTS `tasks`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `tasks`
(
    `id`          int(11) NOT NULL AUTO_INCREMENT,
    `category`    int(11) NOT NULL,
    `codename`    varchar(100)  NOT NULL,
    `description` varchar(500)  NOT NULL DEFAULT '',
    `arguments`   varchar(2000) NOT NULL DEFAULT '',
    `disabled`    tinyint(1) NOT NULL DEFAULT '0',
    `created_at`  datetime      NOT NULL,
    `created_by`  varchar(100)  NOT NULL DEFAULT '',
    `updated_at`  datetime               DEFAULT NULL,
    `updated_by`  varchar(100)  NOT NULL DEFAULT '',
    PRIMARY KEY (`id`),
    UNIQUE KEY `tasks_id_uindex` (`id`),
    UNIQUE KEY `tasks_codename_uindex` (`codename`)
) ENGINE=InnoDB AUTO_INCREMENT=15 DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2021-06-23 17:55:43
