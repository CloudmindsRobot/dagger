function DrawerItems(settings) {
  let alertBar = []
  if (settings.allowSignUp) {
    alertBar = [
      {
        text: '告警事件',
        title: 'loki-event',
        icon: 'ring_volume',
        href: { name: 'loki-event' },
        permission: '',
      },
      {
        text: '组',
        title: 'loki-group',
        icon: 'groups',
        href: { name: 'loki-group' },
        permission: '',
      },
      {
        text: '告警规则',
        title: 'loki-rule',
        icon: 'rule',
        href: { name: 'loki-rule' },
        permission: '',
      },
    ]
  } else {
    alertBar = [
      {
        text: '告警事件',
        title: 'loki-event',
        icon: 'ring_volume',
        href: { name: 'loki-event' },
        permission: '',
      },
      {
        text: '告警规则',
        title: 'loki-rule',
        icon: 'rule',
        href: { name: 'loki-rule' },
        permission: '',
      },
    ]
  }
  const items = {
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

  if (settings.alertEnabled) {
    items.children.push({
      heading: 'ALERTS',
    })
    items.children.push({
      icon: 'keyboard_arrow_up',
      'icon-alt': 'warning',
      text: '告警',
      model: false,
      index: 'loki-alert',
      children: alertBar,
    })
  }

  return items
}

export default DrawerItems
