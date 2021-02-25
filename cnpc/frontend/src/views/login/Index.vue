<template>
    <div class="login">
        <div class="login-con">
            <el-card shadow="never">
                <div slot="header" style="text-align: center">
<!--                    <h3>中石油销售大数据平台</h3>-->
                  <h3>非油业务数据分析系统</h3>
                </div>
                <el-form>
                    <el-form-item>
                        <el-input
                            placeholder="请输入用户名"
                            v-model="username"
                            style="margin-bottom: 18px"
                            @keyup.native.enter="login"
                        >
                        <i slot="prefix" class="fas fa-user-circle"></i>
                        </el-input>
                    </el-form-item>
                    <el-form-item>
                        <el-input
                            placeholder="请输入密码"
                            v-model="password"
                            type="password"
                            style="margin-bottom: 18px"
                            @keyup.native.enter="login"
                        >
                        <i slot="prefix" class="fas fa-lock"></i>
                        </el-input>
                    </el-form-item>
                    <el-form-item>
                        <el-button
                            type="primary" :loading="loginLoading"
                            style="width: 100%;margin-bottom: 18px"
                            @click.native="login"
                        >登录
                        </el-button>
                    </el-form-item>
                </el-form>
            </el-card>
        </div>
    </div>
</template>

<script>
    import { setToken, setUserInfo, setUserID } from '../../utils/dataStorage'
    import { login as loginApi } from '../../api/user'
    import Config from '../../config/index'
    import md5 from 'md5'

    export default {
        data() {
            return {
                username: '',
                password: '',
                remember: true,
                loginLoading: false,
            }
        },
        methods: {
            login() {
                if (this.username.length === 0 || this.password.length === 0) {
                    this.$message({
                        type: 'error',
                        'message': '用户名密码不能为空'
                    })
                } else {
                    loginApi({
                        username: this.username,
                        password: md5(this.password).toUpperCase()
                    }, res => {
                        setToken(res.token);
                        res.version = Config.version
                        setUserInfo(JSON.stringify(res));
                        setUserID(res.userID);
                        this.$router.push({path: Config.firstPage});
                    }, err => {
                        this.$message.error('登陆失败，请稍后再试。');
                    })
                }
            }
        }
    }
</script>

<style lang="scss" scoped>
    @import "Login.scss";
</style>
