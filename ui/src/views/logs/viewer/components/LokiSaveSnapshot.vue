<template>
  <v-dialog v-model="dialog" max-width="500">
    <v-card class="px-1">
      <v-card-title>
        <span class="headline">保存快照</span>
      </v-card-title>
      <v-divider></v-divider>
      <v-card-text>
        <v-subheader>基本信息</v-subheader>
        <p style="margin:0 15px 0;">
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
        <div style="margin: 0 15px;">起始时间:{{ item.start_time }}</div>
        <div style="margin: 0 15px;">终止时间:{{ item.end_time }}</div>
        <v-subheader>快照文件</v-subheader>
        <p>
          <v-form v-model="valid" style="margin: 0 15px;">
            <v-text-field
              v-model="name"
              :thumb-size="20"
              label="快照文件名称"
              thumb-label="always"
            ></v-text-field>
          </v-form>
        </p>
      </v-card-text>
      <v-divider></v-divider>
      <v-card-actions>
        <v-row justify="space-around">
          <v-col cols="6" md="3">
            <v-btn color="primary" @click="createQuerySnapshot(item)" block
              >确定</v-btn
            >
          </v-col>
          <v-col cols="6" md="3">
            <v-btn @click="closeDialog" block>取消</v-btn>
          </v-col>
        </v-row>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script>
import { createQueryResultSnapshot } from '@/api'

export default {
  name: 'CreateLokiQuerySnapshot',
  data: () => ({
    dialog: false,
    loading: false,
    valid: false,
    item: {},
    name: '',
    originname: '',
  }),
  methods: {
    async createQuerySnapshot() {
      this.loading = true
      try {
        const res = await createQueryResultSnapshot(
          Object.assign(this.item, {
            name: this.name,
            tmp_filename: this.originname,
          }),
        )
        if (res.status === 201) {
          this.$store.commit('showSnackBar', {
            text: 'Success: 保存查询结果成功',
            color: 'success',
          })
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
            text: 'Error: 保存查询结果失败',
            color: 'error',
          })
        }
      }
      this.loading = false
      this.dialog = false
    },
    closeDialog() {
      this.name = ''
      this.originname = ''
      this.dialog = false
    },
  },
}
</script>
