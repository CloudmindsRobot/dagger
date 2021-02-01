<template>
  <v-dialog v-model="dialog" max-width="500">
    <v-card class="px-1">
      <v-card-title>
        <span class="headline">删除查询快照</span>
      </v-card-title>
      <v-divider></v-divider>
      <v-card-text>
        <v-list
          class="grey lighten-3"
          dense
          style="border-radius: 5px 5px;margin-top: 10px;"
        >
          <v-list-item tag="span">
            <v-list-item-content>
              <v-list-item-title>快照: {{ item.name }}</v-list-item-title>
            </v-list-item-content>
          </v-list-item>
        </v-list>
        <p class="orange--text" style="margin: 5px">{{ notice }}</p>
      </v-card-text>
      <v-divider></v-divider>
      <v-card-actions>
        <v-row justify="space-around">
          <v-col cols="6" md="3">
            <v-btn color="error" @click="deleteQuerySnapshot(item)" block
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
import { deleteQueryResultSnapshot } from '@/api'

export default {
  name: 'DeleteLokiQuerySnapshot',
  data: () => ({
    dialog: false,
    notice: '注：删除Loki查询快照文件',
    loading: false,
    item: {},
  }),
  methods: {
    async deleteQuerySnapshot() {
      this.loading = true
      try {
        const res = await deleteQueryResultSnapshot({ id: this.item.id })
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
            text: 'Error: 删除查询快照失败',
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
