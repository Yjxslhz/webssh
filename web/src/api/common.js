import axios from 'axios'

export function checkSSH(sshInfo) {
    return axios({
        url: '/check',
        method: 'get',
        params: {
            sshInfo: sshInfo
        }
    })
}

export function getDefaultConfig() {
    return axios({
        url: '/api/default-config',
        method: 'get'
    })
}

export function getServerList() {
    return axios({
        url: '/api/server-list',
        method: 'get'
    })
}
