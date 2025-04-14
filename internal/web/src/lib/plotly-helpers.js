import { formatFunction } from "./format";
import { tooltip } from "./tooltip";

const handleInfoButton = (gd, ev) => {
  let button = ev.currentTarget;
  let val = button.getAttribute("data-val") === "true";

  const instance = tooltip.info(button);
  instance.setContent("<div>" + gd.getAttribute("infoText") + "</div>");
  if (val) {
    instance.hide();
  } else {
    instance.show();
  }
  button.setAttribute("data-val", (!val).toString());
};

export const createPlotlyConfig = (cfg) => {
  return {
    showEditInChartStudio: true,
    plotlyServerURL: "https://chart-studio.plotly.com",
    displaylogo: false,
    modeBarButtonsToRemove: [
      "zoom2d",
      "pan2d",
      "select2d",
      "lasso2d",
      "zoomIn2d",
      "zoomOut2d",
      "autoScale2d",
      "resetScale2d",
      "toggleSpikelines",
    ],
    modeBarButtonsToAdd: [
      {
        name: "info",
        title: "Plot info",
        icon: Plotly.Icons.info,
        click: handleInfoButton,
      },
    ],
    toImageButtonOptions: {
      format: "png", // one of png, svg, jpeg, webp
      filename: cfg.name,
      scale: 2, // Multiply title/legend/axis/canvas sizes by this factor
    },
  };
};

const themeColors = {
  light: {
    paper_bgcolor: "#f8f8f8",
    plot_bgcolor: "#ffffdd",
    font_color: "#434343",
  },
  dark: {
    paper_bgcolor: "#181a1c",
    plot_bgcolor: "#282a2c",
    font_color: "#fff",
  },
};

const currentTheme = themeColors.light;

const copyArrayOrNull = (o) => {
  return (Array.isArray(o) && [...o]) || undefined;
};

export const createPlotlyLayout = (cfg) => {
  // TODO: change based on theme

  let ticktext;
  let tickvals;

  if (cfg.type == "heatmap") {
    // TODO: maybe useless check on tickmode?
    if (cfg.layout.yaxis.tickmode == "array") {
      ticktext = new Array(cfg.layout.yaxis.ticktext.length);
      tickvals = new Array(cfg.layout.yaxis.tickvals.length);

      // Format yaxis ticks
      const formatYUnit = formatFunction(cfg.hover.yunit);
      cfg.layout.yaxis.ticktext.forEach((val, index) => {
        ticktext[index] = formatYUnit(val);
      });
      tickvals = copyArrayOrNull(cfg.layout.yaxis.tickvals);
    }
  }

  return {
    title: {
      y: 0.88,
      font: {
        family: "Roboto",
        size: 18,
      },
      text: cfg.title,
    },
    margin: {
      t: 80,
    },
    paper_bgcolor: currentTheme.paper_bgcolor,
    plot_bgcolor: currentTheme.plot_bgcolor,
    font: {
      color: currentTheme.font_color,
    },
    width: 630,
    height: 450,
    hovermode: "x",
    barmode: cfg.layout.barmode,
    xaxis: {
      tickformat: "%H:%M:%S",
      type: "date",
      fixedrange: true,
      autorange: false,
    },
    yaxis: {
      exponentformat: "SI",
      tickmode: cfg.layout.yaxis.tickmode,
      ticktext: ticktext,
      tickvals: tickvals,
      title: cfg.layout.yaxis.title,
      ticksuffix: cfg.layout.yaxis.ticksuffix,
      fixedrange: true,
    },
    showlegend: true,
    legend: {
      orientation: "h",
    },
  };
};

export const createPlotlyData = (cfg, sliced) => {
  // TODO: all this block can be done once, like it was before. Instead for now
  // we're doing it for every plot update.
  let data = [];
  switch (cfg.type) {
    case "heatmap":
      data.push({
        type: "heatmap",
        // x: null,
        y: cfg.buckets,
        // z: null,
        showlegend: false,
        colorscale: cfg.colorscale,
        custom_data: cfg.custom_data,
      });
      break;
    case "bar":
    case "scatter":
      cfg.subplots.forEach((subplot) => {
        data.push({
          type: cfg.type,
          name: subplot.name,
          hovertemplate: `<b>${subplot.unitfmt}</b>`,
        });
      });
      break;
  }

  switch (cfg.type) {
    case "heatmap":
      data[0].x = sliced.times;
      data[0].z = sliced.data;
      data[0].hoverinfo = "none";
      break;
    case "bar":
    case "scatter":
      // This is done for every data push
      for (let i = 0; i < data.length; i++) {
        data[i].x = sliced.times;
        data[i].y = sliced.data[i];
        data[i].stackgroup = cfg.subplots[i].stackgroup;
        data[i].hoveron = cfg.subplots[i].hoveron;
        data[i].type = cfg.subplots[i].type || cfg.type;
        data[i].marker = {
          color: cfg.subplots[i].color,
        };
      }
      break;
  }
  return data;
};
