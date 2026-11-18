<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import DefaultLayout from '@/layouts/DefaultLayout.vue'
import FormInput from '@/components/FormInput.vue'
import FormAlert from '@/components/FormAlert.vue'
import { ticketService } from '@/services/ticket'
import { ArrowLeft, Loader2 } from 'lucide-vue-next'

const router = useRouter()

const form = reactive({
  title: '',
  description: '',
  category: '',
  priority: '',
})

const errors = reactive({
  title: '',
  description: '',
  priority: '',
})

const loading = ref(false)
const serverError = ref('')
const success = ref(false)

function validate(): boolean {
  let valid = true
  errors.title = ''
  errors.description = ''
  errors.priority = ''

  if (!form.title.trim()) {
    errors.title = 'Judul wajib diisi'
    valid = false
  } else if (form.title.trim().length < 5) {
    errors.title = 'Judul minimal 5 karakter'
    valid = false
  }

  if (!form.description.trim()) {
    errors.description = 'Deskripsi wajib diisi'
    valid = false
  } else if (form.description.trim().length < 10) {
    errors.description = 'Deskripsi minimal 10 karakter'
    valid = false
  }

  if (!form.priority) {
    errors.priority = 'Prioritas wajib dipilih'
    valid = false
  }

  return valid
}

async function handleSubmit() {
  if (!validate()) return

  loading.value = true
  serverError.value = ''

  try {
    const ticket = await ticketService.create({
      title: form.title.trim(),
      description: form.description.trim(),
      category: form.category.trim() || undefined,
      priority: form.priority,
    })

    success.value = true

    setTimeout(() => {
      router.push({ name: 'ticket-detail', params: { id: ticket.id } })
    }, 1500)
  } catch (err: any) {
    serverError.value = err.response?.data?.message || 'Gagal membuat tiket'
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <DefaultLayout>
    <div class="max-w-2xl mx-auto">
      <div class="flex items-center gap-4 mb-6">
        <RouterLink
          to="/"
          class="p-2 rounded-md text-gray-500 hover:bg-gray-100 hover:text-gray-900"
        >
          <ArrowLeft :size="20" />
        </RouterLink>
        <h1 class="text-2xl font-bold text-gray-900">Buat Tiket Baru</h1>
      </div>

      <FormAlert v-if="success" type="success" message="Tiket berhasil dibuat! Mengalihkan ke detail tiket..." />

      <FormAlert v-if="serverError" type="error" :message="serverError" />

      <form novalidate @submit.prevent="handleSubmit" class="bg-white p-6 rounded-lg border border-gray-200 space-y-5">
        <FormInput
          v-model="form.title"
          label="Judul"
          placeholder="Masukkan judul tiket"
          :error="errors.title"
          :required="true"
        />

        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">
            Deskripsi <span class="text-red-500">*</span>
          </label>
          <textarea
            v-model="form.description"
            placeholder="Jelaskan masalah Anda secara detail..."
            rows="5"
            class="w-full px-3 py-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            :class="errors.description ? 'border-red-500' : 'border-gray-300'"
          ></textarea>
          <p v-if="errors.description" class="mt-1 text-sm text-red-600">{{ errors.description }}</p>
        </div>

        <FormInput
          v-model="form.category"
          label="Kategori"
          placeholder="contoh: Teknis, Administrasi, Umum"
        />

        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Prioritas</label>
          <select
            v-model="form.priority"
            class="w-full px-3 py-2 border border-gray-300 rounded-md text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            :class="errors.priority ? 'border-red-500' : 'border-gray-300'"
          >
            <option value="" disabled>-- Pilih Prioritas --</option>
            <option value="low">Low</option>
            <option value="medium">Medium</option>
            <option value="high">High</option>
            <option value="urgent">Urgent</option>
          </select>
          <p v-if="errors.priority" class="mt-1 text-sm text-red-600">{{ errors.priority }}</p>
        </div>

        <div class="flex justify-end gap-3 pt-2">
          <RouterLink
            to="/"
            class="px-4 py-2 border border-gray-300 text-gray-700 text-sm font-medium rounded-md hover:bg-gray-50"
          >
            Batal
          </RouterLink>
          <button
            type="submit"
            :disabled="loading"
            class="inline-flex items-center gap-2 px-4 py-2 bg-blue-600 text-white text-sm font-medium rounded-md hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            <Loader2 v-if="loading" :size="16" class="animate-spin" />
            {{ loading ? 'Membuat...' : 'Buat Tiket' }}
          </button>
        </div>
      </form>
    </div>
  </DefaultLayout>
</template>