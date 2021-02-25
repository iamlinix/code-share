<template>
    <div class="cnpc-bj-promotion">
        <p class="promo-title">
            2020年春节促销活动
        </p>
        <el-tabs stretch type="border-card" v-model="activeTab" @tab-remove="onTabRemove">
            <el-tab-pane name="promoHome" label="统计数据" class="promo-pane" id="xxxxxxxx">
                <el-row :gutter="8">
                    <el-col :span="16">
                        <el-card shadow="hover" class="top-card">
                            <div style="height: 40vh; width: 100%">
                                <v-chart theme="dark" :options="opt1" :autoresize="true"/>
                            </div>
                        </el-card>
                    </el-col>
                    <el-col :span="8">
                        <el-card shadow="hover" class="top-card">
                            <div style="height: 40vh; width: 100%">
                                <v-chart theme="dark" :options="opt2" :autoresize="true"/>
                            </div>
                        </el-card>
                    </el-col>
                    <el-col :span="24">
                        <el-card shadow="hover" class="top-card">
                            <div @click="onCardClick(0)">
                                <el-row>
                                    <el-col class="top-card-left" :span="16">
                                        <p class="top-card-left-title">
                                            <font-awesome-icon icon="yen-sign" color="#49D1E5"/>
                                            销售金额
                                        </p>
                                        <p class="top-card-left-content">{{ saleIncome.value }}</p>
                                    </el-col>
                                    <el-col :span="8">
                                        <el-row>
                                            <el-col class="top-card-right-card" :span="24">
                                                <el-col :span="6">
                                                    <font-awesome-icon v-if="saleIncome.yoyIncreasing" 
                                                        icon="long-arrow-alt-up" 
                                                        style="font-size: 40px" color="#67C23A"/>
                                                    <font-awesome-icon v-else 
                                                        icon="long-arrow-alt-down" 
                                                        style="font-size: 40px" color="#F56C6C"/>
                                                </el-col>
                                                <el-col :span="18">
                                                    <el-col class="top-card-right-title" :span="24">
                                                        同比
                                                    </el-col>
                                                    <el-col class="top-card-right-content" :span="24">
                                                        <span :style="'color: ' + (saleIncome.yoyIncreasing ? '#67C23A' : '#F56C6C')">
                                                            {{ Math.abs(saleIncome.yoy) }}%
                                                        </span>
                                                    </el-col>
                                                </el-col>
                                            </el-col>
                                            <el-col style="margin-top: 6px" class="top-card-right-card" :span="24">
                                                <el-col :span="6">
                                                    <font-awesome-icon v-if="saleIncome.momIncreasing" 
                                                        icon="long-arrow-alt-up" 
                                                        style="font-size: 40px" color="#67C23A"/>
                                                    <font-awesome-icon v-else 
                                                        icon="long-arrow-alt-down" 
                                                        style="font-size: 40px" color="#F56C6C"/>
                                                </el-col>
                                                <el-col :span="18">
                                                    <el-col class="top-card-right-title" :span="24">
                                                        环比
                                                    </el-col>
                                                    <el-col class="top-card-right-content" :span="24">
                                                        <span :style="'color: ' + (saleIncome.momIncreasing ? '#67C23A' : '#F56C6C')">
                                                            {{ Math.abs(saleIncome.mom) }}%
                                                        </span>
                                                    </el-col>
                                                </el-col>
                                            </el-col>
                                        </el-row>
                                    </el-col>
                                </el-row>
                            </div>
                        </el-card>
                    </el-col>
                    <el-col :span="12">
                        <el-card shadow="hover" class="top-card">
                            <div @click="onCardClick(1)">
                                <el-row>
                                    <el-col class="top-card-left" :span="16">
                                        <p class="top-card-left-title">
                                            <font-awesome-icon icon="yen-sign" color="#49D1E5"/>
                                            销售毛利
                                        </p>
                                        <p class="top-card-left-content">{{ saleProfit.value }}%</p>
                                    </el-col>
                                    <el-col :span="8">
                                        <el-row>
                                            <el-col class="top-card-right-card" :span="24">
                                                <el-col :span="6">
                                                    <font-awesome-icon v-if="saleProfit.yoyIncreasing" 
                                                        icon="long-arrow-alt-up" 
                                                        style="font-size: 40px" color="#67C23A"/>
                                                    <font-awesome-icon v-else 
                                                        icon="long-arrow-alt-down" 
                                                        style="font-size: 40px" color="#F56C6C"/>
                                                </el-col>
                                                <el-col :span="18">
                                                    <el-col class="top-card-right-title" :span="24">
                                                        同比
                                                    </el-col>
                                                    <el-col class="top-card-right-content" :span="24">
                                                        <span :style="'color: ' + (saleProfit.yoyIncreasing ? '#67C23A' : '#F56C6C')">
                                                            {{ Math.abs(saleProfit.yoy) }}%
                                                        </span>
                                                    </el-col>
                                                </el-col>
                                            </el-col>
                                            <el-col style="margin-top: 6px" class="top-card-right-card" :span="24">
                                                <el-col :span="6">
                                                    <font-awesome-icon v-if="saleProfit.momIncreasing" 
                                                        icon="long-arrow-alt-up" 
                                                        style="font-size: 40px" color="#67C23A"/>
                                                    <font-awesome-icon v-else 
                                                        icon="long-arrow-alt-down" 
                                                        style="font-size: 40px" color="#F56C6C"/>
                                                </el-col>
                                                <el-col :span="18">
                                                    <el-col class="top-card-right-title" :span="24">
                                                        环比
                                                    </el-col>
                                                    <el-col class="top-card-right-content" :span="24">
                                                        <span :style="'color: ' + (saleProfit.momIncreasing ? '#67C23A' : '#F56C6C')">
                                                            {{ Math.abs(saleProfit.mom) }}%
                                                        </span>
                                                    </el-col>
                                                </el-col>
                                            </el-col>
                                        </el-row>
                                    </el-col>
                                </el-row>
                            </div>
                        </el-card>
                    </el-col>
                    <el-col :span="12">
                        <el-card shadow="hover" class="top-card">
                            <div @click="onCardClick(2)">
                                <el-row>
                                    <el-col class="top-card-left" :span="16">
                                        <p class="top-card-left-title">
                                            <font-awesome-icon icon="boxes" color="#49D1E5"/>
                                            销售数量
                                        </p>
                                        <p class="top-card-left-content">{{ saleQty.value }}</p>
                                    </el-col>
                                    <el-col :span="8">
                                        <el-row>
                                            <el-col class="top-card-right-card" :span="24">
                                                <el-col :span="6">
                                                    <font-awesome-icon v-if="saleQty.yoyIncreasing" 
                                                        icon="long-arrow-alt-up" 
                                                        style="font-size: 40px" color="#67C23A"/>
                                                    <font-awesome-icon v-else 
                                                        icon="long-arrow-alt-down" 
                                                        style="font-size: 40px" color="#F56C6C"/>
                                                </el-col>
                                                <el-col :span="18">
                                                    <el-col class="top-card-right-title" :span="24">
                                                        同比
                                                    </el-col>
                                                    <el-col class="top-card-right-content" :span="24">
                                                        <span :style="'color: ' + (saleQty.yoyIncreasing ? '#67C23A' : '#F56C6C')">
                                                            {{ Math.abs(saleQty.yoy) }}%
                                                        </span>
                                                    </el-col>
                                                </el-col>
                                            </el-col>
                                            <el-col style="margin-top: 6px" class="top-card-right-card" :span="24">
                                                <el-col :span="6">
                                                    <font-awesome-icon v-if="saleQty.momIncreasing" 
                                                        icon="long-arrow-alt-up" 
                                                        style="font-size: 40px" color="#67C23A"/>
                                                    <font-awesome-icon v-else 
                                                        icon="long-arrow-alt-down" 
                                                        style="font-size: 40px" color="#F56C6C"/>
                                                </el-col>
                                                <el-col :span="18">
                                                    <el-col class="top-card-right-title" :span="24">
                                                        环比
                                                    </el-col>
                                                    <el-col class="top-card-right-content" :span="24">
                                                        <span :style="'color: ' + (saleQty.momIncreasing ? '#67C23A' : '#F56C6C')">
                                                            {{ Math.abs(saleQty.mom) }}%
                                                        </span>
                                                    </el-col>
                                                </el-col>
                                            </el-col>
                                        </el-row>
                                    </el-col>
                                </el-row>
                            </div>
                        </el-card>
                    </el-col>
                    <el-col :span="12">
                        <el-card shadow="hover" class="top-card">
                            <div @click="onCardClick(3)">
                                <el-row>
                                    <el-col class="top-card-left" :span="16">
                                        <p class="top-card-left-title">
                                            <font-awesome-icon icon="coins" color="#49D1E5"/>
                                            客单价
                                        </p>
                                        <p class="top-card-left-content">{{ customerAvg.value }}</p>
                                    </el-col>
                                    <el-col :span="8">
                                        <el-row>
                                            <el-col class="top-card-right-card" :span="24">
                                                <el-col :span="6">
                                                    <font-awesome-icon v-if="customerAvg.yoyIncreasing" 
                                                        icon="long-arrow-alt-up" 
                                                        style="font-size: 40px" color="#67C23A"/>
                                                    <font-awesome-icon v-else 
                                                        icon="long-arrow-alt-down" 
                                                        style="font-size: 40px" color="#F56C6C"/>
                                                </el-col>
                                                <el-col :span="18">
                                                    <el-col class="top-card-right-title" :span="24">
                                                        同比
                                                    </el-col>
                                                    <el-col class="top-card-right-content" :span="24">
                                                        <span :style="'color: ' + (customerAvg.yoyIncreasing ? '#67C23A' : '#F56C6C')">
                                                            {{ Math.abs(customerAvg.yoy) }}%
                                                        </span>
                                                    </el-col>
                                                </el-col>
                                            </el-col>
                                            <el-col style="margin-top: 6px" class="top-card-right-card" :span="24">
                                                <el-col :span="6">
                                                    <font-awesome-icon v-if="customerAvg.momIncreasing" 
                                                        icon="long-arrow-alt-up" 
                                                        style="font-size: 40px" color="#67C23A"/>
                                                    <font-awesome-icon v-else 
                                                        icon="long-arrow-alt-down" 
                                                        style="font-size: 40px" color="#F56C6C"/>
                                                </el-col>
                                                <el-col :span="18">
                                                    <el-col class="top-card-right-title" :span="24">
                                                        环比
                                                    </el-col>
                                                    <el-col class="top-card-right-content" :span="24">
                                                        <span :style="'color: ' + (customerAvg.momIncreasing ? '#67C23A' : '#F56C6C')">
                                                            {{ Math.abs(customerAvg.mom) + '%' }}
                                                        </span>
                                                    </el-col>
                                                </el-col>
                                            </el-col>
                                        </el-row>
                                    </el-col>
                                </el-row>
                            </div>
                        </el-card>
                    </el-col>
                    <el-col :span="12">
                        <el-card shadow="hover" class="top-card" @click="onCardClick(3)">
                            <div @click="onCardClick(4)">
                                <el-row>
                                    <el-col class="top-card-left" :span="16">
                                        <p class="top-card-left-title">
                                            <font-awesome-icon icon="random" color="#49D1E5"/>
                                            油非转换率
                                        </p>
                                        <p class="top-card-left-content">{{ transfer.value }}%</p>
                                    </el-col>
                                    <el-col :span="8">
                                        <el-row>
                                            <el-col class="top-card-right-card" :span="24">
                                                <el-col :span="6">
                                                    <font-awesome-icon v-if="transfer.yoyIncreasing" 
                                                        icon="long-arrow-alt-up" 
                                                        style="font-size: 40px" color="#67C23A"/>
                                                    <font-awesome-icon v-else 
                                                        icon="long-arrow-alt-down" 
                                                        style="font-size: 40px" color="#F56C6C"/>
                                                </el-col>
                                                <el-col :span="18">
                                                    <el-col class="top-card-right-title" :span="24">
                                                        同比
                                                    </el-col>
                                                    <el-col class="top-card-right-content" :span="24">
                                                        <span :style="'color: ' + (transfer.yoyIncreasing ? '#67C23A' : '#F56C6C')">
                                                            {{ Math.abs(transfer.yoy) }}%
                                                        </span>
                                                    </el-col>
                                                </el-col>
                                            </el-col>
                                            <el-col style="margin-top: 6px" class="top-card-right-card" :span="24">
                                                <el-col :span="6">
                                                    <font-awesome-icon v-if="transfer.momIncreasing" 
                                                        icon="long-arrow-alt-up" 
                                                        style="font-size: 40px" color="#67C23A"/>
                                                    <font-awesome-icon v-else 
                                                        icon="long-arrow-alt-down" 
                                                        style="font-size: 40px" color="#F56C6C"/>
                                                </el-col>
                                                <el-col :span="18">
                                                    <el-col class="top-card-right-title" :span="24">
                                                        环比
                                                    </el-col>
                                                    <el-col class="top-card-right-content" :span="24">
                                                        <span :style="'color: ' + (transfer.momIncreasing ? '#67C23A' : '#F56C6C')">
                                                            {{ Math.abs(transfer.mom) }}%
                                                        </span>
                                                    </el-col>
                                                </el-col>
                                            </el-col>
                                        </el-row>
                                    </el-col>
                                </el-row>
                            </div>
                        </el-card>
                    </el-col>
                    <el-col :span="8">
                        <el-card shadow="hover" class="top-card">
                            <div class="top-card-left-small" @click="onCardClick(5)">
                                <p class="top-card-left-title">
                                    <font-awesome-icon icon="donate" color="#49D1E5"/>
                                    促销商品综合做贡献率
                                </p>
                                <p class="top-card-left-content">{{ contribution }}%</p>
                            </div>
                        </el-card>
                    </el-col>
                    <el-col :span="8">
                        <el-card shadow="hover" class="top-card">
                            <div class="top-card-left-small" @click="onCardClick(6)">
                                <p class="top-card-left-title">
                                    <font-awesome-icon icon="gas-pump" color="#49D1E5"/>
                                    促销商品在售站数
                                </p>
                                <p class="top-card-left-content">{{ totalStations }}</p>
                            </div>
                        </el-card>
                    </el-col>
                    <el-col :span="8">
                        <el-card shadow="hover" class="top-card">
                            <div class="top-card-left-small" @click="onCardClick(7)">
                                <p class="top-card-left-title">
                                    <font-awesome-icon icon="car" color="#49D1E5"/>
                                    汽服站数
                                </p>
                                <p class="top-card-left-content">{{ stationsWithCarService }}</p>
                            </div>
                        </el-card>
                    </el-col>
                    <el-col :span="6">
                        <el-card shadow="hover" class="top-card">
                            <div class="top-card-left-small" style="color: #67C23A" @click="onCardClick(8)">
                                <p class="top-card-left-title">
                                    <font-awesome-icon icon="award" color="#67C23A"/>
                                    Lv1 站点数
                                </p>
                                <p class="top-card-left-content">{{ lv1 }}</p>
                            </div>
                        </el-card>
                    </el-col>
                    <el-col :span="6">
                        <el-card shadow="hover" class="top-card">
                            <div class="top-card-left-small" style="color: #409EFF" @click="onCardClick(9)">
                                <p class="top-card-left-title">
                                    <font-awesome-icon icon="award" color="#409EFF"/>
                                    Lv2 站点数
                                </p>
                                <p class="top-card-left-content">{{ lv2 }}</p>
                            </div>
                        </el-card>
                    </el-col>
                    <el-col :span="6">
                        <el-card shadow="hover" class="top-card">
                            <div class="top-card-left-small" style="color: #E6A23C" @click="onCardClick(10)">
                                <p class="top-card-left-title">
                                    <font-awesome-icon icon="award" color="#E6A23C"/>
                                    Lv3 站点数
                                </p>
                                <p class="top-card-left-content">{{ lv3 }}</p>
                            </div>
                        </el-card>
                    </el-col>
                    <el-col :span="6">
                        <el-card shadow="hover" class="top-card">
                            <div class="top-card-left-small" style="color: #F56C6C" @click="onCardClick(11)">
                                <p class="top-card-left-title">
                                    <font-awesome-icon icon="award" color="#F56C6C"/>
                                    Lv4 站点数
                                </p>
                                <p class="top-card-left-content">{{ lv4 }}</p>
                            </div>
                        </el-card>
                    </el-col>
                </el-row>
            </el-tab-pane>
            <el-tab-pane closable v-for="(v, i) in tabs" :key="v.key" 
                :label="v.label" :name="v.name" class="promo-pane">
                <el-table size="mini" style="margin-top: 10px" border stripe 
                    :data="details" :ref="v.ref">
                    <el-table-column type="index" />
                    <el-table-column sortable label="加油站" prop="name" />
                    <el-table-column v-if="v.index === 0" sortable label="促销品收入" prop="promoSale" />
                    <el-table-column v-if="v.index === 0" sortable label="总收入" prop="totalSale" />
                    <el-table-column v-if="v.index === 0" sortable label="总收入同比" prop="saleYoy">
                        <template slot-scope="scope">
                            <span>{{ scope.row.saleYoy }}%</span>
                        </template>
                    </el-table-column>
                    <el-table-column v-if="v.index === 0" sortable label="总收入环比" prop="saleYoy">
                        <template slot-scope="scope">
                            <span>{{ scope.row.saleMom }}%</span>
                        </template>
                    </el-table-column>
                    <el-table-column v-if="v.index === 2" sortable label="促销品销量" prop="promoCount" />
                    <el-table-column v-if="v.index === 2" sortable label="总销量" prop="totalCount" />
                    <el-table-column v-if="v.index === 2" sortable label="销量同比" prop="countYoy">
                        <template slot-scope="scope">
                            <span>{{ scope.row.countYoy }}%</span>
                        </template>
                    </el-table-column>
                    <el-table-column v-if="v.index === 2" sortable label="销量环比" prop="countMom">
                        <template slot-scope="scope">
                            <span>{{ scope.row.countMom }}%</span>
                        </template>
                    </el-table-column>
                    <el-table-column v-if="v.index === 1" sortable label="毛利率" prop="profit">
                        <template slot-scope="scope">
                            <span>{{ scope.row.profit }}%</span>
                        </template>
                    </el-table-column>
                    <el-table-column v-if="v.index === 1" sortable label="毛利率同比" prop="profitYoy">
                        <template slot-scope="scope">
                            <span>{{ scope.row.profitYoy }}%</span>
                        </template>
                    </el-table-column>
                    <el-table-column v-if="v.index === 1" sortable label="毛利率环比" prop="profitMom">
                        <template slot-scope="scope">
                            <span>{{ scope.row.profitMom }}%</span>
                        </template>
                    </el-table-column>
                    <el-table-column v-if="v.index === 3" sortable label="客单价" prop="avg" />
                    <el-table-column v-if="v.index === 3" sortable label="客单价同比" prop="avgYoy">
                        <template slot-scope="scope">
                            <span>{{ scope.row.avgYoy }}%</span>
                        </template>
                    </el-table-column>
                    <el-table-column v-if="v.index === 3" sortable label="客单价环比" prop="avgMom">
                        <template slot-scope="scope">
                            <span>{{ scope.row.avgMom }}%</span>
                        </template>
                    </el-table-column>
                    <el-table-column v-if="v.index === 4" sortable label="转化率" prop="transfer">
                        <template slot-scope="scope">
                            <span>{{ scope.row.transfer }}%</span>
                        </template>
                    </el-table-column>
                    <el-table-column v-if="v.index === 4" sortable label="转化率同比" prop="transferYoy">
                        <template slot-scope="scope">
                            <span>{{ scope.row.transferYoy }}%</span>
                        </template>
                    </el-table-column>
                    <el-table-column v-if="v.index === 4" sortable label="转化率环比" prop="transferMom">
                        <template slot-scope="scope">
                            <span>{{ scope.row.transferMom }}%</span>
                        </template>
                    </el-table-column>
                    <el-table-column sortable label="有无汽车服务" prop="hasCarService">
                        <template slot-scope="scope">
                            <span>{{ scope.row.hasCarService }}</span>
                        </template>
                    </el-table-column>
                    <el-table-column sortable label="加油站评级" prop="level">
                        <template slot-scope="scope">
                            <el-tag v-if="scope.row.level == 0" size="small" effect="dark" type="success">Lv1</el-tag>
                            <el-tag v-if="scope.row.level == 1" size="small" effect="dark" type="primary">Lv2</el-tag>
                            <el-tag v-if="scope.row.level == 2" size="small" effect="dark" type="warning">Lv3</el-tag>
                            <el-tag v-if="scope.row.level == 3" size="small" effect="dark" type="danger">Lv4</el-tag>
                        </template>
                    </el-table-column>
                </el-table>
            </el-tab-pane>
        </el-tabs>
        <!--div class="promo-legend">
            <el-tag effect="dark" type="success">Lv1  (收入增长，毛利增长)</el-tag>
            <el-tag effect="dark" type="primary">Lv2  (收入增长，毛利降低)</el-tag>
            <el-tag effect="dark" type="warning">Lv3  (收入降低，毛利增长)</el-tag>
            <el-tag effect="dark" type="danger">Lv4  (收入降低，毛利降低)</el-tag>
        </div-->
        <!--el-table style="margin-top: 10px" border stripe :data="details">
            <el-table-column type="index" />
            <el-table-column sortable label="加油站" prop="name" />
            <el-table-column sortable label="促销品收入" prop="promoSale" />
            <el-table-column sortable label="总收入" prop="totalSale" />
            <el-table-column sortable label="总收入同比" prop="salePercent">
                <template slot-scope="scope">
                    <span>{{ scope.row.salePercent }}%</span>
                </template>
            </el-table-column>
            <el-table-column sortable label="毛利率" prop="profit">
                <template slot-scope="scope">
                    <span>{{ scope.row.profit }}%</span>
                </template>
            </el-table-column>
            <el-table-column sortable label="毛利率同比" prop="profitPercent">
                <template slot-scope="scope">
                    <span>{{ scope.row.profitPercent }}%</span>
                </template>
            </el-table-column>
            <el-table-column sortable label="有无汽车服务" prop="hasCarService">
                <template slot-scope="scope">
                    <span>{{ scope.row.hasCarService }}</span>
                </template>
            </el-table-column>
            <el-table-column sortable label="加油站评级" prop="level">
                <template slot-scope="scope">
                    <el-tag v-if="scope.row.level == 0" size="small" effect="dark" type="success">Lv1</el-tag>
                    <el-tag v-if="scope.row.level == 1" size="small" effect="dark" type="primary">Lv2</el-tag>
                    <el-tag v-if="scope.row.level == 2" size="small" effect="dark" type="warning">Lv3</el-tag>
                    <el-tag v-if="scope.row.level == 3" size="small" effect="dark" type="danger">Lv4</el-tag>
                </template>
            </el-table-column>
        </el-table-->
    </div>
