<template>
  <v-menu
    v-model="menu"
    :close-on-content-click="false"
    :nudge-right="40"
    transition="scale-transition"
    min-width="290px"
    offset-y
  >
    <template v-slot:activator="{ on }">
      <v-text-field
        label="时间选择"
        v-model="dateRangeStr"
        prepend-icon="event"
        readonly
        v-on="on"
        style="width: 270px;"
      ></v-text-field>
    </template>
    <v-card>
      <v-card-text>
        <v-row no-gutters>
          <v-col cols="2">
            <v-list dense>
              <v-subheader>最近时间</v-subheader>
              <v-list-item-group color="primary" v-model="quickSelect">
                <v-list-item
                  @click="handlerChangeQuickTime(item.value)"
                  v-for="(item, i) in quickTimeItems"
                  :key="i"
                >
                  <v-list-item-content>
                    <v-list-item-title v-text="item.text"></v-list-item-title>
                  </v-list-item-content>
                </v-list-item>
              </v-list-item-group>
            </v-list>
          </v-col>
          <v-col cols="10">
            <v-text-field
              v-model="tRangeStr"
              class="date-select"
              label="起止时间(示例09:00:00,11:00:00)"
            ></v-text-field>
            <v-row no-gutters>
              <v-col cols="6">
                <div>起始日期</div>
                <v-date-picker
                  no-title
                  locale="zh-cn"
                  v-model="dStart"
                  @change="handlerChangeDate"
                  style="min-height: 290px !important;margin-right: 5px;"
                ></v-date-picker>
              </v-col>
              <v-col cols="6">
                <div>终止日期</div>
                <v-date-picker
                  no-title
                  locale="zh-cn"
                  v-model="dEnd"
                  @change="handlerChangeDate"
                  style="min-height: 290px !important;"
                ></v-date-picker>
              </v-col>
            </v-row>
          </v-col>
        </v-row>
      </v-card-text>
      <v-card-actions>
        <v-layout row justify-space-around>
          <v-flex xs5>
            <v-btn color="primary" block @click="handlerCheckDatePicker"
              >确定</v-btn
            >
          </v-flex>
          <v-flex xs5>
            <v-btn block @click="handlerCancelDatePicker">取消</v-btn>
          </v-flex>
        </v-layout>
      </v-card-actions>
    </v-card>
  </v-menu>
</template>

<style scoped>
.date-select {
  padding-left: 10px;
  padding-right: 10px;
  padding-top: 10px;
}
</style>

<script>
import { formatDatetime, parserDatetime } from '@/utils/helpers'

