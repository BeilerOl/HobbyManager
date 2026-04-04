import { createRouter, createWebHistory } from 'vue-router'
import WorkList from '../views/WorkList.vue'
import WorkDetail from '../views/WorkDetail.vue'
import WorkForm from '../views/WorkForm.vue'
import WorkImport from '../views/WorkImport.vue'

const routes = [
  { path: '/', name: 'list', component: WorkList },
  { path: '/works/import', name: 'import', component: WorkImport },
  { path: '/works/new', name: 'new', component: WorkForm, props: { isEdit: false } },
  { path: '/works/:id', name: 'detail', component: WorkDetail },
  { path: '/works/:id/edit', name: 'edit', component: WorkForm, props: { isEdit: true } },
]

export default createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes,
})
