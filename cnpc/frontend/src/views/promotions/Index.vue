<template>
    <div class="cnpc-bj-promotion">
        <el-card>
            <div class="promo-title">
                <p>中石油北京销售公司2020年春节促销</p>
                <div class="item-group">
                    <span class="group-title">参与商品</span>
                    <el-tag effect="dark" type="primary" v-for="(v, i) in promotedItems" 
                        :key="'promo-item-tag-' + i" size="mini">{{ v.name }}</el-tag>
                </div>
            </div>
            <el-tabs type="border-card" stretch @tab-click="smrschtDoSmartDelayedResize()">
                <el-tab-pane>
                    <span slot="label"><i class="fas fa-angle-double-up"></i> 销售额排名 Top10</span>
                    <div class="pane-div">
                        <el-row :gutter="8">
                            <el-col :span="24">
                                <el-card style="height: 50vh" 
                                    :body-style="{ width: '100%', height: '100%', position: 'relative' }">
                                    <div class="chart-div" id="promo-chart-top1">
                                    </div>
                                </el-card>
                            </el-col>
                            <el-col :span="12">
                                <el-card style="height: 40vh" 
                                    :body-style="{ width: '100%', height: '100%', position: 'relative' }">
                                    <div class="chart-div" id="promo-chart-top2">
                                    </div>
                                </el-card>
                            </el-col>
                            <el-col :span="12">
                                <el-card style="height: 40vh" 
                                    :body-style="{ width: '100%', height: '100%', position: 'relative' }">
                                    <div class="chart-div" id="promo-chart-top3">
                                    </div>
                                </el-card>
                            </el-col>
                            <el-col :span="8">
                                <el-card style="height: 30vh" 
                                    :body-style="{ width: '100%', height: '100%', position: 'relative' }">
                                    <div class="chart-div" id="promo-chart-top4">
                                    </div>
                                </el-card>
                            </el-col>
                            <el-col :span="8">
                                <el-card style="height: 30vh" 
                                    :body-style="{ width: '100%', height: '100%', position: 'relative' }">
                                    <div class="chart-div" id="promo-chart-top5">
                                    </div>
                                </el-card>
                            </el-col>
                            <el-col :span="8">
                                <el-card style="height: 30vh" 
                                    :body-style="{ width: '100%', height: '100%', position: 'relative' }">
                                    <div class="chart-div" id="promo-chart-top6">
                                    </div>
                                </el-card>
                            </el-col>
                            <el-col :span="6">
                                <el-card style="height: 30vh" :body-style="{ width: '100%', height: '100%', position: 'relative' }">
                                    <div class="chart-div" id="promo-chart-top7">
                                    </div>
                                </el-card>
                            </el-col>
                            <el-col :span="6">
                                <el-card style="height: 30vh" :body-style="{ width: '100%', height: '100%', position: 'relative' }">
                                    <div class="chart-div" id="promo-chart-top8">
                                    </div>
                                </el-card>
                            </el-col>
                            <el-col :span="6">
                                <el-card style="height: 30vh" :body-style="{ width: '100%', height: '100%', position: 'relative' }">
                                    <div class="chart-div" id="promo-chart-top9">
                                    </div>
                                </el-card>
                            </el-col>
                            <el-col :span="6">
                                <el-card style="height: 30vh" :body-style="{ width: '100%', height: '100%', position: 'relative' }">
                                    <div class="chart-div" id="promo-chart-top10">
                                    </div>
                                </el-card>
                            </el-col>
                        </el-row>
                    </div>
                </el-tab-pane>
                <el-tab-pane>
                    <span slot="label"><i class="fas fa-angle-double-down"></i> 销售额排名倒数 Top10</span>
                    <div class="pane-div">
                        <el-row :gutter="8">
                            <el-col :span="24">
                                <el-card style="background-color: #333333; height: 50vh" :body-style="{ width: '100%', height: '100%', position: 'relative' }">
                                    <div class="chart-div" id="promo-chart-rop1">
                                    </div>
                                </el-card>
                            </el-col>
                            <el-col :span="12">
                                <el-card style="background-color: #333333; height: 40vh" :body-style="{ width: '100%', height: '100%', position: 'relative' }">
                                    <div class="chart-div" id="promo-chart-rop2">
                                    </div>
                                </el-card>
                            </el-col>
                            <el-col :span="12">
                                <el-card style="background-color: #333333; height: 40vh" :body-style="{ width: '100%', height: '100%', position: 'relative' }">
                                    <div class="chart-div" id="promo-chart-rop3">
                                    </div>
                                </el-card>
                            </el-col>
                            <el-col :span="8">
                                <el-card style="background-color: #333333; height: 30vh" :body-style="{ width: '100%', height: '100%', position: 'relative' }">
                                    <div class="chart-div" id="promo-chart-rop4">
                                    </div>
                                </el-card>
                            </el-col>
                            <el-col :span="8">
                                <el-card style="background-color: #333333; height: 30vh" :body-style="{ width: '100%', height: '100%', position: 'relative' }">
                                    <div class="chart-div" id="promo-chart-rop5">
                                    </div>
                                </el-card>
                            </el-col>
                            <el-col :span="8">
                                <el-card style="background-color: #333333; height: 30vh" :body-style="{ width: '100%', height: '100%', position: 'relative' }">
                                    <div class="chart-div" id="promo-chart-rop6">
                                    </div>
                                </el-card>
                            </el-col>
                            <el-col :span="6">
                                <el-card style="background-color: #333333; height: 30vh" :body-style="{ width: '100%', height: '100%', position: 'relative' }">
                                    <div class="chart-div" id="promo-chart-rop7">
                                    </div>
                                </el-card>
                            </el-col>
                            <el-col :span="6">
                                <el-card style="background-color: #333333; height: 30vh" :body-style="{ width: '100%', height: '100%', position: 'relative' }">
                                    <div class="chart-div" id="promo-chart-rop8">
                                    </div>
                                </el-card>
                            </el-col>
                            <el-col :span="6">
                                <el-card style="background-color: #333333; height: 30vh" :body-style="{ width: '100%', height: '100%', position: 'relative' }">
                                    <div class="chart-div" id="promo-chart-rop9">
                                    </div>
                                </el-card>
                            </el-col>
                            <el-col :span="6">
                                <el-card style="background-color: #333333; height: 30vh" :body-style="{ width: '100%', height: '100%', position: 'relative' }">
                                    <div class="chart-div" id="promo-chart-rop10">
                                    </div>
                                </el-card>
                            </el-col>
                        </el-row>
                    </div>
                </el-tab-pane>
                <el-tab-pane label="详细信息">
                    <div class="pane-div">
                        <el-table border :data="details" size="mini">
                            <el-table-column type="index"></el-table-column>
                            <el-table-column label="加油站" prop="st"></el-table-column>
                            <el-table-column sortable label="N95口罩销售额" prop="n95"></el-table-column>
                            <el-table-column sortable label="普通口罩销售额" prop="reg"></el-table-column>
                            <el-table-column sortable label="蔬菜包销售额" prop="veg"></el-table-column>
                            <el-table-column sortable label="消毒液销售额" prop="alc"></el-table-column>
                            <el-table-column sortable label="总促销销售额" prop="total"></el-table-column>
                        </el-table>
                    </div>
                </el-tab-pane>
            </el-tabs>
        </el-card>
    </div>
