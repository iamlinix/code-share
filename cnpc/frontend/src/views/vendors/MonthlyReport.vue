<template>
  <div class="vendor-monthly-report">
    <el-dialog :visible.sync="modifyDialog" :title="'修改 ' + modifyVendorName + ' 的 ' + modifyOps[modify.type]" width="40%">
      <el-form size="mini" label-width="120px">
        <el-form-item label="修改人">
          <span style="color: #909399">{{ modifyUser }}</span>
        </el-form-item>
        <el-form-item label="修改值">
          <el-input-number style="width: 160px" v-model="modifyNewValue" />
        </el-form-item>
        <el-form-item label="修改原因">
          <el-input type="textarea" :rows="2" placeholder="请输入修改原因" v-model="modifyReason" />
        </el-form-item>
      </el-form>
      <el-table :data="modifyHistory" size="mini" max-height="250">
        <el-table-column prop="user" label="修改人" width="80"/>
        <el-table-column prop="oldValue" label="修改前"  width="110"/>
        <el-table-column prop="newValue" label="修改后"  width="110"/>
        <el-table-column prop="mtime" label="修改时间" width="150"/>
        <el-table-column prop="reason" label="修改原因"  />
      </el-table>
      <el-button size="mini" round style="width: 100%; margin-top: 8px" type="success" @click="doOverWrite">确认修改</el-button>
    </el-dialog>
    <el-dialog :visible.sync="vendorConfigAddDlg"
               :title="updatingVendorConfig ? `${vendorConfig.vendorName}` : '添加供应商配置'"
               @closed="vendorConfigDlgClosed" width="50%" append-to-body>
      <p slot="title" class="title">
        {{ updatingVendorConfig ? vendorConfig.vendorName : '添加供应商配置' }}
      </p>
      <el-form class="vendor-config-form" size="small" label-position="left"
               label-width="100px" inline>
        <!--el-form-item label="供应商名称">
            <el-input v-model="vendorConfigName" style="width: 80%" />
        </el-form-item-->
        <el-form-item label="供应商编码">
          <el-input v-model="vendorConfig.vendorCode" :disabled="updatingVendorConfig"/>
        </el-form-item>
        <el-divider/>
        <el-form-item label="上月挂账">
          <el-input-number v-model="vendorConfig.lastMonthUnpaid" :disabled="vendorConfig.inactive"
                           :min="0.0" style="width: 80%" :precision="2" :step="0.01"></el-input-number>
        </el-form-item>
        <el-form-item label="上月结余">
          <el-input-number v-model="vendorConfig.lastMonthLeft" :disabled="vendorConfig.inactive"
                           :min="0.0" style="width: 80%" :precision="2" :step="0.01"></el-input-number>
        </el-form-item>
        <el-divider/>
        <el-form-item label="本年采购">
          <el-input-number v-model="vendorConfig.yearPurchase" :disabled="vendorConfig.inactive"
                           :min="0.0" style="width: 80%" :precision="2" :step="0.01"></el-input-number>
        </el-form-item>
        <el-form-item label="本年销售">
          <el-input-number v-model="vendorConfig.yearSales" :disabled="vendorConfig.inactive"
                           :min="0.0" style="width: 80%" :precision="2" :step="0.01"></el-input-number>
        </el-form-item>
        <el-form-item label="本年结算">
          <el-input-number v-model="vendorConfig.yearPaid" :disabled="vendorConfig.inactive"
                           :min="0.0" style="width: 80%" :precision="2" :step="0.01"></el-input-number>
        </el-form-item>
        <el-form-item label="本年挂账">
          <el-input-number v-model="vendorConfig.yearUnpaid" :disabled="vendorConfig.inactive"
                           :min="0.0" style="width: 80%" :precision="2" :step="0.01"></el-input-number>
        </el-form-item>
        <el-divider/>
        <el-form-item label="返点(%)">
          <el-input-number v-model="vendorConfig.rebate" :disabled="vendorConfig.inactive"
                           :min="0.0" :max="100.0" style="width: 80%" :precision="2" :step="0.01"></el-input-number>
        </el-form-item>
        <el-divider/>
        <el-form-item >
          <el-checkbox v-model="vendorConfig.inactive"
                       class="vendor-config-checkbox-frame">非活跃供应商</el-checkbox>
        </el-form-item>
        <el-form-item >
          <el-checkbox v-model="vendorConfig.isNew"
                       :disabled="vendorConfig.inactive"
                       class="vendor-config-checkbox-frame">新供应商</el-checkbox>
        </el-form-item>
        <el-button v-if="updatingVendorConfig" style="width: 100%" type="success"
                   @click="updateVendorConfig">
          更新配置
        </el-button>
        <el-button v-else style="width: 100%" type="primary"
                   @click="addVendorConfig">
          添加配置
        </el-button>
      </el-form>
    </el-dialog>
    <el-dialog :visible.sync="vendorConfigDlg" width="80%">
      <p slot="title" class="title">
        供应商配置
      </p>
      <el-row :gutter="6" style="margin-bottom: 8px">
        <el-col :span="4">
          <el-button size="mini" style="width: 100%"
                     icon="el-icon-download" type="success"
                     @click="downloadVendorConfigCSV">导出供应商配置</el-button>
        </el-col>
        <el-col :span="4">
          <el-upload accept="*" :headers="{'token': getTokenVendor()}"
                     :action="uploadUrl" :show-file-list="false" :multiple="false"
                     :before-upload="vendorConfigUploadBefore"
                     :on-success="vendorConfigUploadSuccess"
                     :on-error="vendorConfigUploadError"
                     name="upload">
            <el-button style="width: 100%" size="mini" icon="el-icon-upload2" type="primary">导入供应商配置</el-button>
          </el-upload>
        </el-col>
        <el-col :span="4">
          <el-button size="mini" style="width: 100%"
                     icon="el-icon-plus" type="warning"
                     @click="showVendorConfigDlg">添加供应商配置</el-button>
        </el-col>
        <el-col :span="4">
          <el-button size="mini" style="width: 100%"
                     icon="el-icon-delete" type="danger"
                     @click="clearVendorConfig">清空供应商配置</el-button>
        </el-col>
      </el-row>

      <smart-table :stripe="true" size="medium" :data="vendorConfigList" :minus="300"
                   :border="true" :loading="configTableLoading">
        <el-table-column width="45" type="index" fixed align="center"/>
        <el-table-column sortable prop="vendorCode" label="供应商编码" fixed width="120"/>
        <el-table-column sortable prop="vendorName" label="供应商名称" width="300"/>
        <el-table-column sortable prop="lastMonthUnpaid" label="上月挂账" width="110"/>
        <el-table-column sortable prop="lastMonthLeft" label="上月结余" width="110"/>
        <el-table-column sortable prop="yearPurchase" label="本年采购" width="110"/>
        <el-table-column sortable prop="yearSales" label="本年销售" width="110"/>
        <el-table-column sortable prop="yearPaid" label="本年结算" width="110"/>
        <el-table-column sortable prop="yearUnpaid" label="本年挂账" width="110"/>
        <el-table-column sortable prop="rebate" label="返点(%)" width="100"/>
        <el-table-column sortable prop="inactive" label="是否活跃" width="110" align="center">
          <template slot-scope="scope">
            <i class="el-icon-error" v-if="scope.row.inactive"
               style="color: #F56C6C; font-size: 20px" />
            <i class="el-icon-success" v-else
               style="color: #67C23A; font-size: 20px" />
          </template>
        </el-table-column>
        <el-table-column sortable prop="inactive" label="新供应商" width="110" align="center">
          <template slot-scope="scope">
            <i class="el-icon-success" v-if="scope.row.isNew"
               style="color: #67C23A; font-size: 20px" />
            <i class="el-icon-error" v-else
               style="color: #F56C6C; font-size: 20px" />
          </template>
        </el-table-column>
        <el-table-column label="操作" width="100" align="center">
          <template slot-scope="scope">
            <el-button size="mini" circle type="primary" icon="el-icon-edit"
                       @click="showVendorConfigDlgUpdate(scope.row)"/>
            <el-button size="mini" circle type="warning" icon="el-icon-close"
                       @click="deleteVendorConfig(scope.row.vendorCode, scope.row.vendorName)"/>
          </template>
        </el-table-column>
      </smart-table>
    </el-dialog>
    <el-dialog :visible.sync="createReportDialog" class="create-report-dialog"
               @closed="createReportDialogClosed" width="40%">

      <p slot="title" class="title">
        创建月度报表
      </p>
      <el-form size="small" label-width="100px" v-model="createReportForm" 
        v-loading="createReportDialogLoading" 
        element-loading-text="报表生成需要较长时间，请耐心等候，不要关闭窗口...">
        <el-form-item label="报表月份" >
          <el-date-picker type="month" v-model="createReportForm.month" />
        </el-form-item>
        <el-form-item label="起止时间">
          <el-date-picker
              v-model="createReportForm.dateRange"
              type="daterange"
              range-separator="-"
              start-placeholder="开始日期"
              end-placeholder="结束日期">
          </el-date-picker>
        </el-form-item>
        <!--el-form-item label="选择供应商">
            <el-transfer :data="vendors" :v-model="createReportForm.excludeVendors"
                :titles="['包含', '不包含']"/>
        </el-form-item!-->
        <el-button type="success" style="width: 100%"
                   @click="createMonthlyReport">创建报表</el-button>
      </el-form>
    </el-dialog>
    <el-tabs v-model="activeTab" type="border-card" stretch @tab-remove="removeTab" @tab-click="tabClick">
      <el-tab-pane name="monthly-report-config">
        <span slot="label"><i class="el-icon-s-tools"></i>报表配置</span>
        <div class="report-list-section">
          <!--p class="section-header" slot="header">
              进销存月度报表
          </p-->
          <el-steps :active="9" align-center>
            <el-step title="创建月度报表"></el-step>
            <el-step title="导出对账单" ></el-step>
            <el-step title="与中央仓对账，核对“本月采购”数据"></el-step>
            <el-step title="核对“本月销售”数据" ></el-step>
            <el-step title="导出中央仓入库汇总表"></el-step>
            <el-step title="与供货商核对“本月开票”数据"></el-step>
            <el-step title="核对“本月结算”数据"></el-step>
            <el-step title="核对“累计挂账”数据"></el-step>
            <el-step title="导出结算台账、导出财务审批"></el-step>
          </el-steps>
          <br/>
          <el-button size="small" type="primary"
                     style="margin-bottom: 12px" icon="el-icon-plus"
                     @click="createReportDialog = true">创建月度报表</el-button>
          <!--el-button size="small" type="success"
              style="margin-bottom: 12px" icon="el-icon-setting"
              @click="vendorConfigDlg = true">供应商配置</el-button-->
          <el-card>
            <x-table size="midium" :stripe="true" :data="monthlyReportList"
                     :minus="180" :border="false"
                     v-loading="reportListLoading">
              <vxe-table-column sortable width="200" title="报表月份" field="month" />
              <vxe-table-column width="200" title="起止时间">
                <template slot-scope="scope">
                  {{ scope.row.beginDate }} ~ {{ scope.row.endDate }}
                </template>
              </vxe-table-column>
              <vxe-table-column title="操作">
                <template slot-scope="scope">
                  <el-button-group style="margin-right: 6px">
                    <el-button size="mini" type="success"
                               @click="viewReport(scope.row, true)">查看活跃</el-button>
                    <el-button size="mini" type="success"
                               @click="viewReport(scope.row, false)">查看所有</el-button>
                  </el-button-group>
                  <el-button-group style="margin-right: 6px">
                    <el-button size="mini" type="primary"
                               @click="exportReport(scope.row, true)">导出活跃</el-button>
                    <el-button size="mini" type="primary"
                               @click="exportReport(scope.row, false)">导出所有</el-button>
                  </el-button-group>
                  <el-button-group>
                    <el-button size="mini" type="danger"
                               @click="deleteReport(scope.row)">删除</el-button>
                  </el-button-group>
                </template>
              </vxe-table-column>
            </x-table>
          </el-card>
        </div>
      </el-tab-pane>
      <el-tab-pane closable v-for="v in reportTabs" :key="v.key" :label="v.label"
                   :name="v.name" v-loading="v.loading">
        <div v-if="!v.detail" class="report-detail-section">
          <!--p class="section-header" slot="header">
              报表详细信息
          </p-->
          <div class="haha">
          <el-steps align-center :active="step" finish-status="success" size="small">
            <el-step title="创建月度报表" ></el-step>
            <el-step title="导出对账单" ></el-step>
            <el-step title="与中央仓对账，核对“本月采购”数据"></el-step>
            <el-step title="核对“本月销售”数据" ></el-step>
            <el-step title="导出中央仓入库汇总表"></el-step>
            <el-step title="与供货商核对“本月开票”数据"></el-step>
            <el-step title="核对“本月结算”数据"></el-step>
            <el-step title="核对“累计挂账”数据"></el-step>
            <el-step title="导出结算台账、导出财务审批"></el-step>
          </el-steps>
          <br/>
          </div>
          <el-tag style="margin-right: 10px" type="primary">{{ v.beginDate + ' ~ ' + v.endDate }}</el-tag>
          <!--导出月度报表-->
          <el-button size="small" type="success"
                     style="margin-bottom: 10px" icon="el-icon-download"
                     @click="exportReport(v, v.activeOnly, true)">导出全部
          </el-button>
          <!--导出商品明细-->
          <el-button size="small" type="success"
                     style="margin-bottom: 10px" icon="el-icon-download"
                     @click="exportDetailReport(v, v.activeOnly)">导出明细
          </el-button>
          <!--导出对账单-->
          <el-button v-show="inactive" size="small" type="success"
                     style="margin-bottom: 10px" icon="el-icon-download"
                     @click="exportStatement(v)">导出对账单
          </el-button>
          <!--导出入库汇总表-->
          <el-button v-show="inactive" size="small" type="success"
                     style="margin-bottom: 10px" icon="el-icon-download"
                     @click="exportInSummary(v)">导出入库汇总
          </el-button>
          <el-button v-show="inactive" size="small" type="success"
                         style="margin-bottom: 10px" icon="el-icon-download"
                         @click="exportSettlement(v)">导出结算台账
              </el-button>
          <!--导出财务审批-->
          <el-button v-show="inactive" size="small" type="success"
                     style="margin-bottom: 10px" icon="el-icon-download"
                     @click="exportFinApproval(v)">导出账务审批
          </el-button>
          <!--数据展示-->
          <vxe-toolbar export refresh custom>
          </vxe-toolbar>
          <vxe-table :border="true" :stripe="true" size="mini" v-resize:debounce.50="resizeTable"
                     :loading="v.loading" :export-config="tableExport"
                     :data="v.report.filter(stationsFilter)" show-footer ref="vendorReportTable"
                     show-overflow="tooltip" show-header-overflow show-footer-overflow highlight-hover-row
                     :footer-method="reportSummary" :highlight-current-row="false"
                     :max-height="tableMaxHeight + 'px'">
            <vxe-table-column fixed="left" width="60" type="seq" align="center"/>
            <vxe-table-column sortable width="260" title="供应商名称" field="vendorName" fixed="left">
              <template slot-scope="r">
                <el-tooltip placement="left" content="查看详细" :hide-after="500">
                  <el-button type="text" size="mini"
                             @click="viewDetailReport(v.month, r.row.vendorCode, r.row.vendorName,
                                            v.beginDate, v.endDate)">
                    {{ r.row.vendorName }}
                  </el-button>
                </el-tooltip>
              </template>
              <template slot="header" slot-scope="scope">
                <el-input
                    v-model="filterName"
                    autosize
                    size="mini"
                    clearable
                    placeholder="搜索供应商名称"/>
              </template>
            </vxe-table-column>
            <vxe-table-column sortable width="135" title="供应商编码" field="vendorCode">
              <!--emplate slot-scope="r">
                  <el-tooltip placement="left" content="查看详细" :hide-after="500">
                      <el-button type="text" size="mini"
                          @click="viewDetailReport(v.month, r.row.vendorCode, r.row.vendorName,
                              v.beginDate, v.endDate)">
                          {{ r.row.vendorCode }}
                      </el-button>
                  </el-tooltip>
              </template-->
              <template slot="header" slot-scope="scope">
                <el-input
                    v-model="filterNumber"
                    size="mini"
                    clearable
                    placeholder="搜索编码"/>
              </template>
            </vxe-table-column>

            <!--vxe-table-column sortable width="120" title="上月挂账" field="lastMonthUnpaid">
                <template slot-scope="p">
                    {{ beautyNum(p.row.lastMonthUnpaid) }}
                </template>
            </vxe-table-column-->
            <vxe-table-column sortable width="120" title="前期挂账" field="lastMonthUnpaid">
              <template slot-scope="p">
                {{ beautyNum(p.row.lastMonthUnpaid) }}
              </template>
            </vxe-table-column>
            <vxe-table-column sortable width="105" title="本月采购(原始)(已含税)" field="monthPurchaseOrig" :show-header-overflow="false">
              <template slot-scope="p">
                {{ beautyNum(p.row.monthPurchaseOrig) }}
              </template>
            </vxe-table-column>
            <vxe-table-column sortable width="130" title="本月采购(核对)" field="monthPurchase">
              <template slot-scope="p">
                <el-button type="text" size="mini"
                           @click="showOverWrite(4, p.row.vendorCode, p.row.vendorName, p.row.monthPurchase)">
                  {{ beautyNum(p.row.monthPurchase) }}
                </el-button>
              </template>
            </vxe-table-column>
            <vxe-table-column sortable width="130" title="本月开票(核对)" field="monthReceipt">
              <template slot-scope="p">
                <el-button type="text" size="mini"
                           @click="showOverWrite(1, p.row.vendorCode, p.row.vendorName, p.row.monthReceipt)">
                  {{ beautyNum(p.row.monthReceipt) }}
                </el-button>
              </template>
            </vxe-table-column>
            <vxe-table-column sortable width="130" title="本月返点(原始)" field="monthRebateOrig">
              <template slot-scope="p">
                {{ beautyNum(p.row.monthRebateOrig) }}
              </template>
            </vxe-table-column>
            <vxe-table-column sortable width="130" title="本月返点(核对)" field="monthRebate">
              <template slot-scope="p">
                <el-button type="text" size="mini"
                           @click="showOverWrite(4, p.row.vendorCode, p.row.vendorName, p.row.monthRebate)">
                  {{ beautyNum(p.row.monthRebate) }}
                </el-button>
              </template>
            </vxe-table-column>
