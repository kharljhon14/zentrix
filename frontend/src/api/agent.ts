import axios, { AxiosError, type AxiosResponse } from 'axios';

import type { Company } from '@/features/company/types/company';

const BASE_URL = 'http://localhost:4000';

const responseBody = <T>(res: AxiosResponse<T>) => res.data;

axios.interceptors.response.use(
  (response) => response,
  (error: AxiosError) => Promise.reject(error)
);

const requests = {
  get: <TResponse>(url: string) => axios.get(url).then(responseBody<TResponse>),
  post: <TResponse, TBody>(url: string, body: TBody) =>
    axios.post(url, body).then(responseBody<TResponse>),
  patch: <TResponse, TBody>(url: string, body: TBody) =>
    axios.patch(url, body).then(responseBody<TResponse>),
  delete: <TResponse>(url: string) => axios.delete(url).then(responseBody<TResponse>)
};

interface DataResponse<T> {
  data: T;
  metadata: Metadata;
}

interface Metadata {
  current_page: number;
  page_size: number;
  first_page: number;
  last_page: number;
  total_records: number;
}

const companies = {
  getById: (id: string) =>
    requests.get<Omit<DataResponse<Company>, 'metadata'>>(`${BASE_URL}/companies/${id}`),
  getCompanies: (page?: number, pageSize?: number, sort?: string) =>
    requests.get<DataResponse<Company[]>>(`${BASE_URL}/companies`),
  create: () => requests.post(`${BASE_URL}/companies`, {}),
  update: () => requests.patch(`${BASE_URL}/companies`, {}),
  delete: (id: string) => requests.delete(`${BASE_URL}/companies/${id}`)
};

export default {
  companies
};
