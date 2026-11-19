<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import DefaultLayout from '@/layouts/DefaultLayout.vue'
import StatusBadge from '@/components/StatusBadge.vue'
import PriorityBadge from '@/components/PriorityBadge.vue'
import FormAlert from '@/components/FormAlert.vue'
import { ticketService, type Ticket } from '@/services/ticket'
import { commentService } from '@/services/comment'
import { userService, type User } from '@/services/user'
import { useAuthStore } from '@/stores/auth'
import { ArrowLeft, Send, Loader2, User as UserIcon, AlertTriangle } from 'lucide-vue-next'

const route = useRoute()
const router = useRouter()
const auth = useAuthStore()

const ticket = ref<Ticket | null>(null)
const loading = ref(false)
const error = ref('')

const commentBody = ref('')
const isInternal = ref(false)
const commentLoading = ref(false)
const commentError = ref('')

const statusLoading = ref(false)
const statusError = ref('')
const statusSuccess = ref('')

const agents = ref<User[]>([])
const selectedAssignee = ref<number | ''>('')
const assignLoading = ref(false)
const assignError = ref('')
const assignSuccess = ref('')

const selectedPriority = ref('')
const priorityLoading = ref(false)
const priorityError = ref('')
const prioritySuccess = ref('')

const ticketId = computed(() => Number(route.params.id))

const allowedTransitions = computed(() => {
  if (!ticket.value) return []
  const status = ticket.value.status
  const role = auth.userRole

  const transitions: { from: string; to: string; label: string; roles: string[] }[] = [
    { from: 'open', to: 'closed', label: 'Tutup Tiket', roles: ['requester', 'admin'] },
    { from: 'resolved', to: 'closed', label: 'Tutup Tiket', roles: ['requester', 'agent', 'admin'] },
    { from: 'resolved', to: 'open', label: 'Buka Kembali', roles: ['requester'] },
    { from: 'open', to: 'in_progress', label: 'Mulai Kerjakan', roles: ['agent', 'admin'] },
    { from: 'in_progress', to: 'resolved', label: 'Selesaikan', roles: ['agent', 'admin'] },
  ]

  return transitions.filter(t => t.from === status && t.roles.includes(role))
})

async function fetchTicket() {
  loading.value = true
  error.value = ''

  try {
    ticket.value = await ticketService.getById(ticketId.value)

    selectedPriority.value = ticket.value.priority
    selectedAssignee.value = ticket.value.assignee_id || ''
  } catch (err: any) {
    error.value = err.response?.data?.message || 'Gagal memuat tiket'
  } finally {
    loading.value = false
  }
}

async function fetchAgents() {
  if (auth.isRequester) return

  try {
    const users = await userService.getAll()
    agents.value = users.filter(u => u.role === 'agent' || u.role === 'admin')
  } catch {
    agents.value = []
  }
}

async function handleAddComment() {
  if (!commentBody.value.trim()) return

  commentLoading.value = true
  commentError.value = ''

  try {
    await commentService.create(ticketId.value, {
      body: commentBody.value.trim(),
      is_internal: isInternal.value || undefined,
    })

    commentBody.value = ''
    isInternal.value = false
    await fetchTicket()
  } catch (err: any) {
    commentError.value = err.response?.data?.message || 'Gagal menambahkan komentar'
  } finally {
    commentLoading.value = false
  }
}

async function handleStatusChange(newStatus: string) {
  statusLoading.value = true
  statusError.value = ''
  statusSuccess.value = ''

  try {
    await ticketService.update(ticketId.value, { status: newStatus })
    statusSuccess.value = 'Status berhasil diubah'
    await fetchTicket()

    setTimeout(() => {
      statusSuccess.value = ''
    }, 3000)
  } catch (err: any) {
    statusError.value = err.response?.data?.message || 'Gagal mengubah status'
  } finally {
    statusLoading.value = false
  }
}

