import { plotMgr, drawPlots } from "./app.js";
import { updateVisibility } from "./nav.js";
import { Popover } from "bootstrap";

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
  const button = ev.currentTarget;
  button.setAttribute("tabindex", "0");
  button.setAttribute("role", "button");

  const popover = Popover.getOrCreateInstance(button, {
    html: true,
    trigger: "focus",
    placement: "bottom",
    container: "body",
    customClass: "plot-info-popover",
    content: gd.infoText,
  });
  popover.show();
};
