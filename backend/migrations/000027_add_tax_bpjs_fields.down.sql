ALTER TABLE payroll_details DROP COLUMN is_employer_borne, DROP COLUMN `group`, DROP COLUMN code;
ALTER TABLE companies DROP COLUMN bpjs_ketenagakerjaan_number, DROP COLUMN bpjs_kesehatan_number;
ALTER TABLE employees DROP COLUMN dependents_count, DROP COLUMN marital_status;
DROP TABLE IF EXISTS ptkp_configs;
DROP TABLE IF EXISTS pph21_term_configs;
DROP TABLE IF EXISTS bpjs_rate_configs;
