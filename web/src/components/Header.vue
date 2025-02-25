<template>
    <el-row>
        <el-col>
            <el-form :inline="true" style="padding-top: 10px;" :model="sshInfo" :rules="checkRules">
                <el-form-item size="small" prop="host">
                    <template slot="label">
                        <el-tooltip effect="dark" placement="left">
                            <div slot="content">
                                <p>Switch Language</p>
                            </div>
                            <span @click="handleSetLanguage()">Host</span>
                        </el-tooltip>
                    </template>
                    <el-input v-model="sshInfo.host" :placeholder="$t('hostTip')" @keyup.enter.native="$emit('ssh-select')"></el-input>
                </el-form-item>
                <el-form-item label="Port" size="small" prop="port">
                    <el-input v-model="sshInfo.port" :placeholder="$t('portTip')" @keyup.enter.native="$emit('ssh-select')" style="width: 100px"></el-input>
                </el-form-item>
                <el-form-item label="Username" size="small" prop="username">
                    <el-input v-model="sshInfo.username" :placeholder="$t('nameTip')" @keyup.enter.native="$emit('ssh-select')" style="width: 110px"></el-input>
                </el-form-item>
                
                <!-- 登录方式选择 -->
                <el-form-item label="Auth Type" size="small">
                    <div>
                        <el-button 
                            size="small" 
                            :type="sshInfo.logintype === 0 ? 'primary' : 'default'" 
                            @click="setLoginType(0)">
                            {{ $t('Password') }}
                        </el-button>
                        <el-button 
                            size="small" 
                            :type="sshInfo.logintype === 1 ? 'primary' : 'default'" 
                            @click="setLoginType(1)">
                            {{ $t('PrivateKey') }}
                        </el-button>
                        <el-tooltip effect="dark" content="Help" placement="top">
                            <el-button type="text" icon="el-icon-question" @click="helpVisible = true"></el-button>
                        </el-tooltip>
                    </div>
                </el-form-item>
                
                <!-- 密码或私钥输入 -->
                <el-form-item size="small" prop="password">
                    <template slot="label">
                        <span>{{ privateKey ? $t('PrivateKey') : $t('Password') }}</span>
                    </template>
                    <div style="display: flex; align-items: center;">
                        <el-input v-model="sshInfo.password" @focus.native="handlePasswordClick" @keyup.enter.native="$emit('ssh-select')" :placeholder="$t('inputTip') + `${privateKey ? $t('privateKey') : $t('password')}`" show-password></el-input>
                        <el-button v-if="privateKey" type="text" icon="el-icon-upload2" @click="textareaVisible = true" style="margin-left: 5px;" :title="$t('UploadKey')"></el-button>
                    </div>
                </el-form-item>
                
                <!-- 私钥对话框 -->
                <el-dialog :title="$t('privateKey')" :visible.sync="textareaVisible" :close-on-click-modal="false" width="60%" @close="handleDialogClose">
                    <el-input :rows="12" v-model="sshInfo.password" type="textarea" :placeholder="$t('keyTip')"></el-input>
                    
                    <!-- 添加私钥密码输入框 -->
                    <div style="margin-top: 10px;">
                        <el-form-item :label="$t('KeyPassphrase')" size="small">
                            <el-input v-model="sshInfo.keypassphrase" type="password" :placeholder="$t('EnterPassphrase')"></el-input>
                        </el-form-item>
                    </div>
                    
                    <div slot="footer" class="dialog-footer">
                        <!-- 选择密钥文件 -->
                        <input ref="pkFile" @change="handleChangePKFile" type="file" style="position: absolute;clip: rect(0 0 0 0)"/>
                        <el-button type="primary" plain @click="$refs.pkFile.click()">{{ $t('SelectFile') }}</el-button>
                        <el-button @click="sshInfo.password=''">{{ $t('Clear') }}</el-button>
                        <el-button type="primary" @click="handleConnectWithKey">{{ $t('Connect') }}</el-button>
                    </div>
                </el-dialog>
                
                <!-- 帮助对话框 -->
                <el-dialog :title="$t('SSHAuthHelp')" :visible.sync="helpVisible" width="60%">
                    <div v-if="$i18n.locale === 'en'">
                        <h3>Password Authentication</h3>
                        <p>Select "Password" as Auth Type and enter your SSH password.</p>
                        
                        <h3>Private Key Authentication</h3>
                        <p>To use SSH key authentication:</p>
                        <ol>
                            <li>Select "Private Key" as Auth Type</li>
                            <li>Click on the password field to open the private key dialog</li>
                            <li>Either paste your private key content or click "Select File" to upload your private key file (e.g., id_rsa)</li>
                            <li>If your private key is passphrase protected, enter the passphrase in the "Key Passphrase" field</li>
                            <li>Click "Connect" to establish the SSH connection</li>
                        </ol>
                        
                        <h4>Notes:</h4>
                        <ul>
                            <li>The private key should be in OpenSSH format</li>
                            <li>Private keys with passphrases are supported - enter the passphrase in the provided field</li>
                            <li>Make sure your private key corresponds to a public key authorized on the server</li>
                        </ul>
                        
                        <h4>Example of a private key format:</h4>
                        <pre style="background-color: #f5f5f5; padding: 10px; border-radius: 5px; overflow-x: auto;">