</template>

<style lang="scss">
.cnpc-bj-promotion {
    .echarts {
        width: 100%;
        height: 100%;
    }

    background-color: white;
    height: calc(100vh - 10px);
    padding: 8px 8px 8px 8px;
    overflow: auto;

    p.promo-title {
        font-weight: bold;
        font-size: 18px;
        padding: 8px 0 10px 0;
        text-align: center;
    }

    .promo-pane {
        overflow: auto;
        padding: 8px 8px 8px 8px;
        max-height: calc(100vh - 140px); 
    }

    div.promo-legend {
        width: 100%;
        border-radius: 5px;
        padding: 20px 20px 20px 20px;
        background-color: #E9E9EB;
        text-align: center;
        
        .el-tag {
            margin-right: 12px;
        }
    }

    .top-card {
        margin-top: 10px;
        cursor: pointer;
        background-color: #1F1F1F;
        color: #808080
    }

    .top-card-left {
        position: relative;
        top: 30%;
        transform: translateY(30%);

        .top-card-left-title {
            font-size: 20px;
        }

        .top-card-left-content {
            font-size: 2.5em;
            font-weight: bold;
        }
    }

    .top-card-left-small {
        position: relative;
        top: 5%;
        transform: translateY(5%);

        .top-card-left-title {
            font-size: 20px;
        }

        .top-card-left-content {
            font-size: 2.5em;
            font-weight: bold;
        }
    }

    .top-card-right-card {
        padding: 8px 8px 8px 8px;

        .top-card-right-title {
            font-size: 12px;
        }

        .top-card-right-content {
            font-size: 19px;
            font-weight: bold;
        }
    }
}
</style>

