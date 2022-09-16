export default async function (context: any) {
  if (context.store.state.authentication.user === null) {
    await context.store.dispatch('authentication/populateUserInfomation')
  }

  const authenticated: boolean = await context.store.dispatch(
    'authentication/isAuthenticated'
  )
  if (!authenticated) {
    context.app.router.push({
      path: '/login',
      query:
        context.route.fullPath !== '/'
          ? { redirect_to: context.route.fullPath }
          : {},
    })
  } else if (Date.now() >= context.store.state.authentication.user.exp * 1000) {
    context.app.router.push({
      path: '/login',
      query:
        context.route.fullPath !== '/'
          ? { redirect_to: context.route.fullPath }
          : {},
    })
  }
}
