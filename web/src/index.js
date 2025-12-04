import { createApp } from 'vue'
import { createRouter, createWebHashHistory } from 'vue-router'
import './style.css'
import App from './App.vue'
import AdminApp from './admin/AdminApp.vue'
import Landing from './Landing.vue'
import NotFound from './NotFound.vue'
import Identities from './admin/Identities.vue'
import NewIdentity from './admin/NewIdentity.vue'
import DiscoveryServices from './admin/DiscoveryServices.vue'
import IssuedCredentials from './admin/credentials/IssuedCredentials.vue';
import RequestCredential from "./admin/credentials/RequestCredential.vue";
import WalletCredentialDetails from "./admin/credentials/WalletCredentialDetails.vue";
import Api from './plugins/api'
import IdentityDetails from "./admin/IdentityDetails.vue";
import IssueCredential from "./admin/credentials/IssueCredential.vue";
import ActivateDiscoveryService from "./admin/ActivateDiscoveryService.vue";
import UploadCredential from "./admin/credentials/UploadCredential.vue";

const routes = [
  {
    path: '/',
    components: {
      default: AdminApp
    },
    children: [
      {
        path: '',
        name: 'home',
        component: Landing,
      },
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
        ],
      },
      {
        path: 'id/:subjectID/credential/:credentialID',
        name: 'admin.credentialDetails',
        component: WalletCredentialDetails,
      },
      {
        path: 'id/:subjectID',
        name: 'admin.identityDetails',
        component: IdentityDetails,
      },
      {
        path: 'id/:subjectID/upload',
        name: 'admin.uploadCredential',
        component: UploadCredential,
      },
      {
        path: 'id/:subjectID/request',
        name: 'admin.requestCredential',
        component: RequestCredential,
      },
      {
        path: 'id/:subjectID/discovery/:discoveryServiceID/activate',
        name: 'admin.activateDiscoveryService',
        component: ActivateDiscoveryService
      },
      {
        path: 'vc/issue/:credentialType?/:subjectDID?',
        name: 'admin.issueCredential',
        component: IssueCredential
      },
      {
        path: 'vc/issuer',
        name: 'admin.issuedCredentials',
        component: IssuedCredentials
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
