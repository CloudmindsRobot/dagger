import axios from 'axios'

export const listLogUser = (data) => axios('loki/user', { params: data })
