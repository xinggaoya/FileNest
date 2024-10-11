import axios from 'axios'

class HttpClient {
  constructor(baseURL, headers = {}) {
    this.axiosInstance = axios.create({
      baseURL,
      headers
    })

    // 请求拦截器
    this.axiosInstance.interceptors.request.use(
      (config) => {
        // 在发送请求之前可以做点什么，比如添加token等
        return config
      },
      (error) => {
        // 对请求错误做些什么
        return Promise.reject(error)
      }
    )

    // 响应拦截器
    this.axiosInstance.interceptors.response.use(
      (response) => {
        // 对响应数据做些什么
        return response.data
      },
      (error) => {
        // 对响应错误做些什么
        return Promise.reject(error)
      }
    )
  }

  // 通用的GET方法
  get(url, params = {}, config = {}) {
    return this.axiosInstance.get(url, { params, ...config })
  }

  // 通用的POST方法
  post(url, data = {}, config = {}) {
    return this.axiosInstance.post(url, data, config)
  }

  // 通用的PUT方法
  put(url, data = {}, config = {}) {
    return this.axiosInstance.put(url, data, config)
  }

  // 通用的DELETE方法
  delete(url, config = {}) {
    return this.axiosInstance.delete(url, config)
  }

  // 上传文件的方法
  upload(url, file, data = {}, config = {}) {
    const formData = new FormData()
    formData.append('file', file)

    // 如果有其他数据也一并添加到FormData中
    for (const key in data) {
      if (Object.prototype.hasOwnProperty.call(data, key)) {
        formData.append(key, data[key])
      }
    }

    return this.axiosInstance.post(url, formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
        ...config.headers
      },
      ...config,
      progress: (event) => {
        console.log(event)
      }
    })
  }
}

export const request = new HttpClient('/api', {
  Authorization: 'Bearer your-token-here'
})
