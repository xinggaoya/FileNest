import { get, post, put, del, download } from '@/utils/request'

// 分享相关接口类型定义
export interface CreateShareRequest {
  filePath: string
  password?: string
  expireHours?: number
  maxDownload?: number
  description?: string
}

export interface ShareInfo {
  shareCode: string
  fileName: string
  fileSize: number
  isDir: boolean
  hasPassword: boolean
  expireTime?: string
  maxDownload: number
  downloaded: number
  description: string
  createTime: string
}

export interface ShareItem {
  id: number
  shareCode: string
  filePath: string
  fileName: string
  fileSize: number
  isDir: boolean
  password?: string
  expireTime?: string
  maxDownload: number
  downloaded: number
  status: string
  description: string
  creatorIP: string
  createTime: string
  updateTime: string
}

export interface MySharesResponse {
  shares: ShareItem[]
  total: number
  page: number
  pageSize: number
}

// 创建分享
export const createShare = (data: CreateShareRequest) => {
  return post<ShareItem>('/share/create', data)
}

// 获取分享信息
export const getShareInfo = (shareCode: string) => {
  return get<ShareInfo>(`/share/info/${shareCode}`)
}

// 验证分享
export const validateShare = (shareCode: string, password?: string) => {
  const formData = new FormData()
  if (password) {
    formData.append('password', password)
  }
  return post(`/share/validate/${shareCode}`, formData)
}

// 下载分享文件
export const downloadSharedFile = (shareCode: string, password?: string) => {
  const params = password ? { password } : {}
  return download(`/share/download/${shareCode}`, params)
}

// 获取我的分享列表
export const getMyShares = (page: number = 1, pageSize: number = 10) => {
  return get<MySharesResponse>('/share/my', { page, pageSize })
}

// 删除分享
export const deleteShare = (shareCode: string) => {
  return del(`/share/${shareCode}`)
}

// 禁用分享
export const disableShare = (shareCode: string) => {
  return put(`/share/disable/${shareCode}`)
}

// 获取分享统计
export const getShareStats = (shareCode: string) => {
  return get<ShareItem>(`/share/stats/${shareCode}`)
}
