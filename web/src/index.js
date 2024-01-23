import { createApp } from 'vue'
import { createRouter, createWebHashHistory } from 'vue-router'
import './style.css'
import App from './App.vue'
import AdminApp from './admin/AdminApp.vue'
import Login from './Login.vue'
import Logout from './Logout.vue'
import NotFound from './NotFound.vue'
import Identities from './admin/Identities.vue'
import NewIdentity from './admin/NewIdentity.vue'
import Api from './plugins/api'

const routes = [
  { path: '/', component: Login },
  {
    name: 'login',
    path: '/login',
    component: Login
  },
  {
    name: 'logout',
    path: '/logout',
    component: Logout
  },
  {
    path: '/admin',
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
          }
        ]
      }
    ],
    meta: { requiresAuth: true }
  },
  { path: '/:pathMatch*', name: 'NotFound', component: NotFound }
]

const router = createRouter({
  // We are using the hash history for simplicity here.
  history: createWebHashHistory(),
  routes // short for `routes: routes`
})

router.beforeEach((to, from) => {
  if (to.meta.requiresAuth) {
    if (localStorage.getItem('session')) {
      return true
    }
    return '/login'
  }
})

const app = createApp(App)

app.use(router)
app.use(Api, { forbiddenRoute: { name: 'logout' } })
app.mount('#app')
