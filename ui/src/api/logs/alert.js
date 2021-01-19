import axios from 'axios'

export const listLogRule = (data) => axios('loki/rule', { params: data })

export const listLogGroup = (data) => axios('loki/group', { params: data })
export const createLogGroup = (data) => axios.post('loki/group/create', data)
export const updateLogGroup = (data) =>
  axios.post(`loki/group/update/${data.id}`, data)
export const deleteLogGroup = (data) =>
  axios.delete(`loki/group/delete/${data.id}`, data)
export const joinLogGroup = (data) => axios.post(`loki/group/join`, data)
export const leaveLogGroup = (data) => axios.post(`loki/group/leave`, data)

export const createLogRule = (data) => axios.post('loki/rule/create', data)
export const updateLogRule = (data) =>
  axios.post(`loki/rule/update/${data.id}`, data)
export const deleteLogRule = (data) =>
  axios.delete(`loki/rule/delete/${data.id}`, data)
export const downloadRuleFile = () => axios.get('loki/rule/download')

export const listLogEvent = (data) => axios(`loki/event`, { params: data })
export const archiveLogEvent = (data) => axios.post(`loki/event/archive`, data)
export const listLogEventDetails = (data) =>
  axios(`loki/event/details/${data.id}`, data)
