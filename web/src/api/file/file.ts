import { request } from '@/api/http'

// 获取文件列表
export const getFileList = (params: any) => {
  return request.get('file/list', params)
}

// 上传文件
export const uploadFile = (file: any, params: any, progress: Function) => {
  return request.upload('file/upload', file, params, {}, progress)
}

// 创建文件夹
export const createFolder = (path: string) => {
  return request.post('file/create', {
    path
  })
}

// 删除文件
export const deleteFile = (path: string) => {
  return request.delete('file/delete', {
    params: {
      path
    }
  })
}
