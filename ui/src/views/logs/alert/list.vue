<template>
  <div>
    <v-breadcrumbs :items="breadcrumbs" divider="/"></v-breadcrumbs>
    <v-card>
      <v-card-title>
        <span class="pr-4">
          <v-btn
            color="primary"
            @click="handlerDownloading"
            :loading="downloading"
          >
            导出规则
          </v-btn>
        </span>
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
        <template v-slot:item.groups="{ item }">
          <v-chip
            color="primary"
            v-for="(group, index) in item.groups"
            :key="index"
            style="margin: 1px 3px;"
          >
            <span class="pr-2">{{ group.log_user_group.group_name }}</span>
          </v-chip>
        </template>
        <template v-slot:item.labels="{ item }">
          <v-chip
            color="primary"
            v-for="(label, index) in item.labels"
            :key="index"
            style="margin: 1px 3px;"
          >
            <span class="pr-2">{{ label.key }}:{{ label.value }}</span>
          </v-chip>
        </template>
        <template v-slot:item.create_at="{ item }">
          <span>{{ new Date(item.create_at).toLocaleString() }}</span>
        </template>
        <template v-slot:item.own="{ item }">
          <span v-if="item.user.username === username">Owner</span>
          <span v-else>Participant</span>
        </template>
        <template v-slot:item.action="{ item }">
          <span class="pr-4">
            <v-btn color="primary" small @click="handleAdvanceQuery(item)">
              查询
            </v-btn>
          </span>
          <span class="pr-4" v-if="item.user.username === username">
            <UpdateRule ref="updateRule" :item.sync="item" @refresh="refresh" />
          </span>
          <span class="pr-4" v-if="item.user.username === username">
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
          @input="listLogRule"
        ></v-pagination>
      </div>
    </v-card>

    <DeleteRule ref="deleteRule" @refresh="refresh" />
  </div>
</template>

<script>
import { listLogRule, downloadRuleFile } from '@/api'
import DeleteRule from './components/DeleteRule'
import UpdateRule from './components/UpdateRule'
import { mapState } from 'vuex'

export default {
  name: 'LogRule',
  components: {
    DeleteRule,
    UpdateRule,
  },
  data: () => ({
    breadcrumbs: [
      { text: 'ALERTS', disabled: true, href: '' },
      { text: '日志告警', disabled: true, href: '' },
      { text: '告警规则', disabled: true },
    ],
    items: [],
    loading: false,
    downloading: false,
    keyCN: {
      app: '应用',
      env: '环境',
      svc: '组件/宿主机',
    },
    pageCount: 0,
    params: {
      page: 1,
      page_size: 10,
      online: 0,
    },
    toggle: 0,
  }),
  computed: {
    ...mapState(['username', 'settings']),
    headers() {
      if (this.settings.allowSignUp) {
        return [
          { text: '规则名称', value: 'name', align: 'start' },
          { text: 'LogQL', value: 'log_ql', align: 'start', width: 500 },
          { text: '组', value: 'groups', align: 'start' },
          { text: '描述', value: 'description', align: 'start' },
          { text: '标签', value: 'labels', align: 'start' },
          { text: '创建时间', value: 'create_at', align: 'start' },
          { text: '权属', value: 'own', align: 'start' },
          { text: '操作', value: 'action', align: 'start', width: 240 },
        ]
      } else {
        return [
          { text: '规则名称', value: 'name', align: 'start' },
          { text: 'LogQL', value: 'log_ql', align: 'start', width: 500 },
          { text: '描述', value: 'description', align: 'start' },
          { text: '标签', value: 'labels', align: 'start' },
          { text: '创建时间', value: 'create_at', align: 'start' },
          { text: '权属', value: 'own', align: 'start' },
          { text: '操作', value: 'action', align: 'start', width: 240 },
        ]
      }
    },
  },
  methods: {
    async listLogRule() {
      this.loading = true
      try {
        this.selected = []
        const res = await listLogRule(this.params)
        if (res.status === 200 && res.data.success) {
          this.items = res.data.data
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
            text: 'Error: 获取规则失败',
            color: 'error',
          })
        }
      }
      this.loading = false
    },
    async handlerDownloading() {
      this.downloading = true
      try {
        const res = await downloadRuleFile()
        if (res.status === 200) {
          const filename = res.data.download
          const link = document.createElement('a')
          link.addEventListener('click', function() {
            link.download = filename
            link.href = '/api/v1/loki/static/rules/' + filename
          })
          const e = document.createEvent('MouseEvents')
          e.initEvent('click', false, false)
          link.dispatchEvent(e)
          this.$store.commit('showSnackBar', {
            text: 'Success: 导出成功',
            color: 'success',
          })
        }
      } catch (err) {
        if (
          err.response &&
          [400, 401, 403, 504].indexOf(err.response.status) === -1
        ) {
          this.$store.commit('showSnackBar', {
            text: 'Error: 获取规则失败',
            color: 'error',
          })
        }
      }
      this.downloading = false
    },
    handlerChangeToggleBtn() {
      this.params.online = this.toggle
      this.listLogRule()
    },
    refresh() {
      this.params.page = 1
      this.listLogRule()
    },
    handleDelete(item) {
      this.$refs.deleteRule.item = item
      this.$refs.deleteRule.dialog = true
    },
    handleAdvanceQuery(item) {
      this.$router.push({
        name: 'loki-viewer',
        query: { logQL: encodeURIComponent(item.log_ql), advanced: true },
      })
    },
  },
  mounted() {
    if (this.$store.state.jwt) {
      this.listLogRule()
    }
  },
}
</script>
