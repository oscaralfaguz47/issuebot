import { Navigate, Outlet } from "react-router-dom";
import { useAuthStore } from "../store/authStore";
import { ROUTES } from "@/routes"
import LoadingScreen from "@/components/LoadingScreen";

function GuestRoute() {
  const session = useAuthStore((s) => s.session);
  const loading = useAuthStore((s) => s.loading);

if (loading) return <LoadingScreen />;

  if (session) {
    return <Navigate to={ROUTES.DASHBOARD} replace />;
  }

  return <Outlet />;
}

export default GuestRoute;