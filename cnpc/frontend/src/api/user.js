import requests from '../utils/requests.js'
import { doRequest } from '../utils/utils'
import { setToken, removeToken } from '../utils/dataStorage'

export function login(params, sucess, fail) {
    return doRequest({
        url: '/login',
        method: 'post',
        loading: true,
        data: params
    }, {
        success: sucess,
        fail: fail,
    })
}

export function getUserList(cb) {
    return doRequest({
        url: '/v1/web/user/list',
        method: 'get'
    }, cb)
}

export function logout() {
    removeToken();
}

export function changePassword(params, cb) {
    return requests({
       url: '/v1/web/user/update',
       method: 'post',
       params: params
    }, cb);
}

export function addUser(params, cb) {
    return requests({
        url: '/v1/web/user/add',
        method: 'post',
        data: params
    }, cb)
}

export function deleteUser(uid, cb) {
    return requests({
        url: '/v1/web/user/del' + uid,
        method: 'get',
    }, cb)
}
