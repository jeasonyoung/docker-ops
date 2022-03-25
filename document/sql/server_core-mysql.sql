-- --------------------------------------------------------------------------------------------------------------------
-- 创建数据库
create database if not exists ops_server_core default charset utf8mb4 collate utf8mb4_general_ci;
-- --------------------------------------------------------------------------------------------------------------------
set names 'utf8mb4';
use ops_server_core;
-- --------------------------------------------------------------------------------------------------------------------
-- 10.字典数据表
drop table if exists tbl_dict_data;
-- 8.部署任务服务器表
drop table if exists tbl_ops_deploy_task_server;
-- 7.部署任务表
drop table if exists tbl_ops_deploy_task;
-- 6.部署表
drop table if exists tbl_ops_deploy;
-- 5.镜像表
drop table if exists tbl_ops_repository_images;
-- 3.删除表顺序
drop table if exists tbl_ops_group_servers;
-- --------------------------------------------------------------------------------------------------------------------
-- 1.服务器表
drop table if exists tbl_ops_server;
create table tbl_ops_server (
    `id`    bigint unsigned  not null   comment '服务器ID',
    `name`  varchar(128)     not null   comment '服务器名称',
    `title` varchar(255) default null   comment '服务器标题(别名)',

    `auth_code`  varchar(64) not null   comment '服务器授权码(用于代理连接标识)',

    `inner_ip_addr` varchar(64) default null comment '内网IP地址',
    `outer_ip_addr` varchar(64) default null comment '外网IP地址',

    `status` tinyint not null default 1 comment '状态(-1:删除,0:停用,1:启用)',

    `last_ping_time` datetime     default null  comment '最后心跳时间',
    `last_ping_msg`  varchar(255) default null  comment '最后心跳消息',

    `create_time` timestamp   not null default current_timestamp comment '创建时间',
    `update_time` timestamp   not null default current_timestamp on update current_timestamp comment '更新时间',

    constraint `pk_ops_server` primary key(`id`),
    constraint `uk_ops_server_name` unique key(`name`),
    constraint `uk_ops_server_abbr` unique key(`auth_code`)
) engine=InnoDB default charset=utf8mb4 comment '服务器表';
-- --------------------------------------------------------------------------------------------------------------------
-- 2.服务器分组表
drop table if exists tbl_ops_group;
create table tbl_ops_group (
    `id`    bigint unsigned  not null   comment '分组ID',
    `name`  bigint unsigned  not null   comment '分组名称',

    `remark` varchar(255) default null  comment '分组描述',

    `status` tinyint not null default 1 comment '状态(-1:删除,0:停用,1:启用)',

    `create_time` timestamp   not null default current_timestamp comment '创建时间',
    `update_time` timestamp   not null default current_timestamp on update current_timestamp comment '更新时间',

    constraint `pk_ops_group` primary key(`id`),
    constraint `uk_ops_group_name` unique key(`name`)
) engine=InnoDB default charset=utf8mb4 comment '服务器分组表';
-- --------------------------------------------------------------------------------------------------------------------
-- 3.服务器分组关联表
drop table if exists tbl_ops_group_servers;
create table tbl_ops_group_servers (
    `group_id`  bigint unsigned  not null   comment '分组ID',
    `server_id` bigint unsigned  not null   comment '服务器ID',

    constraint `pk_ops_group_servers` primary key(`group_id`,`server_id`),
    constraint `fk_ops_group_servers_g` foreign key(group_id) references tbl_ops_group(`id`),
    constraint `fk_ops_group_servers_s` foreign key(`server_id`) references tbl_ops_server(`id`)
) engine=InnoDB default charset=utf8mb4 comment '服务器分组关联表';
-- --------------------------------------------------------------------------------------------------------------------
-- 4.镜像仓库表
drop table if exists tbl_ops_repository;
create table tbl_ops_repository (
    `id`    bigint unsigned  not null   comment '镜像仓库ID',
    `name`  varchar(128)     not null   comment '镜像仓库名称',
    `remark` varchar(255) default null  comment '镜像仓库描述',
    `addr`  varchar(1024) default null  comment '镜像仓库地址',

    `tmp_login_cmd`     varchar(2048) default null comment '临时登录指令',
    `tmp_expire_time`   datetime      default null comment '临时指令到期时间',

    `status` tinyint not null default 1 comment '状态(-1:删除,0:停用,1:启用)',

    `create_time` timestamp   not null default current_timestamp comment '创建时间',
    `update_time` timestamp   not null default current_timestamp on update current_timestamp comment '更新时间',

    constraint `pk_ops_repository` primary key(`id`),
    constraint `uk_ops_repository_name` unique key(`name`)
) engine=InnoDB default charset=utf8mb4 comment '镜像仓库表';
-- --------------------------------------------------------------------------------------------------------------------
-- 5.镜像表
drop table if exists tbl_ops_repository_images;
create table tbl_ops_repository_images (
    `id`    bigint unsigned  not null   comment '镜像ID',
    `title` varchar(128) default null   comment '镜像标题',

    `image_org`     varchar(255) not null   comment '镜像组织',
    `image_name`    varchar(128) not null   comment '镜像名称',
    `image_tag`     varchar(64)  not null   comment '镜像标签(版本)',

    `status` tinyint not null default 1 comment '状态(-1:删除,0:停用,1:启用)',

    `repository_id` bigint unsigned  not null   comment '所属镜像仓库ID',

    `create_time` timestamp   not null default current_timestamp comment '创建时间',
    `update_time` timestamp   not null default current_timestamp on update current_timestamp comment '更新时间',

    constraint `pk_ops_repository_images` primary key(`id`),
    constraint `uk_ops_repository_images_all` unique key(`repository_id`,`image_org`,`image_name`,`image_tag`),
    constraint `fk_ops_repository_images_repository` foreign key(`repository_id`) references tbl_ops_repository(`id`)
) engine=InnoDB default charset=utf8mb4 comment '镜像表';
-- --------------------------------------------------------------------------------------------------------------------
-- 6.部署表
drop table if exists tbl_ops_deploy;
create table tbl_ops_deploy (
    `id`     bigint unsigned  not null  comment '部署ID',
    `name`   varchar(128)     not null  comment '部署名称',
    `remark` varchar(255) default null  comment '部署描述',

    `group_id`      bigint unsigned  not null   comment '服务器分组ID',
    `repository_id` bigint unsigned  not null   comment '所属镜像仓库ID',

    `status` tinyint not null default 1 comment '状态(-1:删除,0:停用,1:启用)',

    `create_time` timestamp   not null default current_timestamp comment '创建时间',
    `update_time` timestamp   not null default current_timestamp on update current_timestamp comment '更新时间',

    constraint `pk_ops_deploy` primary key(`id`),
    constraint `uk_ops_deploy_name` unique key(`name`),
    constraint `fk_ops_deploy_g` foreign key(`group_id`) references tbl_ops_group(`id`),
    constraint `fk_ops_deploy_r` foreign key(`repository_id`) references tbl_ops_repository(`id`)
) engine=InnoDB default charset=utf8mb4 comment '部署表';
-- --------------------------------------------------------------------------------------------------------------------
-- 7.部署任务表
drop table if exists tbl_ops_deploy_task;
create table tbl_ops_deploy_task (
    `id`    bigint unsigned  not null   comment '部署任务ID',
    `name`  varchar(128)     not null   comment '部署任务名称',
    `remark` varchar(255) default null  comment '部署任务描述',

    `deploy_image_id` bigint unsigned  not null comment '部署镜像ID',

    `progress`  tinyint not null default 0   comment '进度(-1:部署失败,0:未部署,1:部署中,2:部署完成)',
    `fail_msg`  varchar(1024) default null   comment '失败消息',

    `status` tinyint not null default 1 comment '状态(-1:删除,0:停用,1:启用)',

    `create_time` timestamp   not null default current_timestamp comment '创建时间',
    `update_time` timestamp   not null default current_timestamp on update current_timestamp comment '更新时间',

    constraint `pk_ops_deploy_task` primary key(`id`),
    constraint `uk_ops_deploy_task_name` unique key(`name`),
    constraint `fk_ops_deploy_task_image` foreign key(`deploy_image_id`) references tbl_ops_repository_images(`id`)
) engine=InnoDB default charset=utf8mb4 comment '部署任务表';
-- --------------------------------------------------------------------------------------------------------------------
-- 8.部署任务服务器表
drop table if exists tbl_ops_deploy_task_server;
create table tbl_ops_deploy_task_server (
    `id`    bigint unsigned  not null   comment '部署任务服务器ID',

    `task_id`   bigint unsigned  not null   comment '所属部署任务ID',
    `server_id` bigint unsigned  not null   comment '所属服务器ID',

    `progress` tinyint not null default 0   comment '进度(-1:部署失败,0:未部署,1:部署中,2:部署完成)',
    `fail_msg` varchar(1024) default null   comment '失败消息',

    `create_time` timestamp   not null default current_timestamp comment '创建时间',
    `update_time` timestamp   not null default current_timestamp on update current_timestamp comment '更新时间',

    constraint `pk_ops_deploy_task_server` primary key(`id`),
    constraint `uk_ops_deploy_task_server_all` unique key(`task_id`,`server_id`),
    constraint `fk_ops_deploy_task_server_t` foreign key(`task_id`) references tbl_ops_deploy_task(`id`),
    constraint `fk_ops_deploy_task_server_s` foreign key(`server_id`) references tbl_ops_server(`id`)
) engine=InnoDB default charset=utf8mb4 comment '部署任务服务器表';
-- --------------------------------------------------------------------------------------------------------------------
-- 9.字典类型表
drop table if exists tbl_dict_type;
create table tbl_dict_type (
    `id`          bigint unsigned not null comment '字典ID',
    `name`        varchar(64)  not null comment '字典名称',
    `type`        varchar(128) not null comment '字典类型',
    `remark`      varchar(255)          default null comment '字典备注',

    `status`      tinyint               default 1 comment '状态(-1:删除,0:停用,1:启用)',
    `create_time` timestamp    not null default current_timestamp comment '创建时间',
    `update_time` timestamp    not null default current_timestamp on update current_timestamp comment '更新时间',

    constraint `pk_dict_type` primary key (`id`),
    constraint `uk_dict_type_type` unique key (`type`)
) engine=InnoDB default charset=utf8mb4 comment '字典类型表';
-- --------------------------------------------------------------------------------------------------------------------
-- 10.字典数据表
drop table if exists tbl_dict_data;
create table tbl_dict_data (
    `id`        bigint unsigned not null    comment '字典数据ID',
    `code`        int unsigned default 0    comment '字典代码(排序)',
    `label`        varchar(128) not null    comment '字典标签',
    `value`        varchar(255) default ''  comment '字典键值',
    `is_default` tinyint unsigned default 0 comment '是否默认(0:否,1:是)',

    `type`        varchar(128) not null     comment '字典类型',

    `css_class`  varchar(128) default null  comment '样式属性',
    `list_class` varchar(128) default null  comment '表格回显样式',

    `remark`  varchar(255) default null     comment '字典数据备注',

    `status` tinyint default 1  comment '状态(-1:删除,0:停用,1:启用)',

    `create_time` timestamp not null default current_timestamp comment '创建时间',
    `update_time` timestamp not null default current_timestamp on update current_timestamp comment '更新时间',

    constraint `pk_dict_data` primary key (`id`),
    constraint `uk_dict_data_value` unique key (`type`,`value`),
    constraint `fk_dict_data_t` foreign key(`type`) references tbl_dict_type(`type`)
) engine=InnoDB default charset=utf8mb4 comment '字典数据表';
-- --------------------------------------------------------------------------------------------------------------------
