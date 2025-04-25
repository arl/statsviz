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

const themeToggle = document.getElementById("themeToggle");

const navbar = document.getElementById("top-navbar");

/**
 * Set light or dark theme
 */
export const updateThemeMode = () => {
  const themeMode = getThemeMode();
  if (themeMode === "dark") {
    document.body.classList.add("dark-mode");
    navbar.setAttribute("data-bs-theme", "dark");
    themeToggle.setAttribute("checked", "");
  } else {
    document.body.classList.remove("dark-mode");
    navbar.setAttribute("data-bs-theme", "light");
    themeToggle.removeAttribute("checked");
  }
};
