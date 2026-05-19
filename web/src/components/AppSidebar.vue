<script setup lang="ts">
import type { SidebarProps } from '@/components/ui/sidebar'
import { ChevronDown, Store } from 'lucide-vue-next'
import SearchForm from '@/components/SearchForm.vue'
import { Collapsible, CollapsibleContent, CollapsibleTrigger } from '@/components/ui/collapsible'
import {
  Sidebar,
  SidebarContent,
  SidebarGroup,
  SidebarHeader,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  SidebarMenuSub,
  SidebarMenuSubButton,
  SidebarMenuSubItem,
  SidebarRail,
} from '@/components/ui/sidebar'
import { computed, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { NAVIGATION_CONFIG } from '@/constant/navigation'

const props = defineProps<SidebarProps>()

const route = useRoute()
const router = useRouter()

const searchQuery = ref('')
const openSectionTitle = ref<string | null>(null)

const filteredNavMain = computed(() => {
  const query = searchQuery.value.toLowerCase().trim()
  if (!query) return NAVIGATION_CONFIG

  return NAVIGATION_CONFIG.map((section) => ({
    ...section,
    items: section.items.filter((item) => item.title.toLowerCase().includes(query)),
  })).filter((section) => section.items.length > 0)
})

function isChildActive(itemUrl: string) {
  if (route.path === itemUrl) return true
  if (itemUrl !== '/' && itemUrl !== '/dashboard' && route.path.startsWith(itemUrl)) {
    return true
  }
  return false
}

function updateOpenSection() {
  const activeSection = NAVIGATION_CONFIG.find((section) =>
    section.items.some((item) => isChildActive(item.url)),
  )
  if (activeSection) {
    openSectionTitle.value = activeSection.title
  }
}

function handleToggle(title: string, open: boolean) {
  if (searchQuery.value) return

  if (open) {
    openSectionTitle.value = title
  } else if (openSectionTitle.value === title) {
    openSectionTitle.value = null
  }
}

const isParentActive = (item: any) => {
  return item.items?.some((child: any) => isChildActive(child.url))
}

function handleParentClick(item: any) {}

watch(
  () => route.path,
  () => {
    updateOpenSection()
  },
  { immediate: true },
)
</script>

<template>
  <Sidebar v-bind="props" class="backdrop-blur-xl">
    <SidebarHeader class="relative overflow-hidden pb-6 pt-8">
      <!-- Subtle Decorative background element -->
      <div class="absolute -right-4 -top-8 size-32 rounded-full blur-3xl" />

      <SidebarMenu class="relative z-10">
        <SidebarMenuItem>
          <SidebarMenuButton size="lg" as-child class="hover:bg-transparent">
            <a href="#" class="flex items-center gap-3">
              <div
                class="flex aspect-square size-10 items-center justify-center rounded-xl bg-primary text-primary-foreground shadow-lg shadow-primary/20 transition-transform"
              >
                <Store class="size-5" />
              </div>
              <div class="flex flex-col gap-0.5 leading-none">
                <span class="text-lg font-bold tracking-tight text-sidebar-foreground"
                  >POS System</span
                >
                <span
                  class="text-[10px] font-medium uppercase tracking-wider text-muted-foreground/60"
                  >Enterprise Edition</span
                >
              </div>
            </a>
          </SidebarMenuButton>
        </SidebarMenuItem>
      </SidebarMenu>
      <div class="px-2 mt-6 relative z-10">
        <SearchForm v-model="searchQuery" />
      </div>
    </SidebarHeader>

    <SidebarContent class="px-2">
      <SidebarGroup>
        <SidebarMenu>
          <Collapsible
            v-for="item in filteredNavMain"
            :key="item.title"
            :open="searchQuery ? true : openSectionTitle === item.title"
            @update:open="(val) => handleToggle(item.title, val)"
            class="group/collapsible mb-2"
          >
            <SidebarMenuItem>
              <CollapsibleTrigger as-child>
                <SidebarMenuButton
                  :is-active="isParentActive(item)"
                  @click="handleParentClick(item)"
                  class="h-11 rounded-xl px-4 transition-all duration-300 hover:bg-sidebar-accent hover:text-sidebar-accent-foreground data-[active=true]:bg-primary data-[active=true]:text-primary-foreground"
                >
                  <component :is="item.icon" v-if="item.icon" class="size-5 opacity-80" />
                  <span
                    :class="[
                      'flex-1 text-sm tracking-tight transition-all',
                      isParentActive(item) ? 'font-bold' : 'font-medium',
                    ]"
                  >
                    {{ item.title }}
                  </span>
                  <ChevronDown
                    v-if="!searchQuery"
                    class="ml-auto size-4 opacity-40 transition-transform duration-300 group-data-[state=open]/collapsible:rotate-180"
                  />
                </SidebarMenuButton>
              </CollapsibleTrigger>
              <CollapsibleContent>
                <SidebarMenuSub class="ml-4 border-l-2 border-primary/10 pl-4 mt-2 space-y-1">
                  <SidebarMenuSubItem v-for="childItem in item.items" :key="childItem.title">
                    <RouterLink :to="childItem.url">
                      <SidebarMenuSubButton
                        :is-active="isChildActive(childItem.url)"
                        class="h-9 rounded-lg"
                      >
                        <span
                          :class="[
                            'text-sm transition-all',
                            isChildActive(childItem.url)
                              ? 'font-semibold'
                              : 'font-medium opacity-70',
                          ]"
                        >
                          {{ childItem.title }}
                        </span>
                      </SidebarMenuSubButton>
                    </RouterLink>
                  </SidebarMenuSubItem>
                </SidebarMenuSub>
              </CollapsibleContent>
            </SidebarMenuItem>
          </Collapsible>
        </SidebarMenu>
      </SidebarGroup>
    </SidebarContent>

    <SidebarRail />
  </Sidebar>
</template>
