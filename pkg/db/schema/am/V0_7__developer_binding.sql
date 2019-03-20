UPDATE role_module_binding
SET is_check_all = 0, bind_id = 'bid-developerm1'
WHERE role_id = 'developer' and module_id = 'm1';

INSERT INTO enable_action_bundle
(enable_id, bind_id, action_bundle_id) VALUES
('eid-developer1', 'bid-developerm1', 'm1.f1.a1'),
('eid-developer2', 'bid-developerm1', 'm1.f1.a2'),
('eid-developer3', 'bid-developerm1', 'm1.f1.a3'),
('eid-developer4', 'bid-developerm1', 'm1.f1.a4'),
('eid-developer5', 'bid-developerm1', 'm1.f1.a5'),
('eid-developer6', 'bid-developerm1', 'm1.f2.a1'),
('eid-developer7', 'bid-developerm1', 'm1.f2.a2')
;
