<script setup lang="ts">
import { computed, reactive, watch, ref } from 'vue'
import { useShiftScheduleStore } from '@/stores/shiftSchedule.store'
import type {
  ShiftSchedule,
  CreateShiftScheduleRequest,
  UpdateShiftScheduleRequest,
} from '@/types/shiftSchedule.types'
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
import { AmountInput } from '@/components/common/form/input/amount'

const props = defineProps<{
  open: boolean
  shiftSchedule?: ShiftSchedule | null
}>()

const emit = defineEmits<{
  (e: 'update:open', value: boolean): void
  (e: 'success'): void
}>()

const shiftScheduleStore = useShiftScheduleStore()
const { setErrors, clearErrors, getErrorMessage, hasError } = useFormErrors()

const isOpen = computed({
  get: () => props.open,
  set: (value) => emit('update:open', value),
})

const isEdit = computed(() => !!props.shiftSchedule)

const form = reactive<CreateShiftScheduleRequest>({
  name: '',
  startTime: '08:00',
  endTime: '17:00',
  toleranceMinutes: 0,
  lateIntervalMinutes: 0,
  lateDeductionAmount: 0,
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
  if (props.shiftSchedule) {
    form.name = props.shiftSchedule.name
    form.startTime = props.shiftSchedule.startTime
    form.endTime = props.shiftSchedule.endTime
    form.toleranceMinutes = props.shiftSchedule.toleranceMinutes
    form.lateIntervalMinutes = props.shiftSchedule.lateIntervalMinutes
    form.lateDeductionAmount = props.shiftSchedule.lateDeductionAmount
  } else {
    form.name = ''
    form.startTime = '08:00'
    form.endTime = '17:00'
    form.toleranceMinutes = 0
    form.lateIntervalMinutes = 0
    form.lateDeductionAmount = 0
  }
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

watch(
  () => props.shiftSchedule,
  () => {
    resetForm()
  },
  { immediate: true },
)

async function handleSubmit() {
  clearErrors()
  try {
    const payload = {
      ...form,
      toleranceMinutes: Number(form.toleranceMinutes),
      lateIntervalMinutes: Number(form.lateIntervalMinutes),
      lateDeductionAmount: Number(form.lateDeductionAmount),
    }

    if (isEdit.value && props.shiftSchedule) {
      await shiftScheduleStore.update(props.shiftSchedule.id, payload as UpdateShiftScheduleRequest)
      toast.success('Shift schedule updated successfully')
    } else {
      await shiftScheduleStore.create(payload as CreateShiftScheduleRequest)
      toast.success('Shift schedule created successfully')
    }
    emit('success')
    isOpen.value = false
  } catch (err: any) {
    if (err?.error && Array.isArray(err.error)) {
      setErrors(err.error as ValidationError[])
    } else {
      toast.error(err?.message || 'Failed to save shift schedule')
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
          <DialogTitle>{{ isEdit ? 'Edit Shift Schedule' : 'Add Shift Schedule' }}</DialogTitle>
          <DialogDescription>
            {{
              isEdit
                ? 'Update the details of your shift schedule here.'
                : 'Create a new shift schedule for your employees.'
            }}
          </DialogDescription>
        </DialogHeader>

        <div class="grid gap-4 py-4">
          <div class="grid gap-2">
            <Label for="name">Name</Label>
            <Input id="name" v-model="form.name" placeholder="Morning Shift" required />
            <span v-if="hasError('Name')" class="text-xs text-destructive">
              {{ getErrorMessage('Name') }}
            </span>
          </div>

          <div class="grid grid-cols-2 gap-4">
            <div class="grid gap-2">
              <Label for="startTime">Start Time</Label>
              <Input id="startTime" v-model="form.startTime" type="time" required />
              <span v-if="hasError('StartTime')" class="text-xs text-destructive">
                {{ getErrorMessage('StartTime') }}
              </span>
            </div>
            <div class="grid gap-2">
              <Label for="endTime">End Time</Label>
              <Input id="endTime" v-model="form.endTime" type="time" required />
              <span v-if="hasError('EndTime')" class="text-xs text-destructive">
                {{ getErrorMessage('EndTime') }}
              </span>
            </div>
          </div>

          <div class="grid gap-2">
            <Label for="toleranceMinutes">Tolerance (Minutes)</Label>
            <Input
              id="toleranceMinutes"
              v-model.number="form.toleranceMinutes"
              type="number"
              placeholder="0"
              required
            />
            <span v-if="hasError('ToleranceMinutes')" class="text-xs text-destructive">
              {{ getErrorMessage('ToleranceMinutes') }}
            </span>
          </div>

          <div class="grid gap-2">
            <Label for="lateIntervalMinutes">Late Interval (Minutes)</Label>
            <Input
              id="lateIntervalMinutes"
              v-model.number="form.lateIntervalMinutes"
              type="number"
              placeholder="0"
              required
            />
            <span v-if="hasError('LateIntervalMinutes')" class="text-xs text-destructive">
              {{ getErrorMessage('LateIntervalMinutes') }}
            </span>
          </div>

          <div class="grid gap-2">
            <Label for="lateDeductionAmount">Late Deduction Amount</Label>
            <AmountInput
              id="lateDeductionAmount"
              v-model="form.lateDeductionAmount"
              placeholder="0"
              required
            />
            <span v-if="hasError('LateDeductionAmount')" class="text-xs text-destructive">
              {{ getErrorMessage('LateDeductionAmount') }}
            </span>
          </div>
        </div>

        <DialogFooter>
          <Button type="button" variant="outline" @click="handleCancel"> Cancel </Button>
          <Button type="submit" :disabled="shiftScheduleStore.loading">
            {{ shiftScheduleStore.loading ? 'Saving...' : 'Save' }}
          </Button>
        </DialogFooter>
      </form>
      <CancelModal v-model:open="showCancelModal" @confirm="confirmCancel" />
    </DialogContent>
  </Dialog>
</template>
