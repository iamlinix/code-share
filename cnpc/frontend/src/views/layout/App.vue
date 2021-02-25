<!--suppress JSUnusedLocalSymbols -->
<template>
    <div class="main">
        <el-dialog :visible.sync="showAbout" title="关于">
            <h2><strong>{{ $Config.siteName }}</strong></h2>
            <br/>
            <h2><strong>V1.0.0</strong></h2>
        </el-dialog>
        <!--div class="header">
            <div class="logo">
                <span id="head-logo" class="min">
                    <img width="40" style="margin-top: 5px" src="../../assets/images/logo-45.png" alt="">
                </span>
            </div>
            <span class="header-btn" @click="hiddenSidebar">
                <i v-if="sidebarHidden" class="el-icon-d-arrow-right"></i>
                <i v-else class="el-icon-d-arrow-left"></i>
            </span>
            <h3 class="header-btn" @click="about"><strong>{{ $Config.siteName }}</strong></h3>
            <div class="right">
                <span class="header-btn" style="font-size: 13px" @click="logout">
                    <i class="el-icon-switch-button"> 退出系统</i>
                </span>
            </div>
        </div-->
        <div class="app">
            <div class="aside">
                <div class="logo">
                    <span id="head-logo" class="min">
                        <a href="/promotions">
                            <img width="40" style="margin-top: 5px" src="../../assets/images/logo-45.png" alt="" />
                        </a>
                    </span>
                </div>
                <div class="menu">
                    <el-menu
                            router
                            background-color="#222d32"
                            text-color="#fff"
                            :default-active="$route.path" class="menu" @open="handleOpen" @close="handleClose"
                            :collapse="isCollapse">
                        <template v-for="(nav_v, nav_k) in navigation">
                            <el-submenu :key="nav_v.key" v-if="shouldShowMenu(nav_v) && nav_v.children" :index="nav_k">
                                <template slot="title">
                                    <i :class="nav_v.icon"></i>
                                    <span>{{ nav_v.name }}</span>
                                </template>
                                <el-menu-item v-for="(ch_v, ch_k) in nav_v.children" :key="ch_k" :index="ch_v.path" 
                                    v-show="shouldShowMenu(ch_v)">
                                    <i :class="ch_v.icon"></i>
                                    <span slot="title">{{ ch_v.name }}</span>
                                </el-menu-item>
                            </el-submenu>
                            <el-menu-item :key="nav_v.key" v-else-if="shouldShowMenu(nav_v)" :index="nav_v.path">
                                <i :class="nav_v.icon"></i>
                                <span slot="title">{{ nav_v.name }}</span>
                            </el-menu-item>
                        </template>
                    </el-menu>
                </div>
                <div class="sidebar-toggle">
                    <div class="icon-left" @click="sidebarToggle">
                        <i class="fas fa-chevron-circle-left"></i>
                    </div>
                    <div id="logout-system" class="log-out" @click="logout">
                        <i class="el-icon-switch-button" style="color: #E6A23C"></i>
                    </div>
                </div>
            </div>
            <div class="hide-aside-toggle-bar" id="hide-all-aside">
                <label class="toggle-button" @click="hiddenSidebar"><i id="toggle-button-icon"></i></label>
                <i class="toggle-slider"/>
            </div>
            <div class="app-body">
                <NavBar id="nav-bar" v-if="switchTabBar"
                        :style="fixedTabBar && switchTabBar?'position: fixed;top: 0;':''"></NavBar>
                <div v-else></div>
                <div id="mainContainer" :style="fixedTabBar && switchTabBar ? 'margin-top: 88px;' : ''"
                     class="main-container">
                    <!--<transition name="fade">-->
                    <router-view></router-view>
                    <!--</transition>-->
                </div>
                <!--EuiFooter><!/EuiFooter-->
            </div>
        </div>
    </div>
</template>

