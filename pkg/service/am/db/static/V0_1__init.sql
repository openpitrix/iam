-- Copyright 2019 The OpenPitrix Authors. All rights reserved.
-- Use of this source code is governed by a Apache license
-- that can be found in the LICENSE file.

create table if not exists `enable_action`
(
	enable_id            varchar(50) not null,
	bind_id              varchar(50) not null,
	action_id            varchar(50) not null,

	PRIMARY KEY(enable_id)
);

create table if not exists `module_api`
(
	api_id               varchar(50) not null,
	module_id            varchar(50) not null,
	module_name          varchar(50) not null,
	feature_id           varchar(50) not null,
	feature_name         varchar(50) not null,
	action_id            varchar(50) not null,
	action_name          varchar(50) not null,
	api_method           varchar(50) not null,
	api_description      varchar(100),
	url_method           varchar(100) not null,
	url                  varchar(255) not null,

	PRIMARY KEY(api_id)
);

create index module_api_module_id_idx on module_api
(
	module_id
);

create index module_api_feature_id_idx on module_api
(
	feature_id
);

create index module_api_action_id_idx on module_api
(
	action_id
);

create table if not exists `role`
(
	role_id              varchar(50) not null,
	role_name            varchar(200) not null,
	description          varchar(255),
	portal               varchar(50) not null comment 'admin,isv,dev,user',
	owner                varchar(50) not null,
	owner_path           varchar(50) not null,
	status               varchar(10) not null comment 'active,disabled,deleted',
	create_time          timestamp default CURRENT_TIMESTAMP,
	update_time          timestamp default CURRENT_TIMESTAMP,
	status_time          timestamp default CURRENT_TIMESTAMP,

	PRIMARY KEY(role_id)
);

create index role_status_idx on role
(
	status
);
create index role_portal_idx on role
(
	portal
);

create index role_owner_idx on role
(
	owner
);

create index role_owner_path_idx on role
(
	owner_path
);

create index role_role_name_idx on role
(
	role_name
);

create table if not exists `role_module_binding`
(
	bind_id              varchar(50) not null,
	role_id              varchar(50) not null,
	module_id            varchar(50) not null,
	data_level           varchar(50) not null comment 'all,group,self',
	create_time          timestamp default CURRENT_TIMESTAMP,
	update_time          timestamp default CURRENT_TIMESTAMP,
	is_check_all         bool not null default false,

	PRIMARY KEY(bind_id)
);

create index role_module_binding_data_level_idx on role_module_binding
(
	data_level
);

create index role_module_module_id_idx on role_module_binding
(
	module_id
);

create table if not exists `user_role_binding`
(
	id                   varchar(50) not null,
	user_id              varchar(50) not null,
	role_id              varchar(50) not null,

	PRIMARY KEY(id)
);

INSERT INTO am.module_api (api_id,module_id,module_name,feature_id,feature_name,action_id,action_name,api_method,api_description,url_method,url) VALUES
('api_0001','m_0001','商店管理','f_0001','应用管理','a_0001','查看全部应用','DescribeApps','','get','/v1/apps')
,('api_0002','m_0001','商店管理','f_0001','应用管理','a_0001','查看全部应用','GetAppStatistics','','get','/v1/apps/statistics')
,('api_0003','m_0001','商店管理','f_0001','应用管理','a_0001','查看全部应用','DescribeActiveApps','','get','/v1/active_apps')
,('api_0004','m_0001','商店管理','f_0001','应用管理','a_0001','查看全部应用','DescribeAppVersions','','get','/v1/app_versions')
,('api_0005','m_0001','商店管理','f_0001','应用管理','a_0001','查看全部应用','DescribeActiveAppVersions','','get','/v1/active_app_versions')
,('api_0006','m_0001','商店管理','f_0001','应用管理','a_0001','查看全部应用','DescribeAppVersionAudits','','get','/v1/app_version_audits')
,('api_0007','m_0001','商店管理','f_0001','应用管理','a_0001','查看全部应用','GetAppVersionPackage','','get','/v1/app_version/package')
,('api_0008','m_0001','商店管理','f_0001','应用管理','a_0001','查看全部应用','GetAppVersionPackageFiles','','get','/v1/app_version/package/files')
,('api_0009','m_0001','商店管理','f_0001','应用管理','a_0002','创建应用','CreateApp','','post','/v1/apps')
,('api_0010','m_0001','商店管理','f_0001','应用管理','a_0002','创建应用','CreateAppVersion','','post','/v1/app_versions')
;

