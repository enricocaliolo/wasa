import axios from 'axios'
import { useUserStore } from '../stores/user'

const api = axios.create({
  baseURL: 'http://localhost:3000',
  timeout: 5000,
  headers: {
    'Content-Type': 'application/json',
    'Access-Control-Allow-Origin': '*',
  },
})

api.interceptors.request.use((config) => {
  const userStore = useUserStore()
  if (userStore.user !== undefined) {
    config.headers.Authorization = `Bearer ${userStore.user.userId}`
  }
  return config
})

export default api
