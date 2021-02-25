<template>
    <div class="div-login-page">
        <el-card>
            <h2 class="login-banner">
                {{ config.lang.welcome }}
            </h2>
            <el-tabs v-model="activeTab">
                <el-tab-pane :label="config.lang.signin" name="signin">
                    <el-form label-width="100px">
                        <el-form-item :label="config.lang.username">
                            <el-input v-model="userName" prefix-icon="el-icon-user"/>
                        </el-form-item>
                        <el-form-item :label="config.lang.password">
                            <el-input v-model="password" prefix-icon="el-icon-unlock" type="password"/>
                        </el-form-item>
                        <el-button style="width: 100%" round type="primary" @click="login">
                            {{ config.lang.login }}
                        </el-button>
                    </el-form>
                </el-tab-pane>
                <el-tab-pane :label="config.lang.signup" name="signup">
                    <el-form label-width="100px">
                        <el-form-item :label="config.lang.username">
                            <el-input v-model="userName" prefix-icon="el-icon-user"/>
                        </el-form-item>
                        <el-form-item :label="config.lang.password">
                            <el-input v-model="password" prefix-icon="el-icon-unlock" type="password"/>
                        </el-form-item>
                        <el-form-item :label="config.lang.passwordConfirm">
                            <el-input v-model="passwordConfirm" prefix-icon="el-icon-unlock" type="password"/>
                        </el-form-item>
                        <el-button style="width: 100%" round type="primary" @click="signup">
                            {{ config.lang.signup }}
                        </el-button>
                    </el-form>
                </el-tab-pane>
            </el-tabs>
        </el-card>
    </div>
</template>

<script>
import Config from '../utils/config'

export default {
    props: {
        loginMethod: {
            type: Function,
            default: undefined
        },
        signupMethod: {
            type: Function,
            default: undefined
        }
    },
    data() {
        return {
            activeTab: 'signin',
            userName: '',
            password: '',
            passwordConfirm: '',
            config: Config
        }
    },
    methods: {
        login() {
            if (this.userName.length === 0 || this.password.length === 0) {
                this.$message({
                    type: 'error',
                    message: Config.lang.errUserPassEmpty
                })
                return
            }

            if (this.loginMethod) {
                this.loginMethod(this.userName, this.password)
            } else {
                console.log("login method not mouonted")
            }
        },
        signup() {
            if (!this.userName || !this.password || !this.passwordConfirm) {
                this.$message({
                    type: 'error',
                    message: Config.lang.errUserPassEmpty
                })
                return
            }

            if (this.password != this.passwordConfirm) {
                this.$message({
                    type: 'error',
                    message: Config.lang.passwordMismatch
                })
                return
            }

            if (this.signupMethod) {
                this.signupMethod(this.userName, this.password, () => {
                    this.activeTab = "signin"
                    this.password = this.passwordConfirm = ''
                })
            } else {
                console.log("signup method not mouonted")
            }
        }
    }
}
</script>

<style scoped>
.div-login-page {
    width: 100%;
    height: 100%;  
}

.div-login-page .login-banner {
    font-family: "europa, sans-serif";
}
</style>