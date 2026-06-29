import { useEffect } from "react";
import { BrowserRouter, Routes, Route } from "react-router-dom";
import { supabase } from "@/lib/supabase";
import { useAuthStore, ProtectedRoute, GuestRoute, getMemberships } from "@/features/auth";
import AppLayout from "@/components/AppLayout";
import RootLayout from "@/components/RootLayout";
import { useThemeStore } from "@/features/theme";
import { ROUTES } from "@/routes"

import { lazy, Suspense } from "react";
import LoadingScreen from "./components/LoadingScreen";

import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { ReactQueryDevtools } from "@tanstack/react-query-devtools";

const Dashboard = lazy(() => import("@/pages/app/Dashboard"));
const Jobs = lazy(() => import("@/pages/app/Jobs"));
const Landing = lazy(() => import("@/pages/public/Landing"));
const Login = lazy(() => import("@/pages/public/Login"));
const NotFound = lazy(() => import("@/pages/public/NotFound"));

const queryClient = new QueryClient();

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
    <QueryClientProvider client={queryClient}>
      <BrowserRouter>
        <Suspense fallback={<LoadingScreen />}>
          <Routes>
            <Route element={<RootLayout />}>
              <Route path={ROUTES.LANDING} element={<Landing />} />
              <Route element={<GuestRoute />}>
                <Route path={ROUTES.LOGIN} element={<Login />} />
              </Route>

              <Route element={<ProtectedRoute />}>      {/* responsability: protect */}
                <Route element={<AppLayout />}>            {/* responsability: visual shell */}
                  <Route path={ROUTES.DASHBOARD} element={<Dashboard />} />
                  <Route path={ROUTES.JOBS} element={<Jobs />} />
                </Route>
              </Route>
              <Route path={ROUTES.NOTFOUND} element={<NotFound />} />
            </Route>
          </Routes>
        </Suspense>
      </BrowserRouter>
      <ReactQueryDevtools initialIsOpen={false} />
    </QueryClientProvider>
  );
}

export default App;