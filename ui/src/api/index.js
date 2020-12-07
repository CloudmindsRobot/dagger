import axios from 'axios'
import router from '@/router'
import store from '@/store'

axios.defaults.baseURL = `${window.location.origin}/api/v1`
axios.defaults.timeout = 1000 * 60

axios.interceptors.request.use(function(config) {
  if (store.state.csrftoken) {
    config.headers['X-CSRFToken'] = store.state.csrftoken
  }
  if (store.state.jwt) {
    config.headers.Authorization = `JWT ${store.state.jwt}`
  }
  if (config.method.toLowerCase() === 'post') {
    config.headers['Content-type'] = 'application/json;charset=utf-8'
  }
  return config
})
axios.interceptors.response.use(
  function(response) {
    return response
  },
  function(error) {
    if (error.response.status === 400) {
      store.commit('clearSnackBar')
      store.commit('showSnackBar', {
        text: 'Error: ' + error.response.data.message,
        color: 'error',
      })
    }
    if (error.response.status === 401) {
      store.commit('clearSnackBar')
      store.commit('showSnackBar', {
        text: 'Error: 请登录后访问',
        color: 'error',
      })
      store.commit('logout')
      router.push({ name: 'login' })
    }
    if (error.response.status === 403) {
      store.commit('showSnackBar', {
        text: 'Error: 无权限获取该数据',
        color: 'error',
      })
    }
    // if (error.response.status === 500) {
    //   store.commit('showSnackBar', {
    //     text: 'Error: 服务器故障',
    //     color: 'error',
    //   })
    // }
    if (error.response.status === 504) {
      store.commit('showSnackBar', {
        text: 'Error: 服务器响应超时',
        color: 'error',
      })
    }
    return Promise.reject(error)
  },
)

export * from './logs'
