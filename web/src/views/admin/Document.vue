<template>
  <div class="document-manage">
    <!-- 工具栏 -->
    <div class="toolbar">
      <el-button type="primary" @click="handleAdd">
        <el-icon><Plus /></el-icon> 新增文档
      </el-button>
      <div class="search-area">
        <el-input
          v-model="keyword"
          placeholder="搜索编号、名称..."
          clearable
          class="search-input"
          @keyup.enter="loadData"
        />
        <el-select v-model="stageFilter" placeholder="当前阶段" clearable class="search-select" @change="loadData">
          <el-option v-for="s in stageOptions" :key="s" :label="s" :value="s" />
        </el-select>
        <el-button type="primary" @click="loadData">搜索</el-button>
      </div>
    </div>

    <!-- PC端：文档表格 -->
    <el-table :data="documents" v-loading="loading" stripe border class="pc-table">
      <el-table-column prop="iso_code" label="标准编号" width="140" />
      <el-table-column prop="name" label="标准名称" min-width="200" show-overflow-tooltip />
      <el-table-column prop="type" label="类型" width="80" align="center" />
      <el-table-column prop="current_stage" label="当前阶段" width="100" align="center" />
      <el-table-column prop="belongs_to" label="标准所属" width="120" show-overflow-tooltip />
      <el-table-column label="地震相关度" width="120" align="center">
        <template #default="{ row }">
          <span class="stars">{{ getStars(row.earthquake_relevance) }}</span>
        </template>
      </el-table-column>
      <el-table-column label="发布日期" width="120" align="center">
        <template #default="{ row }">
          {{ row.publish_date ? row.publish_date.substring(0, 10) : '-' }}
        </template>
      </el-table-column>
      <el-table-column label="附件" width="80" align="center">
        <template #default="{ row }">
          <el-icon v-if="row.attachment" color="#10b981"><Check /></el-icon>
          <el-icon v-else color="#94a3b8"><Close /></el-icon>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="180" fixed="right" align="center">
        <template #default="{ row }">
          <el-button size="small" @click="handleEdit(row)">编辑</el-button>
          <el-button size="small" type="danger" @click="handleDelete(row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <!-- 手机端：卡片列表 -->
    <div class="mobile-list" v-loading="loading">
      <div
        v-for="doc in documents"
        :key="doc.id"
        class="mobile-card"
      >
        <div class="card-badge-row">
          <span class="badge badge-code">{{ doc.iso_code }}</span>
          <span class="badge badge-type">{{ doc.type }}</span>
          <span class="badge badge-stage">{{ doc.current_stage }}</span>
        </div>
        <div class="card-title">{{ doc.name }}</div>
        <div class="card-meta">
          <div class="meta-pair">
            <span>标准所属: {{ doc.standard_belongs_to }}</span>
            <span>发布: {{ doc.publish_date ? doc.publish_date.substring(0, 10) : '-' }}</span>
          </div>
          <div class="meta-pair">
            <span>首发编号: {{ doc.first_publish_code || '-' }}</span>
            <span>所属: {{ doc.belongs_to || '-' }}</span>
          </div>
          <div>相关度: <span class="stars">{{ getStars(doc.earthquake_relevance) }}</span></div>
        </div>
        <div class="card-actions">
          <el-button size="small" @click="handleEdit(doc)">编辑</el-button>
          <el-button size="small" type="danger" @click="handleDelete(doc)">删除</el-button>
        </div>
      </div>
      <el-empty v-if="!loading && documents.length === 0" description="暂无文档" />
    </div>

    <!-- 分页 -->
    <div class="pagination-wrapper">
      <el-pagination
        v-model:current-page="page"
        :page-size="pageSize"
        :total="total"
        layout="total, prev, pager, next"
        @current-change="loadData"
      />
    </div>

    <!-- 新增/编辑弹窗 -->
    <el-dialog
      v-model="dialogVisible"
      :title="isEdit ? '编辑文档' : '新增文档'"
      width="700px"
      :close-on-click-modal="false"
      destroy-on-close
      class="doc-dialog"
    >
      <el-form :model="form" :rules="formRules" ref="formRef" label-position="top" class="doc-form">
        <el-form-item label="国际标准编号" prop="iso_code">
          <el-input v-model="form.iso_code" placeholder="如 ISOxxxx" />
        </el-form-item>
        <el-form-item label="国际标准名称" prop="name">
          <el-input v-model="form.name" placeholder="请输入标准名称" />
        </el-form-item>
        <el-row :gutter="20">
          <el-col :xs="24" :sm="12">
            <el-form-item label="标准类型" prop="type">
              <el-select v-model="form.type" placeholder="选择类型" style="width: 100%">
                <el-option v-for="t in typeOptions" :key="t" :label="t" :value="t" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :xs="24" :sm="12">
            <el-form-item label="当前阶段" prop="current_stage">
              <el-select v-model="form.current_stage" placeholder="选择阶段" style="width: 100%">
                <el-option v-for="s in stageOptions" :key="s" :label="s" :value="s" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item label="地震相关度" prop="earthquake_relevance">
          <el-rate v-model="form.earthquake_relevance" :max="5" />
        </el-form-item>
        <el-row :gutter="20">
          <el-col :xs="24" :sm="12">
            <el-form-item label="标准所属">
              <el-input v-model="form.standard_belongs_to" placeholder="ISO/TC" />
            </el-form-item>
          </el-col>
          <el-col :xs="24" :sm="12">
            <el-form-item label="所属">
              <el-select v-model="form.belongs_to" placeholder="选择所属" clearable style="width: 100%">
                <el-option v-for="b in belongsToOptions" :key="b" :label="b" :value="b" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :xs="24" :sm="12">
            <el-form-item label="发布日期">
              <el-date-picker
                v-model="form.publish_date"
                type="date"
                placeholder="选择日期"
                format="YYYY-MM-DD"
                value-format="YYYY-MM-DD"
                style="width: 100%"
              />
            </el-form-item>
          </el-col>
          <el-col :xs="24" :sm="12">
            <el-form-item label="首次发布编号">
              <el-input v-model="form.first_publish_code" placeholder="首次发布编号" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item label="摘要">
          <el-input v-model="form.summary" type="textarea" :rows="3" placeholder="请输入摘要" />
        </el-form-item>
        <el-form-item label="范围">
          <el-input v-model="form.scope" type="textarea" :rows="3" placeholder="请输入范围" />
        </el-form-item>
        <el-form-item label="附件PDF">
          <el-upload
            ref="uploadRef"
            :auto-upload="false"
            :limit="1"
            accept=".pdf"
            :on-change="handleFileChange"
            :file-list="fileList"
          >
            <el-button size="small" type="primary">选择文件</el-button>
            <template #tip>
              <div class="el-upload__tip">只能上传PDF文件</div>
            </template>
          </el-upload>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitLoading" @click="handleSubmit">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { getDocumentList, createDocument, updateDocument, deleteDocument } from '../../api/admin'

