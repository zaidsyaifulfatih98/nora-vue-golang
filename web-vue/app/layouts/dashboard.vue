<template>
  <div class="flex min-h-screen bg-gray-50">
    <SideBar />
    <div class="flex min-w-0 flex-1 flex-col">
      <DashboardHeader />
      <div class="flex-1 overflow-x-hidden p-4 sm:p-6 lg:p-10">
        <slot />
      </div>
    </div>
  </div>
</template>