INSERT INTO am.module_api (api_id,module_id,module_name,feature_id,feature_name,action_id,action_name,api_method,api_description,url_method,url) VALUES
('api_0011','m_0001','商店管理','f_0001','应用管理','a_0003','修改应用','ModifyApp','','patch','/v1/apps')
,('api_0012','m_0001','商店管理','f_0001','应用管理','a_0003','修改应用','UploadAppAttachment','','patch','/v1/app/attachment')
,('api_0013','m_0001','商店管理','f_0001','应用管理','a_0003','修改应用','ModifyAppVersion','','patch','/v1/app_versions')
,('api_0014','m_0001','商店管理','f_0001','应用管理','a_0004','删除应用','DeleteApps','','delete','/v1/apps')
,('api_0015','m_0001','商店管理','f_0001','应用管理','a_0004','删除应用','DeleteAppVersion','','post','/v1/app_version/action/delete')
,('api_0016','m_0001','商店管理','f_0001','应用管理','a_0005','发布应用','ReleaseAppVersion','','post','/v1/app_version/action/release')
,('api_0017','m_0001','商店管理','f_0001','应用管理','a_0006','下架应用','CancelAppVersion','','post','/v1/app_version/action/cancel')
,('api_0018','m_0001','商店管理','f_0002','应用审核','a_0007','审核提交','SubmitAppVersion','','post','/v1/app_version/action/submit')
,('api_0019','m_0001','商店管理','f_0002','应用审核','a_0008','审核撤销','RecoverAppVersion','','post','/v1/app_version/action/recover')
,('api_0020','m_0001','商店管理','f_0002','应用审核','a_0009','审核通过','PassAppVersion','','post','/v1/app_version/action/{role}/pass')
;

INSERT INTO am.module_api (api_id,module_id,module_name,feature_id,feature_name,action_id,action_name,api_method,api_description,url_method,url) VALUES
('api_0021','m_0001','商店管理','f_0002','应用审核','a_0010','审核拒绝','RejectAppVersion','','post','/v1/app_version/action/{role}/reject')
,('api_0022','m_0001','商店管理','f_0003','应用分类','a_0011','查看全部分类','DescribeCategories','','get','/v1/categories')
,('api_0023','m_0001','商店管理','f_0003','应用分类','a_0012','创建分类','CreateCategory','','post','/v1/categories')
,('api_0024','m_0001','商店管理','f_0003','应用分类','a_0013','修改分类','ModifyCategory','','patch','/v1/categories')
,('api_0025','m_0001','商店管理','f_0003','应用分类','a_0014','删除分类','DeleteCategories','','delete','/v1/categories')
,('api_0026','m_0002','个人中心','f_0004','ssh key 管理','a_0015','创建ssh key','CreateKeyPair','','post','/v1/clusters/key_pairs')
,('api_0027','m_0002','个人中心','f_0004','ssh key 管理','a_0016','查看ssh key','DescribeKeyPairs','','get','/v1/clusters/key_pairs')
,('api_0028','m_0002','个人中心','f_0004','ssh key 管理','a_0017','删除ssh key','DeleteKeyPairs','','delete','/v1/clusters/key_pairs')
,('api_0029','m_0002','个人中心','f_0004','ssh key 管理','a_0018','绑定ssh key','AttachKeyPairs','','post','/v1/clusters/key_pair/attach')
,('api_0030','m_0002','个人中心','f_0004','ssh key 管理','a_0019','解绑ssh key','DetachKeyPairs','','post','/v1/clusters/key_pair/detach')
;

