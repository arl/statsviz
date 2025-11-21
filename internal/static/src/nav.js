import * as theme from "./theme.js";
import { plotMgr } from "./app.js";

export let running = true;
export let gcEnabled = true;
export let timerange = 60;

export function updateVisibility() {
  const tagInputs = Array.from(
    document.querySelectorAll("#navCategories input[data-tag]")
  );
  const searchInput = document.getElementById("plot-search");

  const activeTags = tagInputs
    .filter((i) => i.checked)
    .map((i) => i.dataset.tag);

  const query = searchInput ? searchInput.value.trim() : "";

  plotMgr.plots.forEach((p) => {
    const matchesTag = activeTags.some((tag) => p.hasTag(tag));
    const matchesQuery = p.matches(query);
    p.setVisible(matchesTag && matchesQuery);
  });
}

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
  const searchInput = document.getElementById("plot-search");

  // Ensure initial state: all tags selected
  tagInputs.forEach((input) => (input.checked = true));

  // Listen for tag changes
  tagInputs.forEach((input) => {
    input.addEventListener("click", (e) => {
      if (e.altKey) {
        e.preventDefault();

        // Use setTimeout to ensure changes are applied after browser's default handling
        setTimeout(() => {
          const isSolo = tagInputs.every((i) => i.checked === (i === input));
          if (isSolo) {
            // If already solo, show all
            tagInputs.forEach((i) => (i.checked = true));
          } else {
            // Otherwise, solo this tag
            tagInputs.forEach((i) => (i.checked = i === input));
          }

          updateVisibility();
          onUpdate(true);
        }, 0);
      }
    });

    input.addEventListener("change", () => {
      updateVisibility();
      onUpdate(true);
    });
  });

  // Listen for search input changes
  searchInput.addEventListener("input", () => {
    // If a plot is maximized, minimize it before applying search input filtering.
    plotMgr.plots.forEach((p) => {
      if (p.isMaximized()) {
        p.minimize();
      }
    });
    updateVisibility();
    onUpdate(true);
  });

  // Apply initial tag filter
  updateVisibility();

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
