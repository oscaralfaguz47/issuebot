import { Outlet } from "react-router-dom";
import Toggle from "./Toggle";

function AppLayout() {
  return (
    <div className="mx-auto max-w-[1200px] px-6 py-8">
      {/* Header should go here */}
      <main>
        <Toggle />       <Outlet />
      </main>
    </div>
  );
}

export default AppLayout;