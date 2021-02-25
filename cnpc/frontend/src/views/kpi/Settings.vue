<template>
    <div class="kpi-settings-container">
        <el-row :gutter="8">
            <el-col :span="8">
                <el-card>
                    <p slot="header" style="text-align: center">
                        设定销售目标
                    </p>
                    <el-form label-width="70px" size="small">
                        <el-form-item label="组织机构">
                            <el-select v-model="orgCode" filterable style="width: 100%">
                                <el-option v-for="o in orgs" :key="o.orgCode" :label="o.orgCode + '-' + o.orgText" :value="o.orgCode" />
                            </el-select>
                        </el-form-item>
                        <el-form-item label="指标时限">
                            <el-select v-model="month" style="width: 100%">
                                <el-option label="全年" :value="0"/>
                                <el-option label="一季度" :value="13"/>
                                <el-option label="二季度" :value="14"/>
                                <el-option label="三季度" :value="15"/>
                                <el-option label="四季度" :value="16"/>
                                <el-option label="一月" :value="1"/>
                                <el-option label="二月" :value="2"/>
                                <el-option label="三月" :value="3"/>
                                <el-option label="四月" :value="4"/>
                                <el-option label="五月" :value="5"/>
                                <el-option label="六月" :value="6"/>
                                <el-option label="七月" :value="7"/>
                                <el-option label="八月" :value="8"/>
                                <el-option label="九月" :value="9"/>
                                <el-option label="十月" :value="10"/>
                                <el-option label="十一月" :value="11"/>
                                <el-option label="十二月" :value="12"/>
                            </el-select>
                        </el-form-item>
                        <el-form-item label="金额类型">
                            <el-select v-model="incomeType" filterable style="width: 100%">
                                <el-option :value="0" label="收入"/>
                                <el-option :value="1" label="毛利"/>
                            </el-select>
                        </el-form-item>
                        <!--el-form-item label="商品大类">
                            <el-select v-model="category" filterable style="width: 100%">
                                <el-option value="0000" label="0000-全品类" />
                                <el-option v-for="o in mainCategories" 
                                    :key="o.classCode" :label="o.classText" :value="o.classCode" />
                            </el-select>
                        </el-form-item-->
                        <el-form-item label="销售目标">
                            <el-input-number :min="0" v-model="value" style="width: 100%"/>
                        </el-form-item>
                    </el-form>
                    <el-button type="success" round style="width: 100%" icon="el-icon-plus" @click="updateGoal">设定目标</el-button>
                </el-card>
                <!--el-card style="margin-top: 8px">
                    <p slot="header" style="text-align: center">
                        财年日期设置
                    </p>
                    <el-form label-width="60px" size="mini">
                        <el-form-item :label="monthDisp[i + 1]" v-for="(v, i) in fymonths" :key="'fymonth-' + i">
                            <el-row :gutter="4">
                                <el-col :span="20">
                                    <el-date-picker 
                                        size="mini" style="width: 100%;"
                                        v-model="fymonths[i]"
                                        type="daterange"
                                        range-separator="-"
                                        start-placeholder="开始日期"
                                        end-placeholder="结束日期"
                                        format="MM-dd"
                                        value-format="MM-dd"
                                        @change="onFyMonthChange(i)">
                                    </el-date-picker>
                                </el-col>
                                <el-col :span="3">
                                    <el-button v-if="fyUpdate[i]" circle type="success" size="mini" icon="el-icon-upload2"
                                        @click="updateFyMonth(i)"></el-button>
                                </el-col>
                            </el-row>
                        </el-form-item>
                    </el-form>
                </el-card-->
            </el-col>
            <el-col :span="16">
                <el-card>
                    <smart-table :minus="60" :data="filteredKpis" @cell-click="cellClicked" :highlight-current-row="true">
                        <vxe-table-column width="100px" title="站点ID" field="orgCode" sortable/>
                        <vxe-table-column title="站点名称" field="orgName" show-overflow="tooltip">
                            <template slot="header" slot-scope="r">
                                <el-input v-model="idSearch" placeholder="输入编码或者名称过滤"
                                    size="mini" @change="doFilter()"/>
                            </template>
                        </vxe-table-column>
                        <vxe-table-column width="110px" title="月份" field="month">
                            <template slot="header" slot-scope="r">
                                <el-select v-model="monthSearch" style="width: 100%" multiple 
                                    clearable size="mini"  @change="doFilter()">
                                    <el-option label="全年" :value="0"/>
                                    <el-option label="一季度" :value="13"/>
                                    <el-option label="二季度" :value="14"/>
                                    <el-option label="三季度" :value="15"/>
                                    <el-option label="四季度" :value="16"/>
                                    <el-option label="一月" :value="1"/>
                                    <el-option label="二月" :value="2"/>
                                    <el-option label="三月" :value="3"/>
                                    <el-option label="四月" :value="4"/>
                                    <el-option label="五月" :value="5"/>
                                    <el-option label="六月" :value="6"/>
                                    <el-option label="七月" :value="7"/>
                                    <el-option label="八月" :value="8"/>
                                    <el-option label="九月" :value="9"/>
                                    <el-option label="十月" :value="10"/>
                                    <el-option label="十一月" :value="11"/>
                                    <el-option label="十二月" :value="12"/>
                                </el-select>
                            </template>
                            <template slot-scope="r">
                                {{ monthDisp[r.row.month] }}
                            </template>
                        </vxe-table-column>
                        <vxe-table-column width="120px" title="金额类型" field="incomeType">
                            <template slot="header" slot-scope="r">
                                <el-select v-model="incomeSearch" size="mini" @change="doFilter()">
                                    <el-option :value="-1" label="全部"/>
                                    <el-option :value="0" label="收入"/>
                                    <el-option :value="1" label="毛利"/>
                                </el-select>
                            </template>
                            <template slot-scope="r">
                                {{ r.row.incomeType == 0 ? '收入' : '毛利' }}
                            </template>
                        </vxe-table-column>
                        <vxe-table-column width="140px" title="销售目标" field="value" sortable>
                            <template slot-scope="r">
                                {{ beautyNum(r.row.value) }}
                            </template>
                        </vxe-table-column>
                        <vxe-table-column title="操作" width="50px">
                            <template slot-scope="r">
                                <i class="el-icon-error" style="color: #F56C6C; font-size: 15px; cursor: pointer" @click="deleteGoal(r.row)"/>
                            </template>
                        </vxe-table-column>
                    </smart-table>
                </el-card>
            </el-col>
        </el-row>
    </div>
