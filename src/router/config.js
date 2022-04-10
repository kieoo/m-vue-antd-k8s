import TabsView from '@/layouts/tabs/TabsView'
// import BlankView from '@/layouts/BlankView'
// import PageView from '@/layouts/PageView'

// 路由配置
const options = {
  routes: [
    {
      path: '*',
      name: '404',
      component: () => import('@/pages/exception/404'),
    },
    {
      path: '/403',
      name: '403',
      component: () => import('@/pages/exception/403'),
    },
    {
      path: '/',
      name: '首页',
      component: TabsView,
      redirect: '/k8sExec',
      children: [
        {
          path: 'k8sExec',
          name: 'k8sExec',
          meta: {
            icon: 'form',
            page: {
              cacheAble: false
            }
          },
          component: () => import('@/pages/form/k8sExec')
        },
        {
          path: 'k8sHelmCheck',
          name: 'k8sYaml&HelmCheck',
          meta: {
            icon: 'form',
            page: {
              cacheAble: false
            }
          },
          component: () => import('@/pages/form/k8sYamlCheck2'),
        }
      ]
    },
  ]
}

export default options
