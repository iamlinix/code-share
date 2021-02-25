<template>
    <div class="promo-main-v3" id="cnpc-promo-main-v3">
        <el-dialog :visible.sync="dialogShow" 
            custom-class="transparent-modal-dialog"
            :show-close="false">
        </el-dialog>
        <div class="promo-header-v3">
            <font-awesome-icon v-if="fullScreen" icon="compress-arrows-alt" 
                class="promo-full-screen" @click="exitFullscreen" />
            <font-awesome-icon v-else icon="arrows-alt" 
                class="promo-full-screen" @click="doFullscreen" />
            <p class="promo-title">{{ promotion.title }}</p>
            <span class="promo-timer">{{ currentTimeDesc }}</span>
        </div>
        <grid-layout class="promo-main-body" :layout.sync="gridLayout" 
            :use-css-transforms="true" :vertical-compact="true"
            :is-draggable="false" :is-resizable="false" :row-height="72">
            <grid-item
                   :x="gridLayout[0].x"
                   :y="gridLayout[0].y"
                   :w="gridLayout[0].w"
                   :h="gridLayout[0].h"
                   :i="gridLayout[0].i"
                   :key="gridLayout[0].i">
                   <div class="list-with-ranking">
                       <p class="list-title" @click="dialogShow = true">
                           站级销售排行TOP10
                       </p>
                       <table>
                           <tbody>
                               <tr v-for="(v, i) in promotion.topStations" :key="'tr-top-key-' + i">
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
                                            <span slot="trigger">{{ v.name }}</span>
                                            {{ v.name }}
                                        </tippy>
                                    </td>
                                   <td>
                                       <tippy arrow append-to="parent">
                                            <span slot="trigger">{{ chnNumber(v.amount) }}</span>
                                            {{ v.amount }}
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
                   :key="gridLayout[1].i">
                   <div class="sale-amount-grid">
                       <p class="sale-amount-title">总销售额</p>
                       <p class="sale-amount-number">
                           <font-awesome-icon icon="yen-sign"></font-awesome-icon>
                            {{ beautyNum(promotion.sales.amount) }}
                       </p>
                   </div>
            </grid-item>
            <grid-item
                   :x="gridLayout[2].x"
                   :y="gridLayout[2].y"
                   :w="gridLayout[2].w"
                   :h="gridLayout[2].h"
                   :i="gridLayout[2].i"
                   :key="gridLayout[2].i">
                   <div class="percent-grid">
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
                   :key="gridLayout[3].i">
                   <div class="percent-grid">
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
                   :key="gridLayout[4].i">
                   <div class="list-with-ranking">
                       <p class="list-title">
                           站级销售排行倒数TOP10
                       </p>
                       <table>
                           <tbody>
                               <tr v-for="(v, i) in promotion.ropStations" :key="'tr-rop-key-' + i">
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
                                            <span slot="trigger">{{ v.name }}</span>
                                            {{ v.name }}
                                        </tippy>
                                    </td>
                                   <td>
                                        <tippy arrow append-to="parent">
                                            <span slot="trigger">{{ chnNumber(v.amount) }}</span>
                                            {{ v.amount }}
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
                   :key="gridLayout[5].i">
                   <v-chart theme="macarons" :options="mapOptions" :autoresize="true"/>
            </grid-item>
            <grid-item
                   :x="gridLayout[6].x"
                   :y="gridLayout[6].y"
                   :w="gridLayout[6].w"
                   :h="gridLayout[6].h"
                   :i="gridLayout[6].i"
                   :key="gridLayout[6].i">
                   <div class="list-with-ranking">
                       <p class="list-title">
                           畅销商品TOP10
                       </p>
                       <table>
                           <tbody>
                               <tr v-for="(v, i) in promotion.topGoods" :key="'tr-top-good-key-' + i">
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
                                            <span slot="trigger">{{ v.name }}</span>
                                            {{ v.name }}
                                        </tippy>
                                    </td>
                                   <td>
                                        <tippy arrow append-to="parent">
                                            <span slot="trigger">{{ chnNumber(v.amount) }}</span>
                                            {{ v.amount }}
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
                   :key="gridLayout[7].i">
                   <div class="list-with-ranking">
                       <p class="list-title">
                           滞销商品TOP10
                       </p>
                       <table>
                           <tbody>
                               <tr v-for="(v, i) in promotion.ropGoods" :key="'tr-rop-good-key-' + i">
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
                                            <span slot="trigger">{{ v.name }}</span>
                                            {{ v.name }}
                                        </tippy>
                                    </td>
                                   <td>
                                        <tippy arrow append-to="parent">
                                            <span slot="trigger">{{ chnNumber(v.amount) }}</span>
                                            {{ v.amount }}
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
                   >
                   <div id="donut-auto-resize-transfer" style="width: 100%; height: 100%">
                    <vc-donut :size="donutWidth" unit="%" :thickness="30" 
                            background="#052ACF" 
                            :sections="[{ value: promotion.transfer.amount, color: percentColors[~~(promotion.transfer.amount/20)].color }]"
                            foreground="#C0C4CC">
                            <p style="color: rgba(255, 255, 255, 0.7); font-size: 14px; margin-bottom: 2px">
                                油非转换
                            </p>
                            <p :style="`color: ${percentColors[~~(promotion.transfer.amount/20)].color}; font-size: 21px; font-weight: bold`">
                                {{ promotion.transfer.amount.toFixed(2) }}%
                            <p/>
                    </vc-donut>
                   </div>
                   <!--v-chart theme="macarons" :options="transferOptions" :autoresize="true"/-->
            </grid-item>
            <grid-item
                   :x="gridLayout[9].x"
                   :y="gridLayout[9].y"
                   :w="gridLayout[9].w"
                   :h="gridLayout[9].h"
                   :i="gridLayout[9].i"
                   :key="gridLayout[9].i">
                   <div class="percent-grid">
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
                   :x="gridLayout[10].x"
                   :y="gridLayout[10].y"
                   :w="gridLayout[10].w"
                   :h="gridLayout[10].h"
                   :i="gridLayout[10].i"
                   :key="gridLayout[10].i">
                   <div class="percent-grid">
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
                   :x="gridLayout[11].x"
                   :y="gridLayout[11].y"
                   :w="gridLayout[11].w"
                   :h="gridLayout[11].h"
                   :i="gridLayout[11].i"
                   :key="gridLayout[11].i">
                   <div class="single-num-grid">
                       <p class="single-num-title">
                           客单价
                       </p>
                       <p class="single-num-gold">
                           <font-awesome-icon icon="coins" />
                           {{ promotion.average.amount }}
                       </p>
                   </div>
            </grid-item>
            <grid-item
                   :x="gridLayout[12].x"
                   :y="gridLayout[12].y"
                   :w="gridLayout[12].w"
                   :h="gridLayout[12].h"
                   :i="gridLayout[12].i"
                   :key="gridLayout[12].i">
                   <div class="percent-grid">
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
            </grid-item>
            <grid-item
                   :x="gridLayout[13].x"
                   :y="gridLayout[13].y"
                   :w="gridLayout[13].w"
                   :h="gridLayout[13].h"
                   :i="gridLayout[13].i"
                   :key="gridLayout[13].i">
                   <div class="percent-grid">
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
                   :x="gridLayout[14].x"
                   :y="gridLayout[14].y"
                   :w="gridLayout[14].w"
                   :h="gridLayout[14].h"
                   :i="gridLayout[14].i"
                   :key="gridLayout[14].i">
                   <div class="single-num-grid">
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
                   </div>
            </grid-item>
            <grid-item
                   :x="gridLayout[15].x"
                   :y="gridLayout[15].y"
                   :w="gridLayout[15].w"
                   :h="gridLayout[15].h"
                   :i="gridLayout[15].i"
                   :key="gridLayout[15].i">
                   <div class="single-num-grid">
                       <p class="single-num-title">
                           在售站点数
                       </p>
                       <p class="single-num-normal">
                           <font-awesome-icon icon="gas-pump"/>
                           {{ promotion.activeStationNumber }}
                       </p>
                   </div>
            </grid-item>
            <grid-item
                   :x="gridLayout[16].x"
                   :y="gridLayout[16].y"
                   :w="gridLayout[16].w"
                   :h="gridLayout[16].h"
                   :i="gridLayout[16].i"
                   :key="gridLayout[16].i">
                   <div class="single-num-grid">
                       <p class="single-num-title">
                           汽服站点数
                       </p>
                       <p class="single-num-normal">
                           <font-awesome-icon icon="car"/>
                           {{ promotion.stationWithCarService }}
                       </p>
                   </div>
            </grid-item>
            <grid-item
                   :x="gridLayout[17].x"
                   :y="gridLayout[17].y"
                   :w="gridLayout[17].w"
                   :h="gridLayout[17].h"
                   :i="gridLayout[17].i"
                   :key="gridLayout[17].i">
                   <v-chart theme="macarons" :options="levelTreeOptions" :autoresize="true"/>
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

    .transparent-modal-dialog {
        background-color: rgba(255, 255, 255, 0.2);
    }

    .cdc-container {
        height: 100%;
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
        color: white;
        text-align: center;
        padding-top: 12px;
        font-size: 1.5em;
        
    }

    span.promo-timer {
        color: rgba(255, 255, 255, 0.8);
        position: absolute;
        top: 15px;
        right: 20px;
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

            :hover {
                color: rgba(255, 255, 255, 0.9);
            }
        }
    }

    .promo-main-body {
        padding: 6px 6px 6px 6px;

        .sale-amount-grid {
            .sale-amount-title {
                color:rgba(255, 255, 255, 0.7);
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
                color: #C7CD4F;
                font-size: 4.5em;
                font-weight: bold;
                padding: 20px 0 0 0;
                text-align: center;
            }
        }

        .percent-grid {
            .percent-title {
                color:rgba(255, 255, 255, 0.7);
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
                font-size: 22px;
                margin-top: 12px;
                font-weight: bold;
                color: #67C23A;
            }
            
            .percent-number-down {
                font-size: 22px;
                margin-top: 12px;
                font-weight: bold;
                color: #F56C6C;
            }
        }

        .single-num-grid {
            .single-num-title {
                color:rgba(255, 255, 255, 0.7);
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
                font-size: 40px;
                font-weight: bold;
                margin-top: 40px;
                text-align: center;
            }

            .single-num-normal {
                color: #E6A23C;
                font-size: 40px;
                font-weight: bold;
                margin-top: 40px;
                text-align: center;
            }

            .single-num-small {
                color: #E6A23C;
                font-size: 30px;
                font-weight: bold;
                margin-top: 20px;
                text-align: center;
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
                color: rgba(255, 255, 255, 0.7);
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
            
            table {
                color: rgba(255, 255, 255, 0.5);
                font-size: 16px;
                display: table;
                vertical-align: bottom;
                width: 100%;

                tr {
                    
                    height: 25px;
                    
                    td:nth-child(1) {
                        width: 30px;
                        display: inline-block;
                        vertical-align: middle;
                        text-align: center;
                    }
                    td:nth-child(2) {
                        width: calc(100% - 110px);
                        white-space: nowrap;
                        overflow: hidden;
                        text-overflow: ellipsis;
                        display: inline-block;
                        vertical-align: middle;
                        font-size: 13px;
                    }
                    td:nth-child(3) {
                        padding-left: 10px;
                        width: 80px;
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
}
</style>

<script>
import { toggleKey, beautifyNumberChn, beautifyNumber, randomNumber } from '../../utils/utils'
import { addListener, removeListener } from 'resize-detector'
import moment from 'moment'
import VueGridLayout from 'vue-grid-layout';
import fitty from 'fitty'

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

export default {
    components: {
        'grid-layout': VueGridLayout.GridLayout,
        'grid-item': VueGridLayout.GridItem,
        'v-chart': ECharts
    },
    mounted() {
        this.promoTimer();
        setInterval(this.promoTimer, 1000);
        window.addEventListener('keydown', this.onKeyDown);
        
        
        let ele = document.getElementById('donut-auto-resize-transfer');
        addListener(ele, this.onDonutResize);
        this.onDonutResize(ele);

        let root = document.getElementById('cnpc-promo-main-v3');
        addListener(root, this.onRootResize);

        this.fitties.push(fitty('.sale-amount-number'));
        this.fitties.push(fitty('.single-num-gold'));
    },
    data() {
        return {
            fitties: [],
            dialogShow: false,
            donutWidth: 60,
            ticker: 0,
            fullScreen: false,
            currentTimeDesc: '',
            percentColors: [
                {color: '#f56c6c', percentage: 20},
                {color: '#e6a23c', percentage: 40},
                {color: '#5cb87a', percentage: 60},
                {color: '#1989fa', percentage: 80},
                {color: '#6f7ad3', percentage: 100}
            ],
            promotion: {
                title: '北京销售公司2020春节促销活动',
                topStations: [
                    {
                        name: '北京销售公司第四分公司力克力多加油站',
                        amount: 33222300.23
                    },{
                        name: '北京销售公司第一分公司永定林加油站',
                        amount: 33222.23
                    },{
                        name: '北京销售第四分公司独一处加油站',
                        amount: 33222.23
                    },
                    {
                        name: '北京销售第四分公司双燕加油站',
                        amount: 33222.23
                    },
                    {
                        name: '北京销售第四分公司钰鑫增加油站',
                        amount: 33222.23
                    },
                    {
                        name: '北京销售第一分公司华石加油站',
                        amount: 33222.23
                    },
                    {
                        name: '北京销售第一分公司广仟达加油站',
                        amount: 33222.23
                    },
                    {
                        name: '北京销售第二分公司岳上加油站',
                        amount: 33222.23
                    },
                    {
                        name: '北京销售第一分公司北京鑫海源加油站',
                        amount: 33222.23
                    },
                    {
                        name: '北京销售第四分公司施惠达加油站',
                        amount: 33222.23
                    },
                ],
                ropStations: [
                    {
                        name: '北京销售第二分公司燕鑫兆元加油站',
                        amount: 332223333.23
                    },{
                        name: '北京销售第四分公司兴海源加油站',
                        amount: 33222.23
                    },{
                        name: '北京销售第二分公司富华兆元加油站',
                        amount: 33222.23
                    },
                    {
                        name: '北京销售第三分公司西小营加油站',
                        amount: 33222.23
                    },
                    {
                        name: '北京销售第二分公司军广加油站',
                        amount: 33222.23
                    },
                    {
                        name: '北京销售第四分公司燕兴加油站',
                        amount: 33222.23
                    },
                    {
                        name: '北京销售第四分公司驰华加油站',
                        amount: 33222.23
                    },
                    {
                        name: '北京销售第四分公司华荣加油站',
                        amount: 33222.23
                    },
                    {
                        name: '北京销售第四分公司华星加油站',
                        amount: 33222.23
                    },
                    {
                        name: '北京销售第四分公司欣燕林加油站',
                        amount: 33222.23
                    },
                ],
                topGoods: [
                    {
                        name: '雷达 电热蚊香块(B型)替换装',
                        amount: 33664.6
                    },
                    {
                        name: '揽菊 牌黑蚊香',
                        amount: 33664.6
                    },
                    {
                        name: '卓异 蚊香',
                        amount: 33664.6
                    },
                    {
                        name: '睡得香 蚊香',
                        amount: 33664.6
                    },
                    {
                        name: '枪手 电热蚊香片 小螳螂无味型72片促',
                        amount: 33664.6
                    },
                    {
                        name: '全无敌 电热蚊香器',
                        amount: 33664.6
                    },
                    {
                        name: '卓异 蚊香加热器',
                        amount: 33664.6
                    },
                    {
                        name: '枪手 杀虫气雾剂(清新柠檬)',
                        amount: 33664.6
                    },
                    {
                        name: '安琪 儿面巾纸',
                        amount: 33664.6
                    },
                    {
                        name: 'LUXANA 迷你双层纸巾',
                        amount: 33664.6
                    },
                ],
                ropGoods: [
                    {
                        name: '日晒 自然盐 300G',
                        amount: 123.4
                    },
                    {
                        name: '贝贝熊 卫生桶',
                        amount: 123.4
                    },
                    {
                        name: '昂立 天然元牛乳提取物复合片',
                        amount: 123.4
                    },
                    {
                        name: '恒寿堂 禧礼盒',
                        amount: 123.4
                    },
                    {
                        name: '罗福斯 ”重型订书机B45/3',
                        amount: 123.4
                    },
                    {
                        name: '易达 ”订书针(统一针)',
                        amount: 123.4
                    },
                    {
                        name: '中华 铅笔',
                        amount: 123.4
                    },
                    {
                        name: '中华 铅笔芯0.5',
                        amount: 123.4
                    },
                    {
                        name: '得力 48K复写纸',
                        amount: 123.4
                    },
                    {
                        name: '真彩 5支装中性笔芯',
                        amount: 123.4
                    },
                ],
                sales: {
                    amount: 5632818.22,
                    yoy: 12.33,
                    mom: -22.34,
                },
                transfer: {
                    amount: 33.445,
                    yoy: 22.33,
                    mom: -55.34,
                },
                average: {
                    amount: 34.77,
                    yoy: 44.33,
                    mom: -17.34,
                },
                contribute: {
                    gross: 33.72,
                    profit: 47.65
                },
                activeStationNumber: 117,
                stationWithCarService: 72,
                stationLevels: [

                ]
            },
            gridLayout: [
                {"x":0,"y":0,"w":3,"h":4,"i":"0"},
                {"x":3,"y":0,"w":5,"h":2,"i":"1"},
                {"x":8,"y":0,"w":1,"h":1,"i":"2"},
                {"x":8,"y":1,"w":1,"h":1,"i":"3"},
                {"x":9,"y":0,"w":3,"h":4,"i":"4"},
                {"x":3,"y":2,"w":6,"h":4,"i":"5"},
                {"x":0,"y":4,"w":3,"h":4,"i":"6"},
                {"x":9,"y":4,"w":3,"h":4,"i":"7"},
                {"x":3,"y":8,"w":2,"h":2,"i":"8"},
                {"x":5,"y":8,"w":1,"h":1,"i":"9"},
                {"x":5,"y":9,"w":1,"h":1,"i":"10"},
                {"x":6,"y":8,"w":2,"h":2,"i":"11"},
                {"x":8,"y":8,"w":1,"h":1,"i":"12"},
                {"x":8,"y":9,"w":1,"h":1,"i":"13"},
                {"x":0,"y":10,"w":2,"h":2,"i":"14"},
                {"x":2,"y":10,"w":2,"h":2,"i":"15"},
                {"x":4,"y":10,"w":2,"h":2,"i":"16"},
                {"x":6,"y":10,"w":6,"h":2,"i":"17"},
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
                        symbolSize: '13',
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
                        data: [
                            {
                                name: '北京销售公司第四分公司力克力多加油站',
                                value: [116.416725, 39.699995, 55, 1]
                            },
                            {
                                name: '北京销售公司第一分公司永定林加油站',
                                value: [116.191493, 39.897194, 110, 2]
                            },
                            {
                                name: '北京销售第四分公司独一处加油站',
                                value: [116.219802, 39.808169, 32, 3]   
                            },
                            {
                                name: '北京销售第四分公司双燕加油站',
                                value: [116.119912, 39.765108, 55, 4]
                            },
                            {
                                name: '北京销售第四分公司钰鑫增加油站',
                                value: [115.991702, 39.670391, 55, 5]
                            },
                            {
                                name: '北京销售第一分公司华石加油站',
                                value: [116.123244, 39.943474, 55, 6]
                            },
                            {
                                name: '北京销售第一分公司广仟达加油站',
                                value: [116.291657, 39.810895, 55, 7]
                            },
                            {
                                name: '北京销售第二分公司岳上加油站',
                                value: [116.845252, 39.818971, 55, 8]
                            },
                            {
                                name: '北京销售第一分公司北京鑫海源加油站',
                                value: [116.40989, 39.842482, 55, 9]
                            },
                            {
                                name: '北京销售第四分公司施惠达加油站',
                                value: [116.383885, 39.769992, 55, 10]
                            },
                            {
                                name: '北京销售第二分公司燕鑫兆元加油站',
                                value: [116.490532, 39.90658, 55, -1]
                            },
                            {
                                name: '北京销售第四分公司兴海源加油站',
                                value: [116.407761, 39.75818, 55, -2]
                            },
                            {
                                name: '北京销售第二分公司富华兆元加油站',
                                value: [116.573721, 39.812454, 55, -3]
                            },
                            {
                                name: '北京销售第三分公司西小营加油站',
                                value: [116.56838, 40.196784, 55, -4]
                            },
                            {
                                name: '北京销售第二分公司军广加油站',
                                value: [116.898788, 39.956873, 55, -5]
                            },
                            {
                                name: '北京销售第四分公司燕兴加油站',
                                value: [116.331839, 39.511876, 55, -6]
                            },
                            {
                                name: '北京销售第四分公司驰华加油站',
                                value: [116.332789, 39.63061, 55, -7]
                            },
                            {
                                name: '北京销售第四分公司华荣加油站',
                                value: [116.386238, 39.749236, 55, -8]
                            },
                            {
                                name: '北京销售第四分公司华星加油站',
                                value: [116.329288, 39.51205, 55, -9]
                            },
                            {
                                name: '北京销售第四分公司欣燕林加油站',
                                value: [116.33343, 39.592, 55, -10]
                            }
                        ]
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
                    data: [{
                        name: 'Level1',
                        value: 33,
                    }, {
                        name: 'Level2',
                        value: 44,
                    }, {
                        name: 'Level3',
                        value: 21,
                    }, {
                        name: 'Level4',
                        value: 19,
                    }],
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
            }
        }
    },
    methods: {
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
            if (this.ticker === 5) {
                this.ticker = 0;

                this.promotion.sales.amount = randomNumber(1000000, 9000000000, 2);
                this.promotion.sales.yoy = randomNumber(-100, 100, 2); 
                this.promotion.sales.mom = randomNumber(-100, 100, 2);

                this.promotion.average.amount = randomNumber(20, 1000, 2);
                this.promotion.average.yoy = randomNumber(-100, 100, 2); 
                this.promotion.average.mom = randomNumber(-100, 100, 2);

                this.promotion.transfer.amount = randomNumber(5, 100, 2);
                this.promotion.transfer.yoy = randomNumber(-100, 100, 2); 
                this.promotion.transfer.mom = randomNumber(-100, 100, 2);

                this.promotion.contribute.gross = randomNumber(10, 99, 2);
                this.promotion.contribute.profit = randomNumber(10, 99, 2);

                this.promotion.activeStationNumber = randomNumber(20, 200);
                this.promotion.stationWithCarService = randomNumber(20, 150);

                this.promotion.topStations.forEach(e => {
                    e.amount = randomNumber(20000, 30000, 2) + '';
                })

                this.promotion.ropStations.forEach(e => {
                    e.amount = randomNumber(10000, 20000, 2) + '';
                })

                this.promotion.topGoods.forEach(e => {
                    e.amount = randomNumber(3000, 20000, 2) + '';
                })

                this.promotion.ropGoods.forEach(e => {
                    e.amount = randomNumber(1000, 3000, 2) + '';
                })
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
            this.fitties.forEach(e => {
                e.forEach(i => {
                    i.fit();
                })
            })
        }
    }
}
</script>