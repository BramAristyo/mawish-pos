<script setup lang="ts">
import { useRouter } from 'vue-router'
import AppLayout from '@/layouts/AppLayout.vue'
import CashTransactionForm from './CashTransactionForm.vue'

const router = useRouter()

function handleSuccess() {
  router.push({ name: 'cash-transactions' })
}
</script>

<template>
  <AppLayout>
    <CashTransactionForm @success="handleSuccess" />
  </AppLayout>
</template>
