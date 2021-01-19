<template>
  <v-flex style="padding: 0 !important;">
    <v-chart
      ref="lokiChart"
      :options="options"
      :initOptions="initOptions"
      @datazoom="handlerDataZoom"
      autoresize
      renderer="svg"
    ></v-chart>
  </v-flex>
</template>
<style scoped>
.echarts {
  width: 100%;
  height: 180px;
  padding-top: 5px;
  padding: 5px 20px 0 !important;
}
</style>

<script>
import ECharts from 'vue-echarts/components/ECharts'
import 'echarts/lib/chart/bar'
import 'echarts/lib/component/tooltip'
import 'echarts/lib/component/dataZoom'
import 'echarts/lib/component/legendScroll'

export default {
  name: 'LokiViewerHistogram',
  components: {
    'v-chart': ECharts,
  },
  props: {
    middleStart: {
      type: String,
      default: () => '',
    },
    middleEnd: {
      type: String,
      default: () => '',
    },
  },
  data: () => ({
    chartData: {
      'xAxis-data': [],
      'yAxis-data': {
        info: [],
        debug: [],
        error: [],
        warn: [],
        unknown: [],
      },
    },
    loading: false,
    date: new Date().toISOString().substr(0, 10),
    menu: false,
    start: 0,
    end: 100,
    legendSelected: {
      Info: true,
      Debug: true,
      Warn: true,
      Error: true,
      Unknown: true,
    },
  }),
  computed: {
    initOptions() {
      return { width: 'auto', height: 'auto' }
    },
    options() {
      return {
        color: ['#4caf50', '#5cbbf6', '#fb8c00', '#ff5252', '#607d8b'],
        tooltip: {
          trigger: 'axis',
          axisPointer: {
            type: 'shadow',
          },
        },
        calculable: true,
        xAxis: [
          {
            type: 'category',
            axisTick: { show: false },
            data: this.chartData['xAxis-data'],
          },
        ],
        yAxis: [
          {
            type: 'value',
            minInterval: 1,
          },
        ],
        dataZoom: [
          {
            type: 'slider',
            show: true,
            start: this.start,
            end: this.end,
            realtime: false,
          },
        ],
        grid: {
          x: 30,
          y: 10,
          x2: 30,
          y2: 60,
        },
        series: [
          {
            name: 'Info',
            type: 'bar',
            data: this.chartData['yAxis-data'].info,
          },
          {
            name: 'Debug',
            type: 'bar',
            data: this.chartData['yAxis-data'].debug,
          },
          {
            name: 'Warn',
            type: 'bar',
            data: this.chartData['yAxis-data'].warn,
          },
          {
            name: 'Error',
            type: 'bar',
            data: this.chartData['yAxis-data'].error,
          },
          {
            name: 'Unknown',
            type: 'bar',
            data: this.chartData['yAxis-data'].unknown,
          },
        ],
      }
    },
  },
  methods: {
    handerReset() {
      this.start = 0
      this.end = 100
      this.legendSelected = {
        Info: true,
        Debug: true,
        Warn: true,
        Error: true,
        Unknown: true,
      }
    },
    handlerDataZoom(event) {
      const splitPercent = parseInt(100 / (event.end - event.start))
      const size = splitPercent * 20
      this.start = event.start
      this.end = event.end
      this.$parent.$parent.size = size
      const percent =
        (parseInt(this.middleEnd) - parseInt(this.middleStart)) / 100
      const s = (event.start * percent + parseInt(this.middleStart)).toString()
      const e = (
        parseInt(this.middleEnd) -
        (100 - event.end) * percent
      ).toString()
      this.$parent.$parent.middleStart = s
      this.$parent.$parent.middleEnd = e
      this.$emit('refresh')
    },
  },
}
</script>
