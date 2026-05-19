<script setup lang="ts">
import { onMounted, ref, watch } from 'vue'
import { useRouter } from 'vue-router'
import { useCashTransactionStore } from '@/stores/cashTransaction.store'
import { useCoaStore } from '@/stores/coa.store'
import { useSearch } from '@/composables/common/useSearch'
import { usePagination } from '@/composables/common/usePagination'
import { useFormatter } from '@/composables/common/useFormatter'
import AppLayout from '@/layouts/AppLayout.vue'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Badge } from '@/components/ui/badge'
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
import { MoreHorizontal, Plus, Search, Wallet } from 'lucide-vue-next'
import { TableSkeleton } from '@/components/common/skeleton'
import { CommonPagination } from '@/components/common/pagination'
import { CommonEmpty } from '@/components/common/empty'
import { MEDIUM_SIZE } from '@/constant/pagination.constant'
import type { CashTransaction } from '@/types/cashTransaction.types'
import CashTransactionDeleteDialog from './CashTransactionDeleteDialog.vue'

const router = useRouter()
const cashTransactionStore = useCashTransactionStore()
const coaStore = useCoaStore()
const { formatDate, formatRupiah } = useFormatter()

const isDeleteOpen = ref(false)
const selectedTransaction = ref<CashTransaction | null>(null)

const { page, pageSize, setMeta, goToPage } = usePagination(MEDIUM_SIZE)
const { search } = useSearch(() => {
  page.value = 1
  loadData()
})

async function loadData() {
  await cashTransactionStore.fetchPaginated({
    pageNumber: page.value,
    pageSize: pageSize.value,
    search: search.value || undefined,
  })
  if (cashTransactionStore.meta) {
    setMeta(cashTransactionStore.meta)
  }
}

watch(page, loadData)

onMounted(async () => {
  await Promise.all([loadData(), coaStore.fetchAll()])
})

function getCoaName(coaId: string) {
  return coaStore.coas.find((c) => c.id === coaId)?.name || '-'
}

function handleAdd() {
  router.push({ name: 'cash-transaction-create' })
}

function handleEdit(transaction: CashTransaction) {
  router.push({ name: 'cash-transaction-edit', params: { id: transaction.id } })
}

function handleDelete(transaction: CashTransaction) {
  selectedTransaction.value = transaction
  isDeleteOpen.value = true
}
</script>

<template>
  <AppLayout>
    <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 mb-6">
      <h1 class="text-2xl font-bold">Cash Transactions</h1>
      <div class="flex items-center gap-2 w-full sm:max-w-sm">
        <div class="relative flex-1">
          <Search class="absolute left-2.5 top-2.5 h-4 w-4 text-muted-foreground" />
          <Input v-model="search" placeholder="Search transactions..." class="pl-9" />
        </div>
        <Button @click="handleAdd">
          <Plus class="size-4 mr-2" />
          Add
        </Button>
      </div>
    </div>

    <TableSkeleton
      v-if="cashTransactionStore.loading && cashTransactionStore.transactions.length === 0"
      :column-count="6"
    />

    <CommonEmpty
      v-else-if="!cashTransactionStore.loading && cashTransactionStore.transactions.length === 0"
      title="Cash Transactions"
      description="No cash transactions found. Click add to create one."
      :icon="Wallet"
      :search="search"
      @add="handleAdd"
      @clear-search="search = ''"
    />

    <div v-else class="space-y-4">
      <div class="overflow-x-auto">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>Date</TableHead>
              <TableHead>Description</TableHead>
              <TableHead>Category</TableHead>
              <TableHead>Type</TableHead>
              <TableHead>Amount</TableHead>
              <TableHead class="w-12"></TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            <TableRow
              v-for="transaction in cashTransactionStore.transactions"
              :key="transaction.id"
            >
              <TableCell>{{ formatDate(transaction.date) }}</TableCell>
              <TableCell class="font-medium max-w-75 truncate">
                {{ transaction.description }}
              </TableCell>
              <TableCell>{{ transaction.coaName }}</TableCell>
              <TableCell>
                <Badge
                  :variant="transaction.type === 'in' ? 'success' : 'destructive'"
                  class="capitalize"
                >
                  {{ transaction.type }}
                </Badge>
              </TableCell>
              <TableCell :class="transaction.type === 'in' ? 'text-success' : 'text-destructive'">
                {{ formatRupiah(transaction.amount) }}
              </TableCell>
              <TableCell>
                <DropdownMenu>
                  <DropdownMenuTrigger as-child>
                    <Button variant="ghost" size="icon">
                      <MoreHorizontal class="size-4" />
                    </Button>
                  </DropdownMenuTrigger>
                  <DropdownMenuContent align="end">
                    <DropdownMenuItem @click="handleEdit(transaction)"> Edit </DropdownMenuItem>
                    <DropdownMenuItem
                      class="text-destructive focus:text-destructive"
                      @click="handleDelete(transaction)"
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
        v-if="cashTransactionStore.meta"
        :page="page"
        :page-size="pageSize"
        :total-rows="cashTransactionStore.meta.totalRows"
        :total-pages="cashTransactionStore.meta.totalPages"
        @update:page="goToPage"
      />
    </div>

    <CashTransactionDeleteDialog
      v-model:open="isDeleteOpen"
      :transaction="selectedTransaction"
      @success="loadData"
    />
  </AppLayout>
</template>
