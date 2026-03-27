/**
 * Note: sub-menu only appear when route children.length >= 1
 * Detail see: https://panjiachen.github.io/vue-element-admin-site/guide/essentials/router-and-nav.html
 *
 * hidden: true                   if set true, item will not show in the sidebar(default is false)
 * alwaysShow: true               if set true, will always show the root menu
 *                                if not set alwaysShow, when item has more than one children route,
 *                                it will becomes nested mode, otherwise not show the root menu
 * redirect: noRedirect           if set noRedirect will no redirect in the breadcrumb
 * name:'router-name'             the name is used by <keep-alive> (must set!!!)
 * meta : {
    roles: ['admin','editor']    control the page roles (you can set multiple roles)
    title: 'title'               the name show in sidebar and breadcrumb (recommend set)
    icon: 'svg-name'/'el-icon-x' the icon show in the sidebar
    breadcrumb: false            if set false, the item will hidden in breadcrumb(default is true)
    activeMenu: '/example/list'  if set path, the sidebar will highlight the path you set
  }
 */

/**
 * constantRoutes
 * a base page that does not have permission requirements
 * all roles can be accessed
 */
import LoadingView from '@/pages/LoadingView.vue'
import UpgradeView from '@/pages/UpgradeView.vue';
import WelcomePage from '@/pages/WelcomePage.vue';

import OldPageLayout from '@/layouts/App1/Layout.vue'
// // Key:name
// export const redirectRoute = {
//   "das_topo":"DAS Topo"
// }

export default [
  {
    path: '/loading',
    name: 'loading',
    component: LoadingView,
    meta: { title: 'Loading' }
  },
  {
    path: '/:sub/upgrading',
    name: 'upgrading',
    component: UpgradeView,
    meta: { title: 'Upgrading', }
  },
  {
    path: '/welcome',
    name: 'welcome',
    component: WelcomePage,
    meta: { title: 'Welcome' }
  },
  {
    path: '/',
    redirect: '/loading',
  },
  {
    path: '/elem',
    component: OldPageLayout,
    meta: { title: 'Home' },
    children: [
      {
        path: ":sub/:module/:page",
        name: "app1Page",
        component: () => import('@/pages/LayoutMain/index.vue'),
        meta: { title: "Page" }
      }
    ]
  },
  {
    path: '/404',
    component: () => import('@/pages/404Page.vue'),
    hidden: true
  },
  {
    path: '/login',
    name: 'login',
    component: () => import('@/pages/Login/index.vue'),
    meta: { title: 'login' },
    hidden: false
  },
  { path: '/:pathMatch(.*)*', hidden: true, redirect: '/404' },
]

