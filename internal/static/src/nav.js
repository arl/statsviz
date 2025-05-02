import { updatePlots } from "./app.js";
import * as theme from "./theme.js";

export function initNav(allPlots) {
  // move gcToggle, pauseBtn, themeToggle, rangeInputs handling here
  // exactly as in app.js but using passed-in allPlots

  // Show GC toggle.
  const gcToggle = document.getElementById("gcToggle");

  gcToggle.checked = show_gc;
  gcToggle.addEventListener("change", (e) => {
    show_gc = !show_gc;
    gcToggle.checked = show_gc;
    updatePlots(allPlots, true);
  });

  // Pause/Resume button.
  const pauseBtn = document.getElementById("pauseBtn");
  pauseBtn.addEventListener("click", (e) => {
    paused = !paused;
    pauseBtn.textContent = paused ? "Resume" : "Pause";
    pauseBtn.classList.toggle("active", paused);
    updatePlots(allPlots, true);
  });

  // Dark mode toggle.
  const themeToggle = document.getElementById("themeToggle");
  themeToggle.addEventListener("change", (e) => {
    const themeMode = theme.getThemeMode();
    const newTheme = (themeMode === "dark" && "light") || "dark";
    localStorage.setItem("theme-mode", newTheme);

    theme.updateThemeMode();

    allPlots.forEach((plot) => {
      plot.updateTheme();
    });
  });

  // Time range selection
  const rangeInputs = document.querySelectorAll('input[name="range"]');

  rangeInputs.forEach((r, i) =>
    r.addEventListener("change", () => {
      if (r.checked) {
        rangeInputs[i].checked = true;
        const val = 60 * parseInt(rangeInputs[i].value, 10);
        timerange = val;
        updatePlots(allPlots, true);
      }
    })
  );
  document.getElementById("range1").checked = true;
}

export let paused = false,
  show_gc = true,
  timerange = 60;
