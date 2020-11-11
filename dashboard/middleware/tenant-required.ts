export default function (context: any) {
  if (context.store.state.tenants.tenant === null) {
    let params = ''
    if (context.route.fullPath !== '/') {
      params += '?redirect_to=' + encodeURIComponent(context.route.fullPath)
    }
    context.app.router.push('/login' + params)
  }
}
