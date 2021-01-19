import axios from 'axios'

export const listLabelValues = (data) =>
  axios('loki/label/values', { params: data })
export const listQueryRanges = (data) =>
  axios('loki/query_range', { params: data })
export const exportQueryRanges = (data) =>
  axios('loki/export', { params: data })
export const listContext = (data) => axios('loki/context', { params: data })
export const listLabels = (data) => axios('loki/labels', { params: data })
export const transformLogQL = (data) => axios('loki/logql', { params: data })
export const listLabelAllValues = (data) => axios(`loki/label/values/${data}`)
