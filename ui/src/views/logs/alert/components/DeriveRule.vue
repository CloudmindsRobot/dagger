<template>
  <span>
    <v-dialog v-model="dialog" max-width="500" persistent scrollable>
      <v-card class="px-1">
        <v-card-title>
          <span class="headline">派生告警规则</span>
        </v-card-title>
        <v-divider></v-divider>
        <v-card-text>
          <v-container>
            <v-row>
              <v-col>
                <v-form v-model="valid" ref="deriveRuleForm">
                  <v-text-field
                    v-model="obj.name"
                    :rules="objRules.nameRules"
                    required
                    dense
                    label="名称"
                  ></v-text-field>
                  <v-textarea
                    v-model="obj.log_ql"
                    :rules="objRules.exprRules"
                    auto-grow
                    required
                    readonly
                    dense
                    label="LogQL"
                  >
                  </v-textarea>
                  <v-textarea
                    v-model="obj.description"
                    :rules="objRules.descriptionRules"
                    required
                    auto-grow
                    dense
                    label="描述description（参考loki rule官方定义）"
                  ></v-textarea>
                  <v-combobox
                    v-if="settings && settings.allowSignUp"
                    v-model="obj.groups"
                    :items="groupItems"
                    color="primary"
                    item-text="group_name"
                    item-value="id"
                    label="分发组"
                    full-width
                    clearable
                    dense
                    hide-selected
                    multiple
                    @change="handlerShowGroupUsers"
                  >
                    <template v-slot:selection="{ item }">
                      <v-chip color="primary" style="margin: 1px 3px;">
                        {{ item['group_name'] }}
                      </v-chip>
                    </template>
                  </v-combobox>
                  <v-chip
                    v-for="user in users"
                    :key="user"
                    color="primary"
                    style="margin: 1px 3px;"
                  >
                    {{ user }}
                  </v-chip>
                  <v-switch
                    v-model="labeled"
                    label="自定义Label"
                    @change="handlerChangeLabelSwitch"
                  ></v-switch>
                  <v-row v-if="labeled">
                    <v-text-field
                      class="mx-2"
                      v-model="key"
                      required
                      dense
                      label="Key"
                    ></v-text-field
                    ><v-text-field
                      class="mx-2"
                      v-model="value"
                      required
                      dense
                      label="Value"
                    ></v-text-field>
                    <v-btn
                      class="mx-2"
                      fab
                      dark
                      small
                      color="primary"
                      @click="handlerAddLabel"
                    >
                      <v-icon dark>
                        add
                      </v-icon>
                    </v-btn>
                  </v-row>
                  <template v-if="labeled">
                    <v-chip
                      v-for="(label, index) in obj.labels"
                      :key="index"
                      color="primary"
                      close
                      style="margin: 1px 3px;"
                      @click:close="handleRemoveLabel(label)"
                    >
                      {{ label.key }}:{{ label.value }}
                    </v-chip>
                  </template>
                </v-form>
              </v-col>
            </v-row>
          </v-container>
        </v-card-text>
        <v-divider></v-divider>
        <v-card-actions>
          <v-layout row justify-space-around>
            <v-flex xs5>
              <v-btn color="primary" block @click="handlerCreateLogRule">
                创建
              </v-btn>
            </v-flex>
            <v-flex xs5>
              <v-btn @click="handlerCloseDialog" block>
                关闭
              </v-btn>
            </v-flex>
          </v-layout>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </span>
</template>

<script>
import { createLogRule, listLogGroup } from '@/api'
import { mapState } from 'vuex'

export default {
  name: 'DeriveRule',
  props: {
    logQL: {
      type: String,
      default: () => null,
    },
  },
  data: () => ({
    show: false,
    dialog: false,
    valid: false,
    loading: false,
    labeled: false,
    obj: {
      name: '',
      expr: '',
      for: 0,
      labels: [],
      groups: [],
      description: '',
      log_ql: '',
    },
    objRules: {
      nameRules: [(v) => !!v || '名称必填'],
      exprRules: [(v) => !!v || 'LogQL必填'],
      descriptionRules: [(v) => !!v || '描述必填'],
    },
    users: [],
    groupItems: [],
    key: '',
    value: '',
  }),
  computed: {
    ...mapState(['settings']),
  },
  mounted() {
    this.$nextTick(() => {
      this.listLogGroup()
    })
  },
  methods: {
    async listLogGroup() {
      try {
        const res = await listLogGroup({ size: 1000 })
        if (res.status === 200 && res.data.success) {
          this.groupItems = res.data.data
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
    },
    async handlerCreateLogRule() {
      try {
        if (this.$refs.deriveRuleForm.validate()) {
          const res = await createLogRule(
            Object.assign(this.obj, {
              log_ql: this.logQL,
            }),
          )
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
    handlerShowGroupUsers() {
      this.users = []
      this.obj.groups.forEach((group) => {
        group.users.forEach((user) => {
          if (this.users.indexOf(user.user.username) === -1) {
            this.users.push(user.user.username)
          }
        })
      })
    },
    handlerAddLabel() {
      if (this.key.trim() === '' || this.value.trim() === '') {
        this.$store.commit('showSnackBar', {
          text: 'Warning: 请填写label的key和value',
          color: 'warning',
        })
        return
      }
      this.obj.labels.push({
        key: this.key,
        value: this.value,
      })
      this.key = ''
      this.value = ''
    },
    handleRemoveLabel(label) {
      const index = this.obj.labels.findIndex((l) => {
        return l.key === label.key && l.value === label.value
      })
      this.obj.labels = this.obj.labels
        .slice(0, index)
        .concat(this.obj.labels.slice(index + 1))
    },
    handlerChangeLabelSwitch() {
      if (!this.labeled) {
        this.obj.labels = []
      }
    },
    handlerCloseDialog() {
      this.dialog = false
      this.users = []
      this.$refs.deriveRuleForm.reset()
    },
  },
}
</script>