export default {
  name: 'LokiViewerDatetimeRangePicker',
  components: {},
  data: () => ({
    lastDateRange: [
      formatDatetime(
        new Date(new Date().setMinutes(new Date().getMinutes() - 5)),
        'yyyy-MM-dd hh:mm:ss',
      ),
      formatDatetime(new Date(), 'yyyy-MM-dd hh:mm:ss'),
    ],
    dateRange: [
      formatDatetime(
        new Date(new Date().setMinutes(new Date().getMinutes() - 5)),
        'yyyy-MM-dd hh:mm:ss',
      ),
      formatDatetime(new Date(), 'yyyy-MM-dd hh:mm:ss'),
    ],
    dStart: '',
    dEnd: '',
    tRangeStr: '',
    menu: false,
    quickTime: 5,
    quickSelect: 0,
    quickTimeItems: [
      { text: '最近5分钟', value: 5 },
      { text: '最近15分钟', value: 15 },
      { text: '最近30分钟', value: 30 },
      { text: '最近1小时', value: 60 },
      { text: '最近3小时', value: 180 },
      { text: '最近6小时', value: 360 },
      { text: '最近12小时', value: 720 },
      { text: '最近24小时', value: 1440 },
    ],
    quick: true,
    dateStart: '',
    dateEnd: '',
    timeStart: '',
    timeEnd: '',
  }),
  computed: {
    dateRangeStr() {
      if (this.quick) {
        return this.quickTimeItems[this.quickSelect].text
      }
      if (this.dateRange[0] !== undefined && this.dateRange[1] !== undefined)
        return (
          formatDatetime(parserDatetime(this.dateRange[0]), 'MM-dd hh:mm:ss') +
          '～' +
          formatDatetime(parserDatetime(this.dateRange[1]), 'MM-dd hh:mm:ss')
        )
      else return ''
    },
    tRange() {
      if (
        /^[0-2][0-9]:[0-5][0-9]:[0-5][0-9],[0-2][0-9]:[0-5][0-9]:[0-5][0-9]$/g.test(
          this.tRangeStr.trim(),
        )
      ) {
        return this.tRangeStr.split(',')
      }
      return []
    },
  },
  methods: {
    handlerChangeQuickTime(val) {
      if (val !== undefined) this.quickTime = val
      const end = new Date()
      const start = new Date(
        new Date().setMinutes(new Date().getMinutes() - this.quickTime),
      )

      this.dateRange = [
        formatDatetime(start, 'yyyy-MM-dd hh:mm:ss'),
        formatDatetime(end, 'yyyy-MM-dd hh:mm:ss'),
      ]
      this.dStart = formatDatetime(start, 'yyyy-MM-dd')
      this.dEnd = formatDatetime(end, 'yyyy-MM-dd')
      this.tRangeStr =
        formatDatetime(start, 'hh:mm:ss') +
        ',' +
        formatDatetime(end, 'hh:mm:ss')
      this.$parent.$parent.dateRangeTimestamp = [
        Date.parse(
          new Date(
            new Date().setMinutes(new Date().getMinutes() - this.quickTime),
          ),
        ).toString() + '000000',
        Date.parse(new Date()).toString() + '000000',
      ]
      this.lastDateRange = this.dateRange
      this.quick = true
      this.menu = false
    },
    handlerChangeDate() {
      const dateStart = formatDatetime(
        parserDatetime(this.dStart),
        'yyyy-MM-dd 00:00:00',
      )
      const dateEnd = formatDatetime(
        parserDatetime(this.dEnd),
        'yyyy-MM-dd 00:00:00',
      )
      this.dateRange[0] = dateStart
      this.dateRange[1] = dateEnd
      this.$parent.$parent.dateRangeTimestamp[0] =
        Date.parse(parserDatetime(this.dStart)) + '000000'
      this.$parent.$parent.dateRangeTimestamp[1] =
        Date.parse(parserDatetime(this.dEnd)) + '000000'
    },
    handlerCheckDatePicker() {
      if (
        /^[0-2][0-9]:[0-5][0-9]:[0-5][0-9],[0-2][0-9]:[0-5][0-9]:[0-5][0-9]$/g.test(
          this.tRangeStr.trim(),
        ) === false
      ) {
        this.$store.commit('showSnackBar', {
          text: 'Warning: 时间格式有误',
          color: 'warning',
        })
        return
      }

      const dateStart = this.dateRange[0].substr(0, 10) + ' ' + this.tRange[0]
      const dateEnd = this.dateRange[1].substr(0, 10) + ' ' + this.tRange[1]
      this.dateRange = [dateStart, dateEnd]
      this.$parent.$parent.dateRangeTimestamp = [
        Date.parse(parserDatetime(dateStart)) + '000000',
        Date.parse(parserDatetime(dateEnd)) + '000000',
      ]

      if (
        Date.parse(parserDatetime(this.dateRange[1])) <=
        Date.parse(parserDatetime(this.dateRange[0]))
      ) {
        this.$store.commit('showSnackBar', {
          text: 'Warning: 截止时间超过起始时间，请重新选择',
          color: 'warning',
        })
        return
      }
      this.$set(this.dateRange, 0, this.dateRange[0])
      this.$set(this.dateRange, 1, this.dateRange[1])
      this.lastDateRange = this.dateRange
      this.quickSelect = null
      this.quick = false
      this.menu = false
    },
    handlerCancelDatePicker() {
      this.dateRange = this.lastDateRange
      this.menu = false
    },
  },
  mounted() {},
}
</script>