async function handleAssign() {
  if (selectedAssignee.value === '' || selectedAssignee.value === (ticket.value?.assignee_id || '')) return

  assignLoading.value = true
  assignError.value = ''
  assignSuccess.value = ''

  try {
    await ticketService.update(ticketId.value, { assignee_id: Number(selectedAssignee.value) })
    assignSuccess.value = 'Assignee berhasil diubah'
    await fetchTicket()

    setTimeout(() => {
      assignSuccess.value = ''
    }, 3000)
  } catch (err: any) {
    assignError.value = err.response?.data?.message || 'Gagal mengubah assignee'
  } finally {
    assignLoading.value = false
  }
}

async function handleSelfAssign() {
  if (!auth.user) return
  assignLoading.value = true
  assignError.value = ''
  assignSuccess.value = ''
  try {
    await ticketService.update(ticketId.value, { assignee_id: auth.user.id })
    assignSuccess.value = 'Tiket berhasil diambil dan ditugaskan ke Anda'
    await fetchTicket()
    setTimeout(() => {
      assignSuccess.value = ''
    }, 3000)
  } catch (err: any) {
    assignError.value = err.response?.data?.message || 'Gagal mengambil tiket'
  } finally {
    assignLoading.value = false
  }
}

async function handlePriorityChange() {
  if (!selectedPriority.value || selectedPriority.value === ticket.value?.priority) return

  priorityLoading.value = true
  priorityError.value = ''
  prioritySuccess.value = ''

  try {
    await ticketService.update(ticketId.value, { priority: selectedPriority.value })
    prioritySuccess.value = 'Prioritas berhasil diubah'
    await fetchTicket()

    setTimeout(() => {
      prioritySuccess.value = ''
    }, 3000)
  } catch (err: any) {
    priorityError.value = err.response?.data?.message || 'Gagal mengubah prioritas'
  } finally {
    priorityLoading.value = false
  }
}

function formatDateTime(dateStr: string) {
  return new Date(dateStr).toLocaleDateString('id-ID', {
    day: 'numeric',
    month: 'long',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  })
}

onMounted(() => {
  fetchTicket()
  fetchAgents()
})
</script>

