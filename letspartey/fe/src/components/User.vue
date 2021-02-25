<template>
    <div class="div-user-page">
        <el-card>

                <el-upload
                    class="avatar-uploader"
                    action="/api/v1/upload"
                    list-type="picture-card"
                    :headers="{token: token}"
                    :data="{username: userInfo.username, type: 'user'}"
                    :show-file-list="false"
                    :on-success="handleAvatarSuccess">
                    <img v-if="imageUrl" :src="imageUrl" class="avatar">
                    <i v-else class="el-icon-plus avatar-uploader-icon"></i>
                </el-upload>

                <p style="margin-top: 38px; font-size: 13px; text-align: left; color: #C0C4CC">{{ config.lang.opChangePass }}</p>
                <el-form style="margin-top: 8px" size="mini" label-width="80px">
                    <el-form-item :label="config.lang.password">
                        <el-input type="password" size="mini" v-model="password" prefix-icon="el-icon-unlock" />
                    </el-form-item>
                    <el-form-item :label="config.lang.passwordConfirm">
                        <el-input type="password" size="mini" v-model="passwordConfirm" prefix-icon="el-icon-unlock" />
                    </el-form-item>
                    <el-button round style="width: 100%" type="warning" icon="el-icon-view" size="mini" @click="updatePass">
                        {{ config.lang.OK }}
                    </el-button>
                </el-form>
        </el-card>
        <el-button round style="width: 100%; margin-top: 16px" type="danger" icon="el-icon-switch-button" size="medium" @click="logout">
            {{ config.lang.logout }}
        </el-button>
    </div>
</template>

<script>
import { getUserInfo, getToken, setUserInfo } from '../utils/utils'
import axios from 'axios'
import Config from '../utils/config'
import md5 from 'md5'

export default {
    beforeMount() {
        this.userInfo = getUserInfo()
        this.token = getToken()
        this.imageUrl = '/api/' + this.userInfo.avatar
    },
    props: {
        logoutMethod: {
            type: Function,
            default: null
        }
    },
    data() {
        return {
            imageUrl: null,
            userInfo: null,
            token: '',
            config: Config,
            password: '',
            passwordConfirm: '',
        }
    },
    methods: {
        handleAvatarSuccess(res, file) {
            this.imageUrl = URL.createObjectURL(file.raw);
            axios.get(`/api/v1/user?username=${this.userInfo.username}`)
            .then(res => {
                this.userInfo.avatar = res.data.avatar
                setUserInfo(this.userInfo)
            })
            .catch(err => {
                console.log(err)
            })
        },
        updatePass() {
            if (this.password != this.passwordConfirm || !this.password || !this.passwordConfirm) {
                this.$message({
                    type: 'warning',
                    message: Config.lang.passwordMismatch
                })
                return
            }

            axios.post('/api/v1/password', {
                username: this.userInfo.username,
                password: md5(this.password)
            })
            .then(() => {
                this.$message({
                    type: 'success',
                    message: Config.lang.OK
                })
            })
        },
        logout() {
            this.$confirm(Config.lang.confirmLogout, 'Warning', {
                confirmButtonText: 'OK',
                cancelButtonText: 'Cancel',
                type: 'warning'
            }).then(() => {
                if (this.logoutMethod) {
                    this.logoutMethod()
                }
            })
        }
    }
}
</script>

<style scoped>
.avatar-uploader .el-upload {
    width: 178px;
    height: 178px;
    border: 1px dashed #d9d9d9;
    border-radius: 6px;
    cursor: pointer;
    position: relative;
    overflow: hidden;
  }
  .avatar-uploader .el-upload:hover {
    border-color: #409EFF;
  }
  .avatar-uploader-icon {
    font-size: 28px;
    color: #8c939d;
    width: 148px;
    height: 148px;
    line-height: 148px;
    text-align: center;
  }
  .avatar {
    width: 148px;
    height: 148px;
    display: block;
  }
</style>