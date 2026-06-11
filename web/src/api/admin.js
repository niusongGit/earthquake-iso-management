import request from './request'

// 管理员登录
export function login(data) {
  return request.post('/admin/login', data)
}

// 获取文档列表（后台）
export function getDocumentList(params) {
  return request.get('/admin/documents', { params })
}

// 获取文档详情（后台）
export function getDocumentDetail(id) {
  return request.get(`/admin/documents/${id}`)
}

// 创建文档
export function createDocument(data) {
  return request.post('/admin/documents', data, {
    headers: { 'Content-Type': 'multipart/form-data' }
  })
}

// 更新文档
export function updateDocument(id, data) {
  return request.put(`/admin/documents/${id}`, data, {
    headers: { 'Content-Type': 'multipart/form-data' }
  })
}

// 删除文档
export function deleteDocument(id) {
  return request.delete(`/admin/documents/${id}`)
}
