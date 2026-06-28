import { useEffect } from "react";
import { BrowserRouter, Routes, Route } from "react-router-dom";
import Login from "./pages/Login";
import Dashboard from "./pages/Dashboard";
import ProtectedRoute from "./components/ProtectedRoute";
import { supabase } from "./lib/supabase";
import { useAuthStore } from "./store/authStore";
import Jobs from "./pages/Jobs";
import { getMemberships } from "./lib/api";
import AppLayout from "./components/AppLayout";
import { useThemeStore } from "./store/themeStore";
import Landing from "./pages/Landing";
import RootLayout from "./components/RootLayout";
import GuestRoute from "./components/GuestRoute";
import NotFound from "./pages/NotFound";

function App() {
  const setSession = useAuthStore((s) => s.setSession);
  const setLoading = useAuthStore((s) => s.setLoading);
  const setMemberships = useAuthStore((s) => s.setMemberships);
  const theme = useThemeStore((s) => s.theme);

  useEffect(() => {
    // 1. load the current session on mount
    supabase.auth.getSession().then(({ data }) => {
      setSession(data.session);
      setLoading(false);
      if (data.session) {
        getMemberships().then(setMemberships).catch(console.error);
      }
    });

    // 2. listen session changes (login, logout, refresh)
    const { data: listener } = supabase.auth.onAuthStateChange((_event, session) => {
      setSession(session);
      if (session) {
        getMemberships().then(setMemberships).catch(console.error);
      }
    });

    // 3. clean up the listener when the component unmounts
    return () => listener.subscription.unsubscribe();
  }, [setSession, setLoading, setMemberships]);

  useEffect(() => {
    document.documentElement.classList.toggle("dark", theme === "dark");
  }, [theme]);

  return (
    <BrowserRouter>
      <Routes>
        <Route element={<RootLayout />}>
          <Route path="/" element={<Landing />} />
          <Route element={<GuestRoute />}>
            <Route path="/login" element={<Login />} />
          </Route>

          <Route element={<ProtectedRoute />}>      {/* responsability: protect */}
            <Route element={<AppLayout />}>            {/* responsability: visual shell */}
              <Route path="/dashboard" element={<Dashboard />} />
              <Route path="/jobs" element={<Jobs />} />
            </Route>
          </Route>
          <Route path="*" element={<NotFound />} />
        </Route>
      </Routes>
    </BrowserRouter>
  );
}

export default App;