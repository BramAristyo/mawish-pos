<script setup lang="ts">
import { Button } from '@/components/ui/button'
import { Label } from '@/components/ui/label'

interface Option {
  label: string
  value: any
}

interface Props {
  modelValue: any
  options: Option[]
  label?: string
  description?: string
  disabled?: boolean
}

defineProps<Props>()

const emit = defineEmits<{
  (e: 'update:modelValue', value: any): void
}>()
</script>

<template>
  <div class="space-y-3">
    <div v-if="label || description" class="space-y-0.5">
      <Label v-if="label" class="text-sm font-medium">{{ label }}</Label>
      <p v-if="description" class="text-sm text-muted-foreground">
        {{ description }}
      </p>
    </div>
    <div class="flex items-center gap-2 p-1 border rounded-lg w-fit bg-muted/50">
      <Button
        v-for="option in options"
        :key="String(option.value)"
        type="button"
        :variant="modelValue === option.value ? 'default' : 'ghost'"
        size="sm"
        class="min-w-24 transition-all"
        :disabled="disabled"
        @click="emit('update:modelValue', option.value)"
      >
        {{ option.label }}
      </Button>
    </div>
  </div>
</template>
