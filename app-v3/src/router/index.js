import { createRouter, createWebHashHistory } from "vue-router";
import defRoute from "./defRoute";
import { getPageVersion } from "@/utils/index.js";

const router = createRouter({
  history: createWebHashHistory(),
  routes: defRoute,
  scrollBehavior: () => ({ y: 0 }),
});

import NProgress from "nprogress"; // progress bar
import "nprogress/nprogress.css"; // progress bar style
import { useAuthStore } from "@/stores/auth";
import { useDasDevices } from "@/stores/das-devices";
import { useAppStore } from "@/stores/app";

NProgress.configure({ showSpinner: false }); // NProgress Configuration

router.beforeEach(async (to, from, next) => {
  // start progress bar
  NProgress.start();
  const redirect = to.query?.redirect;
  console.log("to Path", to.path, redirect);

  const authStore = useAuthStore();
  const appStore = useAppStore();

  const isLogin = await authStore.isLogin();
  if (!isLogin) {
    if (to.path !== `/login`) {
      if (to.path && to.path !== "/loading") {
        next(`/login?redirect=${to.path}`);
      } else {
        next(`/login`);
      }
      return;
    } else {
      next();
      return;
    }
  }

  if (!appStore.isAppSetupDone) {
    if (to.path !== "/welcome") {
      if (to.path && to.path !== "/loading") {
        next(`/welcome?redirect=${to.path}`);
      } else {
        next(`/welcome`);
      }
      return;
    } else {
      next();
      return;
    }
  }

  if (to.path == "/loading" || to.path == "/login" || to.path == "/welcome") {
    if (redirect) {
      next(redirect);
      return;
    } else {
      let pageVersion = getPageVersion();
      next(`/${pageVersion}/local/overview/das_topo`);
      return;
    }
  } else {
    const sub = to.params.sub;
    const moduleName = to.params.module;
    const pageName = to.params.page;
    if (sub !== undefined) {
      const dasDevices = useDasDevices();
      await dasDevices.selectDeviceSub(sub);
      const dev = dasDevices.currentDevice;

      if (moduleName && pageName) {
        let doSelect = dev.layout.selectPage(moduleName, pageName);
        if (doSelect == false) {
          next({
            name: "app1Page",
            params: {
              sub: sub,
              module: 'overview',
              page:'element_information'
            } 
          })
          return
        }
      }
      if (dev.upgrade.upgradeState !== undefined) {
        if (to.path !== `/${sub}/upgrading`) {
          if (to.path && to.path !== "/loading") {
            next(`/${sub}/upgrading?redirect=${to.path}`);
          } else {
            next(`/${sub}/upgrading`);
          }
          return;
        }
      }
    }
  }
  // if (to.path === "/loading" || to.path === "/login" || to.path === "/welcome" || to.path === `/${sub}/upgrading`) {
  //   if (to.params.redirect) {
  //     next(to.params.redirect);
  //     return;
  //   } else {
  //     next(`/${sub}/overview/das_topo`);
  //     return;
  //   }
  // }

  if (redirect) {
    next(redirect);
    return;
  }
  next();
});

router.afterEach((to, from, failure) => {
  // finish progress bar
  NProgress.done();

  // if (failure) {
  //   console.error('router error');
  //   console.error(failure);
  // }
});

export function resetRouter() {
  const newRouter = createRouter({
    history: createWebHashHistory(),
    routes: defRoute,
    scrollBehavior: () => ({ y: 0 }),
  });
  router.matcher = newRouter.matcher; // reset router
}

export default router;
