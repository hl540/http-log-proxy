/*
 Navicat Premium Dump SQL

 Source Server         : data
 Source Server Type    : SQLite
 Source Server Version : 3045000 (3.45.0)
 Source Schema         : main

 Target Server Type    : SQLite
 Target Server Version : 3045000 (3.45.0)
 File Encoding         : 65001

 Date: 24/04/2025 10:49:48
*/

PRAGMA foreign_keys = false;

-- ----------------------------
-- Table structure for tb_app
-- ----------------------------
DROP TABLE IF EXISTS "tb_app";
CREATE TABLE "tb_app" (
  "id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
  "key" TEXT,
  "name" TEXT,
  "target" TEXT,
  "create_at" integer,
  "update_at" integer,
  UNIQUE ("key" ASC)
);

-- ----------------------------
-- Table structure for tb_http_log
-- ----------------------------
DROP TABLE IF EXISTS "tb_http_log";
CREATE TABLE "tb_http_log" (
  "request_id" text NOT NULL,
  "app_id" INTEGER,
  "app_key" TEXT,
  "request_url" TEXT,
  "request_method" TEXT,
  "request_header" TEXT,
  "request_body" TEXT,
  "response_code" integer,
  "response_header" TEXT,
  "response_body" TEXT,
  "create_at" integer,
  PRIMARY KEY ("request_id")
);

-- ----------------------------
-- Auto increment value for tb_app
-- ----------------------------
UPDATE "sqlite_sequence" SET seq = 22 WHERE name = 'tb_app';

PRAGMA foreign_keys = true;
