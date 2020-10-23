import axios from 'axios'

export const listQueryHistory = (data) =>
  axios('loki/history', { params: data })
export const createQueryHistoryLabel = (data) =>
  axios.post(`loki/history/create`, data)
export const deleteQueryHistoryLabel = (data) =>
  axios.delete(`loki/history/delete/${data.id}`, data)

export const listQueryResultSnapshot = (data) =>
  axios('loki/snapshot', { params: data })
export const createQueryResultSnapshot = (data) =>
  axios.post(`loki/snapshot/create`, data)
export const deleteQueryResultSnapshot = (data) =>
  axios.delete(`loki/snapshot/delete/${data.id}`, data)
export const detailQueryResultSnapshot = (data) =>
  axios(`loki/snapshot/detail/${data.id}`, data)
