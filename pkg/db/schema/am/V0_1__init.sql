CREATE TABLE IF NOT EXISTS `enable_action_bundle` (
	enable_id            varchar(50) NOT NULL,
	bind_id              varchar(50) NOT NULL,
	action_bundle_id     varchar(50) NOT NULL,
	create_time          timestamp   NOT NULL DEFAULT CURRENT_TIMESTAMP,
	PRIMARY KEY(enable_id)
);

CREATE INDEX enable_action_bundle_bind_id_idx
  ON enable_action_bundle (bind_id);
CREATE INDEX enable_action_bundle_action_bundle_id_idx
  ON enable_action_bundle (action_bundle_id);

CREATE TABLE IF NOT EXISTS `module_api` (
	api_id               varchar(50)  not null,
	module_id            varchar(50)  not null,
	module_name          varchar(50)  not null,
	feature_id           varchar(50)  not null,
	feature_name         varchar(50)  not null,
	action_bundle_id     varchar(50)  not null,
	action_bundle_name   varchar(50)  not null,
	api_method           varchar(50)  not null,
	api_description      varchar(100) null,
	url_method           varchar(100) not null,
	url                  varchar(255) not null,
	global_admin_action_bundle_visibility BOOL NOT NULL DEFAULT false,
	isv_action_bundle_visibility          BOOL NOT NULL DEFAULT false,
	user_action_bundle_visibility         BOOL NOT NULL DEFAULT false,
	PRIMARY KEY(api_id)
);
CREATE INDEX module_api_module_id_idx
	ON module_api (module_id);

CREATE TABLE IF NOT EXISTS `role` (
	role_id              varchar(50)  not null,
	role_name            varchar(200) not null,
	description          varchar(255) null,
	portal               varchar(50)  not null comment 'global_admin,isv,dev,user',
	owner                varchar(50)  not null,
	owner_path           varchar(50)  not null,
	status               varchar(50)  not null comment 'active,disabled,deleted',
	controller           varchar(50)  not null default 'self' comment 'self,pitrix',
	create_time          timestamp default CURRENT_TIMESTAMP,
	update_time          timestamp default CURRENT_TIMESTAMP,
	status_time          timestamp default CURRENT_TIMESTAMP,
	PRIMARY KEY(role_id)
);
create index role_status_idx
	on role (status);
create index role_portal_idx
	on role (portal);
create index role_owner_idx
	on role (owner);
create index role_owner_path_idx
	on role (owner_path);
create index role_role_name_idx
	on role(role_name);
create index role_create_time_idx
	on role(create_time);

CREATE TABLE IF NOT EXISTS `role_module_binding` (
	bind_id              varchar(50) not null,
	role_id              varchar(50) not null,
	module_id            varchar(50) not null,
	data_level           varchar(50) not null comment 'all,group,self',
	create_time          timestamp default CURRENT_TIMESTAMP,
	is_check_all         bool not null default false,

	PRIMARY KEY(bind_id)
);
create index role_module_binding_role_id_idx
	on role_module_binding(role_id);
create index role_module_binding_module_id_idx
	on role_module_binding(module_id);

create table if not exists `user_role_binding` (
	id                   varchar(50) not null,
	user_id              varchar(50) not null,
	role_id              varchar(50) not null,
	create_time          timestamp default CURRENT_TIMESTAMP,

	PRIMARY KEY(id)
);
create index user_role_binding_user_id_idx
	on user_role_binding(user_id);
create index user_role_binding_role_id_idx
	on user_role_binding(role_id);