<template>
  <DefaultLayout>
    <div class="max-w-4xl mx-auto">
      <button
        @click="router.push({ name: 'dashboard' })"
        class="flex items-center gap-2 text-sm text-gray-600 hover:text-gray-900 mb-6"
      >
        <ArrowLeft :size="16" />
        Kembali ke daftar tiket
      </button>

      <div v-if="loading" class="flex items-center justify-center py-12">
        <Loader2 :size="24" class="animate-spin text-blue-600" />
        <span class="ml-2 text-gray-600">Memuat tiket...</span>
      </div>

      <div v-else-if="error" class="bg-red-50 border border-red-200 text-red-700 p-4 rounded-lg text-sm">
        {{ error }}
      </div>

      <div v-else-if="ticket">
        <div class="bg-white p-6 rounded-lg border border-gray-200 mb-6">
          <div class="flex flex-col sm:flex-row sm:items-start sm:justify-between gap-4 mb-4">
            <div>
              <h1 class="text-xl font-bold text-gray-900 mb-2">
                #{{ ticket.id }} — {{ ticket.title }}
              </h1>
              <div class="flex flex-wrap items-center gap-2">
                <StatusBadge :status="ticket.status" />
                <PriorityBadge :priority="ticket.priority" />
                <span
                  v-if="ticket.is_overdue"
                  class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-bold bg-red-600 text-white"
                >
                  <AlertTriangle :size="13" />
                  SLA Overdue
                </span>
                <span v-if="ticket.category" class="text-xs text-gray-500 bg-gray-100 px-2 py-0.5 rounded-full">
                  {{ ticket.category }}
                </span>
              </div>
            </div>
          </div>

          <p class="text-gray-700 whitespace-pre-wrap mb-4">{{ ticket.description }}</p>

          <div class="grid grid-cols-1 sm:grid-cols-2 gap-4 text-sm text-gray-500 border-t border-gray-100 pt-4">
            <div>
              <span class="font-medium text-gray-700">Pembuat:</span>
              {{ ticket.requester.name }} ({{ ticket.requester.email }})
            </div>
            <div>
              <span class="font-medium text-gray-700">Assignee:</span>
              {{ ticket.assignee?.name || 'Belum ditugaskan' }}
            </div>
            <div>
              <span class="font-medium text-gray-700">Dibuat:</span>
              {{ formatDateTime(ticket.created_at) }}
            </div>
            <div>
              <span class="font-medium text-gray-700">Diperbarui:</span>
              {{ formatDateTime(ticket.updated_at) }}
            </div>
          </div>
        </div>

        <div v-if="!auth.isRequester" class="bg-white p-4 rounded-lg border border-gray-200 mb-6 space-y-4">
          <h2 class="text-sm font-semibold text-gray-700">Kelola Tiket</h2>

          <div v-if="agents.length > 0">
            <label class="block text-sm text-gray-600 mb-1">Assign ke</label>
            <FormAlert v-if="assignSuccess" type="success" :message="assignSuccess" />
            <FormAlert v-if="assignError" type="error" :message="assignError" />
            <div class="flex gap-2">
              <select
                v-model="selectedAssignee"
                class="flex-1 px-3 py-2 border border-gray-300 rounded-md text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              >
                <option value="">-- Pilih Agent --</option>
                <option v-for="agent in agents" :key="agent.id" :value="agent.id">
                  {{ agent.name }} ({{ agent.role }})
                </option>
              </select>
              <button
                @click="handleAssign"
                :disabled="assignLoading || selectedAssignee === '' || selectedAssignee === (ticket.assignee_id || '')"
                class="px-4 py-2 bg-blue-600 text-white text-sm font-medium rounded-md hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed"
              >
                {{ assignLoading ? 'Menyimpan...' : 'Assign' }}
              </button>
            </div>
          </div>

          <!-- untuk agent sendiri -->
           <div v-else-if="auth.isAgent && ticket.assignee_id !== auth.user?.id">
            <FormAlert v-if="assignSuccess" type="success" :message="assignSuccess" />
            <FormAlert v-if="assignError" type="error" :message="assignError" />
            <div class="flex items-center justify-between p-3 bg-blue-50 border border-blue-200 rounded-md">
              <div>
                <p class="text-sm font-medium text-blue-900">Tiket belum Anda tangani</p>
                <p class="text-xs text-blue-600">Ambil tiket ini untuk menugaskannya ke akun Anda.</p>
              </div>
              <button
                @click="handleSelfAssign"
                :disabled="assignLoading"
                class="px-3 py-1.5 bg-blue-600 text-white text-xs font-semibold rounded-md hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed"
              >
                {{ assignLoading ? 'Memproses...' : 'Ambil Tiket Ini' }}
              </button>
            </div>
          </div>

          <div>
            <label class="block text-sm text-gray-600 mb-1">Ubah Prioritas</label>
            <FormAlert v-if="prioritySuccess" type="success" :message="prioritySuccess" />
            <FormAlert v-if="priorityError" type="error" :message="priorityError" />
            <div class="flex gap-2">
              <select
                v-model="selectedPriority"
                class="flex-1 px-3 py-2 border border-gray-300 rounded-md text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              >
                <option value="low">Low</option>
                <option value="medium">Medium</option>
                <option value="high">High</option>
                <option value="urgent">Urgent</option>
              </select>
              <button
                @click="handlePriorityChange"
                :disabled="priorityLoading || selectedPriority === ticket.priority"
                class="px-4 py-2 bg-blue-600 text-white text-sm font-medium rounded-md hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed"
              >
                {{ priorityLoading ? 'Menyimpan...' : 'Ubah' }}
              </button>
            </div>
          </div>
        </div>

        <div v-if="allowedTransitions.length > 0" class="bg-white p-4 rounded-lg border border-gray-200 mb-6">
          <h2 class="text-sm font-medium text-gray-700 mb-3">Ubah Status</h2>

          <FormAlert v-if="statusSuccess" type="success" :message="statusSuccess" />
          <FormAlert v-if="statusError" type="error" :message="statusError" />

          <div class="flex flex-wrap gap-2">
            <button
              v-for="transition in allowedTransitions"
              :key="transition.to"
              @click="handleStatusChange(transition.to)"
              :disabled="statusLoading"
              class="px-4 py-2 text-sm font-medium rounded-md border disabled:opacity-50 disabled:cursor-not-allowed"
              :class="{
                'border-red-300 text-red-700 hover:bg-red-50': transition.to === 'closed',
                'border-blue-300 text-blue-700 hover:bg-blue-50': transition.to === 'open',
                'border-yellow-300 text-yellow-700 hover:bg-yellow-50': transition.to === 'in_progress',
                'border-green-300 text-green-700 hover:bg-green-50': transition.to === 'resolved',
              }"
            >
              {{ transition.label }}
            </button>
          </div>
        </div>

        <div class="bg-white rounded-lg border border-gray-200">
          <div class="p-4 border-b border-gray-200">
            <h2 class="text-lg font-semibold text-gray-900">
              Komentar ({{ ticket.comments?.length || 0 }})
            </h2>
          </div>

          <div class="divide-y divide-gray-100">
            <div v-if="!ticket.comments || ticket.comments.length === 0" class="p-6 text-center text-gray-500 text-sm">
              Belum ada komentar
            </div>

            <div
              v-for="comment in ticket.comments"
              :key="comment.id"
              class="p-4"
              :class="comment.is_internal ? 'bg-yellow-50' : ''"
            >
              <div class="flex items-start gap-3">
                <div class="w-8 h-8 rounded-full bg-gray-200 flex items-center justify-center shrink-0">
                  <UserIcon :size="16" class="text-gray-500" />
                </div>
                <div class="flex-1 min-w-0">
                  <div class="flex flex-wrap items-center gap-2 mb-1">
                    <span class="text-sm font-medium text-gray-900">{{ comment.author.name }}</span>
                    <span class="text-xs text-gray-400">{{ formatDateTime(comment.created_at) }}</span>
                    <span
                      v-if="comment.is_internal"
                      class="text-xs font-medium px-2 py-0.5 rounded-full bg-yellow-200 text-yellow-800"
                    >
                      Internal Note
                    </span>
                  </div>
                  <p class="text-sm text-gray-700 whitespace-pre-wrap">{{ comment.body }}</p>
                </div>
              </div>
            </div>
          </div>

          <div class="p-4 border-t border-gray-200">
            <FormAlert v-if="commentError" type="error" :message="commentError" />

            <form @submit.prevent="handleAddComment">
              <div class="flex gap-3">
                <textarea
                  v-model="commentBody"
                  placeholder="Tulis komentar..."
                  rows="2"
                  class="flex-1 px-3 py-2 border border-gray-300 rounded-md text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent resize-none"
                ></textarea>

                <button
                  type="submit"
                  :disabled="commentLoading || !commentBody.trim()"
                  class="self-end px-4 py-2 bg-blue-600 text-white text-sm font-medium rounded-md hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed inline-flex items-center gap-2"
                >
                  <Loader2 v-if="commentLoading" :size="16" class="animate-spin" />
                  <Send v-else :size="16" />
                </button>
              </div>

              <label v-if="!auth.isRequester" class="inline-flex items-center gap-2 mt-2 cursor-pointer">
                <input
                  type="checkbox"
                  v-model="isInternal"
                  class="rounded border-gray-300 text-yellow-600 focus:ring-yellow-500"
                />
                <span class="text-sm text-gray-600">Catatan internal (hanya terlihat oleh agent dan admin)</span>
              </label>
            </form>
          </div>
        </div>
      </div>
    </div>
  </DefaultLayout>
</template>