import request from '@/utils/request'

export const getUserInfo = () =>
  request({
    url: '/users/me',
    method: 'get'
  })