<script>
import { randomNumber } from '../../utils/utils'
import ECharts from 'vue-echarts'
import 'echarts/lib/chart/line'
import 'echarts/lib/chart/bar'
import 'echarts/lib/chart/pie'
import 'echarts/lib/component/title'
import 'echarts/lib/component/tooltip'
import 'echarts/lib/component/legend'
import 'echarts/lib/component/toolbox'

export default {
    components: {
        'v-chart': ECharts
    },
    data() {
        return {
            saleIncome: {
                value: 22345.44,
                yoy: 27.3,
                mom: -17.3,
                yoyIncreasing: false,
                momIncreasing: false
            },
            saleProfit: {
                value: 22345.44,
                yoy: 27.3,
                mom: -17.3,
                yoyIncreasing: false,
                momIncreasing: false
            },
            saleQty: {
                value: 22345.44,
                yoy: 27.3,
                mom: -17.3,
                yoyIncreasing: false,
                momIncreasing: false
            },
            customerAvg: {
                value: 45.2,
                yoy: 22.2,
                mom: 33.3,
                yoyIncreasing: false,
                momIncreasing: false
            },
            transfer: {
                value: 45.22,
                yoy: 23.4,
                mom: 55.5,
                yoyIncreasing: false,
                momIncreasing: false
            },
            contribution: 66.6,
            totalStations: 123,
            stationsWithCarService: 22,
            details: [],
            lv1: 0,
            lv2: 0,
            lv3: 0,
            lv4: 0,

            activeTab: 'promoHome',
            tabs: [],
            tabLabels: [
                '销售额排行',
                '毛利率排行',
                '销售量排行',
                '客单价排行',
                '油非转换率排行',
                '促销贡献度排行',
                '在售油站',
                '汽服油站',
                'Lv1 油站',
                'Lv2 油站',
                'Lv3 油站',
                'Lv4 油站',
            ],
            tabRefs: [
                'cnpc-promo-money',
                'cnpc-promo-profit',
                'cnpc-promo-count',
                'cnpc-promo-average',
                'cnpc-promo-transfer',
                'cnpc-promo-percent',
                'cnpc-promo-on-sale',
                'cnpc-promo-car-service',
                'cnpc-promo-lv1',
                'cnpc-promo-lv2',
                'cnpc-promo-lv3',
                'cnpc-promo-lv4',
            ],
            checkList: [],

            opt1: {
                title: {
                    text: '同环比数据',
                    left: 'center'
                },
                tooltip: {
                    trigger: 'axis'
                },
                legend: {
                    data: ['当前', '同比', '环比'],
                    left: 'left'
                },
                backgroundColor: '#1F1F1F',
                toolbox: {
                    show: false,
                    feature: {
                        magicType: {show: true, type: ['line', 'bar']},
                        saveAsImage: {show: true}
                    }
                },
                calculable: true,
                xAxis: [
                    {
                        type: 'category',
                        data: ['销售金额', '销售数量', '客单价', '油非转换率']
                    }
                ],
                yAxis: [
                    {
                        type: 'value'
                    }
                ],
                series: [
                    {
                        name: '当前',
                        type: 'bar',
                        data: [72, 60, 54.3, 12.2],
                    },
                    {
                        name: '同比',
                        type: 'bar',
                        data: [64, 43, 33.2, 21],
                    },
                    {
                        name: '环比',
                        type: 'bar',
                        data: [57, 44, 23.6, 12],
                    }
                ]},
            opt2: {
                title: {
                    text: '油站分级',
                    left: 'center'
                },
                tooltip: {
                    trigger: 'item',
                    formatter: '{a} <br/>{b}: {c} ({d}%)'
                },
                legend: {
                    orient: 'vertical',
                    left: 10,
                    data: ['Lv1', 'Lv2', 'Lv3', 'Lv4']
                },
                backgroundColor: '#1F1F1F',
                series: [
                    {
                        name: '分级',
                        type: 'pie',
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
                            {value: 335, name: 'Lv1'},
                            {value: 310, name: 'Lv2'},
                            {value: 234, name: 'Lv3'},
                            {value: 135, name: 'Lv4'}
                        ]
                    }
                ]
            }
        }
    },
    mounted() {
        this.doUpdate();
        setInterval(this.doUpdate, 5000);
    },
    methods: {
        doUpdate() {
            this.saleIncome.value = randomNumber(1000000, 5000000, 2);
            this.saleIncome.yoy = randomNumber(-100, 100, 2);
            this.saleIncome.mom = randomNumber(-100, 100, 2);
            this.saleIncome.yoyIncreasing = this.saleIncome.yoy >= 0
            this.saleIncome.momIncreasing = this.saleIncome.mom >= 0

            this.saleProfit.value = randomNumber(10, 40, 2);
            this.saleProfit.yoy = randomNumber(-100, 100, 2);
            this.saleProfit.mom = randomNumber(-100, 100, 2);
            this.saleProfit.yoyIncreasing = this.saleProfit.yoy >= 0
            this.saleProfit.momIncreasing = this.saleProfit.mom >= 0

            this.saleQty.value = randomNumber(10000, 20000);
            this.saleQty.yoy = randomNumber(-100, 100, 2);
            this.saleQty.mom = randomNumber(-100, 100, 2);
            this.saleQty.yoyIncreasing = this.saleIncome.yoy >= 0
            this.saleQty.momIncreasing = this.saleIncome.mom >= 0

            this.customerAvg.value = randomNumber(20, 200, 2);
            this.customerAvg.yoy = randomNumber(-100, 100, 2);
            this.customerAvg.mom = randomNumber(-100, 100, 2);
            this.customerAvg.yoyIncreasing = this.customerAvg.yoy >= 0
            this.customerAvg.momIncreasing = this.customerAvg.mom >= 0


            this.transfer.value = randomNumber(5, 100, 2);
            this.transfer.yoy = randomNumber(-100, 100, 2);
            this.transfer.mom = randomNumber(-100, 100, 2);
            this.transfer.yoyIncreasing = this.transfer.yoy >= 0
            this.transfer.momIncreasing = this.transfer.mom >= 0


            this.contribution = randomNumber(10, 100, 2);
            this.totalStations = randomNumber(100, 300, 0);
            this.stationsWithCarService = randomNumber(0, this.totalStations, 0);

            this.lv1 = randomNumber(0, 100);
            this.lv2 = randomNumber(0, 100);
            this.lv3 = randomNumber(0, 100);
            this.lv4 = randomNumber(0, 100);

            let details = []
            for (let i = 0; i < 20; i ++) {
                let e = {
                    name: `北京销售公司第${i + 1}加油站`,
                    promoSale: randomNumber(100, 1000, 2),
                    promoCount: randomNumber(10000, 20000),
                    totalSale: randomNumber(1000, 2000, 2),
                    totalCount: randomNumber(20000, 30000),
                    saleYoy: randomNumber(-100, 100, 2),
                    saleMom: randomNumber(-100, 100, 2),
                    countYoy: randomNumber(-100, 100, 2),
                    countMom: randomNumber(-100, 100, 2),
                    profit: randomNumber(1, 30, 2),
                    profitYoy: randomNumber(-100, 100, 2),
                    profitMom: randomNumber(-100, 100, 2),
                    hasCarService: parseInt(randomNumber(-1, 1, 0)) > 0 ? '有' : '无',
                    avg: randomNumber(20, 200, 2),
                    avgYoy: randomNumber(-100, 100, 2),
                    avgMom: randomNumber(-100, 100, 2),
                    transfer: randomNumber(10, 80, 2),
                    transferYoy: randomNumber(-100, 100, 2),
                    transferMom: randomNumber(-100, 100, 2),
                };
                if (e.salePercent > 0 && e.profitPercent > 0)
                    e.level = 0
                else if (e.salePercent > 0 && e.profitPercent < 0)
                    e.level = 1
                else if (e.salePercent < 0 && e.profitPercent > 0)
                    e.level = 2
                else
                    e.level = 3
                details.push(e)
            }
            this.details = details;
        },
        onCardClick(index) {
            for(let i = 0; i < this.tabs.length; i ++) {
                if (this.tabs[i].index === index) {
                    this.activeTab = this.tabs[i].name;
                    return;
                }
            }
            this.tabs.push({
                index: index,
                label: this.tabLabels[index],
                name: 'tab' + index,
                key: 'key-' + index,
                ref: this.tabRefs[index]
            })
            this.$nextTick(() => {
                this.activeTab = 'tab' + index;
            })
        },
        onTabRemove(name) {
            for(let i = 0; i < this.tabs.length; i ++) {
                if (this.tabs[i].name === name) {
                    this.tabs.splice(i, 1);
                    break;
                }
            }

            this.activeTab = 'promoHome'
        }
    }
}
</script>