<template>
    <div class="price-ana">
        <el-dialog :visible.sync="showCreateScenarioDialog" width="40%" @closed="scenarioPopoverCleanup" @opened="scenarioPopoverShow" >
            <p slot="title" style="width: 100%; text-align: center">{{ updatingScenario ? '更新查询场景' : '新建查询场景' }}</p>
            <p v-show="showCreateScenarioAlert" class="create-scenario-alert">
                在一个查询场景中选择过多的组织机构或者商品各类，将会使得页面上图表数量过多，可能会使得页面刷新变得缓慢，建议您将图表均匀分布至多个查询场景中。
            </p>
            <el-form label-position="right" label-width="100px" size="small">
                <el-form-item label="场景名称">
                    <el-input v-model="sceneName" />
                </el-form-item>
                <el-form-item style="background-color: #E9E9EB">
                    <span slot="label" style="color: gray; padding-bottom: 12px">时间范围</span>
                </el-form-item>
                <el-form-item label="实时时间">
                    <el-tooltip class="item" effect="light" content="实时时间是指，设定一个时间长度，每次查看页面数据时，均使用当前的实时时间点，向前推进该时间长度，作为数据统计的时间范围" placement="right">
                        <el-switch v-model="useRealTime" active-color="#13ce66" inactive-color="#ff4949" />
                    </el-tooltip>
                </el-form-item>
                <el-form-item label="起始时间">
                    <div v-if="useRealTime">
                        <el-input-number v-model="realTimeCount" :min="1" style="margin-right: 10px"/>
                        <el-select v-model="realTimeMode" style="width: 80px">
                            <el-option label="天" value="d" />
                            <el-option label="月" value="m" />
                        </el-select>
                    </div>
                    <el-date-picker v-else
                        v-model="dateRange"
                        type="daterange"
                        align="right"
                        unlink-panels
                        range-separator="-"
                        start-placeholder="开始日期"
                        end-placeholder="结束日期"
                        :picker-options="pickerOptions"
                        class="stretch-select">
                    </el-date-picker>
                </el-form-item>
                <el-form-item style="background-color: #E9E9EB">
                    <span slot="label" style="color: gray; padding-bottom: 12px">组织机构</span>
                </el-form-item>
                <el-form-item label="北京公司">
                    <el-select filterable multiple class="stretch-select" v-model="orgMajor" placeholder="请选择一级组织机构" @change="createScenarioChangeAlert">
                        <el-option v-for="item in orgListMajor" :key="item.orgCode" :label="item.orgText" :value="item.orgCode"/>
                    </el-select>
                </el-form-item>
                <el-form-item label="分公司">
                    <el-select filterable multiple class="stretch-select" v-model="orgMinor" placeholder="请选择二级组织机构" @change="createScenarioChangeAlert">
                        <el-option v-for="item in orgListMinor" :key="item.orgCode" :label="item.orgText" :value="item.orgCode"/>
                    </el-select>
                </el-form-item>
                <el-form-item label="加油站">
                    <el-select filterable multiple class="stretch-select" v-model="orgMicro" placeholder="请选择三级组织机构" @change="createScenarioChangeAlert">
                        <el-option v-for="item in orgListMicro" :key="item.orgCode" :label="item.orgText" :value="item.orgCode"/>
                    </el-select>
                </el-form-item>
                <el-form-item style="background-color: #E9E9EB">
                    <span slot="label" style="color: gray; padding-bottom: 12px">商品分类</span>
                </el-form-item>
                <el-form-item label="大类">
                    <el-select filterable multiple class="stretch-select" v-model="prodMajor" placeholder="请选择商品大类" @change="createScenarioChangeAlert">
                        <el-option v-for="item in prodListMajor" :key="item.classCode" :label="item.classText" :value="item.classCode"/>
                    </el-select>
                </el-form-item>
                <el-form-item label="中类">
                    <el-select filterable multiple class="stretch-select" v-model="prodMinor" placeholder="请选择商品中类" @change="createScenarioChangeAlert">
                        <el-option v-for="item in prodListMinor" :key="item.classCode" :label="item.classText" :value="item.classCode"/>
                    </el-select>
                </el-form-item>
                <el-form-item label="小类">
                    <el-select filterable multiple class="stretch-select" v-model="prodMicro" placeholder="请选择商品小类" @change="createScenarioChangeAlert">
                        <el-option v-for="item in prodListMicro" :key="item.classCode" :label="item.classText" :value="item.classCode"/>
                    </el-select>
                </el-form-item>
            </el-form>
            <el-row :gutter="10">
                <el-col :span="12">
                    <el-button v-if="!updatingScenario" icon="el-icon-plus" type="primary" style="width: 100%" @click="createScenario(false)">创建场景</el-button>
                    <el-button v-else icon="el-icon-upload" type="success" style="width: 100%" @click="updateScenario">更新场景</el-button>
                </el-col>
                <el-col :span="12">
                    <el-button icon="el-icon-close" type="danger" style="width: 100%" @click="cancelCreateScenario">取消</el-button>
                </el-col>
            </el-row>
            
        </el-dialog>
        <el-button icon="el-icon-plus" type="success" circle id="create-chart-button" @click="showCreateScenarioDialog = true; updatingScenario = false;" />
        <el-popover width="600" placement="left" trigger="click" @after-leave="sceneSearch = ''">
            <el-table :data="validScenarios.filter(data => !sceneSearch || data.name.toLowerCase().includes(sceneSearch.toLowerCase()))" max-height="400px" :stripe="true" @row-dblclick="sceneListRowDbclick">
                <el-table-column fixed>
                    <template slot="header">
                        <el-input v-model="sceneSearch" size="mini" placeholder="输入关键字搜索"/>
                    </template>
                    <template slot-scope="scope">
                        <el-button type="success" size="mini" style="width: 100%" round @click="selectSceneBtnClick(scope.row)">{{ scope.row.name }}</el-button>
                    </template>
                </el-table-column>
                <el-table-column width="50px">
                    <template slot-scope="scope">
                        <i style="color: #409EFF; cursor: pointer" class="el-icon-edit-outline" circle size="mini" @click="editScenario(scope.row)"/>
                    </template>
                </el-table-column>
                <el-table-column width="50px">
                    <template slot-scope="scope">
                        <i style="color: red; cursor: pointer" type="danger" class="el-icon-close" circle size="mini" @click="deleteScenario(scope.row)"/>
                    </template>
                </el-table-column>
            </el-table>
            <el-button slot="reference" icon="el-icon-s-data" type="primary" circle id="show-scene-list-button" />
        </el-popover>
        <el-dialog :visible.sync="showEditPricesDialog" width="40%" @closed="showPriceListSelect = false; priceToEdit = {}">
            <p slot="title" style="font-size: 16px; width: 100%; text-align: center">价格带配置</p>
            <el-row :gutter="8">
                <el-col :span="4">
                    <el-popover width="250" placement="right-start">
                        <div class="price-select-list">
                            <el-button v-for="(p, i) in validPriceList" :key="'price-select-btn-' + i" size="mini" :type="p.useDefault ? 'primary' :'success'" @click="selectPriceList(p)">
                                {{ p.classText }}
                            </el-button>
                        </div>
                        <el-button slot="reference" id="gj-price-list-select-toggle" icon="el-icon-d-arrow-right" size="mini" circle type="primary" @click="showPriceListSelect = !showPriceListSelect"/>
                    </el-popover>
                </el-col>
                <el-col :span="16">
                    <p style="font-size: 22px">{{ priceToEdit.classText }}</p>
                </el-col>
                <el-col :span="4">
                    <el-button v-show="priceToEdit.classCode && !priceToEdit.useDefault && priceToEdit.classCode !== '0000'" size="mini" circle icon="el-icon-refresh-left" type="danger" @click="deletePriceList(priceToEdit)"/>
                </el-col>
            </el-row>
            <el-divider></el-divider>
            <el-form :id="priceEditFormId" label-width="80px" size="mini">
                <el-form-item v-for="(z, i) in validZoneValues" :key="'price-zone-edit-' + i" :label="(i + 1) + ''" v-show="!z.invisible">
                    <el-row :gutter="4" >
                        <el-col :span="6">
                            <el-input-number style="width: 100%" :precision="2" :step="0.1" v-model="z.begin"/>
                        </el-col>
                        <el-col :span="2">
                            <div style="width: 100%; text-align: center">
                            <i class="el-icon-arrow-right" />
                            </div>
                        </el-col>
                        <el-col :span="6">
                            <el-input-number :disabled="z.disableEnd" style="width: 100%" :precision="2" :step="0.1" v-model="z.end"/>
                        </el-col>
                        <el-col :span="2">
                            <i class="el-icon-close" style="color: red; cursor: pointer; padding-left: 8px" @click="removePriceListEntry(i)"/>
                        </el-col>
                        <el-col :span="2">
                            <i class="el-icon-plus" v-if="!z.disableEnd" style="color: green; cursor: pointer; padding-left: 8px" @click="addNewPriceListEntry(i)"/>
                            <i class="el-icon-plus" style="color: white" />
                        </el-col>
                        <el-col :span="2">
                            <el-checkbox v-model="z.inf" @change="infChange(z, i)" />
                        </el-col>
                    </el-row>
                </el-form-item>
                <el-button v-if="validZoneValues" round type="success" style="width: 100%" icon="el-icon-upload" @click="savePriceList" >保存更新</el-button>
                <!--el-form-item v-if="priceToEdit.zones && priceToEdit.zones.values" label="继续添加">
                    <template>
                        <el-button icon="el-icon-plus" round size="mini" type="success" @click="addNewPriceListEntry">添加新价格带</el-button>
                    </template>
                </el-form-item-->
            </el-form>
        </el-dialog>
        <el-button icon="el-icon-edit" type="warning" circle id="edit-prices-button" @click="showEditPricesDialog = true"/>
        <el-button icon="el-icon-setting" type="danger" circle id="edit-scent-button" v-if="validOrgList.length > 0" @click="editScenario(selectedScenario)"/>
        <el-card>
            <div v-if="validOrgList.length > 0">
                <div v-if="showHeader" class="price-chart-title">
                    <i class="el-icon-arrow-up" style="position: absolute; right: 12px; cursor: pointer;" @click="toggleHeader(false)"/>
                    <p>{{ selectedScenario.name }}</p>
                    <el-tag size="mini">{{ validDateRangeDisplay }}</el-tag>
                    <el-tag v-for="org in validOrgList" :key="'tag-' + org.orgCode" type="success" size="mini">{{ org.orgText }}</el-tag>
                    <el-tag v-for="cls in validClassList" :key="'tag-' + cls.classCode" type="warning" size="mini">{{ cls.classText }}</el-tag>
                </div>
                <i v-else class="el-icon-arrow-down" style="position: absolute; right: 40px; top: 35px; cursor: pointer; z-index: 5" @click="toggleHeader(true)"/>
                <el-tabs  class="container-tab" v-model="activeChartTab" :before-leave="beforeTabSwitch" @tab-click="chartTabClick">
                    <el-tab-pane v-for="org in validOrgList" :key="org.orgCode" :label="org.orgText" :name="org.orgCode" :lazy="true">
                        <div v-if="activeChartTab === org.orgCode" class="pane-container">
                            <chart-block v-for="(c, i) in validChartData" :key="org.orgCode + '-chart-block-' + i" :chartKey="org.orgCode + '-' + i" :ref="org.orgCode + '-ref-' + i"
                                :title="c.classText" :rawData="c.pzList" :beginDate="selectedScenario.beginDate" :endDate="selectedScenario.endDate"
                                :orgLevel="org.orgLevel" :orgCode="org.orgCode" :classLevel="c.classLevel" :classCode="c.classCode"></chart-block>
                        </div>
                    </el-tab-pane>
                </el-tabs>
            </div>
            <p v-else style="width: 100%; text-align: center; color: #909399; font-size: 14px;">暂无数据</p>
        </el-card>
    </div>
