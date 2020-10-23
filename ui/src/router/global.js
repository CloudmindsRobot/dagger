import views from '@/views'

// 全局部分
export const global = [
  // 首页
  {
    path: '/',
    name: 'index',
    redirect: { name: 'loki-viewer' },
  },
  // 登录
  {
    path: '/login',
    name: 'login',
    component: views.Login,
    meta: { requireAuth: false },
  },
  // 注册
  {
    path: '/register',
    name: 'register',
    component: views.Register,
    meta: { requireAuth: false },
  },
]