</template>

<style scoped lang="scss">
.cnpc-bj-promotion {
    .chart-div {
        width: 100%;
        height: 100%;
    }

    .pane-div {
        height: calc(100vh - 270px);
        overflow: auto;
        padding: 0 8px 8px 8px;
    }

    .promo-title {
        width: 100%;
        border-radius: 8px;
        background-color: #E9E9EB;
        padding: 10px 10px 10px 10px;
        position: relative;
        margin-bottom: 10px;

        p {
            text-align: center;
            color: #909399;
            font-size: 15px;
            font-weight: bold;
        }

        .item-group {
            span.group-title {
                position: absolute;
                left: 20px;
                top: 28px;
                padding: 4px 4px 4px 4px;
                background-color: #E9E9EB;
                color: #909399;
                font-size: 10px;
            }
            border-radius: 5px;
            border: solid 1px #909399;
            margin-top: 10px;
            padding: 16px 10px 16px 10px;

            .el-tag {
                margin-right: 4px;
                margin-top: 4px;
            }
        }
    }

    .el-card {
        margin-top: 10px;
    }
}
</style>

<script>
import SmartResizeChart from '../mixins/SmartResizeChart'
import echarts from 'echarts';
import 'echarts/theme/macarons'
import 'echarts/theme/dark'
import { randomNumber } from '../../utils/utils'

