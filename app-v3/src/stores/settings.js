import defaultSettings from '@/settings'
import { defineStore } from 'pinia'

const { showSettings, fixedHeader, sidebarLogo } = defaultSettings

export const useSettingsStore = defineStore('settings', {
  state: ()=>({
    showSettings: showSettings,
    fixedHeader: fixedHeader,
    sidebarLogo: sidebarLogo
  }),
  getters: {},
  actions: {
    changeSetting({ key, value }) {
      if ([key] in this) {
        this[key] = value
      }
    }
  }
})
