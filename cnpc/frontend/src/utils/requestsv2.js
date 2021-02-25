import Axios from 'axios'
import Config from '../config/index.js'
import { Notification, Loading  } from 'element-ui';
import { getToken, removeToken, getAddr } from '../utils/dataStorage.js'
import https from 'https'

const httpsAgent = new https.Agent({
    rejectUnauthorized: false,
})

const service = Axios.create({
    baseURL: Config.apiUrl + '/' + Config.apiPrefix,
    headers: {
        'Accept': '*/*'
    },
    timeout: Config.timeout,
    httpsAgent: Config.devEnv ? httpsAgent : null
});
service.defaults.retry = Config.requestRetry;
service.defaults.retryDelay = Config.requestRetryDelay;

service.interceptors.request.use(
    request => {
        // let noParameters = request.url.indexOf('?') === -1;
        request.headers.token = getToken();
        //request.headers['public-addr'] = getAddr();
        // config.url = noParameters ? config.url+'?access_token=' + getToken() : config.url + '&access_token=' + getToken();
        return request
    },
    error => {
        return Promise.reject(error)
    }
);


service.interceptors.response.use(
    response => {
        if (response.status !== 200) {
            Notification({
                type: 'warning',
                title: '网络请求错误',
                message: '错误码：' + response.status
            });
            return Promise.reject('Error Status Code');
        }

        return response;
    },
    error => {
        console.log('request error:', error);
        if (error.response.status === 403) {
            removeToken();
            // Notification({
            //     type: 'error',
            //     title: '用户未登录',
            //     message: '未检测到您的登录状态，请重新登录后使用系统'
            // });
            setTimeout(_=>{
                // window.location.href = '/login';
                Config.router.push({path: '/login'})
            },500);
        }

        return Promise.reject(error.response);//千万不能去掉，，，否则请求超时会进入到then方法，导致逻辑错误。
    }
);

export default service
