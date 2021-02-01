import Vue from 'vue'
import App from './App.vue'
import router from './router'
import store from './store'
import Vuetify from 'vuetify'
import 'vuetify/dist/vuetify.min.css'

// import 'echarts'
// 手动引入 ECharts 各模块来减小打包体积
import 'echarts/lib/chart/pie'
import 'echarts/lib/chart/line'
import 'echarts/lib/chart/tree'
import 'echarts/lib/chart/bar'
import 'echarts/lib/component/tooltip'
import 'echarts/lib/component/title'
import 'echarts/lib/component/legend'
import 'echarts/lib/component/legendScroll'
import ECharts from 'vue-echarts'
import VuetifyConfirm from 'vuetify-confirm'

Vue.component('v-chart', ECharts)

// Vue.use(Vuetify, {
//   iconfont: 'md', // 'md' || 'mdi' || 'fa' || 'fa4'   mdi会出现checkbox不显示问题
// })
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
