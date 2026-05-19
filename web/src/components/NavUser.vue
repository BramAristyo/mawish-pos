<script setup lang="ts">
import { LogOut, User, Settings, CreditCard, Bell } from 'lucide-vue-next'
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
  DropdownMenuShortcut,
} from '@/components/ui/dropdown-menu'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { useAuthStore } from '@/stores/auth.stores'
import { useRouter } from 'vue-router'
import { Button } from '@/components/ui/button'
import { ref } from 'vue'

const props = defineProps<{
  user: {
    name: string
    email: string
    avatar: string
  }
}>()

const authStore = useAuthStore()
const router = useRouter()
const isLogoutDialogOpen = ref(false)

function handleLogout() {
  authStore.logout()
  router.push('/')
}
</script>

<template>
  <div class="flex items-center gap-2">
    <DropdownMenu>
      <DropdownMenuTrigger as-child>
        <Button
          variant="ghost"
          class="h-10 w-10 rounded-full p-0 hover:bg-transparent focus-visible:ring-1 focus-visible:ring-primary"
        >
          <Avatar
            class="h-10 w-10 border border-border shadow-sm transition-transform active:scale-95"
          >
            <AvatarImage :src="user.avatar" :alt="user.name" />
            <AvatarFallback class="bg-primary text-primary-foreground">
              {{ user.name.charAt(0) }}
            </AvatarFallback>
          </Avatar>
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent
        class="w-64 rounded-xl p-2 shadow-xl border-border/50"
        align="end"
        :side-offset="12"
      >
        <DropdownMenuLabel class="p-3 font-normal">
          <div class="flex items-center gap-3 text-left">
            <Avatar class="h-9 w-9 rounded-lg border">
              <AvatarImage :src="user.avatar" :alt="user.name" />
              <AvatarFallback class="rounded-lg"> {{ user.name.charAt(0) }} </AvatarFallback>
            </Avatar>
            <div class="grid flex-1 text-left leading-tight">
              <span class="truncate font-semibold text-sm">{{ user.name }}</span>
              <span class="truncate text-xs text-muted-foreground">{{ user.email }}</span>
            </div>
          </div>
        </DropdownMenuLabel>
        <DropdownMenuSeparator />
        <div class="p-1">
          <DropdownMenuItem
            class="rounded-lg cursor-pointer py-2 px-3 focus:bg-accent focus:text-accent-foreground"
          >
            <User class="mr-3 h-4 w-4 opacity-70" />
            <span class="text-sm">My Profile</span>
            <DropdownMenuShortcut class="text-[10px] opacity-50">⇧⌘P</DropdownMenuShortcut>
          </DropdownMenuItem>
          <DropdownMenuItem
            class="rounded-lg cursor-pointer py-2 px-3 focus:bg-accent focus:text-accent-foreground"
          >
            <Settings class="mr-3 h-4 w-4 opacity-70" />
            <span class="text-sm">Settings</span>
            <DropdownMenuShortcut class="text-[10px] opacity-50">⌘S</DropdownMenuShortcut>
          </DropdownMenuItem>
        </div>
        <DropdownMenuSeparator />
        <div class="p-1">
          <DropdownMenuItem
            @select.prevent="isLogoutDialogOpen = true"
            class="rounded-lg cursor-pointer py-2 px-3 text-destructive!"
          >
            <LogOut class="mr-3 h-4 w-4 text-current" />
            <span class="text-sm font-semibold">Log out</span>
          </DropdownMenuItem>
        </div>
      </DropdownMenuContent>
    </DropdownMenu>

    <!-- Logout Confirmation Dialog -->
    <Dialog v-model:open="isLogoutDialogOpen">
      <DialogContent class="sm:max-w-md">
        <DialogHeader>
          <DialogTitle>Confirm Logout</DialogTitle>
          <DialogDescription>
            Are you sure you want to log out of your account? Any unsaved changes may be lost.
          </DialogDescription>
        </DialogHeader>
        <DialogFooter class="flex sm:justify-end gap-2 pt-4">
          <Button variant="outline" @click="isLogoutDialogOpen = false"> Cancel </Button>
          <Button variant="destructive" @click="handleLogout"> Log out </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  </div>
</template>
