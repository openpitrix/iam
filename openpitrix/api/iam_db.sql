/*==============================================================*/
/* DBMS name:      MySQL 5.0                                    */
/* Created on:     2018/12/21 10:10:07                          */
/*==============================================================*/
CREATE database  IAM;
use IAM;

drop table if exists Action;

drop table if exists Biz_App;

drop table if exists Biz_Runtime;

drop table if exists Role;

drop table if exists feature;

drop table if exists module;

drop table if exists op_group;

drop table if exists role_feature_binding;

drop table if exists user;

/*==============================================================*/
/* Table: Action                                                */
/*==============================================================*/
create table Action
(
   action_id            varchar(50) not null,
   feature_id           varchar(50),
   action_name          varchar(50),
   method               varchar(50),
   description          varchar(100),
   primary key (action_id)
);

/*==============================================================*/
/* Table: Biz_App                                               */
/*==============================================================*/
create table Biz_App
(
   AppID                varchar(50) not null,
   AppName              varchar(50),
   owner                varchar(50),
   owner_path           varchar(50),
   primary key (AppID)
);

/*==============================================================*/
/* Table: Biz_Runtime                                           */
/*==============================================================*/
create table Biz_Runtime
(
   RuntimeID            varchar(50) not null,
   RuntimeName          varchar(50),
   owner                varchar(50),
   owmer_path           varchar(50),
   primary key (RuntimeID)
);

/*==============================================================*/
/* Table: Role                                                  */
/*==============================================================*/
create table Role
(
   role_id              varchar(50) not null,
   role_name            varchar(200),
   description          varchar(255),
   portal               varchar(50) comment '管理员，ISV，开发者，普通用户',
   create_time          timestamp,
   update_time          timestamp,
   owner                varchar(50),
   owner_path           varchar(50),
   primary key (role_id)
);

/*==============================================================*/
/* Table: feature                                               */
/*==============================================================*/
create table feature
(
   feature_id           varchar(50) not null,
   module_id            varchar(50),
   feature_name         varchar(50),
   primary key (feature_id)
);

/*==============================================================*/
/* Table: module                                                */
/*==============================================================*/
create table module
(
   module_id            varchar(50) not null,
   module_name          varchar(50),
   primary key (module_id)
);

/*==============================================================*/
/* Table: op_group                                              */
/*==============================================================*/
create table op_group
(
   group_id             varchar(50) not null,
   group_name           varchar(50),
   parent_group_id      varchar(50),
   group_path           varchar(255),
   level                int,
   seq_order            int,
   create_time          timestamp,
   update_time          timestamp,
   owner                varchar(50),
   owner_path           varchar(50),
   primary key (group_id)
);

/*==============================================================*/
/* Table: role_feature_binding                                  */
/*==============================================================*/
create table role_feature_binding
(
   binding_id           varchar(50) not null,
   role_id              varchar(50),
   feature_id           varchar(50),
   data_level           varchar(50) comment '1 所有数据 2 本部门内 3 仅个人',
   enabled_actions      text comment 'enabled 里面对应具体的action
            $* 或者什么特殊字符表示这个module_id 下全部action
            如果不是全部， 就是action1,action2,action3 这种逗号分割',
   create_time          timestamp,
   update_time          timestamp,
   owner                varchar(50),
   owner_path           varchar(50),
   primary key (binding_id)
);

/*==============================================================*/
/* Table: user                                                  */
/*==============================================================*/
create table user
(
   user_id              varchar(50) not null,
   group_id             varchar(50),
   role_id              varchar(50),
   user_name            varchar(50),
   position             varchar(50),
   email                varchar(50),
   phone_number         varchar(50),
   password             varchar(50),
   old_password         varchar(50),
   description          varchar(200),
   status               varchar(10),
   create_time          timestamp,
   status_time          timestamp,
   update_time          timestamp,
   owner                varchar(50),
   owner_path           varchar(50),
   primary key (user_id)
);
 