<script>
    /* eslint-disable no-unused-vars */
    import Screenfull from 'screenfull/dist/screenfull'
    import NavBar from './NavBar.vue'
    import Navigation from '@/navigation/index'
    import { logout as logoutApi } from '@/api/user'
    import { doRequest, amIAdmin, syncSleep, mergeNavigation } from '../../utils/utils'
    import { setAddr, getUserInfo } from '../../utils/dataStorage'
    import Config from '../../config/index'
    import Axios from 'axios'

    export default {
        data() {
            return {
                fixedTabBar: false,
                switchTabBar: false,
                siteName: this.$Config.siteName,
                isCollapse: false,
                navigation: {},
                sidebarHidden: false,
                showAbout: false,
                isAdmin: false,
            };
        },
        methods: {
            NavBarWidth() {
                let navBar = document.getElementById('nav-bar');
                if (!navBar) return;
                if (!(this.fixedTabBar && this.switchTabBar)) {
                    navBar.style.width = '100%';
                    return;
                }
                let sidebarClose = document.body.classList.contains('sidebar-close');
                if (sidebarClose) {
                    navBar.style.width = '100%';
                    return;
                }
                if (this.isCollapse) {
                    navBar.style.width = 'calc(100% - 64px)';
                } else {
                    navBar.style.width = 'calc(100% - 270px)';
                }

            },
            ToggleGrayMode() {
                document.body.classList.toggle("gray-mode")
            },
            screenfullToggle() {
                if (!Screenfull.enabled) {
                    this.$message({
                        message: '你的浏览器不支持全屏！',
                        type: 'warning'
                    })
                    return false
                }
                Screenfull.toggle();
            },
            saveFixedTabBar(v) {
                v ? localStorage.setItem('fixedTabBar', v) : localStorage.removeItem('fixedTabBar');
                this.NavBarWidth();
            },
            saveSwitchTabBarVal(v) {
                let containerDom = document.getElementById('mainContainer');
                v ? containerDom.style.minHeight = 'calc(100vh - 139px)' : containerDom.style.minHeight = 'calc(100vh - 101px)';
                v ? localStorage.setItem('switchTabBar', v) : localStorage.removeItem('switchTabBar');
                this.NavBarWidth();
            },
            sidebarToggle(e) {
                e.preventDefault();
                if (this.isCollapse) {
                    document.getElementById('logout-system').style.display = '';
                    document.getElementById('hide-all-aside').style.left = '230px';
                    document.body.classList.remove('sidebar-hidden');
                    this.siteName = this.$Config.siteName;
                    this.isCollapse = false;
                } else {
                    document.getElementById('logout-system').style.display = 'none';
                    document.getElementById('hide-all-aside').style.left = '64px';
                    document.body.classList.add('sidebar-hidden');
                    this.isCollapse = true;
                }
                this.NavBarWidth();
                window.dispatchEvent(new Event('resize'));
            },
            hiddenSidebar(e) {
                e.preventDefault();
                document.body.classList.toggle('sidebar-close');
                this.NavBarWidth();
                this.sidebarHidden = !this.sidebarHidden;
                if (this.sidebarHidden) {
                    document.getElementById('hide-all-aside').style.left = '0px';
                    document.getElementById('toggle-button-icon').style.borderLeft = '6px solid #fff';
                    document.getElementById('toggle-button-icon').style.borderRight = '0px';
                } else {
                    if (this.isCollapse)
                        document.getElementById('hide-all-aside').style.left = '64px';
                    else
                        document.getElementById('hide-all-aside').style.left = '230px';
                    document.getElementById('toggle-button-icon').style.borderRight = '6px solid #fff';
                    document.getElementById('toggle-button-icon').style.borderLeft = '0px';
                }
                window.dispatchEvent(new Event('resize'));
            },
            logout() {
                this.$confirm('您确定要退出系统么？', '退出系统', {
                    confirmButtonText: '',
                    cancelButtonText: '',
                    type: 'warning'
                }).then(() => {
                    logoutApi();
                    this.$router.push({path: '/login'});
                }).catch(() => {

                })
            },
            handleOpen(key, keyPath) {
                //console.log(key, keyPath);
            },
            handleClose(key, keyPath) {
                //关闭菜单
            },
            about() {
                this.showAbout = true;
            },
            shouldShowMenu(nav) {
                if (this.isAdmin)
                    return true
                    
                if (nav.show < 0)
                    return false

                if (nav && nav.hide) {
                    return this.isAdmin 
                }

                return true
            }
        },
        beforeMount() {
            let usr = getUserInfo()
            doRequest({
                method: 'GET',
                url: `/v1/web/user/view-perm/${usr.name}`
            }, {
                success: res => {
                    if (res.viewPerm && res.viewPerm.length > 0) {
                        this.navigation = mergeNavigation(Navigation, JSON.parse(res.viewPerm))
                    } else {
                        this.navigation = Navigation
                    }
                }
            })
        },
        mounted: function () {
            this.isAdmin = amIAdmin()
            // Axios.get(Config.echoUrl).then(function (response) {
            //     // handle success
            //     setAddr(response.data.ip)
            // })
            // .catch(function (error) {
            //     // handle error
            //     console.log(error);
            // })

            // console.log(InternalIp.v4.sync())

            this.switchTabBar = localStorage.getItem('switchTabBar') ? true : false;
            this.fixedTabBar = localStorage.getItem('fixedTabBar') ? true : false;
            if (this.switchTabBar) document.getElementById('mainContainer').style.minHeight = 'calc(100vh - 139px)';


            if (!this.isCollapse) {

                document.body.classList.remove('sidebar-hidden')
                this.siteName = this.$Config.siteName
            } else {
                document.body.classList.add('sidebar-hidden')
            }

            setTimeout(() => {
                this.NavBarWidth();
            }, 1000)
        },
        components: {
            NavBar
        },
    }
