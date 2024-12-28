import axios, { type AxiosRequestConfig, type AxiosResponse } from 'axios'
import { createDiscreteApi } from 'naive-ui'

const { message, loadingBar } = createDiscreteApi(['message', 'loadingBar'])

// 创建 axios 实例
const request = axios.create({
    baseURL: '/api',
    timeout: 30000
})

// 请求拦截器
request.interceptors.request.use(
    (config) => {
        loadingBar.start()
        return config
    },
    (error) => {
        console.error('请求错误:', error)
        loadingBar.error()
        return Promise.reject(error)
    }
)

// 响应拦截器
request.interceptors.response.use(
    (response) => {
        loadingBar.finish()
        const { code, message: msg, data } = response.data

        // 这里可以根据后端的响应结构定制
        if (code === 1000 || code === 1005) {
            return response
        }

        // 特定错误码处理
        switch (code) {
            case 401:
                message.error('请先登录')
                break
            case 403:
                message.error('没有权限')
                break
            case 404:
                message.error('资源不存在')
                break
            case 500:
                message.error('服务器错误')
                break
            default:
                message.error(msg || '操作失败')
        }

        return Promise.reject(new Error(msg || '操作失败'))
    },
    (error) => {
        loadingBar.error()
        if (error.response) {
            switch (error.response.status) {
                case 401:
                    message.error('请先登录')
                    break
                case 403:
                    message.error('没有权限')
                    break
                case 404:
                    message.error('资源不存在')
                    break
                case 500:
                    message.error('服务器错误')
                    break
                default:
                    message.error(error.response.data?.message || '网络错误')
            }
        } else if (error.request) {
            message.error('网络连接失败')
        } else {
            message.error('请求配置错误')
        }
        return Promise.reject(error)
    }
)

export interface Response<T = any> {
    code: number
    message: string
    data: T
}

export interface RequestOptions extends AxiosRequestConfig {
    showError?: boolean
    showSuccess?: boolean
}

const defaultOptions: RequestOptions = {
    showError: true,
    showSuccess: false
}

export async function get<T = any>(
    url: string,
    params?: any,
    options?: RequestOptions
): Promise<Response<T>> {
    const finalOptions = { ...defaultOptions, ...options }
    try {
        const response = await request.get<Response<T>>(url, { params, ...finalOptions })
        return response.data
    } catch (error) {
        if (finalOptions.showError) {
            throw error
        }
        return Promise.reject(error)
    }
}

export async function post<T = any>(
    url: string,
    data?: any,
    options?: RequestOptions
): Promise<Response<T>> {
    const finalOptions = { ...defaultOptions, ...options }
    try {
        const response = await request.post<Response<T>>(url, data, finalOptions)
        if (finalOptions.showSuccess) {
            message.success('操作成功')
        }
        return response.data
    } catch (error) {
        if (finalOptions.showError) {
            throw error
        }
        return Promise.reject(error)
    }
}

export async function put<T = any>(
    url: string,
    data?: any,
    options?: RequestOptions
): Promise<Response<T>> {
    const finalOptions = { ...defaultOptions, ...options }
    try {
        const response = await request.put<Response<T>>(url, data, finalOptions)
        if (finalOptions.showSuccess) {
            message.success('操作成功')
        }
        return response.data
    } catch (error) {
        if (finalOptions.showError) {
            throw error
        }
        return Promise.reject(error)
    }
}

export async function del<T = any>(
    url: string,
    params?: any,
    options?: RequestOptions
): Promise<Response<T>> {
    const finalOptions = { ...defaultOptions, ...options }
    try {
        const response = await request.delete<Response<T>>(url, { params, ...finalOptions })
        if (finalOptions.showSuccess) {
            message.success('操作成功')
        }
        return response.data
    } catch (error) {
        if (finalOptions.showError) {
            throw error
        }
        return Promise.reject(error)
    }
}

// 下载文件专用方法
export const download = async (url: string, params: Record<string, any>) => {
    try {
        // 创建一个新的 axios 实例，专门用于下载
        const downloadInstance = axios.create({
            baseURL: '/api',
            timeout: 30000,
            responseType: 'blob'
        })

        const response = await downloadInstance.get(url, {
            params,
            headers: {
                'Content-Type': 'application/json'
            }
        })

        if (response.data instanceof Blob) {
            return response.data
        } else {
            throw new Error('下载失败')
        }
    } catch (error: any) {
        // 如果响应是 Blob 类型但包含错误信息，需要特殊处理
        if (error.response?.data instanceof Blob) {
            const text = await error.response.data.text()
            try {
                const errorData = JSON.parse(text)
                throw new Error(errorData.message || '下载失败')
            } catch {
                throw new Error('下载失败')
            }
        }
        throw error
    }
}

export default request 