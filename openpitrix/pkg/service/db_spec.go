// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package service

const dbInitSql = `
CREATE TABLE IF NOT EXISTS action (
	action_id            varchar(50) not null,
	feature_id           varchar(50),
	action_name          varchar(50),
	method               varchar(50),
	description          varchar(100),
	primary key (action_id)
);

CREATE TABLE IF NOT EXISTS feature (
	feature_id           varchar(50) not null,
	module_id            varchar(50),
	feature_name         varchar(50),
	primary key (feature_id)
);

CREATE TABLE IF NOT EXISTS module (
	module_id            varchar(50) not null,
	module_name          varchar(50),
	primary key (module_id)
);

CREATE TABLE IF NOT EXISTS op_group (
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

CREATE TABLE IF NOT EXISTS role (
	role_id              varchar(50) not null,
	role_name            varchar(200),
	description          varchar(255),
	portal               varchar(50) comment ' admin,isv,dev,normal',
	create_time          timestamp,
	update_time          timestamp,
	owner                varchar(50),
	owner_path           varchar(50),
	primary key (role_id)
);

CREATE TABLE IF NOT EXISTS role_module_binding (
	binding_id           varchar(50) not null,
	role_id              varchar(50),
	module_id            varchar(50),
	data_level           varchar(50) comment 'all,department,onlyself',
	enabled_actions      text comment 'enabled_actions includes specific actions，$* means all the actions under the module_id, if not, it should be action_id strings with comma',
	create_time          timestamp,
	update_time          timestamp,
	owner                varchar(50),
	owner_path           varchar(50),
	primary key (binding_id)
);

CREATE TABLE IF NOT EXISTS user (
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

alter table action add constraint FK_Reference_10 foreign key (feature_id)
	references feature (feature_id) on delete restrict on update restrict;

alter table feature add constraint FK_Reference_8 foreign key (module_id)
	references module (module_id) on delete restrict on update restrict;

alter table op_group add constraint FK_Reference_9 foreign key (parent_group_id)
	references op_group (group_id) on delete restrict on update restrict;

alter table role_module_binding add constraint FK_Reference_11 foreign key (module_id)
	references module (module_id) on delete restrict on update restrict;

alter table role_module_binding add constraint FK_Reference_5 foreign key (role_id)
	references role (role_id) on delete restrict on update restrict;

alter table user add constraint FK_Reference_6 foreign key (group_id)
	references op_group (group_id) on delete restrict on update restrict;

alter table user add constraint FK_Reference_7 foreign key (role_id)
	references role (role_id) on delete restrict on update restrict;
`
