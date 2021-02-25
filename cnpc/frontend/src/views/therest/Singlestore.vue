<template>
    <div class="cnpc-single-store-main" id="cnpc-single-store-main" @click="onMainDiv"
        v-resize:debounce.50="onRootResize">
        <div class="single-store-header">
            <p class="title">
                {{ stationInfo.plantName }}
            </p>
            <p class="sub-title">
                {{ sdate }} - {{ edate }}
            </p>

            <i v-if="fullScreen" class="ri-fullscreen-exit-line single-full-screen"
                @click="exitFullscreen"></i>
            <i v-else class="ri-fullscreen-line single-full-screen" @click="doFullscreen"></i>

            <i class="ri-arrow-right-s-fill single-level-select-fa" @click="toggleShowPanel"></i>

            <font-awesome-icon id="toCheck" class="remove-unchecked" icon="square" v-on:click="getThresholdData"></font-awesome-icon>
            <font-awesome-icon id="checked" class="remove-checked" icon="check-square" v-on:click="getNormalData" style="display: none"></font-awesome-icon>
            <div v-if="showPanel" class="single-level-panel">
                <el-select :popper-append-to-body="false" size="mini" style="width: 100%"
                    v-model="selectedStationCode" filterable placeholder="请选择关注油站">
                    <el-option v-for="s in stations" :key="s.orgCode"
                        :label="`${s.orgText}(${s.orgCode})`"
                        :value="s.orgCode"/>
                </el-select>
                <el-date-picker size="mini" style="width: 100%; margin-top: 6px"
                    type="daterange" v-model="dateRange" range-separator="-"
                    start-placeholder="开始日期"
                    end-placeholder="结束日期"
                    value-format="yyyy-MM-dd" />
                <el-button size="mini" style="width: 100%; margin-top: 6px" type="primary"
                    @click="getAllData">查询</el-button>
            </div>

            <span class="timer">{{ currentTimeDesc }}</span>
        </div>
        <grid-layout class="single-store-grid" :layout.sync="gridLayout"
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
                        畅销商品
                    </p>
                    <div v-if="topGoodDetail" class="detail-display">
                        <i class="el-icon-close close-button" @click="topGoodDetail = false"/>
                        <p class="detail-disp-title">{{ detailTopGood.MaterialTxt }}</p>
                        <el-form size="mini" label-position="right" label-width="70px">
                            <el-form-item label="销售收入">
                                <p>{{ detailTopGood.Metric ?
                                    beautyNum(detailTopGood.Metric.NetvalInv, 2) :
                                    0 }}</p>
                            </el-form-item>
                            <el-form-item label="销售毛利">
                                <p>{{ detailTopGood.Metric ?
                                    beautyNum(detailTopGood.Metric.GrossProfit, 2) :
                                    0 }}</p>
                            </el-form-item>
                            <el-form-item label="毛利率">
                                <p>{{ detailTopGood.Metric ?
                                    beautyNum(detailTopGood.Metric.GrossMargin, 2) :
                                    0 }}%</p>
                            </el-form-item>
                        </el-form>
                    </div>
                    <table v-else>
                        <tbody>
                            <tr v-for="(v, i) in topGoods" :key="'tr-top-good-key-' + i"
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
                        {{ beautyNum(sales.value, 0) }}
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
                <div class="percent-grid" v-if="!isNaN(sales.yoy)">
                    <p class="percent-title">
                        同比
                        <font-awesome-icon v-if="sales.yoy > 0" icon="arrow-up" class="percent-arrow-up" />
                        <font-awesome-icon v-else icon="arrow-down" class="percent-arrow-down" />
                    </p>
                    <p v-if="sales.yoy > 0" class="percent-number-up">
                        {{ sales.yoy.toFixed(2) }}%
                    </p>
                    <p v-else class="percent-number-down">
                        {{ sales.yoy.toFixed(2) }}%
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
                <div class="percent-grid" v-if="!isNaN(sales.mom)">
                    <p class="percent-title">
                        环比
                        <font-awesome-icon v-if="sales.mom > 0" icon="arrow-up" class="percent-arrow-up" />
                        <font-awesome-icon v-else icon="arrow-down" class="percent-arrow-down" />
                    </p>
                    <p v-if="sales.mom > 0" class="percent-number-up">
                        {{ sales.mom.toFixed(2) }}%
                    </p>
                    <p v-else class="percent-number-down">
                        {{ sales.mom.toFixed(2) }}%
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
                        滞销商品
                    </p>
                    <div v-if="ropGoodDetail" class="detail-display">
                        <i class="el-icon-close close-button" @click="ropGoodDetail = false"/>
                        <p class="detail-disp-title">{{ detailRopGood.MaterialTxt }}</p>
                        <el-form size="mini" label-position="right" label-width="70px">
                            <el-form-item label="销售收入">
                                <p>{{ detailRopGood.Metric ?
                                    beautyNum(detailRopGood.Metric.NetvalInv, 2) :
                                    0 }}</p>
                            </el-form-item>
                            <el-form-item label="销售毛利">
                                <p>{{ detailRopGood.Metric ?
                                    beautyNum(detailRopGood.Metric.GrossProfit, 2) :
                                    0 }}</p>
                            </el-form-item>
                            <el-form-item label="毛利率">
                                <p>{{ detailRopGood.Metric ?
                                    beautyNum(detailRopGood.Metric.GrossMargin, 2) :
                                    0 }}%</p>
                            </el-form-item>
                        </el-form>
                    </div>
                    <table v-else>
                        <tbody>
                            <tr v-for="(v, i) in ropGoods" :key="'tr-top-good-key-' + i"
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
                :x="gridLayout[5].x"
                :y="gridLayout[5].y"
                :w="gridLayout[5].w"
                :h="gridLayout[5].h"
                :i="gridLayout[5].i"
                :key="gridLayout[5].i"
                v-loading="gridLayout[5].loading">
                <v-chart ref="focus-line-chart" theme="macarons" :options="trendLineOptions"
                    :autoresize="true"/>
            </grid-item>
            <grid-item
                :x="gridLayout[6].x"
                :y="gridLayout[6].y"
                :w="gridLayout[6].w"
                :h="gridLayout[6].h"
                :i="gridLayout[6].i"
                :key="gridLayout[6].i"
                v-loading="gridLayout[6].loading">
                <div id="donut-auto-resize-transfer" style="width: 100%; height: 100%"
                    v-resize:debounce.50="onDonutResize">
                    <vc-donut :size="donutWidth" unit="%" :thickness="30"
                        background="#1A3A89"
                        :sections="[{ value: transfer.value, color: percentColors[~~(transfer.value/20)].color }]"
                        foreground="#C0C4CC">
                        <p style="color: rgba(255, 255, 255, 1); font-size: 14px; margin-bottom: 2px">
                            油非转换
                        </p>
                        <p :style="`color: ${percentColors[~~(transfer.value/20)].color}; font-size: 21px; font-weight: bold`">
                            {{ transfer.value.toFixed(2) }}%
                        <p/>
                    </vc-donut>
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
                <div class="percent-grid" v-if="!isNaN(transfer.yoy)">
                    <p class="percent-title">
                        同比
                        <font-awesome-icon v-if="transfer.yoy > 0" icon="arrow-up" class="percent-arrow-up" />
                        <font-awesome-icon v-else icon="arrow-down" class="percent-arrow-down" />
                    </p>
                    <p v-if="transfer.yoy > 0" class="percent-number-up">
                        {{ transfer.yoy.toFixed(2) }}%
                    </p>
                    <p v-else class="percent-number-down">
                        {{ transfer.yoy.toFixed(2) }}%
                    </p>
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
                <div class="percent-grid" v-if="!isNaN(transfer.mom)">
                    <p class="percent-title">
                        环比
                        <font-awesome-icon v-if="transfer.mom > 0" icon="arrow-up" class="percent-arrow-up" />
                        <font-awesome-icon v-else icon="arrow-down" class="percent-arrow-down" />
                    </p>
                    <p v-if="transfer.mom > 0" class="percent-number-up">
                        {{ transfer.mom.toFixed(2) }}%
                    </p>
                    <p v-else class="percent-number-down">
                        {{ transfer.mom.toFixed(2) }}%
                    </p>
                </div>
            </grid-item>
            <grid-item
                :x="gridLayout[9].x"
                :y="gridLayout[9].y"
                :w="gridLayout[9].w"
                :h="gridLayout[9].h"
                :i="gridLayout[9].i"
                :key="gridLayout[9].i"
                v-loading="gridLayout[9].loading">
                <div v-if="showCatSel" class="single-cat-panel">
                    <el-button size="mini" style="width: 100%;" type="primary"
                        @click="updateCatType('income')">按销售收入</el-button>
                    <el-button size="mini" style="width: 100%; margin: 6px 0 0 0" type="primary"
                        @click="updateCatType('profit')">按毛利</el-button>
                    <el-button size="mini" style="width: 100%; margin: 6px 0 0 0" type="primary"
                        @click="updateCatType('margin')">按毛利率</el-button>
                </div>
                <i class="ri-arrow-down-s-line cat-select" @click="showCatSel = !showCatSel"></i>
                <v-chart v-show="showCatPie" ref="focus-pie-chart" theme="macarons" :options="catPieOptions"
                    :autoresize="true"/>

                <!--i v-on:click="toCatPie" v-on:mouseover="onCatPie" v-on:mouseout="outCatPie" class="el-icon-arrow-left prev-select"></i>
                <i v-on:click="toPayPie" v-on:mouseover="onPayPie" v-on:mouseout="outPayPie" class="el-icon-arrow-right next-select"></i>
                <v-chart v-show="showPayPie" ref="focus-pie-chart" theme="macarons" :options="payPieOptions"
                         :autoresize="true"/-->
            </grid-item>
            <grid-item
                :x="gridLayout[10].x"
                :y="gridLayout[10].y"
                :w="gridLayout[10].w"
                :h="gridLayout[10].h"
                :i="gridLayout[10].i"
                :key="gridLayout[10].i"
                v-loading="gridLayout[10].loading">
                <div class="sale-amount-grid">
                    <p class="sale-amount-title">销售毛利</p>
                    <p class="sale-amount-normal">
                        <font-awesome-icon icon="funnel-dollar"></font-awesome-icon>
                        {{ beautyNum(profit.value, 0) }}
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
                <div class="percent-grid" v-if="!isNaN(profit.yoy)">
                    <p class="percent-title">
                        同比
                        <font-awesome-icon v-if="profit.yoy > 0" icon="arrow-up" class="percent-arrow-up" />
                        <font-awesome-icon v-else icon="arrow-down" class="percent-arrow-down" />
                    </p>
                    <p v-if="profit.yoy > 0" class="percent-number-up">
                        {{ profit.yoy.toFixed(2) }}%
                    </p>
                    <p v-else class="percent-number-down">
                        {{ profit.yoy.toFixed(2) }}%
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
                <div class="percent-grid" v-if="!isNaN(profit.mom)">
                    <p class="percent-title">
                        环比
                        <font-awesome-icon v-if="profit.mom > 0" icon="arrow-up" class="percent-arrow-up" />
                        <font-awesome-icon v-else icon="arrow-down" class="percent-arrow-down" />
                    </p>
                    <p v-if="profit.mom > 0" class="percent-number-up">
                        {{ profit.mom.toFixed(2) }}%
                    </p>
                    <p v-else class="percent-number-down">
                        {{ profit.mom.toFixed(2) }}%
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
                <div class="sale-amount-grid">
                    <p class="sale-amount-title">毛利率</p>
                    <p class="sale-amount-normal">
                        <font-awesome-icon icon="hand-holding-usd"></font-awesome-icon>
                        {{ margin.value.toFixed(2) }}%
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
                <div class="percent-grid" v-if="!isNaN(margin.yoy)">
                    <p class="percent-title">
                        同比
                        <font-awesome-icon v-if="margin.yoy > 0" icon="arrow-up" class="percent-arrow-up" />
                        <font-awesome-icon v-else icon="arrow-down" class="percent-arrow-down" />
                    </p>
                    <p v-if="margin.yoy > 0" class="percent-number-up">
                        {{ margin.yoy.toFixed(2) }}%
                    </p>
                    <p v-else class="percent-number-down">
                        {{ margin.yoy.toFixed(2) }}%
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
                <div class="percent-grid" v-if="!isNaN(margin.mom)">
                    <p class="percent-title">
                        环比
                        <font-awesome-icon v-if="margin.mom > 0" icon="arrow-up" class="percent-arrow-up" />
                        <font-awesome-icon v-else icon="arrow-down" class="percent-arrow-down" />
                    </p>
                    <p v-if="margin.mom > 0" class="percent-number-up">
                        {{ margin.mom.toFixed(2) }}%
                    </p>
                    <p v-else class="percent-number-down">
                        {{ margin.mom.toFixed(2) }}%
                    </p>
                </div>
            </grid-item>
            <grid-item
                :x="gridLayout[16].x"
                :y="gridLayout[16].y"
                :w="gridLayout[16].w"
                :h="gridLayout[16].h"
                :i="gridLayout[16].i"
                :key="gridLayout[16].i"
                v-loading="gridLayout[16].loading">
                <div class="sale-amount-grid">
                    <p class="sale-amount-title">客单价</p>
                    <p class="sale-amount-normal">
                        <font-awesome-icon icon="coins"></font-awesome-icon>
                        {{ beautyNum(average.value, 2) }}
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
                <div class="percent-grid" v-if="!isNaN(average.yoy)">
                    <p class="percent-title">
                        同比
                        <font-awesome-icon v-if="average.yoy > 0" icon="arrow-up" class="percent-arrow-up" />
                        <font-awesome-icon v-else icon="arrow-down" class="percent-arrow-down" />
                    </p>
                    <p v-if="average.yoy > 0" class="percent-number-up">
                        {{ average.yoy.toFixed(2) }}%
                    </p>
                    <p v-else class="percent-number-down">
                        {{ average.yoy.toFixed(2) }}%
                    </p>
                </div>
            </grid-item>
            <grid-item
                :x="gridLayout[18].x"
                :y="gridLayout[18].y"
                :w="gridLayout[18].w"
                :h="gridLayout[18].h"
                :i="gridLayout[18].i"
                :key="gridLayout[18].i"
                v-loading="gridLayout[18].loading">
                <div class="percent-grid" v-if="!isNaN(average.mom)">
                    <p class="percent-title">
                        环比
                        <font-awesome-icon v-if="average.mom > 0" icon="arrow-up" class="percent-arrow-up" />
                        <font-awesome-icon v-else icon="arrow-down" class="percent-arrow-down" />
                    </p>
                    <p v-if="average.mom > 0" class="percent-number-up">
                        {{ average.mom.toFixed(2) }}%
                    </p>
                    <p v-else class="percent-number-down">
                        {{ average.mom.toFixed(2) }}%
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
                <div class="sale-amount-grid">
                    <p class="sale-amount-title">
                        非油消费人数
                    </p>
                    <p class="sale-amount-normal2">
                        <font-awesome-icon icon="user-friends" />
                        {{ beautyNum(customers.value, 0) }}
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
                <div class="percent-grid" v-if="!isNaN(customers.yoy)">
                    <p class="percent-title">
                        同比
                        <font-awesome-icon v-if="customers.yoy > 0" icon="arrow-up" class="percent-arrow-up" />
                        <font-awesome-icon v-else icon="arrow-down" class="percent-arrow-down" />
                    </p>
                    <p v-if="customers.yoy > 0" class="percent-number-up">
                        {{ customers.yoy.toFixed(2) }}%
                    </p>
                    <p v-else class="percent-number-down">
                        {{ customers.yoy.toFixed(2) }}%
                    </p>
                </div>
            </grid-item>
            <grid-item
                :x="gridLayout[21].x"
                :y="gridLayout[21].y"
                :w="gridLayout[21].w"
                :h="gridLayout[21].h"
                :i="gridLayout[21].i"
                :key="gridLayout[21].i"
                v-loading="gridLayout[21].loading">
                <div class="percent-grid" v-if="!isNaN(customers.mom)">
                    <p class="percent-title">
                        环比
                        <font-awesome-icon v-if="customers.mom > 0" icon="arrow-up" class="percent-arrow-up" />
                        <font-awesome-icon v-else icon="arrow-down" class="percent-arrow-down" />
                    </p>
                    <p v-if="customers.mom > 0" class="percent-number-up">
                        {{ customers.mom.toFixed(2) }}%
                    </p>
                    <p v-else class="percent-number-down">
                        {{ customers.mom.toFixed(2) }}%
                    </p>
                </div>
            </grid-item>
            <grid-item
                :x="gridLayout[22].x"
                :y="gridLayout[22].y"
                :w="gridLayout[22].w"
                :h="gridLayout[22].h"
                :i="gridLayout[22].i"
                :key="gridLayout[22].i"
                v-loading="gridLayout[22].loading">
                <div class="sale-amount-grid">
                    <p class="sale-amount-title">
                        坪效
                    </p>
                    <p class="sale-amount-normal2">
                        <font-awesome-icon icon="home" />
                        {{ beautyNum(area.value, 2) }}
                    </p>
                </div>
            </grid-item>
            <grid-item
                :x="gridLayout[23].x"
                :y="gridLayout[23].y"
                :w="gridLayout[23].w"
                :h="gridLayout[23].h"
                :i="gridLayout[23].i"
                :key="gridLayout[23].i"
                v-loading="gridLayout[23].loading">
                <div class="percent-grid" v-if="!isNaN(area.yoy)">
                    <p class="percent-title">
                        同比
                        <font-awesome-icon v-if="area.yoy > 0" icon="arrow-up" class="percent-arrow-up" />
                        <font-awesome-icon v-else icon="arrow-down" class="percent-arrow-down" />
                    </p>
                    <p v-if="area.yoy > 0" class="percent-number-up">
                        {{ area.yoy.toFixed(2) }}%
                    </p>
                    <p v-else class="percent-number-down">
                        {{ area.yoy.toFixed(2) }}%
                    </p>
                </div>
            </grid-item>
            <grid-item
                :x="gridLayout[24].x"
                :y="gridLayout[24].y"
                :w="gridLayout[24].w"
                :h="gridLayout[24].h"
                :i="gridLayout[24].i"
                :key="gridLayout[24].i"
                v-loading="gridLayout[24].loading">
                <div class="percent-grid" v-if="!isNaN(area.mom)">
                    <p class="percent-title">
                        环比
                        <font-awesome-icon v-if="area.mom > 0" icon="arrow-up" class="percent-arrow-up" />
                        <font-awesome-icon v-else icon="arrow-down" class="percent-arrow-down" />
                    </p>
                    <p v-if="area.mom > 0" class="percent-number-up">
                        {{ area.mom.toFixed(2) }}%
                    </p>
                    <p v-else class="percent-number-down">
                        {{ area.mom.toFixed(2) }}%
                    </p>
                </div>
            </grid-item>
            <grid-item
                :x="gridLayout[25].x"
                :y="gridLayout[25].y"
                :w="gridLayout[25].w"
                :h="gridLayout[25].h"
                :i="gridLayout[25].i"
                :key="gridLayout[25].i"
                v-loading="gridLayout[25].loading">
                <!--div class="info-grid" v-if="area.mom">
                    <p class="info-title">
                        加油站星级
                    </p>
                    <p class="info-desc">
                        <font-awesome-icon icon="star" />
                        {{ stationInfo.ztxtXjms }}
                    </p>
                </div--->

                <span id="saleTarget" class="just-normal-title">销售目标进度</span>
                <!--销售目标进度-->
                <el-row :gutter="16" style="margin-top: 10px" v-show="showSaleTarget">
                    <el-col :span="12">
                        <el-tooltip :content="`${yearProgress.current} / ${yearProgress.target}`" placement="right" effect="light"
                            :hide-after="1000">
                            <el-progress text-inside :percentage="yearProgress.percent" :stroke-width="30"
                                :color="progressStatus[Math.ceil(yearProgress.percent / 25) - 1]" :format="percentFormat"/>
                        </el-tooltip>
                    </el-col>
                    <el-col :span="12">
                        <el-tooltip :content="`${proYearProgress.current} / ${proYearProgress.target}`" placement="right" effect="light"
                            :hide-after="1000">
                            <el-progress text-inside :percentage="proYearProgress.percent" :stroke-width="30"
                                :color="progressStatus[Math.ceil(proYearProgress.percent / 25) - 1]" :format="percentFormat"/>
                        </el-tooltip>
                    </el-col>
                    <el-col :span="12" style="margin-top: 8px">
                        <el-tooltip :content="`${quarterProgress.current} / ${quarterProgress.target}`" placement="right" effect="light"
                            :hide-after="1000">
                            <el-progress text-inside :percentage="quarterProgress.percent" :stroke-width="30"
                                :color="progressStatus[Math.ceil(quarterProgress.percent / 25) - 1]" :format="percentFormat"/>
                        </el-tooltip>
                    </el-col>
                    <el-col :span="12" style="margin-top: 8px">
                        <el-tooltip :content="`${proQuarterProgress.current} / ${proQuarterProgress.target}`" placement="right" effect="light"
                            :hide-after="1000">
                            <el-progress text-inside :percentage="proQuarterProgress.percent" :stroke-width="30"
                                :color="progressStatus[Math.ceil(proQuarterProgress.percent / 25) - 1]" :format="percentFormat"/>
                        </el-tooltip>
                    </el-col>
                    <el-col :span="12" style="margin-top: 8px">
                        <el-tooltip :content="`${monthProgress.current} / ${monthProgress.target}`" placement="right" effect="light"
                            :hide-after="1000">
                            <el-progress text-inside :percentage="monthProgress.percent" :stroke-width="30"
                                :color="progressStatus[Math.ceil(monthProgress.percent / 25) - 1]" :format="percentFormat"/>
                        </el-tooltip>
                    </el-col>
                    <el-col :span="12" style="margin-top: 8px">
                        <el-tooltip :content="`${proMonthProgress.current} / ${proMonthProgress.target}`" placement="right" effect="light"
                            :hide-after="1000">
                            <el-progress text-inside :percentage="proMonthProgress.percent" :stroke-width="30"
                                :color="progressStatus[Math.ceil(proMonthProgress.percent / 25) - 1]" :format="percentFormat"/>
                        </el-tooltip>
                    </el-col>
                </el-row>
                <!--毛利指标进度-->
