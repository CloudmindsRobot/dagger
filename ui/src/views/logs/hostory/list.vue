<template>
  <div>
    <v-breadcrumbs :items="breadcrumbs" divider="/"></v-breadcrumbs>
    <v-card>
      <v-card-title>
        <v-btn icon text color="green" @click="refresh" :loading="loading">
          <v-icon>refresh</v-icon>
        </v-btn>
      </v-card-title>
      <v-data-table
        :headers="headers"
        :items="items"
        :page.sync="params.page"
        :items-per-page="params.page_size"
        hide-default-footer
      >
        <template v-slot:item.label="{ item }">
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
        </template>
        <template v-slot:item.create_at="{ item }">
          <span>{{ new Date(item.create_at).toLocaleString() }}</span>
        </template>
        <template v-slot:item.action="{ item }">
          <span class="pr-4">
            <v-btn color="primary" small @click="handleQuery(item)">查询</v-btn>
          </span>
          <span class="pr-4">
            <v-btn color="error" small @click="handleDelete(item)">删除</v-btn>
          </span>
        </template>
      </v-data-table>
      <div class="text-xs-center pa-2">
        <v-pagination
          v-model="params.page"
          :length="pageCount"
          :total-visible="10"
          circle
        ></v-pagination>
      </div>
    </v-card>
    <DeleteLokiQueryLabel ref="deleteLokiQueryLabel" @refresh="refresh" />
  </div>
</template>

<script>
import { listQueryHistory } from '@/api'
import DeleteLokiQueryLabel from './components/delete'

export default {
  name: 'LogHistory',
  components: {
    DeleteLokiQueryLabel,
  },
  data: () => ({
    breadcrumbs: [
      { text: 'LOGS', disabled: true, href: '' },
      { text: '日志查询', disabled: true },
      { text: '查询历史', disabled: true },
    ],
    items: [],
    loading: false,
    headers: [
      { text: '查询标签', value: 'label', align: 'start' },
      { text: 'LogQL', value: 'log_ql', align: 'start' },
      { text: '创建时间', value: 'create_at', align: 'start' },
      { text: '操作', value: 'action', align: 'start' },
    ],
    pageCount: 0,
    params: {
      page: 1,
      page_size: 10,
    },
  }),
  computed: {},
  methods: {
    async listQueryHistory() {
      this.loading = true
      try {
        const res = await listQueryHistory(this.params)
        if (res.status === 200 && res.data.success) {
          this.items = res.data.data
          this.items.forEach((item) => {
            item.label_json = JSON.parse(item.label_json)
            item.filter_json = JSON.parse(item.filter_json)
          })
          this.pageCount = Math.ceil(res.data.total / res.data.page_size)
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
            text: 'Error: 获取查询标签失败',
            color: 'error',
          })
        }
      }
      this.loading = false
    },
    handleQuery(item) {
      const label = item.label_json
      const filter = item.filter_json
      this.$router.push({
        name: 'loki-viewer',
        query: Object.assign(label, { filters: filter }),
      })
    },
    handleDelete(item) {
      this.$refs.deleteLokiQueryLabel.item = item
      this.$refs.deleteLokiQueryLabel.dialog = true
    },
    refresh() {
      this.params.page = 1
      this.listQueryHistory()
    },
  },
  mounted() {
    if (this.$store.state.jwt) {
      this.listQueryHistory()
    }
  },
}
</script>
