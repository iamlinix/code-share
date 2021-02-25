<template>
    <div class="div-home-page">
        <el-tabs v-model="activeTab" stretch>
            <el-tab-pane :label="config.lang.tabActivity" name="activity">
                <Activity />
            </el-tab-pane>
            <el-tab-pane :label="config.lang.tabPost" name="post">
                <Post />
            </el-tab-pane>
            <el-tab-pane :label="config.lang.tabUser" name="user">
                <User :logoutMethod="logoutMethod"/>
            </el-tab-pane>
        </el-tabs>
    </div>
</template>

<script>
import Config from '../utils/config'
import Activity from './Activity'
import Post from './Post'
import User from './User'

export default {
    components: {
        Activity,
        Post,
        User
    },
    props: {
        logoutMethod: {
            type: Function,
            default: null
        }
    },
    data() {
        return {
            activeTab: 'activity',
            config: Config
        }
    }
}
</script>

<style scoped>
.div-home-page {

}
</style>