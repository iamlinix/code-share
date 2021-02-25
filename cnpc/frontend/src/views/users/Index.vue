<template>  
    
    <div class="cnpc-user-div">
        <el-dialog :title="currentUsername + ' 页面查看权限'" :visible.sync="showViewPerm"
            width="25%">
            <el-tree ref="userViewPermTree" :data="currentViewPerm" node-key="key" show-checkbox 
                :props="{label: 'name', children: 'children'}" 
                :default-checked-keys="currentCheckedKeys"
                @check-change="viewPermCheckChange" default-expand-all/>
            <el-button size="mini" round style="width: 100%; margin-top: 10px" type="success"
                @click="updateUserPerm">
                确定修改
            </el-button>
        </el-dialog>
        <el-dialog :title="currentUsername + ' 机构查看权限'" :visible.sync="showOrgPerm"
            width="40%">
            <el-table ref="userOrgPermTable" :data="currentOrgPerm" @selection-change="orgPermChange"
                max-height="450px">
                <el-table-column type="selection" width="50" />
                <el-table-column prop="orgCode" label="机构编码" width="80"/>
                <el-table-column prop="orgText" label="机构名称" width="400" show-overflow-tooltip/>
            </el-table>
            <el-button size="mini" round style="width: 100%; margin-top: 10px" type="success"
                @click="updateUserOrgPerm">
                确定修改
            </el-button>
        </el-dialog>
        <p class="cnpc-user-name">{{ userInfo }}
            <i style="color: #E6A23C; cursor: pointer; margin-left: 12px" class="el-icon-view" 
                @click="showChangePassword = !showChangePassword"/>
            <i v-if="isAdmin" style="color: #67C23A; cursor: pointer; margin-left: 4px" class="el-icon-user-solid" 
                @click="showNewUser = !showNewUser"/>
        </p>
        <transition name="zoom">
            <el-card v-if="showChangePassword" style="animation-duration: 0.2s">
                <el-form size="mini" label-position="right" label-width="90px">
                    <el-row :gutter="12">
                        <el-col :span="6">
                            <el-form-item style="margin-bottom: 0" label="现有密码">
                                <el-input placeholder="请输入现有密码" v-model="currentPass" size="mini" 
                                    type="password" clearable/>
                            </el-form-item>
                        </el-col>
                        <el-col :span="6">
                            <el-form-item style="margin-bottom: 0" label="新密码">
                                <el-input placeholder="请输入新密码" v-model="newPass" size="mini" 
                                    type="password" clearable/>
                            </el-form-item>
                        </el-col>
                        <el-col :span="6">
                            <el-form-item style="margin-bottom: 0" label="确认新密码">
                                <el-input placeholder="请确认新密码" v-model="newPassConfirm" size="mini" 
                                    type="password" clearable/>
                            </el-form-item>
                        </el-col>
                        <el-col :span="3">
                            <el-button type="primary" size="mini" icon="el-icon-view"
                                @click="resetMyPassword">修改密码</el-button>
                        </el-col>
                    </el-row>
                </el-form>
            </el-card>
        </transition>
        <transition name="zoom">
            <el-card style="margin-top: 10px; animation-duration: 0.2s" v-if="isAdmin && showNewUser">
                <el-form size="mini" label-position="right" label-width="90px">
                    <el-row :gutter="12">
                        <el-col :span="6">
                            <el-form-item style="margin-bottom: 0" label="用户名">
                                <el-input placeholder="请输入用户名" v-model="newUser" size="mini" clearable/>
                            </el-form-item>
                        </el-col>
                        <el-col :span="6">
                            <el-form-item style="margin-bottom: 0" label="密码">
                                <el-input placeholder="请输入密码" v-model="newUserPass" size="mini" 
                                    type="password" clearable/>
                            </el-form-item>
                        </el-col>
                        <el-col :span="3">
                            <el-checkbox v-model="newUserAdmin">管理员用户</el-checkbox>
                        </el-col>
                        <el-col :span="6">
                            <el-button size="mini" type="success" icon="el-icon-user-solid"
                                @click="createUser">添加用户</el-button>
                        </el-col>
                    </el-row>
                </el-form>
            </el-card>
        </transition>
        <el-card style="margin-top: 10px" v-if="isAdmin">
            <s-table :data="userList" :minus="290" :loading="userTableLoading" ref="cnpc-user-table"
                size="small">
                <vxe-table-column title="用户名" field="username" sortable/>
                <vxe-table-column title="创建时间" field="createTime" />
                <vxe-table-column title="操作">
                    <template slot-scope="r">
                        <div v-if="r.row.username != 'admin'">
                            <!--el-input v-if="r.row.resetting" size="mini"
                                v-model="r.row.pass" 
                                placeholder="请输入新密码" clearable
                                style="width: 240px; margin-right: 8px">
                                <el-button slot="append" icon="el-icon-success" size="mini" 
                                    style="color: green; width: 100%"
                                    @click="resetPassword(r.row)"/>
                                <el-button slot="prepend" icon="el-icon-error" size="mini" 
                                    style="color: red; width: 100%"
                                    @click="cancelResetPassword(r.row)"/>
                            </el-input-->
                            <span v-if="r.row.resetting" class="reset-password">{{ r.row.resetPassword }}
                                <i class="el-icon-close" style="cursor: pointer" @click="r.row.resetting = false"/>
                            </span>
                            <el-button v-else type="text" size="mini"
                                @click="resetPassword(r.row)">重置密码</el-button>
                            <el-button type="text" size="mini"
                                style="color: red" @click="deleteUser(r.row.id, r.row.username)">
                                删除用户</el-button>
                            <el-button type="text" size="mini" v-if="r.row.role == 0" style="color: #67C23A" 
                                @click="updateRole(r.row.username, 1)">
                                设为管理员</el-button>
                            <el-button type="text" size="mini" v-else style="color: #E6A23C" 
                                @click="updateRole(r.row.username, 0)">
                                取消管理员</el-button>
                            <el-button type="text" size="mini" style="color: #409EFF" 
                                @click="editViewPerm(r.row.username, r.row.viewPerm)">
                                页面权限</el-button>
                            <el-button type="text" size="mini" style="color: #409EFF" 
                                @click="editOrgPerm(r.row.username, r.row.orgPerm)">
                                机构权限</el-button>
                        </div>
                    </template>
                </vxe-table-column>
            </s-table>
        </el-card>
    </div>