-----BEGIN OPENSSH PRIVATE KEY-----
b3BlbnNzaC1rZXktdjEAAAAABG5vbmUAAAAEbm9uZQAAAAAAAAABAAABlwAAAAdzc2gtcn
NhAAAAAwEAAQAAAYEAvtPDQXr8...
...
-----END OPENSSH PRIVATE KEY-----</pre>
                    </div>
                    <div v-else>
                        <h3>密码认证</h3>
                        <p>选择"密码"作为认证类型，并输入您的SSH密码。</p>
                        
                        <h3>私钥认证</h3>
                        <p>使用SSH密钥认证：</p>
                        <ol>
                            <li>选择"私钥"作为认证类型</li>
                            <li>点击密码字段打开私钥对话框</li>
                            <li>粘贴您的私钥内容或点击"选择文件"上传您的私钥文件（例如id_rsa）</li>
                            <li>如果您的私钥有密码保护，请在"私钥密码"字段中输入密码</li>
                            <li>点击"连接"建立SSH连接</li>
                        </ol>
                        
                        <h4>注意事项：</h4>
                        <ul>
                            <li>私钥应为OpenSSH格式</li>
                            <li>支持带密码的私钥 - 在提供的字段中输入密码</li>
                            <li>确保您的私钥对应于服务器上授权的公钥</li>
                        </ul>
                        
                        <h4>私钥格式示例：</h4>
                        <pre style="background-color: #f5f5f5; padding: 10px; border-radius: 5px; overflow-x: auto;">
-----BEGIN OPENSSH PRIVATE KEY-----
b3BlbnNzaC1rZXktdjEAAAAABG5vbmUAAAAEbm9uZQAAAAAAAAABAAABlwAAAAdzc2gtcn
NhAAAAAwEAAQAAAYEAvtPDQXr8...
...
-----END OPENSSH PRIVATE KEY-----</pre>
                    </div>
                </el-dialog>
                
                <el-form-item size="small">
                    <el-button type="primary" @click="$emit('ssh-select')" plain>{{ $t('Connect') }}</el-button>
                </el-form-item>
                <el-form-item size="small">
                    <file-list></file-list>
                </el-form-item>
                <el-form-item size="small">
                    <el-dropdown @command="handleCommand">
                        <el-button type="primary">
                            {{ $t('History') }}
                        </el-button>
                        <el-dropdown-menu slot="dropdown">
                            <el-dropdown-item
                                v-for="item in sshList"
                                :key="item.host" :command="item" style="padding:0 5px 0 10px">
                                {{item.host}}
                                <i @click="cleanHistory(item)" class="el-icon-close"></i>
                            </el-dropdown-item>
                        </el-dropdown-menu>
                    </el-dropdown>
                </el-form-item>
                <!-- 服务器列表下拉菜单 -->
                <el-form-item size="small" v-if="serverList.length > 0">
                    <el-dropdown @command="handleServerSelect">
                        <el-button type="primary">
                            {{ $t('Servers') }}
                        </el-button>
                        <el-dropdown-menu slot="dropdown">
                            <el-dropdown-item
                                v-for="server in serverList"
                                :key="server.host" 
                                :command="server">
                                {{server.name || server.host}}
                            </el-dropdown-item>
                        </el-dropdown-menu>
                    </el-dropdown>
                </el-form-item>
            </el-form>
        </el-col>
    </el-row>
</template>

<script>
import { getLanguage } from '@/lang/index'
import FileList from '@/components/FileList'
import { mapState } from 'vuex'

