<template>
  <div class="cnpc-daily-avg-div" v-resize:debounce.50="resizeEverything" >
    <el-dialog :visible.sync="editConfigDlg" @closed="cleanupConfig" @open="prepareConfig"
               :close-on-click-modal="false">
      <p slot="title" style="text-align: center">
        评级标准
      </p>
      <x-table :data="configListShadow" :minus="300" ref="smart-config-table"
               show-overflow="tooltip" border="inner">
        <vxe-table-column title="级别">
          <template slot-scope="r">
            {{ r.rowIndex }}
          </template>
        </vxe-table-column>
        <vxe-table-column title="名称" field="name">
          <template slot-scope="r">
            <el-input size="mini" v-if="r.row.editing" v-model="configListShadow[r.$rowIndex].name">
            </el-input>
            <span v-else>
                            {{ r.row.name }}
                        </span>
          </template>
        </vxe-table-column>
        <vxe-table-column title="收入下限" field="begin">
          <template slot-scope="r">
            <el-input size="mini" v-if="r.row.editing" v-model="configListShadow[r.$rowIndex].begin">
            </el-input>
            <span v-else>
                            {{ r.row.begin }}
                        </span>
          </template>
        </vxe-table-column>
        <vxe-table-column title="收入上限" field="end">
          <template slot-scope="r">
            <el-input size="mini" v-if="r.row.editing" v-model="configListShadow[r.$rowIndex].end">
            </el-input>
            <span v-else>
                            {{ r.row.end }}
                        </span>
          </template>
        </vxe-table-column>
        <vxe-table-column title="操作">
          <template slot-scope="r">
            <i class="el-icon-circle-plus"
               style="color: #3C8D8C; font-size: 22px; cursor: pointer"
               @click="newConfigItem(r.$rowIndex)"/>
            <i v-if="r.row.editing" class="el-icon-refresh-right"
               style="color: #F56C6C; font-size: 22px; cursor: pointer"
               @click="cancelConfigItem(r.$rowIndex)"/>
            <i v-if="r.row.editing" class="el-icon-success"
               style="color: #67C23A; font-size: 22px; cursor: pointer"
               @click="confirmConfigItem(r.$rowIndex)"/>
            <i v-else class="el-icon-edit-outline"
               style="color: #67C23A; font-size: 22px; cursor: pointer"
               @click="editConfigItem(r.$rowIndex)"/>
            <i class="el-icon-error"
               style="color: #F56C6C; font-size: 22px; cursor: pointer"
               @click="deleteConfigItem(r.$rowIndex, r.row)"/>
          </template>
        </vxe-table-column>
      </x-table>
      <el-button size="small" type="primary" v-loading.fullscreen.lock="fullscreenLoading"
                 style="width: 100%; margin-top: 10px"
                 @click="postConfig">提交</el-button>
    </el-dialog>
    <el-card>
      <span style="font-size: 14px">统计时间段 </span>
      <el-date-picker size="small"
                      v-model="dateRange"
                      type="daterange"
                      range-separator="-"
                      start-placeholder="开始日期"
                      end-placeholder="结束日期"
                      value-format="yyyy-MM-dd">
      </el-date-picker>
      <el-button size="small" style="margin-left: 12px" type="primary" :disabled="disableFlag"
                 @click="handleCheckedChange()">计算评级</el-button>
      <!--el-button size="small" style="margin-left: 12px; float: right" type="success"
                 @click="editConfigDlg = true" v-loading="configLoading"
                 :disabled="configLoading">评级配置</el-button-->
      <el-button size="small" style="margin-left: 12px; float: right" type="primary"
                 @click="doExport" v-loading="configLoading"
                 :disabled="configLoading">导出数据</el-button>
    </el-card>
    <el-card style="margin-top: 8px">
      <el-row :gutter="6">
        <el-col :span="3" v-for="(v, k) in configList" :key="'summary-card-' + k" style="padding-bottom: 8px">
          <el-card :body-style="{padding: '10px 10px'}" >
            <p style="font-size: 13px; color: #909399">{{ v.name }}</p>
            <p style="font-size: 24px; text-align: center; margin-top: 8px">{{ v.count }}</p>
          </el-card>
        </el-col>
      </el-row>
    </el-card>
    <vxe-checkbox :disabled="disableFlag" v-model="checked" @change="handleCheckedChange" style="float: left; line-height: 60px; background: white;">
      <span class="checkbox-span">
        剔除毛利率小于5%的销售收入
      </span>
    </vxe-checkbox>
    <vxe-toolbar export refresh custom></vxe-toolbar>
    <vxe-table :data="stations.filter(stationsFilter)" :minus="350" border
               size="small" :loading="tableLoading" show-overflow="tooltip" :export-config="tableExport"
               style="margin-top: 10px" ref="smart-station-table" :max-height="maxHeight1">
      <vxe-table-column type="seq" width="45" fixed/>
      <vxe-table-column title='排名' field="index" width="50" fixed align="center"/>
      <vxe-table-column field="plantName" title="站点名称">
        <template slot="header" slot-scope="scope">
          <el-input
              v-model="filterName"
              size="mini"
              clearable
              placeholder="搜索站点名称"/>
        </template>
      </vxe-table-column>
      <vxe-table-column field="plant" title="油站编码" width="140">
        <template slot="header" slot-scope="scope">
          <el-input
              v-model="filterCode"
              size="mini"
              clearable
              placeholder="搜索油站编码"/>
        </template>
      </vxe-table-column>
      <vxe-table-column field="parentOrg" width="140" title="公司名称">
        <template slot="header" slot-scope="scope">
          <el-select
              v-model="filterSubOrg"
              size="mini"
              multiple
              clearable
              placeholder="过滤分公司"
              @change="stationLevelFilterChange">
            <el-option label="一分公司" value="A13"/>
            <el-option label="二分公司" value="A14"/>
            <el-option label="三分公司" value="A15"/>
            <el-option label="四分公司" value="A16"/>
          </el-select>
        </template>
        <template slot-scope="r">
          <span>{{ subOrgMap[r.row.parentOrg] }}</span>
        </template>
      </vxe-table-column>
      <vxe-table-column title="评级" width="180" align="center" field="zoneName">
        <template slot="header" slot-scope="scope">
          <el-select
              v-model="filterLevel"
              size="mini"
              multiple
              clearable
              placeholder="过滤评级"
              @change="stationLevelFilterChange">
            <el-option v-for="(v, i) in configList" :key="'filter-level-' + i"
                       :label='v.name' :value="i"/>
          </el-select>
        </template>
        <template slot-scope="r">
          <el-tooltip effect="dark" placement="top" :hide-after="1000"
                      :content="stationLevelDesc(r.row.level)">
                        <span v-if="r.row.level >= configList.length - 1" class="badge-gold">
                            {{ configList[r.row.level].name }}
                        </span>
            <span v-else-if="r.row.level == configList.length - 2" class="badge-silver">
                            {{ configList[r.row.level].name }}
                        </span>
            <span v-else-if="r.row.level == configList.length - 3" class="badge-bronze">
                            {{ configList[r.row.level].name }}
                        </span>
            <span v-else class="badge-normal">
                            {{ configList[r.row.level].name }}
                        </span>
          </el-tooltip>
          <font-awesome-icon v-if="r.row.ylevel>0" icon="arrow-up" class="percent-arrow-up" />
          <font-awesome-icon v-else icon="arrow-down" class="percent-arrow-down" />
        </template>
      </vxe-table-column>
      <vxe-table-column title="去年同期评级" field="yzoneName" align="center" width="150">
        <template slot-scope="r">
          <el-tooltip effect="dark" placement="top" :hide-after="1000"
                      :content="lastRank[r.row.index-1]">
                        <span v-if="lastRank[r.row.index-1]=== configList[configList.length - 1].name" class="badge-gold">
                            {{ lastRank[r.row.index-1] }}
                        </span>
            <span v-else-if="lastRank[r.row.index-1]=== configList[configList.length - 2].name" class="badge-silver">
                            {{ lastRank[r.row.index-1] }}
                        </span>
            <span v-else-if="lastRank[r.row.index-1]=== configList[configList.length - 3].name" class="badge-bronze">
                            {{ lastRank[r.row.index-1] }}
                        </span>
            <span v-else class="badge-normal">
                            {{ lastRank[r.row.index-1] }}
                        </span>
          </el-tooltip>
        </template>
      </vxe-table-column>
      <vxe-table-column title="日均销售额" width="150" field="sales" sortable>
        <template slot-scope="r">
          {{ beautyNum(r.row.sales) }}
        </template>
      </vxe-table-column>
      <vxe-table-column title="营业天数" width="100" field="saleDays" sortable/>
      <vxe-table-column width="140" align="center" title="所有汽服状态">
        <template slot="header" slot-scope="scope">
          <el-select
              v-model="filterCar"
              size="mini"
              placeholder="过滤汽服">
            <el-option label="所有汽服状态" :value='0'/>
            <el-option label='有汽服' :value='1'/>
            <el-option label='无汽服' :value='2'/>
            <el-option label='无收入' :value='3'/>
          </el-select>
        </template>
        <template slot-scope="r">
          <el-tooltip v-if="r.row.carServSR && r.row.carServLR" :hide-after="1000"
                      placement="top" content="有汽服有收入">
            <font-awesome-icon icon="car"
                               color="#67C23A" size="2x"/>
          </el-tooltip>
          <el-tooltip v-else-if="!r.row.carServSR && !r.row.carServLR" :hide-after="1000"
                      placement="top" content="无汽服无收入">
            <font-awesome-icon icon="car"
                               color="#909399" size="2x"/>
          </el-tooltip>
          <el-tooltip v-else :hide-after="1000" placement="top" content="有汽服无收入">
            <font-awesome-icon icon="car" color="red" size="2x"/>
          </el-tooltip>
        </template>
      </vxe-table-column>
    </vxe-table>
  </div>
