import { BrowserRouter, Routes, Route, Navigate } from "react-router-dom";
import LoginPage from "@/pages/auth/LoginPage";
import { Toaster } from "@/components/ui/sonner";
import DashboardLayout from "./components/layout/DashboardLayout";

// Dummy Pages
const DashboardPage = () => (
  <h1 className="text-2xl font-bold">Dashboard Overview</h1>
);
const AttendancePage = () => (
  <h1 className="text-2xl font-bold">Attendance History</h1>
);

function App() {
  return (
    <BrowserRouter>
      <Routes>
        {/* Public Route - LOGIN */}
        <Route path="/login" element={<LoginPage />} />

        {/* Root redirect to login */}
        <Route path="/" element={<Navigate to="/login" replace />} />

        {/* Protected Routes */}
        <Route path="/dashboard" element={<DashboardLayout />}>
          <Route index element={<DashboardPage />} />
          <Route path="attendance" element={<AttendancePage />} />
          <Route path="*" element={<div>404 Not Found</div>} />
        </Route>
      </Routes>

      <Toaster position="top-right" richColors />
    </BrowserRouter>
  );
}

export default App;
