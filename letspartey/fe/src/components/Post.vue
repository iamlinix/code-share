<template>
    <div class="div-post-page">
        <el-dialog :visible.sync="visible" width="80%" :close-on-click-modal="false">
            <el-input size="mini" :placeholder="config.lang.phNewPostTitle" v-model="title" style="margin-bottom: 8px"/>
            <el-input
                type="textarea"
                :rows="3"
                :placeholder="config.lang.phNewPostContent"
                v-model="content" style="margin-bottom: 8px"/>
            <el-upload ref="postImageList"
                multiple
                :limit="3"
                action="/api/v1/upload"
                list-type="picture-card"
                :headers="{token: token}"
                :data="imageData"
                :auto-upload="false"
                :on-remove="handleRemove"
                :on-success="created">
                <i class="el-icon-plus"></i>
            </el-upload>
            <el-button size="mini" round style="width: 100%; margin-top: 8px" type="warning" icon="el-icon-check" @click="createPost">
                {{ config.lang.OK }}
            </el-button>
        </el-dialog>
        <el-card v-for="post in posts" :key="'post-block-' + post.id" class="post-block">
            <span style="font-size: 18px; color: #909399">{{ post.name }}</span>
            <span style="font-size: 14px; float: right; color: #C0C4CC">created by: {{ post.creator }}</span>
            <p style="font-size: 14px">{{ post.desc }}</p>
            <el-row v-if="post.images"  :gutter="4">
                <el-col v-for="im in post.images" :key="im" :span="8">
                        <el-image
                            style="width: 100px; height: 100px"
                            :src="'/api/' + im"
                            :fit="fit"
                            :preview-src-list="['/api/' + im]"></el-image>
                </el-col>
            </el-row>
        </el-card>
        <el-button style="width: 100%" round size="medium" type="success" icon="el-icon-plus" @click="visible = true">
            {{ config.lang.opAddPost }}
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
        axios.get('/api/v1/posts')
        .then(res => {
            this.posts = res.data.posts
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
            posts: [],
            imageData: {
                lastId: 0,
                type: 'post'
            }
        }
    },
    methods: {
        handleRemove(file, fileList) {
            console.log(file, fileList);
        },
        handlePictureCardPreview(file) {
            this.dialogImageUrl = file.url;
            this.dialogVisible = true;
        },
        createPost() {
            if (!this.title || !this.content) {
                this.$message({
                    type: 'warning',
                    message: Config.lang.errEmptyValues
                })
                return
            }

            axios.post('/api/v1/post', {
                creator: this.userInfo.username,
                name: this.title,
                desc: this.content,
            })
            .then(res => {
                this.imageData.lastId = res.data.id
                this.$refs.postImageList.submit()
                this.title = ''
                this.content = ''
                this.visible = false
                axios.get('/api/v1/posts')
                .then(res => {
                    this.posts = res.data.posts
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
            this.$refs.postImageList.clearFiles();
            axios.get('/api/v1/posts')
            .then(res => {
                this.posts = res.data.posts
            })
            .catch(err => {
                console.log(err)
            })
        }
    }
}
</script>

<style scoped>
.div-post-page .post-block {
    border-radius: 5px;
    text-align: left;
    margin-bottom: 8px;
    padding: 4px 4px
}
</style>