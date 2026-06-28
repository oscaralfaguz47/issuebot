import { Outlet } from "react-router-dom";

function RootLayout() {
  return (
    <div className="min-h-screen bg-bg text-fg">
       <Outlet />
    </div>
  );
}

export default RootLayout;