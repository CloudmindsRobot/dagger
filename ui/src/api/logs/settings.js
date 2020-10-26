import axios from 'axios'

export const loadSettings = () => axios('loki/settings/load', {})
