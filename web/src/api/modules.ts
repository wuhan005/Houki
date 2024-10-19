import axios from 'axios';

export interface ModuleBody {
  Request: {
    on: string;
    transmit: string;
    headers: Record<string, string>;
    body: any;
  };
  Response: {
    on: string;
    statusCode: number;
    headers: Record<string, string>;
    body: any;
  };
}

export interface Module {
  id: number;
  body: ModuleBody;
  createdAt: Date;
}

export interface ListModulesParams {
  enabledOnly?: boolean;
  page?: number;
  pageSize?: number;
}

export function listModules(params: ListModulesParams) {
  return axios.get<{ modules: Module[]; total: number }>('/api/modules', {
    params,
  });
}
