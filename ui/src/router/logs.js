import views from '@/views'

// path和name字段需要分别添加平台标识`/${app}`和`${app}-${routeName}`
// 日志平台部分
export const logs = [
  {
    path: '/logs/viewer',
    name: 'loki-view',
    component: views.Root,
    children: [
      {
        path: 'query',
        name: 'loki-viewer',
        component: views.logs.viewer.LokiViewer,
        meta: {
          requireAuth: true,
        },
      },
      {
        path: 'history',
        name: 'loki-history',
        component: views.logs.history.LokiHistory,
        meta: {
          requireAuth: true,
        },
      },
      {
        path: 'snapshot',
        name: 'loki-snapshot',
        component: views.logs.snapshot.LokiSnapshot,
        meta: {
          requireAuth: true,
        },
      },
    ],
  },
]
