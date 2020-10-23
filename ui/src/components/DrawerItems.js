export default {
  text: 'logs',
  icon: 'notes',
  children: [
    {
      heading: 'LOGS',
    },
    {
      icon: 'keyboard_arrow_up',
      'icon-alt': 'notes',
      text: '日志查询',
      model: false,
      index: 'loki-view',
      children: [
        {
          text: '日志查看器',
          title: 'loki-viewer',
          icon: 'format_list_numbered',
          href: { name: 'loki-viewer' },
          permission: '',
        },
        {
          text: '查询历史',
          title: 'loki-history',
          icon: 'history',
          href: { name: 'loki-history' },
          permission: '',
        },
        {
          text: '日志快照',
          title: 'loki-snapshot',
          icon: 'camera_alt',
          href: { name: 'loki-snapshot' },
          permission: '',
        },
      ],
    },
  ],
}
