import request from './request'

// 获取文档列表
export function getDocumentList(params) {
  return request.get('/front/documents', { params })
}

// 获取文档详情
export function getDocumentDetail(id) {
  return request.get(`/front/documents/${id}`)
}

// 获取附件下载地址
export function getDownloadUrl(id) {
  return `/api/front/documents/${id}/download`
}