<!--            <vxe-table-column sortable width="170" title="本月票面金额" field="monthPurFinal">-->
<!--              <template slot-scope="p">-->
<!--                {{ beautyNum(p.row.monthPurFinal) }}-->
<!--              </template>-->
<!--            </vxe-table-column>-->
            <vxe-table-column sortable width="105" title="本月销售(原始)(已含税)" field="monthSalesOrig" :show-header-overflow="false">
              <template slot-scope="p">
                {{ beautyNum(p.row.monthSalesOrig) }}
              </template>
            </vxe-table-column>
            <vxe-table-column sortable width="130" title="本月销售(核对)" field="monthSales">
              <template slot-scope="p">
                <el-button type="text" size="mini"
                           @click="showOverWrite(2, p.row.vendorCode, p.row.vendorName, p.row.monthSales)">
                  {{ beautyNum(p.row.monthSales) }}
                </el-button>
              </template>
            </vxe-table-column>
            <vxe-table-column sortable width="120" title="本月销售数量" field="monthSalesQty">
              <template slot-scope="p">
                {{ beautyNum(p.row.monthSalesQty) }}
              </template>
            </vxe-table-column>
            <!--vxe-table-column sortable width="120" title="上月结余" field="lastMonthLeft">
                <template slot-scope="p">
                    {{ beautyNum(p.row.lastMonthLeft) }}
                </template>
            </vxe-table-column-->
            <vxe-table-column sortable width="130" title="本月结算(原始)" field="monthPaidOrig">
              <template slot-scope="p">
                {{ beautyNum(p.row.monthPaidOrig) }}
              </template>
            </vxe-table-column>
            <vxe-table-column sortable width="130" title="本月结算(核对)" field="monthPaid">
              <template slot-scope="p">
                <el-button type="text" size="mini"
                           @click="showOverWrite(0, p.row.vendorCode, p.row.vendorName, p.row.monthPaid)">
                  {{ beautyNum(p.row.monthPaid) }}
                </el-button>
                <!--el-input v-else size="mini" style="width: 100%"
                    :placeholder="p.row.monthPaid" clearable
                    v-model="p.row.monthPayEdit">
                    <el-button slot="append" icon="el-icon-success" size="mini"
                        style="color: green; width: 100%"
                        @click="overwriteMonthpay(p.row.vendorCode, p.row.monthPayEdit, p.row)"/>
                    <el-button slot="prepend" icon="el-icon-error" size="mini"
                        style="color: red; width: 100%"
                        @click="$set(p.row, 'editPay', false)"/>
                </el-input-->
              </template>
            </vxe-table-column>
            <vxe-table-column sortable width="120" title="本月挂账" field="monthUnpaid">
              <template slot-scope="p">
                {{ beautyNum(p.row.monthUnpaid) }}
              </template>
            </vxe-table-column>
            <vxe-table-column sortable width="130" title="累计挂账(原始)" field="yearUnpaidOrig">
              <template slot-scope="p">
                {{ beautyNum(p.row.yearUnpaidOrig) }}
              </template>
            </vxe-table-column>
            <vxe-table-column sortable width="130" title="累计挂账(核对)" field="yearUnpaid">
              <template slot-scope="p">
                <el-button type="text" size="mini"
                           @click="showOverWrite(3, p.row.vendorCode, p.row.vendorName, p.row.yearUnpaid)">
                  {{ beautyNum(p.row.yearUnpaid) }}
                </el-button>
              </template>
            </vxe-table-column>
            <!--vxe-table-column sortable width="120" title="本月挂账" field="monthUnpaid">
                <template slot-scope="p">
                    <el-button type="text" size="mini"
                        @click="showOverWrite(3, p.row.vendorCode, p.row.vendorName, p.row.monthUnpaid)">
                        {{ beautyNum(p.row.monthUnpaid) }}
                    </el-button>
                </template>
            </vxe-table-column-->
