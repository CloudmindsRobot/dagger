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
        <template v-slot:item.create_at="{ item }">
          <span>{{ new Date(item.create_at).toLocaleString() }}</span>
        </template>
        <template v-slot:item.start_time="{ item }">
          <span>{{ new Date(item.start_time).toLocaleString() }}</span>
        </template>
        <template v-slot:item.end_time="{ item }">
          <span>{{ new Date(item.end_time).toLocaleString() }}</span>
        </template>
        <template v-slot:item.download_url="{ item }">
          <span>
            <v-btn text color="primary" small @click="handleView(item)">
              {{ item.download_url }}
            </v-btn>
          </span>
        </template>
        <template v-slot:item.action="{ item }">
          <span class="pr-4">
            <v-btn
              color="primary"
              :loading="item.downloading"
              small
              @click="handleDownload(item)"
              >下载</v-btn
            >
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
    <DeleteLokiQuerySnapshot ref="deleteLokiQuerySnapshot" @refresh="refresh" />
  </div>
</template>

<script>
import { listQueryResultSnapshot } from '@/api'
import DeleteLokiQuerySnapshot from './components/delete'

export default {
  name: 'LogSnapshot',
  components: {
    DeleteLokiQuerySnapshot,
  },
  data: () => ({
    breadcrumbs: [
      { text: 'LOGS', disabled: true, href: '' },
      { text: '日志查询', disabled: true },
      { text: '日志快照', disabled: true },
    ],
    loading: false,
    items: [],
    headers: [
      { text: '快照名称', value: 'name', align: 'start' },
      { text: '快照容量（条）', value: 'count', align: 'start' },
      { text: '快照地址', value: 'download_url', align: 'start' },
      { text: '快照起始时间', value: 'start_time', align: 'start' },
      { text: '快照终止时间', value: 'end_time', align: 'start' },
      { text: '保存时间', value: 'create_at', align: 'start' },
      { text: '操作', value: 'action', align: 'start', width: 180 },
    ],
    pageCount: 0,
    params: {
      page: 1,
      page_size: 10,
    },
  }),
  computed: {},
  methods: {
    async listQueryResultSnapshot() {
      this.loading = true
      try {
        const res = await listQueryResultSnapshot(this.params)
        if (res.status === 200 && res.data.success) {
          this.items = res.data.data
          this.items.forEach((item) => {
            item.downloading = false
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
            text: 'Error: 获取查询结果失败',
            color: 'error',
          })
        }
      }
      this.loading = false
    },
    handleDelete(item) {
      this.$refs.deleteLokiQuerySnapshot.item = item
      this.$refs.deleteLokiQuerySnapshot.dialog = true
    },
    handleDetail(item) {
      this.$refs.detailLokiQuerySnapshot.item = item
      this.$refs.detailLokiQuerySnapshot.items = []
      this.$refs.detailLokiQuerySnapshot.detailQueryResultSnapshot()
      this.$refs.detailLokiQuerySnapshot.dialog = true
    },
    handleView(item) {
      window.open(item.download_url)
    },
    handleDownload(item) {
      try {
        item.downloading = true
        const link = document.createElement('a')
        link.addEventListener('click', function() {
          const paths = item.download_url.split('/')
          const filename = paths[paths.length - 1]
          const dir = paths[paths.length - 2]
          link.download = filename
          link.href = `/api/v1/loki/static/snapshot/${dir}/${filename}`
        })
        const e = document.createEvent('MouseEvents')
        e.initEvent('click', false, false)
        link.dispatchEvent(e)
        this.$store.commit('showSnackBar', {
          text: 'Success: 下载成功',
          color: 'success',
        })
      } catch (e) {
        this.$store.commit('showSnackBar', {
          text: 'Error: 下载失败',
          color: 'error',
        })
      }
      item.downloading = false
    },
    refresh() {
      this.params.page = 1
      this.listQueryResultSnapshot()
    },
  },
  mounted() {
    if (this.$store.state.jwt) {
      this.listQueryResultSnapshot()
    }
  },
}
</script>
