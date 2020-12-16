<template>
  <div id="logviewer">
    <v-breadcrumbs :items="breadcrumbs" divider="/"></v-breadcrumbs>

    <v-container
      grid-list-xl
      :style="{ margin: 0, padding: 0, 'max-width': '100%' }"
    >
      <v-layout row wrap>
        <v-flex xs12>
          <v-card>
            <v-card-title style="padding: 10px 20px 0 !important;">
              <loki-filter
                ref="lokiFilter"
                :dateRangeTimestamp.sync="dateRangeTimestamp"
              ></loki-filter>
              <span class="pr-4">
                <loki-datetime-range-picker
                  ref="dateRangePicker"
                ></loki-datetime-range-picker>
              </span>
              <span class="pr-4">
                <v-btn
                  color="primary"
                  :loading="loading"
                  @click="handlerQuerying"
                >
                  查询
                </v-btn>
              </span>
              <span class="pr-4">
                <v-btn
                  color="primary"
                  :loading="saveResultLoading"
                  @click="handleSaveResult"
                >
                  保存
                </v-btn>
              </span>
            </v-card-title>
            <loki-histogram
              ref="lokiHistogram"
              :middleStart.sync="dateRangeTimestamp[0]"
              :middleEnd.sync="dateRangeTimestamp[1]"
              @refresh="listQueryRanges"
            ></loki-histogram>
            <v-flex xs12 md12>
              <v-card id="log-card">
                <v-card-title
                  style="font-size: 13px;font-weight: normal;padding:0 16px !important;height: 60px !important;"
                >
                  限制:
                  <span class="pr-4">
                    <v-text-field
                      type="number"
                      v-model="limit"
                      style="width: 60px;margin-left: 5px;"
                    ></v-text-field>
                  </span>
                  <span class="pr-4">结果: {{ items.length }}</span>
                  <span class="pr-4">
                    级别:
                    <v-tooltip top v-for="item in legends" :key="item.level">
                      <template v-slot:activator="{ on }">
                        <span class="pr-4">
                          <v-btn
                            x-small
                            light
                            v-on="on"
                            :outlined="!item.selected"
                            :color="item.color"
                            @click="handlerFilterLevel(item)"
                            >{{ item.level }}</v-btn
                          >
                        </span>
                      </template>
                      <span>{{ item.intro }}</span>
                    </v-tooltip>
                  </span>
                  时间
                  <span class="pr-4" style="margin-left: 10px;">
                    <v-switch v-model="timestamp"></v-switch>
                  </span>
                  倒序
                  <span class="pr-4" style="margin-left: 10px;">
                    <v-switch
                      v-model="dsc"
                      @change="listQueryRanges"
                    ></v-switch>
                  </span>
                </v-card-title>
                <v-card-title
                  style="padding:0px 16px 15px;font-size: 13px;font-weight: normal;"
                >
                  <span
                    class="pr-4"
                    v-for="item in pods"
                    :key="item.text"
                    style="padding-right: 0px !important;height: 24px !important;"
                  >
                    <v-btn
                      class="ma-2"
                      x-small
                      color="primary"
                      :outlined="!item.selected"
                      @click="handlerToggleSelectPod(item)"
                      >{{ item.text }}</v-btn
                    >
                  </span>
                </v-card-title>
                <v-card-text>
                  <v-data-table
                    :headers="headers"
                    :items="items"
                    :loading="loading"
                    :expanded.sync="expanded"
                    loading-text="载入中..."
                    item-key="info.timestamp"
                    disable-pagination
                    single-expand
                    show-expand
                    hide-default-header
                    hide-default-footer
                    no-data-text="无数据"
                    dense
                  >
                    <template v-slot:item.info="{ item }">
                      <div
                        v-if="timestamp"
                        :class="
                          'v-list-item-html v-list-item-' + item.info.level
                        "
                        v-html="
                          `<b style='color: #5cbbf6;margin-right: 10px;'>` +
                            item.info.timestampstr +
                            '</b>' +
                            item.info.message
                        "
                        :style="item.info.animation"
                      ></div>
                      <div
                        v-else
                        :class="
                          'v-list-item-html v-list-item-' + item.info.level
                        "
                        v-html="item.info.message"
                        :style="item.info.animation"
                      ></div>
                      <loki-context
                        v-if="
                          filtered || (level.length > 0 && level.length < 5)
                        "
                        :loki.sync="item"
                        :timestamp.sync="item.info.timestamp"
                        style="line-height: 18px;"
                      ></loki-context>
                    </template>
                    <template v-slot:expanded-item="{ item }">
                      <td colspan="2">
                        <pre
                          style="white-space: pre-wrap;word-break: break-all;"
                          >{{ item.stream }}</pre
                        >
                      </td>
                    </template>
                  </v-data-table>
                </v-card-text>
              </v-card>
            </v-flex>
          </v-card>
        </v-flex>
      </v-layout>
    </v-container>
    <v-btn
      fab
      color="success"
      class="v-btn v-btn--bottom v-btn--contained v-btn--fab v-btn--fixed v-btn--right v-btn--round theme--dark v-size--middle"
      style="bottom: 80px !important;"
      @click="handlerLiveQuerying"
    >
      <v-icon v-text="icon"></v-icon>
    </v-btn>
    <v-speed-dial
      v-model="btnFloat"
      bottom
      right
      direction="left"
      transition="slide-x-reverse-transition"
    >
      <template v-slot:activator>
        <v-btn v-model="btnFloat" color="primary" fab>
          <v-icon v-if="btnFloat">close</v-icon>
          <v-icon v-else>add</v-icon>
        </v-btn>
      </template>
      <v-btn fab small color="primary" @click="handlerGoTo">
        <v-icon>expand_less</v-icon>
      </v-btn>
      <v-btn
        fab
        small
        color="success"
        :loading="downloading"
        @click="handlerDownloading(false)"
      >
        <v-icon>save_alt</v-icon>
      </v-btn>
    </v-speed-dial>
    <LokiSaveSnapshot ref="lokiSaveSnapshot" />
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
#logviewer .v-speed-dial {
  position: fixed;
}

