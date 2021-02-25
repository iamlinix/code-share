<template>
  <div class="promo-main-v3" id="cnpc-promo-main-v3" @click="onMainDiv" v-resize:debounce.50="onRootResize">
    <el-dialog :visible.sync="dialogShow"
               custom-class="transparent-modal-dialog"
               :show-close="false">
    </el-dialog>
    <div class="promo-header-v3">
      <i v-if="fullScreen" class="ri-fullscreen-exit-line promo-full-screen"
         @click="exitFullscreen"></i>
      <i v-else class="ri-fullscreen-line promo-full-screen" @click="doFullscreen"></i>

      <i class="ri-arrow-right-s-fill promo-level-select-fa" @click="toggleShowPanel"></i>
      <!--el-button icon="el-icon-caret-right" class="promo-level-select" circle size="mini"
          @click="toggleShowPanel"/-->
      <font-awesome-icon id="toCheck" class="remove-unchecked" icon="square" v-on:click="getThresholdData"></font-awesome-icon>
      <font-awesome-icon id="checked" class="remove-checked" icon="check-square" v-on:click="getNormalData" style="display: none"></font-awesome-icon>

      <div v-if="showPanel" class="promo-level-panel">
        <el-select :popper-append-to-body="false" size="mini" style="width: 100%"
                   v-model="selectedOrg">
          <el-option v-for="(r, i) in orgs" :key="r.orgCode" :value="r.orgCode" :label="r.orgText" />
          <!--el-option value="000" label="北京公司" />
          <el-option value="A13" label="北京销售公司一分公司" />
          <el-option value="A14" label="北京销售公司二分公司" />
          <el-option value="A15" label="北京销售公司三分公司" />
          <el-option value="A16" label="北京销售公司四分公司" /-->
        </el-select>
        <el-date-picker size="mini" style="width: 100%; margin-top: 6px"
                        type="daterange" v-model="dateRange" range-separator="-"
                        start-placeholder="开始日期"
                        end-placeholder="结束日期"
                        value-format="yyyy-MM-dd" />
        <el-button size="mini" style="width: 100%; margin-top: 6px" type="primary"
                   @click="getNewData">查询</el-button>
      </div>


      <!--a>
      <el-select class="promo-level-select" size="mini" :popper-append-to-body="false"
          v-model="selectedOrg">
          <el-option value="000" label="总公司" />
          <el-option value="A13" label="北京销售公司一分公司" />
          <el-option value="A14" label="北京销售公司二分公司" />
          <el-option value="A15" label="北京销售公司三分公司" />
          <el-option value="A16" label="北京销售公司四分公司" />
      </el-select>

      <el-date-picker class="promo-level-date" size="mini"
          type="daterange" v-model="dateRange" range-separator="-"
          start-placeholder="开始日期"
          end-placeholder="结束日期"
          value-format="yyyy-MM-dd" />
      </a-->

      <p class="promo-title">{{ promotion.title }} </p>
      <p class="sub-title">{{ sdate }} - {{ edate }}</p>
      <span class="promo-timer">{{ currentTimeDesc }}</span>
    </div>
    <grid-layout class="promo-main-body" :layout.sync="gridLayout"
                 :use-css-transforms="true" :vertical-compact="true"
                 :is-draggable="false" :is-resizable="false" :row-height="gridRowHeight">
      <grid-item
          :x="gridLayout[0].x"
          :y="gridLayout[0].y"
          :w="gridLayout[0].w"
          :h="gridLayout[0].h"
          :i="gridLayout[0].i"
          :key="gridLayout[0].i"
          v-loading="gridLayout[0].loading">
        <div class="list-with-ranking">
          <p class="list-title">
            站级销售排行TOP10
          </p>
          <div v-if="topStationDetail" class="detail-display">
            <i class="el-icon-close close-button" @click="topStationDetail = false"/>
            <p class="detail-disp-title">{{ detailTopStation.OrgText }}</p>
            <el-form size="mini" label-position="right" label-width="70px">
              <el-form-item label="销售收入">
                <p>{{ detailTopStation.Metric ?
                    beautyNum(detailTopStation.Metric.NetvalInv.toFixed(2)) :
                    0 }}</p>
              </el-form-item>
              <el-form-item label="销售毛利">
                <p>{{ detailTopStation.Metric ?
                    beautyNum(detailTopStation.Metric.GrossProfit.toFixed(2)) :
                    0 }}</p>
              </el-form-item>
              <el-form-item label="毛利率">
                <p>{{ detailTopStation.Metric ?
                    beautyNum(detailTopStation.Metric.GrossMargin.toFixed(2)) :
                    0 }}%</p>
              </el-form-item>
            </el-form>
          </div>
          <table v-else>
            <tbody>
            <tr v-for="(v, i) in promotion.topStations" :key="'tr-top-key-' + i"
                @click="showListDetail(v, 'detailTopStation', 'topStationDetail')">
              <td>
                <font-awesome-icon style="font-size: 22px"
                                   v-if="i === 0" icon="medal" color="#FFD700">
                </font-awesome-icon>
                <font-awesome-icon style="font-size: 22px"
                                   v-else-if="i === 1" icon="medal" color="#C0C0C0">
                </font-awesome-icon>
                <font-awesome-icon style="font-size: 22px"
                                   v-else-if="i === 2" icon="medal" color="#B5A642">
                </font-awesome-icon>
                <font-awesome-icon v-else icon="medal"></font-awesome-icon>
              </td>
              <td>
                <tippy arrow append-to="parent">
                  <p class="tb-desc" slot="trigger">{{ v.OrgText }}</p>
                  {{ v.OrgText }}
                </tippy>
              </td>
              <td>
                <tippy arrow append-to="parent">
                  <p class="tb-desc" slot="trigger">{{ chnNumber(v.Metric.NetvalInv) }}</p>
                  {{ v.Metric.NetvalInv }}
                </tippy>
              </td>
            </tr>
            </tbody>
          </table>
        </div>
      </grid-item>
      <grid-item
          :x="gridLayout[1].x"
          :y="gridLayout[1].y"
          :w="gridLayout[1].w"
          :h="gridLayout[1].h"
          :i="gridLayout[1].i"
          :key="gridLayout[1].i"
          v-loading="gridLayout[1].loading">
        <div class="sale-amount-grid">
          <p class="sale-amount-title">销售收入</p>
          <p class="sale-amount-number">
            <font-awesome-icon icon="yen-sign"></font-awesome-icon>
            {{ beautyNum(promotion.sales.value.toFixed(0)) }}
          </p>
        </div>
      </grid-item>
      <grid-item
          :x="gridLayout[2].x"
          :y="gridLayout[2].y"
          :w="gridLayout[2].w"
          :h="gridLayout[2].h"
          :i="gridLayout[2].i"
          :key="gridLayout[2].i"
          v-loading="gridLayout[2].loading">
        <div class="percent-grid" v-if="!isNaN(promotion.sales.yoy)">
          <p class="percent-title">
            同比
            <font-awesome-icon v-if="promotion.sales.yoy > 0" icon="arrow-up" class="percent-arrow-up" />
            <font-awesome-icon v-else icon="arrow-down" class="percent-arrow-down" />
          </p>
          <p v-if="promotion.sales.yoy > 0" class="percent-number-up">
            {{ promotion.sales.yoy.toFixed(2) }}%
          </p>
          <p v-else class="percent-number-down">
            {{ promotion.sales.yoy.toFixed(2) }}%
          </p>
        </div>
      </grid-item>
      <grid-item
          :x="gridLayout[3].x"
          :y="gridLayout[3].y"
          :w="gridLayout[3].w"
          :h="gridLayout[3].h"
          :i="gridLayout[3].i"
          :key="gridLayout[3].i"
          v-loading="gridLayout[3].loading">
        <div class="percent-grid" v-if="!isNaN(promotion.sales.mom)">
          <p class="percent-title">
            环比
            <font-awesome-icon v-if="promotion.sales.mom > 0" icon="arrow-up" class="percent-arrow-up" />
            <font-awesome-icon v-else icon="arrow-down" class="percent-arrow-down" />
          </p>
          <p v-if="promotion.sales.mom > 0" class="percent-number-up">
            {{ promotion.sales.mom.toFixed(2) }}%
          </p>
          <p v-else class="percent-number-down">
            {{ promotion.sales.mom.toFixed(2) }}%
          </p>
        </div>
      </grid-item>
      <grid-item
          :x="gridLayout[4].x"
          :y="gridLayout[4].y"
          :w="gridLayout[4].w"
          :h="gridLayout[4].h"
          :i="gridLayout[4].i"
          :key="gridLayout[4].i"
          v-loading="gridLayout[4].loading">
        <div class="list-with-ranking">
          <p class="list-title">
            站级销售排行LAST10
          </p>
          <div v-if="ropStationDetail" class="detail-display">
            <i class="el-icon-close close-button" @click="ropStationDetail = false"/>
            <p class="detail-disp-title">{{ detailRopStation.OrgText }}</p>
            <el-form size="mini" label-position="right" label-width="70px">
              <el-form-item label="销售收入">
                <p>{{ detailRopStation.Metric ?
                    beautyNum(detailRopStation.Metric.NetvalInv.toFixed(2)) :
                    0 }}</p>
              </el-form-item>
              <el-form-item label="销售毛利">
                <p>{{ detailRopStation.Metric ?
                    beautyNum(detailRopStation.Metric.GrossProfit.toFixed(2)) :
                    0 }}</p>
              </el-form-item>
              <el-form-item label="毛利率">
                <p>{{ detailRopStation.Metric ?
                    beautyNum(detailRopStation.Metric.GrossMargin.toFixed(2)) :
                    0 }}%</p>
              </el-form-item>
            </el-form>
          </div>
          <table v-else>
            <tbody>
            <tr v-for="(v, i) in promotion.ropStations" :key="'tr-rop-key-' + i"
                @click="showListDetail(v, 'detailRopStation', 'ropStationDetail')">
              <td>
                <font-awesome-icon style="font-size: 22px"
                                   v-if="i === 0" icon="frown" color="#FFD700">
                </font-awesome-icon>
                <font-awesome-icon style="font-size: 22px"
                                   v-else-if="i === 1" icon="frown" color="#C0C0C0">
                </font-awesome-icon>
                <font-awesome-icon style="font-size: 22px"
                                   v-else-if="i === 2" icon="frown" color="#B5A642">
                </font-awesome-icon>
                <font-awesome-icon v-else icon="frown"></font-awesome-icon>
              </td>
              <td>
                <tippy arrow append-to="parent">
                  <p class="tb-desc" slot="trigger">{{ v.OrgText }}</p>
                  {{ v.OrgText }}
                </tippy>
              </td>
              <td>
                <tippy arrow append-to="parent">
                  <p class="tb-desc" slot="trigger">{{ chnNumber(v.Metric.NetvalInv) }}</p>
                  {{ v.Metric.NetvalInv }}
                </tippy>
              </td>
            </tr>
            </tbody>
          </table>
        </div>
      </grid-item>
      <grid-item
          :x="gridLayout[5].x"
          :y="gridLayout[5].y"
          :w="gridLayout[5].w"
          :h="gridLayout[5].h"
          :i="gridLayout[5].i"
          :key="gridLayout[5].i"
          v-loading="gridLayout[5].loading">
        <!--v-chart ref="cnpc-promotions-map" theme="macarons" :options="mapOptions" :autoresize="true"/-->
        <div class="list-with-ranking">
          <p class="list-title">销售目标进度</p>
        </div>
        <el-row :gutter="50">
          <el-col :span="12">
            <el-row :gutter="10" id="progress-grid">
              <el-col :span="4" style="margin-bottom: 40px; margin-top: 30px">
                <!--p class="general-title">月度销售目标</p-->
                <el-button circle type="primary">年</el-button>
              </el-col>
              <el-col :span="20" style="margin-bottom: 40px; margin-top: 30px">
                <el-tooltip :content="`${yearProgress.current} / ${yearProgress.target}`" placement="right" effect="light" :hide-after="1000">
                  <el-progress text-inside :percentage="yearProgress.percent" :stroke-width="40" :width="progressWidth"
                               :color="progressStatus[Math.ceil(yearProgress.percent / 25) - 1]" :format="percentFormat"/>
                </el-tooltip>
              </el-col>
              <el-col :span="4">
                <!--p class="general-title">季度销售目标</p-->
                <el-button circle type="primary">季</el-button>
              </el-col>
              <el-col :span="20" style="margin-bottom: 40px">
                <el-tooltip :content="`${quarterProgress.current} / ${quarterProgress.target}`" placement="right" effect="light" :hide-after="1000">
                  <el-progress text-inside :percentage="quarterProgress.percent" :stroke-width="40" :width="progressWidth"
                               :color="progressStatus[Math.ceil(quarterProgress.percent / 25) - 1]" :format="percentFormat"/>
                </el-tooltip>
              </el-col>
              <el-col :span="4">
                <!--p class="general-title">年度销售目标</p-->
                <el-button circle type="primary">月</el-button>
              </el-col>
              <el-col :span="20">
                <el-tooltip :content="`${monthProgress.current} / ${monthProgress.target}`" placement="right" effect="light"
                            :hide-after="1000">
                  <el-progress text-inside :percentage="monthProgress.percent" :stroke-width="40" :width="progressWidth"
                               :color="progressStatus[Math.ceil(monthProgress.percent / 25) - 1]" :format="percentFormat"/>
                </el-tooltip>
              </el-col>
            </el-row>
          </el-col>
          <el-col :span="12">
            <el-row :gutter="10" id="progress-grid">
              <el-col :span="4" style="margin-bottom: 40px; margin-top: 30px">
                <!--p class="general-title">月度销售目标</p-->
                <el-button circle type="primary">年</el-button>
              </el-col>
              <el-col :span="20" style="margin-bottom: 40px; margin-top: 30px">
                <el-tooltip :content="`${yearProfitProgress.current} / ${yearProfitProgress.target}`" placement="right" effect="light" :hide-after="1000">
                  <el-progress text-inside :percentage="yearProfitProgress.percent" :stroke-width="40" :width="progressWidth"
                               :color="progressStatus[Math.ceil(yearProfitProgress.percent / 25) - 1]" :format="percentFormat"/>
                </el-tooltip>
              </el-col>
              <el-col :span="4">
                <!--p class="general-title">季度销售目标</p-->
                <el-button circle type="primary">季</el-button>
              </el-col>
              <el-col :span="20" style="margin-bottom: 40px">
                <el-tooltip :content="`${quarterProfitProgress.current} / ${quarterProfitProgress.target}`" placement="right" effect="light" :hide-after="1000">
                  <el-progress text-inside :percentage="quarterProfitProgress.percent" :stroke-width="40" :width="progressWidth"
                               :color="progressStatus[Math.ceil(quarterProfitProgress.percent / 25) - 1]" :format="percentFormat"/>
                </el-tooltip>
              </el-col>
              <el-col :span="4">
                <!--p class="general-title">年度销售目标</p-->
                <el-button circle type="primary">月</el-button>
              </el-col>
              <el-col :span="20">
                <el-tooltip :content="`${monthProfitProgress.current} / ${monthProfitProgress.target}`" placement="right" effect="light"
                            :hide-after="1000">
                  <el-progress text-inside :percentage="monthProfitProgress.percent" :stroke-width="40" :width="progressWidth"
                               :color="progressStatus[Math.ceil(monthProfitProgress.percent / 25) - 1]" :format="percentFormat"/>
                </el-tooltip>
              </el-col>
            </el-row>
          </el-col>
        </el-row>
      </grid-item>
      <grid-item
          :x="gridLayout[6].x"
          :y="gridLayout[6].y"
          :w="gridLayout[6].w"
          :h="gridLayout[6].h"
          :i="gridLayout[6].i"
          :key="gridLayout[6].i"
          v-loading="gridLayout[6].loading">
        <div class="list-with-ranking">
          <p class="list-title">
            畅销商品TOP10
          </p>
          <div v-if="topGoodDetail" class="detail-display">
            <i class="el-icon-close close-button" @click="topGoodDetail = false"/>
            <p class="detail-disp-title">{{ detailTopGood.MaterialTxt }}</p>
            <el-form size="mini" label-position="right" label-width="70px">
              <el-form-item label="销售收入">
                <p>{{ detailTopGood.Metric ?
                    beautyNum(detailTopGood.Metric.NetvalInv.toFixed(2)) :
                    0 }}</p>
              </el-form-item>
              <el-form-item label="销售毛利">
                <p>{{ detailTopGood.Metric ?
                    beautyNum(detailTopGood.Metric.GrossProfit.toFixed(2)) :
                    0 }}</p>
              </el-form-item>
              <el-form-item label="毛利率">
                <p>{{ detailTopGood.Metric ?
                    beautyNum(detailTopGood.Metric.GrossMargin.toFixed(2)) :
                    0 }}%</p>
              </el-form-item>
            </el-form>
          </div>
          <table v-else>
            <tbody>
            <tr v-for="(v, i) in promotion.topGoods" :key="'tr-top-good-key-' + i"
                @click="showListDetail(v, 'detailTopGood', 'topGoodDetail')">
              <td>
                <font-awesome-icon style="font-size: 22px"
                                   v-if="i === 0" icon="medal" color="#FFD700">
                </font-awesome-icon>
                <font-awesome-icon style="font-size: 22px"
                                   v-else-if="i === 1" icon="medal" color="#C0C0C0">
                </font-awesome-icon>
                <font-awesome-icon style="font-size: 22px"
                                   v-else-if="i === 2" icon="medal" color="#B5A642">
                </font-awesome-icon>
                <font-awesome-icon v-else icon="medal"></font-awesome-icon>
              </td>
              <td>
                <tippy arrow append-to="parent">
                  <p class="tb-desc" slot="trigger">{{ v.MaterialTxt }}</p>
                  {{ v.MaterialTxt }}
                </tippy>
              </td>
              <td>
                <tippy arrow append-to="parent">
                  <p class="tb-desc" slot="trigger">{{ chnNumber(v.Metric.NetvalInv) }}</p>
                  {{ v.Metric.NetvalInv }}
                </tippy>
              </td>
            </tr>
            </tbody>
          </table>
        </div>
      </grid-item>
      <grid-item
          :x="gridLayout[7].x"
          :y="gridLayout[7].y"
          :w="gridLayout[7].w"
          :h="gridLayout[7].h"
          :i="gridLayout[7].i"
          :key="gridLayout[7].i"
          v-loading="gridLayout[7].loading">
        <div class="list-with-ranking">
          <p class="list-title">
            滞销商品TOP10
          </p>
          <div v-if="ropGoodDetail" class="detail-display">
            <i class="el-icon-close close-button" @click="ropGoodDetail = false"/>
            <p class="detail-disp-title">{{ detailRopGood.MaterialTxt }}</p>
            <el-form size="mini" label-position="right" label-width="70px">
              <el-form-item label="销售收入">
                <p>{{ detailRopGood.Metric ?
                    beautyNum(detailRopGood.Metric.NetvalInv.toFixed(2)) :
                    0 }}</p>
              </el-form-item>
              <el-form-item label="销售毛利">
                <p>{{ detailRopGood.Metric ?
                    beautyNum(detailRopGood.Metric.GrossProfit.toFixed(2)) :
                    0 }}</p>
              </el-form-item>
              <el-form-item label="毛利率">
                <p>{{ detailRopGood.Metric ?
                    beautyNum(detailRopGood.Metric.GrossMargin.toFixed(2)) :
                    0 }}%</p>
              </el-form-item>
            </el-form>
          </div>
          <table v-else>
            <tbody>
            <tr v-for="(v, i) in promotion.ropGoods" :key="'tr-rop-good-key-' + i"
                @click="showListDetail(v, 'detailRopGood', 'ropGoodDetail')">
              <td>
                <font-awesome-icon style="font-size: 22px"
                                   v-if="i === 0" icon="frown" color="#FFD700">
                </font-awesome-icon>
                <font-awesome-icon style="font-size: 22px"
                                   v-else-if="i === 1" icon="frown" color="#C0C0C0">
                </font-awesome-icon>
                <font-awesome-icon style="font-size: 22px"
                                   v-else-if="i === 2" icon="frown" color="#B5A642">
                </font-awesome-icon>
                <font-awesome-icon v-else icon="frown"></font-awesome-icon>
              </td>
              <td>
                <tippy arrow append-to="parent">
                  <p class="tb-desc" slot="trigger">{{ v.MaterialTxt }}</p>
                  {{ v.MaterialTxt }}
                </tippy>
              </td>
              <td>
                <tippy arrow append-to="parent">
                  <p class="tb-desc" slot="trigger">{{ chnNumber(v.Metric.NetvalInv) }}</p>
                  {{ v.Metric.NetvalInv }}
                </tippy>
              </td>
            </tr>
            </tbody>
          </table>
        </div>
      </grid-item>
      <grid-item
          :x="gridLayout[8].x"
          :y="gridLayout[8].y"
          :w="gridLayout[8].w"
          :h="gridLayout[8].h"
          :i="gridLayout[8].i"
          :key="gridLayout[8].i"
          v-loading="gridLayout[8].loading">
        <div class="sale-amount-grid">
          <p class="sale-amount-title">销售毛利</p>
          <p class="sale-amount-number">
            <font-awesome-icon icon="funnel-dollar"></font-awesome-icon>
            {{ beautyNum(promotion.profit.value.toFixed(0)) }}
          </p>
        </div>
        <!--v-chart theme="macarons" :options="transferOptions" :autoresize="true"/-->
      </grid-item>
      <grid-item
          :x="gridLayout[9].x"
          :y="gridLayout[9].y"
          :w="gridLayout[9].w"
          :h="gridLayout[9].h"
          :i="gridLayout[9].i"
          :key="gridLayout[9].i"
          v-loading="gridLayout[9].loading">
        <div class="percent-grid" v-if="!isNaN(promotion.profit.yoy)">
          <p class="percent-title">
            同比
            <font-awesome-icon v-if="promotion.profit.yoy > 0" icon="arrow-up" class="percent-arrow-up" />
            <font-awesome-icon v-else icon="arrow-down" class="percent-arrow-down" />
          </p>
          <p v-if="promotion.profit.yoy > 0" class="percent-number-up">
            {{ promotion.profit.yoy.toFixed(2) }}%
          </p>
          <p v-else class="percent-number-down">
            {{ promotion.profit.yoy.toFixed(2) }}%
          </p>
        </div>
      </grid-item>
      <grid-item
          :x="gridLayout[10].x"
          :y="gridLayout[10].y"
          :w="gridLayout[10].w"
          :h="gridLayout[10].h"
          :i="gridLayout[10].i"
          :key="gridLayout[10].i"
          v-loading="gridLayout[10].loading">
        <div class="percent-grid" v-if="!isNaN(promotion.profit.mom)">
          <p class="percent-title">
            环比
            <font-awesome-icon v-if="promotion.profit.mom > 0" icon="arrow-up" class="percent-arrow-up" />
            <font-awesome-icon v-else icon="arrow-down" class="percent-arrow-down" />
          </p>
          <p v-if="promotion.profit.mom > 0" class="percent-number-up">
            {{ promotion.profit.mom.toFixed(2) }}%
          </p>
          <p v-else class="percent-number-down">
            {{ promotion.profit.mom.toFixed(2) }}%
          </p>
        </div>
      </grid-item>
      <grid-item
          :x="gridLayout[11].x"
          :y="gridLayout[11].y"
          :w="gridLayout[11].w"
          :h="gridLayout[11].h"
          :i="gridLayout[11].i"
          :key="gridLayout[11].i"
          v-loading="gridLayout[11].loading">
        <div class="single-num-grid">
          <p class="single-num-title">
            非油消费人数
          </p>
          <p class="single-num-normal">
            <font-awesome-icon icon="user-friends" />
            {{ beautyNum(promotion.customers.value) }}
          </p>
        </div>
      </grid-item>
      <grid-item
          :x="gridLayout[12].x"
          :y="gridLayout[12].y"
          :w="gridLayout[12].w"
          :h="gridLayout[12].h"
          :i="gridLayout[12].i"
          :key="gridLayout[12].i"
          v-loading="gridLayout[12].loading">
        <div class="percent-grid" v-if="!isNaN(promotion.customers.yoy)">
          <p class="percent-title">
            同比
            <font-awesome-icon v-if="promotion.customers.yoy > 0" icon="arrow-up" class="percent-arrow-up" />
            <font-awesome-icon v-else icon="arrow-down" class="percent-arrow-down" />
          </p>
          <p v-if="promotion.customers.yoy > 0" class="percent-number-up">
            {{ promotion.customers.yoy.toFixed(2) }}%
          </p>
          <p v-else class="percent-number-down">
            {{ promotion.customers.yoy.toFixed(2) }}%
          </p>
        </div>
      </grid-item>
      <grid-item
          :x="gridLayout[13].x"
          :y="gridLayout[13].y"
          :w="gridLayout[13].w"
          :h="gridLayout[13].h"
          :i="gridLayout[13].i"
          :key="gridLayout[13].i"
          v-loading="gridLayout[13].loading">
        <div class="percent-grid" v-if="!isNaN(promotion.customers.mom)">
          <p class="percent-title">
            环比
            <font-awesome-icon v-if="promotion.customers.mom > 0" icon="arrow-up" class="percent-arrow-up" />
            <font-awesome-icon v-else icon="arrow-down" class="percent-arrow-down" />
          </p>
          <p v-if="promotion.customers.mom > 0" class="percent-number-up">
            {{ promotion.customers.mom.toFixed(2) }}%
          </p>
          <p v-else class="percent-number-down">
            {{ promotion.customers.mom.toFixed(2) }}%
          </p>
        </div>
      </grid-item>
      <grid-item
          :x="gridLayout[14].x"
          :y="gridLayout[14].y"
          :w="gridLayout[14].w"
          :h="gridLayout[14].h"
          :i="gridLayout[14].i"
          :key="gridLayout[14].i"
          v-loading="gridLayout[14].loading">
        <!--div class="single-num-grid">
               <p class="single-num-title">
                   促销商品贡献度
               </p>
               <div class="single-num-small">
                   <el-row :gutter="8">
                       <el-col :span="4">
                           <font-awesome-icon style="width: 100%" icon="money-bill" />
                       </el-col>
                       <el-col :span="20">
                           <tippy arrow append-to="parent" placement="right">
                        <el-progress slot="trigger" :text-inside="true" :stroke-width="26"
                            :percentage="promotion.contribute.gross"
                            :color="percentColors"/>
                        收入
                    </tippy>
                       </el-col>
                   </el-row>


               </div>
               <div class="single-num-small">
                   <el-row :gutter="8">
                       <el-col :span="4">
                           <font-awesome-icon style="width: 100%" icon="funnel-dollar" />
                       </el-col>
                       <el-col :span="20">
                           <tippy arrow append-to="parent" placement="right">
                        <el-progress slot="trigger" :text-inside="true" :stroke-width="26"
                            :percentage="promotion.contribute.profit"
                            :color="percentColors"/>
                        毛利
                    </tippy>
                       </el-col>
                   </el-row>

               </div>
           </div-->
        <div class="sale-amount-grid">
          <p class="sale-amount-title">客单价</p>
          <p class="sale-amount-normal">
            <font-awesome-icon icon="coins"></font-awesome-icon>
            {{ beautyNum(promotion.average.value.toFixed(2)) }}
          </p>
        </div>
      </grid-item>
      <grid-item
          :x="gridLayout[15].x"
          :y="gridLayout[15].y"
          :w="gridLayout[15].w"
          :h="gridLayout[15].h"
          :i="gridLayout[15].i"
          :key="gridLayout[15].i"
          v-loading="gridLayout[15].loading">
        <div class="percent-grid" v-if="!isNaN(promotion.average.yoy)">
          <p class="percent-title">
            同比
            <font-awesome-icon v-if="promotion.average.yoy > 0" icon="arrow-up" class="percent-arrow-up" />
            <font-awesome-icon v-else icon="arrow-down" class="percent-arrow-down" />
          </p>
          <p v-if="promotion.average.yoy > 0" class="percent-number-up">
            {{ promotion.average.yoy.toFixed(2) }}%
          </p>
          <p v-else class="percent-number-down">
            {{ promotion.average.yoy.toFixed(2) }}%
          </p>
        </div>
        <!--div class="single-num-grid">
            <p class="single-num-title">
                在售站点数
            </p>
            <p class="single-num-normal">
                <font-awesome-icon icon="gas-pump"/>
                {{ promotion.activeStationNumber }}
            </p>
        </div-->
      </grid-item>
      <grid-item
          :x="gridLayout[16].x"
          :y="gridLayout[16].y"
          :w="gridLayout[16].w"
          :h="gridLayout[16].h"
          :i="gridLayout[16].i"
          :key="gridLayout[16].i"
          v-loading="gridLayout[16].loading">
        <!--div class="single-num-grid">
            <p class="single-num-title">
                汽服站点数
            </p>
            <p class="single-num-normal">
                <font-awesome-icon icon="car"/>
                {{ promotion.stationWithCarService }}
            </p>
        </div-->
        <div class="percent-grid" v-if="!isNaN(promotion.average.mom)">
          <p class="percent-title">
            环比
            <font-awesome-icon v-if="promotion.average.mom > 0" icon="arrow-up" class="percent-arrow-up" />
            <font-awesome-icon v-else icon="arrow-down" class="percent-arrow-down" />
          </p>
          <p v-if="promotion.average.mom > 0" class="percent-number-up">
            {{ promotion.average.mom.toFixed(2) }}%
          </p>
          <p v-else class="percent-number-down">
            {{ promotion.average.mom.toFixed(2) }}%
          </p>
        </div>
      </grid-item>
      <grid-item
          :x="gridLayout[17].x"
          :y="gridLayout[17].y"
          :w="gridLayout[17].w"
          :h="gridLayout[17].h"
          :i="gridLayout[17].i"
          :key="gridLayout[17].i"
          v-loading="gridLayout[17].loading">
        <div id="donut-auto-resize-transfer" style="width: 100%; height: 100%"
             v-resize:debounce.50="onDonutResize">
          <vc-donut :size="donutWidth" unit="%" :thickness="30"
                    background="#052ACF"
                    :sections="[{ value: promotion.transfer.value, color: percentColors[~~(promotion.transfer.value/20)].color }]"
                    foreground="#C0C4CC">
            <p style="color: rgba(255, 255, 255, 1); font-size: 14px; margin-bottom: 2px">
              油非转换
            </p>
            <p :style="`color: ${percentColors[~~(promotion.transfer.value/20)].color}; font-size: 21px; font-weight: bold`">
              {{ promotion.transfer.value.toFixed(2) }}%
            <p/>
          </vc-donut>
        </div>
        <!--v-chart theme="macarons" :options="levelTreeOptions" :autoresize="true"/-->
      </grid-item>
      <grid-item
          :x="gridLayout[18].x"
          :y="gridLayout[18].y"
          :w="gridLayout[18].w"
          :h="gridLayout[18].h"
          :i="gridLayout[18].i"
          :key="gridLayout[18].i"
          v-loading="gridLayout[18].loading">
        <div class="percent-grid" v-if="!isNaN(promotion.transfer.yoy)">
          <p class="percent-title">
            同比
            <font-awesome-icon v-if="promotion.transfer.yoy > 0" icon="arrow-up" class="percent-arrow-up" />
            <font-awesome-icon v-else icon="arrow-down" class="percent-arrow-down" />
          </p>
          <p v-if="promotion.transfer.yoy > 0" class="percent-number-up">
            {{ promotion.transfer.yoy.toFixed(2) }}%
          </p>
          <p v-else class="percent-number-down">
            {{ promotion.transfer.yoy.toFixed(2) }}%
          </p>
        </div>
      </grid-item>
      <grid-item
          :x="gridLayout[19].x"
          :y="gridLayout[19].y"
          :w="gridLayout[19].w"
          :h="gridLayout[19].h"
          :i="gridLayout[19].i"
          :key="gridLayout[19].i"
          v-loading="gridLayout[19].loading">
        <div class="percent-grid" v-if="!isNaN(promotion.transfer.mom)">
          <p class="percent-title">
            环比
            <font-awesome-icon v-if="promotion.transfer.mom > 0" icon="arrow-up" class="percent-arrow-up" />
            <font-awesome-icon v-else icon="arrow-down" class="percent-arrow-down" />
          </p>
          <p v-if="promotion.transfer.mom > 0" class="percent-number-up">
            {{ promotion.transfer.mom.toFixed(2) }}%
          </p>
          <p v-else class="percent-number-down">
            {{ promotion.transfer.mom.toFixed(2) }}%
          </p>
        </div>
      </grid-item>
      <grid-item
          :x="gridLayout[20].x"
          :y="gridLayout[20].y"
          :w="gridLayout[20].w"
          :h="gridLayout[20].h"
          :i="gridLayout[20].i"
          :key="gridLayout[20].i"
          v-loading="gridLayout[20].loading">
        <div v-if="showPaySel" class="single-cat-panel">
          <el-button size="mini" style="width: 100%;" type="primary"
                     @click="updatePayType(0)">按金额</el-button>
          <el-button size="mini" style="width: 100%; margin: 6px 0 0 0" type="primary"
                     @click="updatePayType(1)">按数量</el-button>
        </div>
        <i class="ri-arrow-down-s-line cat-select" @click="showPaySel = !showPaySel"></i>
        <v-chart theme="macarons" :options="payPieOptions" :autoresize="true"/>
      </grid-item>
    </grid-layout>
  </div>
