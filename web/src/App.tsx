import { useEffect } from "react";
import { BrowserRouter, Routes, Route } from "react-router-dom";
import Login from "./pages/Login";
import Dashboard from "./pages/Dashboard";
import ProtectedRoute from "./components/ProtectedRoute";
import { supabase } from "./lib/supabase";
import { useAuthStore } from "./store/authStore";

function App() {
  const setSession = useAuthStore((s) => s.setSession);
  const setLoading = useAuthStore((s) => s.setLoading);

  useEffect(() => {
    // 1. cargar la sesión actual al arrancar
    supabase.auth.getSession().then(({ data }) => {
      setSession(data.session);
      setLoading(false); 
    });

    // 2. escuchar cambios de sesión (login, logout, refresh)
    const { data: listener } = supabase.auth.onAuthStateChange((_event, session) => {
      setSession(session);
    });

    // 3. limpiar el listener cuando App se desmonta
    return () => listener.subscription.unsubscribe();
  }, [setSession, setLoading]);

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
      </Routes>
    </BrowserRouter>
  );
}

export default App;