<!--            <vxe-table-column sortable width="120" title="本月结余" field="monthSurplus">-->
<!--              <template slot-scope="p">-->
<!--                {{ beautyNum(p.row.monthSurplus) }}-->
<!--              </template>-->
<!--            </vxe-table-column>-->
            <vxe-table-column sortable width="120" title="本年采购" field="yearPurchase">
              <template slot-scope="p">
                {{ beautyNum(p.row.yearPurchase) }}
              </template>
            </vxe-table-column>
            <vxe-table-column sortable width="120" title="本年销售" field="yearSales">
              <template slot-scope="p">
                {{ beautyNum(p.row.yearSales) }}
              </template>
            </vxe-table-column>
            <vxe-table-column sortable width="120" title="本年销售数量" field="yearSalesQty">
              <template slot-scope="p">
                {{ beautyNum(p.row.yearSalesQty) }}
              </template>
            </vxe-table-column>
            <vxe-table-column sortable width="120" title="本年结算" field="yearPaid">
              <template slot-scope="p">
                {{ beautyNum(p.row.yearPaid) }}
              </template>
            </vxe-table-column>
<!--            <vxe-table-column sortable width="120" title="本年结余" field="yearSurplus">-->
<!--              <template slot-scope="p">-->
<!--                {{ beautyNum(p.row.yearSurplus) }}-->
<!--              </template>-->
<!--            </vxe-table-column>-->
            <!--vxe-table-column sortable width="120" title="年初库存" field="yearOpenStock">
              <template slot-scope="p">
                {{ beautyNum(p.row.yearOpenStock) }}
              </template>
            </vxe-table-column>
            <vxe-table-column sortable width="120" title="期初库存" field="monthOpenStock">
              <template slot-scope="p">
                {{ beautyNum(p.row.monthOpenStock) }}
              </template>
            </vxe-table-column>
            <vxe-table-column sortable width="120" title="期末库存" field="monthCloseStock">
              <template slot-scope="p">
                {{ beautyNum(p.row.monthCloseStock) }}
              </template>
            </vxe-table-column>
            <vxe-table-colgroup title="平均库存" align="center">
              <vxe-table-column sortable width="120" title="本月" field="monthAvgStock">
                <template slot-scope="p">
                  {{ beautyNum(p.row.monthAvgStock) }}
                </template>
              </vxe-table-column>
              <vxe-table-column sortable width="120" title="本年" field="yearAvgStock">
                <template slot-scope="p">
                  {{ beautyNum(p.row.yearAvgStock) }}
                </template>
              </vxe-table-column>
            </vxe-table-colgroup-->
            <vxe-table-colgroup title="周转天数" align="center">
              <vxe-table-column sortable width="120" title="本月" field="monthDaysTo">
                <template slot-scope="p">
                  {{ beautyNum(p.row.monthDaysTo) }}
                </template>
              </vxe-table-column>
              <vxe-table-column sortable width="120" title="本年" field="yearDayTo">
                <template slot-scope="p">
                  {{ beautyNum(p.row.yearDayTo) }}
                </template>
              </vxe-table-column>
            </vxe-table-colgroup>
            <vxe-table-column sortable width="120" title="品效" field="commEff">
              <template slot-scope="p">
                {{ beautyNum(p.row.commEff) }}
              </template>
            </vxe-table-column>
            <vxe-table-column sortable width="120" title="毛利率" field="grossMargin">
              <template slot-scope="scope">
                {{ scope.row.grossMargin }}%
              </template>
            </vxe-table-column>
            <vxe-table-column sortable width="120" title="到货率" field="orderFillRate">
              <template slot-scope="scope">
                <span>{{ scope.row.orderFillRate }}%</span>
              </template>
            </vxe-table-column>
            <vxe-table-column sortable width="200" title="新供应商初次供货时间" field="firstOrder" />
          </vxe-table>
        </div>