const typeOptions = ['IS', 'TS', 'PAS', 'TR', 'IWA', 'Guides']
const stageOptions = ['PWI', 'NP', 'WD', 'CD', 'DIS', 'FDIS', 'IS']
const belongsToOptions = ['SC', 'WG']

const documents = ref([])
const loading = ref(false)
const keyword = ref('')
const stageFilter = ref('')
const page = ref(1)
const pageSize = ref(10)
const total = ref(0)

const dialogVisible = ref(false)
const isEdit = ref(false)
const editId = ref(null)
const submitLoading = ref(false)
const formRef = ref(null)
const uploadRef = ref(null)
const uploadFile = ref(null)
const fileList = ref([])

const form = reactive({
  iso_code: '',
  name: '',
  type: '',
  standard_belongs_to: 'ISO/TC',
  belongs_to: '',
  summary: '',
  scope: '',
  publish_date: '',
  first_publish_code: '',
  current_stage: '',
  earthquake_relevance: 3
})

const formRules = {
  iso_code: [{ required: true, message: '请输入国际标准编号', trigger: 'blur' }],
  name: [{ required: true, message: '请输入国际标准名称', trigger: 'blur' }],
  type: [{ required: true, message: '请选择标准类型', trigger: 'change' }],
  current_stage: [{ required: true, message: '请选择当前阶段', trigger: 'change' }]
}

function getStars(count) {
  return '★'.repeat(count || 0) + '☆'.repeat(5 - (count || 0))
}

