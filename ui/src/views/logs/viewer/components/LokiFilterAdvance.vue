<template>
  <v-flex style="padding: 0;">
    <v-textarea
      v-model="logQL"
      dense
      full-width
      label="高级查询"
      auto-grow
      single-line
      clearable
      rows="1"
      prepend-icon="filter_list"
      @focus="handlerFocus"
      @blur="handlerBlur"
      @keyup="handlerInput"
      @click:clear="handlerClear"
      ref="advanceFilter"
    >
      <template v-slot:append v-if="settings.alertEnabled">
        <v-btn text small color="primary" @click.stop="handlerOpenDeriveDialog">
          派生规则
        </v-btn>
      </template>
    </v-textarea>
    <v-data-table
      v-if="suggestShow"
      :headers="headers"
      :items="filterItems"
      :loading="loading"
      :search="keyword"
      :custom-filter="filterKeyword"
      :style="'top:' + suggestTop + 'px;'"
      loading-text="载入中..."
      item-key="item.value"
      disable-pagination
      hide-default-header
      hide-default-footer
      no-data-text="无提示"
      no-results-text="无匹配"
      class="suggestion-table"
      @click:row="handlerSelect"
      ref="suggestion"
    >
      <template v-slot:item.item="{ item }">
        {{ item.value }}
      </template>
    </v-data-table>
    <DeriveRule
      v-if="settings.alertEnabled"
      :logQL.sync="logQL"
      ref="deriveRule"
    />
  </v-flex>
</template>

<script>
import { listLabels, listLabelAllValues } from '@/api'
import DeriveRule from '@/views/logs/alert/components/DeriveRule'
import { mapState } from 'vuex'

