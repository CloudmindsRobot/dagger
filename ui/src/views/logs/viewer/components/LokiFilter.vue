<template>
  <v-combobox
    :items="filterItems"
    @change="handlerFilterList"
    @blur="handlerClearFilterList"
    @focus="handlerClearFilterList"
    :search-input.sync="filter"
    color="primary"
    chips
    hide-selected
    label="条件过滤"
    multiple
    prepend-icon="filter_list"
    v-model="model"
    dense
    no-data-text="无数据"
    full-width
    ref="filter"
    :loading="loading"
    :hide-no-data="!filter"
    :disabled="disabled"
  >
    <template v-slot:selection="{ attrs, item, parent, selected }">
      <v-chip
        :input-value="selected"
        color="primary"
        label
        small
        v-bind="attrs"
        v-if="item === Object(item)"
      >
        <span class="pr-2">{{ item.text }}</span>
        <v-icon @click="handlerFilterList(item, 'close')" small>close</v-icon>
      </v-chip>
    </template>
    <template v-slot:item="{ index, item }">
      <span v-if="label">{{ item.text }}</span>
      <span v-else>
        <v-chip
          color="primary"
          label
          small
          v-for="(v, k) in item.text.label_json"
          :key="k"
          style="margin:0 5px;"
        >
          <span class="pr-2">标签({{ k }}):{{ v }}</span>
        </v-chip>
        <v-chip
          color="primary"
          label
          small
          v-for="it in item.text.filter_json"
          :key="it"
          style="margin:0 5px;"
        >
          <span class="pr-2">正则(regex):{{ it }}</span>
        </v-chip>
      </span>
    </template>
    <template v-slot:no-data>
      <v-list-item @click="handlerRegexFilterList">
        <span class="subheading">正则(regex):{{ filter }}</span>
      </v-list-item>
    </template>
    <template v-slot:append>
      <v-btn
        text
        color="primary"
        small
        @click="handlerFilterList(null, 'quick_query')"
      >
        快速查询
      </v-btn>
      <v-btn
        color="primary"
        text
        small
        :loading="saveQueryLoading"
        @click="handlerFilterList(null, 'save_query')"
      >
        保存条件
      </v-btn>
    </template>
  </v-combobox>
</template>

<script>
import {
  listLabelValues,
  listQueryHistory,
  listLabels,
  createQueryHistoryLabel,
} from '@/api'

