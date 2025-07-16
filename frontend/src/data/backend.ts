import axios from "axios";
import type { AxiosInstance } from "axios"

class BackendAxiosCaller {
  private axiosInstance: AxiosInstance;
  public get;
  public post;
  public put;
  public delete;
  constructor() {
    this.axiosInstance = axios.create({
      allowAbsoluteUrls: false,
      baseURL: import.meta.env.VITE_APP_BACKEND_URL,
      headers: { "Content-Type": "application/json" }
    })

    this.axiosInstance.interceptors.request.use(
      (config) => {

        const accessToken = this.getAccessToken();
        if (accessToken) {
          config.headers.Authorization = `Bearer ${accessToken}`;
        }
        return config
      },
      (error) => {
        if (error instanceof Error) {
          console.error(error.message)
        }
        Promise.reject(error)
      }
    )

    this.get = this.axiosInstance.get.bind(this.axiosInstance);
    this.post = this.axiosInstance.post.bind(this.axiosInstance);
    this.put = this.axiosInstance.put.bind(this.axiosInstance);
    this.delete = this.axiosInstance.delete.bind(this.axiosInstance);
  }

  private getAccessToken() {
    return localStorage.getItem("accessToken")
  }
}

export const backendClient = new BackendAxiosCaller()