</template>

<style lang="scss">
.cnpc-daily-avg-div {
  width: 100%;
  max-height: calc(100vh - 12px);
  background-color: white;
  padding: 8px 8px 8px 8px;
  overflow: auto;

  .badge-gold {
    padding: 6px 8px 6px 8px;
    border-radius: 6px;
    color: white;
    background-color: #9C27B0;
    font-size: 13px;
  }

  .badge-silver {
    padding: 6px 8px 6px 8px;
    border-radius: 6px;
    color: white;
    background-color: #2196F3;
    font-size: 13px;
  }

  .badge-bronze {
    padding: 6px 8px 6px 8px;
    border-radius: 6px;
    color: white;
    background-color: #FF9800;
    font-size: 13px;
  }

  .badge-normal {
    padding: 6px 8px 6px 8px;
    border-radius: 6px;
    color: white;
    background-color: #B5A642;
    font-size: 13px;
  }

  .checkbox-span {
    opacity: 0.5;
    &:hover {
      opacity: 1;
    }
  }

  .percent-arrow-up {
    margin-left: 10px;
    font-size: 18px;
    font-weight: bold;
    color: #67C23A;
  }

  .percent-arrow-down {
    margin-left: 10px;
    font-size: 18px;
    font-weight: bold;
    color: #F56C6C;
  }
}
</style>

