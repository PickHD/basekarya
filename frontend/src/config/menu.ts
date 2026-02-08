import {
  LayoutDashboard,
  History,
  Users,
  FileSpreadsheet,
  Receipt,
  Calculator,
  CalendarDays,
} from "lucide-react";
import type { MenuItem } from "./types";

export const generalMenu: MenuItem[] = [
  {
    title: "Dashboard",
    href: "/dashboard",
    icon: LayoutDashboard,
  },
  {
    title: "My History",
    href: "/history",
    icon: History,
  },
  {
    title: "My Reimbursement",
    href: "/reimbursement",
    icon: Receipt,
  },
  {
    title: "My Leave Request",
    href: "/leave",
    icon: CalendarDays,
  },
];

export const adminMenu: MenuItem[] = [
  {
    title: "Attendance Recap",
    href: "/admin/recap",
    icon: FileSpreadsheet,
  },
  {
    title: "Employees",
    href: "/admin/employees",
    icon: Users,
  },
  {
    title: "Reimbursements",
    href: "/reimbursement",
    icon: Receipt,
  },
  {
    title: "Payrolls",
    href: "/admin/payrolls",
    icon: Calculator,
  },
  {
    title: "Leave Requests",
    href: "/leave",
    icon: CalendarDays,
  },
];
