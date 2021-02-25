import { doRequest } from '../utils/utils'

export function getProdMajorList(cb) {
    return doRequest({
        url: '/v1/web/class/main',
        method: 'get'
    }, cb)
}

export function getProdMinorList(cb) {
    return doRequest({
        url: '/v1/web/class/mid',
        method: 'get'
    }, cb)
}

export function getProdMicroList(cb) {
    return doRequest({
        url: '/v1/web/class/sub',
        method: 'get'
    }, cb)
}

export function getOrgMinorList(cb) {
    return doRequest({
        url: '/v1/web/org/branch',
        method: 'get'
    }, cb)
}

export function getOrgMicroList(cb) {
    return doRequest({
        url: '/v1/web/org/plant',
        method: 'get'
    }, cb)
}