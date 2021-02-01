<template>
  <v-navigation-drawer
    v-model="localSidebar"
    app
    mobile-breakpoint="10000"
    class="mt-12 pt-5"
    style="height: 100%;"
  >
    <v-list dense>
      <template v-for="item in items">
        <v-row v-if="item.heading" :key="item.heading" row align-center>
          <v-col cols="6" class="pa-1 pl-5">
            <v-subheader v-if="item.heading">{{ item.heading }}</v-subheader>
          </v-col>
        </v-row>
        <v-list-group
          v-else-if="item.children"
          v-model="item.model"
          :key="item.text"
          :prepend-icon="item.model ? item.icon : item['icon-alt']"
          append-icon
        >
          <template v-slot:activator>
            <v-list-item class="pl-0">
              <v-list-item-content>
                <v-list-item-title>{{ item.text }}</v-list-item-title>
              </v-list-item-content>
            </v-list-item>
          </template>
          <template v-for="(child, i) in item.children">
            <v-list-item
              v-if="
                !child.permission ||
                  superuser ||
                  permissions.indexOf(child.permission) > -1
              "
              :key="i"
              :to="child.href"
            >
              <v-list-item-action>
                <v-icon v-if="child.icon">{{ child.icon }}</v-icon>
              </v-list-item-action>
              <v-list-item-content :title="child.title">
                <v-list-item-title>{{ child.text }}</v-list-item-title>
              </v-list-item-content>
            </v-list-item>
          </template>
        </v-list-group>
        <v-divider v-else-if="item.divider" :key="item.text"></v-divider>
        <v-list-item v-else :key="item.text">
          <v-list-item-action>
            <v-icon>{{ item.icon }}</v-icon>
          </v-list-item-action>
          <v-list-item-content>
            <v-list-item-title>{{ item.text }}</v-list-item-title>
          </v-list-item-content>
        </v-list-item>
      </template>
    </v-list>
  </v-navigation-drawer>
</template>

<script>
import { mapState } from 'vuex'
import DrawerItems from '@/components/DrawerItems'
export default {
  name: 'BaseSidebar',
  data: () => ({
    localSidebar: false,
    localApp: null,
  }),
  watch: {
    localSidebar(val) {
      if (this.sidebar !== val) this.$store.commit('setSidebar', val)
    },
    sidebar(val) {
      if (this.localSidebar !== val) this.localSidebar = val
    },
    localApp(val) {
      if (this.app !== val) this.$store.commit('setApp', val)
    },
    app(val) {
      if (this.localApp !== val) this.localApp = val
    },
  },
  computed: {
    ...mapState(['sidebar', 'app', 'settings']),
    items() {
      return DrawerItems(this.settings).children
    },
  },
  mounted() {
    this.localSidebar = this.sidebar
  },
}
</script>