<!--        //////////-->
        <div v-else class="report-detail-section">
          <el-tag type="success">{{ v.vendorName }}</el-tag>
          <el-tag style="margin-bottom: 10px; margin-left: 10px" type="primary">{{ v.beginDate + ' ~ ' + v.endDate }}</el-tag>
          <x-table :minus="180" :border="true" :stripe="true" size="mini"
                   :loading="v.loading"
                   :footer-method="detailSummary" :data="v.report" show-footer
                   show-overflow="tooltip" show-header-overflow show-footer-overflow highlight-hover-row
                   :highlight-current-row="false">
            <vxe-table-column type="seq" width="60" />
            <vxe-table-column title="商品编码" field="material" width="200" sortable />
            <vxe-table-column title="商品名称" field="materialTxt" width="260" sortable />
            <vxe-table-column title="采购(含税)" field="buyinWtax" width="120" sortable>
              <template slot-scope="p">
                {{ beautyNum(p.row.buyinWtax) }}
              </template>
            </vxe-table-column>
            <vxe-table-column title="退货(含税)" field="returnWtax" width="120" sortable>
              <template slot-scope="p">
                {{ beautyNum(p.row.returnWtax) }}
              </template>
            </vxe-table-column>
            <vxe-table-column title="支出(含税)" field="costWtax" width="120" sortable>
              <template slot-scope="p">
                {{ beautyNum(p.row.costWtax) }}
              </template>
            </vxe-table-column>
            <!--vxe-table-column title="期初库存" field="openZinvCost" width="120" sortable>
              <template slot-scope="p">
                {{ beautyNum(p.row.openZinvCost) }}
              </template>
            </vxe-table-column>
            <vxe-table-column title="期末库存" field="closeZinvCost" width="120" sortable>
              <template slot-scope="p">
                {{ beautyNum(p.row.closeZinvCost) }}
              </template>
            </vxe-table-column-->
            <vxe-table-column title="返点" field="rebate" width="100" sortable>
              <template slot-scope="p">
                {{ beautyNum(p.row.rebate) }}%
              </template>
            </vxe-table-column>
          </x-table>
        </div>
      </el-tab-pane>
    </el-tabs>
  </div>