<script>
const moment = require('moment');
import XTable from '../mixins/SmartVxeTable';
import { randomNumber, beautifyNumber, deepCopy, doRequest, allRequests } from '../../utils/utils'
const axios = require('axios');
import { saveAs } from 'file-saver';
import resize from 'vue-resize-directive'

export default {
  directives: {
      resize
  },
  components: {
    'x-table': XTable
  },
  mounted() {
    let now = new Date(), then = new Date();
    then.setDate(then.getDate() - 30);
    this.dateRange.push(moment(then).format('YYYY-MM-DD'));
    this.dateRange.push(moment(now).format('YYYY-MM-DD'));
    this.handleCheckedChange()
    // this.getConfigList();
    // this.getStations();
  },
  data() {
    return {
      checked: false,
      disableFlag: false,
      maxHeight1: '300px',
      summary: {},
      exportBlob: '',
      subOrgMap: {
        'A13': '一分公司',
        'A14': '二分公司',
        'A15': '三分公司',
        'A16': '四分公司',
      },
      dateRange: [],
      configList: [],
      configListShadow: [],
      configItemToDelete: [],
      configLoading: true,
      stations: [],
      tableLoading: false,
      levelTagTypes: [
        'info',
        'danger',
        'warning',
        'primary',
        'success',
      ],
      tagTypeLevels: [],
      filterName: '',
      filterCode: '',
      filterLevel: [],
      filterSubOrg: [],
      filterCar: 0,
      editConfigDlg: false,
      fullscreenLoading: false,
      tableExport: {
        // 默认选中类型
        type: 'xlsx',
        // 自定义类型
        types: ['xlsx', 'csv', 'html', 'xml', 'txt'],
        message: true
      },
      // 去年的日期
      ysdate: '',
      yedate: '',
      // 去年的评级
      lastRank: [],
    }
  },
  methods: {
    doExport() {
      if (this.exportBlob.length > 0) {
        saveAs(new Blob([this.exportBlob], {type: 'text/plain;charset=utf-8'}), `${this.dateRange[0]}-${this.dateRange[1]}.csv`)
      }
    },
    getDiff() {
      let st = this.stations;
      let yst = this.ystations;
      this.lastRank = []
      let  i, j
      for (i = 0; i < st.length; ++i) {
        for (j = 0; j < yst.length; ++j) {
          if (st[i]['plant'] === yst[j]['plant']) {
            st[i]['ylevel'] = st[i]['sales'] - yst[j]['sales'];
            st[i]['yzoneName']=yst[j]['zoneName'];
            this.lastRank.push(yst[j]['zoneName'])
            break;
          }
        }
        if(j == yst.length){
          this.lastRank.push("没有去年数据")
        }
      }
      this.stations = st;
    },
    getStations() {
      this.tableLoading = true;
      this.disableFlag =true;
      let year = this.dateRange[0].substring(0, 4) - 1
      let year1 = this.dateRange[1].substring(0, 4) - 1
      this.ysdate = year + this.dateRange[0].substring(4)
      this.yedate = year1 + this.dateRange[1].substring(4)

      this.configList.forEach( c=> {
        c.count = 0
      })
      let posts = []
      // 今年的销售数据请求
      posts.push({
        url: '/v1/web/plant/sales-zone',
        method: 'POST',
        data: {
          beginDate: this.dateRange[0],
          endDate: this.dateRange[1]
        }
      });
      // 去年的销售数据请求
      posts.push({
        url: '/v1/web/plant/sales-zone',
        method: 'POST',
        data: {
          beginDate: this.ysdate,
          endDate: this.yedate
        }
      });
      if (posts.length > 0) {
        allRequests({
          success: res => {
            console.log(res)
            this.exportBlob = '排名,加油站,油站编码,分公司,评级,去年同期评级,日均销售额,营业天数,汽服状态,\n'
            res[0].data.pszList.sort((a, b) => {
              return b.sales - a.sales;
            });
            res[1].data.pszList.sort((a, b) => {
              return b.sales - a.sales;
            });
            this.stations = res[0].data.pszList;
            this.ystations = res[1].data.pszList;
            for (let j = 0, len = this.stations.length; j < len; j++) {
              this.stations[j]['ylevel'] = 1;
            }
            this.getDiff()
            for (let i = 0; i < res[0].data.pszList.length; i++) {
              let p = res[0].data.pszList[i];
              p.index = i + 1;
              p.level = this.stationLevel(p);
              this.exportBlob += `${i + 1},${p.plantName},${p.plant},${this.subOrgMap[p.parentOrg]},${p.zoneName},${this.lastRank[i]},${p.sales},${p.saleDays},${p.carServSR ? (p.carServLR ? '有汽服有收入' : '有汽服无收入') : (p.carServLR ? '有汽服无收入' : '无汽服无收入')},\n`
            }
          },
          fail: _ => {
            this.exportBlob = ''
          },
          finally: () => {
            this.tableLoading = false;
            this.disableFlag = false;
          }
        },[], posts)
      } else {
        this.tableLoading = false;
      }
    },
    // 获得剔除5%毛利率的加油站数据
    getThresholdStations() {
      this.tableLoading = true;
      this.disableFlag = true;
      let year = this.dateRange[0].substring(0, 4) - 1
      let year1 = this.dateRange[1].substring(0, 4) - 1
      this.ysdate = year + this.dateRange[0].substring(4)
      this.yedate = year1 + this.dateRange[1].substring(4)

      this.configList.forEach( c=> {
        c.count = 0
      })
      let posts = []
      // 今年的销售数据请求
      posts.push({
        url: '/v1/web/plant/sales-zone',
        method: 'POST',
        data: {
          beginDate: this.dateRange[0],
          endDate: this.dateRange[1],
          rateThreshold: 0.05
        }
      });
      // 去年的销售数据请求
      posts.push({
        url: '/v1/web/plant/sales-zone',
        method: 'POST',
        data: {
          beginDate: this.ysdate,
          endDate: this.yedate,
          rateThreshold: 0.05
        }
      });
      if (posts.length > 0) {
        allRequests({
          success: res => {
            console.log(res)
            this.exportBlob = '排名,加油站,油站编码,分公司,评级,去年同期评级,日均销售额,营业天数,汽服状态,\n'
            res[0].data.pszList.sort((a, b) => {
              return b.sales - a.sales;
            });
            res[1].data.pszList.sort((a, b) => {
              return b.sales - a.sales;
            });
            this.stations = res[0].data.pszList;
            this.ystations = res[1].data.pszList;
            for (let j = 0, len = this.stations.length; j < len; j++) {
              this.stations[j]['ylevel'] = 1;
            }
            this.getDiff()
            for (let i = 0; i < res[0].data.pszList.length; i++) {
              let p = res[0].data.pszList[i];
              p.index = i + 1;
              p.level = this.stationLevel(p);
              this.exportBlob += `${i + 1},${p.plantName},${p.plant},${this.subOrgMap[p.parentOrg]},${p.zoneName},${this.lastRank[i]},${p.sales},${p.saleDays},${p.carServSR ? (p.carServLR ? '有汽服有收入' : '有汽服无收入') : (p.carServLR ? '有汽服无收入' : '无汽服无收入')},\n`
            }
          },
          fail: _ => {
            this.exportBlob = ''
          },
          finally: () => {
            this.tableLoading = false;
            this.disableFlag =false;
          }
        },[], posts)
      } else {
        this.tableLoading = false;
      }
    },
    getConfigList() {
      this.configLoading = true;
      doRequest({
        url: '/v1/web/plant/sales-zone/config/list',
        method: 'GET'
      }, {
        success: res => {
          this.configList = res.pszList.sort((a, b) => {
            return a.begin - b.begin;
          })
          this.configList.forEach(c => {
            c.count = 0
          })
          this.configListShadow = deepCopy(this.configList);
          let i = 0, j = 0, gap = 1;
          this.tagTypeLevels = [];

          if (this.configList.length > this.levelTagTypes.length) {
            gap = Math.round(this.configList.length / this.levelTagTypes.length);
          }
          while (i < this.configList.length && j < this.levelTagTypes.length) {
            this.tagTypeLevels.push({
              start: i,
              type: this.levelTagTypes[j]
            })
            i += gap;
            j += 1;
          }
          this.getStations()
        },
        finally: () => {
          this.configLoading = false;
        }
      });
    },
    // 获得剔除5%毛利率的list
    getThresholdConfigList() {
      this.configLoading = true;
      doRequest({
        url: '/v1/web/plant/sales-zone/config/list',
        method: 'GET',
      }, {
        success: res => {
          this.configList = res.pszList.sort((a, b) => {
            return a.begin - b.begin;
          })
          this.configList.forEach(c => {
            c.count = 0
          })
          this.configListShadow = deepCopy(this.configList);
          let i = 0, j = 0, gap = 1;
          this.tagTypeLevels = [];

          if (this.configList.length > this.levelTagTypes.length) {
            gap = Math.round(this.configList.length / this.levelTagTypes.length);
          }
          while (i < this.configList.length && j < this.levelTagTypes.length) {
            this.tagTypeLevels.push({
              start: i,
              type: this.levelTagTypes[j]
            })
            i += gap;
            j += 1;
          }

          this.getThresholdStations()
        },
        finally: () => {
          this.configLoading = false;
        }
      });
    },
    getTagTypeByLevel(level) {
      for (let i = this.tagTypeLevels.length - 1; i >= 0 ; i --) {
        if (level >= this.tagTypeLevels[i].start)
          return this.tagTypeLevels[i].type;
      }
    },
    newConfigItem(index) {
      this.configListShadow.splice(index + 1, 0, {
        id: -1,
        name: '新加项',
        begin: 0,
        end: 0,
        editing: true,
        isNew: true,
      });
    },
    confirmConfigItem(index) {
      if (this.configItemCheck()) {
        this.configListShadow[index].editing = false;
        this.$refs['smart-config-table'].refresh(true);
      }
    },
    cancelConfigItem(index) {
      this.configListShadow[index].editing = false;
      this.$refs['smart-config-table'].refresh(true);
    },
    editConfigItem(index) {
      this.configListShadow[index].editing = true;
      this.configListShadow[index].edited = true;
      this.$refs['smart-config-table'].refresh(true);
    },
    deleteConfigItem(index, item) {
      this.$confirm(`您确定要删除评级配置项 ${item.name} 么?`, '删除评级配置', {
        type: 'warning'
      }).then(() => {
        if (!this.configListShadow[index].isNew)
          this.configItemToDelete.push(this.configListShadow[index].id);
        this.configListShadow.splice(index, 1);
      }).catch(() => {
      })
    },
    configItemCheck(skip = -1) {
      if (this.configListShadow.length === 0) {
        this.$message({
          type: 'error',
          message: '请至少保留一项评级配置!'
        });
        return false;
      }

      let e = this.configListShadow[0].end;
      for (let i = 0; i < this.configListShadow.length; i ++) {
        if (i != skip) {
          if (i >= 1) {
            if (this.configListShadow[i].begin != e) {
              this.$message({
                type: 'error',
                message: `${this.configListShadow[i].name} 的最小值 ` +
                    `${this.configListShadow[i].begin} 与前一项的最大值 ${e} 不相等`
              });
              return false;
            }
            e = this.configListShadow[i].end;
          }

          if (parseInt(this.configListShadow[i].begin) >= parseInt(this.configListShadow[i].end)) {
            this.$message({
              type: 'error',
              message: `当前项 ${this.configListShadow[i].name} 的最小值 ` +
                  `${this.configListShadow[i].begin} 不小于其最大值 ` +
                  `${this.configListShadow[i].end}`
            });
            return false;
          }
        }
      }

      return true;
    },
    prepareConfig() {
      this.configListShadow = deepCopy(this.configList);
    },
    cleanupConfig() {
      this.configItemToDelete = [];
    },
    postConfig() {
      if (!this.configItemCheck())
        return;

      this.fullscreenLoading = true;
      let i = 0, gets = [], posts = [];

      for (i = 0; i < this.configItemToDelete.length; i ++) {
        if (this.configItemToDelete[i] >= 0)
          gets.push({
            url: `/v1/web/plant/sales-zone/config/del/${this.configItemToDelete[i]}`
          })
      }

      for (i = 0; i < this.configListShadow.length; i ++) {
        if (this.configListShadow[i].isNew) {
          posts.push({
            url: '/v1/web/plant/sales-zone/config/add',
            data: {
              name: this.configListShadow[i].name,
              begin: parseInt(this.configListShadow[i].begin),
              end: parseInt(this.configListShadow[i].end)
            }
          });
        }

        if (this.configListShadow[i].edited) {
          posts.push({
            url: '/v1/web/plant/sales-zone/config/update',
            data: {
              id: this.configListShadow[i].id,
              name: this.configListShadow[i].name,
              begin: parseInt(this.configListShadow[i].begin),
              end: parseInt(this.configListShadow[i].end)
            }
          });
        }
      }

      if (gets.length > 0 || posts.length > 0) {
        allRequests({
          success: res => {
            this.fullscreenLoading = false;
            this.editConfigDlg = false;
            this.getConfigList();
          },
          fail: err => {
            console.log(err);
          }
        }, gets, posts)
      } else {
        this.fullscreenLoading = false;
        this.editConfigDlg = false;
      }
    },
    stationsFilter(d) {
      // 按照站点名筛选
      if (this.filterName.length > 0) {
        // if (!d.plantName.toLowerCase().includes(this.filterName.toLowerCase()))
        //     return false;
        let names = this.filterName.split(' ')
        let res = !d.plantName.includes(names[0])
        for(let i=1; i<names.length; i++){
          res = res && (!d.plantName.includes(names[i]))
        }
        if(res){
          return false;
        }
      }
      // 按照编码筛选
      if (this.filterCode.length > 0) {
        // if (!d.plantName.toLowerCase().includes(this.filterName.toLowerCase()))
        //     return false;
        let names = this.filterCode.split(' ')
        let res = !d.plant.includes(names[0].toUpperCase())
        for (let i = 1; i < names.length; i++) {
          res = res && (!d.plant.includes(names[i].toUpperCase()))
        }
        if (res) {
          return false;
        }
      }

      switch(this.filterCar) {
        case 1:
          if (!d.carServSR || !d.carServLR)
            return false;
          break;

        case 2:
          if (d.carServSR || d.carServLR)
            return false;
          break;

        case 3:
          if (!(d.carServLR && !d.carServSR))
            return false;
          break;

        default:
          break;
      }

      let b = false;
      if (this.filterLevel.length > 0) {
        for (let i = 0; i < this.filterLevel.length; i++) {
          if (d.level === this.filterLevel[i]) {
            b = true;
            break;
          }
        }

        if (!b)
          return false;
      }


      if (this.filterSubOrg.length > 0) {
        b = false;
        for (let i = 0; i < this.filterSubOrg.length; i++) {
          if (d.parentOrg === this.filterSubOrg[i]) {
            b = true;
            break;
          }
        }

        if (!b)
          return false;
      }

      // console.log(d.parentOrg, this.filterSubOrg)
      return true;
    },
    stationLevel(v) {
      let i = 0;
      for (; i < this.configList.length; i ++) {
        // if (v.sales >= this.configList[i].begin && v.sales < this.configList[i].end) {
        //     v.level = i;
        //     return i;
        // }
        if (v.zoneID == this.configList[i].id) {
          v.level = i;
          this.configList[i].count += 1
          return i;
        }
      }
      v.level = i - 1;
      return i - 1;
    },
    stationLevelDesc(level) {
      let e = this.configList[this.configList.length - 1];
      if (level >= this.configList.length)
        return `> ${e.end}`;

      e = this.configList[level];
      return `${e.name}: ${e.begin}-${e.end}`;
    },
    stationLevelFilterChange(v) {
      this.$refs['smart-station-table'].refresh();
      //this.$set(this, 'stations', this.stations)
    },
    beautyNum(num) {
      return beautifyNumber(num);
    },
    resizeEverything() {
      this.maxHeight1 = (window.innerHeight - 410) + 'px'
    },
    // 处理剔除毛利选框是否选中
    handleCheckedChange() {
      if (!this.dateRange || !this.dateRange[0] || !this.dateRange[1]) {
        return;
      }
      if(this.checked){
        this.getThresholdConfigList()
      }else {
        this.getConfigList()
      }
    },

  }
}
</script>
