import axios, { RawAxiosRequestConfig } from 'axios';

export const baseURL =
  process.env.NODE_ENV === 'development' ||
  window.location.origin === 'http://localhost:8080'
    ? 'http://localhost:8080'
    : 'https://electron-social-chat-backend.herokuapp.com';

const api = axios.create({
  baseURL,
});

export async function makeRequest(
  url: string,
  options?: RawAxiosRequestConfig
) {
  return api(url, options)
    .then((res: any) => res.data)
    .catch((e: any) =>
      Promise.reject((e.response?.data.msg ?? 'Error').replace('Error: ', ''))
    );
}
