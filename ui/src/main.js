import Vue from 'vue'
import App from './App.vue'
import router from './router'
import store from './store'
import Vuetify from 'vuetify'

import VuetifyConfirm from 'vuetify-confirm'
Vue.use(Vuetify)

Vue.config.productionTip = false

const vuetify = new Vuetify({
  icons: {
    iconfont: 'md',
  },
})

Vue.use(VuetifyConfirm, {
  vuetify,
})

new Vue({
  vuetify: vuetify,
  router,
  store,
  render: (h) => h(App),
}).$mount('#app')
