import request from '@/utils/request'

export function login(data) {
  return request({
    url: '/admin/userCenter/login',
    method: 'post',
    data
  })
}

export function getInfo() {
  return request({
    url: '/admin/userCenter/userInfo',
    method: 'get'
    // params: {  }
  })
}

export function logout() {
  return request({
    url: '/admin/userCenter/logout',
    method: 'post'
  })
}
