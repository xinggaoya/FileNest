import {request} from '@/api/http'

// 获取文件列表
export const getFileList = (params: any) => {
  return request.get('file/list', params)
}

// 上传文件
export const uploadFile = (file: any, params: any) => {
  return request.upload('file/upload', file, params)
}
