// TODO: Pagination, Filters
import { errorForStatus } from './errors'

export const actions = {
  getAll(context: any) {
    return new Promise((resolve, reject) => {
      fetch(
        process.env.baseUrl +
          '/' +
          context.rootState.tenants.tenant.id +
          '/groups',
        {
          headers: new Headers({
            Authorization:
              'Bearer ' + context.rootState.authentication.authToken,
          }),
        }
      )
        .then(async (res) => {
          if (res.status !== 200) {
            reject(errorForStatus(res, 'Error fetching groups from server'))
            return
          }

          const groups = await res.json()
          resolve(groups)
        })
        .catch((err) => {
          console.error(err)
          reject(new Error('An error occurred communicating with the server'))
        })
    })
  },
  getByID(context: any, groupID: string) {
    return new Promise((resolve, reject) => {
      fetch(
        process.env.baseUrl +
          '/' +
          context.rootState.tenants.tenant.id +
          '/group/' +
          encodeURI(groupID),
        {
          headers: new Headers({
            Authorization:
              'Bearer ' + context.rootState.authentication.authToken,
          }),
        }
      )
        .then(async (res) => {
          if (res.status !== 200) {
            reject(errorForStatus(res, 'Error fetching group from server'))
            return
          }

          const group = await res.json()
          resolve(group)
        })
        .catch((err) => {
          console.error(err)
          reject(new Error('An error occurred communicating with the server'))
        })
    })
  },
  createGroup(context: any, createGroupRequest: any) {
    return new Promise((resolve, reject) => {
      fetch(
        process.env.baseUrl +
          '/' +
          context.rootState.tenants.tenant.id +
          '/groups',
        {
          method: 'POST',
          headers: new Headers({
            'Content-Type': 'application/json',
            Authorization:
              'Bearer ' + context.rootState.authentication.authToken,
          }),
          body: JSON.stringify(createGroupRequest),
        }
      )
        .then(async (res) => {
          if (res.status !== 200 && res.status !== 204) {
            reject(errorForStatus(res, 'Error creating new group on server'))
            return
          }

          const body = await res.json()
          resolve(body.group_id)
        })
        .catch((err) => {
          console.error(err)
          reject(new Error('An error occurred communicating with the server'))
        })
    })
  },
}
