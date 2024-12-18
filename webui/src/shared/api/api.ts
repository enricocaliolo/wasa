import axios from 'axios'

const api = axios.create({
  baseURL: 'http://localhost:3000',
  timeout: 5000,
  headers: {
    'Content-Type': 'application/json',
    'Access-Control-Allow-Origin': '*',
  },
})

// api.interceptors.request.use((config) => {
//   const username = localStorage.getItem('username')
//   if (username) {
//     config.headers.Authorization = `Bearer ${username}`
//   }
//   return config
// })

export default api
