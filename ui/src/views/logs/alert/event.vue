<template>
  <div>
    <v-breadcrumbs :items="breadcrumbs" divider="/"></v-breadcrumbs>
    <v-card>
      <v-card-title>
        <span class="pr-4">
          <v-btn
            color="primary"
            @click="handlerArchive"
            :loading="archiveLoading"
          >
            归档
          </v-btn>
        </span>
        <v-btn-toggle
          dense
          v-model="toggle"
          color="primary"
          @change="handlerChangeToggleBtn"
        >
          <v-btn>告警</v-btn>
          <v-btn>已归档</v-btn>
        </v-btn-toggle>
        <v-btn icon text color="green" @click="refresh" :loading="loading">
          <v-icon>refresh</v-icon>
        </v-btn>
        <v-spacer></v-spacer>
        <v-spacer></v-spacer>

        <v-text-field
          append-icon="search"
          label="告警事件"
          single-line
          hide-details
          v-model="params.search"
          @keypress.13="listLogEvent(true)"
        ></v-text-field>
      </v-card-title>
      <v-subheader>共计：{{ total }}条目</v-subheader>
      <v-data-table
        :headers="headers"
        :items="items"
        :loading="loading"
        :expanded.sync="expanded"
        loading-text="载入中..."
        disable-pagination
        single-expand
        show-expand
        hide-default-header
        hide-default-footer
        no-data-text="无数据"
        dense
        @update:expanded="handlerEventDetails"
      >
        <template v-slot:item.info="{ item }">
          <div :class="'v-list-item-html v-list-item-' + item.rule.level">
            <v-btn color="primary" fab x-small dark>{{ item.count }}</v-btn>
            <span class="message">
              {{ new Date(item.create_at).toLocaleString() }}
            </span>
            <v-tooltip top max-width="700">
              <template v-slot:activator="{ on }">
                <span v-on="on">{{ item.rule.name }}</span>
              </template>
              <span style="word-break: break-all;">
                描述：{{ item.rule.description }}<br />
                LogQL：{{ item.rule.log_ql }}
              </span>
            </v-tooltip>
          </div>
        </template>
        <template v-slot:expanded-item="{ item }">
          <td
            style="padding: 10px;word-break: break-all;"
            :colspan="headers.length"
          >
            <div class="v-event-details">
              <div
                :class="'v-list-item-html v-list-item-' + item.rule.level"
                v-for="item in detailItems"
                :key="item.id"
              >
                <span class="message">
                  {{ new Date(item.starts_at).toLocaleString() }}
                </span>
                <span>{{ item.description }}</span>
              </div>
            </div>
          </td>
        </template>
      </v-data-table>
      <div
        v-if="total >= items.length && items.length > 0 && total >= 500"
        style="text-align: center;"
      >
        <v-btn text x-small color="primary" @click="handlerLoadNext">
          点击继续加载...
        </v-btn>
      </div>
    </v-card>
  </div>
</template>

<style scoped>
.v-list-item-html {
  font-size: 13px;
  font-weight: normal;
  padding: 1px 10px;
  margin: 2px 0;
  line-height: 1.2 !important;
  word-break: break-all;
  min-height: 0px !important;
}
.v-list-item-warn {
  border-left: 3px solid #fb8c00 !important;
}
.v-list-item-info {
  border-left: 3px solid #4caf50 !important;
}
.v-list-item-debug {
  border-left: 3px solid #5cbbf6 !important;
}
.v-list-item-error {
  border-left: 3px solid #ff5252 !important;
}
.v-list-item-unknown {
  border-left: 3px solid #607d8b !important;
}
.message {
  margin-right: 10px;
}
.v-btn--fab.v-size--x-small {
  width: 22px;
  height: 22px;
}
.v-event-details {
  max-height: 500px;
  overflow-y: auto;
  padding-left: 15px;
}
</style>

<script>
import { listLogEvent, archiveLogEvent, listLogEventDetails } from '@/api'

export default {
  name: 'LogEvent',
  components: {},
  data: () => ({
    breadcrumbs: [
      { text: 'ALERTS', disabled: true, href: '' },
      { text: '日志告警', disabled: true, href: '' },
      { text: '告警事件', disabled: true },
    ],
    items: [],
    detailItems: [],
    loading: false,
    archiveLoading: false,
    expanded: [],
    headers: [
      { text: '', value: 'info', align: 'start' },
      { text: '', value: 'data-table-expand', align: 'end' },
    ],
    total: 0,
    params: {
      page: 1,
      page_size: 500,
      status: 'firing',
      search: null,
    },
    toggle: 0,
  }),
  computed: {},
  methods: {
    async listLogEvent(refresh) {
      this.loading = true
      try {
        const res = await listLogEvent(this.params)
        if (res.status === 200 && res.data.success) {
          if (refresh) {
            this.items = res.data.data
          } else {
            this.items = this.items.concat(res.data.data)
          }
          this.total = res.data.total
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
            text: 'Error: 获取告警事件失败',
            color: 'error',
          })
        }
      }
      this.loading = false
    },
    handlerChangeToggleBtn() {
      this.params.page = 1
      if (this.toggle === 0) {
        this.params.status = 'firing'
      } else if (this.toggle === 1) {
        this.params.status = 'resolved'
      } else {
        this.params.status = null
      }
      this.listLogEvent(true)
    },
    async handlerArchive() {
      const eventids = []
      this.items.forEach((item) => {
        eventids.push(item.id)
      })
      const res = await this.$confirm('归档告警事件？', {
        buttonTrueText: '确定',
        buttonFalseText: '取消',
        persistent: true,
      })
      if (res) {
        this.archiveLoading = true
        try {
          const res = await archiveLogEvent(eventids)
          if (res.status === 201) {
            this.listLogEvent(true)
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
              text: 'Error: 归档告警事件失败',
              color: 'error',
            })
          }
        }
        this.archiveLoading = false
      }
    },
    handlerLoadNext() {
      this.params.page++
      this.listLogEvent(false)
    },
    refresh() {
      this.params.page = 1
      this.listLogEvent(true)
    },
    async handlerEventDetails(data) {
      try {
        const res = await listLogEventDetails(data[0])
        if (res.status === 200 && res.data.success) {
          this.detailItems = res.data.data
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
            text: 'Error: 获取告警事件详情失败',
            color: 'error',
          })
        }
      }
    },
  },
  mounted() {
    if (this.$store.state.jwt) {
      this.listLogEvent(true)
    }
  },
}
</script>
