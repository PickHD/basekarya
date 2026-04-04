ALTER TABLE permissions
ADD COLUMN permission_group_id INTEGER REFERENCES permission_groups(id),
ADD COLUMN description TEXT NULL,
ADD COLUMN display_name VARCHAR(255) NULL;