</template>

<style lang="scss">
    .price-select-list {
        
        max-height: 400px;
        overflow-y: auto;

        .el-button {
            width: 100%;
            margin: 0px 0px 8px 0px
        }
    }

    .price-ana {
        .price-chart-title {
            width: 100%;
            border-radius: 8px;
            background-color: #E9E9EB;
            padding: 10px 10px 10px 10px;
            position: relative;

            p {
                width: 100%;
                margin-bottom: 8px;
                text-align: center;
                color: #909399;
                font-size: 14px;
                font-weight: bold;
            }

            .el-tag {
                margin-right: 4px;
                margin-top: 4px;
            }
        }

        .create-scenario-alert {
            padding: 8px 8px 8px 8px;
            background-color: #FAECD8;
            color: #E62A3C;
            font-size: 11px;
            border-radius: 5px;
            margin-bottom: 10px;
        }

        .el-tabs{
            height: 100%;
            padding: 0;
            margin: 0;
            display: flex;
            flex-direction: column;
        }
        
        .el-tab-pane {
            height: 100%;
        }
        
        .el-tabs__content {
            flex: 1;
            height: calc(100vh - 110px);
        }

        #show-scene-list-button {
            position: absolute;
            bottom: 80px;
            right: 20px;
            opacity: 0.1;
            z-index: 10;
        }

        #show-scene-list-button:hover {
            opacity: 1;
        }

        #create-chart-button {
            position: absolute;
            bottom: 20px;
            right: 20px;
            opacity: 0.1;
            z-index: 10;
        }

        #create-chart-button:hover {
            opacity: 1;
        }

        #edit-prices-button {
            position: absolute;
            bottom: 140px;
            right: 20px;
            opacity: 0.1;
            z-index: 10;
        }

        #edit-prices-button:hover {
            opacity: 1;
        }

        #edit-scent-button {
            position: absolute;
            bottom: 200px;
            right: 20px;
            opacity: 0.1;
            z-index: 10;
        }

        #edit-scent-button:hover {
            opacity: 1;
        }

        .stretch-select {
            width: 100%;
        }

        .pane-container {
            max-height: 100%;
            overflow:auto;
        }
    }
