<template>
  <div :style="{ transform: 'translate(0%, 50%)' }">
    <v-row align-content="center">
      <v-col cols="12" sm="8" md="6" lg="4" :offset="offset">
        <v-card class="elevation-12">
          <v-main class="primary">
            <v-card-title primary-title class="elevation-2">
              <v-card-text class="text-center">
                <h1 style="color: #FFFFFF">DAGGER</h1>
              </v-card-text>
            </v-card-title>
          </v-main>
          <v-form v-model="valid" @submit.prevent="login">
            <v-card-text>
              <v-text-field
                prepend-icon="person"
                label="用户名"
                :rules="rules.requiredRules"
                required
                placeholder="username"
                v-model="item.username"
              ></v-text-field>
              <v-text-field
                prepend-icon="lock"
                label="密码"
                :rules="rules.requiredRules"
                required
                placeholder="password"
                :type="show ? 'text' : 'password'"
                :append-icon="show ? 'visibility' : 'visibility_off'"
                @click:append="show = !show"
                v-model="item.password"
              ></v-text-field>
            </v-card-text>

            <v-card-actions>
              <v-btn block color="primary" type="submit" dark>登录</v-btn>
            </v-card-actions>
          </v-form>

          <v-card-text class="text-center">
            <router-link :to="{ name: 'register' }">注册账号</router-link>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>
  </div>
</template>

<script>
import { mapState } from 'vuex'
import { getUserSelf, login } from '@/api'
import { rules, validateJWT } from '@/utils/helpers'

export default {
  name: 'Login',
  data: () => ({
    rules,
    valid: false,
    show: false,
    item: {
      username: '',
      password: '',
    },
  }),
  watch: {},
  methods: {
    async login() {
      try {
        if (this.item.username.length === 0) {
          this.$store.commit('showSnackBar', {
            text: 'Warn: 请输入用户名',
            color: 'warning',
          })
          return
        }
        if (this.item.password.length === 0) {
          this.$store.commit('showSnackBar', {
            text: 'Warn: 请输入密码',
            color: 'warning',
          })
          return
        }
        const res = await login(this.item)
        if (res.status === 200 && res.data.success) {
          this.$store.commit('setJwt', res.data.token)
          this.$store.commit('showSnackBar', {
            text: 'Success: 登录成功',
            color: 'success',
          })
          await this.getUserSelf()
          if (this.$route.query.redirect !== undefined) {
            this.$router.push({ path: this.$route.query.redirect })
          } else {
            this.$router.push({ name: 'loki-viewer' })
          }
        } else {
          this.$store.commit('showSnackBar', {
            text: `Error: ${res.data.message}`,
            color: 'error',
          })
          this.item = {}
        }
      } catch (err) {
        if (
          err.response &&
          err.response.status >= 400 &&
          err.response.status < 500
        ) {
          this.$store.commit('showSnackBar', {
            text: 'Error: 登录失败，用户名或密码错误',
            color: 'error',
          })
          this.item = {}
        } else {
          this.$store.commit('showSnackBar', {
            text: 'Error: 登录失败，程序错误',
            color: 'error',
          })
        }
      }
    },
    async getUserSelf() {
      try {
        const res = await getUserSelf()
        if (res.status === 200) {
          const user = res.data.user
          this.$store.commit('setUsername', user.username)
        }
      } catch (err) {
        this.$store.commit('showSnackBar', {
          text: 'Error: 获取用户数据失败',
          color: 'error',
        })
      }
    },
    async init() {
      if (this.jwt && validateJWT(this.jwt)) {
        await this.getUserSelf()
        this.$store.commit('showSnackBar', {
          text: 'Warning: 已登录状态',
          color: 'warning',
        })
        this.$router.push({ name: 'loki-viewer' })
      }
    },
  },
  computed: {
    ...mapState(['jwt']),
    offset() {
      switch (this.$vuetify.breakpoint.name) {
        case 'xs':
          return 0
        case 'sm':
          return 2
        case 'md':
          return 3
        case 'lg':
          return 4
        case 'xl':
          return 4
      }
      return 0
    },
  },
  mounted() {
    this.init()
  },
}
</script>
