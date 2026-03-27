import { useAuthStore } from "@/stores/auth";
import { useDasDevices } from "@/stores/das-devices";
import axios from "axios";
// import { ElMessageBox, ElMessage } from 'element-plus'
// import { getToken } from '@/utils/auth'
// import { useAuthStore } from '../stores/auth'

// create an axios instance
const service = axios.create({
  baseURL: "//", // url = base url + request url
  // withCredentials: true, // send cookies when cross-domain requests
  timeout: 10000, // request timeout
});

// request interceptor
service.interceptors.request.use(
  (config) => {
    // do something before request is sent
    // const token = getToken();
    // if (token) {
    //   config.headers['Authorization'] = "Bearer " + token;
    // }
    return config;
  },
  (error) => {
    // do something with request error
    return Promise.reject({ message: error.message, error: error });
  }
);

// response interceptor
service.interceptors.response.use(
  /**
   * If you want to get http information such as headers or status
   * Please return  response => response
   */

  /**
   * Determine the request status by custom code
   * Here is just an example
   * You can also judge the status by HTTP Status Code
   */
  (response) => {
    let { data, status, config } = response;

    const dasDevices = useDasDevices();
    const url = config?.url;
    const parts = url?.match(/(\/das\/devices|\/proxy)\/(\d+|local)\/.*/);
    const sub = parts?.length === 3 ? parts[2] : "";
    if (sub) {
      const dev = dasDevices.getDevice(sub);
      dev?.setAvailable(true);
    }
    return data;
  },
  (error) => {
    let { message, response } = error;
    let { data, status, config } = response;

    if (status === 501) {
      useAuthStore().logout(data.message);
    } else if (status === 401) {
      useAuthStore().logout(data.message);
    } else {
      const dasDevices = useDasDevices();
      const url = config?.url;
      const parts = url?.match(/(\/das\/devices|\/proxy)\/(\d+|local)\/.*/);
      const sub = parts?.length === 3 ? parts[2] : "";
      if (sub) {
        const dev = dasDevices.getDevice(sub);
        if (response?.status === 503 || response?.status === 502) {
          dev?.setAvailable(false);
        } else if (message.startsWith("timeout")) {
          dev?.setAvailable(false);
        }
      }
    }
    return Promise.reject({ status: status, message: data.message || data, error: error });
  }
);

export default service;

export async function openApiRequest(url, options) {
  return service({ url, ...options });
}

export function exportFileByALink(url, fileName) {
  const link = document.createElement("a");
  const body = document.querySelector("body");

  link.href = url;

  link.download = fileName;

  link.style.display = "none";
  body.appendChild(link);

  link.click();
  body.removeChild(link);
}