</style>

<script>
import { doRequest } from '../../utils/utils'
import { getUserInfo } from '../../utils/dataStorage'
import { getProdMajorList, getProdMinorList, getProdMicroList, getOrgMinorList, getOrgMicroList } from '../../api/general'
import moment from 'moment'
import ChartBlock from './ChartBlock.vue'
import { Base64 } from 'js-base64'

export default {
    components: {
        'chart-block': ChartBlock
    },
    data() {
        return {
        pickerOptions: {
          shortcuts: [{
            text: '最近一周',
            onClick(picker) {
              const end = new Date();
              const start = new Date();
              start.setTime(start.getTime() - 3600 * 1000 * 24 * 7);
              picker.$emit('pick', [start, end]);
            }
          }, {
            text: '最近两周',
            onClick(picker) {
              const end = new Date();
              const start = new Date();
              start.setTime(start.getTime() - 3600 * 1000 * 24 * 7 * 2);
              picker.$emit('pick', [start, end]);
            }
          }, {
            text: '最近三周',
            onClick(picker) {
              const end = new Date();
              const start = new Date();
              start.setTime(start.getTime() - 3600 * 1000 * 24 * 7 * 3);
              picker.$emit('pick', [start, end]);
            }
          }, {
            text: '最近30天',
            onClick(picker) {
              const end = new Date();
              const start = new Date();
              start.setTime(start.getTime() - 3600 * 1000 * 24 * 30);
              picker.$emit('pick', [start, end]);
            }
          }, {
            text: '最近90天',
            onClick(picker) {
              const end = new Date();
              const start = new Date();
              start.setTime(start.getTime() - 3600 * 1000 * 24 * 90);
              picker.$emit('pick', [start, end]);
            }
          }]
        },
        useRealTime: false,
        realTimeCount: 1,
        realTimeMode: 'd',
        dateRange: '',
        orgMajor: [],
        orgMinor: [],
        orgMicro: [],
        prodMajor: [],
        prodMinor: [],
        prodMicro: [],
        orgListMajor: [{
            orgCode: '0000',
            orgText: '北京公司'
        }],
        orgListMinor: [],
        orgListMicro: [],
        prodListMajor: [],
        prodListMinor: [],
        prodListMicro: [],
        orgMap: {
            '0000': {
                orgLevel: 0,
                orgCode: '0000',
                orgText: '北京公司'
            }
        },
        prodMap: {},
        scenarios: [],
        sceneSearch: '',
        selectedScenario: undefined,
        sceneName: '新建场景',
        showCreateScenarioDialog: false,
        updatingScenario: false,
        activeChartTab: '',
        orgChartStatus: {},
        currentChartData: [],
        showEditPricesDialog: false,
        priceList: [],
        defaultPriceZone: {},
        showPriceListSelect: false,
        priceToEdit: {
            zones: {
                
            }
        },
        priceEditFormId: 'price-edit-form',
        user: '',
        priceListShowUntil: -1,
        showCreateScenarioAlert: false,
        showHeader: true,
      };
    },
    mounted: function() {
        getProdMajorList({
            obj: this,
            src: 'classList',
            dst: 'prodListMajor'
        }, {
            success: res => {
                if (this.prodListMajor) {
                    for (let i = 0; i < this.prodListMajor.length; i ++) {
                        this.prodMap[this.prodListMajor[i].classCode] = this.prodListMajor[i].classText;
                    }
                }
            },
            fail: err => {
                console.log(err)
            }
        });
        getProdMinorList({
            obj: this,
            src: 'classList',
            dst: 'prodListMinor'
        }, {
            success: res => {
                if (this.prodListMinor) {
                    for (let i = 0; i < this.prodListMinor.length; i ++) {
                        this.prodMap[this.prodListMinor[i].classCode] = this.prodListMinor[i].classText;
                    }
                }
            },
            fail: err => {
                console.log(err)
            }
        });
        getProdMicroList({
            obj: this,
            src: 'classList',
            dst: 'prodListMicro'
        }, {
            success: res => {
                if (this.prodListMicro) {
                    for (let i = 0; i < this.prodListMicro.length; i ++) {
                        this.prodMap[this.prodListMicro[i].classCode] = this.prodListMicro[i].classText;
                    }
                }
            },
            fail: err => {
                console.log(err)
            }
        });
        getOrgMinorList({
            obj: this,
            src: 'orgList',
            dst: 'orgListMinor'
        }, {
            success: res => {
                if (this.orgListMinor) {
                    for (let i = 0; i < this.orgListMinor.length; i ++) {
                        this.orgMap[this.orgListMinor[i].orgCode] = this.orgListMinor[i];
                    }
                }
            },
            fail: err => {
                console.log(err)
            }
        });
        getOrgMicroList({
            obj: this,
            src: 'orgList',
            dst: 'orgListMicro'
        }, {
            success: res => {
                if (this.orgListMicro) {
                    for (let i = 0; i < this.orgListMicro.length; i ++) {
                        this.orgMap[this.orgListMicro[i].orgCode] = this.orgListMicro[i];
                    }
                }
            },
            fail: err => {
                console.log(err)
            }
        });

        doRequest({
            url: '/v1/web/scenes/list',
            method: 'get'
        }, {
            obj: this,
            src: 'scenes',
            dst: 'scenarios'
        });

        this.getPriceList();
        this.user = getUserInfo().name;
        window.addEventListener('resize', this.toggleHeader);
    },
    computed: {
        validDateRangeDisplay: function() {
            if (this.selectedScenario) {
                if (this.selectedScenario.beginDate === '' || this.selectedScenario.endDate.endsWith('m') || this.selectedScenario.endDate.endsWith('d')) {
                    return `前${this.selectedScenario.endDate.substring(0, this.selectedScenario.endDate.length - 1)}${this.selectedScenario.endDate.endsWith('m') ? '月' : '天'}`;
                } else {
                    return this.selectedScenario.beginDate + '至' + this.selectedScenario.endDate;
                }
            }
            return '';
        },
        validZoneValues: function() {
            if (this.priceToEdit && this.priceToEdit.zones)
                return this.priceToEdit.zones.values;
            return null;
        },
        validOrgList: function() {
            return (this.selectedScenario && this.selectedScenario.orgInfo) ? this.selectedScenario.orgInfo : [];
        },
        validClassList: function() {
            return (this.selectedScenario && this.selectedScenario.classInfo) ? this.selectedScenario.classInfo : [];
        },
        validScenarios: function() {
            return this.scenarios ? this.scenarios : [];
        },
        validChartData: function() {
            return this.currentChartData ? this.currentChartData : [];
        },
        validPriceList: function() {
            let i = 0, j = 0, prod = null, price = null, match = false, priceList = [{
                classLevel: 0,
                classCode: '0000',
                classText: '默认价格带',
                zones: this.deepCopy(this.defaultPriceZone),
            }];
            for (; i < this.prodListMajor.length; i ++) {
                prod = this.prodListMajor[i];
                match = false;
                for (j = 0; j < this.priceList.length; j ++) {
                    price = this.priceList[j];
                    if (prod.classCode === price.classCode) {
                        match = true;
                        priceList.push(this.deepCopy(Object.assign({}, price, {
                            classText: prod.classText
                        })));
                        break;
                    }
                }

                if (!match) {
                    priceList.push({
                        classLevel: 0,
                        classCode: prod.classCode,
                        classText: prod.classText,
                        zones: this.deepCopy(this.defaultPriceZone),
                        useDefault: true
                    })
                }
            }
            priceList.sort((a, b) => {
                let na = parseInt(a.classCode), nb = parseInt(b.classCode);
                return na - nb;
            })
            return priceList;
        }
    },
    methods: {
        deepCopy(obj) {
            return JSON.parse(JSON.stringify(obj));
        },
        selectSceneBtnClick(row) {
            document.getElementById('show-scene-list-button').click();
            this.selectedScenario = null;
            this.orgChartStatus = {};
            if (row.orgInfo && row.orgInfo.length > 0) {
                for (let i = 0; i < row.orgInfo.length; i ++) {
                    if (row.orgInfo[i].orgLevel === 0) {
                        row.orgInfo[i].orgCode = '0000';
                        row.orgInfo[i].orgText = '北京公司';
                    }
                    let orgCode = row.orgInfo[i].orgCode;
                    this.orgChartStatus[row.orgInfo[i].orgCode] = row.orgInfo[i];
                }
                // this.activeChartTab = row.orgInfo[0].orgCode;
                this.getChartData(row.orgInfo[0].orgLevel, row.orgInfo[0].orgCode, row, true);
            } 
        },
        sceneListRowDbclick(row, col, ev) {
            this.selectSceneBtnClick(row);
        },
        updateScenario() {
            this.createScenario(true);
        },
        createScenario(update) {
            let sdate = '', edate = '', orgs = [], prods = [], i = 0;

            if (this.useRealTime) {
                sdate = '';
                edate = this.realTimeCount + this.realTimeMode;
            } else {
                if (this.dateRange.length > 0) {
                    sdate = moment(this.dateRange[0]).format('YYYY-MM-DD');
                    edate = moment(this.dateRange[1]).format('YYYY-MM-DD');
                }
            }
            
            for (i = 0; i < this.orgMajor.length; i ++) {
                orgs.push({
                    orgLevel: 0,
                    orgCode: this.orgMajor[i]
                })
            }
            for (i = 0; i < this.orgMinor.length; i ++) {
                orgs.push({
                    orgLevel: 1,
                    orgCode: this.orgMinor[i]
                })
            }
            for (i = 0; i < this.orgMicro.length; i ++) {
                orgs.push({
                    orgLevel: 2,
                    orgCode: this.orgMicro[i]
                })
            }

            for (i = 0; i < this.prodMajor.length; i ++) {
                prods.push({
                    classLevel: 0,
                    classCode: this.prodMajor[i]
                })
            }
            for (i = 0; i < this.prodMinor.length; i ++) {
                prods.push({
                    classLevel: 1,
                    classCode: this.prodMinor[i]
                })
            }
            for (i = 0; i < this.prodMicro.length; i ++) {
                prods.push({
                    classLevel: 2,
                    classCode: this.prodMicro[i]
                })
            }

            if (edate === '') {
                this.$message({
                    type: 'error',
                    message: '请选择有效的起始日期！'
                })
                return;
            }

            if (orgs.length === 0) {
                this.$message({
                    type: 'error',
                    message: '请选择至少一个组织机构！'
                })
                return;
            }

            if (prods.length === 0) {
                this.$message({
                    type: 'error',
                    message: '请选择至少一个商品分类！'
                })
                return;
            }

            let data = {
                name: this.sceneName,
                beginDate: sdate,
                endDate: edate,
                orgInfo: orgs,
                classInfo: prods,
                user: this.user
            }

            if (edate.length )

            doRequest({
                url: update ? '/v1/web/scenes/update' : '/v1/web/scenes/add',
                method: 'post',
                data: data,
                loading: true,
            }, {
                success: res => {
                    this.showCreateScenarioDialog = false;
                    doRequest({
                        url: '/v1/web/scenes/list',
                        method: 'get'
                    }, {
                        obj: this,
                        src: 'scenes',
                        dst: 'scenarios'
                    })
                },
                fail: err => {
                    this.showCreateScenarioDialog = false;
                }
            })
        },
        cancelCreateScenario() {
            this.showCreateScenarioDialog = false;
        },
        deleteScenario(scene) {
            this.$confirm(`您确定要删除场景 ${scene.name} 么？`, '删除场景', {
                type: 'warning'
            }).then(() => {
                doRequest({
                    url: `/v1/web/scenes/del/${scene.id}`,
                    method: 'get'
                }, {
                    success: res => {
                        doRequest({
                            url: '/v1/web/scenes/list',
                            method: 'get'
                        }, {
                            obj: this,
                            src: 'scenes',
                            dst: 'scenarios'
                        })
                    },
                    fail: err => {
                        console.log(err);
                        this.$message({
                            type: 'error',
                            message: '删除场景出错，请稍后再试！'
                        })
                    }
                })
            }).catch(() => {

            })
        },
        scenarioPopoverCleanup() {
            this.dateRange = [];
            this.orgMajor = [];
            this.orgMinor = [];
            this.orgMicro = [];
            this.prodMajor = [];
            this.prodMinor = [];
            this.prodMicro = [];
            this.sceneName = '';
            this.useRealTime = false;
            this.realTimeCount = 1;
            this.realTimeMode = 'd';
            this.showCreateScenarioAlert = false;
        },
        scenarioPopoverShow() {
            if (!this.updatingScenario) {
                this.useRealTime = false;
                this.realTimeCount = 1;
                this.realTimeMode = 'd';
                this.dateRange = [];
                this.orgMajor = [];
                this.orgMinor = [];
                this.orgMicro = [];
                this.prodMajor = [];
                this.prodMinor = [];
                this.prodMicro = [];
                this.sceneName = '新建场景 ' + new Date().toLocaleString();
                this.showCreateScenarioAlert = false;
            }
        },
        chartTabClick(tab, ev) {
            let org = this.orgChartStatus[tab.name];
            if (org && !org.chartData) {
                this.getChartData(org.orgLevel, org.orgCode, this.selectedScenario)
            } else {
                this.currentChartData = org.chartData;
            }
        },
        beforeTabSwitch(active, old) {
            this.currentChartData = [];
            return true;
        },
        getChartData(orgLevel, orgCode, scene, setActiveTab) {
            doRequest({
                url: '/v1/web/pricezone',
                method: 'post',
                loading: true,
                data: {
                    beginDate: scene.beginDate,
                    endDate: scene.endDate,
                    orgLevel: orgLevel,
                    orgCode: orgCode,
                    classInfo: scene.classInfo
                }
            }, {
                obj: this,
                src: 'classList',
                dst: 'currentChartData',
                success: res => {
                    this.orgChartStatus[orgCode].chartData = res.classList;
                    if (setActiveTab && scene) {
                        this.activeChartTab = orgCode;
                        this.selectedScenario = scene;
                        this.toggleHeader(this.showHeader);
                    }
                },
                fail: err => {
                    this.currentChartData = [];
                    console.log(err);
                }
            })
        },
        selectPriceList(prod, skip) {
            this.showPriceListSelect = false;
            if (!skip)
                document.getElementById('gj-price-list-select-toggle').click();
            let tempList = this.deepCopy(prod);
            if (tempList && tempList.zones && tempList.zones.values) {
                for (let i = 0; i < tempList.zones.values.length; i ++) {
                    if (tempList.zones.values[i].end > 0)
                        tempList.zones.values[i].inf = false;
                    else {
                        tempList.zones.values[i].inf = true;
                        tempList.zones.values[i].disableEnd = true;
                        break;
                    }
                }
            }
            this.priceToEdit = tempList;
        },
        addNewPriceListEntry(i) {
            if (this.priceToEdit && this.priceToEdit.zones && this.priceToEdit.zones.values && this.priceToEdit.zones.values.length > i) {
                this.priceToEdit.zones.values.splice(i + 1, 0, {
                    begin: 0,
                    end: 0,
                    inf: false,
                })
            }
        }, 
        removePriceListEntry(i) {
            if (this.priceToEdit && this.priceToEdit.zones && this.priceToEdit.zones.values && this.priceToEdit.zones.values.length > i) {
                if (this.priceToEdit.zones.values.length === 1) {
                    this.$message({
                        type: 'warning',
                        message: '请至少保留一项价格带配置！'
                    })
                    return;
                }
                let z = this.priceToEdit.zones.values[i];
                if (z.inf) {
                    for (let n = 0; n < this.priceToEdit.zones.values.length; n ++) {
                        this.priceToEdit.zones.values[n].invisible = false;
                    }
                }
                this.priceToEdit.zones.values.splice(i, 1);
            }
        },
        getPriceList() {
            doRequest({
                url: '/v1/web/pricezone/config/list',
                method: 'get'
            }, {
                success: res => {
                    this.priceList = res.pzcList;
                    let i = 0;
                    for (i = 0; i < this.priceList.length; i ++) {
                        this.priceList[i].zones = JSON.parse(Base64.decode(this.priceList[i].zones));
                        if (this.priceList[i].classCode === '0000') {
                            this.defaultPriceZone = this.priceList[i].zones;
                        }
                    }

                    if (this.showEditPricesDialog) {
                        for (i = 0; i < this.validPriceList.length; i ++) {
                            if (this.validPriceList[i].classCode === this.priceToEdit.classCode) {
                                this.selectPriceList(this.validPriceList[i], true);
                            }
                        }
                    }
                },
                fail: err => {
                    console.log(err)
                }
            });
        },
        savePriceList() {
            if (this.priceToEdit && this.priceToEdit.zones && this.priceToEdit.zones.values) {
                this.priceToEdit.zones.count = this.priceToEdit.zones.length;

                // sanity check
                let v = null, s = 0, e = 0, n = 0, values = [];
                for (let i = 0; i < this.priceToEdit.zones.values.length; i ++) {
                    v = this.priceToEdit.zones.values[i];
                    if (v.begin >= v.end && v.end >= 0) {
                        this.$message({
                            type: 'error',
                            message: `取值错误：第 ${i + 1} 项配置，左值 ${v.begin} >= 右值 ${v.end} ！`
                        })
                        return;
                    }

                    if (v.begin != e) {
                        this.$message({
                            type: 'error',
                            message: `取值不连续：第 ${i + 1} 项配置的左值 ${v.begin} 不等于上一项的 右值 ${e} ！`
                        })
                        return;
                    }

                    n += 1;
                    values.push({
                        begin: v.begin,
                        end: v.inf ? -1 : v.end
                    })
                    if (v.inf)
                        break;

                    s = v.begin;
                    e = v.end;
                }

                doRequest({
                    url: this.priceToEdit.useDefault ? '/v1/web/pricezone/config/add' : '/v1/web/pricezone/config/update',
                    method: 'post',
                    data: {
                        classLevel: this.priceToEdit.classLevel,
                        classCode: this.priceToEdit.classCode,
                        zones: Base64.encode(JSON.stringify({
                            count: n,
                            values: values
                        })),
                        user: this.user
                    }
                }, {
                    success: res => {
                        this.$message({
                            type: 'success',
                            message: '更新价格带配置成功！'
                        });
                        this.getPriceList()
                    },
                    fail: err => {
                        this.$message({
                            type: 'error',
                            message: '更新价格带配置失败，请稍后再试！'
                        })
                        console.log(err)
                    }
                })
            }
        },
        deletePriceList(prod) {
            this.$confirm(`您确定要重置 ${prod.classText} 的价格带配置么？`, '重置价格带', {
                type: 'warning'
            }).then(() => {
                doRequest({
                    url: `/v1/web/pricezone/config/del/${prod.classCode}`,
                    method: 'get'
                }, {
                    success: res => {
                        this.$message({
                            type: 'success',
                            message: '重置价格带配置成功！'
                        });
                        this.getPriceList()
                    },
                    fail: err => {
                        this.$message({
                            type: 'error',
                            message: '重置价格带配置失败，请稍后再试！'
                        })
                    }
                })
            }).catch(() => {

            })
        },
        infChange(z, i) {
            let n = 0;
            if (z.inf === true) {
                for (; n < this.priceToEdit.zones.values.length; n ++) {
                    if (i !== n && this.priceToEdit.zones.values[n].end > 0)
                        this.priceToEdit.zones.values[n].inf = false;
                    if (n > i)
                        this.priceToEdit.zones.values[n].invisible = true;
                }
                this.priceToEdit.zones.values[i].disableEnd = true;
            } else {
                for (n = i + 1; n < this.priceToEdit.zones.values.length; n ++) {
                    this.priceToEdit.zones.values[n].invisible = false;
                    if (this.priceToEdit.zones.values[n].end > 0)
                        this.priceToEdit.zones.values[n].disableEnd = false;
                }
                this.priceToEdit.zones.values[i].disableEnd = false;
            }
            
        },
        editScenario(sc) {
            if (sc.endDate.endsWith('d') || sc.endDate.endsWith('m')) {
                this.useRealTime = true;
                this.realTimeCount = sc.endDate.substring(0, sc.endDate.length - 1);
                this.realTimeMode = sc.endDate[sc.endDate.length - 1];
            } else {
                this.dateRange = [
                    new Date(sc.beginDate),
                    new Date(sc.endDate)
                ];
            }
            
            this.sceneName = sc.name;
            let i = 0, org, prod;
            for (; i < sc.orgInfo.length; i ++) {
                org = sc.orgInfo[i];
                if (org.orgLevel === 0)
                    this.orgMajor.push(org.orgCode);
                else if (org.orgLevel === 1)
                    this.orgMinor.push(org.orgCode);
                else
                    this.orgMicro.push(org.orgCode);
            }

            for (i = 0; i < sc.classInfo.length; i ++) {
                prod = sc.classInfo[i];
                if (prod.classLevel === 0)
                    this.prodMajor.push(prod.classCode);
                else if (prod.classLevel === 1)
                    this.prodMinor.push(prod.classCode);
                else
                    this.prodMicro.push(prod.classCode);
            }

            this.updatingScenario = true;
            this.showCreateScenarioDialog = true;
        },
        hidePriceListSelect(ev) {
            if (ev.target.id !== 'gj-price-list-select-toggle' && ev.target.parentNode.id !== 'gj-price-list-select-toggle') {
                this.showPriceListSelect = false;
            }
        }, 
        createScenarioChangeAlert(e) {
            if (this.orgMajor.length + this.orgMinor.length + this.orgMicro.length + this.prodMajor.length + this.prodMinor.length + this.prodMicro.length > 5)
                this.showCreateScenarioAlert = true;
            else
                this.showCreateScenarioAlert = false;
        },
        toggleHeader(show) {
            if (show === undefined)
                show = this.showHeader;
            this.showHeader = show;
            this.$nextTick(function() {
                let tabs = document.querySelectorAll('.price-ana .el-tabs__content'), i = 0, height = 110;
                if (tabs && tabs.length) {
                    if (show) {
                        let header = document.getElementsByClassName('price-chart-title');
                        if (header && header.length) 
                            height += header[0].clientHeight;
                    }
                } 
                for(; i < tabs.length; i ++) {
                    tabs[i].style.height = `calc(100vh - ${height}px)`;
                }
            })
        }
    }
}
</script>