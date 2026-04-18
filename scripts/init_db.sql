-- 创建数据库（如果不存在）
CREATE DATABASE IF NOT EXISTS shortlink DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- 使用数据库
USE shortlink;

-- 创建短链表
CREATE TABLE IF NOT EXISTS short_links (
    id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '主键ID',
    code VARCHAR(10) NOT NULL UNIQUE COMMENT '短码',
    original_url VARCHAR(2048) NOT NULL COMMENT '原始URL',
    short_url VARCHAR(256) COMMENT '短链URL',
    is_enabled BOOLEAN DEFAULT TRUE COMMENT '是否启用',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    INDEX idx_code (code),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='短链表';

-- 创建用户（可选）
CREATE USER IF NOT EXISTS 'shortlink'@'%' IDENTIFIED BY 'shortlink123';
GRANT ALL PRIVILEGES ON shortlink.* TO 'shortlink'@'%';
FLUSH PRIVILEGES;
