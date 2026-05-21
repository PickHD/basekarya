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
  Briefcase,
  GraduationCap,
  Wallet,
  TrendingUp,
} from "lucide-react";
import type { MenuItem } from "./types";
import { PERMISSIONS } from "./permissions";

export const menuItems: MenuItem[] = [
  // Menu Utama
  {
    title: "Dashboard",
    href: "/dashboard",
    icon: LayoutDashboard,
    permission: PERMISSIONS.VIEW_SELF_ATTENDANCE,
    group: "Menu Utama",
  },

  // Kehadiran
  {
    title: "History Attendance",
    href: "/history",
    icon: History,
    permission: [PERMISSIONS.VIEW_SELF_ATTENDANCE],
    group: "Kehadiran",
  },
  {
    title: "Attendance Recap",
    href: "/admin/recap",
    icon: FileSpreadsheet,
    permission: PERMISSIONS.EXPORT_ATTENDANCE,
    group: "Kehadiran",
  },
  {
    title: "Overtime",
    href: "/overtime",
    icon: Clock,
    permission: [PERMISSIONS.VIEW_OVERTIME, PERMISSIONS.VIEW_SELF_OVERTIME],
    group: "Kehadiran",
  },

  // Pengajuan
  {
    title: "Leave Request",
    href: "/leave",
    icon: CalendarDays,
    permission: [PERMISSIONS.VIEW_LEAVE, PERMISSIONS.VIEW_SELF_LEAVE],
    group: "Pengajuan",
  },
  {
    title: "Loan",
    href: "/loan",
    icon: CreditCard,
    permission: [PERMISSIONS.VIEW_LOAN, PERMISSIONS.VIEW_SELF_LOAN],
    group: "Pengajuan",
  },
  {
    title: "Reimbursement",
    href: "/reimbursement",
    icon: Receipt,
    permission: [PERMISSIONS.VIEW_REIMBURSEMENT, PERMISSIONS.VIEW_SELF_REIMBURSEMENT],
    group: "Pengajuan",
  },

  // Keuangan
  {
    title: "Finance",
    href: "/finance",
    icon: Wallet,
    permission: PERMISSIONS.VIEW_FINANCE,
    group: "Keuangan",
  },
  {
    title: "Finance Dashboard",
    href: "/finance/dashboard",
    icon: TrendingUp,
    permission: PERMISSIONS.VIEW_FINANCE_DASHBOARD,
    group: "Keuangan",
  },

  // Karyawan
  {
    title: "Employees",
    href: "/admin/employees",
    icon: Users,
    permission: PERMISSIONS.VIEW_EMPLOYEE,
    group: "Karyawan",
  },
  {
    title: "Contracts",
    href: "/admin/contracts",
    icon: FileText,
    permission: PERMISSIONS.VIEW_CONTRACT,
    group: "Karyawan",
  },
  {
    title: "Payrolls",
    href: "/admin/payrolls",
    icon: Calculator,
    permission: PERMISSIONS.VIEW_PAYROLL,
    group: "Karyawan",
  },

  // Rekrutmen
  {
    title: "Recruitment",
    href: "/admin/requisitions",
    icon: Briefcase,
    permission: PERMISSIONS.VIEW_REQUISITION,
    group: "Rekrutmen",
  },
  {
    title: "Onboarding",
    href: "/admin/onboarding",
    icon: GraduationCap,
    permission: PERMISSIONS.VIEW_ONBOARDING,
    group: "Rekrutmen",
  },
  {
    title: "Onboarding Templates",
    href: "/admin/onboarding/templates",
    icon: GraduationCap,
    permission: PERMISSIONS.MANAGE_ONBOARDING_TEMPLATE,
    group: "Rekrutmen",
  },

  // Pengaturan
  {
    title: "Company Settings",
    href: "/admin/company-settings",
    icon: Settings,
    permission: PERMISSIONS.VIEW_COMPANY,
    group: "Pengaturan",
  },
  {
    title: "Roles & Permissions",
    href: "/admin/roles",
    icon: ShieldAlert,
    permission: PERMISSIONS.VIEW_ROLE,
    group: "Pengaturan",
  },
  {
    title: "Announcements",
    href: "/admin/announcements",
    icon: Megaphone,
    permission: PERMISSIONS.CREATE_ANNOUNCEMENT,
    group: "Pengaturan",
  },
];
