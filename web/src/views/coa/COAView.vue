<script setup lang="ts">
import { onMounted, ref, watch } from 'vue'
import { useCoaStore } from '@/stores/coa.store'
import AppLayout from '@/layouts/AppLayout.vue'
import { Button } from '@/components/ui/button'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import { MoreHorizontal, Plus, Search, BookOpen } from 'lucide-vue-next'
import type { Coa } from '@/types/coa.types'
import COAFormDialog from './COAFormDialog.vue'
import COADeleteDialog from './COADeleteDialog.vue'
import { useSearch } from '@/composables/common/useSearch'
import { usePagination } from '@/composables/common/usePagination'
import { useFormatter } from '@/composables/common/useFormatter'
import { Input } from '@/components/ui/input'
import { CommonPagination } from '@/components/common/pagination'
import { SMALL_SIZE } from '@/constant/pagination.constant'
import { CommonEmpty } from '@/components/common/empty'
import { Badge } from '@/components/ui/badge'
import { TableSkeleton } from '@/components/common/skeleton'

const coaStore = useCoaStore()
const { formatDate } = useFormatter()

const isFormOpen = ref(false)
const isDeleteOpen = ref(false)
const selectedCoa = ref<Coa | null>(null)

const { page, pageSize, setMeta, goToPage } = usePagination(SMALL_SIZE)
const { search } = useSearch(() => {
  page.value = 1
  loadData()
})

async function loadData() {
  await coaStore.fetchPaginated({
    pageNumber: page.value,
    pageSize: pageSize.value,
    search: search.value || undefined,
  })
  if (coaStore.meta) {
    setMeta(coaStore.meta)
  }
}

watch(page, () => {
  loadData()
})

onMounted(() => {
  loadData()
})

function handleAdd() {
  selectedCoa.value = null
  isFormOpen.value = true
}

function handleEdit(coa: Coa) {
  selectedCoa.value = coa
  isFormOpen.value = true
}

function handleDelete(coa: Coa) {
  selectedCoa.value = coa
  isDeleteOpen.value = true
}
</script>

<template>
  <AppLayout>
    <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 mb-6">
      <h1 class="text-2xl font-bold">Chart of Accounts</h1>
      <div class="flex items-center gap-2 w-full sm:max-w-sm">
        <div class="relative flex-1">
          <Search class="absolute left-2.5 top-2.5 h-4 w-4 text-muted-foreground" />
          <Input v-model="search" placeholder="Search accounts..." class="pl-9" />
        </div>
        <Button @click="handleAdd">
          <Plus class="size-4 mr-2" />
          Add
        </Button>
      </div>
    </div>

    <TableSkeleton
      v-if="coaStore.loading && (!coaStore.coas || coaStore.coas.length === 0)"
      :column-count="4"
    />

    <CommonEmpty
      v-else-if="!coaStore.coas || coaStore.coas.length === 0"
      title="Chart of Accounts"
      description="Start by creating your first account to manage your finances."
      :icon="BookOpen"
      :search="search"
      add-button-text="Create Account"
      @add="handleAdd"
      @clear-search="search = ''"
    />

    <div v-else class="space-y-4">
      <div class="overflow-x-auto">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>Name</TableHead>
              <TableHead>Type</TableHead>
              <TableHead>Classification</TableHead>
              <TableHead>Created At</TableHead>
              <TableHead class="w-12.5"></TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            <TableRow v-for="coa in coaStore.coas" :key="coa.id">
              <TableCell class="font-medium">
                <div class="flex items-center gap-2">
                  {{ coa.name }}
                  <Badge v-if="coa.isSystem" variant="secondary" class="text-[10px] px-1.5 py-0 h-4"
                    >System</Badge
                  >
                </div>
              </TableCell>
              <TableCell>
                <span :class="coa.type === 'in' ? 'text-success' : 'text-destructive'">
                  {{ coa.type === 'in' ? 'Income' : 'Expenses' }}
                </span>
              </TableCell>
              <TableCell>
                <Badge :variant="coa.IsOperational ? 'default' : 'secondary'">
                  {{ coa.IsOperational ? 'Operational' : 'General' }}
                </Badge>
              </TableCell>
              <TableCell>{{ formatDate(coa.createdAt) }}</TableCell>
              <TableCell>
                <DropdownMenu>
                  <DropdownMenuTrigger as-child>
                    <Button variant="ghost" size="icon" :disabled="coa.isSystem">
                      <MoreHorizontal class="size-4" />
                    </Button>
                  </DropdownMenuTrigger>
                  <DropdownMenuContent align="end">
                    <DropdownMenuItem @click="handleEdit(coa)"> Edit </DropdownMenuItem>
                    <DropdownMenuItem
                      class="text-destructive focus:text-destructive"
                      @click="handleDelete(coa)"
                    >
                      Delete
                    </DropdownMenuItem>
                  </DropdownMenuContent>
                </DropdownMenu>
              </TableCell>
            </TableRow>
          </TableBody>
        </Table>
      </div>

      <CommonPagination
        v-if="coaStore.meta"
        :page="page"
        :page-size="pageSize"
        :total-rows="coaStore.meta.totalRows"
        :total-pages="coaStore.meta.totalPages"
        @update:page="goToPage"
      />
    </div>

    <COAFormDialog v-model:open="isFormOpen" :coa="selectedCoa" @success="loadData" />

    <COADeleteDialog v-model:open="isDeleteOpen" :coa="selectedCoa" @success="loadData" />
  </AppLayout>
</template>
