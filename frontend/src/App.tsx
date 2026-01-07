import { BrowserRouter, Routes, Route, Navigate } from "react-router-dom";
import LoginPage from "@/pages/auth/LoginPage";
import ProfilePage from "@/pages/profile/ProfilePage";
import { Toaster } from "@/components/ui/sonner";
import DashboardLayout from "./components/layout/DashboardLayout";
import { ProtectedRoute } from "./components/auth/ProtectedRoute";
import { PublicRoute } from "./components/auth/PublicRoute";

// Placeholder Page (Jika belum ada)
const AttendancePage = () => (
  <div className="p-4">Attendance Page Construction</div>
);

function App() {
  return (
    <BrowserRouter>
      <Routes>
        {/* Public Route - LOGIN */}
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

        {/* Protected Routes */}
        <Route
          path="/"
          element={
            <ProtectedRoute>
              <DashboardLayout />
            </ProtectedRoute>
          }
        >
          <Route path="dashboard" element={<div className="p-4">Dashboard Overview</div>} />
          <Route path="dashboard/attendance" element={<AttendancePage />} />
          <Route path="profile" element={<ProfilePage />} />
          <Route path="*" element={<div>404 Not Found</div>} />
        </Route>
      </Routes>

      <Toaster position="top-right" richColors />
    </BrowserRouter>
  );
}

export default App;
