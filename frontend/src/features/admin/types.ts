export interface Employee {
  id: number;
  full_name: string;
  nik: string;
  username: string;
  department_name: string;
  shift_name: string;
  base_salary: number;
}

export interface AttendanceRecap {
  id: number;
  date: string;
  employee_name: string;
  nik: string;
  department: string;
  shift: string;
  check_in_time: string;
  check_out_time: string;
  status: string;
  work_duration: string;
}

export interface CreateEmployeePayload {
  username: string;
  full_name: string;
  nik: string;
  department_id: number;
  shift_id: number;
  base_salary: number;
}

export interface LookupItem {
  id: number;
  name: string;
}

export interface DashboardStats {
  total_employees: number;
  present_today: number;
  late_today: number;
  absent_today: number;
}
