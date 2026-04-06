import {
  LayoutDashboard,
  History,
  Users,
  FileSpreadsheet,
  Receipt,
  Calculator,
  CalendarDays,
  Settings,
  CreditCard,
  Clock,
  ShieldAlert,
  Megaphone,
  FileText,
} from "lucide-react";
import type { MenuItem } from "./types";
import { PERMISSIONS } from "./permissions";

export const menuItems: MenuItem[] = [
  {
    title: "Announcements",
    href: "/admin/announcements",
    icon: Megaphone,
    permission: PERMISSIONS.CREATE_ANNOUNCEMENT
  },
  {
    title: "Attendance Recap",
    href: "/admin/recap",
    icon: FileSpreadsheet,
    permission: PERMISSIONS.EXPORT_ATTENDANCE
  },
  {
    title: "Company Settings",
    href: "/admin/company-settings",
    icon: Settings,
    permission: PERMISSIONS.VIEW_COMPANY
  },
  {
    title: "Contracts",
    href: "/admin/contracts",
    icon: FileText,
    permission: PERMISSIONS.VIEW_CONTRACT
  },
  {
    title: "Dashboard",
    href: "/dashboard",
    icon: LayoutDashboard,
    permission: PERMISSIONS.VIEW_SELF_ATTENDANCE
  },
  {
    title: "Employees",
    href: "/admin/employees",
    icon: Users,
    permission: PERMISSIONS.VIEW_EMPLOYEE
  },
  {
    title: "History Attendance",
    href: "/history",
    icon: History,
    permission: [PERMISSIONS.VIEW_SELF_ATTENDANCE]
  },
  {
    title: "Leave Request",
    href: "/leave",
    icon: CalendarDays,
    permission: [PERMISSIONS.VIEW_LEAVE, PERMISSIONS.VIEW_SELF_LEAVE]
  },
  {
    title: "Loan",
    href: "/loan",
    icon: CreditCard,
    permission: [PERMISSIONS.VIEW_LOAN, PERMISSIONS.VIEW_SELF_LOAN]
  },
  {
    title: "Overtime",
    href: "/overtime",
    icon: Clock,
    permission: [PERMISSIONS.VIEW_OVERTIME, PERMISSIONS.VIEW_SELF_OVERTIME]
  },
  {
    title: "Payrolls",
    href: "/admin/payrolls",
    icon: Calculator,
    permission: PERMISSIONS.VIEW_PAYROLL
  },
  {
    title: "Reimbursement",
    href: "/reimbursement",
    icon: Receipt,
    permission: [PERMISSIONS.VIEW_REIMBURSEMENT, PERMISSIONS.VIEW_SELF_REIMBURSEMENT]
  },
  {
    title: "Roles & Permissions",
    href: "/admin/roles",
    icon: ShieldAlert,
    permission: PERMISSIONS.VIEW_ROLE
  },
];
