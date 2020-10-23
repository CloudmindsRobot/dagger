import { cmdbValidatePassword } from '@/api'

const asyncValidatePassword = {
  data() {
    return {
      password: null,
      passwordErrors: [],
    }
  },
  watch: {
    password() {
      this.cmdbValidatePassword()
    },
  },
  methods: {
    async cmdbValidatePassword() {
      try {
        const res = await cmdbValidatePassword({ password: this.password })
        if (res.status === 200) {
          if (res.data.error) {
            this.passwordErrors = [res.data.error]
          } else {
            this.passwordErrors = []
          }
        }
      } catch (err) {
        if (err.response && err.response.status === 400) {
          this.passwordErrors = ['该字段必填']
        } else {
          this.passwordErrors = ['验证密码请求时发生未知错误']
        }
      }
    },
  },
}

export default asyncValidatePassword
