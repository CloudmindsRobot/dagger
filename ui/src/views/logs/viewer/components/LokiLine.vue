<template>
  <v-flex style="padding: 0 !important;">
    <v-chart
      ref="lokiChart"
      :options="options"
      :initOptions="initOptions"
      autoresize
    ></v-chart>
    <div style="padding: 0px 10px;">
      <div
        class="legend"
        :style="'border-left: 3px solid ' + color[index % 11]"
        v-for="(val, index) in chartData['table-data']"
        :key="index"
      >
        <span class="legend-prefix">
          {{ index }}
        </span>
        {{ val }}
      </div>
    </div>
  </v-flex>
</template>
<style scoped>
.echarts {
  width: 100%;
  /* height: 250px; */
  padding-top: 5px;
  padding: 5px 10px 0 !important;
}
.legend {
  font-size: 13px;
  padding: 0 10px;
  line-height: 15px;
  margin: 5px 10px;
}
.legend-prefix {
  font-weight: bold;
  font-style: italic;
  margin-right: 5px;
  color: blue;
}
</style>

<script>
import ECharts from 'vue-echarts/components/ECharts'
import 'echarts/lib/chart/line'
import 'echarts/lib/component/tooltip'
import 'echarts/lib/component/legend'

export default {
  name: 'LokiLine',
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
      'yAxis-data': {},
    },
    loading: false,
    date: new Date().toISOString().substr(0, 10),
    menu: false,
    start: 0,
    end: 100,
    color: [
      '#c23531',
      '#2f4554',
      '#61a0a8',
      '#d48265',
      '#91c7ae',
      '#749f83',
      '#ca8622',
      '#bda29a',
      '#6e7074',
      '#546570',
      '#c4ccd3',
    ],
  }),
  computed: {
    initOptions() {
      return { width: 'auto', height: 'auto' }
    },
    options() {
      const series = []
      const legend = []
      for (const key in this.chartData['yAxis-data']) {
        series.push({
          name: key,
          type: 'line',
          data: this.chartData['yAxis-data'][key],
          areaStyle: {},
          smooth: true,
        })
        legend.push(key)
      }
      let coefficient = 1
      coefficient = parseInt(legend.length / 30 + 1)
      return {
        color: this.color,
        legend: {
          data: legend,
          top: 'bottom',
        },
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
        grid: {
          x: 50,
          y: 10,
          x2: 30,
          y2: coefficient * 30 + 10,
        },
        series: series,
      }
    },
  },
  methods: {},
}
</script>