<!--                <el-row :gutter="6" style="margin-top: 10px" v-show="showProfit">-->
<!--                  <el-col :span="3">-->
<!--                    <el-button type="primary" circle size="mini">年</el-button>-->
<!--                  </el-col>-->
<!--                  <el-col :span="21">-->
<!--                    <el-tooltip :content="`${proYearProgress.current} / ${proYearProgress.target}`" placement="right" effect="light"-->
<!--                                :hide-after="1000">-->
<!--                      <el-progress text-inside :percentage="proYearProgress.percent" :stroke-width="30"-->
<!--                                   :color="progressStatus[Math.ceil(proYearProgress.percent / 25) - 1]" :format="percentFormat"/>-->
<!--                    </el-tooltip>-->
<!--                  </el-col>-->
<!--                  <el-col :span="3" style="margin-top: 8px">-->
<!--                    <el-button type="primary" circle size="mini">季</el-button>-->
<!--                  </el-col>-->
<!--                  <el-col :span="21" style="margin-top: 8px">-->
<!--                    <el-tooltip :content="`${proQuarterProgress.current} / ${proQuarterProgress.target}`" placement="right" effect="light"-->
<!--                                :hide-after="1000">-->
<!--                      <el-progress text-inside :percentage="proQuarterProgress.percent" :stroke-width="30"-->
<!--                                   :color="progressStatus[Math.ceil(proQuarterProgress.percent / 25) - 1]" :format="percentFormat"/>-->
<!--                    </el-tooltip>-->
<!--                  </el-col>-->
<!--                  <el-col :span="3" style="margin-top: 8px">-->
<!--                    <el-button type="primary" circle size="mini">月</el-button>-->
<!--                  </el-col>-->
<!--                  <el-col :span="21" style="margin-top: 8px">-->
<!--                    <el-tooltip :content="`${proMonthProgress.current} / ${proMonthProgress.target}`" placement="right" effect="light"-->
<!--                                :hide-after="1000">-->
<!--                      <el-progress text-inside :percentage="proMonthProgress.percent" :stroke-width="30"-->
<!--                                   :color="progressStatus[Math.ceil(proMonthProgress.percent / 25) - 1]" :format="percentFormat"/>-->
<!--                    </el-tooltip>-->
<!--                  </el-col>-->
<!--                </el-row>-->
            </grid-item>
            <grid-item
                :x="gridLayout[26].x"
                :y="gridLayout[26].y"
                :w="gridLayout[26].w"
                :h="gridLayout[26].h"
                :i="gridLayout[26].i"
                :key="gridLayout[26].i"
                v-loading="gridLayout[26].loading">
                <!--div class="info-grid" v-if="area.mom">
                    <p class="info-title">
                        加油站类型
                    </p>
                    <p class="info-desc">
                        <font-awesome-icon icon="gas-pump" />
                        {{ stationInfo.ztType }}
                    </p>
                </div-->
                <p class="just-normal-title">销售评级</p>
                <p style="margin-top: 10px; margin-bottom: 8px; padding: 4px 6px; color: #FFFFFF; font-size: 18px; font-weight: bold; border-radius: 5px; background-color: #3c8dbc; border-color: #3c8dbc; text-align: center">
                  {{ saleZone.zoneName }}
                  <font-awesome-icon v-if="this.num1 < this.num2" icon="arrow-up" color="#67C23A"></font-awesome-icon>
                  <font-awesome-icon v-else  icon="arrow-down" color="#F56C6C"></font-awesome-icon>
                </p>
                <p style="margin-bottom: 8px; padding: 4px 6px; color: #FFFFFF; font-size: 12px; font-weight: bold; border-radius: 5px; background-color: #E6A23C; border-color: #E6A23C; text-align: center">
                  {{ lastSaleZone.zoneName }}
                </p>
                <el-row :gutter="6">
                    <el-col :span="12" style="text-align: right">
                        <el-tag effect="dark" type="success">{{ saleZone.carServSR ? (saleZone.carServLR ? '有汽服有收入' : '有汽服无收入') : (saleZone.carServLR ? '有汽服无收入' : '无汽服无收入') }}</el-tag>
                    </el-col>
                    <el-col :span="12" style="text-align: left">
                        <el-tag effect="dark" type="success">营业天数:{{ saleZone.saleDays }}</el-tag>
                    </el-col>
                </el-row>
            </grid-item>
            <grid-item
                :x="gridLayout[27].x"
                :y="gridLayout[27].y"
                :w="gridLayout[27].w"
                :h="gridLayout[27].h"
                :i="gridLayout[27].i"
                :key="gridLayout[27].i"
                v-loading="gridLayout[27].loading">
                <div v-if="showPaySel" class="single-cat-panel">
                    <el-button size="mini" style="width: 100%;" type="primary"
                        @click="updatePayType(0)">按金额</el-button>
                    <el-button size="mini" style="width: 100%; margin: 6px 0 0 0" type="primary"
                        @click="updatePayType(1)">按数量</el-button>
                </div>
                <i class="ri-arrow-down-s-line cat-select" @click="showPaySel = !showPaySel"></i>
                <v-chart theme="macarons" :options="payPieOptions" :autoresize="true"/>
            </grid-item>
            <!--grid-item
                :x="gridLayout[28].x"
                :y="gridLayout[28].y"
                :w="gridLayout[28].w"
                :h="gridLayout[28].h"
                :i="gridLayout[28].i"
                :key="gridLayout[28].i"
                v-loading="gridLayout[28].loading">
                <div class="info-grid" v-if="area.mom">
                    <p class="info-title">
                        便利店面积
                    </p>
                    <p class="info-desc">
                        <font-awesome-icon icon="ruler-combined" />
                        {{ stationInfo.ztxtBld }}
                    </p>
                </div>
            </grid-item>
            <grid-item
                :x="gridLayout[29].x"
                :y="gridLayout[29].y"
                :w="gridLayout[29].w"
                :h="gridLayout[29].h"
                :i="gridLayout[29].i"
                :key="gridLayout[29].i"
                v-loading="gridLayout[29].loading">
                <div class="info-grid" v-if="area.mom">
                    <p class="info-title">
                        油站编号
                    </p>
                    <p class="info-desc">
                        <font-awesome-icon icon="code" />
                        {{ stationInfo.plant }}
                    </p>
                </div>
            </grid-item>
            <grid-item--
                :x="gridLayout[30].x"
                :y="gridLayout[30].y"
                :w="gridLayout[30].w"
                :h="gridLayout[30].h"
                :i="gridLayout[30].i"
                :key="gridLayout[30].i"
                v-loading="gridLayout[30].loading">
                <div class="info-grid" v-if="area.mom">
                    <p class="info-title">
                        资产性质
                    </p>
                    <p class="info-desc">
                        <font-awesome-icon icon="building" />
                        {{ stationInfo.ztxtZcxz }}
                    </p>
                </div>
            </grid-item-->
        </grid-layout>
    </div>
