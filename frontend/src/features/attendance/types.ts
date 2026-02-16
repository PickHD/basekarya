import type { Meta } from "@/types/api";

export interface ClockPayload {
  latitude: number;
  longitude: number;
  image_base64: string;
  address?: string;
  notes?: string;
}

export interface ClockResponse {
  message: string;
  data: {
    type: "CHECK_IN" | "CHECK_OUT";
    status: string;
    time: string;
  };
}

export interface TodayAttendanceResponse {
  status: "ABSENT" | "PRESENT" | "LATE";
  type: "NONE" | "CHECK_IN" | "COMPLETED";
  check_in_time?: string;
  check_out_time?: string;
  work_duration?: string;
}

export interface Shift {
  id: number;
  name: string;
  start_time: string;
  end_time: string;
  created_at: string;
}

export type AttendanceStatus = "PRESENT" | "LATE" | "EXCUSED" | "ABSENT";

export interface AttendanceLog {
  id: number;
  employee_id: number;
  shift_id: number;

  date: string;

  check_in_time: string;
  check_in_lat: number;
  check_in_long: number;
  check_in_image_url: string;
  check_in_address: string;

  check_out_time: string | null;
  check_out_lat: number | null;
  check_out_long: number | null;
  check_out_image_url: string | null;
  check_out_address: string | null;

  status: AttendanceStatus;
  is_suspicious: boolean;
  notes: string;

  created_at: string;
  updated_at: string;

  shift: Shift;
}

export interface HistoryResponse {
  data: AttendanceLog[];
  meta: Meta;
}
