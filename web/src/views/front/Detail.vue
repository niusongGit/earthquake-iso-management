<template>
  <div class="detail-layout">
    <!-- 顶部导航 -->
    <header class="header">
      <div class="header-container">
        <div class="logo-wrapper">
          <img src="/logo.png" alt="Logo" class="logo-img" />
        </div>
        <div class="logo-text">地震国际标准管理平台</div>
      </div>
    </header>

    <main class="main-content" v-loading="loading">
      <!-- 返回按钮 -->
      <div class="back-bar">
        <el-button type="primary" plain round @click="goBack">
          <el-icon><ArrowLeft /></el-icon> 返回列表
        </el-button>
      </div>

      <div class="detail-body" v-if="doc">
        <!-- 详情信息面板 -->
        <div class="details-panel">
          <div class="detail-grid">
            <div class="detail-item">
              <div class="label">国际标准编号</div>
              <div class="value code-value">{{ doc.iso_code }}</div>
            </div>
            <div class="detail-item">
              <div class="label">国际标准名称</div>
              <div class="value">{{ doc.name }}</div>
            </div>
            <div class="detail-item">
              <div class="label">标准类型</div>
              <div class="value"><span class="badge badge-type">{{ doc.type }}</span></div>
            </div>
            <div class="detail-item">
              <div class="label">当前阶段</div>
              <div class="value"><span class="badge badge-stage">{{ doc.current_stage }}</span></div>
            </div>
            <div class="detail-item">
              <div class="label">标准所属</div>
              <div class="value">{{ doc.standard_belongs_to || '-' }}</div>
            </div>
            <div class="detail-item">
              <div class="label">所属</div>
              <div class="value">{{ doc.belongs_to || '-' }}</div>
            </div>
            <div class="detail-item">
              <div class="label">首次发布编号</div>
              <div class="value">{{ doc.first_publish_code || '-' }}</div>
            </div>
            <div class="detail-item">
              <div class="label">发布日期</div>
              <div class="value">{{ doc.publish_date ? formatDate(doc.publish_date) : '-' }}</div>
            </div>
            <div class="detail-item">
              <div class="label">创建时间</div>
              <div class="value muted-value">{{ formatDateTime(doc.created_at) }}</div>
            </div>
            <div class="detail-item">
              <div class="label">地震相关度</div>
              <div class="value stars">{{ getStars(doc.earthquake_relevance) }}</div>
            </div>
            <div class="detail-item">
              <div class="label">摘要</div>
              <div class="value value-block">{{ doc.summary || '-' }}</div>
            </div>
            <div class="detail-item">
              <div class="label">范围</div>
              <div class="value value-block">{{ doc.scope || '-' }}</div>
            </div>
          </div>
        </div>

        <!-- PDF预览面板 -->
        <div class="preview-panel">
          <div class="preview-header">
            <span>附件在线预览 (PDF)</span>
            <el-button
              v-if="doc.attachment"
              type="success"
              size="small"
              @click="handleDownload"
            >下载附件</el-button>
          </div>
          <div class="pdf-viewer">
            <div v-if="doc.attachment" class="pdf-container">
              <!-- 翻页控制 -->
              <div class="pdf-toolbar">
                <el-button size="small" :disabled="currentPage <= 1" @click="prevPage">上一页</el-button>
                <span class="page-info">{{ currentPage }} / {{ totalPages }}</span>
                <el-button size="small" :disabled="currentPage >= totalPages" @click="nextPage">下一页</el-button>
              </div>
              <!-- PDF渲染 -->
              <div class="pdf-canvas-wrapper">
                <div v-if="pdfLoading" class="pdf-loading-overlay">
                  <el-icon class="loading-icon"><Loading /></el-icon>
                  <span>PDF加载中...</span>
                </div>
                <vue-pdf-embed
                  ref="pdfRef"
                  :source="pdfSource"
                  :page="currentPage"
                  :width="pdfWidth"
                  @loaded="onPdfLoaded"
                  @rendered="onPdfRendered"
                  @loading-failed="onPdfLoadingFailed"
                  @rendering-failed="onPdfRenderingFailed"
                  class="pdf-embed"
                />
              </div>
            </div>
            <div v-else class="pdf-empty">
              <p>暂无附件</p>
            </div>
          </div>
        </div>
      </div>
    </main>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onBeforeUnmount } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { getDocumentDetail } from '../../api/front'
import VuePdfEmbed from 'vue-pdf-embed'

const route = useRoute()
const router = useRouter()

const doc = ref(null)
const loading = ref(false)
const pdfRef = ref(null)
const currentPage = ref(1)
const totalPages = ref(0)
const pdfLoading = ref(true)
const windowWidth = ref(window.innerWidth)

// 响应式PDF渲染宽度：PC端使用较大宽度，手机端不设置（使用容器宽度）
const pdfWidth = computed(() => {
  if (windowWidth.value >= 768) {
    // PC端：预览面板约占55%宽度，减去内边距后取较大值
    return Math.min(windowWidth.value * 0.5, 760)
  }
  // 手机端：不指定宽度，由组件使用容器clientWidth
  return undefined
})

function handleResize() {
  windowWidth.value = window.innerWidth
}

const downloadUrl = computed(() => {
  if (!doc.value) return ''
  return `/api/front/documents/${doc.value.id}/download`
})

const previewUrl = computed(() => {
  if (!doc.value || !doc.value.attachment) return ''
  return `/api/front/documents/${doc.value.id}/preview`
})

