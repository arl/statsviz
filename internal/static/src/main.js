import "bootstrap/dist/css/bootstrap.min.css";
import "./style.css";
import { connect } from "./app.js";
import * as theme from "./theme.js";

connect();
theme.updateThemeMode();
