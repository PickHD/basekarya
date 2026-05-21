ALTER TABLE ref_leave_types DROP INDEX name;
ALTER TABLE ref_leave_types ADD UNIQUE INDEX idx_ref_leave_types_name_company_id (name, company_id);
