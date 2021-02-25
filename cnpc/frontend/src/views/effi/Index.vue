<template>
    <div class="cnpn-effi-div">
        <p class="main-title">商品效率参数配置</p>
        <el-divider />
        <el-row :gutter="8" v-loading="mainLoading">
            <el-col v-for="(v, i) in indices" :key="'cpnc-effi-config-card-' + i" :xs="24" :sm="12" :md="8" :lg="6">
                <el-card class="config-card" v-loading="v.loading">
                    <p class="config-title">{{ v.desc }}</p>
                    <el-input-number size="mini" v-model="v.value" />
                    <el-button size="mini" type="primary" @click="updateConfig(v)">更新</el-button>
                </el-card>
            </el-col>
        </el-row>
    </div>
</template>

<style lang="scss">
.cnpn-effi-div {
    width: 100%;
    max-height: calc(100vh - 12px);
    background-color: white;
    padding: 8px 8px 8px 8px;
    overflow: auto;

    .el-divider {
        margin: 10px 0;
    }

    .main-title {
        width: 100%;
        text-align: center;
        font-size: 18px;
        font-weight: bold;
        color: #606266;
    }

    .config-card {
        .el-card__body {
            padding: 8px 8px;
        }

        .el-input-number {
            width: 100%;
            margin-bottom: 8px;
        }

        .el-button {
            width: 100%;
        }

        p.config-title {
            width: 100%;
            text-align: center;
            font-size: 13px;
            color: #909399;
            margin-bottom: 10px;
        }
    }
}
</style>

<script>
import { doRequest, allRequests, message, confirm } from '../../utils/utils'

export default {
    data() {
        return {
            mainLoading: false,
            modules: [ 'EFF001' ],
            indices: [{
                name: 'eff_sar_th',
                module: 'EFF001',
                desc: '动销率阈值'
            }],
            mapping: {}
        }
    },
    mounted() {
        this.indices.forEach(e => {
            if (!this.mapping[e.module])
                this.mapping[e.module] = []

            this.mapping[e.module].push(e)
        })
        this.getConfigs()
    },
    methods: {
        getConfigs() {
            this.mainLoading = true;
            let gets = []
            this.modules.forEach(e => {
                gets.push({
                    url: `/v1/web/eff/metric/config/list/${e}`
                })
            })
            allRequests({
                success: res => {
                    res.forEach(r => {
                        if (r.data.cfgs) {
                            r.data.cfgs.forEach(e => {
                                this.indices.forEach(i => {
                                    if (e.module == i.module && e.name == i.name) {
                                        i.value = e.value;
                                    }
                                })
                            })
                        }
                    })
                },
                fail: err => {
                    message('error', '获取配置列表失败, 请稍后再试')
                },
                finally: _ => {
                    this.mainLoading = false;
                }
            }, gets)
        },
        getConfigForModule(module) {
            this.mapping[module].forEach(e => {
                e.loading = true;
            })

            doRequest({
                url: `/v1/web/eff/metric/config/list/${module}`
            }, {
                success: res => {
                    res.cfgs.forEach(e => {
                        this.indices.forEach(i => {
                            if (e.module == i.module && e.name == i.name) {
                                i.value = e.value;
                            }
                        })
                    })
                },
                fail: err => {
                    message('获取模块配置参数失败, 请稍后再试')
                },
                finally: _ => {
                    this.mapping[module].forEach(e => {
                        e.loading = false;
                    })
                }
            })
        },
        updateConfig(index) {
            this.$set(index, 'loading', true);
            doRequest({
                url: '/v1/web/eff/metric/config/update',
                method: 'POST',
                data: {
                    name: index.name,
                    value: index.value
                }
            }, {
                success: res => {
                    message('success', '参数更新成功')
                    this.getConfigForModule(index.module)
                },
                fail: err => {
                    message('error', '参数更新失败, 请稍后再试')
                }
            })
        }
    }
}
</script>