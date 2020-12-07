// TODO: Pagination, Filters
import { errorForStatus } from './errors'

export const actions = {
  getAll(context: any) {
    return new Promise((resolve, reject) => {
      fetch(
        process.env.baseUrl +
          '/' +
          context.rootState.tenants.tenant.id +
          '/devices',
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
              errorForStatus(context, res, 'Error fetching devices from server')
            )
            return
          }

          const devices = await res.json()
          resolve(devices)
        })
        .catch((err) => {
          console.error(err)
          reject(new Error('An error occurred communicating with the server'))
        })
    })
  },
  getByID(context: any, deviceID: string) {
    return new Promise((resolve, reject) => {
      fetch(
        process.env.baseUrl +
          '/' +
          context.rootState.tenants.tenant.id +
          '/device/' +
          encodeURI(deviceID),
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
              errorForStatus(context, res, 'Error fetching device from server')
            )
            return
          }

          const device = await res.json()
          resolve(device)
        })
        .catch((err) => {
          console.error(err)
          reject(new Error('An error occurred communicating with the server'))
        })
    })
  },
  getInformationByID(context: any, deviceID: string) {
    return new Promise((resolve, reject) => {
      fetch(
        process.env.baseUrl +
          '/' +
          context.rootState.tenants.tenant.id +
          '/device/' +
          encodeURI(deviceID) +
          '/info',
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
                'Error fetching device information from server'
              )
            )
            return
          }

          const deviceInfo = await res.json()
          resolve(deviceInfo)
        })
        .catch((err) => {
          console.error(err)
          reject(new Error('An error occurred communicating with the server'))
        })
    })
  },
  getScopeByID(context: any, deviceID: string) {
    return new Promise((resolve, reject) => {
      fetch(
        process.env.baseUrl +
          '/' +
          context.rootState.tenants.tenant.id +
          '/device/' +
          encodeURI(deviceID) +
          '/scope',
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
                'Error fetching device scope from server'
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
  patchDevice(context: any, params: any) {
    return new Promise((resolve, reject) => {
      fetch(
        process.env.baseUrl +
          '/' +
          context.rootState.tenants.tenant.id +
          '/device/' +
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
              errorForStatus(context, res, 'Error patching device on server')
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
  addToGroup(context: any, params: any) {
    return new Promise((resolve, reject) => {
      fetch(
        process.env.baseUrl +
          '/' +
          context.rootState.tenants.tenant.id +
          '/group/' +
          encodeURI(params.gid) +
          '/devices',
        {
          method: 'POST',
          headers: new Headers({
            'Content-Type': 'application/json',
            Authorization:
              'Bearer ' + context.rootState.authentication.authToken,
          }),
          body: JSON.stringify({
            devices: [params.udid],
          }),
        }
      )
        .then((res) => {
          if (res.status !== 200 && res.status !== 204) {
            reject(
              errorForStatus(
                context,
                res,
                'Error adding device to group on server'
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
  removeFromGroup(context: any, params: any) {
    return new Promise((resolve, reject) => {
      fetch(
        process.env.baseUrl +
          '/' +
          context.rootState.tenants.tenant.id +
          '/group/' +
          encodeURI(params.gid) +
          '/devices',
        {
          method: 'DELETE',
          headers: new Headers({
            'Content-Type': 'application/json',
            Authorization:
              'Bearer ' + context.rootState.authentication.authToken,
          }),
          body: JSON.stringify({
            devices: [params.udid],
          }),
        }
      )
        .then((res) => {
          if (res.status !== 200 && res.status !== 204) {
            reject(
              errorForStatus(
                context,
                res,
                'Error removing device from group on server'
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
  deleteDevice(context: any, id: string) {
    return new Promise((resolve, reject) => {
      fetch(
        process.env.baseUrl +
          '/' +
          context.rootState.tenants.tenant.id +
          '/device/' +
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
              errorForStatus(context, res, 'Error deleting device on server')
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
