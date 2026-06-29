import { useThemeStore } from "@/features/theme/store/themeStore";

function Toggle() {
  // 1. HOOKS arriba: leer stores, useState, useEffect
  const toggleTheme = useThemeStore((s) => s.toggleTheme);
  const theme = useThemeStore((s) => s.theme);

  // 2. RETURN abajo: lo que se ve (JSX)
  return (
    <div>
      <button onClick={toggleTheme}>{theme === "dark" ? "☀️" : "🌙"}</button>
    </div>
  );
}

export default Toggle;