// TODO: Pagination, Filters
import { errorForStatus } from './errors'

export const actions = {
  getAll(context: any) {
    return new Promise((resolve, reject) => {
      fetch(
        process.env.baseUrl +
          '/' +
          context.rootState.tenants.tenant.id +
          '/applications',
        {
          headers: new Headers({
            Authorization:
              'Bearer ' + context.rootState.authentication.authToken,
          }),
        }
      )
        .then(async (res) => {
          if (res.status !== 200) {
            reject(
              errorForStatus(
                context,
                res,
                'Error fetching applications from server'
              )
            )
            return
          }

          resolve(await res.json())
        })
        .catch((err) => {
          console.error(err)
          reject(new Error('An error occurred communicating with the server'))
        })
    })
  },
  getByID(context: any, applicationID: string) {
    return new Promise((resolve, reject) => {
      fetch(
        process.env.baseUrl +
          '/' +
          context.rootState.tenants.tenant.id +
          '/application/' +
          encodeURI(applicationID),
        {
          headers: new Headers({
            Authorization:
              'Bearer ' + context.rootState.authentication.authToken,
          }),
        }
      )
        .then(async (res) => {
          if (res.status !== 200) {
            reject(
              errorForStatus(
                context,
                res,
                'Error fetching application from server'
              )
            )
            return
          }

          resolve(await res.json())
        })
        .catch((err) => {
          console.error(err)
          reject(new Error('An error occurred communicating with the server'))
        })
    })
  },
  create(context: any, createApplicationRequest: any) {
    return new Promise((resolve, reject) => {
      fetch(
        process.env.baseUrl +
          '/' +
          context.rootState.tenants.tenant.id +
          '/applications',
        {
          method: 'POST',
          headers: new Headers({
            'Content-Type': 'application/json',
            Authorization:
              'Bearer ' + context.rootState.authentication.authToken,
          }),
          body: JSON.stringify(createApplicationRequest),
        }
      )
        .then(async (res) => {
          if (res.status !== 200 && res.status !== 204) {
            reject(
              errorForStatus(
                context,
                res,
                'Error creating new application on server'
              )
            )
            return
          }

          const body = await res.json()
          resolve(body.application_id)
        })
        .catch((err) => {
          console.error(err)
          reject(new Error('An error occurred communicating with the server'))
        })
    })
  },
  patch(context: any, params: any) {
    return new Promise((resolve, reject) => {
      fetch(
        process.env.baseUrl +
          '/' +
          context.rootState.tenants.tenant.id +
          '/application/' +
          encodeURI(params.id),
        {
          method: 'PATCH',
          headers: new Headers({
            'Content-Type': 'application/json',
            Authorization:
              'Bearer ' + context.rootState.authentication.authToken,
          }),
          body: JSON.stringify(params.patch),
        }
      )
        .then((res) => {
          if (res.status !== 200 && res.status !== 204) {
            reject(
              errorForStatus(
                context,
                res,
                'Error patching application on server'
              )
            )
            return
          }

          resolve()
        })
        .catch((err) => {
          console.error(err)
          reject(new Error('An error occurred communicating with the server'))
        })
    })
  },
  delete(context: any, id: string) {
    return new Promise((resolve, reject) => {
      fetch(
        process.env.baseUrl +
          '/' +
          context.rootState.tenants.tenant.id +
          '/application/' +
          encodeURI(id),
        {
          method: 'DELETE',
          headers: new Headers({
            'Content-Type': 'application/json',
            Authorization:
              'Bearer ' + context.rootState.authentication.authToken,
          }),
        }
      )
        .then((res) => {
          if (res.status !== 200 && res.status !== 204) {
            reject(
              errorForStatus(
                context,
                res,
                'Error deleting application on server'
              )
            )
            return
          }

          resolve()
        })
        .catch((err) => {
          console.error(err)
          reject(new Error('An error occurred communicating with the server'))
        })
    })
  },
}
