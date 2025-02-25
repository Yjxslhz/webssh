export default {
    SET_PASS: (state, password) => {
        state.sshInfo.password = password
    },
    SET_LIST: (state, sshList) => {
        state.sshList = sshList
        localStorage.setItem('sshList', sshList)
    },
    SET_TERMLIST: (state, termList) => {
        state.termList = termList
    },
    SET_SSH: (state, sshInfo) => {
        state.sshInfo.host = sshInfo.ipaddress
        state.sshInfo.username = sshInfo.username
        state.sshInfo.port = sshInfo.port
        state.sshInfo.password = sshInfo.password
        state.sshInfo.logintype = sshInfo.logintype
        state.sshInfo.keypassphrase = sshInfo.keypassphrase || ''
    },
    SET_TAB: (state, tab) => {
        state.currentTab = tab
    },
    SET_LANGUAGE: (state, language) => {
        state.language = language
        localStorage.setItem('language', language)
    },
    SET_SERVER_LIST: (state, serverList) => {
        state.serverList = serverList
    },
    SET_DEFAULT_CONFIG: (state, config) => {
        console.log('Setting default config:', config);
        // 只有在初始状态下才设置默认值
        if (state.sshInfo.host === '') {
            state.sshInfo.host = config.host || ''
            console.log('Set default host to:', state.sshInfo.host);
        }
        if (state.sshInfo.port === 22 && config.port !== 0) {
            state.sshInfo.port = config.port
        }
        if (state.sshInfo.username === 'root') {
            state.sshInfo.username = config.username || 'root'
        }
        if (state.sshInfo.password === '') {
            state.sshInfo.password = config.password || ''
        }
        if (state.sshInfo.logintype === 0) {
            state.sshInfo.logintype = config.login_type || 0
        }
        if (state.sshInfo.keypassphrase === '') {
            state.sshInfo.keypassphrase = config.key_passphrase || ''
        }
    }
}