</template>

<style lang="scss">
.promo-main-v3 {
  background-image: url('../../assets/images/promo-bg.jpg');
  background-position-x: center;
  background-position-y: top;
  width: 100%;
  height: calc(100vh - 12px);
  overflow: auto;

  .cat-select {
    position: absolute;
    top: 4px;
    left: 4px;
    color: rgba(255, 255, 255, 0.1);
    font-size: 18px;
    cursor: pointer;
    z-index: 999;

    &:hover {
      color: rgba(255, 255, 255, 0.9);
    }
  }

  .single-cat-panel {
    position: absolute;
    width: 150PX;
    padding: 10px 10px 10px 10px;
    border-radius: 5px;
    background-color: #FFF;
    top: 6px;
    left: 25px;
    z-index: 1000;
  }

  .transparent-modal-dialog {
    background-color: rgba(255, 255, 255, 0.2);
  }

  .cdc-container {
    height: 100%;
  }

  .el-progress__text {
    color: #FFF
  }

  .echarts {
    width: 100%;
    height: 100%;
  }

  .el-progress-bar__outer {
    background-color: #C0C4CC;
  }

  .el-progress-bar__innerText {
    margin-bottom: 8px;
  }

  p.promo-title {
    color: #FFF;
    text-align: center;
    padding-top: 4px;
    font-size: 1.5em;
  }

  p.sub-title {
    color: #FFF;
    text-align: center;
    font-size: 12px;
  }

  span.promo-timer {
    color: rgba(255, 255, 255, 0.8);
    position: absolute;
    top: 15px;
    right: 20px;
  }

  p.general-title {
    color:rgba(255, 255, 255, 1);
    font-size: 16px;
    padding-top: 16px;
    padding-right: 10px;
    text-align: right;
  }

  div.promo-header-v3 {
    background-image: url('../../assets/images/promo-head-bg.png');
    background-color: rgba(0, 0, 0, 0);
    height: 60px;
    background-position: center;
    background-repeat: no-repeat;
    position: relative;

    .promo-full-screen {
      color: rgba(255, 255, 255, 0.1);
      position: absolute;
      top: 10px;
      left: 10px;
      font-size: 20px;
      cursor: pointer;

      &:hover {
        color: rgba(255, 255, 255, 0.9);
      }
    }

    .promo-level-select-fa {
      color: rgba(255, 255, 255, 0.1);
      position: absolute;
      top: 0px;
      left: 40px;
      font-size: 35px;
      cursor: pointer;

      &:hover {
        color: rgba(255, 255, 255, 0.9);
      }
    }

    .remove-unchecked {
      color: rgba(255, 255, 255, 0.1);
      position: absolute;
      top: 17px;
      left: 90px;
      font-size: 16px;
      cursor: pointer;

      &:hover {
        color: rgba(255, 255, 255, 0.9);
      }
    }

    .remove-checked {
      color: rgba(255, 255, 255, 0.9);
      position: absolute;
      top: 17px;
      left: 90px;
      font-size: 16px;
      cursor: pointer;
    }

    .promo-level-select {
      opacity: 0;
      top: 6px;
      left: 60px;
      position: absolute;

      &:hover {
        opacity: 1;
      }
    }

    .promo-level-panel {
      position: absolute;
      width: 250PX;
      padding: 10px 10px 10px 10px;
      border-radius: 5px;
      background-color: #FFF;
      top: 6px;
      left: 80px;
      z-index: 999;
    }

    .promo-level-date {
      opacity: 0;
      top: 8px;
      left: 160px;
      position: absolute;

      &:hover {
        opacity: 1;
      }
    }
  }

  .promo-main-body {
    padding: 6px 6px 6px 6px;

    .sale-amount-grid {
      text-align: center;

      .sale-amount-title {
        text-align: left;
        color:rgba(255, 255, 255, 1);
        font-size: 18px;
        padding-left: 12px;
      }

      .sale-amount-title::before {
        content: "";
        width: 5px;
        height: 20px;
        border-radius: 18px;
        background-color: #409EFF;
        display: block;
        position: absolute;
        top: 11px;
        left: 8px;
      }

      .sale-amount-number {
        font-family: 'Manrope';
        color: #C7CD4F;
        font-size: 20px;
        font-weight: bold;
        padding: 5px 0 0 0;
        text-align: center;
      }

      .sale-amount-normal {
        font-family: 'Manrope';
        color: #C7CD4F;
        font-size: 20px;
        font-weight: bold;
        padding: 5px 0 0 0;
        text-align: center;
      }
    }

    .percent-grid {
      .percent-title {
        color:rgba(255, 255, 255, 1);
        font-size: 13px;
        padding-left: 12px;
      }

      .percent-title::before {
        content: "";
        width: 5px;
        height: 15px;
        border-radius: 15px;
        background-color: #409EFF;
        display: block;
        position: absolute;
        top: 14px;
        left: 8px;
      }

      .percent-arrow-up {
        font-size: 22px;
        font-weight: bold;
        color: #67C23A;
      }

      .percent-arrow-down {
        font-size: 22px;
        font-weight: bold;
        color: #F56C6C;
      }

      .percent-number-up {
        font-family: 'Manrope';
        font-size: 22px;
        margin-top: 4px;
        font-weight: bold;
        color: #67C23A;
      }

      .percent-number-down {
        font-family: 'Manrope';
        font-size: 22px;
        margin-top: 4px;
        font-weight: bold;
        color: #F56C6C;
      }
    }

    .single-num-grid {
      .single-num-title {
        color:rgba(255, 255, 255, 1);
        font-size: 16px;
        padding-left: 12px;
      }

      .single-num-title::before {
        content: "";
        width: 5px;
        height: 18px;
        border-radius: 18px;
        background-color: #409EFF;
        display: block;
        position: absolute;
        top: 10px;
        left: 8px;
      }

      .single-num-gold {
        color: #C7CD4F;
        font-size: 20px;
        font-weight: bold;
        margin-top: 40px;
        text-align: center;
        font-family: 'Manrope';
      }

      .single-num-normal {
        color: #E6A23C;
        font-size: 30px;
        font-weight: bold;
        margin-top: 20px;
        text-align: center;
        font-family: 'Manrope';
      }

      .single-num-small {
        color: #E6A23C;
        font-size: 30px;
        font-weight: bold;
        margin-top: 20px;
        text-align: center;
        font-family: 'Manrope';
      }
    }

    .vue-grid-item {
      background-color: rgba(255, 255, 255, 0.02);
      border-radius: 4px;
      padding: 8px 8px 8px 8px;
    }

    .list-with-ranking {
      .list-title {
        font-size: 16px;
        color: rgba(255, 255, 255, 1);
        padding-left: 12px;
        margin-bottom: 8px;
      }

      .list-title::before {
        content: "";
        width: 5px;
        height: 18px;
        border-radius: 18px;
        background-color: #409EFF;
        display: block;
        position: absolute;
        top: 10px;
        left: 8px;
      }

      .detail-display {
        top: 15%;
        left: 5%;
        width: 90%;
        height: 80%;
        position: absolute;
        padding: 8px 8px 8px 8px;

        color: #FFF;
        border-radius: 5px;

        .close-button {
          font-size: 10px;
          position: absolute;
          right: 4px;
          cursor: pointer;
          color: #AAA;
        }

        p.detail-disp-title {
          font-size: 14px;
          font-weight: bold;
          text-align: center;
          margin-bottom: 12px;
          padding: 0 6px 0 6px;
        }

        form {
          position: absolute;
          width: 50%;
          left: 25%;
          p {
            font-weight: bold;
          }
        }
      }

      p.tb-desc {
        max-width: 100%;
        width: 100%;
        white-space: nowrap;
        overflow: hidden;
        text-overflow: ellipsis;
      }

      table {
        color: rgba(255, 255, 255, 1);
        font-size: 16px;
        display: table;
        vertical-align: bottom;
        width: 100%;
        table-layout: fixed;

        tr {

          height: 25px;
          cursor: pointer;

          td:nth-child(1) {
            width: 30px;
            display: inline-block;
            vertical-align: middle;
            text-align: center;
          }
          td:nth-child(2) {
            width: calc(100% - 120px);
            white-space: nowrap;
            overflow: hidden;
            text-overflow: ellipsis;
            display: inline-block;
            vertical-align: middle;
            font-size: 13px;
          }
          td:nth-child(3) {
            padding-left: 10px;
            width: 90px;
            white-space: nowrap;
            overflow: hidden;
            text-overflow: ellipsis;
            display: inline-block;
            vertical-align: middle;
          }
        }
      }
    }
  }

  .el-progress-bar__innerText {
    font-size: 20px;
  }
}
</style>