</template>

<script>
import { doRequest, message, confirm, beautifyNumber, deepCopy } from '../../utils/utils'
import SmartTable from '../mixins/SmartVxeTable'

export default {
    components: {
        SmartTable
    },
    data() {
        return {
            idSearch: '',
            nameSearch: '',
            monthSearch: [],
            incomeSearch: -1,
            filteredKpis: [],
            fymonths: [
                [],
                [],
                [],
                [],
                [],
                [],
                [],
                [],
                [],
                [],
                [],
                [],
            ],
            fyUpdate: [
                false,
                false,
                false,
                false,
                false,
                false,
                false,
                false,
                false,
                false,
                false,
                false,
            ],
            kpis: [],
            hardcodeOrgs: [{
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
            }],
            orgs: [],
            orgCode: '',
            orgName: '',
            month: 0,
            incomeType: 0,
            category: '0000',
            value: 0,
            mainCategories: [],
            monthDisp: {
                0: '全年',
                1: '一月',
                2: '二月',
                3: '三月',
                4: '四月',
                5: '五月',
                6: '六月',
                7: '七月',
                8: '八月',
                9: '九月',
                10: '十月',
                11: '十一月',
                12: '十二月',
                13: '一季度',
                14: '二季度',
                15: '三季度',
                16: '四季度',
            }
        }
    },
    mounted() {
        doRequest({
            method: 'GET',
            url: '/v1/web/kpi/list',
            loading: true,
        }, {
            success: res => {
                this.kpis = res.kpis
                this.filteredKpis = deepCopy(this.kpis)
            },
            fail: _ => {
                message('error', '获取销售目标列表失败，请稍后再试')
            }
        })

        doRequest({
            method: 'GET',
            url: '/v1/web/org/plant'
        }, {
            success: res => {
                this.orgs = this.hardcodeOrgs.concat(res.orgList)
            },
            fail: _ => {
                message('error', '获取组织机构列表失败，请稍后再试')
            }
        })

        doRequest({
            method: 'GET',
            url: '/v1/web/kpi/fystart'
        }, {
            success: res => {
                if (res.value.length > 0)
                    this.fystart = res.value
                else {
                    let d = new Date()
                    this.fystart = this.fystart2 = `${d.getFullYear()}-01-01`
                }
            },
            fail: _ => {
                message('error', '获取财年设置失败，请稍后再试')
            }
        })

        doRequest({
            method: 'GET',
            url: '/v1/web/class/main'
        }, {
            success: res => {
                this.mainCategories = res.classList
            },
            fail: _ => {
                message('error', '获取商品大类信息失败，请稍后再试')
            }
        })

        this.getFyMonths()
    },
    methods: {
        doFilter() {
            console.log("do filter")
            let kpis = []
            this.kpis.forEach(e => {
                if (this.idSearch.length > 0) {
                    if (!e.orgCode.includes(this.idSearch.toUpperCase()) && !e.orgName.toUpperCase().includes(this.idSearch)) {
                        return
                    }
                }

                if (this.monthSearch.length > 0) {
                    if (!this.monthSearch.includes(e.month)) {
                        return
                    }
                }

                if (this.incomeSearch >= 0) {
                    if (e.incomeType != this.incomeSearch) {
                        return 
                    }
                }

                kpis.push(e)
            })

            this.filteredKpis = kpis
        },
        updateGoal() {
            let orgName = ''
            this.orgs.forEach(e => {
                if (e.orgCode == this.orgCode) {
                    orgName = e.orgText
                    return;
                }
            })
            doRequest({
                method: 'POST',
                url: '/v1/web/kpi/update',
                loading: true,
                data: {
                    orgCode: this.orgCode,
                    orgName: orgName,
                    month: this.month,
                    incomeType: this.incomeType,
                    value: this.value
                }
            }, {
                success: _ => {
                    doRequest({
                        method: 'GET',
                        url: '/v1/web/kpi/list',
                        loading: true,
                    }, {
                        success: res => {
                            this.kpis = res.kpis
                        },
                        fail: _ => {
                            message('error', '获取销售目标列表失败，请稍后再试')
                        }
                    })
                },
                fail: _ => {
                    message('error', '设置销售目标失败，请稍后再试')
                }
            })
        },
        deleteGoal(row) {
            confirm('warning', '删除目标', `确定要删除 ${row.orgName} 的 ${this.monthDisp[row.month]} 目标么?`, _ => {
                doRequest({
                    method: 'POST',
                    url: '/v1/web/kpi/delete',
                    loading: true,
                    data: {
                        orgCode: row.orgCode,
                        month: row.month,
                        incomeType: row.incomeType,
                    }
                }, {
                    success: _ => {
                        doRequest({
                            method: 'GET',
                            url: '/v1/web/kpi/list',
                            loading: true,
                        }, {
                            success: res => {
                                this.kpis = res.kpis
                            },
                            fail: _ => {
                                message('error', '获取销售目标列表失败，请稍后再试')
                            }
                        })
                    },
                    fail: _ => {
                        message('error', '删除销售目标失败，请稍后再试')
                    }
                })
            })
        },
        updateFyStart() {
            doRequest({
                method: 'POST',
                url: '/v1/web/kpi/fystart',
                data: {
                    value: this.fystart2
                }
            }, {
                success: _ => {
                    message('success', '更新财年设置成功')
                    this.fystart = this.fystart2
                },
                fail: _ => {
                    message('error', '更新财年设置失败，请稍后再试')
                }
            })
        },
        updateFyMonth(i) {
            let mon = this.fymonths[i]
            if (!mon || mon.length <= 0) {
                message('error', '月份设置不能为空')
                return
            }
            doRequest({
                method: 'POST',
                url: '/v1/web/kpi/fmonth',
                data: {
                    name: `fymon${(i + '').padStart(2, '0')}`,
                    value: `${mon[0]}~${mon[1]}`
                },
                loading: true
            }, {
                success: _ => {
                    message('success', '设置财年日期成功')
                    this.$set(this.fyUpdate, i, false)
                },
                fail: _ => {
                    message('error', '设置财年日期失败，请稍后再试')
                }
            })
        },
        getFyMonths() {
            doRequest({
                method: 'GET',
                url: '/v1/web/kpi/fmonths',
                loading: true
            }, {
                success: res => {
                    if (res.values) {
                        res.values.forEach(e => {
                            let p = e.value.split('~')
                            this.$set(this.fymonths, parseInt(e.name.substring(5)), [p[0], p[1]])
                        })
                    }
                },
                fail: _ => {
                    message('error', '获取财年日期设置失败，请稍后再试')
                }
            })
        },
        onFyMonthChange(i) {
            if (!this.fymonths[i] || this.fymonths[i].length <= 0) {
                message('error', '月份设置不可为空')
                return
            }

            for (let k = 1; k < this.fymonths.length; k ++) {
                let prev = this.fymonths[k - 1], now = this.fymonths[k]
                if (prev.length > 0 && prev[1] >= now[0]) {
                    message('warning', '时间范围重叠，请注意')
                }
            }
            this.$set(this.fyUpdate, i, true)
        },
        beautyNum(num) {
            return beautifyNumber(num)
        },
        cellClicked(obj, ev) {
            this.orgCode = obj.row.orgCode
            this.orgName = obj.row.orgText
            this.month = obj.row.month
            this.incomeType = obj.row.incomeType
            this.value = obj.row.value
        }
    }
}
</script>

<style scoped>
.kpi-settings-container {

}
</style>