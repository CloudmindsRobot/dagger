<template>
  <div :style="{ transform: 'translate(0%, 30%)' }">
    <v-row align-content="center">
      <v-col cols="12" sm="8" md="6" lg="4" :offset="offset">
        <v-card class="elevation-12">
          <v-main class="primary">
            <v-card-title primary-title class="elevation-2">
              <v-card-text class="text-center">
                <h1 style="color: #FFFFFF">创建用户</h1>
              </v-card-text>
            </v-card-title>
          </v-main>
          <v-form v-model="valid">
            <v-card-text>
              <v-text-field
                label="用户名"
                :rules="rules.requiredRules"
                required
                v-model="item.username"
              >
              </v-text-field>
              <v-text-field
                label="密码"
                :rules="rules.requiredRules"
                required
                :type="show ? 'text' : 'password'"
                :append-icon="show ? 'visibility' : 'visibility_off'"
                @click:append="show = !show"
                v-model="password"
              >
              </v-text-field>
              <v-text-field
                label="确认密码"
                :rules="rules.requiredRules.concat([contrastPassword])"
                required
                :type="show ? 'text' : 'password'"
                :append-icon="show ? 'visibility' : 'visibility_off'"
                @click:append="show = !show"
                v-model="confirmed_password"
              >
              </v-text-field>
              <v-text-field
                label="邮箱地址"
                :rules="rules.requiredRules.concat(rules.emailRules)"
                required
                v-model="item.email"
              >
              </v-text-field>
            </v-card-text>
            <v-card-actions>
              <v-row justify="space-around">
                <v-col cols="6" md="3">
                  <v-btn color="primary" @click="register" block
                    >确认创建</v-btn
                  >
                </v-col>
                <v-col cols="6" md="3">
                  <v-btn @click="login" block>返回登录</v-btn>
                </v-col>
              </v-row>
            </v-card-actions>
          </v-form>
        </v-card>
      </v-col>
    </v-row>
  </div>
</template>

<script>
import { register } from '@/api'
import { rules } from '@/utils/helpers'

export default {
  name: 'Register',
  data: () => ({
    rules,
    show: false,
    valid: false,
    item: {},
    password: null,
    confirmed_password: null,
  }),
  computed: {
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
  methods: {
    login() {
      this.$router.push({ name: 'login' })
    },
    async register() {
      try {
        this.item.password = this.password
        const res = await register(this.item)
        if (res.status === 201) {
          this.item = {}
          this.dialog = false
          this.$store.commit('showSnackBar', {
            text: 'Success: 创建用户成功',
            color: 'success',
          })
          this.$router.push({ name: 'login' })
        } else {
          this.$store.commit('showSnackBar', {
            text: `Error: ${res.data.message}`,
            color: 'error',
          })
        }
      } catch (err) {
        this.$store.commit('showSnackBar', {
          text: 'Error: 创建用户失败',
          color: 'error',
        })
      }
    },
    contrastPassword() {
      if (this.password !== this.confirmed_password) {
        return '两次输入密码不一致'
      }
      return true
    },
  },
}
</script>