<script>
import { toggleKey, beautifyNumberChn, beautifyNumber, randomNumber,
  doRequest, allRequests, message } from '../../utils/utils'
import { getUserInfo } from '../../utils/dataStorage'
import moment from 'moment'
import VueGridLayout from 'vue-grid-layout';
import fitty from 'fitty'
import resize from 'vue-resize-directive'
import 'remixicon/fonts/remixicon.css'
import { addListener, removeListener } from 'resize-detector'
import Config from '../../config/index'

import ECharts from 'vue-echarts'
import 'echarts/lib/chart/line'
import 'echarts/lib/chart/bar'
import 'echarts/lib/chart/pie'
import 'echarts/lib/component/title'
import 'echarts/lib/component/tooltip'
import 'echarts/lib/component/legend'
import 'echarts/lib/component/toolbox'
import 'echarts/map/js/province/beijing'
import 'echarts/map/json/province/beijing'
import 'echarts/theme/macarons'
var self = null

export default {
  directives: {
    resize
  },
  components: {
    'grid-layout': VueGridLayout.GridLayout,
    'grid-item': VueGridLayout.GridItem,
    'v-chart': ECharts
  },
  beforeDestroy() {
    clearInterval(this.timerInterval);
  },
  mounted() {
    self = this
    let earlier = new Date();
    earlier.setDate(earlier.getDate() - 30);
    this.dateRange.push(moment(earlier).format('YYYY-MM-DD'));
    this.dateRange.push(moment(Date.now()).format('YYYY-MM-DD'));

    this.promoTimer();
    this.timerInterval = setInterval(this.promoTimer, 1000);
    window.addEventListener('keydown', this.onKeyDown);

    let ele = document.getElementById('donut-auto-resize-transfer');
    // addListener(ele, this.onDonutResize);
    this.onDonutResize(ele);

    let root = document.getElementById('cnpc-promo-main-v3');
    addListener(root, this.onRootResize);

    this.fitties.push(fitty('.sale-amount-number', {
      maxSize: 50
    }));
    this.fitties.push(fitty('.sale-amount-normal', {
      maxSize: 40
    }));
    this.fitties.push(fitty('.single-num-normal', {
      maxSize: 40
    }));

    // this.getAllData();
    doRequest({
      method: 'GET',
      url: `/v1/web/user/org-perm/${getUserInfo().name}`,
      loading: true,
    }, {
      success: res => {
        let perms = []
        if (res.orgPerm.length > 0) {
          perms = JSON.parse(res.orgPerm)
        }

        let orgs = [{
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
        }]

        if (perms.length > 0) {
          orgs.forEach(e => {
            perms.forEach(p => {
              if (e.orgCode == p.orgCode) {
                if (p.show || (p.show == null)) {
                  this.orgs.push(e)
                }
              }
            })
          })
        } else {
          this.orgs = orgs
        }
      },
      fail: _ => {
        message('error', '获取用户机构权限失败，请稍后再试')
      }
    })
  },
  data() {
    return {
      orgs: [],
      paySource: null,
      showPaySel: false,
      gridRowHeight: 72,
      yearProgress: {
        target: 0,
        current: 0,
        percent: 0,
      },
      monthProgress: {
        target: 0,
        current: 0,
        percent: 0,
      },
      quarterProgress: {
        target: 0,
        current: 0,
        percent: 0,
      },
      yearProfitProgress: {
        target: 0,
        current: 0,
        percent: 0,
      },
      monthProfitProgress: {
        target: 0,
        current: 0,
        percent: 0,
      },
      quarterProfitProgress: {
        target: 0,
        current: 0,
        percent: 0,
      },
      progressWidth: 200,
      progressStatus: [
        '#F56C6C',
        '#E6A23C',
        '#409EFF',
        '#67C23A',
      ],

      fitties: [],
      dialogShow: false,
      donutWidth: 60,
      ticker: 0,
      fullScreen: false,
      currentTimeDesc: '',

      topStationDetail: false,
      detailTopStation: {},
      ropStationDetail: false,
      detailRopStation: {},
      topGoodDetail: false,
      detailTopGood: {},
      ropGoodDetail: false,
      detailRopGood: {},

      timerInterval: null,

      selectedOrg: '',
      selectedOrgTxt: '',
      subCorps: [],
      dateRange: [],
      sdate: '',
      edate: '',
      dispOrgLevel: 0,
      dispOrgCode: '',
      showPanel: false,
      levelName: {
        '000': '北京公司',
        'A13': '北京销售公司一分公司',
        'A14': '北京销售公司二分公司',
        'A15': '北京销售公司三分公司',
        'A16': '北京销售公司四分公司',
      },

      percentColors: [
        {color: '#f56c6c', percentage: 20},
        {color: '#e6a23c', percentage: 40},
        {color: '#5cb87a', percentage: 60},
        {color: '#1989fa', percentage: 80},
        {color: '#6f7ad3', percentage: 100}
      ],
      promotion: {
        title: '北京公司销售数据',
        topStations: [],
        ropStations: [],
        topGoods: [],
        ropGoods: [],
        sales: {
          value: 0
        },
        transfer: {
          value: 0
        },
        average: {
          value: 0
        },
        contribute: {},
        profit: {
          value: 0
        },
        customers: {},
        activeStationNumber: 117,
        stationWithCarService: 72,
        stationLevels: [

        ]
      },
      gridLayout: [
        // top10 station
        {"x":0,"y":0,"w":3,"h":4,"i":"0", loading: false},

        // sale income
        {"x":3,"y":0,"w":5,"h":2,"i":"1", loading: false},
        {"x":8,"y":0,"w":1,"h":1,"i":"2", loading: false},
        {"x":8,"y":1,"w":1,"h":1,"i":"3", loading: false},

        // rop10 station
        {"x":9,"y":0,"w":3,"h":4,"i":"4", loading: false},

        // map
        {"x":3,"y":2,"w":6,"h":4,"i":"5", loading: false},

        // top10 goods
        {"x":0,"y":4,"w":3,"h":4,"i":"6", loading: false},

        // rop10 goods
        {"x":9,"y":4,"w":3,"h":4,"i":"7", loading: false},
        // {"x":3,"y":8,"w":2,"h":2,"i":"8"},
        // {"x":5,"y":8,"w":1,"h":1,"i":"9"},
        // {"x":5,"y":9,"w":1,"h":1,"i":"10"},
        // {"x":6,"y":8,"w":2,"h":2,"i":"11"},

        // profit
        {"x":1,"y":10,"w":4,"h":2,"i":"8", loading: false},
        {"x":0,"y":10,"w":1,"h":1,"i":"9", loading: false},
        {"x":0,"y":11,"w":1,"h":1,"i":"10", loading: false},

        // customer
        {"x":4,"y":8,"w":2,"h":2,"i":"11", loading: false},
        {"x":3,"y":8,"w":1,"h":1,"i":"12", loading: false},
        {"x":3,"y":9,"w":1,"h":1,"i":"13", loading: false},

        // average
        {"x":8,"y":10,"w":3,"h":2,"i":"14", loading: false},
        {"x":11,"y":10,"w":1,"h":1,"i":"15", loading: false},
        {"x":11,"y":11,"w":1,"h":1,"i":"16", loading: false},

        // transfer
        {"x":6,"y":8,"w":2,"h":2,"i":"17", loading: false},
        {"x":8,"y":8,"w":1,"h":1,"i":"18", loading: false},
        {"x":8,"y":9,"w":1,"h":1,"i":"19", loading: false},

        // payments
        {"x":5,"y":10,"w":3,"h":2,"i":"20", loading: false},

        // {"x":0,"y":10,"w":2,"h":2,"i":"14"},
        // {"x":2,"y":10,"w":2,"h":2,"i":"15"},
        // {"x":4,"y":10,"w":2,"h":2,"i":"16"},
        // {"x":6,"y":10,"w":6,"h":2,"i":"17"},
      ],
      mapOptions: {
        backgroundColor: 'transparent',
        title: {
          text: '重点油站分布',
          left: 'center',
          textStyle: {
            color: 'rgba(255, 255, 255, 0.7)'
          }
        },

        tooltip: {
          enterable: true,
          formatter: function(params) {
            var value = params.value;
            var rank = value[3] > 0 ? '' : '倒数';
            return `${params.name} (排名 ${rank} #${Math.abs(value[3])}):<br/>销售额: ${value[2]}`;
          }
        },

        geo: {
          zoom: 1.8,
          map: '北京',
          roam: true,
          itemStyle: {
            color: '#004981',
            borderColor: 'rgb(54,192,118)'
          }
        },

        series: [
          {
            type: 'effectScatter',
            coordinateSystem: 'geo',
            symbol: 'pin',
            symbolSize: '18',
            itemStyle: {                   // 配置每个数据点的样式
              color: function(params) {
                var color = '';
                var value = params.value;
                if (value[3] > 0) {
                  color = '#67C23A';
                } else {
                  color = '#F56C6C';
                }
                return color;
              }
            },
            data: []
          }
        ]
      },
      transferOptions: {
        tooltip: {
          trigger: 'item',
          formatter: '{a}'
        },
        series: [{
          color:["#00F5FF","#909399"],
          name: '油非转换率',
          type: 'pie',
          radius: ['50%', '95%'],
          avoidLabelOverlap: false,
          label: {
            normal: {
              show: true,
              position: 'center'
            },
            emphasis: {
              show: true,
              textStyle: {
                fontSize: '20',
              }
            }
          },
          labelLine: {
            normal: {
              show: false
            }
          },
          data: [{
            value: 32.45,
            name: '32.45%',
            label: {
              normal: {
                textStyle: {
                  fontSize: '16',
                }
              }
            },
          },
            {
              value: 67.55,
            }
          ]
        }]
      },
      levelTreeOptions: {
        color: ['#67C23A', '#409EFF', '#E6A23C', '#F56C6C'],
        tooltip: {
          trigger: 'item',
          axisPointer: {            // 坐标轴指示器，坐标轴触发有效
            type: 'shadow'        // 默认为直线，可选为：'line' | 'shadow'
          },
          formatter: function (params) {
            let tooltip = params.name + '站点数' ;
            switch(params.dataIndex) {
              case 1:
                tooltip += `(收入增加, 毛利增加):<br/>${params.value}`
                break;

              case 2:
                tooltip += `(收入增加, 毛利降低):<br/>${params.value}`
                break;

              case 3:
                tooltip += `(收入降低, 毛利增加):<br/>${params.value}`
                break;

              case 4:
                tooltip += `(收入降低, 毛利降低):<br/>${params.value}`
                break;

              default:
                console.log("SHOULD NOT BE HERE!");
                break;
            }
            return tooltip;
          }
        },
        series: [{
          type: 'treemap',
          name: '站点分级',
          data: [],
          sort: null
        }]
      },
      levelOptions: {
        color: ['#3398DB'],
        tooltip: {
          trigger: 'axis',
          axisPointer: {            // 坐标轴指示器，坐标轴触发有效
            type: 'shadow'        // 默认为直线，可选为：'line' | 'shadow'
          },
          formatter: '{b}站点数:<br/>{c}'
        },
        grid: {
          left: '3%',
          right: '4%',
          height: '100%',
          bottom: '3%',
          containLabel: true
        },
        xAxis: [
          {
            type: 'category',
            data: ['Lv1', 'Lv2', 'Lv3', 'Lv4'],
            axisLabel: {
              color: 'white'
            }
          }
        ],
        yAxis: [
          {
            type: 'value',
            axisLabel: {
              color: 'white'
            }
          }
        ],
        series: [
          {
            name: 'StationNum',
            type: 'treemap',
            barWidth: '20%',
            data: [37, 52, 66, 29]
          }
        ]
      },
      payPieOptions: {
        title: {
        },
        tooltip: {
          trigger: 'item',
          formatter: function(args) {
            return `<p>${args.seriesName}</p> <p>${args.marker} ${args.data.name}: ${beautifyNumber(args.data.value.toFixed(2))} (${args.percent}%)</p>`
          },
        },
        legend: {
          textStyle: {
            color: '#FFF'
          },
          pageIconInactiveColor: '#2f4554',
          pageIconColor : '#aaa',
          pageTextStyle: {
            color: '#FFF'
          },
          type: 'scroll',
          orient: 'vertical',
          right: 10,
          top: 20,
          bottom: 20,
          data: [],
          tooltip: {
            show: true,
            trigger: 'item',
            formatter: function(args) {
              let percent = 0
              self.paySource.forEach(e => {
                if (e.name == args.name) {
                  percent = (e.value * 100 / self.paySource.total).toFixed(2)
                }
              })
              return `<p>${args.name}: ${percent}%</p>`
            },
          }
        },
        series: [
          {
            name: '支付方式',
            type: 'pie',
            radius: '95%',
            center: ['30%', '50%'],
            data: [],
            labelLine: {
              show: false,
            },
            label: {
              show: false,
            },
            emphasis: {
              itemStyle: {
                shadowBlur: 10,
                shadowOffsetX: 0,
                shadowColor: 'rgba(0, 0, 0, 0.5)'
              }
            }
          }
        ]
      },
      paymentsMoney: [],
      paymentsCount: [],
    }
  },
  methods: {
    percentFormat(percent) {
      let p = percent + ''
      if (p.length > 5) {
        let q = p.substring(5)
        return `${q.substring(0, 3)}.${q.substring(3, 5)}%`
      } else {
        return `${percent}%`
      }
    },
    showListDetail(row, p1, p2) {
      this[p1] = row;
      this[p2] = true;
    },
    onMainDiv(ele) {
      if (ele.path[0].className != 'el-range-input') {
        this.showPanel = false;
      }
    },
    toggleShowPanel(ev) {
      this.showPanel = !this.showPanel;
      ev.stopPropagation();
    },
    getNewData(ev) {
      ev.stopPropagation();
      this.getAllData();
    },
    getAllData(loading = true) {
      if (loading && (!this.dateRange || this.dateRange.length != 2 || this.dateRange[0].length <= 0 ||
          this.dateRange[1].length <= 0)) {
        message('error', '请选择有效的时间范围');
        return;
      }

      if (!loading && (this.sdate.length <= 0 || this.edate.length <= 0)) {
        message('error', '选择的时间范围无效');
        return;
      }

      if (loading) {
        this.sdate = this.dateRange[0];
        this.edate = this.dateRange[1];
        this.dispOrgLevel = this.selectedOrg == '000' ? 0 : 1;
        this.dispOrgCode = this.selectedOrg;
        this.promotion.title = `${this.levelName[this.selectedOrg]}非油销售数据`;
      }

      let checked = document.querySelector('#checked')
      if(checked.style.display === 'block') {
        this.getThresholdKpi(loading)
        this.getThresholdStations(loading)
        this.getThresholdTopGoods(loading)
        this,this.getThresholdRopGoods(loading)
      }else{
        this.getKpi(loading);
        this.getStations(loading);
        this.getTopGoods(loading);
        this.getRopGoods(loading);
      }
      // this.getKpi(loading);
      //this.getKpiByMonth(loading);
      // this.getStations(loading);
      // this.getTopStations(loading);
      // this.getRopStations(loading);
      // this.getTopGoods(loading);
      // this.getRopGoods(loading);
      this.getPayments(loading);
      this.showPanel = false;

    },
    updatePayType(t) {
      if (t == 0) {
        this.paySource = this.paymentsMoney
        this.payPieOptions.legend.data = this.paymentsMoney.legend
        this.payPieOptions.series[0].data = this.paymentsMoney
        this.payPieOptions.series[0].name = '支付方式(金额)'
      } else {
        this.paySource = this.paymentsCount
        this.payPieOptions.legend.data = this.paymentsCount.legend
        this.payPieOptions.series[0].data = this.paymentsCount
        this.payPieOptions.series[0].name = '支付方式(数量)'
      }
      this.showPaySel = false
    },
    getPayments(loading) {
      if (loading) {
        this.gridLayout[20].loading = true
      }

      doRequest({
        url: '/v1/web/kpi/payments',
        method: 'POST',
        data: {
          beginDate: this.dateRange[0],
          endDate: this.dateRange[1],
          orgCode: this.selectedOrg,
        }
      }, {
        success: res => {
          if (res.payments) {
            this.payPieOptions.series[0].data = []
            this.payPieOptions.legend.data = []

            this.paymentsMoney = []
            this.paymentsCount = []
            this.paymentsMoney.total = 0
            this.paymentsCount.total = 0
            this.paymentsMoney.legend = []
            this.paymentsCount.legend = []

            res.payments.forEach(p => {
              this.paymentsMoney.total += p.money
              this.paymentsMoney.push({
                name: Config.payments[p.id],
                value: p.money
              })

              this.paymentsCount.total += p.count
              this.paymentsCount.push({
                name: Config.payments[p.id],
                value: p.count
              })
            })

            this.paymentsMoney.sort((a, b) => {
              return b.value - a.value
            })
            this.paymentsCount.sort((a, b) => {
              return b.value - a.value
            })
            this.paymentsMoney.forEach(e => {
              this.paymentsMoney.legend.push(e.name)
            })
            this.paymentsCount.forEach(e => {
              this.paymentsCount.legend.push(e.name)
            })
            this.payPieOptions.legend.data = this.paymentsMoney.legend
            this.payPieOptions.series[0].data = this.paymentsMoney
            this.payPieOptions.series[0].name = '支付方式(金额)'
            this.paySource = this.paymentsMoney
          }
        },
        finally: _ => {
          if (loading) {
            this.gridLayout[20].loading = false
          }
        }
      })
    },
    getKpiByMonth(loading) {
      if (loading) {
        this.gridLayout[5].loading = true;
      }

      let ld = 2

      doRequest({
        url: '/v1/web/kpi/bymonth',
        method: 'POST',
        data: {
          startDate: this.dateRange[0],
          endDate: this.dateRange[1],
          orgCode: this.selectedOrg,
          incomeType: 0,
        }

      }, {
        success: res => {
          this.monthProgress.current = this.promotion.sales.value
          this.monthProgress.target = res.monthValue
          this.yearProgress.current = this.promotion.sales.value
          this.yearProgress.target = res.yearValue
          this.quarterProgress.current = this.promotion.sales.value
          this.quarterProgress.target = res.quarterValue
          let percent = 0

          if (this.monthProgress.target <= 0) {
            this.monthProgress.percent = 100
          } else {
            percent = parseFloat((parseFloat(this.monthProgress.current) * 100 / this.monthProgress.target).toFixed(2))
            if (percent > 100) {
              this.monthProgress.percent = 99.99 + (percent / 100000)
            } else
              this.monthProgress.percent = percent
          }

          if (this.yearProgress.target <= 0) {
            this.yearProgress.percent = 100
          } else {
            percent = parseFloat((parseFloat(this.yearProgress.current) * 100 / this.yearProgress.target).toFixed(2))
            if (percent > 100)
              this.yearProgress.percent = 99.99 + (percent / 100000)
            else
              this.yearProgress.percent = percent
          }

          if (this.quarterProgress.target <= 0) {
            this.quarterProgress.percent = 100
          } else {
            percent = parseFloat((parseFloat(this.quarterProgress.current) * 100 / this.quarterProgress.target).toFixed(2))
            if (percent > 100)
              this.quarterProgress.percent = 99.99 + (percent / 100000)
            else
              this.quarterProgress.percent = percent
          }
        },
        fail: err => {
          console.log(err)
          message('error', '加载销售目标数据失败，请稍后再试')
        },
        finally: _ => {
          if (loading) {
            ld --
            if (ld == 0)
              this.gridLayout[5].loading = false;
          }
        }
      })

      doRequest({
        url: '/v1/web/kpi/bymonth',
        method: 'POST',
        data: {
          startDate: this.dateRange[0],
          endDate: this.dateRange[1],
          orgCode: this.selectedOrg,
          incomeType: 1,
        }
      }, {
        success: res => {
          this.monthProfitProgress.current = this.promotion.profit.value
          this.monthProfitProgress.target = res.monthValue
          this.yearProfitProgress.current = this.promotion.profit.value
          this.yearProfitProgress.target = res.yearValue
          this.quarterProfitProgress.current = this.promotion.profit.value
          this.quarterProfitProgress.target = res.quarterValue
          let percent = 0

          if (this.monthProfitProgress.target <= 0) {
            this.monthProfitProgress.percent = 100
          } else {
            percent = parseFloat((parseFloat(this.monthProfitProgress.current) * 100 / this.monthProfitProgress.target).toFixed(2))
            if (percent > 100) {
              this.monthProfitProgress.percent = 99.99 + (percent / 100000)
            } else
              this.monthProfitProgress.percent = percent
          }

          if (this.yearProfitProgress.target <= 0) {
            this.yearProfitProgress.percent = 100
          } else {
            percent = parseFloat((parseFloat(this.yearProfitProgress.current) * 100 / this.yearProfitProgress.target).toFixed(2))
            if (percent > 100)
              this.yearProfitProgress.percent = 99.99 + (percent / 100000)
            else
              this.yearProfitProgress.percent = percent
          }

          if (this.quarterProfitProgress.target <= 0) {
            this.quarterProfitProgress.percent = 100
          } else {
            percent = parseFloat((parseFloat(this.quarterProfitProgress.current) * 100 / this.quarterProfitProgress.target).toFixed(2))
            if (percent > 100)
              this.quarterProfitProgress.percent = 99.99 + (percent / 100000)
            else
              this.quarterProfitProgress.percent = percent
          }
        },
        fail: err => {
          console.log(err)
          message('error', '加载销售目标数据失败，请稍后再试')
        },
        finally: _ => {
          if (loading) {
            ld --
            if (ld == 0)
              this.gridLayout[5].loading = false;
          }
        }
      })
    },
    getKpi(loading = true) {
      if (loading) {
        for (let i = 1; i < 20; i ++) {
          if (i <= 3 || i >=8 || i == 5) {
            this.gridLayout[i].loading = true;
          }
        }
      }

      doRequest({
        url: '/v1/web/sales/kpi',
        method: 'POST',
        data: {
          beginDate: this.sdate,
          endDate: this.edate,
          orgLevel: this.dispOrgLevel,
          orgCode: this.dispOrgCode
        }
      }, {
        success: res => {
          this.promotion.customers = res.kpi.nonFuelCount;
          this.promotion.transfer = res.kpi.fnConversionRate;
          this.promotion.sales = res.kpi.netIncome;
          this.promotion.profit = res.kpi.grossMargin;
          this.promotion.average = res.kpi.avgTrxValue;
          this.promotion.y2nIncome = res.kpi.y2nIncome;
          this.promotion.q2nIncome = res.kpi.q2nIncome;
          this.promotion.m2nIncome = res.kpi.m2nIncome
          this.getKpiByMonth(loading)
        },
        fail: err => {
          message('error', '加载销售数据失败，请稍后再试')
        },
        finally: _ => {
          if (loading) {
            for (let i = 1; i < 20; i ++) {
              if (i <= 3 || i >=8) {
                this.gridLayout[i].loading = false;
              }
            }
          }
        }
      })
    },
    // 获得剔除5%毛利率的kpi
    getThresholdKpi(loading = true) {
      if (loading) {
        for (let i = 1; i < 20; i ++) {
          if (i <= 3 || i >=8 || i == 5) {
            this.gridLayout[i].loading = true;
          }
        }
      }

      doRequest({
        url: '/v1/web/sales/kpi',
        method: 'POST',
        data: {
          beginDate: this.sdate,
          endDate: this.edate,
          orgLevel: this.dispOrgLevel,
          orgCode: this.dispOrgCode,
          rateThreshold: 0.05
        }
      }, {
        success: res => {
          this.promotion.customers = res.kpi.nonFuelCount;
          this.promotion.transfer = res.kpi.fnConversionRate;
          this.promotion.sales = res.kpi.netIncome;
          this.promotion.profit = res.kpi.grossMargin;
          this.promotion.average = res.kpi.avgTrxValue;
          this.promotion.y2nIncome = res.kpi.y2nIncome;
          this.promotion.q2nIncome = res.kpi.q2nIncome;
          this.promotion.m2nIncome = res.kpi.m2nIncome
          this.getKpiByMonth(loading)
        },
        fail: err => {
          message('error', '加载销售数据失败，请稍后再试')
        },
        finally: _ => {
          if (loading) {
            for (let i = 1; i < 20; i ++) {
              if (i <= 3 || i >=8) {
                this.gridLayout[i].loading = false;
              }
            }
          }
        }
      })
    },
    getStations(loading) {
      if (loading) {
        this.gridLayout[0].loading = true;
        this.gridLayout[4].loading = true;
        //this.gridLayout[5].loading = true;
      }

      allRequests({
        success: res => {
          for (let i = 0; i < res.length; i ++) {
            if (res[i].config.data.includes('"DESC"')) {
              this.promotion.topStations = res[i].data.plantList;
            } else {
              this.promotion.ropStations = res[i].data.plantList;
            }
          }

          if (!this.promotion.topStations || !this.promotion.ropStations)
            return;

          let a = null, b = null, c = this.mapOptions.series[0].data;
          for (let j = 0; j < 10; j ++) {
            a = null;
            b = null;
            for (let k = 0; k < c.length; k ++) {
              if (a == null && this.promotion.topStations[j].OrgText == c[k].name) {
                a = c[k];
              } else if (b == null && this.promotion.ropStations[j].OrgText == c[k].name) {
                b = c[k];
              }

              if (a && b)
                break;
            }

            if (!a || !b) {
              if (this.$refs['cnpc-promotions-map']) {
                this.$refs['cnpc-promotions-map'].clear();
                this.mapOptions.series[0].data = [];
                for (j = 0; j < 10; j ++) {
                  a = this.promotion.topStations[j];
                  b = this.promotion.ropStations[j];

                  if (a)
                    this.mapOptions.series[0].data.push({
                      name: a.OrgText,
                      value: [a.PosX, a.PosY, a.Metric.NetvalInv, j + 1]
                    })

                  if (b)
                    this.mapOptions.series[0].data.push({
                      name: b.OrgText,
                      value: [b.PosX, b.PosY, b.Metric.NetvalInv, -j - 1]
                    })
                }
                this.$refs['cnpc-promotions-map'].mergeOptions(this.mapOptions, true, false);
              }
              break;
            }
          }
        },
        fail: err => {
          console.log(err)
          message('error', '获取油站排名失败，请稍后再试');
        },
        finally: _ => {
          if (loading) {
            this.gridLayout[0].loading = false;
            this.gridLayout[4].loading = false;
            //this.gridLayout[5].loading = false;
          }
        }
      }, [], [{
        url: '/v1/web/sales/plant/rank',
        method: 'POST',
        data: {
          beginDate: this.sdate,
          endDate: this.edate,
          orgLevel: this.dispOrgLevel,
          orgCode: this.dispOrgCode,
          sortBy: 'DESC',
          limit: 10
        }
      }, {
        url: '/v1/web/sales/plant/rank',
        method: 'POST',
        data: {
          beginDate: this.sdate,
          endDate: this.edate,
          orgLevel: this.dispOrgLevel,
          orgCode: this.dispOrgCode,
          sortBy: 'ASC',
          limit: 10
        }
      }])
    },
    // 获得剔除5%毛利率的站级排名
    getThresholdStations(loading) {
      if (loading) {
        this.gridLayout[0].loading = true;
        this.gridLayout[4].loading = true;
        //this.gridLayout[5].loading = true;
      }

      allRequests({
        success: res => {
          for (let i = 0; i < res.length; i ++) {
            if (res[i].config.data.includes('"DESC"')) {
              this.promotion.topStations = res[i].data.plantList;
            } else {
              this.promotion.ropStations = res[i].data.plantList;
            }
          }

          if (!this.promotion.topStations || !this.promotion.ropStations)
            return;

          let a = null, b = null, c = this.mapOptions.series[0].data;
          for (let j = 0; j < 10; j ++) {
            a = null;
            b = null;
            for (let k = 0; k < c.length; k ++) {
              if (a == null && this.promotion.topStations[j].OrgText == c[k].name) {
                a = c[k];
              } else if (b == null && this.promotion.ropStations[j].OrgText == c[k].name) {
                b = c[k];
              }

              if (a && b)
                break;
            }

            if (!a || !b) {
              if (this.$refs['cnpc-promotions-map']) {
                this.$refs['cnpc-promotions-map'].clear();
                this.mapOptions.series[0].data = [];
                for (j = 0; j < 10; j ++) {
                  a = this.promotion.topStations[j];
                  b = this.promotion.ropStations[j];

                  if (a)
                    this.mapOptions.series[0].data.push({
                      name: a.OrgText,
                      value: [a.PosX, a.PosY, a.Metric.NetvalInv, j + 1]
                    })

                  if (b)
                    this.mapOptions.series[0].data.push({
                      name: b.OrgText,
                      value: [b.PosX, b.PosY, b.Metric.NetvalInv, -j - 1]
                    })
                }
                this.$refs['cnpc-promotions-map'].mergeOptions(this.mapOptions, true, false);
              }
              break;
            }
          }
        },
        fail: err => {
          console.log(err)
          message('error', '获取油站排名失败，请稍后再试');
        },
        finally: _ => {
          if (loading) {
            this.gridLayout[0].loading = false;
            this.gridLayout[4].loading = false;
            //this.gridLayout[5].loading = false;
          }
        }
      }, [], [{
        url: '/v1/web/sales/plant/rank',
        method: 'POST',
        data: {
          beginDate: this.sdate,
          endDate: this.edate,
          orgLevel: this.dispOrgLevel,
          orgCode: this.dispOrgCode,
          rateThreshold: 0.05,
          sortBy: 'DESC',
          limit: 10
        }
      }, {
        url: '/v1/web/sales/plant/rank',
        method: 'POST',
        data: {
          beginDate: this.sdate,
          endDate: this.edate,
          orgLevel: this.dispOrgLevel,
          orgCode: this.dispOrgCode,
          rateThreshold: 0.05,
          sortBy: 'ASC',
          limit: 10
        }
      }])
    },
    getTopStations(loading) {
      if (loading)
        this.gridLayout[0].loading = true;
      doRequest({
        url: '/v1/web/sales/plant/rank',
        method: 'POST',
        data: {
          beginDate: this.sdate,
          endDate: this.edate,
          orgLevel: this.dispOrgLevel,
          orgCode: this.dispOrgCode,
          sortBy: 'DESC',
          limit: 10
        }
      }, {
        success: res => {
          let e = null;
          this.promotion.topStations = res.plantList;
          for (let j = 0; j < this.promotion.topStations.length; j ++) {
            e = null;
            for (let k = 0; k < this.mapOptions.series[0].data.length; k ++) {
              if (this.promotion.topStations[j].OrgText == this.mapOptions.series[0].data[k].name) {
                e = this.mapOptions.series[0].data[k];
                break;
              }
            }
            if (!e) {
              break;
            }
          }

          if (!e) {
            for (let i = 0; i < this.promotion.topStations.length; i ++) {
              e = this.promotion.topStations[i];
              this.mapOptions.series[0].data.push({
                name: e.OrgText,
                value: [e.PosX, e.PosY, e.Metric.NetvalInv, i + 1]
              })
            }
            this.$refs['cnpc-promotions-map'].mergeOptions(this.mapOptions, true);
          }
        },
        fail: err => {
          message('error', '获取油站排名失败，请稍后再试')
        },
        finally: _ => {
          if (loading)
            this.gridLayout[0].loading = false;
        }
      })
    },
    getRopStations(loading) {
      if (loading)
        this.gridLayout[4].loading = true;
      doRequest({
        url: '/v1/web/sales/plant/rank',
        method: 'POST',
        data: {
          beginDate: this.sdate,
          endDate: this.edate,
          orgLevel: this.dispOrgLevel,
          orgCode: this.dispOrgCode,
          sortBy: 'ASC',
          limit: 10
        }
      }, {
        success: res => {
          let e = null;
          this.promotion.ropStations = res.plantList;
          for (let j = 0; j < this.promotion.ropStations.length; j ++) {
            e = null;
            for (let k = 0; k < this.mapOptions.series[0].data.length; k ++) {
              if (this.promotion.ropStations[j].OrgText == this.mapOptions.series[0].data[k].name) {
                e = this.mapOptions.series[0].data[k];
                break;
              }
            }
            if (!e) {
              break;
            }
          }

          if (!e) {
            for (let i = 0; i < this.promotion.ropStations.length; i ++) {
              e = this.promotion.ropStations[i];
              this.mapOptions.series[0].data.push({
                name: e.OrgText,
                value: [e.PosX, e.PosY, e.Metric.NetvalInv, -i - 1]
              })
            }
            this.$refs['cnpc-promotions-map'].mergeOptions(this.mapOptions, true);
          }
        },
        fail: err => {
          message('error', '获取油站排名失败，请稍后再试')
        },
        finally: _ => {
          if (loading)
            this.gridLayout[4].loading = false;
        }
      })
    },
    getTopGoods(loading) {
      if (loading)
        this.gridLayout[6].loading = true;
      doRequest({
        url: '/v1/web/sales/material/rank',
        method: 'POST',
        data: {
          beginDate: this.sdate,
          endDate: this.edate,
          orgLevel: this.dispOrgLevel,
          orgCode: this.dispOrgCode,
          sortBy: 'DESC',
          limit: 10
        }
      }, {
        success: res => {
          this.promotion.topGoods = res.matlList;
        },
        fail: err => {
          message('error', '获取商品排名失败，请稍后再试')
        },
        finally: _ => {
          if (loading)
            this.gridLayout[6].loading = false;
        }
      })
    },
    // 获取剔除5%毛利率的畅销售品
    getThresholdTopGoods(loading) {
      if (loading)
        this.gridLayout[6].loading = true;
      doRequest({
        url: '/v1/web/sales/material/rank',
        method: 'POST',
        data: {
          beginDate: this.sdate,
          endDate: this.edate,
          orgLevel: this.dispOrgLevel,
          orgCode: this.dispOrgCode,
          rateThreshold: 0.05,
          sortBy: 'DESC',
          limit: 10
        }
      }, {
        success: res => {
          this.promotion.topGoods = res.matlList;
        },
        fail: err => {
          message('error', '获取商品排名失败，请稍后再试')
        },
        finally: _ => {
          if (loading)
            this.gridLayout[6].loading = false;
        }
      })
    },
    getRopGoods(loading) {
      if (loading)
        this.gridLayout[7].loading = true;
      doRequest({
        url: '/v1/web/sales/material/rank',
        method: 'POST',
        data: {
          beginDate: this.sdate,
          endDate: this.edate,
          orgLevel: this.dispOrgLevel,
          orgCode: this.dispOrgCode,
          sortBy: 'ASC',
          limit: 10
        }
      }, {
        success: res => {
          this.promotion.ropGoods = res.matlList;
        },
        fail: err => {
          message('error', '获取商品排名失败，请稍后再试')
        },
        finally: _ => {
          if (loading)
            this.gridLayout[7].loading = false;
        }
      })
    },
    // 获取剔除5%毛利率的滞销售品
    getThresholdRopGoods(loading) {
      if (loading)
        this.gridLayout[7].loading = true;
      doRequest({
        url: '/v1/web/sales/material/rank',
        method: 'POST',
        data: {
          beginDate: this.sdate,
          endDate: this.edate,
          orgLevel: this.dispOrgLevel,
          orgCode: this.dispOrgCode,
          rateThreshold: 0.05,
          sortBy: 'ASC',
          limit: 10
        }
      }, {
        success: res => {
          this.promotion.ropGoods = res.matlList;
        },
        fail: err => {
          message('error', '获取商品排名失败，请稍后再试')
        },
        finally: _ => {
          if (loading)
            this.gridLayout[7].loading = false;
        }
      })
    },
    doFullscreen() {
      let ele = document.getElementById("cnpc-promo-main-v3");
      if (ele) {
        if (!document.fullscreenElement) {
          ele.requestFullscreen();
        }
        this.fullScreen = true;
      }
    },
    exitFullscreen() {
      if (document.fullscreenElement) {
        document.exitFullscreen();
      }
      this.fullScreen = false;
    },
    promoTimer() {
      this.currentTimeDesc = moment(Date.now()).format('YYYY/MM/DD HH:mm:ss');
      if (!document.fullscreenElement)
        this.fullScreen = false;

      this.ticker += 1;
      if (this.ticker === 3000) {
        this.ticker = 0;
        this.getAllData(false);
      }
    },
    onKeyDown(ev) {
      if (ev.keyCode === 27) {
        this.fullScreen = false;
      }

      if (ev.keyCode === 113) {
        // F2, do fullscreen
        this.doFullscreen();
        ev.stopPropagation();
      }
    },
    chnNumber(num) {
      return beautifyNumberChn(num, 10000);
    },
    beautyNum(num) {
      return beautifyNumber(num);
    },
    onDonutResize(ele) {
      let h = ele.offsetHeight, w = ele.offsetWidth;
      let m = Math.min(h, w);
      m -= 8;
      this.donutWidth = ~~(m * 100 / w);
    },
    onRootResize(root) {
      //this.$refs['cnpc-promotions-map'].clear();
      //this.$refs['cnpc-promotions-map'].mergeOptions(this.mapOptions, true, false);
      this.gridRowHeight = window.innerHeight * 0.08
      fitty.fitAll();
      // let div = document.getElementById("progress-grid").parentElement
      // let a = div.offsetHeight - 60, b = (div.offsetWidth - 32 ) / 2
      // this.progressWidth = Math.min(a, b)
      // this.fitties.forEach(e => {
      //     e.forEach(i => {
      //         console.log(i)
      //         i.fit();
      //     })
      // })
    },
    // 获得剔除5%以下的数据
    getThresholdData(loading = true) {
      let square = document.querySelector('#toCheck')
      let checked = document.querySelector('#checked')
      square.style.display = 'none'
      checked.style.display = 'block'

      if (!this.dispOrgCode || !this.dateRange || !this.dateRange[0] || !this.dateRange[1]) {
        return;
      }

      this.getThresholdKpi(loading)
      this.getThresholdStations(loading)
      this.getThresholdTopGoods(loading)
      this,this.getThresholdRopGoods(loading)
    },
    // 获得完整的数据
    getNormalData(loading) {
      let square = document.querySelector('#toCheck')
      let checked = document.querySelector('#checked')
      square.style.display = 'block'
      checked.style.display = 'none'

      if (!this.dispOrgCode || !this.dateRange || !this.dateRange[0] || !this.dateRange[1]) {
        return;
      }
      this.getKpi(loading);
      this.getStations(loading);
      this.getTopGoods(loading);
      this.getRopGoods(loading);
    }
  }
}
</script>
