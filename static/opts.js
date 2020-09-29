// opts exports the plots options
var opts = (function () {


    function tooltipsPlugin(opts) {
        function init(u, opts, data) {
            let plot = u.root.querySelector(".u-over");

            let ttc = u.cursortt = document.createElement("div");
            ttc.className = "tooltip";
            ttc.textContent = "(x,y)";
            ttc.style.pointerEvents = "none";
            ttc.style.position = "absolute";
            ttc.style.background = "rgba(0,0,255,0.1)";
            plot.appendChild(ttc);

            u.seriestt = opts.series.map((s, i) => {
                if (i == 0) return;

                let tt = document.createElement("div");
                tt.className = "tooltip";
                tt.textContent = "Tooltip!";
                tt.style.pointerEvents = "none";
                tt.style.position = "absolute";
                tt.style.background = "rgba(0,0,0,0.1)";
                tt.style.color = s.color;
                tt.style.display = s.show ? null : "none";
                plot.appendChild(tt);
                return tt;
            });

            function hideTips() {
                ttc.style.display = "none";
                u.seriestt.forEach((tt, i) => {
                    if (i == 0) return;

                    tt.style.display = "none";
                });
            }

            function showTips() {
                ttc.style.display = null;
                u.seriestt.forEach((tt, i) => {
                    if (i == 0) return;

                    let s = u.series[i];
                    tt.style.display = s.show ? null : "none";
                });
            }

            plot.addEventListener("mouseleave", () => {
                if (!u.cursor.locked) {
                    //	u.setCursor({left: -10, top: -10});
                    hideTips();
                }
            });

            plot.addEventListener("mouseenter", () => {
                showTips();
            });

            hideTips();
        }

        function setCursor(u) {
            const { left, top, idx } = u.cursor;

            // this is here to handle if initial cursor position is set
            // not great (can be optimized by doing more enter/leave state transition tracking)
            //	if (left > 0)
            //		u.cursortt.style.display = null;

            u.cursortt.style.left = left + "px";
            u.cursortt.style.top = top + "px";
            u.cursortt.textContent = "(" + u.posToVal(left, "x").toFixed(2) + ", " + u.posToVal(top, "y").toFixed(2) + ")";

            // can optimize further by not applying styles if idx did not change
            u.seriestt.forEach((tt, i) => {
                if (i == 0) return;

                let s = u.series[i];

                if (s.show) {
                    // this is here to handle if initial cursor position is set
                    // not great (can be optimized by doing more enter/leave state transition tracking)
                    //	if (left > 0)
                    //		tt.style.display = null;

                    let xVal = u.data[0][idx];
                    let yVal = u.data[i][idx];

                    tt.textContent = "(" + xVal + ", " + yVal + ")";

                    tt.style.left = Math.round(u.valToPos(xVal, 'x')) + "px";
                    tt.style.top = Math.round(u.valToPos(yVal, s.scale)) + "px";
                }
            });
        }

        return {
            hooks: {
                init,
                setCursor,
                setScale: [
                    (u, key) => {
                        console.log('setScale', key);
                    }
                ],
                setSeries: [
                    (u, idx) => {
                        console.log('setSeries', idx);
                    }
                ],
            },
        };
    }


    // Removes the legend at the plot bottom and shows values on the right.
    function legendAsTooltipPlugin({ className, style = { backgroundColor: "rgba(255, 249, 196, 0.92)", color: "black" } } = {}) {
        let legendEl;

        function init(u, opts) {
            legendEl = u.root.querySelector(".u-legend");

            legendEl.classList.remove("u-inline");
            className && legendEl.classList.add(className);

            uPlot.assign(legendEl.style, {
                textAlign: "left",
                pointerEvents: "none",
                display: "none",
                position: "absolute",
                left: 0,
                top: 0,
                zIndex: 100,
                boxShadow: "2px 2px 10px rgba(0,0,0,0.5)",
                ...style
            });

            // hide series color markers
            const idents = legendEl.querySelectorAll(".u-marker");

            for (let i = 0; i < idents.length; i++)
                idents[i].style.display = "none";

            const overEl = u.root.querySelector(".u-over");
            overEl.style.overflow = "visible";

            // move legend into plot bounds
            overEl.appendChild(legendEl);

            // show/hide tooltip on enter/exit
            overEl.addEventListener("mouseenter", () => { legendEl.style.display = null; });
            overEl.addEventListener("mouseleave", () => { legendEl.style.display = "none"; });

            // let tooltip exit plot
            overEl.style.overflow = "visible";
        }

        function update(u) {
            const { left, top } = u.cursor;
            legendEl.style.transform = "translate(" + left + "px, " + top + "px)";
        }

        return {
            hooks: {
                init: init,
                setCursor: update,
            }
        };
    }

    function humanBytes(bytes, si = false, dp = 1) {
        const thresh = si ? 1000 : 1024;

        if (Math.abs(bytes) < thresh) {
            return bytes + ' B';
        }

        const units = si
            ? ['kB', 'MB', 'GB', 'TB', 'PB', 'EB', 'ZB', 'YB']
            : ['KiB', 'MiB', 'GiB', 'TiB', 'PiB', 'EiB', 'ZiB', 'YiB'];
        let u = -1;
        const r = 10 ** dp;

        do {
            bytes /= thresh;
            ++u;
        } while (Math.round(Math.abs(bytes) * r) / r >= thresh && u < units.length - 1);

        return bytes.toFixed(dp) + ' ' + units[u];
    }

    function valueFmt(u, v) {
        return humanBytes(v, true);
    }

    const cursorOpts = {
        lock: true,
        focus: {
            prox: 16,
        },
        points: {
            show: false,
        },
        sync: {
            key: "ts",
            setSeries: true,
        },
    };

    function gcLinesPlugin() {
        return {
            hooks: {
                draw: u => {
                    const { ctx } = u;
                    const gcs = stats.lastGCs;
                    const mints = u.data[0][0];
                    const maxts = u.data[0][u.data[0].length - 1];
                    const y0 = u.valToPos(u.scales['b'].min, 'b', true);
                    const y1 = u.valToPos(u.scales['b'].max, 'b', true);

                    for (let i = 0, n = gcs.length; i < n; i++) {
                        const ts = gcs[i];
                        if (ts < mints || ts > maxts) {
                            continue;
                        }

                        const x = u.valToPos(ts, 'x', true);
                        ctx.beginPath();
                        ctx.moveTo(x, y0);
                        ctx.lineTo(x, y1);
                        ctx.setLineDash([5, 10]);
                        ctx.lineWidth = 1;
                        ctx.strokeStyle = 'grey';
                        ctx.stroke();
                    }
                }
            }
        };
    }

    function getSize() {
        return {
            width: 950,
            height: 400,
        }
    }

    function humanBytesValues(u, sidx, idx) {
        let val = " -- ";
        if (u.data != undefined) {
            val = humanBytes(u.data[sidx][idx], true);
        }
        return {
            "bytes": val,
        };
    }

    const opts1 = {
        title: "Heap",
        ...getSize(),
        cursor: cursorOpts,
        plugins: [
            gcLinesPlugin(),
            // legendAsTooltipPlugin(),
            // tooltipsPlugin(),
        ],
        series: [
            {},
            {
                label: "HeapAlloc",
                scale: "b",
                value: valueFmt,
                values: humanBytesValues,
                stroke: "red",
                points: {
                    show: true,
                    size: 3,
                    fill: "red",
                },
            },
            {
                label: "HeapSys",
                scale: "b",
                value: valueFmt,
                values: humanBytesValues,
                stroke: "blue",
                points: {
                    show: true,
                    size: 3,
                    fill: "blue",
                },
            },
            {
                label: "HeapIdle",
                scale: "b",
                value: valueFmt,
                values: humanBytesValues,
                stroke: "green",
                points: {
                    show: true,
                    size: 3,
                    fill: "green",
                },
            },
            {
                label: "HeapInuse",
                scale: "b",
                value: valueFmt,
                values: humanBytesValues,
                stroke: "orange",
                points: {
                    show: true,
                    size: 3,
                    fill: "orange",
                },
            },
        ],
        axes: [
            {
                values: (u, vals, space) => vals.map(v => formatAxisTimestamp(v)),
                rotate: 50,
                size: 80,
            },
            {
                scale: 'b',
                values: (u, vals, space) => vals.map(v => humanBytes(v, true, 0)),
                size: 90,
            },
        ],
    };

    // formatAxisTimestamp formats a given Unix epoch timestamp for printing
    // along an axis. It prints:
    //  - `hh:mm:ss` if ss is a multiple of 5s,
    //  - only `ss` otherwise
    function formatAxisTimestamp(ts) {
        let d = new Date(ts * 1000);
        let s = d.getSeconds()
        let ss = s.toString().padStart(2, '0')

        if (s % 5 != 0) {
            return ss
        }

        let hh = d.getHours().toString().padStart(2, '0')
        let mm = d.getMinutes().toString().padStart(2, '0')
        return hh + ':' + mm + ':' + ss
    }

    function heatmapPlugin() {
        // let global min/max
        function fillStyle(count, maxCount) {
            const norm = count / maxCount;

            // salmon
            // const r = 254 - (24 * norm);
            // const g = 230 - (145 * norm);
            // const b = 206 - (193 * norm);

            // purple
            const r = 239 - (122 * norm);
            const g = 237 - (120 * norm);
            const b = 245 - (68 * norm);
            return `rgba(${r}, ${g}, ${b}, 1)`;
        }

        return {
            hooks: {
                draw: u => {
                    const { ctx, data } = u;

                    let yData = data[3];
                    let yQtys = data[4];

                    let iMin = u.scales.x.min;
                    let iMax = u.scales.x.max;

                    const rectw = u.bbox.width / (iMax - iMin);
                    const recth = u.bbox.height / (u.scales.y.max - u.scales.y.min);

                    let maxCount = -Infinity;

                    yQtys.forEach(qtys => {
                        maxCount = Math.max(maxCount, Math.max.apply(null, qtys));
                    });

                    yData.forEach((yVals, xi) => {
                        let xPos = u.valToPos(data[0][xi], 'x', true);
                        xPos = xPos - rectw;

                        yVals.forEach((yVal, yi) => {
                            const count = yQtys[xi][yi];
                            if (count == 0) {
                                // Skip empty size classes
                                return;
                            }
                            const yPos = Math.round(u.valToPos(yVal, 'y', true));
                            ctx.fillStyle = fillStyle(count, maxCount);
                            ctx.fillRect(xPos, yPos, rectw, recth);
                            ctx.strokeRect(xPos, yPos, rectw, recth);
                        });
                    });
                }
            }
        };
    }

    // column-highlights the hovered x index
    function columnHighlightPlugin({ className, style = { backgroundColor: "rgba(51, 51, 51, 0.1)" } } = {}) {
        let over, hlColumn, hlSizeClass, currIdx, currTop;

        function init(u) {
            over = u.root.querySelector(".u-over");

            hlColumn = document.createElement("div");
            hlSizeClass = document.createElement("div");

            className && hlColumn.classList.add(className);
            uPlot.assign(hlColumn.style, {
                pointerEvents: "none",
                display: "none",
                position: "absolute",
                left: 0,
                top: 0,
                height: "100%",
                ...style
            });

            className && hlSizeClass.classList.add(className);
            uPlot.assign(hlSizeClass.style, {
                pointerEvents: "none",
                display: "none",
                position: "absolute",
                left: 0,
                top: 0,
                // TODO: doesn't feel right... we do have 'stats.classSizes.length'
                // size classes, however the full height allows to place more than 
                // that. We probably don't see any real difference because we're 
                // talking about a discrepancy of some fraction of a pixel. However 
                // if the devicePixelRatio changes and/or we zoom in, I suspect the
                // height of the highlighting element class won't exactly match that
                // of the highlighted size class.
                height: 100 / stats.classSizes.length + "%",
                backgroundColor: "rgba(51, 51, 51, 0.3)",
            });


            over.appendChild(hlColumn);
            over.appendChild(hlSizeClass);

            // show/hide highlight on enter/exit
            over.addEventListener("mouseenter", () => {
                hlColumn.style.display = null;
                hlSizeClass.style.display = null;
            });
            over.addEventListener("mouseleave", () => {
                hlColumn.style.display = "none";
                hlSizeClass.style.display = "none";
            });
        }

        function update(u) {
            if (currIdx !== u.cursor.idx || currTop !== u.cursor.top) {
                currIdx = u.cursor.idx;
                currTop = u.cursor.top;
                const dx = u.scales.x.max - u.scales.x.min;
                const width = (u.bbox.width / dx) / devicePixelRatio;
                const xVal = u.data[0][currIdx];
                const left = u.valToPos(xVal, "x") - width;

                hlColumn.style.transform = "translateX(" + Math.round(left) + "px)";
                hlColumn.style.width = Math.round(width) + "px";

                // size class rectangle height
                const recth = u.bbox.height / (u.scales.y.max - u.scales.y.min);

                const yVal = u.posToVal(u.cursor.top, "y");
                const scIdx = Math.floor((u.scales.y.max - u.scales.y.min) - yVal).toFixed(0);
                const tx = scIdx * recth;
                hlSizeClass.style.transform = "translate(" + Math.round(left) + "px, " + Math.round(tx) + "px)";
                hlSizeClass.style.width = Math.round(width) + "px";
            }
        }

        return {
            opts: (u, opts) => {
                uPlot.assign(opts, {
                    cursor: {
                        x: false,
                        y: false,
                    }
                });
            },
            hooks: {
                init: init,
                setCursor: update,
            }
        };
    }

    const opts2 = {
        title: "Size classes Heatmap",
        ...getSize(),
        cursor: cursorOpts,
        legend: { show: false },
        plugins: [
            heatmapPlugin(),
            columnHighlightPlugin(),
        ],
        gutters: {
            y: 10,
        },
        series: [
            {
                scale: 'x',
            },
            {
                paths: () => null,
                points: { show: false },
                scale: 'y',
                label: "size class",
                /*
                values: (u, sidx, idx) => {
                    let val = " -- ";
                    if (u.data != undefined) {
                        const { top } = u.cursor;
                        if (top !== null) {
                            // console.log("left: " + left + " top: " + top + " idx: " + idx);
                            // console.log(`bbox: left:${u.bbox.left} top:${u.bbox.top} width:${u.bbox.width} height:${u.bbox.height}`);
                            // let idx = Math.floor(u.posToIdx(left, 'x'));
                            let j = Math.floor(u.posToVal(top, 'y'));
                            if (j >= 0 && j < stats.classSizes.length) {
                                val = stats.classSizes[j];
                            }
                        }
                    }
                    return {
                        "value": val + ' - ' + '??',
                    };
                }
                */
            },
            {
                paths: () => null,
                points: { show: false },
                scale: 'y',
                label: "count",
                /*
                values: (u, sidx, idx) => {
                    let val = " -- ";
                    if (u.data != undefined) {
                        const { top, idx } = u.cursor;
                        if (idx !== null && top !== null) {
                            let j = Math.floor(u.posToVal(top, 'y'));
                            if (j >= 0 && j < u.data[4][idx].length) {
                                val = u.data[4][idx][j];
                            }
                        }
                    }
                    return {
                        "value": val,
                    };
                }
                */
            },
        ],
        axes: [
            {
                scale: 'x',
                values: (u, vals, space) => vals.map(v => formatAxisTimestamp(v)),
                rotate: 50,
                size: 80,
            },
            {
                scale: 'y',
                values: (u, vals, space) => vals.map(function (i) {
                    if (i > stats.classSizes.length - 1) {
                        return '';
                    }
                    return humanBytes(stats.classSizes[i], true, 0);
                }),
                size: 90,
            },
        ],
    };

    return {
        heap: opts1,
        bySizes: opts2,
    }
}());
