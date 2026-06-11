<template>
  <div class="front-layout">
    <!-- 顶部导航 -->
    <header class="header">
      <div class="header-container">
        <div class="logo-wrapper">
          <img src="/logo.png" alt="Logo" class="logo-img" />
        </div>
        <div class="logo-text">地震国际标准管理平台</div>
        <div class="header-right">
          <a href="javascript:;" class="admin-link" @click="goAdmin">
            <el-icon><Setting /></el-icon>
            <span>管理</span>
          </a>
        </div>
      </div>
    </header>

    <!-- 主体内容 -->
    <main class="main-content">
      <!-- 搜索与筛选区 -->
      <section class="search-filter-section">
        <div class="search-group">
          <div class="search-box">
            <el-input
              v-model="keyword"
              placeholder="搜索编号、名称、摘要、范围..."
              clearable
              @keyup.enter="handleSearch"
            />
          </div>
          <el-button type="primary" class="btn-search" @click="handleSearch">搜索</el-button>
        </div>
        <div class="filter-box">
          <el-select v-model="stage" placeholder="当前阶段" clearable @change="handleFilterChange">
            <el-option label="PWI" value="PWI" />
            <el-option label="NP" value="NP" />
            <el-option label="WD" value="WD" />
            <el-option label="CD" value="CD" />
            <el-option label="DIS" value="DIS" />
            <el-option label="FDIS" value="FDIS" />
            <el-option label="IS" value="IS" />
          </el-select>
        </div>
      </section>

      <!-- 工具栏 -->
      <div class="toolbar">
        <div class="tabs">
          <div
            class="tab-item"
            :class="{ active: sort === 'createTime' }"
            @click="switchTab('createTime')"
          >创建时间 ⬇</div>
          <div
            class="tab-item"
            :class="{ active: sort === 'publishDate' }"
            @click="switchTab('publishDate')"
          >发布日期 ⬇</div>
        </div>
        <div class="doc-count">共 <span>{{ total }}</span> 项</div>
      </div>

      <!-- 文档列表 -->
      <section class="doc-list" v-loading="loading">
        <div
          v-for="doc in documents"
          :key="doc.id"
          class="doc-card"
          @click="goDetail(doc.id)"
        >
          <div class="doc-badge-row">
            <span class="badge badge-code">{{ doc.iso_code }}</span>
            <span class="badge badge-type">{{ doc.type }}</span>
            <span class="badge badge-stage">{{ doc.current_stage }}</span>
          </div>
          <div class="doc-title">{{ doc.name }}</div>
          <div class="doc-meta-row">
            <div class="meta-pair">
              <span>标准所属: {{ doc.standard_belongs_to }}</span>
              <span>发布日期: {{ doc.publish_date ? formatDate(doc.publish_date) : '-' }}</span>
            </div>
            <div class="meta-pair">
              <span>首发编号: {{ doc.first_publish_code || '-' }}</span>
              <span>所属: {{ doc.belongs_to || '-' }}</span>
            </div>
            <div>地震相关度: <span class="stars">{{ getStars(doc.earthquake_relevance) }}</span></div>
          </div>
        </div>

        <el-empty v-if="!loading && documents.length === 0" description="暂无文档" />
      </section>

      <!-- 分页 -->
      <div class="pagination-wrapper" v-if="total > 0">
        <el-pagination
          v-model:current-page="page"
          :page-size="pageSize"
          :total="total"
          layout="prev, pager, next"
          @current-change="loadData"
        />
      </div>
    </main>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { Setting } from '@element-plus/icons-vue'
import { getDocumentList } from '../../api/front'

const router = useRouter()

const keyword = ref('')
const stage = ref('')
const sort = ref('createTime')
const page = ref(1)
const pageSize = ref(10)
const total = ref(0)
const documents = ref([])
const loading = ref(false)

function formatDate(dateStr) {
  if (!dateStr) return '-'
  return dateStr.substring(0, 10)
}

function getStars(count) {
  return '★'.repeat(count || 0) + '☆'.repeat(5 - (count || 0))
}

async function loadData() {
  loading.value = true
  try {
    const res = await getDocumentList({
      keyword: keyword.value,
      stage: stage.value,
      sort: sort.value,
      page: page.value,
      pageSize: pageSize.value
    })
    if (res.code === 0) {
      documents.value = res.data.list || []
      total.value = res.data.total || 0
    }
  } catch (e) {
    console.error('加载文档列表失败', e)
  } finally {
    loading.value = false
  }
}

function handleSearch() {
  page.value = 1
  loadData()
}

function handleFilterChange() {
  page.value = 1
  loadData()
}

function switchTab(key) {
  sort.value = key
  page.value = 1
  loadData()
}

function goDetail(id) {
  router.push(`/detail/${id}`)
}

function goAdmin() {
  const token = localStorage.getItem('token')
  if (token) {
    router.push('/admin/documents')
  } else {
    router.push('/admin/login')
  }
}

onMounted(() => {
  loadData()
})
</script>

