import ky, { type Options } from 'ky'

class HttpService {
  private readonly instance: any

  constructor(baseURL: string) {
    this.instance = ky.extend({
      prefixUrl: baseURL,
      // 可以在这里添加默认的请求头
      headers: {
        'Content-Type': 'application/json'
      },
      timeout: 5000,
      retry: {
        limit: 3
      },
      hooks: {
        beforeRequest: [
          (request) => {
            // 在请求发送前添加请求头
            request.headers.set('Authorization', 'Bearer your_token')
          }
        ],
        afterResponse: [
          (request, options, response) => {
            // 在响应返回后进行一些处理，例如检查响应状态码等
            if (response.status === 401) {
              // 响应状态码为401，表示需要重新登录
              // 可以在这里进行重新登录逻辑，例如跳转到登录页面等
            }
          }
        ]
      }
    })
  }

  private async request<T>(method: string, url: string, options?: Options): Promise<T> {
    try {
      const response = await this.instance(url, { method, ...options })
      return await response.json()
    } catch (error) {
      // 处理错误（可以在这里添加更多的错误处理逻辑）
      console.error('Request failed:', error)
      throw error
    }
  }

  public get<T>(url: string, options?: Options): Promise<T> {
    return this.request<T>('GET', url, options)
  }

  public post<T>(url: string, body: any, options?: Options): Promise<T> {
    return this.request<T>('POST', url, { json: body, ...options })
  }

  public put<T>(url: string, body: any, options?: Options): Promise<T> {
    return this.request<T>('PUT', url, { json: body, ...options })
  }

  public delete<T>(url: string, options?: Options): Promise<T> {
    return this.request<T>('DELETE', url, options)
  }

  // 新增的文件上传方法
  public async uploadFile<T>(url: string, file: File, options?: Options): Promise<T> {
    const formData = new FormData();
    formData.append('file', file);

    try {
      const response = await this.instance(url, {
        method: 'POST',
        body: formData,
        ...options,
      });
      return await response.json();
    } catch (error) {
      console.error('File upload failed:', error);
      throw error;
    }
  }
}

// 使用示例
const httpService = new HttpService('http://localhost:3000/api')

// 导出实例
export default httpService
