DELETE FROM module_api;

INSERT INTO module_api
(module_id,module_name,feature_id,feature_name,action_bundle_id,action_bundle_name,global_admin_action_bundle_visibility,isv_action_bundle_visibility,user_action_bundle_visibility,api_id,api_method,url_method,url) VALUES
('m0','默认权限','m0.f1','用户管理','m0.f1.a3','查看用户',0,0,0,'m0.f1.a3.api1','DescribeUsers','get','/v1/users'),
('m0','默认权限','m0.f1','用户管理','m0.f1.a3','查看用户',0,0,0,'m0.f1.a3.api2','DescribeUsersDetail','get','/v1/users_detail'),
('m0','默认权限','m0.f1','用户管理','m0.f1.a4','修改用户',0,0,0,'m0.f1.a4.api1','ModifyUser','patch','/v1/users'),
('m0','默认权限','m0.f1','用户管理','m0.f1.a4','修改用户',0,0,0,'m0.f1.a4.api2','ChangePassword','post','/v1/users/password:change'),
('m0','默认权限','m0.f1','用户管理','m0.f1.a4','修改用户',0,0,0,'m0.f1.a4.api3','CreatePasswordReset','post','/v1/users/password:reset'),
('m0','默认权限','m0.f1','用户管理','m0.f1.a4','修改用户',0,0,0,'m0.f1.a4.api4','ValidateUserPassword','post','/v1/users/password:validate'),
('m0','默认权限','m0.f1','用户管理','m0.f1.a4','修改用户',0,0,0,'m0.f1.a4.api5','GetPasswordReset','get','/v1/users/password:reset'),
('m0','默认权限','m0.f2','角色管理','m0.f2.a1','查看角色',0,0,0,'m0.f2.a1.api1','DescribeRoles','get','/v1/roles'),
('m0','默认权限','m0.f2','角色管理','m0.f2.a1','查看角色',0,0,0,'m0.f2.a1.api2','GetRole','get','/v1/role'),
('m0','默认权限','m0.f2','角色管理','m0.f2.a1','查看角色',0,0,0,'m0.f2.a1.api3','GetRoleModule','get','/v1/roles:module'),
('m0','默认权限','m0.f3','用户组管理','m0.f3.a1','查看用户组',0,0,0,'m0.f3.a1.api1','DescribeGroups','get','/v1/groups'),
('m0','默认权限','m0.f3','用户组管理','m0.f3.a1','查看用户组',0,0,0,'m0.f3.a1.api2','DescribeGroupsDetail','get','/v1/groups_detail'),
('m0','默认权限','m0.f4','附件管理','m0.f4.a1','获取附件',0,0,0,'m0.f4.a1.api1','GetAttachment','get','/v1/attachments'),
('m0','默认权限','m0.f5','服务设置','m0.f5.a1','查看服务设置',0,0,0,'m0.f5.a1.api1','GetServiceConfig','post','/v1/service_configs/get'),
('m0','默认权限','m0.f6','应用分类','m0.f6.a1','查看分类',0,0,0,'m0.f6.a1.api1','DescribeCategories','get','/v1/categories'),
('m0','默认权限','m0.f7','Job管理','m0.f7.a1','查看Job',0,0,0,'m0.f7.a1.api1','DescribeJobs','get','/v1/jobs'),
('m0','默认权限','m0.f8','Task管理','m0.f8.a1','查看Task',0,0,0,'m0.f8.a1.api1','DescribeTasks','get','/v1/tasks'),
('m0','默认权限','m0.f9','应用管理','m0.f9.a1','查看应用',0,0,0,'m0.f9.a1.api1','DescribeApps','get','/v1/apps'),
('m0','默认权限','m0.f9','应用管理','m0.f9.a1','查看应用',0,0,0,'m0.f9.a1.api3','DescribeActiveApps','get','/v1/active_apps'),
('m0','默认权限','m0.f9','应用管理','m0.f9.a1','查看应用',0,0,0,'m0.f9.a1.api4','DescribeAppVersions','get','/v1/app_versions'),
('m0','默认权限','m0.f9','应用管理','m0.f9.a1','查看应用',0,0,0,'m0.f9.a1.api5','DescribeActiveAppVersions','get','/v1/active_app_versions'),
('m0','默认权限','m0.f9','应用管理','m0.f9.a1','查看应用',0,0,0,'m0.f9.a1.api6','DescribeAppVersionAudits','get','/v1/app_version_audits'),
('m0','默认权限','m0.f9','应用管理','m0.f9.a1','查看应用',0,0,0,'m0.f9.a1.api7','GetAppVersionPackage','get','/v1/app_version/package'),
('m0','默认权限','m0.f9','应用管理','m0.f9.a1','查看应用',0,0,0,'m0.f9.a1.api8','GetAppVersionPackageFiles','get','/v1/app_version/package/files'),
('m0','默认权限','m0.f10','应用实例管理','m0.f10.a1','查看应用实例',0,0,0,'m0.f10.a1.api1','DescribeClusters','get','/v1/clusters'),
('m0','默认权限','m0.f10','应用实例管理','m0.f10.a1','查看应用实例',0,0,0,'m0.f10.a1.api2','DescribeClusterNodes','get','/v1/clusters/nodes'),
('m0','默认权限','m0.f10','应用实例管理','m0.f10.a1','查看应用实例',0,0,0,'m0.f10.a1.api3','DescribeAppClusters','get','/v1/clusters/apps'),
('m0','默认权限','m0.f10','应用实例管理','m0.f10.a1','查看应用实例',0,0,0,'m0.f10.a1.api4','DescribeDebugClusters','get','/v1/debug_clusters'),
('m0','默认权限','m0.f11','环境管理','m0.f11.a1','查看环境',0,0,0,'m0.f11.a1.api1','DescribeRuntimes','get','/v1/runtimes'),
('m0','默认权限','m0.f11','环境管理','m0.f11.a1','查看环境',0,0,0,'m0.f11.a1.api2','DescribeDebugRuntimes','get','/v1/debug_runtimes'),
('m0','默认权限','m0.f12','授权信息管理','m0.f12.a1','查看授权信息',0,0,0,'m0.f12.a1.api1','DescribeRuntimeCredentials','get','/v1/runtimes/credentials'),
('m0','默认权限','m0.f12','授权信息管理','m0.f12.a1','查看授权信息',0,0,0,'m0.f12.a1.api2','DescribeDebugRuntimeCredentials','get','/v1/debug_runtimes/credentials'),
('m1','商店管理','m1.f1','应用管理','m1.f1.a1','查看应用统计信息',1,1,1,'m1.f1.a1.api2','GetAppStatistics','get','/v1/apps/statistics'),
('m1','商店管理','m1.f1','应用管理','m1.f1.a2','创建应用',1,1,0,'m1.f1.a2.api1','CreateApp','post','/v1/apps'),
('m1','商店管理','m1.f1','应用管理','m1.f1.a2','创建应用',1,1,0,'m1.f1.a2.api2','CreateAppVersion','post','/v1/app_versions'),
('m1','商店管理','m1.f1','应用管理','m1.f1.a2','创建应用',1,1,0,'m1.f1.a2.api3','ValidatePackage','post','/v1/apps/validate_package'),
('m1','商店管理','m1.f1','应用管理','m1.f1.a2','创建应用',1,1,0,'m1.f1.a2.api4','UploadAppAttachment','patch','/v1/app/attachment'),
('m1','商店管理','m1.f1','应用管理','m1.f1.a3','修改应用',1,1,0,'m1.f1.a3.api1','ModifyApp','patch','/v1/apps'),
('m1','商店管理','m1.f1','应用管理','m1.f1.a3','修改应用',1,1,0,'m1.f1.a3.api3','ModifyAppVersion','patch','/v1/app_versions'),
('m1','商店管理','m1.f1','应用管理','m1.f1.a4','删除应用',1,1,0,'m1.f1.a4.api1','DeleteApps','delete','/v1/apps'),
('m1','商店管理','m1.f1','应用管理','m1.f1.a4','删除应用',1,1,0,'m1.f1.a4.api2','DeleteAppVersion','post','/v1/app_version/action/delete'),
('m1','商店管理','m1.f1','应用管理','m1.f1.a5','发布应用',1,1,0,'m1.f1.a5.api1','ReleaseAppVersion','post','/v1/app_version/action/release'),
('m1','商店管理','m1.f1','应用管理','m1.f1.a6','下架应用',1,0,0,'m1.f1.a6.api1','SuspendAppVersion','post','/v1/app_version/action/suspend'),
('m1','商店管理','m1.f1','应用管理','m1.f1.a7','恢复应用',1,0,0,'m1.f1.a7.api1','RecoverAppVersion','post','/v1/app_version/action/recover'),
('m1','商店管理','m1.f2','应用审核','m1.f2.a1','提交审核',1,1,0,'m1.f2.a1.api1','SubmitAppVersion','post','/v1/app_version/action/submit'),
('m1','商店管理','m1.f2','应用审核','m1.f2.a2','撤销审核',1,1,0,'m1.f2.a2.api1','CancelAppVersion','post','/v1/app_version/action/cancel'),
('m1','商店管理','m1.f2','应用审核','m1.f2.a3','商务审核',1,0,0,'m1.f2.a3.api1','BusinessAdminPassAppVersion','post','/v1/app_version/action/business_admin/pass'),
('m1','商店管理','m1.f2','应用审核','m1.f2.a3','商务审核',1,0,0,'m1.f2.a3.api2','BusinessAdminRejectAppVersion','post','/v1/app_version/action/business_admin/reject'),
('m1','商店管理','m1.f2','应用审核','m1.f2.a3','商务审核',1,0,0,'m1.f2.a3.api3','BusinessAdminReviewAppVersion','post','/v1/app_version/action/business_admin/review'),
('m1','商店管理','m1.f2','应用审核','m1.f2.a3','商务审核',1,0,0,'m1.f2.a3.api4','DescribeAppVersionReviews','get','/v1/app_version_reviews'),
('m1','商店管理','m1.f2','应用审核','m1.f2.a4','技术审核',1,0,0,'m1.f2.a4.api1','DevelopAdminPassAppVersion','post','/v1/app_version/action/develop_admin/pass'),
('m1','商店管理','m1.f2','应用审核','m1.f2.a4','技术审核',1,0,0,'m1.f2.a4.api2','DevelopAdminRejectAppVersion','post','/v1/app_version/action/develop_admin/reject'),
('m1','商店管理','m1.f2','应用审核','m1.f2.a4','技术审核',1,0,0,'m1.f2.a4.api3','DevelopAdminReviewAppVersion','post','/v1/app_version/action/develop_admin/review'),
('m1','商店管理','m1.f2','应用审核','m1.f2.a4','技术审核',1,0,0,'m1.f2.a4.api4','DescribeAppVersionReviews','get','/v1/app_version_reviews'),
('m1','商店管理','m1.f2','应用审核','m1.f2.a5','ISV审核',0,1,0,'m1.f2.a5.api1','IsvRejectAppVersion','post','/v1/app_version/action/isv/reject'),
('m1','商店管理','m1.f2','应用审核','m1.f2.a5','ISV审核',0,1,0,'m1.f2.a5.api2','IsvReviewAppVersion','post','/v1/app_version/action/isv/review'),
('m1','商店管理','m1.f2','应用审核','m1.f2.a5','ISV审核',0,1,0,'m1.f2.a5.api3','IsvPassAppVersion','post','/v1/app_version/action/isv/pass'),
('m1','商店管理','m1.f2','应用审核','m1.f2.a5','ISV审核',0,1,0,'m1.f2.a5.api4','DescribeAppVersionReviews','get','/v1/app_version_reviews'),
('m1','商店管理','m1.f3','应用分类','m1.f3.a2','创建分类',1,0,0,'m1.f3.a2.api1','CreateCategory','post','/v1/categories'),
('m1','商店管理','m1.f3','应用分类','m1.f3.a3','修改分类',1,0,0,'m1.f3.a3.api1','ModifyCategory','patch','/v1/categories'),
('m1','商店管理','m1.f3','应用分类','m1.f3.a4','删除分类',1,0,0,'m1.f3.a4.api1','DeleteCategories','delete','/v1/categories'),
('m2','个人中心','m2.f1','ssh key管理','m2.f1.a1','创建ssh key',1,1,1,'m2.f1.a1.api1','CreateKeyPair','post','/v1/clusters/key_pairs'),
('m2','个人中心','m2.f1','ssh key管理','m2.f1.a2','查看ssh key',1,1,1,'m2.f1.a2.api1','DescribeKeyPairs','get','/v1/clusters/key_pairs'),
('m2','个人中心','m2.f1','ssh key管理','m2.f1.a3','删除ssh key',1,1,1,'m2.f1.a3.api1','DeleteKeyPairs','delete','/v1/clusters/key_pairs'),
('m2','个人中心','m2.f1','ssh key管理','m2.f1.a4','绑定ssh key',1,1,1,'m2.f1.a4.api1','AttachKeyPairs','post','/v1/clusters/key_pair/attach'),
('m2','个人中心','m2.f1','ssh key管理','m2.f1.a5','解绑ssh key',1,1,1,'m2.f1.a5.api1','DetachKeyPairs','post','/v1/clusters/key_pair/detach'),
('m3','我的实例','m3.f1','应用实例管理','m3.f1.a1','创建应用实例',1,1,1,'m3.f1.a1.api1','CreateCluster','post','/v1/clusters/create'),
('m3','我的实例','m3.f1','应用实例管理','m3.f1.a1','创建应用实例',1,1,0,'m3.f1.a1.api2','CreateDebugCluster','post','/v1/debug_clusters/create'),
('m3','我的实例','m3.f1','应用实例管理','m3.f1.a1','创建应用实例',1,1,1,'m3.f1.a1.api3','DescribeSubnets','get','/v1/clusters/subnets'),
('m3','我的实例','m3.f1','应用实例管理','m3.f1.a2','修改应用实例',1,1,1,'m3.f1.a2.api1','ModifyClusterAttributes','post','/v1/clusters/modify'),
('m3','我的实例','m3.f1','应用实例管理','m3.f1.a2','修改应用实例',1,1,1,'m3.f1.a2.api2','ModifyClusterNodeAttributes','post','/v1/clusters/modify_nodes'),
('m3','我的实例','m3.f1','应用实例管理','m3.f1.a3','删除应用实例',1,1,1,'m3.f1.a3.api1','DeleteClusters','post','/v1/clusters/delete'),
('m3','我的实例','m3.f1','应用实例管理','m3.f1.a4','升级应用实例',1,1,1,'m3.f1.a4.api1','UpgradeCluster','post','/v1/clusters/upgrade'),
('m3','我的实例','m3.f1','应用实例管理','m3.f1.a5','回滚应用实例',1,1,1,'m3.f1.a5.api1','RollbackCluster','post','/v1/clusters/rollback'),
('m3','我的实例','m3.f1','应用实例管理','m3.f1.a6','纵向伸缩应用实例',1,1,1,'m3.f1.a6.api1','ResizeCluster','post','/v1/clusters/resize'),
('m3','我的实例','m3.f1','应用实例管理','m3.f1.a7','横向伸缩应用实例',1,1,1,'m3.f1.a7.api1','AddClusterNodes','post','/v1/clusters/add_nodes'),
('m3','我的实例','m3.f1','应用实例管理','m3.f1.a7','横向伸缩应用实例',1,1,1,'m3.f1.a7.api2','DeleteClusterNodes','post','/v1/clusters/delete_nodes'),
('m3','我的实例','m3.f1','应用实例管理','m3.f1.a8','更新环境变量',1,1,1,'m3.f1.a8.api1','UpdateClusterEnv','patch','/v1/clusters/update_env'),
('m3','我的实例','m3.f1','应用实例管理','m3.f1.a9','查看应用实例统计信息',1,1,1,'m3.f1.a9.api4','GetClusterStatistics','get','/v1/clusters/statistics'),
('m3','我的实例','m3.f1','应用实例管理','m3.f1.a10','关闭应用实例',1,1,1,'m3.f1.a10.api1','StopClusters','post','/v1/clusters/stop'),
('m3','我的实例','m3.f1','应用实例管理','m3.f1.a11','启动应用实例',1,1,1,'m3.f1.a11.api1','StartClusters','post','/v1/clusters/start'),
('m3','我的实例','m3.f1','应用实例管理','m3.f1.a12','销毁应用实例',1,1,1,'m3.f1.a12.api1','CeaseClusters','post','/v1/clusters/cease'),
('m3','我的实例','m3.f1','应用实例管理','m3.f1.a13','恢复应用实例',1,1,1,'m3.f1.a13.api1','RecoverClusters','post','/v1/clusters/recover'),
('m4','账户与权限','m4.f1','用户管理','m4.f1.a1','管理员创建用户',1,0,0,'m4.f1.a1.api1','CreateUser','post','/v1/users'),
('m4','账户与权限','m4.f1','用户管理','m4.f1.a2','ISV创建用户',0,1,0,'m4.f1.a2.api1','IsvCreateUser','post','/v1/isv_users'),
('m4','账户与权限','m4.f1','用户管理','m4.f1.a5','修改用户角色',1,1,1,'m4.f1.a5.api1','BindUserRole','post','/v1/user:role'),
('m4','账户与权限','m4.f1','用户管理','m4.f1.a5','修改用户角色',1,1,1,'m4.f1.a5.api2','UnbindUserRole','delete','/v1/user:role'),
('m4','账户与权限','m4.f1','用户管理','m4.f1.a6','删除用户',1,1,0,'m4.f1.a6.api1','DeleteUsers','delete','/v1/users'),
('m4','账户与权限','m4.f2','角色管理','m4.f2.a2','删除角色',1,1,0,'m4.f2.a2.api1','DeleteRoles','delete','/v1/roles'),
('m4','账户与权限','m4.f2','角色管理','m4.f2.a3','创建角色',1,1,0,'m4.f2.a3.api1','CreateRole','post','/v1/roles'),
('m4','账户与权限','m4.f2','角色管理','m4.f2.a4','修改角色',1,1,0,'m4.f2.a4.api1','ModifyRole','patch','/v1/roles'),
('m4','账户与权限','m4.f2','角色管理','m4.f2.a5','修改角色权限',1,1,0,'m4.f2.a5.api1','ModifyRoleModule','post','/v1/roles:module'),
('m4','账户与权限','m4.f3','用户组管理','m4.f3.a1','创建用户组',1,1,1,'m4.f3.a1.api1','CreateGroup','post','/v1/groups'),
('m4','账户与权限','m4.f3','用户组管理','m4.f3.a3','修改用户组',1,1,1,'m4.f3.a3.api1','ModifyGroup','patch','/v1/groups'),
('m4','账户与权限','m4.f3','用户组管理','m4.f3.a4','删除用户组',1,1,1,'m4.f3.a4.api1','DeleteGroups','delete','/v1/groups'),
('m4','账户与权限','m4.f3','用户组管理','m4.f3.a5','加入用户组',1,1,1,'m4.f3.a5.api1','JoinGroup','post','/v1/groups:join'),
('m4','账户与权限','m4.f3','用户组管理','m4.f3.a6','离开用户组',1,1,1,'m4.f3.a6.api1','LeaveGroup','post','/v1/groups:leave'),
('m5','平台设置','m5.f1','仓库管理','m5.f1.a1','创建仓库',1,0,0,'m5.f1.a1.api1','CreateRepo','post','/v1/repos'),
('m5','平台设置','m5.f1','仓库管理','m5.f1.a1','创建仓库',1,0,0,'m5.f1.a1.api2','ValidateRepo','get','/v1/repos/validate'),
('m5','平台设置','m5.f1','仓库管理','m5.f1.a2','查看仓库',1,1,1,'m5.f1.a2.api1','DescribeRepos','get','/v1/repos'),
('m5','平台设置','m5.f1','仓库管理','m5.f1.a3','修改仓库',1,0,0,'m5.f1.a3.api1','ModifyRepo','patch','/v1/repos'),
('m5','平台设置','m5.f1','仓库管理','m5.f1.a4','删除仓库',1,0,0,'m5.f1.a4.api1','DeleteRepos','delete','/v1/repos'),
('m5','平台设置','m5.f1','仓库管理','m5.f1.a5','同步应用',1,0,0,'m5.f1.a5.api1','IndexRepo','post','/v1/repos/index'),
('m5','平台设置','m5.f1','仓库管理','m5.f1.a6','查看同步事件',1,0,0,'m5.f1.a6.api1','DescribeRepoEvents','get','/v1/repo_events'),
('m5','平台设置','m5.f2','服务设置','m5.f2.a2','修改服务设置',1,0,0,'m5.f2.a2.api1','SetServiceConfig','post','/v1/service_configs/set'),
('m6','我的环境/个人中心-测试环境','m6.f1','环境管理','m6.f1.a1','创建环境',1,1,1,'m6.f1.a1.api1','CreateRuntime','post','/v1/runtimes'),
('m6','我的环境/个人中心-测试环境','m6.f1','环境管理','m6.f1.a1','创建环境',1,1,1,'m6.f1.a1.api2','DescribeRuntimeProviderZones','get','/v1/runtimes/zones'),
('m6','我的环境/个人中心-测试环境','m6.f1','环境管理','m6.f1.a1','创建环境',1,1,0,'m6.f1.a1.api3','CreateDebugRuntime','post','/v1/debug_runtimes'),
('m6','我的环境/个人中心-测试环境','m6.f1','环境管理','m6.f1.a2','查看环境统计信息',1,1,1,'m6.f1.a2.api1','GetRuntimeStatistics','get','/v1/runtimes/statistics'),
('m6','我的环境/个人中心-测试环境','m6.f1','环境管理','m6.f1.a3','修改环境',1,1,1,'m6.f1.a3.api1','ModifyRuntime','patch','/v1/runtimes'),
('m6','我的环境/个人中心-测试环境','m6.f1','环境管理','m6.f1.a4','删除环境',1,1,1,'m6.f1.a4.api1','DeleteRuntimes','delete','/v1/runtimes'),
('m6','我的环境/个人中心-测试环境','m6.f2','授权信息管理','m6.f2.a1','创建授权信息',1,1,1,'m6.f2.a1.api1','CreateRuntimeCredential','post','/v1/runtimes/credentials'),
('m6','我的环境/个人中心-测试环境','m6.f2','授权信息管理','m6.f2.a1','创建授权信息',1,1,1,'m6.f2.a1.api2','ValidateRuntimeCredential','post','/v1/runtimes/credentials:validate'),
('m6','我的环境/个人中心-测试环境','m6.f2','授权信息管理','m6.f2.a1','创建授权信息',1,1,0,'m6.f2.a1.api3','CreateDebugRuntimeCredential','post','/v1/debug_runtimes/credentials'),
('m6','我的环境/个人中心-测试环境','m6.f2','授权信息管理','m6.f2.a3','修改授权信息',1,1,1,'m6.f2.a2.api1','ModifyRuntimeCredential','patch','/v1/runtimes/credentials'),
('m6','我的环境/个人中心-测试环境','m6.f2','授权信息管理','m6.f2.a4','删除授权信息',1,1,1,'m6.f2.a3.api1','DeleteRuntimeCredentials','delete','/v1/runtimes/credentials'),
('m7','应用服务商管理','m7.f1','申请成为服务商','m7.f1.a1','提交服务商认证',0,1,1,'m7.f1.a1.api1','SubmitVendorVerifyInfo','post','/v1/app_vendors'),
('m7','应用服务商管理','m7.f2','入驻申请','m7.f1.a2','通过服务商认证',1,0,0,'m7.f1.a2.api1','PassVendorVerifyInfo','post','/v1/app_vendors/pass'),
('m7','应用服务商管理','m7.f2','入驻申请','m7.f1.a3','拒绝服务商认证',1,0,0,'m7.f1.a3.api1','RejectVendorVerifyInfo','post','/v1/app_vendors/reject'),
('m7','应用服务商管理','m7.f2','入驻申请','m7.f1.a4','查看服务商认证',1,1,1,'m7.f1.a4.api1','DescribeVendorVerifyInfos','get','/v1/app_vendors'),
('m7','应用服务商管理','m7.f2','入驻申请','m7.f1.a4','查看服务商认证',1,1,1,'m7.f1.a4.api2','DescribeAppVendorStatistics','get','/v1/app_vendors/app_vendor_statistics'),
('m8','其它','m8.f2','Task管理','m8.f2.a2','重试Task',1,1,0,'m8.f2.a2.api1','RetryTasks','post','/v1/tasks/retry')
;
