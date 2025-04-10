set NAMES utf8mb4;
set FOREIGN_KEY_CHECKS = 0;

CREATE DATABASE IF NOT EXISTS coauth DEFAULT CHARSET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE coauth;

CREATE TABLE app (
    id VARCHAR(36) NOT NULL,          -- 客户端ID
    user_id VARCHAR(36) NOT NULL,      -- 用户ID
    `secret` VARCHAR(64) NOT NULL,        -- 客户端密钥
    `name` VARCHAR(32) NOT NULL,          -- 客户端名称
    redirect_url VARCHAR(256) NOT NULL,  -- 回调地址
    home_page VARCHAR(256) NOT NULL,               -- 主页地址
    intro VARCHAR(64),               -- 简介
    scopes VARCHAR(256),              -- 权限范围
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- 创建时间
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP, -- 更新时间
    deleted_at TIMESTAMP NULL DEFAULT NULL, -- 删除时间
    PRIMARY KEY (id),
    UNIQUE KEY (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

INSERT INTO `coauth`.`user`
(`id`,
`username`,
`password`,
`salt`,
`gender`,
`nickname`,
`avatar`,
`email`,
`phone`)
VALUES
(uuid(),'admin1','c886db7436266ee8afb0a8c40f8d65a4aebec45ad47c5f4cb05fca46c1f6ecef','saltsalt',1,'admin1',null,null,null);

CREATE TABLE user (
    id VARCHAR(36) NOT NULL,          -- 用户ID
    username VARCHAR(32) NOT NULL,      -- 用户名
    `password` VARCHAR(64) NOT NULL,      -- 密码
    salt VARCHAR(16) NOT NULL,          -- 盐值
    gender TINYINT(1) DEFAULT 0 NOT NULL COMMENT '性别：0-保密，1-男，2-女',
    nickname VARCHAR(32),      -- 昵称
    avatar VARCHAR(256),                 -- 头像
    email VARCHAR(128),                  -- 邮箱
    phone VARCHAR(20),                   -- 手机号
    `status` TINYINT(1) DEFAULT 1 NOT NULL COMMENT '状态：1-正常，0-禁用',
    last_login_at bigint NULL,        -- 最后登录时间
    last_login_ip VARCHAR(16),           -- 最后登录IP
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- 创建时间
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP, -- 更新时间
    deleted_at TIMESTAMP NULL DEFAULT NULL, -- 删除时间
    PRIMARY KEY (id),
    UNIQUE KEY (username),
    Index idx_status (`status`),
    Index idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE auth_code (
    id bigint NOT NULL AUTO_INCREMENT,
    code VARCHAR(36) NOT NULL,        -- 授权码
    client_id VARCHAR(36) NOT NULL,      -- 客户端ID
    user_id VARCHAR(36) NOT NULL,        -- 用户ID
    redirect_url VARCHAR(256) NOT NULL,  -- 回调地址
    expires_at TIMESTAMP NOT NULL,       -- 过期时间
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- 创建时间
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP, -- 更新时间
    deleted_at TIMESTAMP NULL DEFAULT NULL, -- 删除时间
    PRIMARY KEY (id),
    INDEX idx_code (code),
    INDEX idx_client_id (client_id),
    INDEX idx_user_id (user_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE token (
    id bigint NOT NULL AUTO_INCREMENT,
    access_token VARCHAR(36) NOT NULL, -- 访问令牌
    refresh_token VARCHAR(36) NOT NULL,            -- 刷新令牌
    client_id VARCHAR(36) NOT NULL,       -- 客户端ID
    user_id VARCHAR(36) NOT NULL,         -- 用户ID
    expires_at TIMESTAMP NOT NULL,        -- 过期时间
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- 创建时间
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP, -- 更新时间
    deleted_at TIMESTAMP NULL DEFAULT NULL, -- 删除时间
    PRIMARY KEY (id),
    INDEX idx_client_id (client_id),
    INDEX idx_user_id (user_id),
    INDEX idx_uid_cli (user_id, client_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
