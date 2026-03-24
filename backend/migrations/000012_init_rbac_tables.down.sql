ALTER TABLE users ADD COLUMN role ENUM('SUPERADMIN', 'EMPLOYEE') DEFAULT 'EMPLOYEE';

UPDATE users u
JOIN roles r ON u.role_id = r.id
SET u.role = r.name;

ALTER TABLE users DROP FOREIGN KEY fk_user_role;
ALTER TABLE users DROP COLUMN role_id;

DROP TABLE IF EXISTS role_permissions;
DROP TABLE IF EXISTS permissions;
DROP TABLE IF EXISTS roles;
