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
  name: string;
  body: ModuleBody;
  createdAt: Date;
}

export interface ListModulesParams {
  enabledOnly?: boolean;
  page?: number;
  pageSize?: number;
}

export function listModules(params: ListModulesParams) {
  return axios.get<
    { modules: Module[]; total: number },
    { modules: Module[]; total: number }
  >('/api/modules', {
    params,
  });
}

export function getModule(id: string | number) {
  return axios.get<Module, Module>(`/api/modules/${id}`);
}

export interface CreateModuleData {
  name: string;
  body: ModuleBody;
}

export function createModule(data: CreateModuleData) {
  return axios.post<string, string>('/api/modules', data);
}

export interface UpdateModuleData {
  name: string;
  body: ModuleBody;
}

export function updateModule(id: string | number, data: UpdateModuleData) {
  return axios.put<string, string>(`/api/modules/${id}`, data);
}

export function deleteModule(id: string | number) {
  return axios.delete(`/api/modules/${id}`);
}
