import axios from 'axios'

export const login = (data) => axios.post(`loki/auth/login`, data)
export const register = (data) => axios.post(`loki/auth/register`, data)
export const getUserSelf = () => axios.get(`loki/auth/userinfo`)
