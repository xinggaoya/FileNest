import httpService from '@/api/http'

// 获取文件列表
export const getFileList = (params: any) => {
  return httpService.get('file/list', {
    searchParams: params
  })
}

// 上传文件
export const uploadFile = (file: any, params: any) => {
  return httpService.uploadFile('file/upload', file, {
    json: params
  })
}
