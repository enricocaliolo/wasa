<script setup lang="ts">
import { userAPI } from '@/modules/auth/api/user-api'
import { ref } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()

const username = ref('')

const handleLogin = async () => {
  try {
    await userAPI.login(username.value)
    router.push('/')
  } catch (error) {
    alert(error)
    console.error('Login failed:', error)
  }
}
</script>

<template>
  <main>
    <div class="login-box">
      <div>
        <h1>Login</h1>
        <div class="input-form">
          <h2>Username</h2>
          <input v-model="username" />
          <button @click="handleLogin">LOGIN</button>
        </div>
      </div>
    </div>
  </main>
</template>

<style scoped>
main {
  height: 100vh;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  width: 90vw;
  margin: 0 auto;
  max-width: 1110px;
}

.login-box {
  border: 1px solid black;
  border-radius: 1.5em;
  padding: 1em;
}
</style>
