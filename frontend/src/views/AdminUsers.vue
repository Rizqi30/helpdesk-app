<script setup lang="ts">
import { ref, onMounted } from 'vue'
import DefaultLayout from '@/layouts/DefaultLayout.vue'
import FormAlert from '@/components/FormAlert.vue'
import { userService, type User } from '@/services/user'
import { Loader2 } from 'lucide-vue-next'

const users = ref<User[]>([])
const loading = ref(false)
const error = ref('')

const editingUserId = ref<number | null>(null)
const editingRole = ref('')
const roleLoading = ref(false)
const roleError = ref('')
const roleSuccess = ref('')

async function fetchUsers() {
  loading.value = true
  error.value = ''

  try {
    users.value = await userService.getAll()
  } catch (err: any) {
    error.value = err.response?.data?.message || 'Gagal memuat daftar pengguna'
  } finally {
    loading.value = false
  }
}

function startEdit(user: User) {
  editingUserId.value = user.id
  editingRole.value = user.role
  roleError.value = ''
  roleSuccess.value = ''
}

function cancelEdit() {
  editingUserId.value = null
  editingRole.value = ''
}

async function handleSaveRole(userId: number) {
  roleLoading.value = true
  roleError.value = ''
  roleSuccess.value = ''

  try {
    await userService.updateRole(userId, editingRole.value)
    roleSuccess.value = 'Role berhasil diubah'
    editingUserId.value = null
    await fetchUsers()

    setTimeout(() => {
      roleSuccess.value = ''
    }, 3000)
  } catch (err: any) {
    roleError.value = err.response?.data?.message || 'Gagal mengubah role'
  } finally {
    roleLoading.value = false
  }
}

function getRoleBadgeClass(role: string) {
  switch (role) {
    case 'admin':
      return 'bg-red-100 text-red-700'
    case 'agent':
      return 'bg-blue-100 text-blue-700'
    default:
      return 'bg-gray-100 text-gray-700'
  }
}

function formatDate(dateStr: string) {
  return new Date(dateStr).toLocaleDateString('id-ID', {
    day: 'numeric',
    month: 'short',
    year: 'numeric',
  })
}

onMounted(fetchUsers)
</script>

<template>
  <DefaultLayout>
    <div>
      <h1 class="text-2xl font-bold text-gray-900 mb-6">Kelola Pengguna</h1>

      <FormAlert v-if="roleSuccess" type="success" :message="roleSuccess" />
      <FormAlert v-if="roleError" type="error" :message="roleError" />

      <div v-if="error" class="bg-red-50 border border-red-200 text-red-700 p-4 rounded-lg mb-6 text-sm">
        {{ error }}
      </div>

      <div v-if="loading" class="flex items-center justify-center py-12">
        <Loader2 :size="24" class="animate-spin text-blue-600" />
        <span class="ml-2 text-gray-600">Memuat pengguna...</span>
      </div>

      <div v-else-if="users.length === 0" class="text-center py-12 bg-white rounded-lg border border-gray-200">
        <p class="text-gray-500">Belum ada pengguna</p>
      </div>

      <div v-else class="bg-white rounded-lg border border-gray-200 overflow-hidden">
        <div class="overflow-x-auto">
          <table class="min-w-full divide-y divide-gray-200">
            <thead class="bg-gray-50">
              <tr>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">ID</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Nama</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Email</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Role</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Terdaftar</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Aksi</th>
              </tr>
            </thead>
            <tbody class="bg-white divide-y divide-gray-200">
              <tr v-for="user in users" :key="user.id">
                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{{ user.id }}</td>
                <td class="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">{{ user.name }}</td>
                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{{ user.email }}</td>
                <td class="px-6 py-4 whitespace-nowrap">
                  <select
                    v-if="editingUserId === user.id"
                    v-model="editingRole"
                    class="px-2 py-1 border border-gray-300 rounded-md text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                  >
                    <option value="requester">Requester</option>
                    <option value="agent">Agent</option>
                    <option value="admin">Admin</option>
                  </select>

                  <span
                    v-else
                    class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium"
                    :class="getRoleBadgeClass(user.role)"
                  >
                    {{ user.role }}
                  </span>
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{{ formatDate(user.created_at) }}</td>
                <td class="px-6 py-4 whitespace-nowrap text-sm">
                  <div v-if="editingUserId === user.id" class="flex items-center gap-2">
                    <button
                      @click="handleSaveRole(user.id)"
                      :disabled="roleLoading || editingRole === user.role"
                      class="px-3 py-1 bg-blue-600 text-white text-xs font-medium rounded-md hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed"
                    >
                      {{ roleLoading ? 'Menyimpan...' : 'Simpan' }}
                    </button>

                    <button
                      @click="cancelEdit"
                      class="px-3 py-1 border border-gray-300 text-gray-700 text-xs font-medium rounded-md hover:bg-gray-50"
                    >
                      Batal
                    </button>
                  </div>
                  
                  <button
                    v-else
                    @click="startEdit(user)"
                    class="px-3 py-1 border border-gray-300 text-gray-700 text-xs font-medium rounded-md hover:bg-gray-50"
                  >
                    Ubah Role
                  </button>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>
  </DefaultLayout>
</template>