#logviewer .v-btn--floating {
  position: relative;
}
</style>

<script>
import { listQueryRanges, exportQueryRanges, listLabels } from '@/api'
import { mapState } from 'vuex'
import LokiHistogram from './components/LokiHistogram'
import LokiDatetimeRangePicker from './components/LokiDateTimeRangePicker'
import LokiContext from './components/LokiContext'
import LokiFilter from './components/LokiFilter'
import LokiSaveSnapshot from './components/LokiSaveSnapshot'
import { formatDatetime, parserDatetime } from '@/utils/helpers'

export default {
  name: 'LokiViewer',
  components: {
    LokiHistogram,
    LokiDatetimeRangePicker,
    LokiContext,
    LokiFilter,
    LokiSaveSnapshot,
  },
  data: () => ({
    breadcrumbs: [
      { text: 'LOGS', disabled: true, href: '' },
      { text: '日志查询', disabled: true },
      { text: '日志查看器', disabled: true },
    ],
    items: [],
    headers: [
      { text: '', value: 'info', align: 'start' },
      { text: '', value: 'data-table-expand', align: 'end' },
    ],
    expanded: [],
    loading: false,
    saveResultLoading: false,
    disabled: false,
    downloading: false,
    level: [],
    filtered: false,
    dateRangeTimestamp: [],
    limit: 2000,
    size: 20,
    middleStart: '',
    middleEnd: '',
    pods: [],
    legends: [
      {
        level: 'Info',
        color: 'success',
        selected: false,
        intro: '匹配元素:[I],[info],【info】,info,level=info,忽略大小写',
      },
      {
        level: 'Debug',
        color: 'primary',
        selected: false,
        intro: '匹配元素:[D],[debug],【debug】,debug,level=debug,忽略大小写',
      },
      {
        level: 'Warn',
        color: 'warning',
        selected: false,
        intro:
          '匹配元素:[W],[warn],[warning],【warn】,【warning】,warn,warning,level=warn,level=warning,忽略大小写',
      },
      {
        level: 'Error',
        color: 'error',
        selected: false,
        intro: '匹配元素:[E],[error],【error】,error,level=error,忽略大小写',
      },
      {
        level: 'Unknown',
        color: 'blue-grey',
        selected: false,
        intro: '无法匹配前面四个日志级别的其他日志',
      },
    ],
    duration: 500,
    offset: 0,
    easing: 'easeInOutCubic',
    pod: '',
    icon: 'play_circle_outline',
    websocket: null,
    timeoutHandler: null,
    btnFloat: null,
    timestamp: false,
    dsc: true,
    filters: [],
    labels: [],
  }),
  computed: {
    ...mapState(['username']),
  },
  methods: {
    async listQueryRanges() {
      if (this.limit > 50000) {
        this.$store.commit('showSnackBar', {
          text: 'Warn: 最大支持单次50000条日志输出',
          color: 'warning',
        })
        return
      }
      this.loading = true
      try {
        const filterData = {}
        this.filters = []
        this.$refs.lokiFilter.model.forEach((item) => {
          if (item.value === 'filter') {
            this.filters.push(item.text.substr(item.text.indexOf(':') + 1))
          } else {
            filterData[item.value] = item.text.substr(
              item.text.indexOf(':') + 1,
            )
          }
        })
        const data = Object.assign(filterData, {
          start: this.dateRangeTimestamp[0],
          end: this.dateRangeTimestamp[1],
          level: this.level.join(','),
          limit: this.limit,
          size: this.size,
          middleStart: this.middleStart,
          middleEnd: this.middleEnd,
          pod: this.pod,
          dsc: this.dsc,
          filters: this.filters,
        })
        this.filtered =
          data.hasOwnProperty('filters') && data.filters.length > 0
        this.items = []
        const res = await listQueryRanges(data)
        if (res.status === 200) {
          if (res.data === null) {
            this.$refs.lokiHistogram.chartData = {
              'xAxis-data': [],
              'yAxis-data': {
                info: [],
                debug: [],
                error: [],
                warn: [],
                unknown: [],
              },
            }
            this.loading = false
            return
          }
          this.items =
            res.data.query === null
              ? []
              : res.data.query.sort((a, b) => {
                  if (this.dsc) return b.info.timestamp - a.info.timestamp
                  else return a.info.timestamp - b.info.timestamp
                })
          if ((this.items.length * 1.0) / this.limit > 0.8) {
            this.$store.commit('showSnackBar', {
              text: 'Warn: 条目较多，请精确查询条件',
              color: 'warning',
            })
          }
          const xAxisData = []
          res.data.chart['xAxis-data'].forEach((item) => {
            xAxisData.push(
              formatDatetime(new Date(item), 'yyyy-MM-dd hh:mm:ss'),
            )
          })
          res.data.chart['xAxis-data'] = xAxisData
          this.$refs.lokiHistogram.chartData = res.data.chart
          if (res.data.pod !== undefined) {
            this.pods = res.data.pod
          }
          const pod = this.$route.query['pod']
          if (pod !== undefined) {
            this.pod = pod
            this.pods = [{ text: pod, selected: true }]
          }
        } else {
          this.$store.commit('showSnackBar', {
            text: `Error: ${res.data.message}`,
            color: 'error',
          })
        }
      } catch (err) {
        this.$store.commit('showSnackBar', {
          text: 'Error: 获取数据失败',
          color: 'error',
        })
      }
      this.loading = false
    },
    async handlerDownloading(saving) {
      if (this.$refs.lokiFilter.model.length === 0) {
        this.$store.commit('showSnackBar', {
          text: 'Warn: 请进行条件过滤',
          color: 'warning',
        })
        return
      }
      this.downloading = true
      try {
        const filterData = {}
        this.filters = []
        this.$refs.lokiFilter.model.forEach((item) => {
          if (item.value === 'filter') {
            this.filters.push(item.text.substr(item.text.indexOf(':') + 1))
          } else {
            filterData[item.value] = item.text.substr(
              item.text.indexOf(':') + 1,
            )
          }
        })
        const data = Object.assign(filterData, {
          start: this.dateRangeTimestamp[0],
          end: this.dateRangeTimestamp[1],
          level: this.level.join(','),
          pod: this.pod,
          dsc: this.dsc,
          filters: this.filters,
        })
        const res = await exportQueryRanges(data)
        if (res.status === 200) {
          if (res.data === null) {
            this.downloading = false
            return null
          }
          if (res.data.exist) {
            const filename = res.data.download
            if (saving) {
              this.downloading = false
              return filename
            }
            const link = document.createElement('a')
            link.addEventListener('click', function() {
              link.download = filename
              link.href = '/api/v1/loki/static/export/' + filename
            })
            const e = document.createEvent('MouseEvents')
            e.initEvent('click', false, false)
            link.dispatchEvent(e)
            this.$store.commit('showSnackBar', {
              text: 'Success: 导出成功',
              color: 'success',
            })
          } else {
            this.$store.commit('showSnackBar', {
              text: 'Warn: 该区间内没有数据',
              color: 'warning',
            })
            this.downloading = false
            return null
          }
        } else {
          this.$store.commit('showSnackBar', {
            text: `Error: ${res.data.message}`,
            color: 'error',
          })
        }
      } catch (err) {
        this.$store.commit('showSnackBar', {
          text: 'Error: 获取数据失败',
          color: 'error',
        })
      }
      this.downloading = false
    },
    handlerFilterLevel(item) {
      if (this.loading) return
      const index = this.legends.findIndex((i) => i.level === item.level)
      this.$set(this.legends, index, {
        level: item.level,
        color: item.color,
        selected: !item.selected,
        intro: item.intro,
      })
      this.level = []
      const histogramLegends = {
        Info: false,
        Debug: false,
        Warn: false,
        Error: false,
        Unknown: false,
      }
      this.legends.forEach((item) => {
        if (item.selected) {
          this.level.push(item.level)
          histogramLegends[item.level] = true
        }
      })
      if (this.level.length === 0) {
        this.$refs.lokiHistogram.legendSelected = {
          Info: true,
          Debug: true,
          Warn: true,
          Error: true,
          Unknown: true,
        }
      } else {
        this.$refs.lokiHistogram.legendSelected = histogramLegends
      }
      this.listQueryRanges()
    },
    async handleSaveResult() {
      if (this.saveResultLoading) return
      if (this.items.length === 0) {
        this.$store.commit('showSnackBar', {
          text: 'Warn: 无查询结果',
          color: 'warning',
        })
        return
      }
      this.saveResultLoading = true
      this.$store.commit('showSnackBar', {
        text: 'Success: 正在生成结果临时文件，请等待...',
        color: 'success',
      })
      const filename = await this.handlerDownloading(true)
      if (filename !== null) {
        const filters = []
        const labels = {}
        this.$refs.lokiFilter.model.forEach((item) => {
          if (item.value === 'filter') {
            filters.push(item.text.substr(item.text.indexOf(':') + 1))
          } else {
            labels[item.value] = item.text.substr(item.text.indexOf(':') + 1)
          }
        })
        const item = {
          label_json: labels,
          filter_json: filters,
          start_time: this.$refs.dateRangePicker.dateRange[0],
          end_time: this.$refs.dateRangePicker.dateRange[1],
        }
        this.$refs.lokiSaveSnapshot.item = item
        this.$refs.lokiSaveSnapshot.originname = filename
        this.$refs.lokiSaveSnapshot.name = ''
        this.$refs.lokiSaveSnapshot.dialog = true
      }
      this.saveResultLoading = false
    },
    handlerQuerying() {
      if (this.loading) return
      if (this.$refs.dateRangePicker.quick) {
        this.$refs.dateRangePicker.handlerChangeQuickTime()
      }
      this.middleStart = ''
      this.middleEnd = ''
      this.size = 20
      this.$refs.lokiHistogram.start = 0
      this.$refs.lokiHistogram.end = 100
      this.$refs.lokiHistogram.legendSelected = {
        Info: true,
        Debug: true,
        Warn: true,
        Error: true,
        Unknown: true,
      }
      this.level = []
      this.pod = ''
      this.pods = []
      this.legends.forEach((item) => {
        item.selected = false
      })
      this.listQueryRanges()
    },
    handlerGotoOptions() {
      return {
        duration: this.duration,
        offset: this.offset,
        easing: this.easing,
      }
    },
    handlerGoTo() {
      this.$vuetify.goTo(0, this.handlerGotoOptions())
    },
    handlerToggleSelectPod(item) {
      if (this.loading) return
      const index = this.pods.findIndex((i) => i.text === item.text)
      this.$set(this.pods, index, { text: item.text, selected: !item.selected })
      const filterPod = []
      this.pods.forEach((item) => {
        if (item.selected) {
          filterPod.push(item.text)
        }
      })
      this.pod = filterPod.join('|')
      this.listQueryRanges()
    },
    handlerParams() {
      const params = {}
      this.filters = []
      this.$refs.lokiFilter.model.forEach((item) => {
        if (item.value === 'filter') {
          this.filters.push(item.text.substr(item.text.indexOf(':') + 1))
        } else {
          params[item.value] = item.text.substr(item.text.indexOf(':') + 1)
        }
      })
      const data = Object.assign(params, {
        start: Date.parse(new Date()).toString() + '000000',
        level: this.level.join(','),
        pod: this.pod,
        filters: this.filters,
      })
      const paramArray = []
      for (var item in data) {
        paramArray.push(item + '=' + data[item])
      }
      return paramArray.join('&')
    },
    handlerInitWebSocket() {
      const params = this.handlerParams()
      const protocol = window.location.protocol === 'https:' ? 'wss' : 'ws'
      const host = window.location.host
      let wsuri = `${protocol}://${host}/ws/loki/tail?${params}`
      this.websocket = new WebSocket(wsuri)
      this.websocket.onmessage = this.handlerWebsocketOnmessage
      this.websocket.onopen = this.handlerWebsocketOnopen
      this.websocket.onerror = this.handlerWebsocketOnerror
      this.websocket.onclose = this.handlerWebsocketClose
    },
    handlerWebsocketOnopen() {
      this.loading = true
      this.icon = 'pause_circle_outline'
      this.$refs.lokiFilter.disabled = true
      this.$vuetify.goTo('#log-card', this.handlerGotoOptions())
    },
    handlerWebsocketOnerror() {
      this.$store.commit('showSnackBar', {
        text: 'Warning: websocket连接失败',
        color: 'warning',
      })
    },
    handlerWebsocketOnmessage(e) {
      const data = JSON.parse(e.data)
      data.sort((a, b) => {
        return a.info.timestamp - b.info.timestamp
      })
      const instance = this
      data.forEach((item) => {
        this.items.unshift(item)
        if (this.items.length > this.limit) {
          this.items.pop()
        }
      })
      this.timeoutHandler = setTimeout(() => {
        for (var index = 0; index < data.length; index++) {
          instance.items[index].info.animation =
            'transition: background-color 2s;'
        }
        clearTimeout(this.timeoutHandler)
      }, 1)
    },
    handlerWebsocketClose() {
      this.loading = false
      this.icon = 'play_circle_outline'
      this.$refs.lokiFilter.disabled = false
      this.items.forEach((item) => {
        if (item.info !== undefined) item.info.animation = ''
      })
    },
    handlerClose() {
      if (this.websocket) {
        this.websocket.send('close')
        this.websocket.close()
        this.websocket = null
        this.loading = false
        this.icon = 'play_circle_outline'
        this.$refs.lokiFilter.disabled = false
      }
    },
    handlerLiveQuerying() {
      if (!this.loading) {
        this.handlerInitWebSocket()
      } else {
        this.handlerClose()
      }
    },
    async listLabels(data) {
      try {
        const res = await listLabels(data)
        if (res.status === 200) {
          this.labels = res.data
        }
      } catch (err) {
        this.$store.commit('showSnackBar', {
          text: 'Error: 获取数据失败',
          color: 'error',
        })
      }
    },
  },
  created() {
    window.addEventListener('beforeunload', this.handlerClose)
    window.addEventListener('unload', this.handlerClose)
  },
  destroyed() {
    this.handlerClose()
    window.removeEventListener('beforeunload', this.handlerClose)
    window.removeEventListener('unload', this.handlerClose)
  },
  async mounted() {
    if (this.$store.state.jwt) {
      this.$refs.dateRangePicker.handlerChangeQuickTime(5)
      await this.listLabels({})
      const pod = this.$route.query['pod']
      const start = this.$route.query['start']
      const end = this.$route.query['end']
      let filters = this.$route.query['filters']
      for (var key in this.$route.query) {
        if (this.labels.indexOf(key) > -1) {
          this.$refs.lokiFilter.model.push({
            text: `标签(${key}):${this.$route.query[key]}`,
            value: key,
          })
        }
      }
      if (filters !== undefined) {
        if (typeof filters === 'string') {
          filters = [filters]
        }
        filters.forEach((item) => {
          this.$refs.lokiFilter.model.push({
            text: '正则(regex):' + item,
            value: 'filter',
          })
        })
      }
      if (this.$refs.lokiFilter.model.length > 0) {
        if (start !== undefined && end !== undefined) {
          const timeReg = /^[0-9][0-9][0-9][0-9]-[0-1][0-9]-[0-3][0-9] [0-2][0-9]:[0-5][0-9]:[0-5][0-9]$/g
          if (timeReg.test(start) && timeReg.test(end)) {
            this.$refs.dateRangePicker.dateRange = [start, end]
            this.dateRangeTimestamp = [
              Date.parse(parserDatetime(start)).toString() + '000000',
              Date.parse(parserDatetime(end)).toString() + '000000',
            ]
            this.$refs.dateRangePicker.quick = false
          } else {
            this.$store.commit('showSnackBar', {
              text: 'Warning: 参数start，end格式错误',
              color: 'warning',
            })
            return
          }
        }
        if (pod !== undefined) {
          this.pod = pod
          this.pods = [{ text: pod, selected: true }]
        }
        this.listQueryRanges()
      }
    }
  },
}
</script>
