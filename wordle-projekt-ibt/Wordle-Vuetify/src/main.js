import './assets/main.css'
import 'vuetify/styles'


import { createApp } from 'vue'
import { createVuetify } from 'vuetify'
import App from './App.vue'
import router from './router'

import * as components from 'vuetify/components'
import * as directives from 'vuetify/directives'
import '@mdi/font/css/materialdesignicons.css'

const vuetify = createVuetify({
  components,
  directives,
  theme: {
    defaultTheme: 'dark',
    themes: {
      light: {
        colors: {
          theme_color_1: '#f5f5f5', // hellgrau/weiß
        },
      },
      dark: {
        colors: {
          theme_color_1: '#2f2f2f', // dunkelgrau
        },
      },
    },
  },
})

createApp(App)
  .use(vuetify)
  .use(router)
  .mount('#app')


