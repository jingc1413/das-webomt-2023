import { createI18n } from 'vue-i18n'

import en from "./locale/en.json";

const i18n = createI18n({
  legacy: false,
  locale: localStorage.getItem('lang') || 'en',
  fallbackLocale: 'en',
  messages: {
    en
  }
})

export default i18n

export function translator(title, options=null) {
  const {t, te} = i18n.global;

  if (te(`${title}`)) return t(`${title}`, options)
  return title
}