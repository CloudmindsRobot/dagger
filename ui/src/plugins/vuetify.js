import Vue from 'vue'
import Vuetify from 'vuetify/lib'
import '@/scss/vuetify/overrides.scss'
import VuetifyConfirm from 'vuetify-confirm'
Vue.use(Vuetify)

const theme = {
  primary: '#1e88e5',
  info: '#1e88e5',
  success: '#21c1d6',
  accent: '#fc4b6c',
  default: '#563dea',
}

const vuetify = new Vuetify({
  icons: {
    iconfont: 'md',
  },
  theme: {
    themes: {
      dark: theme,
      light: theme,
    },
  },
})

Vue.use(VuetifyConfirm, {
  vuetify,
})

export default vuetify
