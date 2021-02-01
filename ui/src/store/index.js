import Vue from 'vue'
import Vuex from 'vuex'
Vue.use(Vuex)

import { getCookie, setCookie, delCookie } from '@/utils/cookie'

const CSRFTOKENNAME = 'csrftoken'
const JWTName = 'DaggerJWT'
const TimeOut = 6000
const Username = 'username'
const Settings = 'settings'

export default new Vuex.Store({
  state: {
    // global
    sidebar: false,
    drawer: true,
    app: null,
    snackBarItems: [],
    csrftoken: getCookie(CSRFTOKENNAME),
    jwt: getCookie(JWTName),
    username: getCookie(Username),
    settings: getCookie(Settings),
  },
  mutations: {
    // global
    switchSidebar(state) {
      state.sidebar = !state.sidebar
    },
    setSidebar(state, value) {
      state.sidebar = !!value
    },
    switchDrawer(state) {
      state.drawer = !state.drawer
    },
    setDrawer(state, value) {
      state.drawer = !!value
    },
    setApp(state, value) {
      state.app = value
    },
    showSnackBar(state, { text, color }) {
      if (!color) {
        color = 'success'
      }
      const index = state.snackBarItems.push({ text, color, value: true }) - 1
      setTimeout(() => {
        state.snackBarItems[index].value = false
      }, TimeOut)
    },
    closeSnackBar(state, index) {
      state.snackBarItems[state.snackBarItems.length - 1 - index].value = false
    },
    clearSnackBar(state) {
      state.snackBarItems.forEach((item) => {
        item.value = false
      })
    },
    setJwt(state, jwt) {
      state.jwt = jwt
      setCookie(
        JWTName,
        jwt,
        JSON.parse(window.atob(jwt.split('.')[1])).exp * 1000,
      )
    },
    logout(state) {
      state.jwt = null
      delCookie(JWTName)
      delCookie(Username)
    },
    setUsername(state, username) {
      state.username = username
      setCookie(Username, username)
    },
    setSettings(state, settings) {
      state.settings = settings
      setCookie(Settings, settings)
    },
  },
  actions: {},
})
