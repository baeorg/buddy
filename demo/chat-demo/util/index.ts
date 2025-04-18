import axios from "./axios";
import apiController from "./api";

function setBaseUrl(url: string) {
  window.localStorage.setItem("IM_BASE_URL", url);
  axios.defaults.baseURL = url;
}

export const decodeBase64Response = (base64: string) => {
  if (!base64) {
    console.error("Empty base64 string");
    return null;
  }

  try {
    const decoded = atob(base64);

    try {
      return JSON.parse(decoded);
    } catch (jsonError) {
      return decoded;
    }
  } catch (error) {
    console.error("Failed to decode base64:", error);
    return null;
  }
};

export const api = {
  ...apiController,
  setBaseUrl,
};
