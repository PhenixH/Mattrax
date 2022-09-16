// TODO: Pagination, Filters
import { errorForStatus } from './errors'

export const actions = {
  getAll(context: any) {
    return new Promise((resolve, reject) => {
      fetch(
        process.env.baseUrl +
          '/' +
          context.rootState.tenants.tenant.id +
          '/policies',
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
                'Error fetching policies from server'
              )
            )
            return
          }

          const policies = await res.json()
          resolve(policies)
        })
        .catch((err) => {
          console.error(err)
          reject(new Error('An error occurred communicating with the server'))
        })
    })
  },
  getByID(context: any, policyID: string) {
    return new Promise((resolve, reject) => {
      fetch(
        process.env.baseUrl +
          '/' +
          context.rootState.tenants.tenant.id +
          '/policy/' +
          encodeURI(policyID),
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
              errorForStatus(context, res, 'Error fetching policy from server')
            )
            return
          }

          const policy = await res.json()
          resolve(policy)
        })
        .catch((err) => {
          console.error(err)
          reject(new Error('An error occurred communicating with the server'))
        })
    })
  },
  createPolicy(context: any, createPolicyRequest: any) {
    return new Promise((resolve, reject) => {
      fetch(
        process.env.baseUrl +
          '/' +
          context.rootState.tenants.tenant.id +
          '/policies',
        {
          method: 'POST',
          headers: new Headers({
            'Content-Type': 'application/json',
            Authorization:
              'Bearer ' + context.rootState.authentication.authToken,
          }),
          body: JSON.stringify(createPolicyRequest),
        }
      )
        .then(async (res) => {
          if (res.status !== 200 && res.status !== 204) {
            reject(
              errorForStatus(
                context,
                res,
                'Error creating new policy on server'
              )
            )
            return
          }

          const body = await res.json()
          resolve(body.policy_id)
        })
        .catch((err) => {
          console.error(err)
          reject(new Error('An error occurred communicating with the server'))
        })
    })
  },
  getScopeByID(context: any, policyID: string) {
    return new Promise((resolve, reject) => {
      fetch(
        process.env.baseUrl +
          '/' +
          context.rootState.tenants.tenant.id +
          '/policy/' +
          encodeURI(policyID) +
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
                'Error fetching policy scope from server'
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
  patchPolicy(context: any, params: any) {
    return new Promise((resolve, reject) => {
      fetch(
        process.env.baseUrl +
          '/' +
          context.rootState.tenants.tenant.id +
          '/policy/' +
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
              errorForStatus(context, res, 'Error patching policy on server')
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
          '/policies',
        {
          method: 'POST',
          headers: new Headers({
            'Content-Type': 'application/json',
            Authorization:
              'Bearer ' + context.rootState.authentication.authToken,
          }),
          body: JSON.stringify({
            policies: [params.pid],
          }),
        }
      )
        .then((res) => {
          if (res.status !== 200 && res.status !== 204) {
            reject(
              errorForStatus(
                context,
                res,
                'Error adding policy to group on server'
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
          '/policies',
        {
          method: 'DELETE',
          headers: new Headers({
            'Content-Type': 'application/json',
            Authorization:
              'Bearer ' + context.rootState.authentication.authToken,
          }),
          body: JSON.stringify({
            policies: [params.pid],
          }),
        }
      )
        .then((res) => {
          if (res.status !== 200 && res.status !== 204) {
            reject(
              errorForStatus(
                context,
                res,
                'Error removing policy from group on server'
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
  deletePolicy(context: any, id: string) {
    return new Promise((resolve, reject) => {
      fetch(
        process.env.baseUrl +
          '/' +
          context.rootState.tenants.tenant.id +
          '/policy/' +
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
              errorForStatus(context, res, 'Error deleting policy on server')
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
