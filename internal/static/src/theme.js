/**
 * Get color theme based on previous user choice or browser theme
 */
export const getThemeMode = () => {
  let themeMode = localStorage.getItem("theme-mode");

  if (themeMode === null) {
    const isDark =
      window.matchMedia &&
      window.matchMedia("(prefers-color-scheme: dark)").matches;
    themeMode = (isDark && "dark") || "light";

    localStorage.setItem("theme-mode", themeMode);
  }

  return themeMode;
};

/**
 * Set light or dark theme
 */
export const updateThemeMode = () => {
  const themeMode = getThemeMode();
  console.log("themeMode", themeMode);
  if (themeMode === "dark") {
    document.body.classList.add("dark-mode");
    document.getElementById("navbar").setAttribute("data-bs-theme", "dark");
    document.getElementById("dark_mode_switch").setAttribute("checked", "");
  } else {
    document.body.classList.remove("dark-mode");
    document.getElementById("navbar").setAttribute("data-bs-theme", "light");
    document.getElementById("dark_mode_switch").removeAttribute("checked");
  }
};