</template>

<style lang="scss">
.vendor-config-checkbox-frame {
  background-color: #D9ECFF;
  border-radius: 7px;
  padding: 4px 12px 4px 12px;
  margin-right: 20px;
}

.vendor-monthly-report {
  background-color: white;
  max-height: calc(100vh - 10px);
  overflow: auto;
  padding: 8px 8px 8px 8px;
  width: 100%;
  position: relative;

  .create-report-dialog {
    p.title {
      font-size: 15px;
      font-weight: bold;
    }
  }

  .report-list-section {
    height: calc(100vh - 100px);

    .el-card__header {
      padding: 10px 10px;
    }

    p.section-header {
      font-size: 15px;
      font-weight: bold;
    }
  }
.haha{
  //height: calc(20vh - 20px);
}
  .report-detail-section {

    p.section-header {
      font-size: 15px;
      font-weight: bold;
    }
  }

  .el-input-group__append {
    padding: 0 10px;
  }

  .el-input-group__prepend {
    padding: 0 10px;
  }
}
</style>

<script>
import { doRequest, downloadFile, beautifyNumber, doRequestv2, deepCopy, message, confirm } from '../../utils/utils'
import { getUserInfo } from '../../utils/dataStorage'
import moment from 'moment'
import SmartTable from '../mixins/SmartMaxHeightTable'
import XTable from '../mixins/SmartVxeTable'
import { getToken } from '../../utils/dataStorage.js'
import { Loading } from 'element-ui'
import Config from '../../config/index'
import fileDownload from 'js-file-download';
import XEUtils from 'xe-utils'
import resize from 'vue-resize-directive'

