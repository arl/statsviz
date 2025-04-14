import tippy, { followCursor } from "tippy.js";
import "tippy.js/dist/tippy.css";

let heatmapInstance;

const heatmapOptions = {
  followCursor: true,
  trigger: "manual",
  allowHTML: true,
  plugins: [followCursor],
};

const infoOptions = {
  allowHTML: true,
  trigger: "click",
};

export const tooltip = {
  heatmap: function get() {
    if (!heatmapInstance) {
      heatmapInstance = tippy(document.body, heatmapOptions);
    }
    return heatmapInstance;
  },
  info: function get(element) {
    return tippy(element, infoOptions);
  },
};
