interface State {
  error: Error | null
  errorTraceID: string | null
  menuActive: Boolean
  editting: Boolean | null
  deletable: Boolean
}

export const state = (): State => ({
  error: null,
  errorTraceID: null,
  menuActive: sessionStorage.getItem('menuActive') !== 'false',
  editting: null,
  deletable: false,
})

export const mutations = {
  setError(state: State, error: Error) {
    state.error = error
  },

  setErrorTraceID(state: State, errorTraceID: string) {
    state.errorTraceID = errorTraceID
  },

  clearError(state: State) {
    state.error = null
    state.errorTraceID = null
  },

  toggleMenuActive(state: State) {
    state.menuActive = !state.menuActive
    sessionStorage.setItem('menuActive', JSON.stringify(state.menuActive))
  },

  setEditting(state: State, editting: Boolean | null) {
    state.editting = editting
  },

  setDeletable(state: State, deletable: Boolean) {
    state.deletable = deletable
  },
}
