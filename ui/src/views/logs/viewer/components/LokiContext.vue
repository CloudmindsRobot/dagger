<template>
  <div>
    <v-btn text x-small color="primary" @click.stop="handlerShowLokiContext"
      >显示上下文</v-btn
    >
    <v-dialog v-model="dialog" max-width="90%">
      <v-card>
        <v-card-text>
          <v-btn
            style="margin-top: 15px;"
            text
            x-small
            :loading="loadingPreview"
            color="primary"
            @click="handlerLoadingPreview"
            >向前加载10条数据</v-btn
          >
          <v-data-table
            :headers="headers"
            :items="itemsPreview"
            loading-text="载入中..."
            item-key="info.timestamp"
            disable-pagination
            hide-default-header
            hide-default-footer
            no-data-text="无数据"
            dense
          >
            <template v-slot:item.info="{ item }">
              <div
                :class="'v-list-item-html v-list-item-' + item.level"
                v-html="item.message"
              ></div>
            </template>
          </v-data-table>
          <v-data-table
            :headers="headers"
            :items="items"
            item-key="info.timestamp"
            disable-pagination
            hide-default-header
            hide-default-footer
            no-data-text="无数据"
            dense
          >
            <template v-slot:item.info="{ item }">
              <div
                :class="'v-list-item-html v-list-item-' + item.level"
                v-html="item.message"
                style="color: #1976d2 !important;"
              ></div>
            </template>
          </v-data-table>
          <v-data-table
            :headers="headers"
            :items="itemsNext"
            loading-text="载入中..."
            item-key="info.timestamp"
            disable-pagination
            hide-default-header
            hide-default-footer
            no-data-text="无数据"
            dense
          >
            <template v-slot:item.info="{ item }">
              <div
                :class="'v-list-item-html v-list-item-' + item.level"
                v-html="item.message"
              ></div>
            </template>
          </v-data-table>
          <v-btn
            text
            x-small
            :loading="loadingNext"
            color="primary"
            @click="handlerLoadingNext"
            >向后加载10条数据</v-btn
          >
        </v-card-text>
      </v-card>
    </v-dialog>
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
.v-list--dense .v-list-item,
.v-list-item--dense {
  min-height: 0px !important;
}
</style>

<script>
import { listContext } from '@/api'

export default {
  name: 'LokiViewerContext',
  props: {
    timestamp: {
      type: String,
      default: () => '',
    },
    loki: {
      type: Object,
      default: () => {},
    },
    logQL: {
      type: String,
      default: () => '',
    },
  },
  data: () => ({
    show: false,
    dialog: false,
    valid: false,
    loadingPreview: false,
    loadingNext: false,
    itemsPreview: [],
    itemsNext: [],
    items: [],
    headers: [{ text: '', value: 'info', align: 'start' }],
  }),
  methods: {
    async listContext(data) {
      try {
        const res = await listContext(data)
        if (res.status === 200) {
          if (data.direction === 'preview') {
            res.data.sort((a, b) => {
              return a.timestamp - b.timestamp
            })
            if (data.append) {
              res.data.reverse()
              this.itemsPreview.reverse()
              res.data.forEach((item) => {
                this.itemsPreview.push(item)
              })
              this.itemsPreview.reverse()
            } else {
              this.itemsPreview = res.data
            }
          } else if (data.direction === 'next') {
            res.data.sort((a, b) => {
              return a.timestamp - b.timestamp
            })
            if (data.append) {
              res.data.forEach((item) => {
                this.itemsNext.push(item)
              })
            } else {
              this.itemsNext = res.data
            }
          }
        } else {
          this.$store.commit('showSnackBar', {
            text: `Warn: ${res.data.message}`,
            color: 'warning',
          })
        }
      } catch (err) {
        this.$store.commit('showSnackBar', {
          text: 'Error: 获取数据失败',
          color: 'error',
        })
      }
    },
    async handlerShowLokiContext() {
      this.handlerLoadingPreview()
      this.handlerLoadingNext()
      this.dialog = true
    },
    async handlerLoadingPreview() {
      this.loadingPreview = true
      let timestamp = (' ' + this.timestamp).slice(1)
      if (this.itemsPreview.length > 0)
        timestamp = (' ' + this.itemsPreview[0].timestamp).slice(1)

      const previewStart =
        Date.parse(
          new Date(
            new Date(parseInt(timestamp.substr(0, 13))).setHours(
              new Date(parseInt(timestamp.substr(0, 13))).getHours() - 3,
            ),
          ),
        ) + '000000'
      const previewEnd = (parseInt(timestamp) - 100000).toString()
      await this.listContext({
        logql: this.logQL,
        direction: 'preview',
        start: previewStart,
        end: previewEnd,
        append: true,
      })
      this.loadingPreview = false
    },
    async handlerLoadingNext() {
      this.loadingNext = true
      let timestamp = (' ' + this.timestamp).slice(1)
      if (this.itemsNext.length > 0)
        timestamp = (
          ' ' + this.itemsNext[this.itemsNext.length - 1].timestamp
        ).slice(1)
      const nextStart = (parseInt(timestamp) + 100000).toString()
      const nextEnd =
        Date.parse(
          new Date(
            new Date(parseInt(timestamp.substr(0, 13))).setHours(
              new Date(parseInt(timestamp.substr(0, 13))).getHours() + 3,
            ),
          ),
        ) + '000000'
      await this.listContext({
        logql: this.logQL,
        direction: 'next',
        start: nextStart,
        end: nextEnd,
        append: true,
      })
      this.loadingNext = false
    },
  },
  mounted() {
    this.items = [
      {
        message: this.loki.info.message,
        timestamp: this.loki.info.timestamp,
        level: this.loki.info.level,
      },
    ]
  },
}
</script>