export default {
  directives: {
    resize
  },
  components: {
    'smart-table': SmartTable,
    'x-table': XTable
  },
  mounted() {
    this.getReportList();
    this.getVendorConfigList();
    this.modifyUser = getUserInfo().name
  },
  data() {
    return {
      step:1,
      tableMaxHeight: 500,
      modifyUser: '',
      modifyDialog: false,
      modifyOps: ['结算金额', '开票金额', '销售金额', '累计挂账', '采购金额'],
      modifyReason: '',
      modifyVendor: '',
      modifyVendorName: '',
      modifyNewValue: 0,
      modify: {
        type: 0,
        old: 0,
      },
      modifyHistory: [],
      createReportDialog: false,
      createReportDialogLoading: false,
      vendorConfigDlg: false,
      monthlyReportList: [],
      createReportForm: {
        month: null,
        dateRange: [],
        excludeVendors: []
      },
      vendors: [],
      selectedReport: [],
      reportListLoading: false,
      activeTab: 'monthly-report-config',
      reportTabs: [],

      vendorConfigList: [],
      configTableLoading: false,
      vendorConfigAddDlg: false,
      updatingVendorConfig: false,
      vendorConfig: {
        vendorCode: '',
        vendorName: '',
        lastMonthUnpaid: 0.0,
        lastMonthLeft: 0.0,
        yearPurchase: 0.0,
        yearSales: 0.0,
        yearPaid: 0.0,
        yearUnpaid: 0.0,
        rebate: 0.0,
        inactive: false,
        isNew: false,
      },

      inactive : true,
      tableExport: {
        // 默认选中类型
        type: 'xlsx',
        // 自定义类型
        types: ['xlsx', 'csv', 'html', 'xml', 'txt']
      },
      filterName: '',
      filterNumber: ''
    }
  },
  methods: {
    tabClick(tab) {
      if(tab.name.indexOf('true') != -1){
        this.inactive = false
      }
      if(tab.name.indexOf('false') != -1) {
        this.inactive = true
      }
    },
    createReportDialogClosed() {
      this.createReportForm.month = null;
      this.createReportForm.dateRange = [],
          this.createReportForm.excludeVendors = []
    },
    getReportList() {
      this.reportListLoading = true;
      doRequest({
        url: '/v1/web/erp/vendors/report/list',
        method: 'get'
      }, {
        obj: this,
        src: 'vendorsReportList',
        dst: 'monthlyReportList',
        success: res => {
          this.reportListLoading = false;
        },
        fail: err => {
          this.$message({
            type: 'error',
            message: '获取报表列表失败'
          })
          this.reportListLoading = false;
        }
      })
    },
    createMonthlyReport() {
      if (!this.createReportForm.month || this.createReportForm.dateRange.length < 2) {
        this.$message({
          type: 'warning',
          message: '请输入正确的时间信息'
        })
        return;
      }

      this.createReportDialogLoading = true;
      doRequest({
        url: '/v1/web/erp/vendors/report/add',
        method: 'post',
        data: {
          month: moment(this.createReportForm.month).format('YYYYMM'),
          beginDate: moment(this.createReportForm.dateRange[0]).format('YYYY-MM-DD'),
          endDate: moment(this.createReportForm.dateRange[1]).format('YYYY-MM-DD'),
        }
      }, {
        success: res => {
          this.$message({
            type: 'success',
            message: '创建月度报表成功'
          });
          this.getReportList();
          this.createReportDialog = false;
          this.createReportDialogLoading = false;
        },
        fail: err => {
          this.$message({
            type: 'error',
            message: `创建月度报表失败，请稍后再试`
          });
          this.createReportDialog = false;
          this.createReportDialogLoading = false;
        }
      })
    },
    viewDetailReport(month, vendor, vendorName, bdate, edate) {
      let tabName = `${month}-${vendor}`;

      for (let i = 0; i < this.reportTabs.length; i ++) {
        if (this.reportTabs[i].name == tabName) {
          this.activeTab = tabName;
          return;
        }
      }

      doRequest({
        url: `/v1/web/erp/materials/report/detail/${month}/${vendor}`,
        method: 'GET',
        loading: true,
      }, {
        success: res => {
          this.reportTabs.push({
            month: month,
            name: tabName,
            key: tabName,
            vendorName: vendorName,
            label: `${vendor} 明细`,
            report: res.materialReportList,
            beginDate: bdate,
            endDate: edate,
            loading: false,
            detail: true,
          });
          this.$nextTick(() => {
            this.activeTab = tabName;
          })
        },
        fail: err => {
          console.log(err);
          message('error', '获取商品明细失败，请稍后再试')
        }
      })
    },
    viewReport(b, c, force = false) {
      let tabName = `${b.month}-${c}`;
      this.inactive = !c
      for (let i = 0; i < this.reportTabs.length; i ++) {
        if (this.reportTabs[i].name == tabName) {
          this.activeTab = tabName;
          return;
        }
      }

      this.reportListLoading = true;
      doRequest({
        url: `/v1/web/erp/vendors/report/detail/${b.month}/${c ? 'active' : 'all'}`,
        method: 'get'
      }, {
        success: res => {
          this.reportTabs.push({
            month: b.month,
            name: tabName,
            key: tabName,
            label: `${b.month} 月度报表 (${c ? '仅活跃供应商' : '全供应商'})`,
            detailLabel: `${b.month} 商品明细 (${c ? '仅活跃供应商' : '全供应商'})`,
            statement: `${b.month} 对账单 (${c ? '仅活跃供应商' : '全供应商'})`,
            inSummary: `${b.month} 入库汇总表 (${c ? '仅活跃供应商' : '全供应商'})`,
            settlement:  `${b.month} 结算台账 (${c ? '仅活跃供应商' : '全供应商'})`,
            finApproval: `${b.month} 财务审批 (${c ? '仅活跃供应商' : '全供应商'})`,
            report: res.vendorsList,
            beginDate: b.beginDate,
            endDate: b.endDate,
            activeOnly: c,
            loading: false,
          });
          this.$nextTick(() => {
            this.activeTab = tabName;
          })
        },
        fail: err => {
          console.log(err);
          this.$message({
            type: 'error',
            message: '获取报表失败，请稍后再试'
          })
        },
        finally: _ => {
          this.reportListLoading = false;
        }
      })
    },
    deleteReport(b) {
      let self = this;
      this.$confirm(`确定要删除 ${b.month} 的月度报表么?`, '删除报表').then(
          () => {
            self.reportListLoading = true;
            doRequest({
              url: `/v1/web/erp/vendors/report/del/${b.month}`,
              method: 'get'
            }, {
              success: (res) => {
                this.reportListLoading = false;
                this.getReportList();
              },
              fail: err => {
                this.reportListLoading = false;
                console.log(err);
                this.$message({
                  type: 'error',
                  message: '删除报表失败，请稍后再试'
                })
              }
            })
          }
      ).catch(
          () => {

          }
      )
    },
    exportDetailReport(b, c) {
      b.loading = true;
      downloadFile({
        url: `/v1/web/erp/materials/report/export/${b.month}/${c ? 'active' : 'all'}`,
        method: 'GET',
        loading: true
      }, {
        fail: err => {
          console.log(err);
          this.$message({
            type: 'error',
            message: '导出明细报表失败，请稍后再试'
          });
        },
        finally: () => {
          b.loading = false;
        }
      })
    },
    exportReport(b, c, d) {
      if (!d)
        this.reportListLoading = true;
      else
        b.loading = true;
      downloadFile({
        url: '/v1/web/erp/vendors/report/export',
        method: 'post',
        loading: true,
        data: {
          month: b.month,
          beginDate: b.beginDate,
          endDate: b.endDate,
          activeOnly: c
        }
      }, {
        fail: err => {
          console.log(err);
          this.$message({
            type: 'error',
            message: '导出报表失败，请稍后再试'
          });
        },
        finally: () => {
          if (!d)
            this.reportListLoading = false;
          else
            b.loading = false;
        }
      })
    },
    exportStatement(b) {
      this.step=4;
      // v1/web/erp/vendors/report/account-statement/2020-04-25/2020-04-27
      b.loading = true;
      const loading = this.$loading({
        lock: true,
        text: 'Loading',
        spinner: 'el-icon-loading',
        background: 'rgba(0, 0, 0, 0.7)'
      });
      downloadFile({
        url: encodeURI(`/v1/web/erp/vendors/report/account-statement/${b.beginDate}/${b.endDate}`),
        method: 'GET',
        // loading: true
      }, {
        fail: err => {
          console.log(err);
          this.$message({
            type: 'error',
            message: '导出对账单失败，请稍后再试'
          });
        },
        finally: () => {
          b.loading = false;
          loading.close();
        }
      })
    },
    exportInSummary(b) {
      this.step=8;
      // alert("导出入库汇总表")
      // /v1/web/erp/vendors/report/endorsement-form/:begin-date/:end-date
      b.loading = true;
      const loading = this.$loading({
        lock: true,
        text: 'Loading',
        spinner: 'el-icon-loading',
        background: 'rgba(0, 0, 0, 0.7)'
      });
      downloadFile({
        url: `/v1/web/erp/vendors/report/endorsement-form/${b.beginDate}/${b.endDate}`,
        method: 'GET',
        // loading: true
      }, {
        fail: err => {
          console.log(err);
          this.$message({
            type: 'error',
            message: '导出入库汇总表失败，请稍后再试'
          });
        },
        finally: () => {
          b.loading = false;
          loading.close();
        }
      })
    },
    exportSettlement(b) {

      // this.step=9;
      // v1/web/erp/vendors/report/standing-book/:month
      b.loading = true;
      const loading = this.$loading({
        lock: true,
        text: 'Loading',
        spinner: 'el-icon-loading',
        background: 'rgba(0, 0, 0, 0.7)'
      });
      downloadFile({
        url: `/v1/web/erp/vendors/report/standing-book/${b.month}`,
        method: 'GET',
        // loading: true
      }, {
        fail: err => {
          console.log(err);
          this.$message({
            type: 'error',
            message: '导出结算台账失败，请稍后再试'
          });
        },
        finally: () => {
          b.loading = false;
          loading.close();
        }
      })
    },
    exportFinApproval(b) {
      // this.step=9;
      // alert("导出财务审批")
      // /v1/web/erp/vendors/report/payment-details/:month
      b.loading = true;
      const loading = this.$loading({
        lock: true,
        text: 'Loading',
        spinner: 'el-icon-loading',
        background: 'rgba(0, 0, 0, 0.7)'
      });
      downloadFile({
        url: `/v1/web/erp/vendors/report/payment-details/${b.month}`,
        method: 'GET',
        // loading: true
      }, {
        fail: err => {
          console.log(err);
          this.$message({
            type: 'error',
            message: '导出财务审批失败，请稍后再试'
          });
        },
        finally: () => {
          b.loading = false;
          loading.close();
        }
      })
    },
    removeTab(name) {
      for (let i = 0; i < this.reportTabs.length; i ++) {
        if (this.reportTabs[i].name == name) {

          this.reportTabs.splice(i, 1);
          if (this.reportTabs.length > 0)
            this.activeTab = this.reportTabs[this.reportTabs.length - 1].name;
          else
            this.activeTab = 'monthly-report-config';
          return;
        }
      }
    },
    detailSummary({ columns, data }) {
      const clen = columns.length;
      return [
        columns.map((column, columnIndex) => {
          if (columnIndex === 0) {
            return '合计'
          }

          if (columnIndex > 2 && columnIndex < clen - 1)
            return beautifyNumber(parseFloat(XEUtils.sum(data, column.property)).toFixed(2))
          else
            return ''
        })
      ]
    },
    reportSummary(param) {
      const { columns, data } = param;
      const clen = columns.length;
      return [
        columns.map((column, columnIndex) => {
          if (columnIndex === 0) {
            return '合计'
          }

          if (columnIndex > 2 && columnIndex < clen - 6)
            return beautifyNumber(parseFloat(XEUtils.sum(data, column.property)).toFixed(2))
          else
            return ''
        })
      ]


      // const sums = [];
      // const clen = columns.length;
      // columns.forEach((column, index) => {
      //     if (index === 0) {
      //         sums[index] = '合计';
      //         return;
      //     }

      //     if (index === 1 || index === 2 || index >= clen - 6) {
      //         sums[index] = '';
      //         return;
      //     }

      //     const values = data.map(item => Number(item[column.property]));
      //     if (!values.every(value => isNaN(value))) {
      //         sums[index] = values.reduce((prev, curr) => {
      //         const value = Number(curr);
      //         if (!isNaN(value)) {
      //             return prev + curr;
      //         } else {
      //             return prev;
      //         }
      //         }, 0).toFixed(2);
      //         sums[index] = beautifyNumber(sums[index]);
      //     }
      // });

      // return sums;
    },
    beautyNum(val) {
      return beautifyNumber(val);
    },
    getTokenVendor() {
      return getToken();
    },
    getVendorConfigList() {
      this.configTableLoading = true;
      doRequest({
        url: '/v1/web/erp/vendor/config/list',
        method: 'get'
      }, {
        obj: this,
        src: 'cfgs',
        dst: 'vendorConfigList',
        finally: _ => {
          this.configTableLoading = false;
        }
      })
    },
    showVendorConfigDlg(updating) {
      this.updatingVendorConfig = false;
      this.vendorConfigAddDlg = true;
    },
    showVendorConfigDlgUpdate(config) {
      this.updatingVendorConfig = true;
      this.vendorConfig = deepCopy(config);
      this.vendorConfigAddDlg = true;
    },
    addVendorConfig() {
      // if (this.vendorConfigName.length <= 0 || this.vendorConfigCode.length <= 0) {
      if (this.vendorConfigCode.length <= 0) {
        this.$message({
          type: 'warning',
          message: '请输入有效的供应商信息'
        });
        return;
      }

      doRequest({
        url: '/v1/web/erp/vendor/config/add',
        method: 'post',
        loading: true,
        data: this.vendorConfig
      }, {
        success: _ => {
          this.$message({
            type: 'success',
            message: '添加供应商配置成功'
          });
          this.getVendorConfigList();
          this.vendorConfigAddDlg = false;
        },
        fail: err => {
          console.log(err);
          this.$message({
            type: 'error',
            message: '添加供应商配置失败，请稍后再试'
          });
        }
      })
    },
    clearVendorConfig() {
      this.$confirm(`您确定要删除所有供应商配置么？`, '清空供应商配置', {
        type: 'warning'
      }).then(() => {
        doRequest({
          url: '/v1/web/erp/vendor/config/clear',
          method: 'get'
        }, {
          success: res => {
            this.$message({
              type: 'success',
              message: '清空供应商配置成功！'
            });
            this.getVendorConfigList();
          },
          fail: err => {
            console.log(err);
            this.$message({
              type: 'error',
              message: '清空供应商配置失败，请稍后再试！'
            })
          }
        })
      }).catch(() => {

      })
    },
    updateVendorConfig() {
      if (this.vendorConfig.vendorCode.length <= 0) {
        this.$message({
          type: 'warning',
          message: '请输入有效的供应商信息'
        });
        return;
      }

      doRequest({
        url: '/v1/web/erp/vendor/config/update',
        method: 'post',
        loading: true,
        data: this.vendorConfig,
      }, {
        success: _ => {
          this.$message({
            type: 'success',
            message: '更新供应商配置成功'
          });
          this.getVendorConfigList();
          this.vendorConfigAddDlg = false;
        },
        fail: err => {
          console.log(err);
          this.$message({
            type: 'error',
            message: '更新供应商配置失败，请稍后再试'
          });
        }
      })
    },
    deleteVendorConfig(vendorCode, vendorName) {
      this.$confirm(`您确定要删除供应商 ${vendorName} 的配置么？`, '删除供应商配置', {
        type: 'warning'
      }).then(() => {
        doRequest({
          url: `/v1/web/erp/vendor/config/del/${vendorCode}`,
          method: 'get'
        }, {
          success: res => {
            this.$message({
              type: 'success',
              message: '删除供应商配置成功！'
            });
            this.getVendorConfigList();
          },
          fail: err => {
            console.log(err);
            this.$message({
              type: 'error',
              message: '删除供应商配置失败，请稍后再试！'
            })
          }
        })
      }).catch(() => {

      })
    },
    downloadVendorConfigCSV() {
      doRequestv2({
        url: '/v1/web/erp/vendor/config/export',
        method: 'get',
        responseType: 'blob',
        loading: true
      }, {
        success: res => {
          let fn = res.headers['content-disposition'].substring('attachment;filename='.length);
          if (fn && res.data) {
            fileDownload(res.data, fn);
          }
        },
        fail: err => {
          console.log(err);
        }
      })
    },
    vendorConfigUploadBefore(file) {
      if (!this.localLoading) {
        this.localLoading = Loading.service({
          text: `${file.name} 上传中...`
        })
      }
    },
    vendorConfigUploadSuccess(res, file, fileList) {
      if (this.localLoading) {
        this.localLoading.close();
        this.localLoading = null;
      }

      this.$message({
        type: "success",
        message: `${file.name} 上传成功`
      });
      this.getVendorConfigList();
    },
    vendorConfigUploadError(err, file, fileList) {
      if (this.localLoading) {
        this.localLoading.close();
        this.localLoading = null;
      }

      this.$message({
        type: "error",
        message: `${file.name} 上传失败，请稍后再试`
      });
    },
    vendorConfigDlgClosed() {
      this.vendorConfig.vendorCode = '';
      this.vendorConfig.rebate = 0.0;
      this.vendorConfig.inactive = false;
      this.vendorConfig.isNew = false;
      this.vendorConfig.lastMonthUnpaid = 0.0;
      this.vendorConfig.lastMonthLeft = 0.0;
      this.vendorConfig.yearPurchase = 0.0;
      this.vendorConfig.yearSales = 0.0;
      this.vendorConfig.yearPaid = 0.0;
      this.vendorConfig.yearUnpaid = 0.0;
    },
    showMonthPay(row) {
      //row['monthPayEdit'] = `${row.monthPaid}`;
      this.$set(row, 'editPay', true);

    },
    showOverWrite(type, vendor, vendorName, old) {
      let tab = null;
      for (let i = 0; i < this.reportTabs.length;  i++) {
        if (this.activeTab == this.reportTabs[i].name) {
          tab = this.reportTabs[i];
          break;
        }
      }
      if (!tab) {
        message('error', '未找到合法的时间');
        return;
      }

      doRequest({
        method: 'GET',
        url: `/v1/web/erp/vendors/report/history/${type}/${tab.month}/${vendor}`,
        loading: true
      }, {
        success: res => {
          this.modifyHistory = res.history
        },
        fail: err => {
          console.log(err)
          message('error', '获取修改历史失败')
        },
        finally: _ => {
          this.modify.type = type
          this.modify.old = old
          this.modifyNewValue = old
          this.modifyVendor = vendor
          this.modifyVendorName = vendorName
          this.modifyDialog = true
        }
      })
    },
    doOverWrite(){
      if (this.modifyReason.length == 0) {
        message('warning', '修改原因不能为空')
        return
      }

      let tab = null;
      for (let i = 0; i < this.reportTabs.length;  i++) {
        if (this.activeTab == this.reportTabs[i].name) {
          tab = this.reportTabs[i];
          break;
        }
      }
      if (!tab) {
        message('error', '未找到合法的时间');
        return;
      }

      confirm('warning', '修改金额', `确认要修改 ${this.modifyOps[this.modify.type]} 值么?`, _ => {
        doRequest({
          method: 'POST',
          url: '/v1/web/erp/vendors/report/update',
          loading: true,
          data: {
            type: this.modify.type,
            month: tab.month,
            vendor: this.modifyVendor,
            user: this.modifyUser,
            reason: this.modifyReason,
            oldValue: this.modify.old,
            newValue: this.modifyNewValue
          }
        }, {
          success: res => {
            message('success', '更新结算金额成功');
            this.$set(tab, 'loading', true);
            doRequest({
              url: `/v1/web/erp/vendors/report/detail/${tab.month}/${tab.activeOnly ? 'active' : 'all'}`,
              method: 'GET',
              loading: true,
            }, {
              success: r => {
                tab.report = r.vendorsList
                this.modifyReason = ''
                this.modifyDialog = false
              },
              fail: err => {
                console.log(err);
                this.$message({
                  type: 'error',
                  message: '获取报表失败，请稍后再试'
                })
              },
              finally: _ => {
                this.$set(tab, 'loading', false);
              }
            })
          },
          fail: err => {
            console.log(err);
            message('error', '更新结算金额失败')
          }
        })
      })
    },
    overwriteMonthpay(vendor, val, row) {
      if (isNaN(val)) {
        message('error', `${val} 不是合法的数值`);
        return;
      }
      let tab = null;
      for (let i = 0; i < this.reportTabs.length;  i++) {
        if (this.activeTab == this.reportTabs[i].name) {
          tab = this.reportTabs[i];
          break;
        }
      }
      if (!tab) {
        message('error', '未找到合法的时间');
        return;
      }

      doRequest({
        url: '/v1/web/erp/vendors/report/update',
        method: 'POST',
        loading: true,
        data: {
          month: tab.month,
          vendorCode: vendor,
          monthPaid: parseFloat(val)
        }
      }, {
        success: res => {
          message('success', '更新结算金额成功');
          this.$set(tab, 'loading', true);
          doRequest({
            url: `/v1/web/erp/vendors/report/detail/${tab.month}/${tab.activeOnly ? 'active' : 'all'}`,
            method: 'GET',
            loading: true,
          }, {
            success: res => {
              tab.report = res.vendorsList
            },
            fail: err => {
              console.log(err);
              this.$message({
                type: 'error',
                message: '获取报表失败，请稍后再试'
              })
            },
            finally: _ => {
              this.$set(tab, 'loading', false);
            }
          })
        },
        fail: err => {
          console.log(err);
          message('error', '更新结算金额失败')
        }
      })
    },
    resizeTable() {
      this.tableMaxHeight = window.innerHeight - 350
    },
    stationsFilter(d) {
      if (this.filterName.length > 0) {
        let names = this.filterName.split(' ')
        let res = !d.vendorName.includes(names[0])
        for(let i=1; i<names.length; i++){
          res = res && (!d.vendorName.includes(names[i]))
        }
        if(res){
          return false;
        }
        // if (!d.vendorName.toLowerCase().includes(this.filterName.toLowerCase()))
        //     return false;
      }

      if (this.filterNumber.length > 0) {
        let nums = this.filterNumber.split(' ')
        let res = !d.vendorCode.includes(nums[0])
        for(let i=1; i<nums.length; i++){
          res = res && (!d.vendorCode.includes(nums[i]))
        }
        if(res){
          return false;
        }
        // if (!d.vendorCode.toLowerCase().includes(this.filterNumber.toLowerCase()))
        //     return false;
      }

      return true;
    }
  },
  computed: {
    uploadUrl: function() {
      return Config.apiUrl + '/v1/web/erp/vendor/config/import';
    }
  }
}
</script>
