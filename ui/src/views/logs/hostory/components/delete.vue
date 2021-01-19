<template>
  <v-dialog v-model="dialog" max-width="500">
    <v-card class="px-1">
      <v-card-title>
        <span class="headline">删除查询标签</span>
      </v-card-title>
      <v-divider></v-divider>
      <v-card-text>
        <p style="margin-top: 10px;">
          <v-chip
            color="primary"
            label
            small
            v-for="(v, k) in item.label_json"
            :key="k"
            style="margin:5px 5px 0;"
          >
            <span class="pr-2">标签({{ k }}):{{ v }}</span>
          </v-chip>
          <v-chip
            color="primary"
            label
            small
            v-for="it in item.filter_json"
            :key="it"
            style="margin:5px 5px 0;"
          >
            <span class="pr-2">正则(regex):{{ it }}</span>
          </v-chip>
        </p>
        <p class="orange--text" style="margin: 5px">{{ notice }}</p>
      </v-card-text>
      <v-divider></v-divider>
      <v-card-actions>
        <v-row justify="space-around">
          <v-col cols="6" md="3">
            <v-btn color="error" @click="deleteQueryHistoryLabel(item)" block
              >确定删除</v-btn
            >
          </v-col>
          <v-col cols="6" md="3">
            <v-btn @click="dialog = false" block>取消</v-btn>
          </v-col>
        </v-row>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script>
import { deleteQueryHistoryLabel } from '@/api'

export default {
  name: 'DeleteLokiHistoryQuery',
  data: () => ({
    dialog: false,
    notice: '注：删除Loki查询Label',
    loading: false,
    item: {},
  }),
  methods: {
    async deleteQueryHistoryLabel() {
      this.loading = true
      try {
        const res = await deleteQueryHistoryLabel({ id: this.item.id })
        if (res.status === 204) {
          this.$emit('refresh')
        } else {
          this.$store.commit('showSnackBar', {
            text: `Warn: ${res.data.message}`,
            color: 'warning',
          })
        }
      } catch (err) {
        if (
          err.response &&
          [400, 401, 403, 504].indexOf(err.response.status) === -1
        ) {
          this.$store.commit('showSnackBar', {
            text: 'Error: 删除查询历史失败',
            color: 'error',
          })
        }
      }
      this.loading = false
      this.dialog = false
    },
  },
}
</script>
