import * as theme from "./theme.js";
import { plotMgr } from "./app.js";

export let running = true;
export let gcEnabled = true;
export let timerange = 60;

export function initNav(onUpdate) {
  // Show GC toggle.
  const gcbtn = document.getElementById("btn-gc-events");

  gcbtn.checked = gcEnabled;
  gcbtn.addEventListener("change", (e) => {
    gcEnabled = !gcEnabled;
    gcbtn.checked = gcEnabled;
    onUpdate(true);
  });

  // Pause/Resume button.
  const playbtn = document.getElementById("btn-play");
  playbtn.addEventListener("click", (e) => {
    running = !running;
    playbtn.checked = running;
    onUpdate(true);
  });

  // Dark mode toggle.
  const themebtn = document.getElementById("btn-darkmode");
  themebtn.addEventListener("change", (e) => {
    const themeMode = theme.getThemeMode();
    const newTheme = (themeMode === "dark" && "light") || "dark";
    localStorage.setItem("theme-mode", newTheme);

    theme.updateThemeMode();
    plotMgr.plots.forEach((p) => p.updateTheme());
    onUpdate(true);
  });

  // Plot tags toggling.
  const tagInputs = Array.from(
    document.querySelectorAll("#navCategories input[data-tag]")
  );

  // Ensure initial state: all tags selected
  tagInputs.forEach((input) => (input.checked = true));

  // Update plot visibility based on selected tags
  const updateByTags = () => {
    const activeTags = tagInputs
      .filter((i) => i.checked)
      .map((i) => i.dataset.tag);

    plotMgr.plots.forEach((p) => {
      p.setVisible(activeTags.some((tag) => p.hasTag(tag)));
    });
  };

  // Listen for tag changes
  tagInputs.forEach((input) => {
    input.addEventListener("change", () => {
      updateByTags();
      onUpdate(true);
    });
  });

  // Apply initial tag filter
  updateByTags();

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
