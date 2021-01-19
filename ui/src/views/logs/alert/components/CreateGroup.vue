<template>
  <span>
    <v-btn color="primary" @click.stop="handlerOpenDialog">创建组</v-btn>
    <v-dialog v-model="dialog" max-width="500" persistent scrollable>
      <v-card class="px-1">
        <v-card-title>
          <span class="headline">创建组</span>
        </v-card-title>
        <v-divider></v-divider>
        <v-card-text>
          <v-container>
            <v-row>
              <v-col>
                <v-form v-model="valid" ref="createGroup">
                  <v-text-field
                    v-model="obj.group_name"
                    :rules="objRules.groupNameRules"
                    required
                    dense
                    label="名称"
                  ></v-text-field>
                  <v-combobox
                    v-model="obj.users"
                    :items="users"
                    color="primary"
                    item-text="username"
                    item-value="id"
                    label="用户"
                    full-width
                    clearable
                    dense
                    hide-selected
                    multiple
                  >
                    <template v-slot:selection="{ item }">
                      <v-chip color="primary" style="margin: 1px 3px;">
                        {{ item['username'] }}
                      </v-chip>
                    </template>
                  </v-combobox>
                </v-form>
              </v-col>
            </v-row>
          </v-container>
        </v-card-text>
        <v-divider></v-divider>
        <v-card-actions>
          <v-layout row justify-space-around>
            <v-flex xs5>
              <v-btn color="primary" block @click="handlerCreateUserGroup">
                创建
              </v-btn>
            </v-flex>
            <v-flex xs5>
              <v-btn @click="handlerCloseDialog" block>关闭</v-btn>
            </v-flex>
          </v-layout>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </span>
</template>

<script>
import { createLogGroup, listLogUser } from '@/api'

export default {
  name: 'CreateGroup',
  data: () => ({
    show: false,
    dialog: false,
    valid: false,
    loading: false,
    obj: {
      group_name: '',
      users: [],
    },
    objRules: {
      groupNameRules: [(v) => !!v || '名称必填'],
      usersRules: [(v) => !!v || '用户必填'],
    },
    users: [],
  }),
  methods: {
    async handlerCreateUserGroup() {
      try {
        if (this.$refs.createGroup.validate()) {
          const res = await createLogGroup(this.obj)
          if (res.status === 201) {
            this.$emit('refresh')
            this.$store.commit('showSnackBar', {
              text: 'Success: 创建成功',
              color: 'success',
            })
          }
          this.handlerCloseDialog()
        }
      } catch (err) {
        if (
          err.response &&
          [400, 401, 403, 504].indexOf(err.response.status) === -1
        ) {
          this.$store.commit('showSnackBar', {
            text: 'Error: 创建失败',
            color: 'error',
          })
        }
      }
    },
    async habdlerListUser() {
      try {
        const res = await listLogUser({ page_size: 500 })
        if (res.status === 200) {
          this.users = res.data.data
        }
      } catch (err) {
        if (
          err.response &&
          [400, 401, 403, 504].indexOf(err.response.status) === -1
        ) {
          this.$store.commit('showSnackBar', {
            text: 'Error: 获取用户失败',
            color: 'error',
          })
        }
      }
    },
    handlerCloseDialog() {
      this.dialog = false
      this.$refs.createGroup.reset()
    },
    handlerOpenDialog() {
      this.habdlerListUser()
      this.dialog = true
    },
  },
  mounted() {},
}
</script>
