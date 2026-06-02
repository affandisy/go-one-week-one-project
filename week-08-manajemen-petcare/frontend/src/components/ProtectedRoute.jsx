import { Navigate, Outlet } from 'react-router-dom';
import { useAuthStore } from '../store/useAuthStore';

export default function ProtectedRoute({ allowedRoles }) {
  const { isAuthenticated, userRole } = useAuthStore();

  // 1. Jika belum login, tendang kembali ke halaman login
  if (!isAuthenticated) {
    return <Navigate to="/login" replace />;
  }

  // 2. (Opsional) Validasi Role based access
  if (allowedRoles && !allowedRoles.includes(userRole)) {
    // Jika role tidak diizinkan, arahkan ke dasbor bawaan atau halaman error
    return <Navigate to="/" replace />;
  }

  // 3. Jika aman, render halaman anak (Outlet)
  return <Outlet />;
}