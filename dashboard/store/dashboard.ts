interface State {
  error: Error | null
  menuActive: Boolean
  editting: Boolean | null
  deletable: Boolean
}

export const state = (): State => ({
  error: null,
  menuActive: sessionStorage.getItem('menuActive') !== 'false',
  editting: null,
  deletable: false,
})

export const mutations = {
  setError(state: State, error: Error) {
    state.error = error
  },

  clearError(state: State) {
    state.error = null
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
