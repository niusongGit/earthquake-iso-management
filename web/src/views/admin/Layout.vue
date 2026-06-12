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
          <el-dropdown @command="handleCommand">
            <span class="user-dropdown">
              <el-icon><UserFilled /></el-icon>
              <span class="dropdown-text">{{ username }}</span>
              <el-icon class="el-icon--right"><ArrowDown /></el-icon>
            </span>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="password">
                  <el-icon><Lock /></el-icon>修改密码
                </el-dropdown-item>
                <el-dropdown-item command="logout" divided>
                  <el-icon><SwitchButton /></el-icon>退出登录
                </el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </el-header>
      <el-main class="admin-main">
        <router-view />
      </el-main>
    </el-container>

    <!-- 修改密码弹窗 -->
    <el-dialog
      v-model="pwdDialogVisible"
      title="修改密码"
      width="420px"
      :close-on-click-modal="false"
      destroy-on-close
      class="pwd-dialog"
    >
      <el-form :model="pwdForm" :rules="pwdRules" ref="pwdFormRef" label-position="top">
        <el-form-item label="旧密码" prop="old_password">
          <el-input v-model="pwdForm.old_password" type="password" show-password placeholder="请输入旧密码" />
        </el-form-item>
        <el-form-item label="新密码" prop="new_password">
          <el-input v-model="pwdForm.new_password" type="password" show-password placeholder="请输入新密码（至少6位）" />
        </el-form-item>
        <el-form-item label="确认新密码" prop="confirm_password">
          <el-input v-model="pwdForm.confirm_password" type="password" show-password placeholder="请再次输入新密码" />
        </el-form-item>
        <div v-if="pwdErrorMsg" class="pwd-error-msg">{{ pwdErrorMsg }}</div>
      </el-form>
      <template #footer>
        <el-button @click="pwdDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="pwdLoading" @click="handleChangePassword">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, computed } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { changePassword } from '../../api/admin'

const router = useRouter()

const username = computed(() => {
  try {
    const token = localStorage.getItem('token')
    if (!token) return '管理员'
    const payload = JSON.parse(atob(token.split('.')[1]))
    return payload.username || '管理员'
  } catch {
    return '管理员'
  }
})

function goFront() {
  router.push('/')
}

function handleCommand(command) {
  if (command === 'logout') {
    ElMessageBox.confirm('确定要退出登录吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }).then(() => {
      localStorage.removeItem('token')
      ElMessage.success('已退出登录')
      router.push('/admin/login')
    }).catch(() => {})
  } else if (command === 'password') {
    openPwdDialog()
  }
}

// 修改密码相关
const pwdDialogVisible = ref(false)
const pwdLoading = ref(false)
const pwdErrorMsg = ref('')
const pwdFormRef = ref(null)

const pwdForm = reactive({
  old_password: '',
  new_password: '',
  confirm_password: ''
})

const pwdRules = {
  old_password: [{ required: true, message: '请输入旧密码', trigger: 'blur' }],
  new_password: [
    { required: true, message: '请输入新密码', trigger: 'blur' },
    { min: 6, message: '密码至少6位', trigger: 'blur' }
  ],
  confirm_password: [
    { required: true, message: '请确认新密码', trigger: 'blur' },
    {
      validator: (rule, value, callback) => {
        if (value !== pwdForm.new_password) {
          callback(new Error('两次输入的密码不一致'))
        } else {
          callback()
        }
      },
      trigger: 'blur'
    }
  ]
}

function openPwdDialog() {
  pwdForm.old_password = ''
  pwdForm.new_password = ''
  pwdForm.confirm_password = ''
  pwdErrorMsg.value = ''
  pwdDialogVisible.value = true
}

async function handleChangePassword() {
  const valid = await pwdFormRef.value.validate().catch(() => false)
  if (!valid) return

  pwdErrorMsg.value = ''
  pwdLoading.value = true
  try {
    const res = await changePassword({
      old_password: pwdForm.old_password,
      new_password: pwdForm.new_password
    })
    if (res.code === 0) {
      ElMessage.success('密码修改成功，请重新登录')
      pwdDialogVisible.value = false
      localStorage.removeItem('token')
      router.push('/admin/login')
    } else {
      pwdErrorMsg.value = res.message || '修改密码失败'
    }
  } catch (e) {
    const msg = e.response?.data?.message || '修改密码失败'
    pwdErrorMsg.value = msg
  } finally {
    pwdLoading.value = false
  }
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

.header-right {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.user-dropdown {
  display: flex;
  align-items: center;
  gap: 0.3rem;
  cursor: pointer;
  color: var(--primary-color);
  font-size: 0.9rem;
  padding: 0.4rem 0.6rem;
  border-radius: 6px;
  transition: background-color 0.2s;
}

.user-dropdown:hover {
  background-color: rgba(30, 58, 138, 0.06);
}

.dropdown-text {
  max-width: 100px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.admin-main {
  padding: 1.5rem;
  max-width: 1400px;
  margin: 0 auto;
  width: 100%;
}

.pwd-error-msg {
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

  .dropdown-text {
    display: none;
  }
}
</style>

<style>
/* 修改密码弹窗 - 非scoped以穿透el-dialog */
.pwd-dialog .el-dialog__footer {
  text-align: left;
}

@media (max-width: 768px) {
  .pwd-dialog .el-dialog {
    width: calc(100vw - 20px) !important;
    margin: 10px auto !important;
  }

  .pwd-dialog .el-dialog__header {
    padding: 1rem 1rem 0.5rem;
  }

  .pwd-dialog .el-dialog__body {
    padding: 1rem;
  }

  .pwd-dialog .el-dialog__footer {
    padding: 0.5rem 1rem 1rem;
  }
}
</style>
