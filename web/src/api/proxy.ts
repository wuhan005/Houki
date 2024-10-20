import axios from 'axios';

export interface ProxyStatusResp {
  forward: {
    status: string;
    address: number;
  };
  reverse: {
    status: string;
    address: number;
  };
}

export function proxyStatus() {
  return axios.get<ProxyStatusResp, ProxyStatusResp>('/api/proxy/status');
}

export function startForwardProxy(address: string) {
  return axios.post<string, string>('/api/proxy/forward/start', {
    address,
  });
}

export function shutdownForwardProxy() {
  return axios.post<string, string>('/api/proxy/forward/shutdown');
}

export function startReverseProxy(address: string) {
  return axios.post<string, string>('/api/proxy/reverse/start', {
    address,
  });
}

export function shutdownReverseProxy() {
  return axios.post<string, string>('/api/proxy/reverse/shutdown');
}
