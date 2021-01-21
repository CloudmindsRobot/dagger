<template>
  <div>
    <v-breadcrumbs :items="breadcrumbs" divider="/"></v-breadcrumbs>
    <v-card>
      <v-card-title>
        <span class="pr-4">
          <CreateGroup @refresh="refresh" />
        </span>
        <v-btn icon text color="green" @click="refresh" :loading="loading">
          <v-icon>refresh</v-icon>
        </v-btn>
      </v-card-title>
      <v-data-table
        :headers="headers"
        :items="items"
        :page.sync="params.page"
        :items-per-page="params.page_size"
        hide-default-footer
      >
        <template v-slot:item.create_at="{ item }">
          <span>{{ new Date(item.create_at).toLocaleString() }}</span>
        </template>
        <template v-slot:item.users="{ item }">
          <v-chip
            v-for="user in item.users"
            :key="user.user.id"
            color="primary"
            style="margin: 1px 3px;"
          >
            {{ user.user.username }}
          </v-chip>
        </template>
        <template v-slot:item.action="{ item }">
          <span class="pr-4" v-if="item.user.username === username">
            <UpdateGroup
              ref="updateGroup"
              :item.sync="item"
              @refresh="refresh"
            />
          </span>
          <span class="pr-4">
            <v-btn color="primary" small @click="handleJoin(item)">
              加入组
            </v-btn>
          </span>
          <span class="pr-4">
            <v-btn color="error" small @click="handleLeave(item)">
              移出组
            </v-btn>
          </span>
          <span class="pr-4" v-if="item.user.username === username">
            <v-btn color="error" small @click="handleDelete(item)">删除</v-btn>
          </span>
        </template>
      </v-data-table>
      <div class="text-xs-center pa-2">
        <v-pagination
          v-model="params.page"
          :length="pageCount"
          :total-visible="10"
          circle
          @input="listLogGroup"
        ></v-pagination>
      </div>
    </v-card>

    <DeleteGroup ref="deleteGroup" @refresh="refresh" />
  </div>
</template>

<script>
import { listLogGroup, joinLogGroup, leaveLogGroup } from '@/api'
import CreateGroup from './components/CreateGroup'
import DeleteGroup from './components/DeleteGroup'
import UpdateGroup from './components/UpdateGroup'
import { mapState } from 'vuex'

export default {
  name: 'LogGroup',
  components: {
    CreateGroup,
    DeleteGroup,
    UpdateGroup,
  },
  data: () => ({
    breadcrumbs: [
      { text: 'ALERTS', disabled: true, href: '' },
      { text: '日志告警', disabled: true, href: '' },
      { text: '组', disabled: true },
    ],
    items: [],
    loading: false,
    headers: [
      { text: '组', value: 'group_name', align: 'start' },
      { text: '成员', value: 'users', align: 'start', width: 1000 },
      { text: '创建时间', value: 'create_at', align: 'start' },
      { text: '操作', value: 'action', align: 'start' },
    ],
    pageCount: 0,
    params: {
      page: 1,
      page_size: 10,
    },
  }),
  computed: {
    ...mapState(['username']),
  },
  methods: {
    async listLogGroup() {
      this.loading = true
      try {
        const res = await listLogGroup(this.params)
        if (res.status === 200 && res.data.success) {
          this.items = res.data.data
          this.pageCount = Math.ceil(res.data.total / res.data.page_size)
        } else {
          this.$store.commit('showSnackBar', {
            text: `Warn: ${res.data.message}`,
            color: 'warning',
          })
        }
      } catch (err) {
        if (
          err.response &&
          [400, 401, 403, 504].indexOf(err.response.status) === -1
        ) {
          this.$store.commit('showSnackBar', {
            text: 'Error: 获取分发组失败',
            color: 'error',
          })
        }
      }
      this.loading = false
    },
    async handleJoin(item) {
      const res = await this.$confirm('加入该组？', {
        buttonTrueText: '确定',
        buttonFalseText: '取消',
        persistent: true,
      })
      if (res) {
        try {
          const res = await joinLogGroup({ group_id: item.id })
          if (res.status === 201) {
            this.listLogGroup()
          } else {
            this.$store.commit('showSnackBar', {
              text: `Warn: ${res.data.message}`,
              color: 'warning',
            })
          }
        } catch (err) {
          if (
            err.response &&
            [400, 401, 403, 504].indexOf(err.response.status) === -1
          ) {
            this.$store.commit('showSnackBar', {
              text: 'Error: 加入组失败',
              color: 'error',
            })
          }
        }
      }
    },
    async handleLeave(item) {
      const res = await this.$confirm('移出该组？', {
        buttonTrueText: '确定',
        buttonFalseText: '取消',
        persistent: true,
      })
      if (res) {
        try {
          const res = await leaveLogGroup({ group_id: item.id })
          if (res.status === 204) {
            this.listLogGroup()
          } else {
            this.$store.commit('showSnackBar', {
              text: `Warn: ${res.data.message}`,
              color: 'warning',
            })
          }
        } catch (err) {
          if (
            err.response &&
            [400, 401, 403, 504].indexOf(err.response.status) === -1
          ) {
            this.$store.commit('showSnackBar', {
              text: 'Error: 移出组失败',
              color: 'error',
            })
          }
        }
      }
    },
    refresh() {
      this.params.page = 1
      this.listLogGroup()
    },
    handleDelete(item) {
      this.$refs.deleteGroup.item = item
      this.$refs.deleteGroup.dialog = true
    },
  },
  mounted() {
    if (this.$store.state.jwt) {
      this.listLogGroup()
    }
  },
}
</script>
