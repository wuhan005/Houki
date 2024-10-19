import axios from 'axios';
import type { AxiosRequestConfig, AxiosResponse } from 'axios';
import { Message } from '@arco-design/web-vue';

export interface HttpResponse<T = unknown> {
  msg: string;
  data: T;
}

if (import.meta.env.VITE_API_BASE_URL) {
  axios.defaults.baseURL = import.meta.env.VITE_API_BASE_URL;
}

axios.interceptors.response.use(
  (response: AxiosResponse<HttpResponse>) => {
    return response.data.data;
  },
  (error) => {
    const response = error.response.data;
    Message.error({
      content: response.msg || 'Failed to request',
      duration: 5 * 1000,
    });
    return Promise.reject(error);
  }
);
