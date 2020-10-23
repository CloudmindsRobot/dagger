import Vue from 'vue'
import Router from 'vue-router'
import store from '@/store'

import { global } from './global'
import { logs } from './logs'

Vue.use(Router)

const router = new Router({
  routes: global.concat(logs),
})

router.beforeEach(async function(to, from, next) {
  if (to.meta.requireAuth && !store.state.jwt) {
    store.commit('showSnackBar', {
      text: 'Error: 请登录后访问',
      color: 'error',
    })
    next({ name: 'login', query: { redirect: to.fullPath } })
  } else {
    next()
  }
})

export default router
