-- 管理员表
CREATE TABLE IF NOT EXISTS `sys_users` (
    `id` char(36) NOT NULL,
    `username` varchar(32) NOT NULL COMMENT '用户名',
    `password` char(60) NOT NULL COMMENT '密码',
    `nickname` varchar(32) DEFAULT NULL COMMENT '昵称',
    `avatar` varchar(255) DEFAULT NULL COMMENT '头像',
    `is_admin` tinyint(1) NOT NULL DEFAULT '1' COMMENT '是否管理员 0:否 1:是',
    `status` tinyint(1) NOT NULL DEFAULT '1' COMMENT '状态 0:禁用 1:启用',
    `last_login_at` datetime DEFAULT NULL COMMENT '最后登录时间',
    `last_login_ip` varchar(39) DEFAULT NULL COMMENT '最后登录IP',
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_username` (`username`),
    KEY `idx_status` (`status`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '管理员表';

-- 普通用户表
CREATE TABLE IF NOT EXISTS `normal_users` (
    `id` char(36) NOT NULL,
    `username` varchar(32) NOT NULL COMMENT '用户名',
    `password` char(60) NOT NULL COMMENT '密码',
    `nickname` varchar(32) DEFAULT NULL COMMENT '昵称',
    `avatar` varchar(255) DEFAULT NULL COMMENT '头像',
    `phone` char(11) DEFAULT NULL COMMENT '手机号',
    `status` tinyint(1) NOT NULL DEFAULT '1' COMMENT '状态 0:禁用 1:启用',
    `last_login_at` datetime DEFAULT NULL COMMENT '最后登录时间',
    `last_login_ip` varchar(39) DEFAULT NULL COMMENT '最后登录IP',
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_username` (`username`),
    UNIQUE KEY `uk_phone` (`phone`),
    KEY `idx_status` (`status`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '普通用户表';

-- 权限表
CREATE TABLE IF NOT EXISTS `permissions` (
    `id` char(36) NOT NULL,
    `name` varchar(50) NOT NULL COMMENT '权限名称',
    `code` varchar(50) NOT NULL COMMENT '权限编码',
    `type` tinyint(4) NOT NULL COMMENT '权限类型 1:菜单 2:按钮 3:接口',
    `parent_id` char(36) NOT NULL DEFAULT '0' COMMENT '父权限ID',
    `path` varchar(100) DEFAULT NULL COMMENT '权限路径',
    `status` tinyint(1) NOT NULL DEFAULT '1' COMMENT '状态 0:禁用 1:启用',
    `remark` varchar(255) DEFAULT NULL COMMENT '备注',
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_code` (`code`),
    KEY `idx_parent_status` (`parent_id`, `status`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '权限表';

-- 用户权限关联表
CREATE TABLE IF NOT EXISTS `user_permissions` (
    `id` char(36) NOT NULL,
    `user_id` char(36) NOT NULL COMMENT '用户ID',
    `user_type` tinyint(1) NOT NULL COMMENT '用户类型 0:普通用户 1:管理员',
    `permission_id` char(36) NOT NULL COMMENT '权限ID',
    `operator_id` char(36) NOT NULL COMMENT '操作人ID',
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_user_perm` (`user_id`, `user_type`, `permission_id`),
    KEY `idx_permission` (`permission_id`),
    CONSTRAINT `fk_up_permission` FOREIGN KEY (`permission_id`) REFERENCES `permissions` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT `fk_up_sys_user` FOREIGN KEY (`user_id`) REFERENCES `sys_users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT `fk_up_normal_user` FOREIGN KEY (`user_id`) REFERENCES `normal_users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT `fk_up_operator` FOREIGN KEY (`operator_id`) REFERENCES `sys_users` (`id`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '用户权限关联表';

-- 角色表
CREATE TABLE IF NOT EXISTS `roles` (
    `id` char(36) NOT NULL,
    `name` varchar(32) NOT NULL COMMENT '角色名称',
    `code` varchar(32) NOT NULL COMMENT '角色编码',
    `status` tinyint(1) NOT NULL DEFAULT '1' COMMENT '状态 0:禁用 1:启用',
    `remark` varchar(255) DEFAULT NULL COMMENT '备注',
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_code` (`code`),
    KEY `idx_status` (`status`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '角色表';

-- 用户角色关联表
CREATE TABLE IF NOT EXISTS `user_roles` (
    `id` char(36) NOT NULL,
    `user_id` char(36) NOT NULL COMMENT '用户ID',
    `user_type` tinyint(1) NOT NULL COMMENT '用户类型 0:普通用户 1:管理员',
    `role_id` char(36) NOT NULL COMMENT '角色ID',
    `operator_id` char(36) NOT NULL COMMENT '操作人ID',
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_user_role` (`user_id`, `user_type`, `role_id`),
    KEY `idx_role` (`role_id`),
    CONSTRAINT `fk_ur_role` FOREIGN KEY (`role_id`) REFERENCES `roles` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT `fk_ur_sys_user` FOREIGN KEY (`user_id`) REFERENCES `sys_users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT `fk_ur_normal_user` FOREIGN KEY (`user_id`) REFERENCES `normal_users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT `fk_ur_operator` FOREIGN KEY (`operator_id`) REFERENCES `sys_users` (`id`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '用户角色关联表';

-- 系统日志表
CREATE TABLE IF NOT EXISTS `system_logs` (
    `id` char(36) NOT NULL,
    `level` varchar(10) NOT NULL COMMENT '日志级别',
    `content` text NOT NULL COMMENT '日志内容',
    `trace_id` varchar(32) DEFAULT NULL COMMENT '追踪ID',
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    KEY `idx_created_at` (`created_at`),
    KEY `idx_level` (`level`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '系统日志表';

-- 操作日志表
CREATE TABLE IF NOT EXISTS `operation_logs` (
    `id` char(36) NOT NULL,
    `user_id` char(36) NOT NULL COMMENT '用户ID',
    `user_type` tinyint(1) NOT NULL COMMENT '用户类型 0:普通用户 1:管理员',
    `operation` varchar(32) NOT NULL COMMENT '操作类型',
    `method` varchar(10) NOT NULL COMMENT '请求方法',
    `path` varchar(100) NOT NULL COMMENT '请求路径',
    `params` text COMMENT '请求参数',
    `ip` varchar(39) DEFAULT NULL COMMENT '操作IP',
    `status` int(11) NOT NULL COMMENT '操作状态',
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    KEY `idx_user_created` (`user_id`, `user_type`, `created_at`),
    KEY `idx_operation` (`operation`),
    CONSTRAINT `fk_ol_sys_user` FOREIGN KEY (`user_id`) REFERENCES `sys_users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT `fk_ol_normal_user` FOREIGN KEY (`user_id`) REFERENCES `normal_users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '操作日志表';
