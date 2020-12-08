<template>
  <v-app id="inspire">
    <v-container fluid>
      <router-view></router-view>
    </v-container>
    <v-snackbar
      v-for="(snackBar, index) in filterSnackBarItems"
      top
      :multi-line="snackBar.text.length > 50"
      :style="{ 'margin-top': `${index * 50}px` }"
      :key="index"
      :color="snackBar.color"
      :value="snackBar.value"
      :timeout="-1"
    >
      {{ snackBar.text }}
      <template v-slot:action="{ attrs }">
        <v-btn
          text
          color="white"
          v-bind="attrs"
          @click="() => closeSnackBar(index)"
        >
          x
        </v-btn>
      </template>
    </v-snackbar>
  </v-app>
</template>
<script>
import { mapState } from 'vuex'
import { loadSettings } from '@/api'
export default {
  name: 'App',
  computed: {
    ...mapState(['snackBarItems']),
    filterSnackBarItems() {
      return this.snackBarItems.filter((item) => item.value).reverse()
    },
  },
  methods: {
    closeSnackBar(index) {
      this.$store.commit('closeSnackBar', index)
    },
    async loadSettings() {
      try {
        const res = await loadSettings()
        if (res.status === 200) {
          const settings = res.data.data
          this.$store.commit('setSettings', settings)
        } else {
          this.$store.commit('showSnackBar', {
            text: `Error: ${res.data.message}`,
            color: 'error',
          })
        }
      } catch (err) {
        this.$store.commit('showSnackBar', {
          text: 'Error: 加载配置失败',
          color: 'error',
        })
      }
    },
  },
  async mounted() {
    await this.loadSettings()
  },
}
</script>