INSERT INTO am.module_api (api_id,module_id,module_name,feature_id,feature_name,action_id,action_name,api_method,api_description,url_method,url) VALUES
('api_0031','m_0003','我的实例','f_0005','应用实例管理','a_0020','创建应用实例','CreateCluster','','post','/v1/clusters/create')
,('api_0032','m_0003','我的实例','f_0005','应用实例管理','a_0021','创建应用实例','DescribeSubnets','','get','/v1/clusters/subnets')
,('api_0033','m_0003','我的实例','f_0005','应用实例管理','a_0022','修改应用实例','ModifyClusterAttributes','','post','/v1/clusters/modify')
,('api_0034','m_0003','我的实例','f_0005','应用实例管理','a_0022','修改应用实例','ModifyClusterNodeAttributes','','post','/v1/clusters/modify_nodes')
,('api_0035','m_0003','我的实例','f_0005','应用实例管理','a_0023','删除应用实例','DeleteClusters','','post','/v1/clusters/delete')
,('api_0036','m_0003','我的实例','f_0005','应用实例管理','a_0024','升级应用实例','UpgradeCluster','','post','/v1/clusters/upgrade')
,('api_0037','m_0003','我的实例','f_0005','应用实例管理','a_0025','回滚应用实例','RollbackCluster','','post','/v1/clusters/rollback')
,('api_0038','m_0003','我的实例','f_0005','应用实例管理','a_0026','纵向伸缩应用实例','ResizeCluster','','post','/v1/clusters/resize')
,('api_0039','m_0003','我的实例','f_0005','应用实例管理','a_0027','横向伸缩应用实例','AddClusterNodes','','post','/v1/clusters/add_nodes')
,('api_0040','m_0003','我的实例','f_0005','应用实例管理','a_0027','横向伸缩应用实例','DeleteClusterNodes','','post','/v1/clusters/delete_nodes')
;

INSERT INTO am.module_api (api_id,module_id,module_name,feature_id,feature_name,action_id,action_name,api_method,api_description,url_method,url) VALUES
('api_0041','m_0003','我的实例','f_0005','应用实例管理','a_0028','更新环境变量','UpdateClusterEnv','','patch','/v1/clusters/update_env')
,('api_0042','m_0003','我的实例','f_0005','应用实例管理','a_0029','查看全部应用实例','DescribeClusters','','get','/v1/clusters')
,('api_0043','m_0003','我的实例','f_0005','应用实例管理','a_0029','查看全部应用实例','DescribeClusterNodes','','get','/v1/clusters/nodes')
,('api_0044','m_0003','我的实例','f_0005','应用实例管理','a_0029','查看全部应用实例','GetClusterStatistics','','get','/v1/clusters/statistics')
,('api_0045','m_0003','我的实例','f_0005','应用实例管理','a_0030','关闭应用实例','StopClusters','','post','/v1/clusters/stop')
,('api_0046','m_0003','我的实例','f_0005','应用实例管理','a_0031','启动应用实例','StartClusters','','post','/v1/clusters/start')
,('api_0047','m_0004','账户与权限','f_0006','用户管理','a_0032','创建用户','CreateUser','','post','/v1/users')
,('api_0048','m_0004','账户与权限','f_0006','用户管理','a_0033','查看全部用户','DescribeUsers','','get','/v1/users')
,('api_0049','m_0004','账户与权限','f_0006','用户管理','a_0034','修改用户','ModifyUser','','patch','/v1/users')
,('api_0050','m_0004','账户与权限','f_0006','用户管理','a_0034','修改用户','ChangePassword','','post','/v1/users/password:change')
;

