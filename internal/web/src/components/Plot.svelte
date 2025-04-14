<script>
  import { onMount, onDestroy } from "svelte";
  import { metricStore } from "../stores/metrics";
  import { configStore } from "../stores/config";
  import Plotly from "plotly.js-dist-min";
  import {
    createPlotlyConfig,
    createPlotlyData,
    createPlotlyLayout,
  } from "../lib/plotly-helpers";

  import { formatFunction } from "../lib/format";
  import { tooltip } from "../lib/tooltip";

  export let name;
  let plotDiv;

  let unsubConfig;
  let unsubMetrics;

  onMount(() => {
    unsubConfig?.();

    unsubConfig = configStore.subscribe((cfg) => {
      if (!cfg) return;
      const serieCfg = cfg.series.find((s) => s.name === name);
      if (!serieCfg) return;

      const plotlyConfig = createPlotlyConfig(serieCfg);
      const plotlyLayout = createPlotlyLayout(serieCfg);

      // Directly attach tooltip to the dom element.
      plotDiv.setAttribute(
        "infoText",
        serieCfg.infoText
          .split("\n")
          .map((line) => `<p>${line}</p>`)
          .join("")
      );

      Plotly.newPlot(plotDiv, [], plotlyLayout, plotlyConfig);

      unsubMetrics?.();
      unsubMetrics = metricStore(name).subscribe((data) => {
        const pdata = createPlotlyData(serieCfg, data);
        plotlyLayout.xaxis.range = data.xrange;
        Plotly.react(plotDiv, pdata, plotlyLayout, plotlyConfig);
        if (serieCfg.type === "heatmap") {
          installHeatmapTooltip(serieCfg);
        }
      });
    });
  });

  function installHeatmapTooltip(cfg) {
    const instance = tooltip.heatmap();
    const hover = cfg.hover;
    const formatYUnit = formatFunction(hover.yunit);

    // Remove any previously attached Plotly listeners.
    plotDiv.removeAllListeners("plotly_hover");
    plotDiv.removeAllListeners("plotly_unhover");

    plotDiv.on("plotly_hover", function (data) {
      const pt2txt = (d) => {
        let bucket;
        if (d.y == 0) {
          const yhigh = formatYUnit(d.data.custom_data[d.y]);
          bucket = `(-Inf, ${yhigh})`;
        } else if (d.y == d.data.custom_data.length - 1) {
          const ylow = formatYUnit(d.data.custom_data[d.y]);
          bucket = `[${ylow}, +Inf)`;
        } else {
          const ylow = formatYUnit(d.data.custom_data[d.y - 1]);
          const yhigh = formatYUnit(d.data.custom_data[d.y]);
          bucket = `[${ylow}, ${yhigh})`;
        }

        return `
  <div class="tooltip-table tooltip-style">
      <div class="tooltip-row">
          <div class="tooltip-label">${hover.yname}</div>
          <div class="tooltip-value">${bucket}</div>
      </div>
      <div class="tooltip-row">
          <div class="tooltip-label">${hover.zname}</div>
          <div class="tooltip-value">${d.z}</div>
      </div>
  </div>`;
      };
      instance.setContent(data.points.map(pt2txt)[0]);
      instance.show();
    });

    plotDiv.on("plotly_unhover", function (ev) {
      instance.hide();
    });

    // Fallback: hide tooltip when the mouse leaves the plot area.
    plotDiv.addEventListener("mouseout", () => {
      instance.hide();
    });
  }

  onDestroy(() => {
    Plotly.purge(plotDiv);
    unsubConfig?.();
    unsubMetrics?.();
  });
</script>

<div>
  <div bind:this={plotDiv} style="width:100%; height:100%"></div>
</div>
