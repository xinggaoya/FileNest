import axios, { type AxiosRequestConfig, type AxiosResponse } from 'axios'

class HttpClient {
  axiosInstance: any

  constructor(baseURL: string, headers = {}) {
    this.axiosInstance = axios.create({
      baseURL,
      headers
    })

    // 请求拦截器
    this.axiosInstance.interceptors.request.use(
      (config: AxiosRequestConfig) => {
        // 在发送请求之前可以做点什么，比如添加token等
        return config
      },
      (error: any) => {
        // 对请求错误做些什么
        return Promise.reject(error)
      }
    )

    // 响应拦截器
    this.axiosInstance.interceptors.response.use(
      (response: AxiosResponse) => {
        // 对响应数据做些什么
        if (response.data.code !== 1000) {
          throw new Error(response.data.message)
        }
        return response.data
      },
      (error: any) => {
        // 对响应错误做些什么
        console.error(error)
        return Promise.reject(error)
      }
    )
  }

  // 通用的GET方法
  get(url: string, params = {}, config: AxiosRequestConfig = {}) {
    return this.axiosInstance.get(url, { params, ...config })
  }

  // 通用的POST方法
  post(url: string, data = {}, config: AxiosRequestConfig = {}) {
    return this.axiosInstance.post(url, data, config)
  }

  // 通用的PUT方法
  put(url: string, data = {}, config: AxiosRequestConfig = {}) {
    return this.axiosInstance.put(url, data, config)
  }

  // 通用的DELETE方法
  delete(url: string, config: AxiosRequestConfig = {}) {
    return this.axiosInstance.delete(url, config)
  }

  // 上传文件的方法
  upload(
    url: string,
    file: File,
    data: any = {},
    config: AxiosRequestConfig = {},
    progress?: Function
  ) {
    const formData = new FormData()
    formData.append('file', file)

    // 如果有其他数据也一并添加到FormData中
    for (const key in data) {
      if (Object.prototype.hasOwnProperty.call(data, key)) {
        formData.append(key, data[key])
      }
    }

    // 计算表单总大小
    return this.axiosInstance.post(url, formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
        ...config.headers
      },
      ...config,
      onUploadProgress: (event: ProgressEvent) => {
        if (progress) {
          progress(event)
        }
      }
    })
  }
}

export const request = new HttpClient('/api', {
  Authorization: 'Bearer your-token-here'
})
