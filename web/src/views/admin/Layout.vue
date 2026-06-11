<template>
  <div class="admin-layout">
    <el-container>
      <el-header class="admin-header">
        <div class="header-left">
          <div class="logo-wrapper">
            <img src="/logo.png" alt="Logo" class="logo-img" />
          </div>
          <span class="header-title">地震国际标准后台管理</span>
        </div>
        <div class="header-right">
          <el-button text @click="goFront">前台首页</el-button>
          <el-button text type="danger" @click="handleLogout">退出登录</el-button>
        </div>
      </el-header>
      <el-main class="admin-main">
        <router-view />
      </el-main>
    </el-container>
  </div>
</template>

<script setup>
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'

const router = useRouter()

function goFront() {
  router.push('/')
}

function handleLogout() {
  localStorage.removeItem('token')
  ElMessage.success('已退出登录')
  router.push('/admin/login')
}
</script>

<style scoped>
.admin-layout {
  min-height: 100vh;
  background: var(--bg-color);
}

.admin-header {
  background: linear-gradient(180deg, #e3effd 0%, #f8fafc 100%);
  border-bottom: 1px solid rgba(30, 58, 138, 0.1);
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0 1.5rem;
  height: 60px;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 0.8rem;
}

.logo-wrapper {
  width: 36px;
  height: 36px;
  overflow: hidden;
}

.logo-img {
  width: 100%;
  height: 100%;
  object-fit: contain;
}

.header-title {
  font-size: 1.1rem;
  font-weight: bold;
  color: var(--primary-color);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

@media (max-width: 768px) {
  .admin-header {
    padding: 0 0.8rem;
  }

  .header-left {
    gap: 0.5rem;
    flex: 1;
    min-width: 0;
  }

  .header-title {
    font-size: 0.9rem;
  }

  .header-right {
    flex-shrink: 0;
  }

  .header-right .el-button {
    padding: 4px 8px;
    font-size: 0.8rem;
  }
}

.header-right {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.admin-main {
  padding: 1.5rem;
  max-width: 1400px;
  margin: 0 auto;
  width: 100%;
}
</style>
