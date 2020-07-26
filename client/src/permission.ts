import router from './router'
import NProgress from 'nprogress'
import 'nprogress/nprogress.css'
import { Message } from 'element-ui'
import { Route } from 'vue-router'
import { UserModule } from '@/store/modules/user'

NProgress.configure({ showSpinner: false })

router.beforeEach(async(to: Route, _: Route, next: any) => {
  // Start progress bar
  NProgress.start()

  // Determine whether the user has logged in
  if (UserModule.name) {
    next()
  } else {
    try {
      // Get user info, including roles
      await UserModule.GetUserInfo()
      // Set the replace: true, so the navigation will not leave a history record
      next({ ...to, replace: true })
    } catch (err) {
      next()
      NProgress.done()
    }
  }
})

router.afterEach((to: Route) => {
  // Finish progress bar
  NProgress.done()

  // set page title
  document.title = `${to.meta.title} - oiscon`
})
