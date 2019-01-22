-- -------------------------------------------------------------
-- TablePlus 1.0(166)
--
-- https://tableplus.com/
--
-- Database: am2
-- Generation Time: 2019-01-22 23:24:56.2930
-- -------------------------------------------------------------


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;


CREATE TABLE IF NOT EXISTS `user_role_binding` (
  `id` varchar(50) NOT NULL,
  `user_id` varchar(50) DEFAULT NULL,
  `role_id` varchar(50) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS `role_module_binding` (
  `bind_id` varchar(50) NOT NULL,
  `role_id` varchar(50) DEFAULT NULL,
  `module_id` varchar(50) DEFAULT NULL,
  `data_level` varchar(50) DEFAULT NULL COMMENT 'all,department,onlyself',
  `create_time` timestamp NULL DEFAULT NULL,
  `update_time` timestamp NULL DEFAULT NULL,
  `owner` varchar(50) DEFAULT NULL,
  `is_feature_all_checked` tinyint(1) NOT NULL DEFAULT '0',
  PRIMARY KEY (`bind_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS `role` (
  `role_id` varchar(50) NOT NULL,
  `role_name` varchar(200) DEFAULT NULL,
  `description` varchar(255) DEFAULT NULL,
  `portal` varchar(50) DEFAULT NULL COMMENT ' admin,isv,dev,normal',
  `create_time` timestamp NULL DEFAULT NULL,
  `update_time` timestamp NULL DEFAULT NULL,
  `owner` varchar(50) DEFAULT NULL,
  `owner_path` varchar(50) DEFAULT NULL,
  PRIMARY KEY (`role_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS `module_api` (
  `api_id` varchar(50) NOT NULL,
  `module_id` varchar(50) DEFAULT NULL,
  `module_name` varchar(50) DEFAULT NULL,
  `feature_id` varchar(50) DEFAULT NULL,
  `feature_name` varchar(50) DEFAULT NULL,
  `action_id` varchar(50) DEFAULT NULL,
  `action_name` varchar(50) DEFAULT NULL,
  `api_method` varchar(50) DEFAULT NULL,
  `api_description` varchar(100) DEFAULT NULL,
  `url_method` varchar(20) DEFAULT NULL,
  `url` varchar(500) DEFAULT NULL,
  PRIMARY KEY (`api_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS `enable_action` (
  `enable_id` varchar(50) NOT NULL,
  `bind_id` varchar(50) DEFAULT NULL,
  `action_id` varchar(50) DEFAULT NULL,
  PRIMARY KEY (`enable_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

INSERT INTO `user_role_binding` (`id`, `user_id`, `role_id`) VALUES ('urbind0001', 'uid-PYu7bdqa', 'role_0001'),
('xid-MP4dkvPG', 'chai', 'role_0002');

INSERT INTO `role_module_binding` (`bind_id`, `role_id`, `module_id`, `data_level`, `create_time`, `update_time`, `owner`, `is_feature_all_checked`) VALUES ('bind_00001', 'role_0001', 'm_0001', 'all', NULL, NULL, 'system', '1'),
('bind_00002', 'role_0001', 'm_0002', 'department', NULL, NULL, 'system', '1'),
('bind_00003', 'role_0001', 'm_0003', 'all', NULL, NULL, 'system', '1'),
('bind_00004', 'role_0001', 'm_0004', 'all', NULL, NULL, 'system', '1'),
('bind_00005', 'role_0001', 'm_0005', 'all', NULL, NULL, 'system', '1'),
('bind_00006', 'role_0001', 'm_0006', 'all', NULL, NULL, 'system', '1'),
('bind_00007', 'role_0001', 'm_0007', 'all', NULL, NULL, 'system', '1'),
('bind_00008', 'role_0001', 'm_0008', 'all', NULL, NULL, 'system', '1'),
('bind_00009', 'role_0001', 'm_00010', 'all', NULL, NULL, 'system', '1'),
('bind_00010', 'role_0002', 'm_0001', 'all', NULL, NULL, 'system', '1'),
('bind_00011', 'role_0002', 'm_0002', 'department', NULL, NULL, 'system', '1'),
('bind_00012', 'role_0002', 'm_0003', 'all', NULL, NULL, 'system', '1'),
('bind_00013', 'role_0002', 'm_0004', 'all', NULL, NULL, 'system', '1'),
('bind_00014', 'role_0002', 'm_0005', 'all', NULL, NULL, 'system', '1'),
('bind_00015', 'role_0002', 'm_0008', 'all', NULL, NULL, 'system', '1'),
('bind_00016', 'role_0002', 'm_0009', 'all', NULL, NULL, 'system', '1'),
('bind_00017', 'role_0002', 'm_0010', 'all', NULL, NULL, 'system', '1');

INSERT INTO `role` (`role_id`, `role_name`, `description`, `portal`, `create_time`, `update_time`, `owner`, `owner_path`) VALUES ('role_0001', '超级管理员', 'Portal是admin的超级管理员', 'admin', NULL, NULL, 'system', 'system.'),
('role_0002', '超级管理员', 'Portal是isv的超级管理员', 'isv', NULL, NULL, 'system', 'system.'),
('role_0003', '超级管理员', 'Portal是dev的超级管理员', 'dev', NULL, NULL, 'system', 'system.'),
('role_0004', '超级管理员', 'Portal是normal的超级管理员', 'normal', NULL, NULL, 'system', 'system.');

INSERT INTO `module_api` (`api_id`, `module_id`, `module_name`, `feature_id`, `feature_name`, `action_id`, `action_name`, `api_method`, `api_description`, `url_method`, `url`) VALUES ('api_0001', 'm_0001', '商店管理', 'f_0001', '应用管理', 'a_0001', '查看全部应用', 'DescribeApps', '', 'get', '/v1/apps'),
('api_0002', 'm_0001', '商店管理', 'f_0001', '应用管理', 'a_0001', '查看全部应用', 'GetAppStatistics', '', 'get', '/v1/apps/statistics'),
('api_0003', 'm_0001', '商店管理', 'f_0001', '应用管理', 'a_0001', '查看全部应用', 'DescribeActiveApps', '', 'get', '/v1/active_apps'),
('api_0004', 'm_0001', '商店管理', 'f_0001', '应用管理', 'a_0001', '查看全部应用', 'DescribeAppVersions', '', 'get', '/v1/app_versions'),
('api_0005', 'm_0001', '商店管理', 'f_0001', '应用管理', 'a_0001', '查看全部应用', 'DescribeActiveAppVersions', '', 'get', '/v1/active_app_versions'),
('api_0006', 'm_0001', '商店管理', 'f_0001', '应用管理', 'a_0001', '查看全部应用', 'DescribeAppVersionAudits', '', 'get', '/v1/app_version_audits'),
('api_0007', 'm_0001', '商店管理', 'f_0001', '应用管理', 'a_0001', '查看全部应用', 'GetAppVersionPackage', '', 'get', '/v1/app_version/package'),
('api_0008', 'm_0001', '商店管理', 'f_0001', '应用管理', 'a_0001', '查看全部应用', 'GetAppVersionPackageFiles', '', 'get', '/v1/app_version/package/files'),
('api_0009', 'm_0001', '商店管理', 'f_0001', '应用管理', 'a_0002', '创建应用', 'CreateApp', '', 'post', '/v1/apps'),
('api_0010', 'm_0001', '商店管理', 'f_0001', '应用管理', 'a_0002', '创建应用', 'CreateAppVersion', '', 'post', '/v1/app_versions'),
('api_0011', 'm_0001', '商店管理', 'f_0001', '应用管理', 'a_0003', '修改应用', 'ModifyApp', '', 'patch', '/v1/apps'),
('api_0012', 'm_0001', '商店管理', 'f_0001', '应用管理', 'a_0003', '修改应用', 'UploadAppAttachment', '', 'patch', '/v1/app/attachment'),
('api_0013', 'm_0001', '商店管理', 'f_0001', '应用管理', 'a_0003', '修改应用', 'ModifyAppVersion', '', 'patch', '/v1/app_versions'),
('api_0014', 'm_0001', '商店管理', 'f_0001', '应用管理', 'a_0004', '删除应用', 'DeleteApps', '', 'delete', '/v1/apps'),
('api_0015', 'm_0001', '商店管理', 'f_0001', '应用管理', 'a_0004', '删除应用', 'DeleteAppVersion', '', 'post', '/v1/app_version/action/delete'),
('api_0016', 'm_0001', '商店管理', 'f_0001', '应用管理', 'a_0005', '发布应用', 'ReleaseAppVersion', '', 'post', '/v1/app_version/action/release'),
('api_0017', 'm_0001', '商店管理', 'f_0001', '应用管理', 'a_0006', '下架应用', 'CancelAppVersion', '', 'post', '/v1/app_version/action/cancel'),
('api_0018', 'm_0001', '商店管理', 'f_0002', '应用审核', 'a_0007', '审核提交', 'SubmitAppVersion', '', 'post', '/v1/app_version/action/submit'),
('api_0019', 'm_0001', '商店管理', 'f_0002', '应用审核', 'a_0008', '审核撤销', 'RecoverAppVersion', '', 'post', '/v1/app_version/action/recover'),
('api_0020', 'm_0001', '商店管理', 'f_0002', '应用审核', 'a_0009', '审核通过', 'PassAppVersion', '', 'post', '/v1/app_version/action/{role}/pass'),
('api_0021', 'm_0001', '商店管理', 'f_0002', '应用审核', 'a_0010', '审核拒绝', 'RejectAppVersion', '', 'post', '/v1/app_version/action/{role}/reject'),
('api_0022', 'm_0001', '商店管理', 'f_0003', '应用分类', 'a_0011', '查看全部分类', 'DescribeCategories', '', 'get', '/v1/categories'),
('api_0023', 'm_0001', '商店管理', 'f_0003', '应用分类', 'a_0012', '创建分类', 'CreateCategory', '', 'post', '/v1/categories'),
('api_0024', 'm_0001', '商店管理', 'f_0003', '应用分类', 'a_0013', '修改分类', 'ModifyCategory', '', 'patch', '/v1/categories'),
('api_0025', 'm_0001', '商店管理', 'f_0003', '应用分类', 'a_0014', '删除分类', 'DeleteCategories', '', 'delete', '/v1/categories'),
('api_0026', 'm_0002', '个人中心', 'f_0004', 'ssh key 管理', 'a_0015', '创建ssh key', 'CreateKeyPair', '', 'post', '/v1/clusters/key_pairs'),
('api_0027', 'm_0002', '个人中心', 'f_0004', 'ssh key 管理', 'a_0016', '查看ssh key', 'DescribeKeyPairs', '', 'get', '/v1/clusters/key_pairs'),
('api_0028', 'm_0002', '个人中心', 'f_0004', 'ssh key 管理', 'a_0017', '删除ssh key', 'DeleteKeyPairs', '', 'delete', '/v1/clusters/key_pairs'),
('api_0029', 'm_0002', '个人中心', 'f_0004', 'ssh key 管理', 'a_0018', '绑定ssh key', 'AttachKeyPairs', '', 'post', '/v1/clusters/key_pair/attach'),
('api_0030', 'm_0002', '个人中心', 'f_0004', 'ssh key 管理', 'a_0019', '解绑ssh key', 'DetachKeyPairs', '', 'post', '/v1/clusters/key_pair/detach'),
('api_0031', 'm_0003', '我的实例', 'f_0005', '应用实例管理', 'a_0020', '创建应用实例', 'CreateCluster', '', 'post', '/v1/clusters/create'),
('api_0032', 'm_0003', '我的实例', 'f_0005', '应用实例管理', 'a_0021', '创建应用实例', 'DescribeSubnets', '', 'get', '/v1/clusters/subnets'),
('api_0033', 'm_0003', '我的实例', 'f_0005', '应用实例管理', 'a_0022', '修改应用实例', 'ModifyClusterAttributes', '', 'post', '/v1/clusters/modify'),
('api_0034', 'm_0003', '我的实例', 'f_0005', '应用实例管理', 'a_0022', '修改应用实例', 'ModifyClusterNodeAttributes', '', 'post', '/v1/clusters/modify_nodes'),
('api_0035', 'm_0003', '我的实例', 'f_0005', '应用实例管理', 'a_0023', '删除应用实例', 'DeleteClusters', '', 'post', '/v1/clusters/delete'),
('api_0036', 'm_0003', '我的实例', 'f_0005', '应用实例管理', 'a_0024', '升级应用实例', 'UpgradeCluster', '', 'post', '/v1/clusters/upgrade'),
('api_0037', 'm_0003', '我的实例', 'f_0005', '应用实例管理', 'a_0025', '回滚应用实例', 'RollbackCluster', '', 'post', '/v1/clusters/rollback'),
('api_0038', 'm_0003', '我的实例', 'f_0005', '应用实例管理', 'a_0026', '纵向伸缩应用实例', 'ResizeCluster', '', 'post', '/v1/clusters/resize'),
('api_0039', 'm_0003', '我的实例', 'f_0005', '应用实例管理', 'a_0027', '横向伸缩应用实例', 'AddClusterNodes', '', 'post', '/v1/clusters/add_nodes'),
('api_0040', 'm_0003', '我的实例', 'f_0005', '应用实例管理', 'a_0027', '横向伸缩应用实例', 'DeleteClusterNodes', '', 'post', '/v1/clusters/delete_nodes'),
('api_0041', 'm_0003', '我的实例', 'f_0005', '应用实例管理', 'a_0028', '更新环境变量', 'UpdateClusterEnv', '', 'patch', '/v1/clusters/update_env'),
('api_0042', 'm_0003', '我的实例', 'f_0005', '应用实例管理', 'a_0029', '查看全部应用实例', 'DescribeClusters', '', 'get', '/v1/clusters'),
('api_0043', 'm_0003', '我的实例', 'f_0005', '应用实例管理', 'a_0029', '查看全部应用实例', 'DescribeClusterNodes', '', 'get', '/v1/clusters/nodes'),
('api_0044', 'm_0003', '我的实例', 'f_0005', '应用实例管理', 'a_0029', '查看全部应用实例', 'GetClusterStatistics', '', 'get', '/v1/clusters/statistics'),
('api_0045', 'm_0003', '我的实例', 'f_0005', '应用实例管理', 'a_0030', '关闭应用实例', 'StopClusters', '', 'post', '/v1/clusters/stop'),
('api_0046', 'm_0003', '我的实例', 'f_0005', '应用实例管理', 'a_0031', '启动应用实例', 'StartClusters', '', 'post', '/v1/clusters/start'),
('api_0047', 'm_0004', '账户与权限', 'f_0006', '用户管理', 'a_0032', '创建用户', 'CreateUser', '', 'post', '/v1/users'),
('api_0048', 'm_0004', '账户与权限', 'f_0006', '用户管理', 'a_0033', '查看全部用户', 'DescribeUsers', '', 'get', '/v1/users'),
('api_0049', 'm_0004', '账户与权限', 'f_0006', '用户管理', 'a_0034', '修改用户', 'ModifyUser', '', 'patch', '/v1/users'),
('api_0050', 'm_0004', '账户与权限', 'f_0006', '用户管理', 'a_0034', '修改用户', 'ChangePassword', '', 'post', '/v1/users/password:change'),
('api_0051', 'm_0004', '账户与权限', 'f_0006', '用户管理', 'a_0034', '修改用户', 'CreatePasswordReset', '', 'post', '/v1/users/password:reset'),
('api_0052', 'm_0004', '账户与权限', 'f_0006', '用户管理', 'a_0034', '修改用户', 'GetPasswordReset', '', 'get', '/v1/users/password/reset/{reset_id}'),
('api_0053', 'm_0004', '账户与权限', 'f_0006', '用户管理', 'a_0035', '删除用户', 'DeleteUsers', '', 'delete', '/v1/users'),
('api_0054', 'm_0004', '账户与权限', 'f_0007', '用户组管理', 'a_0036', '创建用户组', 'CreateGroup', '', 'post', '/v1/groups'),
('api_0055', 'm_0004', '账户与权限', 'f_0007', '用户组管理', 'a_0037', '查看全部用户组', 'DescribeGroups', '', 'get', '/v1/groups'),
('api_0056', 'm_0004', '账户与权限', 'f_0007', '用户组管理', 'a_0038', '修改用户组', 'ModifyGroup', '', 'patch', '/v1/groups'),
('api_0057', 'm_0004', '账户与权限', 'f_0007', '用户组管理', 'a_0039', '删除用户组', 'DeleteGroups', '', 'delete', '/v1/groups'),
('api_0058', 'm_0004', '账户与权限', 'f_0007', '用户组管理', 'a_0040', '加入用户组', 'JoinGroup', '', 'post', '/v1/groups:join'),
('api_0059', 'm_0004', '账户与权限', 'f_0007', '用户组管理', 'a_0041', '踢出用户组', 'LeaveGroup', '', 'post', '/v1/groups:leave'),
('api_0060', 'm_0005', '平台设置', 'f_0008', 'Job 管理', 'a_0042', '查看全部Job', 'DescribeJobs', '', 'get', '/v1/jobs'),
('api_0061', 'm_0006', '平台设置', 'f_0009', 'Task 管理', 'a_0043', '查看全部Task', 'DescribeTasks', '', 'get', '/v1/tasks'),
('api_0062', 'm_0006', '平台设置', 'f_0009', 'Task 管理', 'a_0044', '重试 Task', 'RetryTasks', '', 'post', '/v1/tasks/retry'),
('api_0063', 'm_0007', '平台设置', 'f_0010', '仓库管理', 'a_0045', '创建仓库', 'CreateRepo', '', 'post', '/v1/repos'),
('api_0064', 'm_0007', '平台设置', 'f_0010', '仓库管理', 'a_0045', '创建仓库', 'ValidateRepo', '', 'get', '/v1/repos/validate'),
('api_0065', 'm_0007', '平台设置', 'f_0010', '仓库管理', 'a_0046', '查看全部仓库', 'DescribeRepos', '', 'get', '/v1/repos'),
('api_0066', 'm_0007', '平台设置', 'f_0010', '仓库管理', 'a_0047', '修改仓库', 'ModifyRepo', '', 'patch', '/v1/repos'),
('api_0067', 'm_0007', '平台设置', 'f_0010', '仓库管理', 'a_0048', '删除仓库', 'DeleteRepos', '', 'delete', '/v1/repos'),
('api_0068', 'm_0007', '平台设置', 'f_0010', '仓库管理', 'a_0049', '同步应用', 'IndexRepo', '', 'post', '/v1/repos/index'),
('api_0069', 'm_0007', '平台设置', 'f_0010', '仓库管理', 'a_0050', '查看同步事件', 'DescribeRepoEvents', '', 'get', '/v1/repo_events'),
('api_0070', 'm_0008', '我的环境/个人中心-测试环境', 'f_0011', '环境管理', 'a_0051', '创建环境', 'CreateRuntime', '', 'post', '/v1/runtimes'),
('api_0071', 'm_0008', '我的环境/个人中心-测试环境', 'f_0011', '环境管理', 'a_0051', '创建环境', 'DescribeRuntimeProviderZones', '', 'get', '/v1/runtimes/zones'),
('api_0072', 'm_0008', '我的环境/个人中心-测试环境', 'f_0011', '环境管理', 'a_0051', '创建环境', 'GetRuntimeStatistics', '', 'get', '/v1/runtimes/statistics'),
('api_0073', 'm_0008', '我的环境/个人中心-测试环境', 'f_0011', '环境管理', 'a_0052', '查看全部环境', 'DescribeRuntimes', '', 'get', '/v1/runtimes'),
('api_0074', 'm_0008', '我的环境/个人中心-测试环境', 'f_0011', '环境管理', 'a_0053', '修改环境', 'ModifyRuntime', '', 'patch', '/v1/runtimes'),
('api_0075', 'm_0008', '我的环境/个人中心-测试环境', 'f_0011', '环境管理', 'a_0054', '删除环境', 'DeleteRuntimes', '', 'delete', '/v1/runtimes'),
('api_0076', 'm_0008', '我的环境/个人中心-测试环境', 'f_0012', '授权信息管理', 'a_0055', '创建授权信息', 'CreateRuntimeCredential', '', 'post', '/v1/runtimes/credentials'),
('api_0077', 'm_0008', '我的环境/个人中心-测试环境', 'f_0012', '授权信息管理', 'a_0056', '查看全部授权信息', 'DescribeRuntimeCredentials', '', 'get', '/v1/runtimes/credentials'),
('api_0078', 'm_0008', '我的环境/个人中心-测试环境', 'f_0012', '授权信息管理', 'a_0057', '修改授权信息', 'ModifyRuntimeCredential', '', 'patch', '/v1/runtimes/credentials'),
('api_0079', 'm_0008', '我的环境/个人中心-测试环境', 'f_0012', '授权信息管理', 'a_0058', '删除授权信息', 'DeleteRuntimeCredentials', '', 'delete', '/v1/runtimes/credentials'),
('api_0080', 'm_0009', '申请成为服务商', 'f_0013', '申请成为服务商', 'a_0059', '提交服务商认证', 'SubmitVendorVerifyInfo', '', 'post', '/v1/app_vendors'),
('api_0081', 'm_0009', '申请成为服务商', 'f_0013', '申请成为服务商', 'a_0059', '提交服务商认证', 'UploadVendorVerifyAttachment', '', '', ''),
('api_0082', 'm_0010', '应用服务商管理', 'f_0014', '入驻申请', 'a_0060', '通过服务商认证', 'PassVendorVerifyInfo', '', 'post', '/v1/app_vendors/pass'),
('api_0083', 'm_0010', '应用服务商管理', 'f_0014', '入驻申请', 'a_0061', '拒绝服务商认证', 'RejectVendorVerifyInfo', '', 'post', '/v1/app_vendors/reject'),
('api_0084', 'm_0010', '应用服务商管理', 'f_0014', '入驻申请', 'a_0062', '查看全部服务商认证', 'DescribeVendorVerifyInfos', '', 'get', '/v1/app_vendors'),
('api_0085', 'm_0010', '应用服务商管理', 'f_0014', '入驻申请', 'a_0062', '查看全部服务商认证', 'DescribeAppVendorStatistics', '', 'get', '/v1/app_vendors/app_vendor_statistics');

INSERT INTO `enable_action` (`enable_id`, `bind_id`, `action_id`) VALUES ('enable_00001', 'bind_00001', 'a_0001'),
('enable_00002', 'bind_00001', 'a_0006'),
('enable_00003', 'bind_00001', 'a_0008'),
('enable_00004', 'bind_00001', 'a_0009'),
('enable_00005', 'bind_00001', 'a_0010'),
('enable_00006', 'bind_00001', 'a_0011'),
('enable_00007', 'bind_00001', 'a_0012'),
('enable_00008', 'bind_00001', 'a_0013'),
('enable_00009', 'bind_00001', 'a_0014'),
('enable_00010', 'bind_00002', 'a_0015'),
('enable_00011', 'bind_00002', 'a_0016'),
('enable_00012', 'bind_00002', 'a_0017'),
('enable_00013', 'bind_00002', 'a_0018'),
('enable_00014', 'bind_00002', 'a_0019'),
('enable_00015', 'bind_00003', 'a_0020'),
('enable_00016', 'bind_00003', 'a_0021'),
('enable_00017', 'bind_00003', 'a_0022'),
('enable_00018', 'bind_00003', 'a_0023'),
('enable_00019', 'bind_00003', 'a_0024'),
('enable_00020', 'bind_00003', 'a_0025'),
('enable_00021', 'bind_00003', 'a_0026'),
('enable_00022', 'bind_00003', 'a_0027'),
('enable_00023', 'bind_00003', 'a_0028'),
('enable_00024', 'bind_00003', 'a_0029'),
('enable_00025', 'bind_00003', 'a_0030'),
('enable_00026', 'bind_00003', 'a_0031'),
('enable_00027', 'bind_00004', 'a_0032'),
('enable_00028', 'bind_00004', 'a_0033'),
('enable_00029', 'bind_00004', 'a_0034'),
('enable_00030', 'bind_00004', 'a_0035'),
('enable_00031', 'bind_00004', 'a_0036'),
('enable_00032', 'bind_00004', 'a_0037'),
('enable_00033', 'bind_00004', 'a_0038'),
('enable_00034', 'bind_00004', 'a_0039'),
('enable_00035', 'bind_00004', 'a_0040'),
('enable_00036', 'bind_00004', 'a_0041'),
('enable_00037', 'bind_00005', 'a_0042'),
('enable_00038', 'bind_00006', 'a_0043'),
('enable_00039', 'bind_00006', 'a_0044'),
('enable_00040', 'bind_00007', 'a_0045'),
('enable_00041', 'bind_00007', 'a_0046'),
('enable_00042', 'bind_00007', 'a_0047'),
('enable_00043', 'bind_00007', 'a_0048'),
('enable_00044', 'bind_00007', 'a_0049'),
('enable_00045', 'bind_00007', 'a_0050'),
('enable_00046', 'bind_00008', 'a_0051'),
('enable_00047', 'bind_00008', 'a_0052'),
('enable_00048', 'bind_00008', 'a_0053'),
('enable_00049', 'bind_00008', 'a_0054'),
('enable_00050', 'bind_00008', 'a_0055'),
('enable_00051', 'bind_00008', 'a_0056'),
('enable_00052', 'bind_00008', 'a_0057'),
('enable_00053', 'bind_00008', 'a_0058'),
('enable_00054', 'bind_00009', 'a_0060'),
('enable_00055', 'bind_00009', 'a_0061'),
('enable_00056', 'bind_00009', 'a_0062');



/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;
/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;
