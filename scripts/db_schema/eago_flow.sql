-- MySQL dump 10.13  Distrib 5.7.22, for macos10.13 (x86_64)
--
-- Host: localhost    Database: eago_flow
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
-- Table structure for table `categories`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `categories`
(
    `id`         int(11) unsigned NOT NULL AUTO_INCREMENT,
    `name`       varchar(100) NOT NULL DEFAULT '',
    `created_at` datetime     NOT NULL,
    `created_by` varchar(100) NOT NULL,
    `updated_at` datetime              DEFAULT NULL,
    `updated_by` varchar(100) NOT NULL DEFAULT '',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `flows`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `flows`
(
    `id`             int(11) unsigned NOT NULL AUTO_INCREMENT,
    `categories_id`  int(11) unsigned DEFAULT NULL,
    `name`           varchar(100) NOT NULL,
    `instance_title` varchar(200) NOT NULL,
    `disabled`       tinyint(1) NOT NULL DEFAULT '0',
    `description`    varchar(100) NOT NULL DEFAULT '',
    `form_id`        int(11) unsigned DEFAULT NULL,
    `first_node_id`  int(11) unsigned DEFAULT NULL,
    `created_at`     datetime     NOT NULL,
    `created_by`     varchar(100) NOT NULL,
    `updated_at`     datetime              DEFAULT NULL,
    `updated_by`     varchar(100) NOT NULL DEFAULT '',
    PRIMARY KEY (`id`),
    UNIQUE KEY `flows_id_uindex` (`id`),
    UNIQUE KEY `flows_name_uindex` (`name`),
    KEY              `flows_form_id_fk` (`form_id`),
    KEY              `flows_first_node_id_fk` (`first_node_id`),
    KEY              `flows_categories_id_fk` (`categories_id`),
    CONSTRAINT `flows_categories_id_fk` FOREIGN KEY (`categories_id`) REFERENCES `categories` (`id`),
    CONSTRAINT `flows_first_node_id_fk` FOREIGN KEY (`first_node_id`) REFERENCES `nodes` (`id`),
    CONSTRAINT `flows_form_id_fk` FOREIGN KEY (`form_id`) REFERENCES `forms` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `forms`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `forms`
(
    `id`          int(11) unsigned NOT NULL AUTO_INCREMENT,
    `name`        varchar(100) NOT NULL,
    `disabled`    tinyint(1) NOT NULL DEFAULT '0',
    `description` varchar(500) NOT NULL DEFAULT '',
    `body`        json         NOT NULL,
    `created_at`  datetime     NOT NULL,
    `created_by`  varchar(100) NOT NULL,
    `updated_at`  datetime              DEFAULT NULL,
    `updated_by`  varchar(100) NOT NULL DEFAULT '',
    PRIMARY KEY (`id`),
    UNIQUE KEY `forms_id_uindex` (`id`),
    UNIQUE KEY `forms_name_uindex` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `instances`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `instances`
(
    `id`                 int(11) unsigned NOT NULL AUTO_INCREMENT,
    `name`               varchar(200)  NOT NULL,
    `status`             int(11) NOT NULL DEFAULT '1',
    `form_id`            int(11) unsigned NOT NULL,
    `form_data`          json          NOT NULL,
    `flow_chain`         json          NOT NULL,
    `current_step`       int(11) NOT NULL DEFAULT '1',
    `assignees_required` int(11) NOT NULL,
    `current_assignees`  varchar(2000) NOT NULL DEFAULT '',
    `passed_assignees`   text          NOT NULL,
    `created_at`         datetime      NOT NULL,
    `created_by`         varchar(100)  NOT NULL,
    `updated_at`         datetime               DEFAULT NULL,
    `updated_by`         varchar(100)  NOT NULL DEFAULT '',
    PRIMARY KEY (`id`),
    UNIQUE KEY `instance_id_uindex` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `logs`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `logs`
(
    `id`          bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    `instance_id` int(11) unsigned NOT NULL,
    `result`      tinyint(1) NOT NULL,
    `content`     varchar(500) NOT NULL DEFAULT '',
    `created_at`  datetime     NOT NULL,
    `created_by`  varchar(100) NOT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `logs_id_uindex` (`id`),
    KEY           `logs_instance_id_index` (`instance_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `node_triggers`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `node_triggers`
(
    `id`         int(11) unsigned NOT NULL AUTO_INCREMENT,
    `node_id`    int(11) unsigned NOT NULL,
    `trigger_id` int(11) unsigned NOT NULL,
    `created_at` datetime     NOT NULL,
    `created_by` varchar(100) NOT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `node_triggers_id_uindex` (`id`),
    UNIQUE KEY `node_triggers_node_id_trigger_id_uindex` (`node_id`,`trigger_id`),
    KEY          `node_triggers_node_id_index` (`node_id`),
    KEY          `node_triggers_trigger_id_index` (`trigger_id`),
    CONSTRAINT `node_triggers_node_id_fk` FOREIGN KEY (`node_id`) REFERENCES `nodes` (`id`),
    CONSTRAINT `node_triggers_trigger_id_fk` FOREIGN KEY (`trigger_id`) REFERENCES `triggers` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `nodes`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `nodes`
(
    `id`                 int(11) unsigned NOT NULL AUTO_INCREMENT,
    `parent_id`          int(11) unsigned DEFAULT NULL,
    `name`               varchar(100)  NOT NULL DEFAULT '',
    `category`           int(11) NOT NULL DEFAULT '0',
    `entry_condition`    varchar(2000) NOT NULL DEFAULT '{}',
    `assignee_condition` varchar(2000) NOT NULL DEFAULT '{}',
    `visible_fields`     varchar(2000) NOT NULL DEFAULT '',
    `editable_fields`    varchar(2000) NOT NULL DEFAULT '',
    `created_at`         datetime      NOT NULL,
    `created_by`         varchar(100)  NOT NULL,
    `updated_at`         datetime               DEFAULT NULL,
    `updated_by`         varchar(100)  NOT NULL DEFAULT '',
    PRIMARY KEY (`id`),
    UNIQUE KEY `nodes_id_uindex` (`id`),
    KEY                  `nodes_parent_id_fk` (`parent_id`),
    CONSTRAINT `nodes_nodes_id_fk` FOREIGN KEY (`parent_id`) REFERENCES `nodes` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `triggers`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `triggers`
(
    `id`            int(11) unsigned NOT NULL AUTO_INCREMENT,
    `name`          varchar(100) NOT NULL DEFAULT '',
    `description`   varchar(500) NOT NULL DEFAULT '',
    `task_codename` varchar(100) NOT NULL DEFAULT '',
    `arguments`     json         NOT NULL,
    `created_at`    datetime     NOT NULL,
    `created_by`    varchar(100) NOT NULL,
    `updated_at`    datetime              DEFAULT NULL,
    `updated_by`    varchar(100) NOT NULL DEFAULT '',
    PRIMARY KEY (`id`),
    UNIQUE KEY `triggers_id_uindex` (`id`),
    UNIQUE KEY `triggers_name` (`name`)
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

-- Dump completed on 2021-07-16 18:05:05