export default {
  name: 'FilterAdvanced',
  components: {
    DeriveRule,
  },
  props: {
    queryed: {
      type: Boolean,
      default: () => false,
    },
    resultType: {
      type: String,
      default: () => '',
    },
  },
  computed: {
    ...mapState(['settings']),
  },
  data: () => ({
    suggestShow: false,
    model: [],
    filterItems: [],
    originRangeFunctions: [
      { text: 'rate', value: 'rate', t: 'normal' },
      { text: 'count_over_time', value: 'count_over_time', t: 'time' },
      { text: 'sum_over_time', value: 'sum_over_time', t: 'time' },
      { text: 'avg_over_time', value: 'avg_over_time', t: 'time' },
      { text: 'max_over_time', value: 'max_over_time', t: 'time' },
      { text: 'min_over_time', value: 'min_over_time', t: 'time' },
      { text: 'topk', value: 'topk', t: 'top' },
      { text: 'sum', value: 'sum', t: 'normal' },
      { text: 'min', value: 'min', t: 'normal' },
      { text: 'max', value: 'max', t: 'normal' },
      { text: 'avg', value: 'avg', t: 'normal' },
      { text: 'stddev', value: 'stddev', t: 'normal' },
      { text: 'stdvar', value: 'stdvar', t: 'normal' },
      { text: 'count', value: 'count', t: 'normal' },
      { text: 'bottomk', value: 'bottomk', t: 'top' },
      { text: 'stdvar_over_time', value: 'stdvar_over_time', t: 'time' },
      { text: 'stddev_over_time', value: 'stddev_over_time', t: 'time' },
      { text: 'quantile_over_time', value: 'quantile_over_time', t: 'time' },
    ],
    originLabels: [],
    originValues: [],
    originPiplines: [
      { text: 'json', value: 'json', t: 'none' },
      { text: 'line_format', value: 'line_format', t: 'quote' },
      { text: 'logfmt', value: 'logfmt', t: 'none' },
      { text: 'regexp', value: 'regexp', t: 'quote' },
      { text: 'label_format', value: 'label_format', t: 'space' },
      { text: 'unwrap', value: 'unwrap', t: 'space' },
    ],
    headers: [{ text: '', value: 'item', align: 'start' }],
    loading: false,
    tmpFilters: [],
    logQL: '',
    keyword: '',
    position: 0,
    label: '',
    suggestTop: 62,
  }),
  mounted() {
    this.$refs.advanceFilter.$refs.input.style.lineHeight = '40px'
  },
  methods: {
    async handlerLabels() {
      try {
        this.loading = true
        const res = await listLabels({})
        if (res.status === 200) {
          this.originLabels = []
          res.data.data.forEach((item) => {
            if (item.indexOf('__') === -1) {
              this.originLabels.push({
                text: item,
                value: item,
                t: 'label',
              })
            }
          })
        }
      } catch (err) {
        if (
          err.response &&
          [400, 401, 403, 504].indexOf(err.response.status) === -1
        ) {
          this.$store.commit('showSnackBar', {
            text: 'Error: 获取数据失败',
            color: 'error',
          })
        }
      }
      this.loading = false
    },
    async handlerLabelValues(label) {
      try {
        this.loading = true
        const res = await listLabelAllValues(label)
        if (res.status === 200) {
          this.originValues = []
          res.data.data.forEach((item) => {
            this.originValues.push({
              text: item,
              value: item,
              t: 'val',
            })
          })
        }
      } catch (err) {
        if (
          err.response &&
          [400, 401, 403, 504].indexOf(err.response.status) === -1
        ) {
          this.$store.commit('showSnackBar', {
            text: 'Error: 获取数据失败',
            color: 'error',
          })
        }
      }
      this.loading = false
    },
    handlerClear() {
      this.logQL === ''
      this.keyword = ''
      this.suggestShow = false
      this.suggestTop = 62
      this.$refs.advanceFilter.blur()
      this.filterItems = this.originRangeFunctions
      this.tmpFilters = this.originRangeFunctions
    },
    handlerFocus() {
      const vue = this
      setTimeout(() => {
        vue.suggestShow = true
      }, 350)
      if (this.logQL === '') {
        this.keyword = ''
        this.filterItems = this.originRangeFunctions
        this.tmpFilters = this.originRangeFunctions
      }
    },
    handlerBlur(e) {
      this.position = e.srcElement.selectionStart
      const vue = this
      setTimeout(() => {
        vue.suggestShow = false
      }, 300)
      if (
        this.logQL.indexOf('}') < this.logQL.lastIndexOf('|') &&
        this.logQL.lastIndexOf('|') > -1
      )
        this.$emit('updateAdvanceFilter', true)
      else this.$emit('updateAdvanceFilter', false)
    },
    handlerOpenDeriveDialog() {
      this.suggestShow = false
      if (this.logQL === '') {
        this.$store.commit('showSnackBar', {
          text: 'Warn: 请输入查询条件',
          color: 'warning',
        })
      } else if (this.queryed === false) {
        this.$store.commit('showSnackBar', {
          text: 'Warn: 请查询以验证LogQL语句',
          color: 'warning',
        })
      } else if (this.resultType !== 'matrix') {
        this.$store.commit('showSnackBar', {
          text: 'Warn: 请输入正确的Metric LogQL语句',
          color: 'warning',
        })
      } else {
        this.$refs.deriveRule.obj.log_ql = this.logQL
        this.$refs.deriveRule.dialog = true
      }
      this.$refs.advanceFilter.blur()
    },
    handlerSelect(item) {
      this.$emit('updateQueryed', false)
      this.logQL = `${this.logQL.substr(
        0,
        this.position - this.keyword.length,
      )}${this.logQL.substr(this.position)}`
      this.position = this.position - this.keyword.length
      switch (item.t) {
        case 'normal':
          this.logQL = `${item.value}(${this.logQL})`
          break
        case 'top':
          this.logQL = `${item.value}(10, ${this.logQL})`
          break
        case 'time':
          this.logQL = `${item.value}(${this.logQL}[5m])`
          break
        case 'none':
          this.logQL = `${this.logQL.substr(0, this.position)}${
            item.value
          }${this.logQL.substr(this.position)}`
          break
        case 'space':
          this.logQL = `${this.logQL.substr(0, this.position)}${
            item.value
          } ${this.logQL.substr(this.position)}`
          break
        case 'quote':
          this.logQL = `${this.logQL.substr(0, this.position)}${
            item.value
          } \`${this.logQL.substr(this.position)}`
          break
        case 'label':
          this.logQL = `${this.logQL.substr(0, this.position)}${
            item.value
          }=${this.logQL.substr(this.position)}`
          this.label = item.value
          break
        case 'val':
          this.logQL = `${this.logQL.substr(0, this.position)}${
            item.value
          }"${this.logQL.substr(this.position)}`
          break
      }
      this.keyword = ''
      this.suggestShow = false
      this.filterItems = []
      this.$refs.advanceFilter.focus()
    },
    async handlerInput(e) {
      this.suggestTop = this.$refs.advanceFilter.$el.clientHeight + 4
      this.$emit('updateQueryed', false)
      const p = e.srcElement.selectionStart
      if (e.keyCode > 32 && e.keyCode <= 200) {
        this.keyword += e.key
      }
      if (this.logQL === '') {
        this.keyword = ''
        this.filterItems = []
        this.filterItems = this.originRangeFunctions
        this.tmpFilters = this.originRangeFunctions
      } else {
        let key = e.key
        if (e.key === 'Backspace') {
          key = this.logQL[p - 1]
          this.keyword = this.keyword.substr(0, this.keyword.length - 1)
        }
        switch (key) {
          case '{':
          case ',':
            if (
              (this.logQL.indexOf('{') < p &&
                this.logQL.indexOf('}') >= p &&
                this.logQL[p - 1] !== '{') ||
              this.logQL === '{' ||
              (this.logQL.indexOf('{') < p && this.logQL.indexOf('}') === -1)
            ) {
              // label show
              this.keyword = ''
              await this.handlerLabels()
              this.filterItems = this.originLabels
              this.tmpFilters = this.originLabels
            }
            break
          case '"':
            if (
              this.logQL.indexOf('}') === -1 ||
              (this.logQL.indexOf('{') < p && this.logQL.indexOf('}') > p) ||
              this.logQL
                .substr(this.logQL.indexOf('{'), this.logQL.indexOf('}'))
                .match(new RegExp('"', 'g')).length %
                2 ===
                1
            ) {
              // label value
              this.keyword = ''
              await this.handlerLabelValues(this.label)
              this.filterItems = this.originValues
              this.tmpFilters = this.originValues
            }
            break
          case '|':
            // pipelinex
            this.keyword = ''
            this.filterItems = this.originPiplines
            this.tmpFilters = this.originPiplines
            break
          default:
            if (
              this.logQL.indexOf('{') > -1 &&
              this.logQL.indexOf('{') === p - 1 &&
              this.logQL.length > p
            ) {
              this.keyword = ''
              await this.handlerLabels()
              this.filterItems = this.originLabels
              this.tmpFilters = this.originLabels
            } else if (
              (this.logQL.indexOf('}') > -1 &&
                this.logQL.indexOf('}') === p - 1) ||
              p === 0
            ) {
              this.keyword = ''
              this.filterItems = this.originRangeFunctions
              this.tmpFilters = this.originRangeFunctions
            }
        }
        this.filterItems = this.tmpFilters.filter((item) => {
          return item.value.indexOf(this.keyword) > -1
        })
      }
    },
    filterKeyword(value, search, item) {
      return (
        typeof item.value === 'string' &&
        item.value.toString().indexOf(search) !== -1
      )
    },
  },
}
</script>
<style scoped>
.suggestion-table {
  min-width: 200px !important;
  box-shadow: 2px 3px 3px gray !important;
  position: absolute;
  z-index: 1;
  left: 60px;
  max-height: 300px;
  overflow-y: auto;
  font-weight: normal;
}
</style>
