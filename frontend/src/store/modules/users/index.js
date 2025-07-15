// src/store/modules/users/index.js
// Vereinfachter Users Store - nur gRPC-Web

import mutations from './mutations.js'
import actions from './actions.js'
import getters from './getters.js'

export default {
  namespaced: true,

  state() {
    return {
      users: [],
      loading: false,
      creating: false,
      error: null
    }
  },

  mutations,
  actions,
  getters
}