INSERT INTO am.module_api (api_id,module_id,module_name,feature_id,feature_name,action_id,action_name,api_method,api_description,url_method,url) VALUES
('api_0051','m_0004','账户与权限','f_0006','用户管理','a_0034','修改用户','CreatePasswordReset','','post','/v1/users/password:reset')
,('api_0052','m_0004','账户与权限','f_0006','用户管理','a_0034','修改用户','GetPasswordReset','','get','/v1/users/password/reset/{reset_id}')
,('api_0053','m_0004','账户与权限','f_0006','用户管理','a_0035','删除用户','DeleteUsers','','delete','/v1/users')
,('api_0054','m_0004','账户与权限','f_0007','用户组管理','a_0036','创建用户组','CreateGroup','','post','/v1/groups')
,('api_0055','m_0004','账户与权限','f_0007','用户组管理','a_0037','查看全部用户组','DescribeGroups','','get','/v1/groups')
,('api_0056','m_0004','账户与权限','f_0007','用户组管理','a_0038','修改用户组','ModifyGroup','','patch','/v1/groups')
,('api_0057','m_0004','账户与权限','f_0007','用户组管理','a_0039','删除用户组','DeleteGroups','','delete','/v1/groups')
,('api_0058','m_0004','账户与权限','f_0007','用户组管理','a_0040','加入用户组','JoinGroup','','post','/v1/groups:join')
,('api_0059','m_0004','账户与权限','f_0007','用户组管理','a_0041','踢出用户组','LeaveGroup','','post','/v1/groups:leave')
,('api_0060','m_0005','其它','f_0008','Job 管理','a_0042','查看全部Job','DescribeJobs','','get','/v1/jobs')
;

INSERT INTO am.module_api (api_id,module_id,module_name,feature_id,feature_name,action_id,action_name,api_method,api_description,url_method,url) VALUES
('api_0061','m_0006','其它','f_0009','Task 管理','a_0043','查看全部Task','DescribeTasks','','get','/v1/tasks')
,('api_0062','m_0006','其它','f_0009','Task 管理','a_0044','重试 Task','RetryTasks','','post','/v1/tasks/retry')
,('api_0063','m_0007','平台设置','f_0010','仓库管理','a_0045','创建仓库','CreateRepo','','post','/v1/repos')
,('api_0064','m_0007','平台设置','f_0010','仓库管理','a_0045','创建仓库','ValidateRepo','','get','/v1/repos/validate')
,('api_0065','m_0007','平台设置','f_0010','仓库管理','a_0046','查看全部仓库','DescribeRepos','','get','/v1/repos')
,('api_0066','m_0007','平台设置','f_0010','仓库管理','a_0047','修改仓库','ModifyRepo','','patch','/v1/repos')
,('api_0067','m_0007','平台设置','f_0010','仓库管理','a_0048','删除仓库','DeleteRepos','','delete','/v1/repos')
,('api_0068','m_0007','平台设置','f_0010','仓库管理','a_0049','同步应用','IndexRepo','','post','/v1/repos/index')
,('api_0069','m_0007','平台设置','f_0010','仓库管理','a_0050','查看同步事件','DescribeRepoEvents','','get','/v1/repo_events')
,('api_0070','m_0008','我的环境/个人中心-测试环境','f_0011','环境管理','a_0051','创建环境','CreateRuntime','','post','/v1/runtimes')
;