</template>

<style lang="scss">
.cnpc-user-div {
    background-color: #FFFFFF;
    max-height: calc(100vh - 10px);
    overflow: auto;
    padding: 8px 8px 8px 8px;
    width: 100%;
    position: relative;

    p.cnpc-user-name {
        font-size: 20px;
        font-weight: bold;
        padding: 10px 0 20px 10px;
    }

    .el-input-group__prepend {
        padding: 0 8px;
    }

    .el-input-group__append {
        padding: 0 8px;
    }
    
    span.reset-password {
        margin-right: 8px;
        font-size: 13px;
        color: #909399;
        padding: 4px 8px;
        border-radius: 5px;
        background-color: #DCDFE6;
    }
}
</style>

<script>
import { getUserInfo, getUserID } from '../../utils/dataStorage'
import SmartTable from '../mixins/SmartVxeTable'
import { doRequest, message, confirm, amIAdmin, deepCopy, mergeNavigation } from '../../utils/utils'
import { logout } from '../../api/user'
import md5 from 'md5'
require('vue2-animate/dist/vue2-animate.min.css');
import Config from '../../config/index'
import Navigation from '@/navigation/index'

export default {
    components: {
        's-table': SmartTable
    },
    mounted() {
        this.getOrgList()
        this.userInfo = getUserInfo().name;
        this.userId = getUserID();
        if (this.userId < 0) {
            message('error', '您的登录信息有误, 请重新登录')
        }
        this.isAdmin = amIAdmin();

        if (this.isAdmin)
            this.getUserList();
    },
    data() {
        return {
            userId: -1,
            userInfo: null,
            newUser: '',
            newUserPass: '',
            newUserAdmin: false,
            currentPass: '',
            newPass: '',
            newPassConfirm: '',
            userList: [],
            isAdmin: false,
            orgList: [],

            userTableLoading: false,
            
            showChangePassword: true,
            showNewUser: true,

            currentUsername: '',
            currentViewPerm: [],
            currentOrgPerm: [],
            changedOrgPerm: [],
            currentCheckedKeys: [],
            realViewPerm: null,
            realOrgPerm: null,
            showViewPerm: false,
            showOrgPerm: false,
        }
    },
    methods: {
        getOrgList() {
            doRequest({
                method: 'GET',
                url: '/v1/web/org/plant',
                loading: true
            }, {
                success: res => {
                    this.orgList = [{
                        orgCode: '000',
                        orgText: '北京公司'
                    }, {
                        orgCode: 'A13',
                        orgText: '北京销售公司一分公司'
                    }, {
                        orgCode: 'A14',
                        orgText: '北京销售公司二分公司'
                    }, {
                        orgCode: 'A15',
                        orgText: '北京销售公司三分公司'
                    }, {
                        orgCode: 'A16',
                        orgText: '北京销售公司四分公司'
                    }]
                    res.orgList.forEach(e => {
                        this.orgList.push(e)
                    })
                }
            })
        },
        orgPermChange(sel) {
            this.changedOrgPerm = sel
        },
        editOrgPerm(username, perm) {
            this.currentUsername = username
            this.currentOrgPerm = deepCopy(this.orgList)
            if (perm.length > 0) {
                this.currentOrgPerm = JSON.parse(perm)
            }
           
            // console.log(this.$refs)
            // this.currentOrgPerm.forEach(e => {
            //     this.$refs.userOrgPermTable.toggleRowSelection(e, e.show || (e.show == null))
            // })
            this.showOrgPerm = true
            this.$nextTick(() => {
                this.$refs.userOrgPermTable.clearSelection()
                this.currentOrgPerm.forEach(e => {
                    this.$refs.userOrgPermTable.toggleRowSelection(e, e.show || (e.show == null))
                })
            })
        },
        updateUserOrgPerm() {
            confirm('warning', '修改机构权限', `您确定要修改 ${this.currentUsername} 的机构查看权限么？`, _ => {
                let orgPerm = deepCopy(this.currentOrgPerm)
                orgPerm.forEach(p => {
                    let b = false
                    this.changedOrgPerm.forEach(c => {
                        if (c.orgCode == p.orgCode) {
                            p.show = 1
                            b = true
                        }
                    })
                    if (!b) {
                        p.show = 0
                    }
                })
                doRequest({
                    method: 'POST',
                    url: '/v1/web/user/org-perm',
                    data: {
                        username: this.currentUsername,
                        orgPerm: JSON.stringify(orgPerm)
                    },
                    loading: true
                }, {
                    success: _ => {
                        message('success', '修改机构权限成功')
                        this.showOrgPerm = false
                        this.getUserList()
                    },
                    fail: _ => {
                        message('success', '修改机构权限失败, 请稍后再试')
                    }
                })
            })
        },
        updateUserPerm() {
            confirm('warning', '修改页面权限', `您确定要修改 ${this.currentUsername} 的页面查看权限么？`, _ => {
                doRequest({
                    method: 'POST',
                    url: '/v1/web/user/view-perm',
                    data: {
                        username: this.currentUsername,
                        viewPerm: JSON.stringify(this.realViewPerm)
                    },
                    loading: true
                }, {
                    success: _ => {
                        message('success', '修改页面权限成功')
                        this.showViewPerm = false
                        this.getUserList()
                    },
                    fail: _ => {
                        message('success', '修改页面权限失败, 请稍后再试')
                    }
                })
            })
        },
        viewPermCheckChange(item, me, kids) {
            console.log(item, me, kids)
            if (item.children) {
                if (kids)
                    item.show = me ? 1 : 0
                else
                    item.show = me ? 1 : -1
            } else {
                item.show = me ? 1 : -1
            }

            this.currentViewPerm.forEach(p => {
                if (p.children) {
                    let n = p.children.length, m = 0
                    p.children.forEach(c => {
                        if (c.show == 1) {
                            m ++
                        }
                    })

                    if (m == 0) {
                        p.show = -1
                    } else if (m == n){
                        p.show = 1
                    } else {
                        p.show = 0
                    }
                }
                // if (p.show >= 0 && p.children) {
                //     let b = false
                //     p.children.forEach(k => {
                //         if (k.show == 1) {
                //             b = true
                //         }
                //     })
                //     if (!b)
                //         p.show = -1;
                // }
            })
        },
        editViewPerm(username, perm) {
            this.currentUsername = username
            let viewPerm = Navigation
            if (perm.length > 0) {
                viewPerm = mergeNavigation(Navigation, JSON.parse(perm))
            }
            this.currentViewPerm = []
            this.realViewPerm = viewPerm
            for (let k in viewPerm) {
                this.currentViewPerm.push(viewPerm[k])
            }

            this.currentCheckedKeys = []
            this.currentViewPerm.forEach(e => {
                if (e.show == 1)
                    this.currentCheckedKeys.push(e.key)
                if (e.children) {
                    e.children.forEach(c => {
                        if (c.show == 1)
                            this.currentCheckedKeys.push(c.key)
                    })
                }
            })
            this.showViewPerm = true
        },
        updateRole(name, role) {
            confirm('warning', '更改管理员', role == 0 ? `您确认要取消 ${name} 的管理员权限么？` : `您确认要将 ${name} 设为管理员么？`, _ => {
                doRequest({
                    method: 'POST',
                    url: '/v1/web/user/role-update',
                    data: {
                        username: name,
                        role: role
                    },
                    loading: true
                }, {
                    success: _ => {
                        message('success', '设置用户权限成功')
                        this.getUserList()
                    },
                    fail: _ => {
                        message('error', '设置用户权限失败，请稍后再试')
                    }
                })
            })
        },
        getUserList() {
            this.userTableLoading = true;
            doRequest({
                url: '/v1/web/user/list',
                method: 'GET'
            }, {
                obj: this,
                src: 'users',
                dst: 'userList',
                fail: err => {
                    message('error', '获取用户列表失败')
                },
                finally: d => {
                    this.userTableLoading = false;
                }
            })
        },
        createUser() {
            if (this.newUser.length > 0 && this.newUserPass.length > 0) {
                doRequest({
                    url: '/v1/web/user/add',
                    method: 'POST',
                    loading: true,
                    data: {
                        username: this.newUser,
                        password: md5(this.newUserPass).toUpperCase(),
                        role: this.newUserAdmin ? 1 : 0
                    }
                }, {
                    success: res => {
                        message('success', `用户 ${this.newUser} 创建成功`)
                        this.newUser = '';
                        this.newUserPass = '';
                        this.newUserAdmin = false;
                        this.getUserList();
                    },
                    fail: err => {
                        message('error', `用户 ${this.newUser} 创建失败`)
                    }
                })
            } else {
                message('error', '请输入有效的用户名密码')
            }
        },
        showResetPassword(row) {
            this.$set(row, 'resetting', true);
        },
        cancelResetPassword(row) {
            this.$set(row, 'resetting', false);
        },
        resetPassword(row) {
            confirm('warning', '重置密码', `确定要重置用户 ${row.username} 的密码么?`, () => {
                this.userTableLoading = true;
                doRequest({
                    url: '/v1/web/user/reset', 
                    method: 'POST',
                    data: {
                        id: row.id,
                        username: row.username,
                    }
                }, {
                    success: res => {
                        message('success', '密码重置成功')
                        row.resetPassword = res.message;
                        this.$set(row, 'resetting', true);
                    },
                    fail: err => {
                        message('error', '密码重置失败, 请稍后再试')
                    },
                    finally: _ => {
                        this.userTableLoading = false;
                    }
                })
            })
        },
        deleteUser(userId, username) {
            confirm('warning', '删除用户', `确定要删除用户 ${username} 么?`, () => {
                doRequest({
                    url: `/v1/web/user/del/${userId}`,
                    method: 'GET',
                    loading: true
                }, {
                    success: res => {
                        message('success', `用户 ${username} 删除成功`)
                        this.getUserList();
                    },
                    fail: err => {
                        message('error', `用户 ${username} 删除失败`)
                    }
                })   
            })
        },
        resetMyPassword() {
            if (this.currentPass.length > 0 && this.newPass.length > 0 && this.newPassConfirm.length > 0) {
                if (this.newPass !== this.newPassConfirm) {
                    message('error', '新密码不一致，请再次确认')
                } else {
                    doRequest({
                        url: '/v1/web/user/update',
                        method: 'POST',
                        loading: true,
                        data: {
                            id: this.userId,
                            username: this.userInfo,
                            password: md5(this.currentPass).toUpperCase(),
                            newPassword: md5(this.newPass).toUpperCase()
                        }
                    }, {
                        success: res => {
                            message('success', `密码更新成功，请重新登录系统`)
                            setTimeout(() => {
                                logout()
                                // window.location.href = '/login'
                                Config.router.push({path: '/login'})
                            }, 2000);
                        },
                        fail: err => {
                            message('error', `密码更新失败`)
                        }
                    })
                }
            } else {
                message('error', '请输入有效的密码')
            }
        }
    }
}
</script>