const pdfSource = computed(() => {
  if (!doc.value || !doc.value.attachment) return null
  const origin = window.location.origin
  return {
    url: `${origin}${previewUrl.value}`,
    wasmUrl: `${origin}/wasm/`,
    cMapUrl: `${origin}/cmaps/`,
    standardFontDataUrl: `${origin}/standard_fonts/`,
    cMapPacked: true,
    useWasm: false
  }
})

function formatDate(dateStr) {
  if (!dateStr) return '-'
  return dateStr.substring(0, 10)
}

function formatDateTime(dateStr) {
  if (!dateStr) return '-'
  return dateStr.replace('T', ' ').substring(0, 19)
}

function getStars(count) {
  return '★'.repeat(count || 0) + '☆'.repeat(5 - (count || 0))
}

function goBack() {
  router.push('/')
}

function handleDownload() {
  if (doc.value && doc.value.attachment) {
    window.open(downloadUrl.value, '_blank')
  }
}

function onPdfLoaded(pdf) {
  // pdf.document 返回 PDFDocumentProxy，可获取总页数
  if (pdf && pdf.numPages) {
    totalPages.value = pdf.numPages
  } else if (pdf && pdf.document && pdf.document.numPages) {
    totalPages.value = pdf.document.numPages
  }
}

function onPdfRendered() {
  // 渲染完成，隐藏加载动画
  pdfLoading.value = false
}

function onPdfLoadingFailed(error) {
  console.error('[PDF] 加载失败:', error)
  if (error && error.message) {
    console.error('[PDF] 错误详情:', error.message)
  }
  if (error && error.stack) {
    console.error('[PDF] 堆栈:', error.stack)
  }
  pdfLoading.value = false
}

function onPdfRenderingFailed(error) {
  console.error('[PDF] 渲染失败:', error)
  if (error && error.message) {
    console.error('[PDF] 错误详情:', error.message)
  }
  if (error && error.stack) {
    console.error('[PDF] 堆栈:', error.stack)
  }
  pdfLoading.value = false
}

function prevPage() {
  if (currentPage.value > 1) {
    pdfLoading.value = true
    currentPage.value--
  }
}

function nextPage() {
  if (currentPage.value < totalPages.value) {
    pdfLoading.value = true
    currentPage.value++
  }
}

async function loadDetail() {
  loading.value = true
  try {
    const id = route.params.id
    const res = await getDocumentDetail(id)
    if (res.code === 0) {
      doc.value = res.data
    }
  } catch (e) {
    console.error('加载文档详情失败', e)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  window.addEventListener('resize', handleResize)
  loadDetail()
})

onBeforeUnmount(() => {
  window.removeEventListener('resize', handleResize)
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

.main-content {
  max-width: 1200px;
  margin: 1rem auto;
  padding: 0 1rem;
}

.back-bar {
  margin-bottom: 1rem;
}

.back-bar :deep(.el-button) {
  font-weight: 600;
  padding: 0.5rem 1.2rem;
  font-size: 0.9rem;
  box-shadow: 0 2px 6px rgba(37, 99, 235, 0.2);
}

.detail-body {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.details-panel {
  width: 100%;
  padding: 1rem;
  background: white;
  border-radius: 8px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.05);
}

.detail-grid {
  display: grid;
  grid-template-columns: 1fr;
  gap: 1rem;
}

.detail-item .label {
  font-size: 0.8rem;
  color: var(--text-muted);
  margin-bottom: 0.15rem;
  font-weight: bold;
}

.detail-item .value {
  font-size: 0.9rem;
  color: var(--text-main);
}

.code-value {
  font-weight: bold;
  font-family: monospace;
}

.muted-value {
  font-size: 0.85rem;
  color: var(--text-muted);
}

.value-block {
  background: #f1f5f9;
  padding: 0.6rem;
  border-radius: 6px;
  font-size: 0.85rem;
  white-space: pre-wrap;
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

.stars {
  color: var(--star-color);
}

.preview-panel {
  width: 100%;
  padding: 1rem;
  background: #475569;
  border-radius: 8px;
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
  min-height: 300px;
}

.preview-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  color: white;
  font-size: 0.9rem;
}

.pdf-viewer {
  flex: 1;
  min-height: 400px;
}

.pdf-container {
  display: flex;
  flex-direction: column;
  height: 100%;
}

.pdf-toolbar {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 1rem;
  padding: 0.5rem;
  background: #334155;
  border-radius: 4px 4px 0 0;
}

.page-info {
  color: white;
  font-size: 0.85rem;
  min-width: 60px;
  text-align: center;
}

.pdf-canvas-wrapper {
  flex: 1;
  overflow: auto;
  background: #52525b;
  border-radius: 0 0 4px 4px;
  padding: 0.5rem;
  display: flex;
  justify-content: center;
  position: relative;
  min-height: 400px;
}

.pdf-embed {
  display: flex;
  justify-content: center;
}

.pdf-embed :deep(canvas) {
  max-width: 100%;
  height: auto;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.3);
}

.pdf-loading-overlay {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 0.75rem;
  background: #52525b;
  color: #cbd5e1;
  font-size: 0.9rem;
  z-index: 10;
}

.loading-icon {
  font-size: 2rem;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}

.pdf-empty {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 200px;
  background: white;
  border-radius: 4px;
  color: var(--text-muted);
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

  .detail-body {
    flex-direction: row;
  }

  .details-panel {
    width: 45%;
    overflow-y: auto;
    border-right: 1px solid var(--border-color);
    padding: 1.5rem;
  }

  .preview-panel {
    width: 55%;
    padding: 1.5rem;
  }
}
</style>
