import { getLanguage } from '@/lang/index'

export default {
    sshInfo: {
        host: '',
        username: 'root',
        port: 22,
        password: '',
        logintype: 0,
        keypassphrase: ''
    },
    sshList: Object.prototype.hasOwnProperty.call(localStorage, 'sshList') ? localStorage.getItem('sshList') : null,
    serverList: [],
    termList: [],
    currentTab: {},
    language: getLanguage()
}
