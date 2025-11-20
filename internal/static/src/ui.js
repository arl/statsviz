import { plotMgr, drawPlots } from "./app.js";
import { updateVisibility } from "./nav.js";
import tippy from "tippy.js";

export const onClickPlotMaximize = (cfg) => (_gd, _ev) => {
  const clicked = plotMgr.plots.find((p) => p.name() === cfg.name);

  if (clicked.isMaximized()) {
    // Restore plots visibility based on current filters.
    updateVisibility();
    clicked.minimize();
  } else {
    // Hide all plots except the clicked one.
    plotMgr.plots.forEach((p) => {
      if (p !== clicked) p.setVisible(false);
    });
    clicked.maximize();
  }

  drawPlots(true);
};

export const onClickPlotInfo = (gd, ev) => {
  let button = ev.currentTarget;
  let val = button.getAttribute("data-val") === "true";

  const options = {
    allowHTML: true,
    trigger: "click",
  };

  const instance = tippy(ev.currentTarget, options);
  instance.setContent("<div>" + gd.infoText + "</div>");
  if (val) {
    instance.hide();
  } else {
    instance.show();
  }
  button.setAttribute("data-val", !val);
};