export default {
    mixins: [ SmartResizeChart ],
    data() {
        return {
            promotedItems: [{
                name: 'N95口罩'
            },{
                name: '普通一次性口罩'
            },{
                name: '蔬菜包'
            },{
                name: '消毒液'
            }],
            top10Charts: [],
            rop10Charts: [],

            details: [{
                st: '北京销售第二分公司方深加油站',
                n95: 123,
                reg: 342,
                veg: 1223,
                alc: 90,
                total: 2332
            }, {
                st: '北京销售第四分公司钰鑫增加油站',
                n95: 34,
                reg: 555,
                veg: 231,
                alc: 433,
                total: 3456
            }, {
                st: '北京销售第三分公司京顺路加油站',
                n95: 432,
                reg: 22,
                veg: 333,
                alc: 555,
                total: 4345
            }, {
                st: '北京销售第三分公司庆铁加油站',
                n95: 334,
                reg: 87,
                veg: 666,
                alc: 543,
                total: 2332
            }, {
                st: '北京销售第三分公司顺兴祥加油站',
                n95: 123,
                reg: 342,
                veg: 1223,
                alc: 90,
                total: 2332
            }, {
                st: '北京销售第一分公司南下路加油站',
                n95: 123,
                reg: 342,
                veg: 1223,
                alc: 90,
                total: 2332
            }, {
                st: '北京销售第三分公司利易泰加油站',
                n95: 123,
                reg: 342,
                veg: 1223,
                alc: 90,
                total: 2332
            }, {
                st: '北京销售第三分公司快鑫加油站',
                n95: 123,
                reg: 342,
                veg: 1223,
                alc: 90,
                total: 2332
            }, {
                st: '北京销售第三分公司东方动力加油站',
                n95: 123,
                reg: 342,
                veg: 1223,
                alc: 90,
                total: 2332
            }, {
                st: '北京销售第二分公司通惠加油站',
                n95: 123,
                reg: 342,
                veg: 1223,
                alc: 90,
                total: 2332
            }, {
                st: '北京销售第一分公司南营加油站',
                n95: 123,
                reg: 342,
                veg: 1223,
                alc: 90,
                total: 2332
            }]
        }
    },
    mounted() {
        this.details.forEach(e => {
            e.n95 = randomNumber(50, 2000);
            e.reg = randomNumber(50, 2000);
            e.veg = randomNumber(50, 2000);
            e.alc = randomNumber(50, 2000);
            e.total = e.n95 + e.reg + e.veg + e.alc;
        })
        let titles = [{
            t: '北京销售第二分公司方深加油站',
            r: '北京销售第四分公司钰鑫增加油站'
        }, {
            t: '北京销售第三分公司京顺路加油站',
            r: '北京销售第三分公司庆铁加油站'
        }, {
            t: '北京销售第三分公司顺兴祥加油站',
            r: '北京销售第一分公司南下路加油站'
        }, {
            t: '北京销售第三分公司利易泰加油站',
            r: '北京销售第三分公司快鑫加油站'
        }, {
            t: '北京销售第三分公司东方动力加油站',
            r: '北京销售第二分公司通惠加油站'
        }, {
            t: '北京销售第一分公司南营加油站',
            r: '北京销售第一分公司新景都市加油站'
        }, {
            t: '北京销售第四分公司时顺苑加油站',
            r: '北京销售第一分公司南平庄加油站'
        }, {
            t: '北京销售第二分公司北苑加油站',
            r: '北京销售第二分公司前苇沟加油站'
        }, {
            t: '北京销售第二分公司康凯旋加油站',
            r: '北京销售第二分公司奥东加油站'
        }, {
            t: '北京销售第二分公司天诚金龙加油站',
            r: '北京销售第二分公司常营加油站'
        }]
        let option = {
            tooltip: {
                trigger: 'item',
                formatter: '{a} <br/>{b}: {c} ({d}%)'
            },
            legend: {
                orient: 'vertical',
                left: 10,
                data: ['N95口罩', '普通一次性口罩', '蔬菜包', '消毒液']
            },
            title: {
                left: 'center',
                textStyle: {
                    fontSize: 13
                }
            },
            series: [
                {
                    name: '促销商品',
                    type: 'pie',
                    radius: ['50%', '70%'],
                    avoidLabelOverlap: false,
                    label: {
                        normal: {
                            show: false,
                            position: 'center'
                        },
                        emphasis: {
                            show: true,
                            textStyle: {
                                fontSize: '20',
                                fontWeight: 'bold'
                            }
                        }
                    },
                    labelLine: {
                        normal: {
                            show: false
                        }
                    },
                    data: [
                        {value: 100, name: 'N95口罩'},
                        {value: 310, name: '普通一次性口罩'},
                        {value: 234, name: '蔬菜包'},
                        {value: 135, name: '消毒液'}
                    ]
                }
            ]
        };
        setTimeout(() => {
            let t = null, r = null;
            for (let i = 1; i <= 10; i ++) {
                t = echarts.init(document.getElementById(`promo-chart-top${i}`), 'macarons');
                r = echarts.init(document.getElementById(`promo-chart-rop${i}`), 'dark')
                this.top10Charts.push(t)
                this.rop10Charts.push(r)

                if (i === 4)
                    delete option['legend'];

                option.title.text = titles[i - 1].t
                t.setOption(option, true);
                option.title.text = titles[i - 1].r
                r.setOption(option, true);

                this.smrschtAddChart(t);
                this.smrschtAddChart(r);

                this.smrschtMonitorResizeById(`promo-chart-top${i}`);
                this.smrschtMonitorResizeById(`promo-chart-rop${i}`);
            }
            
        }, 500)
        
    },
    methods: {

    }
}
</script>