async function loadData() {
  loading.value = true
  try {
    const res = await getDocumentList({
      keyword: keyword.value,
      stage: stageFilter.value,
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

function resetForm() {
  form.iso_code = ''
  form.name = ''
  form.type = ''
  form.standard_belongs_to = 'ISO/TC'
  form.belongs_to = ''
  form.summary = ''
  form.scope = ''
  form.publish_date = ''
  form.first_publish_code = ''
  form.current_stage = ''
  form.earthquake_relevance = 3
  uploadFile.value = null
  fileList.value = []
}

function handleAdd() {
  resetForm()
  isEdit.value = false
  editId.value = null
  dialogVisible.value = true
}

function handleEdit(row) {
  resetForm()
  isEdit.value = true
  editId.value = row.id
  form.iso_code = row.iso_code
  form.name = row.name
  form.type = row.type
  form.standard_belongs_to = row.standard_belongs_to || 'ISO/TC'
  form.belongs_to = row.belongs_to || ''
  form.summary = row.summary || ''
  form.scope = row.scope || ''
  form.publish_date = row.publish_date ? row.publish_date.substring(0, 10) : ''
  form.first_publish_code = row.first_publish_code || ''
  form.current_stage = row.current_stage
  form.earthquake_relevance = row.earthquake_relevance || 3
  if (row.attachment) {
    fileList.value = [{ name: row.attachment, url: '' }]
  }
  dialogVisible.value = true
}

function handleFileChange(file) {
  uploadFile.value = file.raw
}

async function handleSubmit() {
  submitLoading.value = true
  const valid = await formRef.value.validate().catch(() => false)
  if (!valid) {
    submitLoading.value = false
    return
  }

  try {
    const formData = new FormData()
    formData.append('iso_code', form.iso_code.trim())
    formData.append('name', form.name)
    formData.append('type', form.type)
    formData.append('standard_belongs_to', form.standard_belongs_to)
    formData.append('belongs_to', form.belongs_to)
    formData.append('summary', form.summary)
    formData.append('scope', form.scope)
    formData.append('publish_date', form.publish_date || '')
    formData.append('first_publish_code', form.first_publish_code)
    formData.append('current_stage', form.current_stage)
    formData.append('earthquake_relevance', String(form.earthquake_relevance))
    if (uploadFile.value) {
      formData.append('attachment', uploadFile.value)
    }

    let res
    if (isEdit.value) {
      res = await updateDocument(editId.value, formData)
    } else {
      res = await createDocument(formData)
    }

    if (res.code === 0) {
      ElMessage.success(isEdit.value ? '更新成功' : '创建成功')
      dialogVisible.value = false
      loadData()
    } else {
      ElMessage.error(res.message || '操作失败！')
    }
  } catch (e) {
    ElMessage.error('操作失败')
  } finally {
    submitLoading.value = false
  }
}

async function handleDelete(row) {
  try {
    await ElMessageBox.confirm(`确定要删除文档 "${row.iso_code}" 吗？`, '确认删除', {
      type: 'warning'
    })
    const res = await deleteDocument(row.id)
    if (res.code === 0) {
      ElMessage.success('删除成功')
      loadData()
    } else {
      ElMessage.error(res.message || '删除失败')
    }
  } catch {
    // 取消删除
  }
}

onMounted(() => {
  loadData()
})
</script>

<style scoped>
.document-manage {
  background: white;
  border-radius: 8px;
  padding: 1.5rem;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.05);
}

.toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1rem;
  flex-wrap: wrap;
  gap: 0.5rem;
}

.search-area {
  display: flex;
  gap: 0.5rem;
  align-items: center;
  flex-wrap: wrap;
}

.search-input {
  width: 250px;
}

.search-select {
  width: 150px;
}

.stars {
  color: var(--star-color);
}

.pagination-wrapper {
  display: flex;
  justify-content: center;
  margin-top: 1rem;
}

/* 手机端卡片样式 */
.mobile-list {
  display: none;
  flex-direction: column;
  gap: 0.75rem;
}

.mobile-card {
  border: 1px solid var(--border-color);
  border-radius: 8px;
  padding: 1rem;
}

.card-badge-row {
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

.card-title {
  font-size: 1rem;
  font-weight: 600;
  color: #1e293b;
  margin-bottom: 0.5rem;
  line-height: 1.4;
}

.card-meta {
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
  font-size: 0.8rem;
  color: var(--text-muted);
  margin-bottom: 0.75rem;
}

.meta-pair {
  display: flex;
  gap: 1.5rem;
}

.card-actions {
  display: flex;
  gap: 0.5rem;
  justify-content: flex-end;
}

/* 手机端适配 */
@media (max-width: 768px) {
  .document-manage {
    padding: 0.75rem;
  }

  .toolbar {
    flex-direction: column;
    align-items: stretch;
  }

  .search-area {
    flex-direction: column;
  }

  .search-input,
  .search-select {
    width: 100% !important;
  }

  .pc-table {
    display: none;
  }

  .mobile-list {
    display: flex;
  }

  .doc-form :deep(.el-form-item__label) {
    font-size: 0.85rem;
    padding-bottom: 2px;
  }

  .doc-form :deep(.el-col) {
    max-width: 100%;
    flex: 0 0 100%;
  }
}
</style>

<style>
/* 弹窗样式 - 非 scoped，因为 el-dialog 传送到 body */
.doc-dialog .el-dialog__footer {
  text-align: left;
}

@media (max-width: 768px) {
  .doc-dialog .el-dialog {
    width: calc(100vw - 20px) !important;
    margin: 5px auto !important;
  }

  .doc-dialog .el-dialog__header {
    padding: 10px 12px;
  }

  .doc-dialog .el-dialog__body {
    padding: 10px 8px;
    max-height: 72vh;
    overflow-y: auto;
  }

  .doc-dialog .el-dialog__footer {
    padding: 8px 12px;
    text-align: left;
  }
}
</style>
