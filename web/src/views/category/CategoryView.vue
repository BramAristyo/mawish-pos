<script setup lang="ts">
import { onMounted, ref, watch } from 'vue'
import { useCategoryStore } from '@/stores/category.stores'
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
import { MoreHorizontal, Plus, Search, Tag } from 'lucide-vue-next'
import type { Category } from '@/types/category.types'
import CategoryFormDialog from './CategoryFormDialog.vue'
import CategoryDeleteDialog from './CategoryDeleteDialog.vue'
import { useSearch } from '@/composables/common/useSearch'
import { usePagination } from '@/composables/common/usePagination'
import { useFormatter } from '@/composables/common/useFormatter'
import { Input } from '@/components/ui/input'
import { CommonPagination } from '@/components/common/pagination'
import { SMALL_SIZE } from '@/constant/pagination.constant'
import { CommonEmpty } from '@/components/common/empty'
import { Skeleton } from '@/components/ui/skeleton'

const categoryStore = useCategoryStore()
const { formatDate } = useFormatter()

const isFormOpen = ref(false)
const isDeleteOpen = ref(false)
const selectedCategory = ref<Category | null>(null)

const { page, pageSize, setMeta, goToPage } = usePagination(SMALL_SIZE)
const { search } = useSearch(() => {
  page.value = 1
  loadData()
})

async function loadData() {
  await categoryStore.fetchPaginated({
    pageNumber: page.value,
    pageSize: pageSize.value,
    search: search.value || undefined,
  })
  if (categoryStore.meta) {
    setMeta(categoryStore.meta)
  }
}

watch(page, () => {
  loadData()
})

onMounted(() => {
  loadData()
})

function handleAdd() {
  selectedCategory.value = null
  isFormOpen.value = true
}

function handleEdit(category: Category) {
  selectedCategory.value = category
  isFormOpen.value = true
}

function handleDelete(category: Category) {
  selectedCategory.value = category
  isDeleteOpen.value = true
}
</script>

<template>
  <AppLayout>
    <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 mb-6">
      <h1 class="text-2xl font-bold">Category</h1>
      <div class="flex items-center gap-2 w-full sm:max-w-sm">
        <div class="relative flex-1">
          <Search class="absolute left-2.5 top-2.5 h-4 w-4 text-muted-foreground" />
          <Input v-model="search" placeholder="Search categories..." class="pl-9" />
        </div>
        <Button @click="handleAdd">
          <Plus class="size-4 mr-2" />
          Add
        </Button>
      </div>
    </div>

    <TableSkeleton v-if="categoryStore.loading && categoryStore.categories.length === 0" :column-count="4" />

    <CommonEmpty
      v-else-if="categoryStore.categories.length === 0"
      title="Categories"
      description="Start by creating your first category to organize your products."
      :icon="Tag"
      :search="search"
      add-button-text="Create Category"
      @add="handleAdd"
      @clear-search="search = ''"
    />

    <div v-else class="space-y-4">
      <div class="overflow-x-auto">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>Name</TableHead>
              <TableHead>Description</TableHead>
              <TableHead>Created At</TableHead>
              <TableHead class="w-12.5"></TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            <TableRow v-for="category in categoryStore.categories" :key="category.id">
              <TableCell class="font-medium">{{ category.name }}</TableCell>
              <TableCell>{{ category.description }}</TableCell>
              <TableCell>{{ formatDate(category.createdAt) }}</TableCell>
              <TableCell>
                <DropdownMenu>
                  <DropdownMenuTrigger as-child>
                    <Button variant="ghost" size="icon">
                      <MoreHorizontal class="size-4" />
                    </Button>
                  </DropdownMenuTrigger>
                  <DropdownMenuContent align="end">
                    <DropdownMenuItem @click="handleEdit(category)"> Edit </DropdownMenuItem>
                    <DropdownMenuItem
                      class="text-destructive focus:text-destructive"
                      @click="handleDelete(category)"
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
        v-if="categoryStore.meta"
        :page="page"
        :page-size="pageSize"
        :total-rows="categoryStore.meta.totalRows"
        :total-pages="categoryStore.meta.totalPages"
        @update:page="goToPage"
      />
    </div>

    <CategoryFormDialog
      v-model:open="isFormOpen"
      :category="selectedCategory"
      @success="loadData"
    />

    <CategoryDeleteDialog
      v-model:open="isDeleteOpen"
      :category="selectedCategory"
      @success="loadData"
    />
  </AppLayout>
</template>
