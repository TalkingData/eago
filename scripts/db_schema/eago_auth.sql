-- MySQL dump 10.13  Distrib 5.7.22, for macos10.13 (x86_64)
--
-- Host: localhost    Database: eago_auth
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
-- Table structure for table `departments`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `departments`
(
    `id`         int(11) NOT NULL AUTO_INCREMENT,
    `name`       varchar(100) NOT NULL,
    `parent_id`  int(11) DEFAULT NULL,
    `created_at` datetime     NOT NULL,
    `updated_at` datetime DEFAULT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `departments_id_uindex` (`id`),
    UNIQUE KEY `departments_name_uindex` (`name`),
    KEY          `departments_departments_id_fk` (`parent_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `groups`
--

DROP TABLE IF EXISTS `groups`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `groups`
(
    `id`          int(11) NOT NULL AUTO_INCREMENT,
    `name`        varchar(100) NOT NULL,
    `description` varchar(500) NOT NULL DEFAULT '',
    `created_at`  datetime     NOT NULL,
    `updated_at`  datetime              DEFAULT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `groups_id_uindex` (`id`),
    UNIQUE KEY `groups_name_uindex` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `products`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `products`
(
    `id`          int(11) NOT NULL AUTO_INCREMENT,
    `name`        varchar(200) NOT NULL,
    `alias`       varchar(200) NOT NULL,
    `disabled`    tinyint(1) NOT NULL DEFAULT '0',
    `description` varchar(500) NOT NULL DEFAULT '',
    `created_at`  datetime     NOT NULL,
    `updated_at`  datetime              DEFAULT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `products_id_uindex` (`id`),
    UNIQUE KEY `products_name_uindex` (`name`),
    UNIQUE KEY `products_alias_uindex` (`alias`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `roles`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `roles`
(
    `id`          int(11) NOT NULL AUTO_INCREMENT,
    `name`        varchar(100) NOT NULL,
    `description` varchar(500) NOT NULL DEFAULT '',
    PRIMARY KEY (`id`),
    UNIQUE KEY `roles_id_uindex` (`id`),
    UNIQUE KEY `roles_name_uindex` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `user_departments`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `user_departments`
(
    `id`            int(11) NOT NULL AUTO_INCREMENT,
    `user_id`       int(11) NOT NULL,
    `department_id` int(11) NOT NULL,
    `is_owner`      tinyint(1) NOT NULL DEFAULT '0',
    `joined_at`     datetime NOT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `user_departments_id_uindex` (`id`),
    UNIQUE KEY `user_departments_user_id_uindex` (`user_id`),
    KEY             `user_departments_department_id_index` (`department_id`),
    CONSTRAINT `user_departments_departments_id_fk` FOREIGN KEY (`department_id`) REFERENCES `departments` (`id`),
    CONSTRAINT `user_departments_users_id_fk` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `user_groups`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `user_groups`
(
    `id`        int(11) NOT NULL AUTO_INCREMENT,
    `user_id`   int(11) NOT NULL,
    `group_id`  int(11) NOT NULL,
    `is_owner`  tinyint(1) NOT NULL DEFAULT '0',
    `joined_at` datetime NOT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `user_groups_id_uindex` (`id`),
    UNIQUE KEY `user_groups_user_id_group_id_uindex` (`user_id`,`group_id`),
    KEY         `user_groups_user_id_index` (`user_id`),
    KEY         `user_groups_group_id_index` (`group_id`),
    CONSTRAINT `user_groups_groups_id_fk` FOREIGN KEY (`group_id`) REFERENCES `groups` (`id`),
    CONSTRAINT `user_groups_users_id_fk` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `user_products`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `user_products`
(
    `id`         int(11) NOT NULL AUTO_INCREMENT,
    `user_id`    int(11) NOT NULL,
    `product_id` int(11) NOT NULL,
    `is_owner`   tinyint(1) NOT NULL DEFAULT '0',
    `joined_at`  datetime NOT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `user_products_id_uindex` (`id`),
    UNIQUE KEY `user_products_user_id_product_id_uindex` (`user_id`,`product_id`),
    KEY          `user_products_user_id_index` (`user_id`),
    KEY          `user_products_product_id_index` (`product_id`),
    CONSTRAINT `user_products_products_id_fk` FOREIGN KEY (`product_id`) REFERENCES `products` (`id`),
    CONSTRAINT `user_products_users_id_fk` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `user_roles`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `user_roles`
(
    `id`        int(11) NOT NULL AUTO_INCREMENT,
    `user_id`   int(11) NOT NULL,
    `role_id`   int(11) NOT NULL,
    `joined_at` datetime NOT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `user_roles_id_uindex` (`id`),
    UNIQUE KEY `user_roles_user_id_role_id_pk` (`user_id`,`role_id`),
    KEY         `user_roles_roles_id_fk` (`role_id`),
    CONSTRAINT `user_roles_roles_id_fk` FOREIGN KEY (`role_id`) REFERENCES `roles` (`id`),
    CONSTRAINT `user_roles_users_id_fk` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `users`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `users`
(
    `id`           int(11) NOT NULL AUTO_INCREMENT,
    `username`     varchar(100) NOT NULL,
    `password`     varchar(200) NOT NULL DEFAULT '',
    `email`        varchar(100) NOT NULL DEFAULT '',
    `phone`        varchar(20)  NOT NULL DEFAULT '',
    `is_superuser` tinyint(1) NOT NULL DEFAULT '0',
    `disabled`     tinyint(1) NOT NULL DEFAULT '0',
    `last_login`   datetime              DEFAULT NULL,
    `created_at`   datetime     NOT NULL,
    `updated_at`   datetime              DEFAULT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `users_id_uindex` (`id`),
    UNIQUE KEY `users_username_uindex` (`username`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2021-01-27 11:44:54
