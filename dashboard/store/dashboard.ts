interface State {
  error: Error | null
  menuActive: Boolean
}

export const state = (): State => ({
  error: null,
  menuActive: sessionStorage.getItem('menuActive') !== 'false',
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
}
