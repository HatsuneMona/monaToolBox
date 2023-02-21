import request from '@/utils/request'

export function login(data) {
  return request({
    url: '/usercenter/login',
    method: 'post',
    data
  })
}

export function getInfo() {
  return request({
    url: '/usercenter/userInfo',
    method: 'get',
    // params: {  }
  })
}

export function logout() {
  return request({
    url: '/usercenter/logout',
    method: 'post'
  })
}
