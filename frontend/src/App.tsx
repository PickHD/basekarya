import { BrowserRouter, Routes, Route, Navigate } from "react-router-dom";
import { Toaster } from "@/components/ui/sonner";

// Auth Components
import { ProtectedRoute } from "./components/auth/ProtectedRoute";
import { PublicRoute } from "./components/auth/PublicRoute";
import DashboardLayout from "./components/layout/DashboardLayout";

// Pages - Auth & General
import LoginPage from "@/pages/auth/LoginPage";
import DashboardPage from "@/pages/dashboard/DashboardPage";
import ProfilePage from "@/pages/profile/ProfilePage";
import AttendanceHistoryPage from "@/pages/dashboard/AttendanceHistoryPage";

import ReimbursementListPage from "@/pages/reimbursement/ReimbursementListPage";

// Pages - Admin
import EmployeeListPage from "@/pages/admin/EmployeeListPage";
import AttendanceRecapPage from "@/pages/admin/AttendanceRecapPage";

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
          {/* EMPLOYEE ROUTES */}
          <Route path="dashboard" element={<DashboardPage />} />
          <Route path="profile" element={<ProfilePage />} />

          <Route path="history" element={<AttendanceHistoryPage />} />

          <Route path="reimbursement">
            <Route index element={<ReimbursementListPage />} />
          </Route>

          {/* SUPERADMIN ROUTES */}
          <Route element={<ProtectedRoute allowedRoles={["SUPERADMIN"]} />}>
            <Route path="admin/employees" element={<EmployeeListPage />} />
            <Route path="admin/recap" element={<AttendanceRecapPage />} />
          </Route>

          {/* 404 Inside Layout */}
          <Route path="*" element={<div className="p-10">404 Not Found</div>} />
        </Route>
      </Routes>

      <Toaster position="top-right" richColors />
    </BrowserRouter>
  );
}

export default App;
