import _ from 'lodash'
import { merge } from '@/utils/helpers'

const asyncDataTable = {
  data() {
    return {
      items: [],
      selected: [],
      loading: false,
      search: '',
      pagination: {
        page: 1,
        itemsPerPage: 20,
      },
      totalItems: 0,
      params: {},
    }
  },
  watch: {
    pagination: {
      async handler({ page, itemsPerPage, sortBy, sortDesc }) {
        this.params.page_size = itemsPerPage
        this.params.page = page
        const ordering = []
        if (sortDesc.length > 0 && sortBy.length > 0) {
          sortBy.forEach((item, index) => {
            if (sortDesc[index]) {
              ordering.push(`-${item}`)
            } else {
              ordering.push(item)
            }
          })
        } else if (sortDesc.length === 0 && sortBy.length > 0) {
          sortBy.forEach((item) => ordering.push(item))
        }
        if (ordering.length > 0) this.params.ordering = ordering.join(',')
        await this.getDataFromAPI()
      },
      deep: true,
    },
    search(val) {
      if (val) {
        this.params.search = val
      } else {
        delete this.params.search
      }
      this.debounceSearch()
    },
  },
  methods: {
    async getDataFromAPI() {
      this.loading = true
      try {
        // CAUTION: mapping this function to the real API function
        // e.g., listAPI: listService
        const resp = await this.listAPI(merge(this.params, this.$route.query))
        const items = []
        resp.data.data.forEach(function(item) {
          items.push(item)
        })
        this.items = items
        this.totalItems = resp.data.total
      } catch (e) {
        this.$store.commit('showSnackBar', {
          text: '获取数据失败',
          color: 'error',
        })
      } finally {
        this.loading = false
      }
    },
    refresh() {
      this.getDataFromAPI()
      this.selected = []
    },
    // eslint-disable-next-line
    debounceSearch: _.debounce(function() {
      this.getDataFromAPI()
    }, 500),
    updatePagination(val) {
      this.pagination = val
    },
    updateSearch(val) {
      this.search = val
    },
  },
}

export default asyncDataTable
