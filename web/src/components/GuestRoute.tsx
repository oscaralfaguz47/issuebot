import { Navigate, Outlet } from "react-router-dom";
import { useAuthStore } from "../store/authStore";

function GuestRoute() {
  const session = useAuthStore((s) => s.session);
  const loading = useAuthStore((s) => s.loading);

  if (loading) {
    return <p>Loading...</p>;
  }

  if (session) {
    return <Navigate to="/dashboard" replace />;
  }

  return <Outlet />;
}

export default GuestRoute;