INSERT INTO am.module_api (api_id,module_id,module_name,feature_id,feature_name,action_id,action_name,api_method,api_description,url_method,url) VALUES
('api_0071','m_0008','我的环境/个人中心-测试环境','f_0011','环境管理','a_0051','创建环境','DescribeRuntimeProviderZones','','get','/v1/runtimes/zones')
,('api_0072','m_0008','我的环境/个人中心-测试环境','f_0011','环境管理','a_0051','创建环境','GetRuntimeStatistics','','get','/v1/runtimes/statistics')
,('api_0073','m_0008','我的环境/个人中心-测试环境','f_0011','环境管理','a_0052','查看全部环境','DescribeRuntimes','','get','/v1/runtimes')
,('api_0074','m_0008','我的环境/个人中心-测试环境','f_0011','环境管理','a_0053','修改环境','ModifyRuntime','','patch','/v1/runtimes')
,('api_0075','m_0008','我的环境/个人中心-测试环境','f_0011','环境管理','a_0054','删除环境','DeleteRuntimes','','delete','/v1/runtimes')
,('api_0076','m_0008','我的环境/个人中心-测试环境','f_0012','授权信息管理','a_0055','创建授权信息','CreateRuntimeCredential','','post','/v1/runtimes/credentials')
,('api_0077','m_0008','我的环境/个人中心-测试环境','f_0012','授权信息管理','a_0056','查看全部授权信息','DescribeRuntimeCredentials','','get','/v1/runtimes/credentials')
,('api_0078','m_0008','我的环境/个人中心-测试环境','f_0012','授权信息管理','a_0057','修改授权信息','ModifyRuntimeCredential','','patch','/v1/runtimes/credentials')
,('api_0079','m_0008','我的环境/个人中心-测试环境','f_0012','授权信息管理','a_0058','删除授权信息','DeleteRuntimeCredentials','','delete','/v1/runtimes/credentials')
,('api_0080','m_0009','申请成为服务商','f_0013','申请成为服务商','a_0059','提交服务商认证','SubmitVendorVerifyInfo','','post','/v1/app_vendors')
;
INSERT INTO am.module_api (api_id,module_id,module_name,feature_id,feature_name,action_id,action_name,api_method,api_description,url_method,url) VALUES
('api_0081','m_0009','申请成为服务商','f_0013','申请成为服务商','a_0059','提交服务商认证','UploadVendorVerifyAttachment','','patch','/v1/app_vendors/attachment')
,('api_0082','m_0010','应用服务商管理','f_0014','入驻申请','a_0060','通过服务商认证','PassVendorVerifyInfo','','post','/v1/app_vendors/pass')
,('api_0083','m_0010','应用服务商管理','f_0014','入驻申请','a_0061','拒绝服务商认证','RejectVendorVerifyInfo','','post','/v1/app_vendors/reject')
,('api_0084','m_0010','应用服务商管理','f_0014','入驻申请','a_0062','查看全部服务商认证','DescribeVendorVerifyInfos','','get','/v1/app_vendors')
,('api_0085','m_0010','应用服务商管理','f_0014','入驻申请','a_0062','查看全部服务商认证','DescribeAppVendorStatistics','','get','/v1/app_vendors/app_vendor_statistics')
;

INSERT INTO `role` (role_id,role_name,description,portal,owner,owner_path,status) VALUES
('developer','developer','','dev','system',':system','active')
,('global_admin','global_admin','','admin','system',':system','active')
,('isv','isv','','isv','system',':system','active')
,('user','user','','user','system',':system','active')
;

INSERT INTO role_module_binding (bind_id,role_id,module_id,data_level,is_check_all) VALUES
('bind_0001','global_admin','m_0001','all',1)
,('bind_0002','global_admin','m_0002','all',1)
,('bind_0003','global_admin','m_0003','all',1)
,('bind_0004','global_admin','m_0004','all',1)
,('bind_0005','global_admin','m_0005','all',1)
,('bind_0006','global_admin','m_0006','all',1)
,('bind_0007','global_admin','m_0007','all',1)
,('bind_0008','global_admin','m_0008','all',1)
,('bind_0009','global_admin','m_0010','all',1)
,('bind_0010','isv','m_0001','group',1)
;

INSERT INTO role_module_binding (bind_id,role_id,module_id,data_level,is_check_all) VALUES
('bind_0011','isv','m_0002','group',1)
,('bind_0012','isv','m_0003','group',1)
,('bind_0013','isv','m_0004','group',1)
,('bind_0014','isv','m_0008','group',1)
,('bind_0015','isv','m_0009','group',1)
,('bind_0016','isv','m_0010','group',1)
,('bind_0017','developer','m_0001','self',1)
,('bind_0018','developer','m_0002','self',1)
,('bind_0019','developer','m_0003','self',1)
,('bind_0020','developer','m_0008','self',1)
;

INSERT INTO role_module_binding (bind_id,role_id,module_id,data_level,is_check_all) VALUES
('bind_0021','user','m_0001','self',1)
,('bind_0022','user','m_0002','self',1)
,('bind_0023','user','m_0003','self',1)
,('bind_0024','user','m_0004','self',1)
,('bind_0025','user','m_0008','self',1)
,('bind_0026','user','m_0009','self',1)
;
