import { getLanguage } from '@/lang/index'
import { getDefaultConfig, getServerList } from '@/api/common'

export default {
    setLanguage({ commit }, language) {
        commit('SET_LANGUAGE', language)
    },
    loadDefaultConfig({ commit }) {
        return new Promise((resolve, reject) => {
            getDefaultConfig().then(response => {
                const defaultConfig = response.data
                console.log('Received default config from server:', defaultConfig)
                commit('SET_DEFAULT_CONFIG', defaultConfig)
                resolve(defaultConfig)
            }).catch(error => {
                console.error('Failed to load default config:', error)
                reject(error)
            })
        })
    },
    loadServerList({ commit }) {
        return new Promise((resolve, reject) => {
            getServerList().then(response => {
                const serverList = response.data
                commit('SET_SERVER_LIST', serverList)
                resolve(serverList)
            }).catch(error => {
                console.error('Failed to load server list:', error)
                reject(error)
            })
        })
    }
}
