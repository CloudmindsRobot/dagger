<template>
  <div>
    <BaseSidebar />
    <BaseToolbar />
    <v-main class="grey lighten-4">
      <router-view></router-view>
    </v-main>
  </div>
</template>

<script>
import { mapState } from 'vuex'
import BaseToolbar from '@/components/BaseToolbar.vue'
import BaseSidebar from '@/components/BaseSidebar.vue'
import { validateJWT } from '@/utils/helpers'

export default {
  name: 'Root',
  components: {
    BaseToolbar,
    BaseSidebar,
  },
  computed: {
    ...mapState(['jwt']),
  },
  mounted() {
    if (!validateJWT(this.jwt)) {
      this.$store.commit('showSnackBar', {
        text: 'Error: 登录已过期',
        color: 'error',
      })
      this.$router.push({ name: 'login' })
    }
  },
}
</script>
<style>
.v-data-table table tbody tr td {
  font-size: 11px;
}
</style>
