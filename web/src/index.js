import { createApp } from 'vue'
import { createRouter, createWebHashHistory } from 'vue-router'
import './style.css'
import App from './App.vue'
import AdminApp from './admin/AdminApp.vue'
import NotFound from './NotFound.vue'
import Identities from './admin/Identities.vue'
import NewIdentity from './admin/NewIdentity.vue'
import DiscoveryServices from './admin/DiscoveryServices.vue'
import Api from './plugins/api'
import IdentityDetails from "./admin/IdentityDetails.vue";

const routes = [
  {
    path: '/',
    components: {
      default: AdminApp
    },
    children: [
      {
        path: 'identities',
        name: 'admin.identities',
        component: Identities,
        children: [
          {
            path: 'new',
            name: 'admin.newIdentity',
            components: {
              modal: NewIdentity
            }
          },
        ]
      },
      {
        path: 'id/:id',
        name: 'admin.identityDetails',
        component: IdentityDetails
      },
      {
        path: 'discovery',
        name: 'admin.discovery',
        component: DiscoveryServices
      }
    ],
  },
  { path: '/:pathMatch*', name: 'NotFound', component: NotFound }
]

const router = createRouter({
  // We are using the hash history for simplicity here.
  history: createWebHashHistory(),
  routes // short for `routes: routes`
})

const app = createApp(App)

app.use(router)
app.use(Api)
app.mount('#app')
