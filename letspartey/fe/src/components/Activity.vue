<template>
    <div class="div-activity-page">
        <el-dialog :visible.sync="visible" width="80%" :close-on-click-modal="false">
            <el-input size="mini" :placeholder="config.lang.phNewActTitle" v-model="title" style="margin-bottom: 8px"/>
            <el-input
                type="textarea"
                :rows="3"
                :placeholder="config.lang.phNewActContent"
                v-model="content" style="margin-bottom: 8px"/>
            <el-date-picker
                v-model="startTime"
                type="datetime"
                size="mini"
                :placeholder="config.lang.phNewActDate"
                style="width: 100%; margin-bottom: 8px">
                </el-date-picker>
            <el-upload ref="actImageList"
                multiple
                :limit="3"
                action="/api/v1/upload"
                list-type="picture-card"
                :headers="{token: token}"
                :data="imageData"
                :auto-upload="false"
                :on-remove="handleRemove"
                :on-success="created"
                :on-exceed="handleExceed">
                <i class="el-icon-plus"></i>
            </el-upload>
            <el-button size="mini" round style="width: 100%; margin-top: 8px" type="warning" icon="el-icon-check" @click="createActivity">
                {{ config.lang.OK }}
            </el-button>
        </el-dialog>
        <el-card v-for="act in activities" :key="'act-block-' + act.id" class="activity-block">
            <span style="font-size: 18px; color: #909399">{{ act.name }}</span>
            <span style="font-size: 14px; float: right; color: #C0C4CC">created by: {{ act.creator }}</span>
            
            <p style="font-size: 13px; color: #C0C4CC">Start Time: {{ act.startTime }}</p>
            <p style="font-size: 14px">{{ act.desc }}</p>
            <el-row v-if="act.images"  :gutter="4">
                <el-col v-for="im in act.images" :key="im" :span="8">
                        <el-image
                            style="width: 100px; height: 100px"
                            :src="'/api/' + im"
                            :fit="fit"
                            :preview-src-list="['/api/' + im]"></el-image>
                </el-col>
            </el-row>
        </el-card>
        <el-button style="width: 100%;" round size="medium" type="primary" icon="el-icon-plus" @click="visible = true">
            {{ config.lang.opAddActivity }}
        </el-button>
    </div>
</template>

<script>
import Config from '../utils/config'
import { getToken, getUserInfo } from '../utils/utils'
import axios from 'axios'

export default {
    beforeMount() {
        this.token = getToken()
        this.userInfo = getUserInfo()
    },
    mounted() {
        axios.get('/api/v1/activities')
        .then(res => {
            this.activities = res.data.activities
        })
        .catch(err => {
            console.log(err)
        })
    },
    data() {
        return {
            userInfo: null,
            config: Config,
            visible: false,
            title: '',
            content: '',
            token: '',
            startTime: '',
            activities: [],
            imageData: {
                lastId: 0,
                type: 'activity'
            }
        }
    },
    methods: {
        handleExceed() {
            this.$message.warning(Config.lang.errThreeImageLimit);
        },
        handleRemove(file, fileList) {
            console.log(file, fileList);
        },
        handlePictureCardPreview(file) {
            this.dialogImageUrl = file.url;
            this.dialogVisible = true;
        },
        createActivity() {
            if (!this.title || !this.content || !this.startTime) {
                this.$message({
                    type: 'warning',
                    message: Config.lang.errEmptyValues
                })
                return
            }
            axios.post('/api/v1/activity', {
                creator: this.userInfo.username,
                name: this.title,
                desc: this.content,
                startTime: this.startTime,
            })
            .then(res => {
                this.imageData.lastId = res.data.id
                this.$refs.actImageList.submit()
                this.title = ''
                this.content = ''
                this.startTime = ''
                this.visible = false
                axios.get('/api/v1/activities')
                .then(res => {
                    this.activities = res.data.activities
                })
                .catch(err => {
                    console.log(err)
                })
            })
            .catch(err => {
                console.log(err)
            })
        },
        created() {
            this.$refs.actImageList.clearFiles();
            axios.get('/api/v1/activities')
            .then(res => {
                this.activities = res.data.activities
            })
            .catch(err => {
                console.log(err)
            })
        }
    }
}
</script>

<style scoped>
.div-activity-page .activity-block {
    border-radius: 5px;
    text-align: left;
    margin-bottom: 8px;
    padding: 4px 4px
}
</style>