import axios from "./axios";
import apiController from "./api";

function setBaseUrl(url: string) {
  window.localStorage.setItem("IM_BASE_URL", url);
  axios.defaults.baseURL = url;
}

export const api = {
  ...apiController,
  setBaseUrl,
};
