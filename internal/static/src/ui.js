import * as app from "./app.js";
import tippy from "tippy.js";

export const onClickPlotMaximize = (cfg) => (gd, ev) => {
  const clicked = app.allPlots.find((p) => p.name() === cfg.name);
  const isOnlyVisible = app.allPlots.every(
    (p) => p === clicked || !p.isVisible()
  );

  if (isOnlyVisible) {
    // Restore all plots.
    app.allPlots.forEach((p) => p.show());
  } else {
    // Hide all plots except the clicked one.
    app.allPlots.forEach((p) => {
      if (p !== clicked) p.hide();
    });
  }
  if (isOnlyVisible) {
    clicked.minimize();
  } else {
    clicked.maximize();
    app.updatePlots([clicked], true);
  }
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
