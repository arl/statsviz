import "bootstrap/dist/js/bootstrap.bundle.min.js";
import "bootstrap/dist/css/bootstrap.min.css";

import "./style.css";
import { startApp } from "./app.js";

document.querySelector(
  "#app"
).innerHTML = `<nav id="navbar" class="navbar navbar-expand-lg bg-body-tertiary">
  <div class="container-fluid">
    <a class="navbar-brand">Statsviz</a>
    <button class="navbar-toggler" type="button" data-bs-toggle="collapse" data-bs-target="#navbarSupportedContent" aria-controls="navbarSupportedContent" aria-expanded="false" aria-label="Toggle navigation">
      <span class="navbar-toggler-icon"></span>
    </button>
    <div class="collapse navbar-collapse" id="navbarSupportedContent">
      <ul class="navbar-nav mb-2 mb-lg-0 ms-auto">
        <li class="nav-item">
          <a class="bi bi-github nav-link" href="https://github.com/arl/statsviz"> GitHub</a>
        </li>
        <li class="nav-item">
          <a class="nav-link" href="https://pkg.go.dev/runtime/metrics#hdr-Supported_metrics">runtime/metrics</a>
        </li>
        <li class="nav-item">
          <select class="form-select" id="select-timerange" style="width: 120px;">
            <option value="60">1 minute</option>
            <option value="300">5 minutes</option>
            <option value="600">10 minutes</option>
          </select>
        </li>
        <li class="nav-item d-flex align-items-center" style="margin-left: 25px;">
          <div class="form-check form-switch">
            <input class="form-check-input" type="checkbox" id="dark_mode_switch" name="darkmode" value="yes">
            <label class="form-check-label" for="darkmode">Dark Mode</label>
          </div>
        </li>
        <li class="nav-item d-flex align-items-center" style="margin-left: 25px;">
          <div class="form-check form-switch">
            <input class="form-check-input" type="checkbox" id="show_gc_switch" name="showgc" value="yes">
            <label class="form-check-label" for="showgc">Show GC</label>
          </div>
        </li>
      </ul>
    </div>
  </div>
</nav>

<div id="plots" class="plots-wrapper"></div>
`;
startApp();