<style scoped>
.header {
  background: linear-gradient(180deg, #e3effd 0%, #f8fafc 100%);
  color: var(--primary-color);
  padding: 1.2rem 1rem;
  border-bottom: 1px solid rgba(30, 58, 138, 0.1);
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.02);
}

.header-container {
  max-width: 1200px;
  margin: 0 auto;
  display: flex;
  align-items: center;
  gap: 0.8rem;
}

.logo-wrapper {
  width: 46px;
  height: 46px;
  overflow: hidden;
  display: flex;
  align-items: center;
  justify-content: center;
}

.logo-img {
  width: 100%;
  height: 100%;
  object-fit: contain;
}

.logo-text {
  font-size: 1.35rem;
  font-weight: bold;
  letter-spacing: 0.5px;
  color: var(--primary-color);
  white-space: nowrap;
}

.header-right {
  margin-left: auto;
}

.admin-link {
  display: inline-flex;
  align-items: center;
  gap: 0.3rem;
  font-size: 0.85rem;
  color: var(--secondary-color);
  text-decoration: none;
  padding: 0.35rem 0.75rem;
  background-color: rgba(37, 99, 235, 0.08);
  border-radius: 4px;
  white-space: nowrap;
  transition: all 0.2s;
}

.admin-link:hover {
  background-color: var(--secondary-color);
  color: white;
}

.main-content {
  max-width: 1200px;
  margin: 1rem auto;
  padding: 0 1rem;
}

.search-filter-section {
  background: white;
  padding: 1rem;
  border-radius: 8px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.05);
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
  margin-bottom: 1.5rem;
}

.search-group {
  display: flex;
  gap: 0.5rem;
  width: 100%;
}

.search-box {
  flex: 1;
}

.btn-search {
  background-color: var(--secondary-color);
  border-color: var(--secondary-color);
  white-space: nowrap;
}

.btn-search:hover {
  background-color: var(--primary-color);
  border-color: var(--primary-color);
}

.filter-box {
  width: 100%;
}

.filter-box :deep(.el-select) {
  width: 100%;
}

.toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  border-bottom: 2px solid var(--border-color);
  margin-bottom: 1rem;
}

.tabs {
  display: flex;
  gap: 1rem;
}

.tab-item {
  padding: 0.5rem 0;
  font-size: 0.9rem;
  font-weight: 500;
  color: var(--text-muted);
  cursor: pointer;
  position: relative;
}

.tab-item.active {
  color: var(--secondary-color);
  font-weight: bold;
}

.tab-item.active::after {
  content: '';
  position: absolute;
  bottom: -2px;
  left: 0;
  width: 100%;
  height: 2px;
  background-color: var(--secondary-color);
}

.doc-count {
  color: var(--text-muted);
  font-size: 0.8rem;
  white-space: nowrap;
}

.doc-list {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
  min-height: 200px;
}

.doc-card {
  background: white;
  border: 1px solid var(--border-color);
  border-radius: 8px;
  padding: 1rem;
  cursor: pointer;
  transition: background-color 0.2s;
}

.doc-card:hover {
  background-color: #f1f5f9;
}

.doc-badge-row {
  display: flex;
  gap: 0.4rem;
  margin-bottom: 0.5rem;
  flex-wrap: wrap;
}

.badge {
  font-size: 0.7rem;
  padding: 0.15rem 0.4rem;
  border-radius: 4px;
  font-weight: bold;
}

.badge-type {
  background-color: #e0f2fe;
  color: #0369a1;
}

.badge-stage {
  background-color: #fef3c7;
  color: #d97706;
}

.badge-code {
  background-color: #f1f5f9;
  color: #334155;
  font-family: monospace;
}

.doc-title {
  font-size: 1.05rem;
  font-weight: 600;
  color: #1e293b;
  margin-bottom: 0.5rem;
  line-height: 1.4;
}

.doc-meta-row {
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
  font-size: 0.8rem;
  color: var(--text-muted);
}

.meta-pair {
  display: flex;
  gap: 1.5rem;
}

.stars {
  color: var(--star-color);
}

.pagination-wrapper {
  display: flex;
  justify-content: center;
  margin-top: 1.5rem;
  padding-bottom: 2rem;
}

/* PC端适配 */
@media (min-width: 768px) {
  .header {
    padding: 1rem 2rem;
  }

  .header-container {
    gap: 1rem;
  }

  .logo-wrapper {
    width: 48px;
    height: 48px;
  }

  .logo-text {
    font-size: 1.5rem;
  }

  .main-content {
    margin: 2rem auto;
  }

  .search-filter-section {
    flex-direction: row;
    align-items: center;
    padding: 1.5rem;
  }

  .search-group {
    flex: 1;
  }

  .filter-box {
    width: auto;
  }

  .filter-box :deep(.el-select) {
    width: 200px;
  }

  .doc-meta-row {
    flex-direction: row;
    gap: 1.5rem;
  }

  .doc-card {
    padding: 1.5rem;
  }

  .doc-title {
    font-size: 1.2rem;
  }
}
</style>
