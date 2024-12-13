<script setup lang="ts">
import { ref } from 'vue'

const username = ref('')

function login() {
  alert(username.value)
}
</script>

<template>
  <body>
    <main>
      <div class="login-box">
        <div>
          <h1>Login</h1>
          <div class="input-form">
            <h2>Username</h2>
            <input v-model="username" />
            <button @click="login">LOGIN</button>
          </div>
        </div>
      </div>
    </main>
  </body>
</template>

<style>
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;

  /* border: 1px solid red; */
}

body {
  background-color: lightblue;
}

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
