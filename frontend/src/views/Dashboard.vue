<script setup lang="ts">
import { ref, watch, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import DefaultLayout from '@/layouts/DefaultLayout.vue'
import StatusBadge from '@/components/StatusBadge.vue'
import PriorityBadge from '@/components/PriorityBadge.vue'
import Pagination from '@/components/Pagination.vue'
import { ticketService, type TicketFilter, type Ticket } from '@/services/ticket'
import { Plus, Search, Loader2 } from 'lucide-vue-next'

const router = useRouter()

const tickets = ref<Ticket[]>([])
const loading = ref(false)
const error = ref('')

const page = ref(1)
const totalPages = ref(1)
const total = ref(0)

const statusFilter = ref('')
const priorityFilter = ref('')
const searchQuery = ref('')

async function fetchTickets() {
  loading.value = true
  error.value = ''

  try {
    const filter: TicketFilter = {
      page: page.value,
      limit: 10,
    }

    if (statusFilter.value) filter.status = statusFilter.value
    if (priorityFilter.value) filter.priority = priorityFilter.value
    if (searchQuery.value) filter.search = searchQuery.value

    const result = await ticketService.getAll(filter)
    tickets.value = result.tickets || []

    total.value = result.total
    totalPages.value = result.total_pages
  } catch (err: any) {
    error.value = err.response?.data?.message || 'Gagal memuat daftar tiket'
  } finally {
    loading.value = false
  }
}

function onSearch() {
  page.value = 1
  fetchTickets()
}

function onFilterChange() {
  page.value = 1
  fetchTickets()
}

function goToDetail(id: number) {
  router.push({ name: 'ticket-detail', params: { id } })
}

function formatDate(dateStr: string) {
  return new Date(dateStr).toLocaleDateString('id-ID', {
    day: 'numeric',
    month: 'short',
    year: 'numeric',
  })
}

watch(page, fetchTickets)

onMounted(fetchTickets)
</script>

<template>
  <DefaultLayout>
    <div>
      <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4 mb-6">
        <h1 class="text-2xl font-bold text-gray-900">Daftar Tiket</h1>

        <RouterLink
          to="/tickets/create"
          class="inline-flex items-center gap-2 px-4 py-2 bg-blue-600 text-white text-sm font-medium rounded-md hover:bg-blue-700"
        >
          <Plus :size="16" />
          Buat Tiket
        </RouterLink>
      </div>

      <div class="bg-white p-4 rounded-lg border border-gray-200 mb-6">
        <div class="grid grid-cols-1 sm:grid-cols-3 gap-4">
          <div class="relative">
            <Search :size="16" class="absolute left-3 top-1/2 -translate-y-1/2 text-gray-400" />
            <input
              v-model="searchQuery"
              @keyup.enter="onSearch"
              type="text"
              placeholder="Cari judul tiket..."
              class="w-full pl-9 pr-3 py-2 border border-gray-300 rounded-md text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            />
          </div>

          <select
            v-model="statusFilter"
            @change="onFilterChange"
            class="w-full px-3 py-2 border border-gray-300 rounded-md text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
          >
            <option value="">Semua Status</option>
            <option value="open">Open</option>
            <option value="in_progress">In Progress</option>
            <option value="resolved">Resolved</option>
            <option value="closed">Closed</option>
          </select>

          <select
            v-model="priorityFilter"
            @change="onFilterChange"
            class="w-full px-3 py-2 border border-gray-300 rounded-md text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
          >
            <option value="">Semua Priority</option>
            <option value="low">Low</option>
            <option value="medium">Medium</option>
            <option value="high">High</option>
            <option value="urgent">Urgent</option>
          </select>
        </div>
      </div>

      <div v-if="error" class="bg-red-50 border border-red-200 text-red-700 p-4 rounded-lg mb-6 text-sm">
        {{ error }}
      </div>

      <div v-if="loading" class="flex items-center justify-center py-12">
        <Loader2 :size="24" class="animate-spin text-blue-600" />
        <span class="ml-2 text-gray-600">Memuat tiket...</span>
      </div>

      <div v-else-if="tickets.length === 0" class="text-center py-12 bg-white rounded-lg border border-gray-200">
        <p class="text-gray-500 text-lg mb-2">Belum ada tiket</p>
        <p class="text-gray-400 text-sm">Buat tiket pertama Anda dengan menekan tombol "Buat Tiket"</p>
      </div>

      <div v-else class="bg-white rounded-lg border border-gray-200 overflow-hidden">
        <div class="overflow-x-auto">
          <table class="min-w-full divide-y divide-gray-200">
            <thead class="bg-gray-50">
              <tr>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">ID</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Judul</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Kategori</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Status</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Prioritas</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Tanggal</th>
              </tr>
            </thead>
            <tbody class="bg-white divide-y divide-gray-200">
              <tr
                v-for="ticket in tickets"
                :key="ticket.id"
                @click="goToDetail(ticket.id)"
                class="hover:bg-gray-50 cursor-pointer"
              >
                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">#{{ ticket.id }}</td>
                <td class="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">{{ ticket.title }}</td>
                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{{ ticket.category || '-' }}</td>
                <td class="px-6 py-4 whitespace-nowrap">
                  <StatusBadge :status="ticket.status" />
                </td>
                <td class="px-6 py-4 whitespace-nowrap">
                  <PriorityBadge :priority="ticket.priority" />
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{{ formatDate(ticket.created_at) }}</td>
              </tr>
            </tbody>
          </table>
        </div>

        <div class="px-6 py-4">
          <Pagination
            :page="page"
            :total-pages="totalPages"
            :total="total"
            @update:page="page = $event"
          />
        </div>
      </div>
    </div>
  </DefaultLayout>
</template>