export default {
    components: {
        'file-list': FileList
    },
    data() {
        return {
            textareaVisible: false,
            helpVisible: false,
            checkRules: {
                host: [
                    { required: true, trigger: 'blur' }
                ],
                port: [
                    { required: true, trigger: 'blur', type: 'number', transform(value) { return Number(value) } }
                ],
                username: [
                    { required: true, trigger: 'blur' }
                ],
                password: [
                    { required: true, trigger: 'blur', message: 'value is required' }
                ]
            }
        }
    },
    methods: {
        setLoginType(type) {
            // 确保logintype是数字类型
            this.sshInfo.logintype = Number(type);
            console.log('Set login type to:', this.sshInfo.logintype, typeof this.sshInfo.logintype);
            
            if (this.sshInfo.logintype === 1) {
                // 私钥模式
                console.log('Set to key mode');
                this.textareaVisible = true;
                // 清空密码字段
                this.sshInfo.password = '';
                // 清空密钥密码
                this.sshInfo.keypassphrase = '';
            } else {
                // 密码模式
                console.log('Set to password mode');
                // 清空密钥密码
                this.sshInfo.keypassphrase = '';
            }
        },
        handleSetLanguage() {
            const oldLang = getLanguage()
            const lang = oldLang === 'zh' ? 'en' : 'zh'
            this.$i18n.locale = lang
            this.$store.dispatch('setLanguage', lang)
        },
        handleCommand(command) {
            this.$store.commit('SET_SSH', command)
            if (command.password === undefined) {
                this.$store.commit('SET_PASS', '')
            }
            // 新开窗口
            this.$emit('ssh-select')
        },
        cleanHistory(command) {
            const sshListObj = this.sshList
            sshListObj.forEach((v, i) => {
                if (v.host === command.host) {
                    sshListObj.splice(i, 1)
                }
            })
            this.$store.commit('SET_LIST', window.btoa(JSON.stringify(sshListObj)))
        },
        handleChangePKFile(event) {
            const file = event.target.files[0]
            if (file) {
                const sshInfo = this.sshInfo
                const reader = new FileReader()
                reader.onload = e => {
                    sshInfo.password = e.target.result
                }
                reader.readAsText(file)
            }
        },
        handlePasswordClick() {
            // 如果是私钥模式，点击密码框时打开私钥输入对话框
            if (this.privateKey) {
                this.textareaVisible = true
                console.log('Opening private key dialog, logintype:', this.sshInfo.logintype)
            }
        },
        handleConnectWithKey() {
            console.log('Connecting with key, passphrase length:', this.sshInfo.keypassphrase.length);
            this.textareaVisible = false;
            this.$emit('ssh-select');
        },
        handleDialogClose() {
            console.log('Dialog closed, keypassphrase:', this.sshInfo.keypassphrase);
        },
        handleServerSelect(server) {
            // 将服务器信息转换为sshInfo格式
            const sshInfo = {
                ipaddress: server.host,
                username: server.username,
                port: server.port,
                password: server.password || '',
                logintype: server.login_type || 0,
                keypassphrase: server.key_passphrase || ''
            }
            
            // 设置SSH信息
            this.$store.commit('SET_SSH', sshInfo)
            
            // 新开窗口
            this.$emit('ssh-select')
        }
    },
    mounted() {
        console.log('Header mounted, current sshInfo:', JSON.stringify(this.sshInfo));
        
        if (this.sshList.length > 0) {
            const latestSSH = this.sshList[this.sshList.length - 1]
            this.$store.commit('SET_SSH', latestSSH)
            if (latestSSH.password === undefined) {
                this.$store.commit('SET_PASS', '')
            }
        }
    },
    computed: {
        ...mapState(['sshInfo']),
        privateKey() {
            return this.sshInfo.logintype === 1
        },
        sshList() {
            const sshList = this.$store.state.sshList
            if (sshList === null) {
                return []
            } else {
                return JSON.parse(window.atob(sshList))
            }
        },
        serverList() {
            return this.$store.state.serverList
        }
    },
    watch: {
        'sshInfo.logintype': function(newVal) {
            // 当切换到私钥模式时，清空密码字段
            console.log('logintype changed:', newVal, typeof newVal);
            if (parseInt(newVal) === 1) {
                this.sshInfo.password = '';
                // 自动打开私钥对话框
                this.textareaVisible = true;
            }
        }
    }
}
</script>

<style scoped>
pre {
    white-space: pre-wrap;
    word-wrap: break-word;
}
</style>
