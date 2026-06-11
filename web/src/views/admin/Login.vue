<template>
  <div class="login-page">
    <div class="login-card">
      <div class="login-header">
        <div class="logo-wrapper">
          <img src="/logo.png" alt="Logo" class="logo-img" />
        </div>
        <h2>后台管理系统</h2>
      </div>
      <el-form :model="form" :rules="rules" ref="formRef" @submit.prevent="handleLogin">
        <el-form-item prop="username">
          <el-input
            v-model="form.username"
            placeholder="请输入用户名"
            prefix-icon="User"
            size="large"
          />
        </el-form-item>
        <el-form-item prop="password">
          <el-input
            v-model="form.password"
            type="password"
            placeholder="请输入密码"
            prefix-icon="Lock"
            size="large"
            show-password
            @keyup.enter="handleLogin"
          />
        </el-form-item>
        <el-form-item>
          <el-button
            type="primary"
            native-type="button"
            size="large"
            class="login-btn"
            :loading="loading"
            @click="handleLogin"
          >登 录</el-button>
        </el-form-item>
        <div v-if="errorMsg" class="error-msg">{{ errorMsg }}</div>
      </el-form>
      <div class="back-link">
        <router-link to="/">返回前台</router-link>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { login } from '../../api/admin'

const router = useRouter()
const formRef = ref(null)
const loading = ref(false)
const errorMsg = ref('')

// 已登录则直接跳转后台
onMounted(() => {
  const token = localStorage.getItem('token')
  if (token) {
    router.replace('/admin/documents')
  }
})

const form = reactive({
  username: '',
  password: ''
})

const rules = {
  username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
  password: [{ required: true, message: '请输入密码', trigger: 'blur' }]
}

async function handleLogin() {
  const valid = await formRef.value.validate().catch(() => false)
  if (!valid) return

  errorMsg.value = ''
  loading.value = true
  try {
    const res = await login(form)
    if (res.code === 0) {
      localStorage.setItem('token', res.data.token)
      ElMessage.success('登录成功')
      router.push('/admin/documents')
    } else {
      errorMsg.value = res.message || '登录失败，请检查用户名和密码'
    }
  } catch (e) {
    const msg = e.response?.data?.message || '登录失败，请检查用户名和密码'
    errorMsg.value = msg
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login-page {
  min-height: 100vh;
  display: flex;
  justify-content: center;
  align-items: center;
  background: linear-gradient(135deg, #e3effd 0%, #f8fafc 100%);
}

.login-card {
  width: 400px;
  max-width: 90vw;
  padding: 2.5rem;
  background: white;
  border-radius: 12px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
}

.login-header {
  text-align: center;
  margin-bottom: 2rem;
}

.logo-wrapper {
  width: 56px;
  height: 56px;
  margin: 0 auto 1rem;
  overflow: hidden;
}

.logo-img {
  width: 100%;
  height: 100%;
  object-fit: contain;
}

.login-header h2 {
  color: var(--primary-color);
  font-size: 1.3rem;
}

.login-btn {
  width: 100%;
  background-color: var(--secondary-color);
  border-color: var(--secondary-color);
}

.login-btn:hover {
  background-color: var(--primary-color);
  border-color: var(--primary-color);
}

.back-link {
  text-align: center;
  margin-top: 1rem;
}

.back-link a {
  color: var(--secondary-color);
  text-decoration: none;
  font-size: 0.85rem;
}

.back-link a:hover {
  text-decoration: underline;
}

.error-msg {
  color: #ef4444;
  font-size: 0.85rem;
  text-align: center;
  padding: 0.5rem 0.75rem;
  background-color: #fef2f2;
  border: 1px solid #fecaca;
  border-radius: 6px;
  margin-top: 0.25rem;
  word-break: break-all;
  line-height: 1.5;
}
</style>
