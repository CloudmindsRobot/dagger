module.exports = {
  configureWebpack: {
    devServer: {
      proxy: {
        '/api/v1/loki': {
          target: 'http://127.0.0.1:8088',
          changeOrigin: true,
          pathRewrite: {
            '^/api/v1/loki': '/api/v1/loki',
          },
        },
        '/ws/loki': {
          target: 'http://127.0.0.1:8088',
          changeOrigin: true,
          pathRewrite: {
            '^/ws/loki': '/ws',
          },
        },
      },
    },
  },
  // Vue-ECharts 默认在 webpack 环境下会引入未编译的源码版本，如果你正在使用官方的 Vue CLI 来创建项目，
  // 可能会遇到默认配置把 node_modules 中的文件排除在 Babel 转译范围以外的问题。请按如下方法修改配置：
  transpileDependencies: ['vue-echarts', 'resize-detector'],
}
