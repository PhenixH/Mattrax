import { errorForStatus } from './errors'

export interface CreateTenantRequest {
  display_name: string
  primary_domain: string
}

interface Tenant {
  id: string
  display_name: string
  primary_domain: string
}

interface State {
  tenant: Tenant | null
  tenants: Tenant[] | null
}

export const state = (): State => ({
  tenant:
    sessionStorage.getItem('tenant') !== null
      ? JSON.parse(sessionStorage.getItem('tenant') || '')
      : null,
  tenants: null,
})

export const mutations = {
  set(state: State, tenant: Tenant | null) {
    state.tenant = tenant
    if (tenant !== null) {
      sessionStorage.setItem('tenant', JSON.stringify(tenant))
    }
  },

  setTenants(state: State, tenants: Tenant[]) {
    state.tenants = tenants
  },
}

export const actions = {
  create(context: any, tenant: CreateTenantRequest) {
    return new Promise((resolve, reject) => {
      fetch(process.env.baseUrl + '/tenants', {
        method: 'POST',
        headers: new Headers({
          'Content-Type': 'application/json',
          Authorization: 'Bearer ' + context.rootState.authentication.authToken,
        }),
        body: JSON.stringify(tenant),
      })
        .then(async (res) => {
          if (res.status !== 200) {
            reject(
              errorForStatus(res, 'The create tenant request was rejected')
            )
            return
          }

          const data = await res.json()
          const t: Tenant = {
            id: data.tenant_id,
            display_name: tenant.display_name,
            primary_domain: tenant.primary_domain,
          }
          context.commit('set', t)
          resolve()
        })
        .catch((err) => {
          console.error(err)
          reject(new Error('An error occurred communicating with the server'))
        })
    })
  },

  getAll(context: any) {
    return new Promise((resolve, reject) => {
      fetch(process.env.baseUrl + '/tenants', {
        headers: new Headers({
          Authorization: 'Bearer ' + context.rootState.authentication.authToken,
        }),
      })
        .then(async (res) => {
          if (res.status !== 200) {
            reject(errorForStatus(res, 'The get tenants request was rejected'))
            return
          }

          const data = await res.json()
          console.log(data)
          context.commit('setTenants', data)
          resolve()
        })
        .catch((err) => {
          console.error(err)
          reject(new Error('An error occurred communicating with the server'))
        })
    })
  },
}
