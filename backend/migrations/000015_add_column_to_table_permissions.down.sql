ALTER TABLE permissions
DROP COLUMN IF EXISTS permission_group_id,
DROP COLUMN IF EXISTS description,
DROP COLUMN IF EXISTS display_name;