</script>
<style lang="scss">
    @import "../../assets/css/variables.scss";

    .menu {
        .el-menu-item {
            svg{
                margin-right: 12px;
                width: 25px !important;
            }
        }

        .el-submenu {
            svg{
                margin-right: 12px;
                width: 25px !important;
            }
        }
    }
    

    .el-dialog {
        margin-top: 5vh !important;
    }

    .sidebar-hidden {
        .header {
            .logo {
                background: #222d32;

                .big {
                    display: none;
                }

                .min {
                    display: block;
                }

                width: 64px;
            }
        }

        .aside {
            .sidebar-toggle {
                .icon-left {
                    transform: rotate(180deg);
                }
            }
        }

        .main {
            .app-body {
                margin-left: 64px;
            }
        }
    }

    .sidebar-close {
        .header {
            .logo {
                width: 0;
                overflow: hidden;
            }
        }

        .aside {
            margin-left: -230px;
        }

        .main {
            .app-body {
                margin-left: 0;
            }
        }

        .hide-aside-toggle-bar {
            left: 0;
        }
    }

    .sidebar-hidden.sidebar-close {
        .aside {
            margin-left: -64px;
        }

        .hide-aside-toggle-bar {
            left: 64px;
        }
    }


    .main {
        display: flex;

        .el-menu:not(.el-menu--collapse) {
            width: 230px;
        }

        .app {
            width: 100%;
            background-color: #ecf0f5;
        }

        .aside {
            position: fixed;
            margin-top: 0px;
            z-index: 10;
            background-color: #222d32;
            transition: all 0.3s ease-in-out;

            .logo {
                text-align: center;
                margin-top: 10px;
                margin-bottom: 50px;
            }

            .menu {
                overflow-y: auto;
                height: calc(100vh - 157px);
            }

            .sidebar-toggle {
                position: relative;
                width: 100%;
                height: 50px;
                color: #fff;
                

                .icon-left {
                    position: absolute;
                    display: flex;
                    align-items: center;
                    justify-content: center;
                    right: 0;
                    width: 64px;
                    height: 100%;
                    font-size: 20px;
                    transition: all 0.3s ease-in-out;
                    cursor: pointer;
                }

                .log-out {
                    position: absolute;
                    display: flex;
                    align-items: center;
                    justify-content: center;
                    left: 0;
                    width: 64px;
                    height: 100%;
                    font-size: 20px;
                    cursor: pointer;
                }
            }
        }

        .hide-aside-toggle-bar {
            z-index: 2;
            display: block;
            position: absolute;
            top: 50%;
            width: 20px;
            height: 150px;
            left: 230px;
            margin-top: -95px;
            transition: all 0.3s ease-in-out;

            .toggle-slider {
                top: 0;
                width: 3px;
                height: 100000em;
                position: absolute;
                left: 0px;
                margin-top: -50000em;
            }

            .toggle-button {
                cursor: pointer;
                display: none;
                height: 155px;
                width: 75px;
                border-bottom-right-radius: 30%;
                border-bottom-left-radius: 30%;
                border-top-left-radius: 30%;
                border-top-right-radius: 30%;
                position: relative;
                margin-left: -60px;
                top: 25%;

                i {
                    border-right: 6px solid #fff;
                    border-top: 7px dashed transparent;
                    border-bottom: 7px dashed transparent;
                    display: inline-block;
                    width: 0;
                    height: 0;
                    vertical-align: center;
                    position: absolute;
                    top: 50%;
                    right: 6px;
                    margin-top: -7px;
                }
            }

            .toggle-button:hover {
                background-color: #ff5d23 !important;
            }

            .toggle-button:hover~.toggle-slider {
                background-color: #ff5d23 !important;
            }

            &:hover .toggle-slider, &:hover .toggle-button {
                display: block;
                background-color: #999;
            }
        }

        .app-body {
            margin-left: 230px;
            -webkit-transition: margin-left 0.3s ease-in-out;
            transition: margin-left 0.3s ease-in-out;
        }

        .main-container {
            //margin-top: 50px;
            padding: 6px;
            min-height: calc(100vh);
        }
    }

    .header {
        width: 100%;
        position: fixed;
        display: flex;
        height: 50px;
        background-color: $--color-primary;
        z-index: 10;

        .logo {
            .min {
                display: block;
            }

            width: 230px;
            height: 50px;
            text-align: center;
            overflow: hidden;
            line-height: 50px;
            color: #fff;
            background-color: rgb(34, 45, 50);
            -webkit-transition: width 0.35s;
            transition: all 0.3s ease-in-out;
        }

        .right {
            position: absolute;
            right: 0;
        }

        .header-btn {
            .el-badge__content {
                top: 14px;
                right: 7px;
                text-align: center;
                font-size: 9px;
                padding: 0 3px;
                background-color: #00a65a;
                color: #fff;
                border: none;
                white-space: nowrap;
                vertical-align: baseline;
                border-radius: .25em;
            }

            overflow: hidden;
            height: 50px;
            display: inline-block;
            text-align: center;
            line-height: 50px;
            cursor: pointer;
            padding: 0 14px;
            color: #fff;

            &:hover {
                background-color: mix(#000, $--color-primary, 10%);
            }
        }

    }

    .menu {
        border-right: none;
    }

    .el-menu--vertical {
        min-width: 190px;
    }

    .setting-category {
        padding: 10px 0;
        border-bottom: 1px solid #eee;
    }
</style>
