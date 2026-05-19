<script setup lang="ts">
import { computed, reactive, onMounted, ref, watch } from 'vue'
import { usePayrollStore } from '@/stores/payroll.store'
import { useEmployeeStore } from '@/stores/employee.store'
import type { CreatePayrollRequest } from '@/types/payroll.types'
import type { ValidationError } from '@/types/common.types'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { CancelModal } from '@/components/common/cancel'
import { toast } from 'vue-sonner'
import { useFormErrors } from '@/composables/common/useFormErrors'

const props = defineProps<{
  open: boolean
}>()

const emit = defineEmits<{
  (e: 'update:open', value: boolean): void
  (e: 'success'): void
}>()

const payrollStore = usePayrollStore()
const employeeStore = useEmployeeStore()
const { setErrors, clearErrors, getErrorMessage, hasError } = useFormErrors()

const isOpen = computed({
  get: () => props.open,
  set: (value) => emit('update:open', value),
})

const form = reactive<CreatePayrollRequest>({
  employeeID: '',
  periodStart: '',
  periodEnd: '',
})

const initialForm = ref<string>('')
const showCancelModal = ref(false)

const isDirty = computed(() => {
  return JSON.stringify(form) !== initialForm.value
})

function handleCancel() {
  if (isDirty.value) {
    showCancelModal.value = true
  } else {
    isOpen.value = false
  }
}

function confirmCancel() {
  isOpen.value = false
}

function resetForm() {
  clearErrors()
  form.employeeID = ''
  form.periodStart = ''
  form.periodEnd = ''
  initialForm.value = JSON.stringify(form)
}

watch(
  () => props.open,
  (open) => {
    if (open) {
      resetForm()
    }
  },
)

onMounted(async () => {
  await employeeStore.ensureDataLoaded()
})

async function handleSubmit() {
  clearErrors()
  try {
    await payrollStore.create(form)
    toast.success('Payroll created successfully')
    emit('success')
    isOpen.value = false
    // Reset form
    form.employeeID = ''
    form.periodStart = ''
    form.periodEnd = ''
  } catch (err: any) {
    if (err?.error && Array.isArray(err.error)) {
      setErrors(err.error as ValidationError[])
    } else {
      toast.error(err?.message || 'Failed to create payroll')
    }
  }
}
</script>

<template>
  <Dialog v-model:open="isOpen">
    <DialogContent 
      class="sm:max-w-[425px]"
      @pointer-down-outside="(e) => isDirty && e.preventDefault()"
      @escape-key-down="(e) => isDirty && e.preventDefault()"
    >
      <form @submit.prevent="handleSubmit">
        <DialogHeader>
          <DialogTitle>Create Payroll</DialogTitle>
          <DialogDescription>
            Generate a new payroll record for an employee for a specific period.
          </DialogDescription>
        </DialogHeader>

        <div class="grid gap-4 py-4">
          <div class="grid gap-2">
            <Label for="employeeID">Employee</Label>
            <select
              id="employeeID"
              v-model="form.employeeID"
              class="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
              required
            >
              <option value="" disabled>Select employee</option>
              <option v-for="emp in employeeStore.employees" :key="emp.id" :value="emp.id">
                {{ emp.name }} ({{ emp.code }})
              </option>
            </select>
            <span v-if="hasError('EmployeeID')" class="text-xs text-destructive">
              {{ getErrorMessage('EmployeeID') }}
            </span>
          </div>

          <div class="grid gap-2">
            <Label for="periodStart">Period Start</Label>
            <Input
              id="periodStart"
              v-model="form.periodStart"
              type="date"
              required
            />
            <span v-if="hasError('PeriodStart')" class="text-xs text-destructive">
              {{ getErrorMessage('PeriodStart') }}
            </span>
          </div>

          <div class="grid gap-2">
            <Label for="periodEnd">Period End</Label>
            <Input
              id="periodEnd"
              v-model="form.periodEnd"
              type="date"
              required
            />
            <span v-if="hasError('PeriodEnd')" class="text-xs text-destructive">
              {{ getErrorMessage('PeriodEnd') }}
            </span>
          </div>
        </div>

        <DialogFooter>
          <Button type="button" variant="outline" @click="handleCancel"> Cancel </Button>
          <Button type="submit" :disabled="payrollStore.loading">
            {{ payrollStore.loading ? 'Creating...' : 'Create Payroll' }}
          </Button>
        </DialogFooter>
      </form>
      <CancelModal v-model:open="showCancelModal" @confirm="confirmCancel" />
    </DialogContent>
  </Dialog>
</template>
