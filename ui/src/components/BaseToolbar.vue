<template>
  <v-app-bar
    :clipped-left="$vuetify.breakpoint.lgAndUp"
    color="#326de6"
    dark
    app
    style="z-index: 7"
  >
    <v-app-bar-nav-icon @click.stop="switchSidebar" />
    <span class="hidden-sm-and-down"
      ><v-img src="/logo.png" height="35" width="35" contain></v-img
    ></span>
    <v-toolbar-title>
      <span class="hidden-sm-and-down pl-2">{{ toolbarTitle }}</span>
      <span class="hidden-sm-and-down pl-2">{{
        app && app.toUpperCase()
      }}</span>
    </v-toolbar-title>
    <v-btn icon class="hidden-lg-and-up" @click.stop="switchDrawer"
      ><v-icon>more_vert</v-icon></v-btn
    >
    <v-spacer></v-spacer>

    <template v-for="item in items">
      <v-menu
        offset-y
        v-if="item.children"
        :key="item.title"
        attach
        content-class="zoom-menu"
      >
        <template v-slot:activator="{ on }">
          <v-btn v-on="on" icon :title="item.title">
            <v-icon>{{ item.icon }}</v-icon>
          </v-btn>
        </template>
        <v-list v-if="item.children && item.children.length > 0">
          <v-list-item v-for="(child, index) in item.children" :key="index">
            <v-list-item-title v-if="child.outer"
              ><a :href="child.href">{{ child.text }}</a></v-list-item-title
            >
            <v-list-item-title v-else-if="child.href">
              <router-link :to="child.href" class="black--text">{{
                child.text
              }}</router-link>
            </v-list-item-title>
            <v-list-item-title v-else>
              <span @click="logout">{{ child.text }}</span>
            </v-list-item-title>
          </v-list-item>
        </v-list>
      </v-menu>
      <v-btn icon v-else :title="item.title" :key="item.title">
        <a
          v-if="item.outer"
          target="_blank"
          :href="item.href"
          class="white--text"
        >
          <v-icon>{{ item.icon }}</v-icon>
        </a>
        <router-link v-else :to="item.href" class="white--text">
          <v-badge v-if="item.badge && badge" color="error" overlap>
            <template v-slot:badge>
              <span>{{ badge }}</span>
            </template>
            <v-icon>
              {{ item.icon }}
            </v-icon>
          </v-badge>
          <v-icon v-else>{{ item.icon }}</v-icon>
        </router-link>
      </v-btn>
    </template>
  </v-app-bar>
</template>

<script>
import { mapState } from 'vuex'
import ToolbarItems from './ToolbarItems.js'

export default {
  name: 'BaseToolBar',
  data: () => ({
    indexUrl: { name: 'index' },
    toolbarTitle: 'DAGGER',
  }),
  computed: {
    ...mapState(['app']),
    items() {
      return ToolbarItems
    },
    badge() {
      return ''
    },
  },
  methods: {
    switchSidebar() {
      this.$store.commit('switchSidebar')
    },
    switchDrawer() {
      this.$store.commit('switchDrawer')
    },
    logout() {
      this.$store.commit('logout')
      this.$router.push({ name: 'login' })
    },
  },
  mounted() {},
}
</script>

<style scoped>
a {
  text-decoration: none;
  color: inherit;
}
.zoom-menu {
  right: 10px;
  left: auto !important;
}
</style>
