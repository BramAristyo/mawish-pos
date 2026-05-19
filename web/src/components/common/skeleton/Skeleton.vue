<script setup lang="ts">
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'
import { Skeleton } from '@/components/ui/skeleton'

interface Props {
  columnCount: number
  rowCount?: number
}

withDefaults(defineProps<Props>(), {
  rowCount: 5,
})
</script>

<template>
  <div class="border rounded-lg overflow-hidden">
    <Table>
      <TableHeader>
        <TableRow>
          <TableHead v-for="i in columnCount" :key="i">
            <Skeleton class="h-4 w-24" />
          </TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        <TableRow v-for="i in rowCount" :key="i">
          <TableCell v-for="j in columnCount" :key="j">
            <div class="flex items-center" :class="{ 'justify-end': j === columnCount }">
              <Skeleton
                v-if="j === columnCount"
                class="h-8 w-8 rounded-md"
              />
              <Skeleton
                v-else
                :class="[
                  'h-4 w-full',
                  j === 1 ? 'max-w-[120px]' : (i % 2 === 0 ? 'max-w-[180px]' : 'max-w-[150px]')
                ]"
              />
            </div>
          </TableCell>
        </TableRow>
      </TableBody>
    </Table>
  </div>
</template>
