<script setup lang="ts">
import LoginForm from '@/components/LoginForm.vue'
</script>

<template>
  <div class="grid min-h-svh lg:grid-cols-2 overflow-hidden bg-background">
    <!-- Left Panel: Dark Mode Minimalist -->
    <div
      class="relative hidden lg:flex flex-col justify-between bg-zinc-950 p-16 overflow-hidden select-none"
    >
      <!-- Background Patterns -->
      <div class="absolute inset-0 z-0">
        <!-- Fine Dot Grid -->
        <div
          class="absolute inset-0 opacity-20"
          style="
            background-image: radial-gradient(circle at 1px 1px, #3f3f46 1px, transparent 0);
            background-size: 32px 32px;
          "
        />

        <!-- Abstract Intersecting Wireframe Patterns -->
        <div
          class="absolute top-1/4 -left-20 w-[600px] h-[600px] border border-white/5 rounded-full"
        />
        <div
          class="absolute -top-32 right-1/4 w-96 h-96 border border-white/5 rotate-12 rounded-2xl"
        />
        <div
          class="absolute bottom-1/4 -right-20 w-[500px] h-[500px] border border-white/5 rounded-full"
        />
        <div
          class="absolute -bottom-20 left-1/3 w-80 h-80 border border-white/5 -rotate-12 rounded-3xl"
        />
      </div>

      <!-- Center Hero Content -->
      <div class="relative z-10 flex flex-col justify-center flex-1 max-w-xl">
        <h1 class="text-6xl font-bold text-white tracking-tight leading-[1.1] mb-6">
          The All-in-One <br />POS System.
        </h1>
        <p class="text-xl text-white font-normal leading-relaxed opacity-80">
          Manage your sales, track inventory, and monitor analytics easily in one place.
        </p>
      </div>

      <!-- Bottom Area: Empty -->
      <div class="relative z-10 h-10" />
    </div>

    <!-- Right Panel: Focused Auth -->
    <div class="flex items-center justify-center p-8 bg-background lg:p-24">
      <div class="w-full max-w-[400px]">
        <LoginForm />
      </div>
    </div>
  </div>
</template>
