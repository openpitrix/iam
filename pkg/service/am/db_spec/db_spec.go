// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package db_spec

import "regexp"

var (
	reUid       = regexp.MustCompile(`^[a-zA-Z0-9-_]{2,64}$`)
	reGid       = regexp.MustCompile(`^[a-zA-Z0-9-_]{2,64}$`)
	reGroupPath = regexp.MustCompile(`^[a-zA-Z0-9_.-]{2,255}$`)
)

const (
	UserTableName      = "user"
	UserPrimaryKeyName = "user_id"

	UserGroupTableName      = "user_group"
	UserGroupPrimaryKeyName = "group_id"

	UserGroupBindingTableName      = "user_group_binding"
	UserGroupBindingPrimaryKeyName = "id"
)

var TableMap = map[string]struct{ Name, PrimaryKey, Sql string }{
	UserTableName: {
		Name:       UserTableName,
		PrimaryKey: UserPrimaryKeyName,
		Sql: `CREATE TABLE IF NOT EXISTS ` + UserTableName + ` (
			user_id      varchar(50) not null,
			user_name    varchar(50),
			email        varchar(50),
			phone_number varchar(50),
			description  varchar(1000),
			password     varchar(128),
			status       varchar(10),
			create_time  timestamp,
			update_time  timestamp default CURRENT_TIMESTAMP,
			status_time  timestamp,
			extra        json,

			primary key(` + UserPrimaryKeyName + `)
		);`,
	},
	UserGroupTableName: {
		Name:       UserGroupTableName,
		PrimaryKey: UserGroupPrimaryKeyName,
		Sql: `CREATE TABLE IF NOT EXISTS ` + UserGroupTableName + ` (
			group_id        varchar(50) not null,
			group_path      varchar(255),
			group_name      varchar(50),
			description     varchar(1000),
			status          varchar(10),
			create_time     timestamp,
			update_time  timestamp default CURRENT_TIMESTAMP,
			status_time     timestamp,
			extra           json,

			parent_group_id varchar(50),
			level           int,

			primary key(` + UserGroupPrimaryKeyName + `)
		);`,
	},
	UserGroupBindingTableName: {
		Name:       UserGroupBindingTableName,
		PrimaryKey: UserGroupBindingPrimaryKeyName,
		Sql: `CREATE TABLE IF NOT EXISTS ` + UserGroupBindingTableName + ` (
			id          varchar(50) not null,
			user_id     varchar(50),
			group_id    varchar(50),
			create_time timestamp default CURRENT_TIMESTAMP,

			primary key(` + UserGroupBindingPrimaryKeyName + `)
		);`,
	},
}

const sql = `
 -- create database am;
 -- use am;

/*==============================================================*/
/* Table: action                                                */
/*==============================================================*/
create table action
(
   action_id            varchar(50) not null,
   feature_id           varchar(50),
   action_name          varchar(50),
   method               varchar(50),
   description          varchar(100),
   url                  varchar(500),
   url_method           varchar(20)
);

alter table action
   add primary key (action_id);

/*==============================================================*/
/* Table: action2                                               */
/*==============================================================*/
create table action2
(
   module_id            varchar(50),
   module_name          varchar(50),
   feature_id           varchar(50),
   feature_name         varchar(50),
   action_id            varchar(50) not null,
   action_name          varchar(50),
   method               varchar(50),
   description          varchar(100),
   url_method           varchar(20),
   url                  varchar(500)
);

alter table action2
   add primary key (action_id);

/*==============================================================*/
/* Table: feature                                               */
/*==============================================================*/
create table feature
(
   feature_id           varchar(50) not null,
   module_id            varchar(50),
   feature_name         varchar(50)
);

alter table feature
   add primary key (feature_id);

/*==============================================================*/
/* Table: module                                                */
/*==============================================================*/
create table module
(
   module_id            varchar(50) not null,
   module_name          varchar(50)
);

alter table module
   add primary key (module_id);

/*==============================================================*/
/* Table: role                                                  */
/*==============================================================*/
create table role
(
   role_id              varchar(50) not null,
   role_name            varchar(200),
   description          varchar(255),
   portal               varchar(50) comment ' admin,isv,dev,normal',
   create_time          timestamp,
   update_time          timestamp,
   owner                varchar(50),
   owner_path           varchar(50)
);

alter table role
   add primary key (role_id);

/*==============================================================*/
/* Table: role_module_binding                                   */
/*==============================================================*/
create table role_module_binding
(
   id                   varchar(50) not null,
   role_id              varchar(50),
   module_id            varchar(50),
   data_level           varchar(50) comment 'all,department,onlyself',
   enabled_actions      text comment 'enabled_actions includes specific actions��$* means all the actions under the module_id, if not, it should be action_id strings with comma',
   create_time          timestamp,
   update_time          timestamp,
   owner                varchar(50)
);

alter table role_module_binding
   add primary key (id);

/*==============================================================*/
/* Table: user_role_binding                                     */
/*==============================================================*/
create table user_role_binding
(
   id                   varchar(50) not null,
   user_id              varchar(50),
   role_id              varchar(50)
);

alter table user_role_binding
   add primary key (id);

-- alter table action add constraint FK_Reference_10 foreign key (feature_id)
--       references feature (feature_id) on delete restrict on update restrict;
--
-- alter table feature add constraint FK_Reference_8 foreign key (module_id)
--       references module (module_id) on delete restrict on update restrict;
--
-- alter table role_module_binding add constraint FK_Reference_11 foreign key (module_id)
--       references module (module_id) on delete restrict on update restrict;
--
-- alter table role_module_binding add constraint FK_Reference_5 foreign key (role_id)
--       references role (role_id) on delete restrict on update restrict;
--
-- alter table user_role_binding add constraint FK_Reference_19 foreign key (user_id)
--       references user (user_id) on delete restrict on update restrict;
--
-- alter table user_role_binding add constraint FK_Reference_20 foreign key (role_id)
--       references role (role_id) on delete restrict on update restrict;
`
