// TODO: Pagination, Filters
import { errorForStatus } from './errors'

export const actions = {
  get(context: any, objectID: string) {
    return new Promise((resolve, reject) => {
      fetch(
        process.env.baseUrl +
          '/' +
          context.rootState.tenants.tenant.id +
          '/object/' +
          encodeURI(objectID),
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
              errorForStatus(context, res, 'Error fetching object from server')
            )
            return
          }

          const blob = await res.blob()
          const filename = res.headers.get('X-Filename') || encodeURI(objectID)
          resolve({ blob, filename })
        })
        .catch((err) => {
          console.error(err)
          reject(new Error('An error occurred communicating with the server'))
        })
    })
  },
  upload(context: any, params: any) {
    return new Promise((resolve, reject) => {
      const formData = new FormData()
      formData.append('file', params.file)

      fetch(
        process.env.baseUrl +
          '/' +
          context.rootState.tenants.tenant.id +
          '/object/' +
          encodeURI(params.id),
        {
          method: 'POST',
          headers: new Headers({
            Authorization:
              'Bearer ' + context.rootState.authentication.authToken,
          }),
          body: formData,
        }
      )
        .then((res) => {
          if (res.status !== 200 && res.status !== 204) {
            reject(
              errorForStatus(context, res, 'Error uploading object to server')
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
