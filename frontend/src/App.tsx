import { BrowserRouter, Routes, Route, Navigate } from "react-router-dom";
import { Toaster } from "@/components/ui/sonner";

// Auth Components
import { ProtectedRoute } from "@/components/auth/ProtectedRoute";
import { PublicRoute } from "@/components/auth/PublicRoute";
import DashboardLayout from "@/components/layout/DashboardLayout";

// Pages - Auth & General
import LoginPage from "@/pages/auth/LoginPage";
import DashboardPage from "@/pages/dashboard/DashboardPage";
import ProfilePage from "@/pages/profile/ProfilePage";
import AttendanceHistoryPage from "@/pages/dashboard/AttendanceHistoryPage";

import ReimbursementListPage from "@/pages/reimbursement/ReimbursementListPage";
import LoanListPage from "@/pages/loan/LoanListPage";
import OvertimeListPage from "@/pages/overtime/OvertimeListPage";

// Pages - Admin
import EmployeeListPage from "@/pages/admin/EmployeeListPage";
import AttendanceRecapPage from "@/pages/admin/AttendanceRecapPage";
import PayrollListPage from "@/pages/payroll/PayrollListPage";
import LeaveListPage from "@/pages/leave/LeaveListPage";
import CompanySettingsPage from "@/pages/admin/CompanySettingsPage";
import RoleListPage from "@/pages/admin/RoleListPage";

import { PERMISSIONS } from "@/config/permissions";

function App() {
  return (
    <BrowserRouter>
      <Routes>
        {/* === PUBLIC ROUTES === */}
        <Route
          path="/login"
          element={
            <PublicRoute>
              <LoginPage />
            </PublicRoute>
          }
        />

        {/* Root redirect to login */}
        <Route path="/" element={<Navigate to="/login" replace />} />

        {/* === PROTECTED ROUTES (Global) === */}
        <Route
          element={
            <ProtectedRoute>
              <DashboardLayout />
            </ProtectedRoute>
          }
        >
          <Route path="dashboard" element={<DashboardPage />} />
          <Route path="profile" element={<ProfilePage />} />

          <Route element={<ProtectedRoute requiredPermissions={[PERMISSIONS.VIEW_ATTENDANCE, PERMISSIONS.VIEW_SELF_ATTENDANCE]} />}>
            <Route path="history" element={<AttendanceHistoryPage />} />
          </Route>

          <Route element={<ProtectedRoute requiredPermissions={[PERMISSIONS.VIEW_REIMBURSEMENT, PERMISSIONS.VIEW_SELF_REIMBURSEMENT]} />}>
            <Route path="reimbursement">
              <Route index element={<ReimbursementListPage />} />
            </Route>
          </Route>

          <Route element={<ProtectedRoute requiredPermissions={[PERMISSIONS.VIEW_LOAN, PERMISSIONS.VIEW_SELF_LOAN]} />}>
            <Route path="loan">
              <Route index element={<LoanListPage />} />
            </Route>
          </Route>

          <Route element={<ProtectedRoute requiredPermissions={[PERMISSIONS.VIEW_OVERTIME, PERMISSIONS.VIEW_SELF_OVERTIME]} />}>
            <Route path="overtime">
              <Route index element={<OvertimeListPage />} />
            </Route>
          </Route>

          <Route element={<ProtectedRoute requiredPermissions={[PERMISSIONS.VIEW_LEAVE, PERMISSIONS.VIEW_SELF_LEAVE]} />}>
            <Route path="leave">
              <Route index element={<LeaveListPage />} />
            </Route>
          </Route>

          {/* ADMINISTRATIVE ROUTES */}
          <Route 
            element={<ProtectedRoute requiredPermissions={[
              PERMISSIONS.VIEW_EMPLOYEE, 
              PERMISSIONS.VIEW_ATTENDANCE,
              PERMISSIONS.VIEW_PAYROLL,
              PERMISSIONS.VIEW_COMPANY,
              PERMISSIONS.VIEW_ROLE
            ]} />}
          >
            <Route path="admin/employees" element={<EmployeeListPage />} />
            <Route path="admin/roles" element={<RoleListPage />} />
            <Route path="admin/recap" element={<AttendanceRecapPage />} />
            <Route path="admin/payrolls" element={<PayrollListPage />} />
            <Route
              path="admin/company-settings"
              element={<CompanySettingsPage />}
            />

            {/* 404 Inside Layout */}
            <Route
              path="*"
              element={<div className="p-10">404 Not Found</div>}
            />
          </Route>
        </Route>
      </Routes>

      <Toaster position="top-right" richColors />
    </BrowserRouter>
  );
}

export default App;
