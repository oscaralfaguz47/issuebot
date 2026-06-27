import { useEffect } from "react";
import { BrowserRouter, Routes, Route } from "react-router-dom";
import Login from "./pages/Login";
import Dashboard from "./pages/Dashboard";
import ProtectedRoute from "./components/ProtectedRoute";
import { supabase } from "./lib/supabase";
import { useAuthStore } from "./store/authStore";
import Jobs from "./pages/Jobs";
import { getMemberships } from "./lib/api";

function App() {
  const setSession = useAuthStore((s) => s.setSession);
  const setLoading = useAuthStore((s) => s.setLoading);
  const setMemberships = useAuthStore((s) => s.setMemberships);
  useEffect(() => {
    // 1. cargar la sesión actual al arrancar
    supabase.auth.getSession().then(({ data }) => {
      setSession(data.session);
      setLoading(false);
      if (data.session) {
        getMemberships().then(setMemberships).catch(console.error);
      }
    });

    // 2. escuchar cambios de sesión (login, logout, refresh)
    const { data: listener } = supabase.auth.onAuthStateChange((_event, session) => {
      setSession(session);
      if (session) {
        getMemberships().then(setMemberships).catch(console.error);
      }
    });

    // 3. limpiar el listener cuando App se desmonta
    return () => listener.subscription.unsubscribe();
  }, [setSession, setLoading, setMemberships]);

  return (
    <BrowserRouter>
      <Routes>
        <Route path="/login" element={<Login />} />
        <Route
          path="/"
          element={
            <ProtectedRoute>
              <Dashboard />
            </ProtectedRoute>
          }
        />
        <Route
          path="/jobs"
          element={
            <ProtectedRoute>
              <Jobs />
            </ProtectedRoute>
          }
        />
      </Routes>
    </BrowserRouter>
  );
}

export default App;