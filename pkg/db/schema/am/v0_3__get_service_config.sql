DELETE FROM module_api where api_method in ('GetAttachment', 'GetServiceConfig');
INSERT INTO module_api
(module_id,module_name,feature_id,feature_name,action_bundle_id,action_bundle_name,global_admin_action_bundle_visibility,isv_action_bundle_visibility,user_action_bundle_visibility,api_id,api_method,url_method,url) VALUES
('m0','默认权限','m0.f4','附件管理','m0.f4.a1','获取附件',0,0,0,'m0.f4.a1.api1','GetAttachment','get','/v1/attachments'),
('m0','默认权限','m0.f5','服务设置','m0.f5.a1','查看服务设置',0,0,0,'m0.f5.a1.api1','GetServiceConfig','post','/v1/service_configs/get')
;
