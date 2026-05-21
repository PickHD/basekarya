ALTER TABLE ref_leave_types DROP INDEX idx_ref_leave_types_name_company_id;
ALTER TABLE ref_leave_types ADD UNIQUE INDEX name (name);
