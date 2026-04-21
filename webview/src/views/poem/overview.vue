<script lang="ts">
import { Component, Vue, toNative } from 'vue-facing-decorator'
import { api } from '@/service/api'

interface HealthCheck {
  status: string
  message: string
}

interface SystemInfo {
  os: string
  hostname: string
  poem_version: string
}

interface IsrvdInfo {
  status: string
  version: string
  services: string[]
}

@Component
class PoemOverview extends Vue {
  // 数据属性
  healthCheck: HealthCheck | null = null
  systemInfo: SystemInfo | null = null
  isrvdInfo: IsrvdInfo | null = null
  loading = true
  error = ''

  // 方法
  async fetchPoemData() {
    this.loading = true
    this.error = ''
    
    try {
      // 并行请求所有数据
      const [healthRes, systemRes, isrvdRes] = await Promise.all([
        fetch('http://localhost:3000/health'),
        fetch('http://localhost:3000/system'),
        fetch('http://localhost:3000/isrvd')
      ])

      if (!healthRes.ok || !systemRes.ok || !isrvdRes.ok) {
        throw new Error('Failed to fetch Poem service data')
      }

      this.healthCheck = await healthRes.json()
      this.systemInfo = await systemRes.json()
      this.isrvdInfo = await isrvdRes.json()
    } catch (err) {
      this.error = err instanceof Error ? err.message : 'Unknown error'
      console.error('Error fetching Poem data:', err)
    } finally {
      this.loading = false
    }
  }

  // 生命周期
  mounted() {
    this.fetchPoemData()
  }
}

export default toNative(PoemOverview)
</script>

<template>
  <div class="page-container">
    <div class="page-header">
      <h1 class="text-2xl font-bold">Poem 服务</h1>
      <button 
        @click="fetchPoemData"
        class="px-4 py-2 bg-blue-500 text-white rounded-lg hover:bg-blue-600 transition-colors"
        :disabled="loading"
      >
        <i class="fas fa-refresh mr-2" :class="{ 'animate-spin': loading }"></i>
        {{ loading ? '加载中...' : '刷新' }}
      </button>
    </div>

    <!-- 错误信息 -->
    <div v-if="error" class="mt-4 p-4 bg-red-50 border border-red-200 rounded-lg text-red-700">
      <i class="fas fa-exclamation-circle mr-2"></i>
      {{ error }}
    </div>

    <!-- 加载状态 -->
    <div v-if="loading" class="mt-8 flex justify-center">
      <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-500"></div>
    </div>

    <!-- 数据展示 -->
    <div v-else class="mt-8 grid grid-cols-1 md:grid-cols-3 gap-6">
      <!-- 健康检查 -->
      <div class="bg-white rounded-xl shadow-sm border border-slate-200 p-6">
        <div class="flex items-center mb-4">
          <div class="w-10 h-10 rounded-lg bg-green-100 flex items-center justify-center">
            <i class="fas fa-heartbeat text-green-500"></i>
          </div>
          <h2 class="ml-3 text-lg font-semibold">健康检查</h2>
        </div>
        <div v-if="healthCheck" class="space-y-2">
          <div class="flex justify-between">
            <span class="text-slate-500">状态:</span>
            <span class="font-medium" :class="healthCheck.status === 'ok' ? 'text-green-600' : 'text-red-600'">
              {{ healthCheck.status }}
            </span>
          </div>
          <div class="flex justify-between">
            <span class="text-slate-500">消息:</span>
            <span class="font-medium">{{ healthCheck.message }}</span>
          </div>
        </div>
      </div>

      <!-- 系统信息 -->
      <div class="bg-white rounded-xl shadow-sm border border-slate-200 p-6">
        <div class="flex items-center mb-4">
          <div class="w-10 h-10 rounded-lg bg-blue-100 flex items-center justify-center">
            <i class="fas fa-server text-blue-500"></i>
          </div>
          <h2 class="ml-3 text-lg font-semibold">系统信息</h2>
        </div>
        <div v-if="systemInfo" class="space-y-2">
          <div class="flex justify-between">
            <span class="text-slate-500">操作系统:</span>
            <span class="font-medium">{{ systemInfo.os }}</span>
          </div>
          <div class="flex justify-between">
            <span class="text-slate-500">主机名:</span>
            <span class="font-medium">{{ systemInfo.hostname }}</span>
          </div>
          <div class="flex justify-between">
            <span class="text-slate-500">Poem 版本:</span>
            <span class="font-medium">{{ systemInfo.poem_version }}</span>
          </div>
        </div>
      </div>

      <!-- Isrvd 信息 -->
      <div class="bg-white rounded-xl shadow-sm border border-slate-200 p-6">
        <div class="flex items-center mb-4">
          <div class="w-10 h-10 rounded-lg bg-purple-100 flex items-center justify-center">
            <i class="fas fa-cloud text-purple-500"></i>
          </div>
          <h2 class="ml-3 text-lg font-semibold">Isrvd 信息</h2>
        </div>
        <div v-if="isrvdInfo" class="space-y-2">
          <div class="flex justify-between">
            <span class="text-slate-500">状态:</span>
            <span class="font-medium" :class="isrvdInfo.status === 'running' ? 'text-green-600' : 'text-red-600'">
              {{ isrvdInfo.status }}
            </span>
          </div>
          <div class="flex justify-between">
            <span class="text-slate-500">版本:</span>
            <span class="font-medium">{{ isrvdInfo.version }}</span>
          </div>
          <div>
            <span class="text-slate-500">服务:</span>
            <div class="mt-2 flex flex-wrap gap-2">
              <span v-for="service in isrvdInfo.services" :key="service" class="px-2 py-1 bg-slate-100 rounded-full text-xs font-medium">
                {{ service }}
              </span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 服务说明 -->
    <div class="mt-8 bg-white rounded-xl shadow-sm border border-slate-200 p-6">
      <h2 class="text-lg font-semibold mb-4">关于 Poem 服务</h2>
      <p class="text-slate-600 mb-4">
        Poem 是一个由 Rust 编程语言构建的全功能、易用的 Web 框架。它将程序的优雅与诗歌的韵律融合在一起，赋予 Web 开发新的生命。
      </p>
      <p class="text-slate-600">
        本集成展示了如何在 Isrvd 项目中添加 Poem 服务支持，通过 Rust 与 Go 的协作，为用户提供更加丰富和高效的服务体验。
      </p>
    </div>
  </div>
</template>

<style scoped>
.page-container {
  @apply p-6 bg-slate-50 min-h-screen;
}

.page-header {
  @apply flex justify-between items-center mb-6;
}
</style>
