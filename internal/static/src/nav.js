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
  const pauseToggle = document.getElementById("pauseToggle");
  pauseToggle.addEventListener("click", (e) => {
    paused = !paused;
    pauseToggle.classList.toggle("active", paused);
    pauseToggle == paused;
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

  // Plot tags toggling
  const tagInputs = document.querySelectorAll("#navCategories input[data-tag]");

  // Ensure initial state: all checked → all plots shown
  tagInputs.forEach((input) => {
    input.checked = true; // redundant if HTML has checked, but safe
  });

  tagInputs.forEach((input) => {
    const tag = input.dataset.tag;

    // On each toggle, show or hide matching plots
    input.addEventListener("change", () => {
      if (input.checked) {
        plotMgr.plots.forEach((p) => {
          if (p.hasTag(tag)) p.show();
        });
      } else {
        plotMgr.plots.forEach((p) => {
          if (p.hasTag(tag)) p.hide();
        });
      }
      // Redraw after visibility change
      onUpdate(true);
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
        onUpdate(true);
      }
    })
  );
  document.getElementById("range1").checked = true;
}