</template>

<style lang="scss">
.cnpc-single-store-main {
    background-image: url('../../assets/images/single-bg.jpg');
    background-position-x: center;
    background-position-y: top;
    width: 100%;
    height: calc(100vh - 12px);
    overflow: auto;

    .just-normal-title, .just-normal-subTitle {
        text-align: left;
        color:rgba(255, 255, 255, 1);
        font-size: 18px;
        padding-left: 12px;
    }

    .just-normal-title::before {
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
    .just-normal-subTitle::before{
        content: "";
        width: 5px;
        height: 20px;
        border-radius: 18px;
        background-color: #409EFF;
        display: block;
        position: absolute;
        top: 11px;
        left: 160px;
    }
    .just-text {
        color: #FFFFFF;
        text-align: center;
        margin-top: 8px;
        font-size: 18px;
    }

    .echarts {
        width: 100%;
        height: 100%;
    }

    .tippy-popper {
        max-width: 600px;
    }

    .single-level-panel {
        position: absolute;
        width: 250PX;
        padding: 10px 10px 10px 10px;
        border-radius: 5px;
        background-color: #FFF;
        top: 6px;
        left: 80px;
        z-index: 999;
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

    .single-store-header {
        width: 100%;
        height: 60px;
        position: relative;

        .title {
            font-size: 1.5em;
            width: 100%;
            text-align: center;
            color: #FFF;
        }

        .sub-title {
            color: #FFF;
            text-align: center;
            font-size: 12px;
        }

        span.timer {
            color: rgba(255, 255, 255, 0.8);
            position: absolute;
            top: 15px;
            right: 20px;
        }

        .single-full-screen {
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

        .single-level-select-fa {
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
    }

    .single-store-grid {
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
                padding: 25px 0 0 0;
                text-align: center;
            }

            .sale-amount-normal2 {
                font-family: 'Manrope';
                color: #E6A23C;
                font-size: 20px;
                font-weight: bold;
                padding: 25px 0 0 0;
                text-align: center;
            }
        }

        .info-grid {
            .info-title {
                color:rgba(255, 255, 255, 1);
                font-size: 13px;
                padding-left: 12px;
            }

            .info-title::before {
                content: "";
                width: 5px;
                height: 13px;
                border-radius: 15px;
                background-color: #409EFF;
                display: block;
                position: absolute;
                top: 10px;
                left: 8px;
            }

            .info-desc {
                color: rgb(19, 204, 173);
                margin-top: 8px;
                font-size: 17px;
                font-weight: bold;
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

    .vue-grid-item {
        background-color: rgba(255, 255, 255, 0.02);
        border-radius: 4px;
        padding: 8px 8px 8px 8px;
    }

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

    .prev-select {
      position: absolute;
      top: calc(100% /2 - 2px);
      left: 2px;
      color: rgba(255, 255, 255, 0.1);
      font-size: 15px;
      cursor: pointer;
      z-index: 999;

      &:hover {
        color: rgba(255, 255, 255, 0.9);
      }
    }

    .next-select {
      position: absolute;
      top: calc(100% /2 - 2px);
      left: calc(100% - 17px);
      color: rgba(255, 255, 255, 0.1);
      font-size: 15px;
      cursor: pointer;
      z-index: 999;

      &:hover {
        color: rgba(255, 255, 255, 0.9);
      }
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
import { Loading } from 'element-ui';
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
        'v-chart': ECharts
    },
    beforeDestroy() {
        clearInterval(this.timer)
    },
    mounted() {
        self = this
        this.getStations()

        let earlier = new Date()
        earlier.setDate(earlier.getDate() - 30)
        this.dateRange.push(moment(earlier).format('YYYY-MM-DD'))
        this.dateRange.push(moment(Date.now()).format('YYYY-MM-DD'))
        this.sdate = this.dateRange[0]
        this.edate = this.dateRange[1]

        this.timer = setInterval(this.singleTimer, 1000)
        this.fitties.push(fitty('.sale-amount-number', {
            maxSize: 50
        }));
        this.fitties.push(fitty('.sale-amount-normal', {
            maxSize: 40
        }));
        this.fitties.push(fitty('.sale-amount-normal2', {
            maxSize: 40
        }));
        this.fitties.push(fitty('.single-num-normal', {
            maxSize: 50
        }));

        let ele = document.getElementById('donut-auto-resize-transfer');
        this.onDonutResize(ele);
    },
    data() {
        return {
            showPaySel: false,
            showSaleTarget: true,
            showProfit: false,
            showCatPie: true,
            showPayPie: false,
            progressStatus: [
                '#F56C6C',
                '#E6A23C',
                '#409EFF',
                '#67C23A',
            ],
            gridRowHeight: 72,
            currentTimeDesc: '',
            fullScreen: false,
            ticker: 0,
            fitties: [],
            dateRange: [],
            sdate: '',
            edate: '',
            // 去年的起始和结束时间
            ysdate: '',
            yedate: '',
            sales: {
                value: 0
            },

            topGoodDetail: false,
            ropGoodDetail: false,
            detailTopGood: {},
            detailRopGood: {},
            topGoods: [],
            ropGoods: [],

            transfer: {
                value: 0
            },
            donutWidth: 60,

            profit: {
                value: 0
            },

            margin: {
                value: 0
            },

            average: {
                value: 0
            },

            area: {
                value: 0
            },

            customers: {
                value: 0
            },

            percentColors: [
                {color: '#f56c6c', percentage: 20},
                {color: '#e6a23c', percentage: 40},
                {color: '#5cb87a', percentage: 60},
                {color: '#1989fa', percentage: 80},
                {color: '#6f7ad3', percentage: 100}
            ],

            catIncomeData: [],
            catProfitData: [],
            catMarginData: [],
            catType: 'income',
            catDatasource: null,
            showCatSel: false,
            catPieOptions: {
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
                            self.catDatasource.forEach(e => {
                                if (e.name == args.name) {
                                    percent = (e.value * 100 / self.catDatasource.total).toFixed(2)
                                }
                            })
                            return `<p>${args.name}: ${percent}%</p>`
                        },
                    },
                },
                series: [
                    {
                        name: '销售比例(销售收入)',
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
            // 支付方式配置
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
            trendLineOptions: {
                title: {
                    text: '销售趋势',
                    textStyle: {
                        fontSize: 18,
                        color: 'rgba(255, 255, 255, 0.7)'
                    }
                },
                tooltip: {
                    trigger: 'axis',
                    formatter: function(args) {
                        let tooltip = `<p>${args[0].axisValue}</p> `;
                        if (args[0])
                            tooltip += `<p>${args[0].marker} ${args[0].seriesName}:
                                ${beautifyNumber(args[0].value.toFixed(2))}</p>`;
                        if (args[1])
                            tooltip += `<p>${args[1].marker} ${args[1].seriesName}:
                                ${beautifyNumber(args[1].value.toFixed(2))}</p>`;
                        if (args[2])
                            tooltip += `<p>${args[2].marker} ${args[2].seriesName}:
                                ${args[2].value.toFixed(2)}%</p>`;
                        return tooltip;
                    },
                    axisPointer: {
                        type: 'cross',
                        label: {
                            backgroundColor: '#6a7985'
                        }
                    }
                },
                legend: {
                    data: ['销售收入', '销售毛利', '销售毛利率'],
                    icon: 'roundRect',
                    textStyle: {
                        color: '#FFF'
                    }
                },
                toolbox: {
                    feature: {
                        saveAsImage: {
                            name: '销售趋势'
                        }
                    }
                },
                dataZoom: [
                    {
                        type: 'inside',
                    }, {
                        start: 0,
                        bottom: 0,
                        height: 15,
                        textStyle: {
                            color: '#FFF'
                        }
                    }
                ],
                grid: {
                    left: '3%',
                    right: '4%',
                    bottom: 20,
                    containLabel: true,
                },
                xAxis: [
                    {
                        type: 'category',
                        boundaryGap: false,
                        data: [],
                        splitLine:{//remove grid lines
                            show:false
                        },
                    }
                ],
                yAxis: [
                    {
                        type: 'value',
                        name: '收入',
                        splitArea: {
                            "show": false
                        },
                        splitLine:{//remove grid lines
                            show:false
                        },
                    }, {
                        type: 'value',
                        name: '毛利',
                        position: 'left',
                        offset: 60,
                        splitArea: {
                            "show": false
                        },
                        splitLine:{//remove grid lines
                            show:false
                        },
                    }, {
                        type: 'value',
                        name: '毛利率',
                        splitArea: {
                            "show": false
                        },
                        splitLine:{//remove grid lines
                            show:false
                        },
                    }
                ],
                series: [
                    {
                        name: '销售收入',
                        type: 'line',
                        lineStyle: {
                            // color: '#D87A80',
                            color: '#ff7c7c',
                            width: 3
                        },
                        itemStyle: {
                            // color: '#D87A80'
                            color: '#ff7c7c'
                        },
                        data: []
                    },
                    {
                        name: '销售毛利',
                        type: 'line',
                        symbol: 'circle',
                        symbolSize: 7,
                        lineStyle: {
                            // color: '#FFB980',
                            color: '#63c2ff',
                            width: 3,
                        },
                        itemStyle: {
                            // color: '#FFB980'
                            color: '#63c2ff',
                            borderWidth: 2,
                            borderColor: '#c8b2f4',
                        },
                        yAxisIndex: 1,
                        data: []
                    },
                    {
                        name: '销售毛利率',
                        type: 'line',
                        symbol: 'triangle',
                        symbolSize: 3,
                        lineStyle: {
                            // color: '#5AB1EF',
                            color: '#5bc49f',
                            width: 3,
                            type: 'dashed'
                        },
                        itemStyle: {
                            // color: '#5AB1EF'
                            color: '#5bc49f',
                            borderWidth: 3,
                            borderColor: '#f8cb7f',
                        },
                        yAxisIndex: 2,
                        data: []
                    }
                ]
            },

            gridLayout: [
                // top10 goods
                {"x":0,"y":2,"w":3,"h":3,"i":"single-store-0", loading: false},

                // sale income
                {"x":5,"y":2,"w":3,"h":2,"i":"single-store-1", loading: false},
                {"x":8,"y":2,"w":1,"h":1,"i":"single-store-2", loading: false},
                {"x":8,"y":3,"w":1,"h":1,"i":"single-store-3", loading: false},

                // rop10 goods
                {"x":0,"y":5,"w":3,"h":3,"i":"single-store-4", loading: false},

                // sale trend chart
                {"x":3,"y":4,"w":9,"h":4,"i":"single-store-5", loading: false},

                // transfer
                {"x":3,"y":0,"w":2,"h":2,"i":"single-store-6", loading: false},
                {"x":5,"y":0,"w":1,"h":1,"i":"single-store-7", loading: false},
                {"x":5,"y":1,"w":1,"h":1,"i":"single-store-8", loading: false},

                // categories
                {"x":9,"y":2,"w":3,"h":2,"i":"single-store-9", loading: false},

                // profit
                {"x":6,"y":8,"w":2,"h":2,"i":"single-store-10", loading: false},
                {"x":8,"y":8,"w":1,"h":1,"i":"single-store-11", loading: false},
                {"x":8,"y":9,"w":1,"h":1,"i":"single-store-12", loading: false},

                // margin
                {"x":3,"y":8,"w":2,"h":2,"i":"single-store-13", loading: false},
                {"x":5,"y":8,"w":1,"h":1,"i":"single-store-14", loading: false},
                {"x":5,"y":9,"w":1,"h":1,"i":"single-store-15", loading: false},

                // average
                {"x":9,"y":8,"w":2,"h":2,"i":"single-store-16", loading: false},
                {"x":11,"y":8,"w":1,"h":1,"i":"single-store-17", loading: false},
                {"x":11,"y":9,"w":1,"h":1,"i":"single-store-18", loading: false},

                // customers
                {"x":6,"y":0,"w":2,"h":2,"i":"single-store-19", loading: false},
                {"x":8,"y":0,"w":1,"h":1,"i":"single-store-20", loading: false},
                {"x":8,"y":1,"w":1,"h":1,"i":"single-store-21", loading: false},

                // area
                {"x":9,"y":0,"w":2,"h":2,"i":"single-store-22", loading: false},
                {"x":11,"y":0,"w":1,"h":1,"i":"single-store-23", loading: false},
                {"x":11,"y":1,"w":1,"h":1,"i":"single-store-24", loading: false},

                // station info --> progress
                {"x":0,"y":0,"w":3,"h":2,"i":"single-store-25", loading: false},
                {"x":3,"y":3,"w":2,"h":2,"i":"single-store-26", loading: false},

                // payments
                {"x":0,"y":8,"w":3,"h":2,"i":"single-store-27", loading: false},
                // {"x":0,"y":1,"w":1,"h":1,"i":"single-store-27", loading: false},
                // {"x":1,"y":1,"w":1,"h":1,"i":"single-store-28", loading: false},
                // {"x":2,"y":0,"w":1,"h":1,"i":"single-store-29", loading: false},
                // {"x":2,"y":1,"w":1,"h":1,"i":"single-store-30", loading: false},
            ],

            stations: [],
            selectedStationCode: '',
            stationInfo: {},
            loadingInst: null,

            timer: null,

            showPanel: false,
            monthProgress: {
                current: 0,
                target: 0,
                percent: 0
            },
            quarterProgress: {
                current: 0,
                target: 0,
                percent: 0
            },
            yearProgress: {
                current: 0,
                target: 0,
                percent: 0
            },
            // 毛利指标进度
            proMonthProgress: {
              current: 0,
              target: 0,
              percent: 0
            },
            proQuarterProgress: {
              current: 0,
              target: 0,
              percent: 0
            },
            proYearProgress: {
              current: 0,
              target: 0,
              percent: 0
            },
            saleZone: {

            },
            lastSaleZone: {},
            // 用于计算评级同比 num2:今年的销售  num1:去年的销售
            num1: 0,
            num2: 0
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
        singleTimer() {
            this.currentTimeDesc = moment(Date.now()).format('YYYY/MM/DD HH:mm:ss');
            if (!document.fullscreenElement)
                this.fullScreen = false;

            this.ticker += 1;
            if (this.ticker === 3000) {
                this.ticker = 0;
                this.getAllData();
            }
        },
        getStations() {
            doRequest({
                method: 'GET',
                url: `/v1/web/user/org-perm/${getUserInfo().name}`,
                loading: true
            }, {
                success: res => {
                    let perms = []
                    if (res.orgPerm.length > 0) {
                        perms = JSON.parse(res.orgPerm)
                    }

                    doRequest({
                        method: 'GET',
                        url: '/v1/web/org/plant',
                        loading: true
                    }, {
                        success: res => {
                            let orgs = []
                            res.orgList.forEach(e => {
                                orgs.push(e)
                            })

                            if (perms.length > 0) {
                                orgs.forEach(e => {
                                    perms.forEach(p => {
                                        if (p.orgCode == e.orgCode) {
                                            if (p.show || (p.show == null))
                                                this.stations.push(e)
                                        }
                                    })
                                })
                            } else {
                                this.stations = orgs
                            }
                        },
                        fail: _ => {
                            message('error', '获取油站列表失败，请稍后再试')
                        }
                    })
                },
                fail: _ => {
                    message('error', '获取用户机构权限失败, 请稍后再试')
                }
            })
        },
        getStationInfo() {
            if (!this.selectedStationCode) {
                message('warning', '请选择有效的油站')
                return
            }

            // this.gridLayout[25].loading = true
            // this.gridLayout[26].loading = true
            // this.gridLayout[27].loading = true
            // this.gridLayout[28].loading = true
            // this.gridLayout[29].loading = true
            // this.gridLayout[30].loading = true

            doRequest({
                url: `/v1/web/basic/plant/${this.selectedStationCode}`
            }, {
                obj: this,
                src: 'plant',
                dst: 'stationInfo',
                fail: _ => {
                    message('error', '获取站点基本信息失败')
                },
                finally: _ => {
                    // this.gridLayout[25].loading = false
                    // this.gridLayout[26].loading = false
                    // this.gridLayout[27].loading = false
                    // this.gridLayout[28].loading = false
                    // this.gridLayout[29].loading = false
                    // this.gridLayout[30].loading = false
                }
            })
        },
        getSaleZone() {
            if (!this.selectedStationCode || !this.dateRange || !this.dateRange[0] || !this.dateRange[1]) {
                message('warning', '请选择有效的油站与时间')
                return
            }

            this.gridLayout[26].loading = true
            let ysdate = new Date(this.sdate)
            ysdate.setFullYear(ysdate.getFullYear() - 1)
            let yedate = new Date(this.edate)
            yedate.setFullYear(yedate.getFullYear() - 1)

            let posts = []
            // 今年的销售数据
            posts[0] = {
              url: '/v1/web/plant/single-sales-zone',
              method: 'POST',
              data: {
                beginDate: this.sdate,
                endDate: this.edate,
                orgLevel: 2,
                orgCode: this.selectedStationCode
              }
            }
            posts[1] = {
              url: '/v1/web/plant/single-sales-zone',
              method: 'POST',
              data: {
                beginDate: moment(ysdate).format('YYYY-MM-DD'),
                endDate: moment(yedate).format('YYYY-MM-DD'),
                orgLevel: 2,
                orgCode: this.selectedStationCode
              }
            }
            if (posts.length > 0){
              allRequests({
                success: res=> {
                  // 今年的销售评级
                  this.saleZone = res[0].data.psz
                  this.num2 = res[0].data.psz.sales
                  // 去年的销售评级
                  this.lastSaleZone = res[1].data.psz
                  this.num1 = res[1].data.psz.sales
                },
                fail: _ => {
                  message('error', '获取销售评级同比失败，请稍后再试')
                },
                finally: _ => {
                  this.gridLayout[26].loading = false;
                }
              }, [], posts)
            }
        },
        // 获得剔除5%毛利率的销售评级
        getThresholdSaleZone() {
        if (!this.selectedStationCode || !this.dateRange || !this.dateRange[0] || !this.dateRange[1]) {
          message('warning', '请选择有效的油站与时间')
          return
        }

        this.gridLayout[26].loading = true
        let ysdate = new Date(this.sdate)
        ysdate.setFullYear(ysdate.getFullYear() - 1)
        let yedate = new Date(this.edate)
        yedate.setFullYear(yedate.getFullYear() - 1)

        let posts = []
        // 今年的销售数据
        posts[0] = {
          url: '/v1/web/plant/single-sales-zone',
          method: 'POST',
          data: {
            beginDate: this.sdate,
            endDate: this.edate,
            orgLevel: 2,
            orgCode: this.selectedStationCode,
            rateThreshold: 0.05
          }
        }
        posts[1] = {
          url: '/v1/web/plant/single-sales-zone',
          method: 'POST',
          data: {
            beginDate: moment(ysdate).format('YYYY-MM-DD'),
            endDate: moment(yedate).format('YYYY-MM-DD'),
            orgLevel: 2,
            orgCode: this.selectedStationCode,
            rateThreshold: 0.05
          }
        }
        if (posts.length > 0){
          allRequests({
            success: res=> {
              // 今年的销售评级
              this.saleZone = res[0].data.psz
              this.num2 = res[0].data.psz.sales
              // 去年的销售评级
              this.lastSaleZone = res[1].data.psz
              this.num1 = res[1].data.psz.sales
            },
            fail: _ => {
              message('error', '获取销售评级同比失败，请稍后再试')
            },
            finally: _ => {
              this.gridLayout[26].loading = false;
            }
          }, [], posts)
        }
      },
        getKpi() {
            if (!this.selectedStationCode || !this.dateRange || !this.dateRange[0] || !this.dateRange[1]) {
                message('warning', '请选择有效的油站与时间')
                return
            }

            this.gridLayout[1].loading = true
            this.gridLayout[2].loading = true
            this.gridLayout[3].loading = true
            this.gridLayout[6].loading = true
            this.gridLayout[7].loading = true
            this.gridLayout[8].loading = true
            for (let i = 10; i <= 25; i ++)
                this.gridLayout[i].loading = true

            doRequest({
                url: '/v1/web/sales/plant/kpi',
                method: 'POST',
                data: {
                    beginDate: this.sdate,
                    endDate: this.edate,
                    orgLevel: 2,
                    orgCode: this.selectedStationCode
                }
            }, {
                success: res => {
                    this.sales = res.kpi.netIncome
                    this.profit = res.kpi.grossProfit
                    this.margin = res.kpi.grossMargin
                    this.customers = res.kpi.nonFuelCount
                    this.transfer = res.kpi.fnConversionRate
                    this.average = res.kpi.avgTrxValue
                    this.area = res.kpi.salesPerM2
                    this.y2nIncome = res.kpi.y2nIncome
                    this.q2nIncome = res.kpi.q2nIncome
                    this.m2nIncome = res.kpi.m2nIncome
                    this.getKpiByMonth()
                },
                fail: _ => {
                    message('error', '获取 KPI 数据失败')
                },
                finally: _ => {
                    this.gridLayout[1].loading = false
                    this.gridLayout[2].loading = false
                    this.gridLayout[3].loading = false
                    this.gridLayout[6].loading = false
                    this.gridLayout[7].loading = false
                    this.gridLayout[8].loading = false
                    for (let i = 10; i <= 24; i ++)
                        this.gridLayout[i].loading = false
                }
            })
        },
        // 获得剔除5%毛利率的kpi
        getThresholdKpi() {
          this.gridLayout[1].loading = true
          this.gridLayout[2].loading = true
          this.gridLayout[3].loading = true
          this.gridLayout[6].loading = true
          this.gridLayout[7].loading = true
          this.gridLayout[8].loading = true
          for (let i = 10; i <= 25; i ++)
            this.gridLayout[i].loading = true

          doRequest({
            url: '/v1/web/sales/plant/kpi',
            method: 'POST',
            data: {
              beginDate: this.sdate,
              endDate: this.edate,
              orgLevel: 2,
              orgCode: this.selectedStationCode,
              rateThreshold: 0.05
            }
          }, {
            success: res => {
              console.log(res)
              this.sales = res.kpi.netIncome
              this.profit = res.kpi.grossProfit
              this.margin = res.kpi.grossMargin
              this.customers = res.kpi.nonFuelCount
              this.transfer = res.kpi.fnConversionRate
              this.average = res.kpi.avgTrxValue
              this.area = res.kpi.salesPerM2
              this.y2nIncome = res.kpi.y2nIncome
              this.q2nIncome = res.kpi.q2nIncome
              this.m2nIncome = res.kpi.m2nIncome
              this.getKpiByMonth()
            },
            fail: _ => {
              message('error', '获取 KPI 数据失败')
            },
            finally: _ => {
              this.gridLayout[1].loading = false
              this.gridLayout[2].loading = false
              this.gridLayout[3].loading = false
              this.gridLayout[6].loading = false
              this.gridLayout[7].loading = false
              this.gridLayout[8].loading = false
              for (let i = 10; i <= 24; i ++)
                this.gridLayout[i].loading = false
            }
          })
        },
        getTopGoods() {
            if (!this.selectedStationCode || !this.dateRange || !this.dateRange[0] || !this.dateRange[1]) {
                message('warning', '请选择有效的油站与时间')
                return
            }

            this.gridLayout[0].loading = true
            doRequest({
                url: '/v1/web/sales/plant/material/rank',
                method: 'POST',
                data: {
                    beginDate: this.sdate,
                    endDate: this.edate,
                    orgLevel: 2,
                    orgCode: this.selectedStationCode,
                    sortBy: 'DESC',
                    limit: 10
                }
            }, {
                success: res => {
                    this.topGoods = res.matlList.slice(0, 6)
                },
                fail: _ => {
                    message('error', '获取畅销商品列表失败')
                },
                finally: _ => {
                    this.gridLayout[0].loading = false
                }
            })
        },
        // 获得剔除5%毛利率的畅销商品
        getThresholdTopGoods() {
          this.gridLayout[0].loading = true
          doRequest({
            url: '/v1/web/sales/plant/material/rank',
            method: 'POST',
            data: {
              beginDate: this.sdate,
              endDate: this.edate,
              orgLevel: 2,
              orgCode: this.selectedStationCode,
              rateThreshold: 0.05,
              sortBy: 'DESC',
              limit: 10
            }
          }, {
            success: res => {
              this.topGoods = res.matlList.slice(0, 6)
            },
            fail: _ => {
              message('error', '获取畅销商品列表失败')
            },
            finally: _ => {
              this.gridLayout[0].loading = false
            }
          })
        },
        getRopGoods() {
            if (!this.selectedStationCode || !this.dateRange || !this.dateRange[0] || !this.dateRange[1]) {
                message('warning', '请选择有效的油站与时间')
                return
            }

            this.gridLayout[4].loading = true
            doRequest({
                url: '/v1/web/sales/plant/material/rank',
                method: 'POST',
                data: {
                    beginDate: this.sdate,
                    endDate: this.edate,
                    orgLevel: 2,
                    orgCode: this.selectedStationCode,
                    sortBy: 'ASC',
                    limit: 10
                }
            }, {
                success: res => {
                    this.ropGoods = res.matlList.slice(0, 6)
                },
                fail: _ => {
                    message('error', '获取滞销商品列表失败')
                },
                finally: _ => {
                    this.gridLayout[4].loading = false
                }
            })
        },
        // 获得剔除5%毛利率的滞销商品
        getThresholdRopGoods() {
          this.gridLayout[4].loading = true
          doRequest({
            url: '/v1/web/sales/plant/material/rank',
            method: 'POST',
            data: {
              beginDate: this.sdate,
              endDate: this.edate,
              orgLevel: 2,
              orgCode: this.selectedStationCode,
              rateThreshold: 0.05,
              sortBy: 'ASC',
              limit: 10
            }
          }, {
            success: res => {
              this.ropGoods = res.matlList.slice(0, 6)
            },
            fail: _ => {
              message('error', '获取滞销商品列表失败')
            },
            finally: _ => {
              this.gridLayout[4].loading = false
            }
          })
        },
        getTrend() {
            if (!this.selectedStationCode || !this.sdate || !this.edate) {
                message('warning', '请选择有效的油站与时间')
                return
            }

            this.gridLayout[5].loading = true
            doRequest({
                url: '/v1/web/sales/plant/date',
                method: 'POST',
                data: {
                    beginDate: this.sdate,
                    endDate: this.edate,
                    orgLevel: 2,
                    orgCode: this.selectedStationCode
                }
            }, {
                success: res => {
                    this.trendLineOptions.xAxis[0].data = []
                    this.trendLineOptions.series[0].data = []
                    this.trendLineOptions.series[1].data = []
                    this.trendLineOptions.series[2].data = []
                    res.dateList.forEach(e => {
                        this.trendLineOptions.xAxis[0].data.push(e.Date)
                        this.trendLineOptions.series[0].data.push(e.Metric.NetvalInv)
                        this.trendLineOptions.series[1].data.push(e.Metric.GrossProfit)
                        this.trendLineOptions.series[2].data.push(e.Metric.GrossMargin)
                    })
                },
                fail: _ => {
                    message('error', '获取销售趋势失败')
                },
                finally: _ => {
                    this.gridLayout[5].loading = false
                }
            })
        },
        // 获取剔除5%毛利率的销售趋势
        getThresholdTrend() {
          this.gridLayout[5].loading = true
          doRequest({
            url: '/v1/web/sales/plant/date',
            method: 'POST',
            data: {
              beginDate: this.sdate,
              endDate: this.edate,
              orgLevel: 2,
              orgCode: this.selectedStationCode,
              rateThreshold: 0.05
            }
          }, {
            success: res => {
              this.trendLineOptions.xAxis[0].data = []
              this.trendLineOptions.series[0].data = []
              this.trendLineOptions.series[1].data = []
              this.trendLineOptions.series[2].data = []
              res.dateList.forEach(e => {
                this.trendLineOptions.xAxis[0].data.push(e.Date)
                this.trendLineOptions.series[0].data.push(e.Metric.NetvalInv)
                this.trendLineOptions.series[1].data.push(e.Metric.GrossProfit)
                this.trendLineOptions.series[2].data.push(e.Metric.GrossMargin)
              })
            },
            fail: _ => {
              message('error', '获取销售趋势失败')
            },
            finally: _ => {
              this.gridLayout[5].loading = false
            }
          })
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
        getPayments() {
            this.gridLayout[20].loading = true

            doRequest({
                url: '/v1/web/kpi/payments',
                method: 'POST',
                data: {
                    beginDate: this.sdate,
                    endDate: this.edate,
                    orgCode: this.selectedStationCode
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
                    this.gridLayout[20].loading = false
                }
            })
        },
        updateCatType(catType) {
            switch (catType) {
                case 'income':
                    if (this.catIncomeData.length > 0) {
                        this.catPieOptions.legend.data = this.catIncomeData.legend
                        this.catPieOptions.series[0].name = '销售比例(销售收入)'
                        this.catPieOptions.series[0].data = this.catIncomeData
                    }
                    this.catDatasource = this.catIncomeData
                    break

                case 'margin':
                    if (this.catMarginData.length > 0){
                        this.catPieOptions.legend.data = this.catMarginData.legend
                        this.catPieOptions.series[0].name = '销售比例(毛利率)'
                        this.catPieOptions.series[0].data = this.catMarginData
                    }
                    this.catDatasource = this.catMarginData
                    break

                case 'profit':
                    if (this.catProfitData.length > 0){
                        this.catPieOptions.legend.data = this.catProfitData.legend
                        this.catPieOptions.series[0].name = '销售比例(毛利)'
                        this.catPieOptions.series[0].data = this.catProfitData
                    }
                    this.catDatasource = this.catProfitData
                    break
            }
        },
        getCategories() {
            if (!this.selectedStationCode || !this.sdate || !this.edate) {
                message('warning', '请选择有效的油站与时间')
                return
            }

            this.catType = 'income'
            this.gridLayout[9].loading = true
            doRequest({
                url: '/v1/web/sales/plant/class',
                method: 'POST',
                data: {
                    beginDate: this.sdate,
                    endDate: this.edate,
                    orgLevel: 2,
                    orgCode: this.selectedStationCode
                }
            }, {
                success: res => {
                    this.catPieOptions.legend.data = []
                    this.catPieOptions.series[0].data = []
                    this.catIncomeData = []
                    this.catProfitData = []
                    this.catMarginData = []
                    let s1 = 0, s2 = 0, s3 = 0
                    res.classList.forEach(e => {
                        this.catIncomeData.push({
                            name: e.ClassText,
                            value: e.Metric.NetvalInv
                        })
                        this.catMarginData.push({
                            name: e.ClassText,
                            value: e.Metric.GrossMargin
                        })
                        this.catProfitData.push({
                            name: e.ClassText,
                            value: e.Metric.GrossProfit
                        })
                        s1 += e.Metric.NetvalInv
                        s2 += e.Metric.GrossMargin
                        s3 += e.Metric.GrossProfit
                    })
                    this.catIncomeData.total = s1
                    this.catMarginData.total = s2
                    this.catProfitData.total = s3
                    this.catIncomeData.legend = []
                    this.catMarginData.legend = []
                    this.catProfitData.legend = []
                    this.catIncomeData.sort((a, b) => {
                        return b.value - a.value
                    })
                    this.catMarginData.sort((a, b) => {
                        return b.value - a.value
                    })
                    this.catProfitData.sort((a, b) => {
                        return b.value - a.value
                    })
                    this.catIncomeData.forEach(e => {
                        this.catIncomeData.legend.push(e.name)
                    })
                    this.catMarginData.forEach(e => {
                        this.catMarginData.legend.push(e.name)
                    })
                    this.catProfitData.forEach(e => {
                        this.catProfitData.legend.push(e.name)
                    })
                    this.catDatasource = this.catIncomeData
                    this.catPieOptions.series[0].data = this.catIncomeData
                    this.catPieOptions.legend.data = this.catIncomeData.legend
                },
                fail: _ => {
                    message('error', '获取销售分类失败')
                },
                finally: _ => {
                    this.gridLayout[9].loading = false
                }
            })
        },
        // 获取剔除5%毛利率的大类比例数据
        getThresholdCategories() {
          this.catType = 'income'
          this.gridLayout[9].loading = true
          doRequest({
            url: '/v1/web/sales/plant/class',
            method: 'POST',
            data: {
              beginDate: this.sdate,
              endDate: this.edate,
              orgLevel: 2,
              orgCode: this.selectedStationCode,
              rateThreshold: 0.05
            }
          }, {
            success: res => {
              this.catPieOptions.legend.data = []
              this.catPieOptions.series[0].data = []
              this.catIncomeData = []
              this.catProfitData = []
              this.catMarginData = []
              let s1 = 0, s2 = 0, s3 = 0
              res.classList.forEach(e => {
                this.catIncomeData.push({
                  name: e.ClassText,
                  value: e.Metric.NetvalInv
                })
                this.catMarginData.push({
                  name: e.ClassText,
                  value: e.Metric.GrossMargin
                })
                this.catProfitData.push({
                  name: e.ClassText,
                  value: e.Metric.GrossProfit
                })
                s1 += e.Metric.NetvalInv
                s2 += e.Metric.GrossMargin
                s3 += e.Metric.GrossProfit
              })
              this.catIncomeData.total = s1
              this.catMarginData.total = s2
              this.catProfitData.total = s3
              this.catIncomeData.legend = []
              this.catMarginData.legend = []
              this.catProfitData.legend = []
              this.catIncomeData.sort((a, b) => {
                return b.value - a.value
              })
              this.catMarginData.sort((a, b) => {
                return b.value - a.value
              })
              this.catProfitData.sort((a, b) => {
                return b.value - a.value
              })
              this.catIncomeData.forEach(e => {
                this.catIncomeData.legend.push(e.name)
              })
              this.catMarginData.forEach(e => {
                this.catMarginData.legend.push(e.name)
              })
              this.catProfitData.forEach(e => {
                this.catProfitData.legend.push(e.name)
              })
              this.catDatasource = this.catIncomeData
              this.catPieOptions.series[0].data = this.catIncomeData
              this.catPieOptions.legend.data = this.catIncomeData.legend
            },
            fail: _ => {
              message('error', '获取销售分类失败')
            },
            finally: _ => {
              this.gridLayout[9].loading = false
            }
          })
        },
        getKpiByMonth() {
            this.gridLayout[25].loading = true
            let ld = 2

            doRequest({
                url: '/v1/web/kpi/bymonth',
                method: 'POST',
                data: {
                    startDate: this.sdate,
                    endDate: this.edate,
                    orgCode: this.selectedStationCode,
                    incomeType: 0
                }
            }, {
                success: res => {
                    this.monthProgress.current = this.sales.value
                    this.monthProgress.target = res.monthValue
                    this.yearProgress.current = this.sales.value
                    this.yearProgress.target = res.yearValue
                    this.quarterProgress.current = this.sales.value
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
                    ld --
                    if (ld == 0)
                        this.gridLayout[25].loading = false
                }
            })

            doRequest({
                url: '/v1/web/kpi/bymonth',
                method: 'POST',
                data: {
                    startDate: this.sdate,
                    endDate: this.edate,
                    orgCode: this.selectedStationCode,
                    incomeType: 1
                }
            }, {
                success: res => {
                    this.proMonthProgress.current = this.profit.value
                    this.proMonthProgress.target = res.monthValue
                    this.proYearProgress.current = this.profit.value
                    this.proYearProgress.target = res.yearValue
                    this.proQuarterProgress.current = this.profit.value
                    this.proQuarterProgress.target = res.quarterValue
                    let percent = 0

                    if (this.proMonthProgress.target <= 0) {
                        this.proMonthProgress.percent = 100
                    } else {
                        percent = parseFloat((parseFloat(this.proMonthProgress.current) * 100 / this.proMonthProgress.target).toFixed(2))
                        if (percent > 100) {
                            this.proMonthProgress.percent = 99.99 + (percent / 100000)
                        } else
                            this.proMonthProgress.percent = percent
                    }

                    if (this.proYearProgress.target <= 0) {
                        this.proYearProgress.percent = 100
                    } else {
                        percent = parseFloat((parseFloat(this.proYearProgress.current) * 100 / this.proYearProgress.target).toFixed(2))
                        if (percent > 100)
                            this.proYearProgress.percent = 99.99 + (percent / 100000)
                        else
                            this.proYearProgress.percent = percent
                    }

                    if (this.proQuarterProgress.target <= 0) {
                        this.proQuarterProgress.percent = 100
                    } else {
                        percent = parseFloat((parseFloat(this.proQuarterProgress.current) * 100 / this.proQuarterProgress.target).toFixed(2))
                        if (percent > 100)
                            this.proQuarterProgress.percent = 99.99 + (percent / 100000)
                        else
                            this.proQuarterProgress.percent = percent
                    }
                },
                fail: err => {
                    console.log(err)
                    message('error', '加载销售目标数据失败，请稍后再试')
                },
                finally: _ => {
                    ld --
                    if (ld == 0)
                        this.gridLayout[25].loading = false
                }
            })
        },
        getAllData() {
            if (!this.selectedStationCode || !this.dateRange || !this.dateRange[0] || !this.dateRange[1]) {
                message('warning', '请选择有效的油站与时间')
                return
            }
            this.sdate = this.dateRange[0]
            this.edate = this.dateRange[1]
            this.getStationInfo()
            let checked = document.querySelector('#checked')
            if(checked.style.display === 'block'){
              this.getThresholdData()
            }else {
              this.getNormalData()
            }
            // this.getKpi()
            // this.getTopGoods()
            // this.getRopGoods()
            // this.getTrend()
            // this.getCategories()
            // this.getSaleZone()
            this.getPayments()
        },
        doFullscreen() {
            let ele = document.getElementById("cnpc-single-store-main");
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
        onMainDiv(ele) {
            if (ele.path[0].className != 'el-range-input') {
                this.showPanel = false;

            }

            if (ele.path[0].className != 'ri-arrow-down-s-line cat-select') {
                this.showCatSel = false
            }
        },
        toggleShowPanel(ev) {
            this.showPanel = !this.showPanel;
            ev.stopPropagation();
        },
        showListDetail(row, p1, p2) {
            this[p1] = row;
            this[p2] = true;
        },
        beautyNum(num, pricision=2) {
            if (isNaN(num))
                return num + '';
            return beautifyNumber(num.toFixed(pricision));
        },
        onDonutResize(ele) {
            if (!ele)
                return
            let h = ele.offsetHeight, w = ele.offsetWidth;
            let m = Math.min(h, w);
            m -= 8;
            this.donutWidth = ~~(m * 100 / w);
        },
        onRootResize(root) {
            // this.$refs['cnpc-promotions-map'].clear();
            // this.$refs['cnpc-promotions-map'].mergeOptions(this.mapOptions, true, false);
            //fitty.fitAll();
            // this.fitties.forEach(e => {
            //     e.forEach(i => {
            //         console.log(i)
            //         i.fit();
            //     })
            // })
            this.gridRowHeight = window.innerHeight * 0.08
        },
        chnNumber(num) {
            return beautifyNumberChn(num, 10000)
        },
        // 切换到销售目标进度
        // toSaleTarget(){
        //     this.showSaleTarget = true;
        //     this.showProfit = false;
        //     let saleTarget = document.querySelector('#saleTarget')
        //     let profit = document.querySelector('#profit')
        //     saleTarget.style.color = '#E6A23C'
        //     profit.style.color = '#fff'
        // },
        // 切换到毛利指标进度
        // toProfit(){
        //   this.showSaleTarget = false;
        //   this.showProfit = true;
        //   let saleTarget = document.querySelector('#saleTarget')
        //   let profit = document.querySelector('#profit')
        //   saleTarget.style.color = '#fff'
        //   profit.style.color = '#E6A23C'
        // },
        // 切换到大类比例图
        // toCatPie() {
        //   this.showCatPie = true;
        //   this.showPayPie = false;
        //   let prev = document.querySelector('.prev-select')
        //   prev.style.color = "rgba(255,255,255,0.1)"
        // },
        // 鼠标移动到prev
        // onCatPie() {
        //   let prev = document.querySelector('.prev-select')
        //   if(this.showCatPie){
        //     prev.style.color = "rgba(255,255,255,0.1)"
        //   }else {
        //     prev.style.color = "rgba(255,255,255,0.9)"
        //   }
        // },
        // 鼠标离开prev
        // outCatPie() {
        //   let prev = document.querySelector('.prev-select')
        //   prev.style.color = "rgba(255,255,255,0.1)"
        // },
        // // 切换到支付方式比例图
        // toPayPie() {
        //   this.showCatPie = false;
        //   this.showPayPie = true;
        //   let next = document.querySelector('.next-select')
        //   next.style.color = "rgba(255,255,255,0.1)"
        // },
        // 鼠标移动到next
        // onPayPie() {
        //   let next = document.querySelector('.next-select')
        //   if(this.showPayPie){
        //     next.style.color = "rgba(255,255,255,0.1)"
        //   }else {
        //     next.style.color = "rgba(255,255,255,0.9)"
        //   }
        // },
        // 鼠标离开next
        // outPayPie() {
        //   let next = document.querySelector('.next-select')
        //   next.style.color = "rgba(255,255,255,0.1)"
        // },

        // 获得剔除5%以下的数据
        getThresholdData() {
          let square = document.querySelector('#toCheck')
          let checked = document.querySelector('#checked')
          square.style.display = 'none'
          checked.style.display = 'block'

          if (!this.selectedStationCode || !this.dateRange || !this.dateRange[0] || !this.dateRange[1]) {
            // message('warning', '请选择有效的油站与时间')
            return
          }
          this.getThresholdKpi()
          this.getThresholdTopGoods()
          this.getThresholdRopGoods()
          this.getThresholdTrend()
          this.getThresholdCategories()
          this.getThresholdSaleZone()
        },
        // 获得完整的数据
        getNormalData() {
          let square = document.querySelector('#toCheck')
          let checked = document.querySelector('#checked')
          square.style.display = 'block'
          checked.style.display = 'none'
          this.getKpi()
          this.getTopGoods()
          this.getRopGoods()
          this.getTrend()
          this.getCategories()
          this.getSaleZone()
        }
    }
}
</script>
