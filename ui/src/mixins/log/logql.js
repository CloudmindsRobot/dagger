import { transformLogQL } from '@/api'

const logql = {
  methods: {
    async handlerLogQL(model, pod) {
      const filterData = {}
      const filters = []
      model.forEach((item) => {
        if (item.value === 'filter') {
          filters.push(item.text.substr(item.text.indexOf(':') + 1))
        } else {
          filterData[item.value] = item.text.substr(item.text.indexOf(':') + 1)
        }
      })
      try {
        const res = await transformLogQL(
          Object.assign(
            Object.assign({}, { filters: filters, pod: pod }),
            filterData,
          ),
        )
        if (res.status === 200) {
          return { logql: res.data.data, filters: filters, labels: filterData }
        }
      } catch (err) {
        return {}
      }
    },
  },
}

export default logql
