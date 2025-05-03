import * as theme from "./theme.js";
import { plotMgr } from "./app.js";

export let paused = false;
export let gcEnabled = true;
export let timerange = 60;

export function initNav(onUpdate) {
  // Show GC toggle.
  const gcToggle = document.getElementById("gcToggle");

  gcToggle.checked = gcEnabled;
  gcToggle.addEventListener("change", (e) => {
    gcEnabled = !gcEnabled;
    gcToggle.checked = gcEnabled;
    onUpdate(true);
  });

  // Pause/Resume button.
  const pauseBtn = document.getElementById("pauseBtn");
  pauseBtn.addEventListener("click", (e) => {
    paused = !paused;
    pauseBtn.textContent = paused ? "Resume" : "Pause";
    pauseBtn.classList.toggle("active", paused);
    onUpdate(true);
  });

  // Dark mode toggle.
  const themeToggle = document.getElementById("themeToggle");
  themeToggle.addEventListener("change", (e) => {
    const themeMode = theme.getThemeMode();
    const newTheme = (themeMode === "dark" && "light") || "dark";
    localStorage.setItem("theme-mode", newTheme);

    theme.updateThemeMode();
    plotMgr.plots.forEach((p) => p.updateTheme());
    onUpdate(true);
  });

  // Time range selection
  const rangeInputs = document.querySelectorAll('input[name="range"]');

  rangeInputs.forEach((r, i) =>
    r.addEventListener("change", () => {
      if (r.checked) {
        rangeInputs[i].checked = true;
        const val = 60 * parseInt(rangeInputs[i].value, 10);
        timerange = val;
        onUpdate(true);
      }
    })
  );
  document.getElementById("range1").checked = true;
}
