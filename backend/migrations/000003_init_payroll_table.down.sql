DROP TABLE IF EXISTS payroll_details;
DROP TABLE IF EXISTS payrolls;

ALTER TABLE employees
DROP COLUMN base_salary,
DROP COLUMN bank_name,
DROP COLUMN bank_account_number,
DROP COLUMN bank_account_holder,
DROP COLUMN npwp;