export default {
  name: 'LokiViewerFilterCombobox',
  props: {
    dateRangeTimestamp: {
      type: Array,
      default: () => [],
    },
  },
  data: () => ({
    model: [],
    filterItems: [],
    originItems: [],
    filterDict: {
      label: {
        value: '',
        text: '',
        items: [],
      },
    },
    loading: false,
    saveQueryLoading: false,
    filter: '',
    filters: [],
    disabled: false,
    label: true,
    querys: [],
  }),
  methods: {
    async listLabelValues(data) {
      try {
        this.loading = true
        const res = await listLabelValues(data)
        if (res.status === 200) {
          this.filterDict.label.value = 'label'
          this.filterDict.label.text = `标签(${data.label})`
          const items = []
          res.data.forEach((item) => {
            items.push({ text: item, value: item, sec: true, parent: 'label' })
          })
          this.filterDict.label.items = items
        }
      } catch (err) {
        this.$store.commit('showSnackBar', {
          text: 'Error: 获取数据失败',
          color: 'error',
        })
      }
      this.loading = false
    },
    async listLabels(data) {
      try {
        this.loading = true
        const res = await listLabels(data)
        if (res.status === 200) {
          const items = []
          res.data.forEach((item) => {
            items.push({ text: `标签(${item})`, value: item })
          })
          this.filterItems = JSON.parse(JSON.stringify(items))
          this.originItems = JSON.parse(JSON.stringify(items))
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
    handlerClearFilterList() {
      this.label = true
      this.filterItems = []
      this.originItems.forEach((item) => {
        this.filterItems.push({
          text: item.text.split(':') > 1 ? item.text.split(':')[0] : item.text,
          value: item.value,
        })
      })
      this.model = this.model.filter((item) => item.text.indexOf(':') > -1)
    },
    handlerRegexFilterList() {
      this.model.push({ text: '正则(regex):' + this.filter, value: 'filter' })
      this.filters.push(this.filter)
      this.handlerSearch()
    },
    handlerSearch() {
      this.$refs.filter.blur()
    },
    async handleQuickqueryList() {
      try {
        const res = await listQueryHistory({ page: 1, page_size: 500 })
        if (res.status === 200 && res.data.success) {
          this.querys = res.data.data
          this.querys.forEach((item) => {
            item.label_json = JSON.parse(item.label_json)
            item.filter_json = JSON.parse(item.filter_json)
          })
        } else {
          this.$store.commit('showSnackBar', {
            text: 'Error: 获取历史查询标签失败',
            color: 'error',
          })
        }
      } catch (err) {
        if (err.response && err.response.status !== 401) {
          this.$store.commit('showSnackBar', {
            text: 'Error: 获取历史查询标签失败',
            color: 'error',
          })
        }
      }
    },
    async handlerSaveQuerying() {
      if (this.model.length === 0) {
        this.$store.commit('showSnackBar', {
          text: 'Warn: 请输入查询',
          color: 'warning',
        })
        return
      }
      if (this.saveQueryLoading) return
      this.saveQueryLoading = true
      const filters = []
      const labels = {}
      this.model.forEach((item) => {
        if (item.value === 'filter') {
          filters.push(item.text.substr(item.text.indexOf(':') + 1))
        } else {
          labels[item.value] = item.text.substr(item.text.indexOf(':') + 1)
        }
      })
      try {
        const res = await createQueryHistoryLabel({
          label_json: JSON.stringify(labels),
          filter_json: JSON.stringify(filters),
        })
        if (res.status === 201) {
          this.$store.commit('showSnackBar', {
            text: 'Success: 保存查询条件成功',
            color: 'success',
          })
        } else {
          this.$store.commit('showSnackBar', {
            text: `Error: ${res.data.message}`,
            color: 'error',
          })
        }
      } catch (err) {
        this.$store.commit('showSnackBar', {
          text: 'Error: 保存查询条件失败',
          color: 'error',
        })
      }
      this.saveQueryLoading = false
    },
    async handlerFilterList(im, op) {
      if (op === 'close') {
        this.model = this.model.filter(
          (item) => item.value !== im.value || item.text !== im.text,
        )
        this.filters = []
        this.model.forEach((item) => {
          if (item.value === 'filter') {
            this.filters.push(item.text.split(':')[1])
          }
        })
        this.handlerSearch()
      } else if (op === 'quick_query') {
        this.label = false
        await this.handleQuickqueryList()
        this.filterItems = []
        this.querys.forEach((item, index) => {
          this.filterItems.push({
            text: item,
            value: `quick${index}`,
          })
        })
      } else if (op === 'save_query') {
        this.$refs.filter.blur()
        await this.handlerSaveQuerying()
      } else {
        let label = null
        im.forEach((item) => {
          if (
            item.hasOwnProperty('value') &&
            item.value.indexOf('quick') > -1
          ) {
            label = item
          }
        })
        if (label !== null) {
          this.model = []
          for (var item in label.text.label_json) {
            this.model.push({
              text: `标签(${item}):${label.text.label_json[item]}`,
              value: item,
            })
          }
          label.text.filter_json.forEach((item) => {
            this.model.push({ text: `正则(regex):${item}`, value: 'filter' })
          })
          this.$refs.filter.blur()
        } else {
          this.label = true
          if (this.model.length === 0) return
          this.model = this.model.filter((item) => typeof item !== 'string')
          var lastModel = this.model[this.model.length - 1]
          this.filterItems = []
          if (lastModel.text.indexOf(':') > -1) {
            this.$refs.filter.blur()
            return
          }
          if (lastModel.sec) {
            const lastSecModel = this.model[this.model.length - 2]
            if (lastSecModel.text.indexOf(':') === -1) {
              lastSecModel.text = lastSecModel.text + ':' + lastModel.value
            }
            this.model = this.model.filter((item) => !item.sec)

            this.handlerSearch()
          } else {
            await this.listLabelValues({
              label: lastModel.value,
              start: this.dateRangeTimestamp[0],
              end: this.dateRangeTimestamp[1],
            })
            this.filterItems = []
            this.filterDict.label.items.forEach((i) => {
              this.filterItems.push(i)
            })
          }
        }
      }
    },
  },
  mounted() {
    if (this.$store.state.jwt) {
      this.$nextTick(() => {
        this.listLabels({
          start: this.dateRangeTimestamp[0],
          end: this.dateRangeTimestamp[1],
        })
      })
    }
  },
}
</script>
