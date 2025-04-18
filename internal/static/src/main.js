import "bootstrap/dist/css/bootstrap.min.css";
import "./style.css";
import { startApp } from "./app.js";

document.querySelector("#app").innerHTML = `
    <nav id="navbar" class="navbar navbar-expand-lg navbar-light bg-light">
        <div class="container-xxl">
            <a class="navbar-brand" href="#">Statsviz</a>
            <button class="navbar-toggler" type="button">
                <span class="navbar-toggler-icon"></span>
            </button>
            <div class="collapse navbar-collapse">
                <ul class="navbar-nav me-auto mb-2 mb-lg-0">
                    <li class="nav-item">
                        <a class="nav-link" href="https://github.com/arl/statsviz">Github</a>
                    </li>
                    <li class="nav-item">
                        <a class="nav-link" href="https://pkg.go.dev/runtime/metrics#hdr-Supported_metrics">runtime/metrics</a>
                    </li>
                </ul>
            </div>
            <span title="show/hide GC events">
                <input id="show_gc" type="checkbox" data-toggle="toggle" data-onstyle="default" data-on="<span style='color:black'><i class='fa-solid fa-broom'></i></span>" data-off="<span style='color:lightgrey'><i class='fa-solid fa-broom'></i></span>"
            data-size="mini" data-bs-toggle="tooltip">
            </span>
            <select data-bs-toggle="tooltip" title="time range" id="select_timerange">
                <option selected value="60">1 minute</option>
                <option value="300">5 minute</option>
                <option value="600">10 minutes</option>
              </select>
            <span class="action" title="play/pause">
                <input id="play_pause" type="checkbox" data-toggle="toggle" data-onstyle="default" data-on="<i class='fa fa-play'></i>" data-off="<i class='fa fa-pause'></i>" data-size="mini" data-bs-toggle="tooltip">
            </span>
            <span class="action" title="color theme">
                <input id="color_theme_sw" type="checkbox" data-toggle="toggle" data-onstyle="default" data-on="<i class='fa fa-circle-half-stroke'></i>" data-off="<i class='fa fa-circle-half-stroke'></i>" data-size="mini" data-bs-toggle="tooltip">
            </span>
        </div>
    </nav>

    <div id="plots" class="plots-wrapper"></div>`;

startApp();
